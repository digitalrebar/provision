package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/client/contents"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/spf13/cobra"
)

type ContentOps struct{ CommonOps }

func (be ContentOps) GetType() interface{} {
	return &models.Content{}
}

func (be ContentOps) GetId(obj interface{}) (string, error) {
	content, ok := obj.(*models.Content)
	if !ok || content.Name == nil {
		return "", fmt.Errorf("Invalid type passed to content create")
	}
	return *content.Name, nil
}

func (be ContentOps) GetIndexes() map[string]string {
	return map[string]string{}
}

func (be ContentOps) List(parms map[string]string) (interface{}, error) {
	if len(parms) > 0 {
		return nil, fmt.Errorf("Does not support filtering")
	}
	params := contents.NewListContentsParams()
	d, e := session.Contents.ListContents(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ContentOps) Get(id string) (interface{}, error) {
	d, e := session.Contents.GetContent(contents.NewGetContentParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ContentOps) Create(obj interface{}) (interface{}, error) {
	content, ok := obj.(*models.Content)
	if !ok {
		profName, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to content create")
		}
		content = &models.Content{Name: &profName}
	}
	d, e := session.Contents.CreateContent(contents.NewCreateContentParams().WithBody(content), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ContentOps) Update(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(*models.Content)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to content update")
	}
	d, e := session.Contents.UploadContent(contents.NewUploadContentParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ContentOps) Delete(id string) (interface{}, error) {
	_, e := session.Contents.DeleteContent(contents.NewDeleteContentParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return "Good", nil
}

func init() {
	tree := addContentCommands()
	App.AddCommand(tree)
}

func addContentCommands() (res *cobra.Command) {
	singularName := "content"
	name := "contents"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &ContentOps{CommonOps{Name: name, SingularName: singularName}}
	commands := commonOps(mo)
	res.AddCommand(commands...)
	return res
}
