package provisioner

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/rackn/rocket-skates/models"
	"github.com/rackn/rocket-skates/restapi/operations/templates"
)

// GREG FIX multi-return code test

// Template represents a template that will be associated with a boot environment.
type Template struct {
	models.TemplateOutput
	parsedTmpl *template.Template
}

func CastTemplate(t1 *models.TemplateInput) *Template {
	return &Template{models.TemplateOutput{*t1, make([]string, 0, 0)}, nil}
}

func NewTemplate(uuid string) *Template {
	return &Template{models.TemplateOutput{models.TemplateInput{UUID: uuid}, make([]string, 0, 0)}, nil}
}

func TemplateList(params templates.ListTemplatesParams, p *models.Principal) middleware.Responder {
	allthem, err := listThings(&Template{})
	if err != nil {
		return templates.NewListTemplatesInternalServerError().WithPayload(err)
	}
	data := make([]*models.TemplateOutput, 0, 0)
	for _, j := range allthem {
		original, ok := j.(models.TemplateOutput)
		if ok {
			data = append(data, &original)
		}
	}
	return templates.NewListTemplatesOK().WithPayload(data)
}

func TemplatePost(params templates.PostTemplateParams, p *models.Principal) middleware.Responder {
	item, code, err := createThing(CastTemplate(params.Body))
	if err != nil {
		return templates.NewPostTemplateConflict().WithPayload(err)
	}
	original, ok := item.(models.TemplateOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "failed to cast template")
		return templates.NewPostTemplateInternalServerError().WithPayload(e)
	}
	if code == http.StatusOK {
		return templates.NewPostTemplateOK().WithPayload(&original)
	}
	return templates.NewPostTemplateCreated().WithPayload(&original)
}

func TemplateReplace(params templates.ReplaceTemplateParams, p *models.Principal) middleware.Responder {
	finalStatus := http.StatusCreated
	oldThing := NewTemplate(params.UUID)
	newThing := NewTemplate(params.UUID)
	if err := load(oldThing); err == nil {
		finalStatus = http.StatusOK
	} else {
		oldThing = nil
	}
	buf, err := ioutil.ReadAll(params.Body)
	if err != nil {
		e := NewError(http.StatusExpectationFailed,
			"template: failed to read request body")
		return templates.NewReplaceTemplateExpectationFailed().WithPayload(e)
	}
	newThing.Contents = string(buf)
	if err := save(newThing, oldThing); err != nil {
		e := NewError(http.StatusConflict, err.Error())
		return templates.NewReplaceTemplateConflict().WithPayload(e)
	}

	original, ok := interface{}(newThing).(models.TemplateOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Failed to convert template")
		return templates.NewGetTemplateInternalServerError().WithPayload(e)
	}
	if finalStatus == http.StatusOK {
		return templates.NewReplaceTemplateOK().WithPayload(&original)
	}
	return templates.NewReplaceTemplateCreated().WithPayload(&original)
}

func TemplateGet(params templates.GetTemplateParams, p *models.Principal) middleware.Responder {
	newThing := NewTemplate(params.UUID)
	item, err := getThing(newThing)
	if err != nil {
		return templates.NewGetTemplateNotFound().WithPayload(err)
	}
	original, ok := item.(models.TemplateOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Failed to convert template")
		return templates.NewGetTemplateInternalServerError().WithPayload(e)
	}
	return templates.NewGetTemplateOK().WithPayload(&original)
}

func TemplatePut(params templates.PutTemplateParams, p *models.Principal) middleware.Responder {
	item, err := putThing(CastTemplate(params.Body))
	if err != nil {
		if err.Code == http.StatusNotFound {
			return templates.NewPutTemplateNotFound().WithPayload(err)
		}
		return templates.NewPutTemplateConflict().WithPayload(err)
	}
	original, ok := item.(models.TemplateOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Failed to convert template")
		return templates.NewPutTemplateInternalServerError().WithPayload(e)
	}
	return templates.NewPutTemplateOK().WithPayload(&original)
}

