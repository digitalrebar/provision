package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/contents"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/digitalrebar/store"
	"github.com/spf13/cobra"
)

type ContentOps struct{ CommonOps }

func (be ContentOps) GetType() interface{} {
	return &models.Content{}
}

func (be ContentOps) GetId(obj interface{}) (string, error) {
	content, ok := obj.(*models.Content)
	if !ok || content.Meta.Name == nil {
		return "", fmt.Errorf("Invalid type passed to content create")
	}
	return *content.Meta.Name, nil
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
		content = &models.Content{Meta: &models.ContentMetaData{Name: &profName}}
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

func findOrFake(field string, args map[string]string) *string {
	buf, err := ioutil.ReadFile(fmt.Sprintf("._%s.meta", field))
	if err == nil {
		s := string(buf)
		return &s
	}
	if p, ok := args[field]; !ok {
		s := "Unspecified"
		return &s
	} else {
		return &p
	}
}

var typeToObject = map[string](func() store.KeySaver){
	"machines":     (&backend.Machine{}).New,
	"params":       (&backend.Param{}).New,
	"profiles":     (&backend.Profile{}).New,
	"users":        (&backend.User{}).New,
	"templates":    (&backend.Template{}).New,
	"bootenvs":     (&backend.BootEnv{}).New,
	"leases":       (&backend.Lease{}).New,
	"reservations": (&backend.Reservation{}).New,
	"subnets":      (&backend.Subnet{}).New,
	"tasks":        (&backend.Task{}).New,
	"jobs":         (&backend.Job{}).New,
	"plugins":      (&backend.Plugin{}).New,
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

	commands = append(commands, &cobra.Command{
		Use:   "bundle [file] [meta fields]",
		Short: "Bundle a directory into a single file, specifed by [file].  [meta fields] allows for the specification of the meta data.",
		Long:  "Bundle assumes that the directories are the object types of the system.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide a file")
			}
			params := map[string]string{}
			for i := 1; i < len(args); i++ {
				if !strings.ContainsAny(args[i], "=") {
					return fmt.Errorf("Meta fields must have '=' in them")
				}
				arrs := strings.SplitN(args[i], "=", 2)
				params[arrs[0]] = arrs[1]
			}
			dumpUsage = false
			filename := args[0]

			content := &models.Content{Meta: &models.ContentMetaData{}}

			content.Meta.Name = findOrFake("Name", params)
			content.Meta.Description = *findOrFake("Description", params)
			content.Meta.Version = *findOrFake("Version", params)
			content.Meta.Source = *findOrFake("Source", params)

			content.Sections = models.Sections{}

			// for each valid content type, load it
			for prefix, fn := range typeToObject {
				objs := map[string]interface{}{}

				err := filepath.Walk(fmt.Sprintf("./%s", prefix), func(filepath string, info os.FileInfo, err error) error {
					if info != nil && !info.IsDir() {
						ext := path.Ext(filepath)
						codec := store.DefaultCodec
						if ext == ".yaml" || ext == ".yml" {
							codec = store.YamlCodec
						}

						obj := fn()

						if buf, err := ioutil.ReadFile(filepath); err != nil {
							return err
						} else {
							if err := codec.Decode(buf, obj); err != nil {
								return err
							}
						}

						objs[obj.Key()] = obj
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("Failed to process content: %v", err)
				}

				if len(objs) > 0 {
					content.Sections[prefix] = objs
				}
			}

			if data, err := prettyPrintBuf(content); err != nil {
				return err
			} else {
				if err := ioutil.WriteFile(filename, data, 0640); err != nil {
					return fmt.Errorf("Failed to write file: %v", err)
				}
			}
			return nil
		},
	})

	res.AddCommand(commands...)
	return res
}
