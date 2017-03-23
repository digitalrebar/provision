package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/rackn/rocket-skates/client/templates"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type TemplateOps struct{}

func (be TemplateOps) GetType() interface{} {
	return &models.Template{}
}

func (be TemplateOps) List() (interface{}, error) {
	d, e := Session.Templates.ListTemplates(templates.NewListTemplatesParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Get(id string) (interface{}, error) {
	d, e := Session.Templates.GetTemplate(templates.NewGetTemplateParams().WithName(id))
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
	d, e := Session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(template))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Put(id string, obj interface{}) (interface{}, error) {
	template, ok := obj.(*models.Template)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to template put")
	}
	d, e := Session.Templates.PutTemplate(templates.NewPutTemplateParams().WithName(id).WithBody(template))
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
	d, e := Session.Templates.PatchTemplate(templates.NewPatchTemplateParams().WithName(id).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TemplateOps) Delete(id string) (interface{}, error) {
	d, e := Session.Templates.DeleteTemplate(templates.NewDeleteTemplateParams().WithName(id))
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
	d, e := Session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(tmpl))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
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
