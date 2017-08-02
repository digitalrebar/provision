package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ghodss/yaml"

	"github.com/digitalrebar/provision/backend"
	bootenvs "github.com/digitalrebar/provision/client/boot_envs"
	"github.com/digitalrebar/provision/client/isos"
	"github.com/digitalrebar/provision/client/templates"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type BootEnvOps struct{ CommonOps }

func (be BootEnvOps) GetType() interface{} {
	return &models.BootEnv{}
}

func (be BootEnvOps) GetId(obj interface{}) (string, error) {
	bootenv, ok := obj.(*models.BootEnv)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to bootenv create")
	}
	return *bootenv.Name, nil
}

func (be BootEnvOps) GetIndexes() map[string]string {
	b := &backend.BootEnv{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be BootEnvOps) List(parms map[string]string) (interface{}, error) {
	params := bootenvs.NewListBootEnvsParams()
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
		case "Available":
			params = params.WithAvailable(&v)
		case "OnlyUnknown":
			params = params.WithOnlyUnknown(&v)
		case "Name":
			params = params.WithName(&v)
		}
	}

	d, e := session.BootEnvs.ListBootEnvs(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Get(id string) (interface{}, error) {
	d, e := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Create(obj interface{}) (interface{}, error) {
	bootenv, ok := obj.(*models.BootEnv)
	if !ok {
		name, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to bootenv create")
		}
		bootenv = &models.BootEnv{Name: &name}
	}
	d, e := session.BootEnvs.CreateBootEnv(bootenvs.NewCreateBootEnvParams().WithBody(bootenv), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Update(id string, obj interface{}) (interface{}, error) {
	bootenv, ok := obj.(*models.BootEnv)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to bootenv update")
	}
	d, e := session.BootEnvs.PutBootEnv(bootenvs.NewPutBootEnvParams().WithName(id).WithBody(bootenv), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to bootenv patch")
	}
	d, e := session.BootEnvs.PatchBootEnv(bootenvs.NewPatchBootEnvParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Delete(id string) (interface{}, error) {
	d, e := session.BootEnvs.DeleteBootEnv(bootenvs.NewDeleteBootEnvParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addBootEnvCommands()
	App.AddCommand(tree)
}

var installSkipDownloadIsos = true

func uploadTemplateFile(tid string) error {
	_, err := session.Templates.GetTemplate(
		templates.NewGetTemplateParams().WithName(tid), basicAuth)
	if err == nil {
		return nil
	}
	log.Printf("Installing template %s", tid)
	tmpl := &models.Template{}
	tmpl.ID = &tid
	tmplName := path.Join("templates", tid)
	buf, err := ioutil.ReadFile(tmplName)
	if err != nil {
		return generateError(err, "Unable to find template: %s", tid)
	}
	tmplContents := string(buf)
	tmpl.Contents = &tmplContents
	if _, err := session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(tmpl), basicAuth); err != nil {
		return generateError(err, "Unable to create new template: %s", tid)
	}
	return nil
}

func addBootEnvCommands() (res *cobra.Command) {
	singularName := "bootenv"
	name := "bootenvs"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	commands := commonOps(&BootEnvOps{CommonOps{Name: name, SingularName: singularName}})

	installCmd := &cobra.Command{
		Use:   "install [bootenvFile] [isoPath]",
		Short: "Install a bootenv along with everything it requires",
		Long: `bootenvs install assumes you are in a directory with two subdirectories:
   bootenvs/
   templates/

bootenvs must contain [bootenvFile]
templates must contain any templates that the requested bootenv refers to.

bootenvs install will try to upload any required ISOs if they are not already
present in DigitalRebar Provision.  If [isoPath] is specified, it will use that
directory to to check and download ISOs into, otherwise it will use isos/  If the
ISO is not present, we will try to download it if the bootenv specifies a location
to download the ISO from.  If we cannot find an ISO to upload, then the bootenv
will still be uploaded, but it will not be available until the ISO is uploaded
using isos upload.git `,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v needs at least 1 arg", c.UseLine())
			}
			if len(args) > 2 {
				return fmt.Errorf("%v has Too many args", c.UseLine())
			}
			dumpUsage = false
			isoCache := "isos"
			if len(args) == 2 {
				isoCache = args[1]
			}
			if bs, err := os.Stat("bootenvs"); err != nil {
				return fmt.Errorf("Error determining whether bootenvs dir exists: %s", err)
			} else if !bs.IsDir() {
				return fmt.Errorf("bootenvs is not a directory")
			}
			var err error
			var bootEnvBuf []byte
			bootEnvBuf, err = ioutil.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("No bootenv %s", args[0])
			}
			bootEnv := &models.BootEnv{}
			err = yaml.Unmarshal(bootEnvBuf, bootEnv)
			if err != nil {
				return fmt.Errorf("Invalid %v object: %v\n", singularName, err)
			}
			// Upload any required templates if needed.  This includes inline templates
			for _, ti := range bootEnv.Templates {
				if ti.ID == "" {
					continue
				}
				err = uploadTemplateFile(ti.ID)
				if err != nil {
					return err
				}
			}
			// Upload all templates in the templates directory - from subtemplate inclusion
			files, err := ioutil.ReadDir("templates")
			if err == nil {
				for _, f := range files {
					err = uploadTemplateFile(f.Name())
					if err != nil {
						return err
					}
				}
			}

			if err = os.MkdirAll(isoCache, 0755); err != nil {
				return fmt.Errorf("Error ensuring ISO cache exists: %s", err)
			}
			// Upload the bootenv
			log.Printf("Installing bootenv %s", *bootEnv.Name)
			resp, err := session.BootEnvs.CreateBootEnv(bootenvs.NewCreateBootEnvParams().WithBody(bootEnv), basicAuth)
			if err != nil {
				return generateError(err, "Unable to create new %v", singularName)
			}
			if bootEnv.OS.IsoFile == "" {
				return prettyPrint(resp.Payload)
			}
			// See if we need to install the ISO
			isoResp, err := session.Isos.ListIsos(isos.NewListIsosParams(), basicAuth)
			if err != nil {
				return generateError(err, "Error listing isos")
			}
			for _, isoName := range isoResp.Payload {
				if bootEnv.OS.IsoFile == isoName {
					return prettyPrint(resp.Payload)
				}
			}
			// We need to install the ISO
			isoPath := path.Join(isoCache, bootEnv.OS.IsoFile)
			if _, err := os.Stat(isoPath); err != nil {
				isoUrl := bootEnv.OS.IsoURL.String()
				if installSkipDownloadIsos {
					log.Printf("Skipping ISO download as requested")
					log.Printf("Upload with `drpcli isos upload %s as %s` when you have it", bootEnv.OS.IsoFile, bootEnv.OS.IsoFile)
					return prettyPrint(resp.Payload)
				}
				err = func() error {
					// It is not present locally, we need to download it
					if isoUrl == "" {
						return fmt.Errorf("Unable to automatically download %s", isoUrl)
					}
					log.Printf("Downloading %s to %s", isoUrl, isoPath)
					isoTarget, err := os.Create(isoPath)
					defer isoTarget.Close()
					if err != nil {
						return fmt.Errorf("Unable to create %s to download ISO into: %v", isoPath, err)
					}
					isoDlResp, err := http.Get(isoUrl)
					if err != nil {
						return fmt.Errorf("Unable to connect to %s: %v", isoUrl, err)
					}
					defer isoDlResp.Body.Close()
					if isoDlResp.StatusCode >= 300 {
						return fmt.Errorf("Unable to initiate download of %s: %s", isoUrl, isoDlResp.Status)
					}
					byteCount, err := io.Copy(isoTarget, isoDlResp.Body)
					if err != nil {
						return fmt.Errorf("Download of %s aborted: %v", isoUrl, err)
					}
					log.Printf("Downloaded %d bytes", byteCount)
					return nil
				}()
				if err != nil {
					return err
				}
			}
			// We have the ISO now.
			log.Printf("Uploading %s to DigitalRebar Provision", isoPath)
			isoTarget, err := os.Open(isoPath)
			if err != nil {
				return fmt.Errorf("Unable to open %s for upload: %v", isoPath, err)
			}
			defer isoTarget.Close()
			params := isos.NewUploadIsoParams()
			params.Path = bootEnv.OS.IsoFile
			params.Body = isoTarget
			if _, err := session.Isos.UploadIso(params, basicAuth); err != nil {
				return generateError(err, "Error uploading %s", isoPath)
			}
			if resp, err := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(*bootEnv.Name), basicAuth); err != nil {
				return generateError(err, "Failed to fetch %v: %v", singularName, *bootEnv.Name)
			} else {
				return prettyPrint(resp.Payload)
			}
		},
	}
	installCmd.Flags().BoolVar(&installSkipDownloadIsos, "skip-download", false, "Whether to try to download ISOs from their upstream")
	commands = append(commands, installCmd)

	res.AddCommand(commands...)
	return res
}