func TemplatePatch(params templates.PatchTemplateParams, p *models.Principal) middleware.Responder {
	newThing := NewTemplate(params.UUID)
	patch, _ := json.Marshal(params.Body)
	item, err := patchThing(newThing, patch)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return templates.NewPatchTemplateNotFound().WithPayload(err)
		}
		if err.Code == http.StatusConflict {
			return templates.NewPatchTemplateConflict().WithPayload(err)
		}
		return templates.NewPatchTemplateExpectationFailed().WithPayload(err)
	}
	original, ok := item.(models.TemplateOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Failed to convert template")
		return templates.NewPatchTemplateInternalServerError().WithPayload(e)
	}
	return templates.NewPatchTemplateOK().WithPayload(&original)
}

func TemplateDelete(params templates.DeleteTemplateParams, p *models.Principal) middleware.Responder {
	err := deleteThing(NewTemplate(params.UUID))
	if err != nil {
		return templates.NewDeleteTemplateConflict().WithPayload(err)
	}
	return templates.NewDeleteTemplateNoContent()
}

func (t *Template) prefix() string {
	return "templates"
}

func (t *Template) key() string {
	return path.Join(t.prefix(), t.UUID)
}

func (t *Template) typeName() string {
	return "TEMPLATE"
}

func (t *Template) newIsh() keySaver {
	res := NewTemplate(t.UUID)
	return keySaver(res)
}

// Parse checks to make sure the template contents are valid according to text/template.
func (t *Template) Parse() (err error) {
	parsedTmpl, err := template.New(t.UUID).Parse(t.Contents)
	if err != nil {
		return err
	}
	t.parsedTmpl = parsedTmpl.Option("missingkey=error")
	return nil
}

func (t *Template) onChange(oldThing interface{}) error {
	if t.Contents == "" || t.UUID == "" {
		return fmt.Errorf("template: Illegal template %+v", t)
	}
	if err := t.Parse(); err != nil {
		return fmt.Errorf("template: %s does not compile: %v", t.UUID, err)
	}

	if old, ok := oldThing.(*Template); ok && old != nil && old.UUID != t.UUID {
		return fmt.Errorf("template: Cannot change UUID of %s", t.UUID)
		machine := &Machine{}
		machines, err := machine.List()
		if err == nil {
			for _, machine := range machines {
				reRender := false
				bootEnv := NewBootenv(machine.BootEnv)
				if err := load(bootEnv); err == nil {
					for ii, template := range bootEnv.Templates {
						ti := bootEnv.TemplateInfo[ii]
						if template.UUID == t.UUID {
							reRender = true
							ti.contents = t
							break
						}
					}
				}
				if reRender {
					bootEnv.RenderTemplates(machine)
				}
			}
		}
	}
	return nil
}

func (t *Template) onDelete() error {
	bootenv := &BootEnv{}
	bootEnvs, err := bootenv.List()
	if err == nil {
		for _, bootEnv := range bootEnvs {
			for _, tmpl := range bootEnv.Templates {
				if tmpl.UUID == t.UUID {
					return fmt.Errorf("template: %s is in use by bootenv %s (template %s", t.UUID, bootEnv.Name, tmpl.Name)
				}
			}
		}
	}
	return err
}

// Render executes the template with params writing the results to dest
func (t *Template) Render(dest io.Writer, params interface{}) error {
	if t.parsedTmpl == nil {
		if err := t.Parse(); err != nil {
			return fmt.Errorf("template: %s does not compile: %v", t.UUID, err)
		}
	}
	if err := t.parsedTmpl.Execute(dest, params); err != nil {
		return fmt.Errorf("template: cannot execute %s: %v", t.UUID, err)
	}
	return nil
}

func (t *Template) RebuildRebarData() error {
	return nil
}
