package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/digitalrebar/provision/client/contents"
	models "github.com/digitalrebar/provision/genmodels"
	prmodels "github.com/digitalrebar/provision/models"
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
	filepath := fmt.Sprintf("._%s.meta", field)
	fmt.Printf("Processing Meta: %s\n", filepath)
	buf, err := ioutil.ReadFile(filepath)
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

func writeMetaFile(field string, data *string) error {
	if data == nil {
		return nil
	}
	fname := fmt.Sprintf("._%s.meta", field)
	fmt.Printf("Writing Meta: %s\n", fname)
	return ioutil.WriteFile(fname, []byte(*data), 0640)
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
			files, _ := ioutil.ReadDir("./")
			for _, f := range files {
				prefix := f.Name()

				if _, err := prmodels.New(prefix); err != nil {
					// Skip things we can instantiate
					continue
				}
				objs := map[string]interface{}{}

				err := filepath.Walk(fmt.Sprintf("./%s", prefix), func(filepath string, info os.FileInfo, err error) error {
					if info != nil && !info.IsDir() {
						ext := path.Ext(filepath)
						codec := store.DefaultCodec
						if ext == ".yaml" || ext == ".yml" {
							codec = store.YamlCodec
						}

						obj, _ := prmodels.New(prefix)

						fmt.Printf("Processing: %s\n", filepath)
						if buf, err := ioutil.ReadFile(filepath); err != nil {
							return err
						} else {
							if err := codec.Decode(buf, obj); err != nil {
								if prefix == "templates" {
									// Templates could be plain.
									id := path.Base(filepath)
									tmpl := &prmodels.Template{ID: id}
									tmpl.Contents = string(buf)
									obj = tmpl
								} else {
									return err
								}
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

	commands = append(commands, &cobra.Command{
		Use:   "unbundle [file]",
		Short: "Unbundle a [file] into the local directory.",
		Long:  "Unbundle assumes that the current directory is the target",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must provide a file")
			}
			dumpUsage = false
			filename := args[0]

			ext := path.Ext(filename)
			codec := store.DefaultCodec
			if ext == ".yaml" || ext == ".yml" {
				codec = store.YamlCodec
			}

			content := &models.Content{}
			fmt.Printf("Processing: %s\n", filename)
			if buf, err := ioutil.ReadFile(filename); err != nil {
				return err
			} else {
				if err := codec.Decode(buf, content); err != nil {
					return err
				}
			}

			// Record Meta fields
			if err := writeMetaFile("Name", content.Meta.Name); err != nil {
				return err
			}
			s := content.Meta.Source
			if err := writeMetaFile("Source", &s); err != nil {
				return err
			}
			s = content.Meta.Description
			if err := writeMetaFile("Description", &s); err != nil {
				return err
			}
			s = content.Meta.Version
			if err := writeMetaFile("Version", &s); err != nil {
				return err
			}

			// Write sections
			for prefix, data := range content.Sections {
				if err := os.MkdirAll(prefix, 0750); err != nil {
					return err
				}

				for name, obj := range data {
					fname := fmt.Sprintf("%s/%s%s", prefix, name, ext)
					if prefix == "templates" {
						mobj := obj.(map[string]interface{})
						name := mobj["ID"].(string)
						contents := mobj["Contents"].(string)
						fname = fmt.Sprintf("%s/%s", prefix, name)
						fmt.Printf("Writing: %s\n", fname)
						if err := ioutil.WriteFile(fname, []byte(contents), 0640); err != nil {
							return err
						}
					} else {
						if jobj, err := codec.Encode(obj); err != nil {
							return err
						} else {
							fmt.Printf("Writing: %s\n", fname)
							if err := ioutil.WriteFile(fname, jobj, 0640); err != nil {
								return err
							}
						}
					}
				}
			}

			return nil
		},
	})

	res.AddCommand(commands...)
	return res
}
