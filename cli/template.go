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

func (be TemplateOps) GetId(obj interface{}) (string, error) {
	template, ok := obj.(*models.Template)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to template create")
	}
	return *template.ID, nil
}

func (be TemplateOps) List() (interface{}, error) {
	d, e := session.Templates.ListTemplates(templates.NewListTemplatesParams(), basicAuth)
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
	d, e := session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(tmpl), basicAuth)
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
