package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/digitalrebar/provision/client/templates"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type TemplateOps struct{}

func (be TemplateOps) GetType() interface{} {
	return &models.Template{}
}

func (be TemplateOps) GetId(obj interface{}) (string, error) {
	template, ok := obj.(*models.Template)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to template create")
	}
	return *template.ID, nil
}

func (be TemplateOps) GetIndexes() map[string]string {
	return map[string]string{"ID": "string"}
}

func (be TemplateOps) List(parms map[string]string) (interface{}, error) {
	params := templates.NewListTemplatesParams()
	if listLimit != -1 {
		t1 := int64(listLimit)
		params = params.WithLimit(&t1)
	}
	if listOffset != -1 {
		t1 := int64(listOffset)
		params = params.WithOffset(&t1)
	}
	for k, v := range parms {
		switch k {
		case "ID":
			params = params.WithID(&v)
		}
	}
	d, e := session.Templates.ListTemplates(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Get(id string) (interface{}, error) {
	d, e := session.Templates.GetTemplate(templates.NewGetTemplateParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Create(obj interface{}) (interface{}, error) {
	template, ok := obj.(*models.Template)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to template create")
	}
	d, e := session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(template), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to template patch")
	}
	d, e := session.Templates.PatchTemplate(templates.NewPatchTemplateParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Delete(id string) (interface{}, error) {
	d, e := session.Templates.DeleteTemplate(templates.NewDeleteTemplateParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Upload(id string, f *os.File) (interface{}, error) {
	tmpl := &models.Template{ID: &id}
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, f)
	if err != nil {
		return nil, err
	}
	str := string(buf.Bytes())
	tmpl.Contents = &str

	_, err = be.Get(id)
	if err == nil {
		d, e := session.Templates.PutTemplate(templates.NewPutTemplateParams().WithName(id).WithBody(tmpl), basicAuth)
		if e != nil {
			return nil, e
		}
		return d.Payload, nil
	} else {
		d, e := session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(tmpl), basicAuth)
		if e != nil {
			return nil, e
		}
		return d.Payload, nil
	}
}

func init() {
	tree := addTemplateCommands()
	App.AddCommand(tree)
}

func addTemplateCommands() (res *cobra.Command) {
	singularName := "template"
	name := "templates"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := commonOps(singularName, name, &TemplateOps{})
	res.AddCommand(commands...)
	return res
}
