package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ghodss/yaml"

	bootenvs "github.com/rackn/rocket-skates/client/boot_envs"
	"github.com/rackn/rocket-skates/client/isos"
	"github.com/rackn/rocket-skates/client/templates"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type BootEnvOps struct{}

func (be BootEnvOps) GetType() interface{} {
	return &models.BootEnv{}
}

func (be BootEnvOps) List() (interface{}, error) {
	d, e := session.BootEnvs.ListBootEnvs(bootenvs.NewListBootEnvsParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Get(id string) (interface{}, error) {
	d, e := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(id))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Create(obj interface{}) (interface{}, error) {
	bootenv, ok := obj.(*models.BootEnv)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to bootenv create")
	}
	d, e := session.BootEnvs.CreateBootEnv(bootenvs.NewCreateBootEnvParams().WithBody(bootenv))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Put(id string, obj interface{}) (interface{}, error) {
	bootenv, ok := obj.(*models.BootEnv)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to bootenv put")
	}
	d, e := session.BootEnvs.PutBootEnv(bootenvs.NewPutBootEnvParams().WithName(id).WithBody(bootenv))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.([]*models.JSONPatchOperation)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to bootenv patch")
	}
	d, e := session.BootEnvs.PatchBootEnv(bootenvs.NewPatchBootEnvParams().WithName(id).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be BootEnvOps) Delete(id string) (interface{}, error) {
	d, e := session.BootEnvs.DeleteBootEnv(bootenvs.NewDeleteBootEnvParams().WithName(id))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addBootEnvCommands()
	app.AddCommand(tree)
}

func addBootEnvCommands() (res *cobra.Command) {
	singularName := "bootenv"
	name := "bootenvs"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	commands := commonOps(singularName, name, &BootEnvOps{})

	installDownloadIsos := true
	installCmd := &cobra.Command{
		Use:   "install [bootenvFile] [isoPath]",
		Short: "Install a bootenv along with everything it requires",
		Long: `bootenvs install assumes you are in a directory with two subdirectories:
   bootenvs/
   templates/

bootenvs must contain [bootenvFile]
templates must contain any templates that the requested bootenv refers to.

bootenvs install will try to upload any required ISOs if they are not already
present in RocketSkates.  If [isoPath] is specified, it will use that directory
to to check and download ISOs into, otherwise it will use isos/  If the ISO
is not present, we will try to download it if the bootenv specifies a location
to download the ISO from.  If we cannot find an ISO to upload, then the bootenv
will still be uploaded, but it will not be available until the ISO is uploaded
using isos upload.git `,
		Run: func(c *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatalf("bootenvs install needs at least 1 arg")
			}
			if len(args) > 2 {
				log.Fatalf("Too many args to bootenvs install")
			}
			isoCache := "isos"
			if len(args) == 2 {
				isoCache = args[1]
			}
			var err error
			if err = os.MkdirAll(isoCache, 0755); err != nil {
				log.Fatalf("Error ensuring ISO cache exists: %s", err)
			}
			if bs, err := os.Stat("bootenvs"); err != nil {
				log.Fatalf("Error determining whether bootenvs dir exists: %s", err)
			} else if !bs.IsDir() {
				log.Fatalf("bootenvs is not a directory")
			}
			var bootEnvBuf []byte
			bootEnvBuf, err = ioutil.ReadFile(args[0])
			if err != nil {
				log.Fatalf("No bootenv %s", args[0])
			}
			bootEnv := &models.BootEnv{}
			err = yaml.Unmarshal(bootEnvBuf, bootEnv)
			if err != nil {
				log.Fatalf("Invalid %v object: %v\n", singularName, err)
			}
			// Upload any required templates if needed.
			for _, ti := range bootEnv.Templates {
				if ti.ID == "" {
					continue
				}
				_, err = session.Templates.GetTemplate(
					templates.NewGetTemplateParams().WithName(ti.ID))
				if err == nil {
					continue
				}
				log.Printf("Installing template %s", ti.ID)
				tmpl := &models.Template{}
				tmpl.ID = &ti.ID
				tmplName := path.Join("templates", ti.ID)
				buf, err := ioutil.ReadFile(tmplName)
				if err != nil {
					log.Fatalf("%s requires template %s, but we cannot find it in %s", *bootEnv.Name, ti.ID, tmplName)
				}
				tmplContents := string(buf)
				tmpl.Contents = &tmplContents
				if _, err := session.Templates.CreateTemplate(templates.NewCreateTemplateParams().WithBody(tmpl)); err != nil {
					log.Fatalf("Unable to create new template: %v\n", err)
				}
			}
			// Upload the bootenv
			log.Printf("Installing bootenv %s", *bootEnv.Name)
			resp, err := session.BootEnvs.CreateBootEnv(bootenvs.NewCreateBootEnvParams().WithBody(bootEnv))
			if err != nil {
				log.Fatalf("Unable to create new %v: %v\n", singularName, err)
			}
			if bootEnv.OS.IsoFile == "" {
				fmt.Println(pretty(resp.Payload))
				return
			}
			// See if we need to install the ISO
			isoResp, err := session.Isos.ListIsos(isos.NewListIsosParams())
			if err != nil {
				log.Fatalf("Error listing isos: %v", err)
			}
			for _, isoName := range isoResp.Payload {
				if bootEnv.OS.IsoFile == isoName {
					fmt.Println(pretty(resp.Payload))
					return
				}
			}
			// We need to install the ISO
			isoPath := path.Join(isoCache, bootEnv.OS.IsoFile)
			if _, err := os.Stat(isoPath); err != nil {
				isoUrl := bootEnv.OS.IsoURL.String()
				if !installDownloadIsos {
					log.Printf("Skipping ISO download as requested")
					log.Printf("Upload with `rscli isos upload %s as %s` when you have it", bootEnv.OS.IsoFile, bootEnv.OS.IsoFile)
					fmt.Println(pretty(resp.Payload))
					return
				}
				func() {
					// It is not present locally, we need to download it
					if isoUrl == "" {
						log.Fatalf("Unable to automatically download %s", isoUrl)
					}
					log.Printf("Downloading %s to %s", isoUrl, isoPath)
					isoTarget, err := os.Create(isoPath)
					defer isoTarget.Close()
					if err != nil {
						log.Fatalf("Unable to create %s to download ISO into: %v", isoPath, err)
					}
					isoDlResp, err := http.Get(isoUrl)
					if err != nil {
						log.Fatalf("Unable to connect to %s: %v", isoUrl, err)
					}
					defer isoDlResp.Body.Close()
					if isoDlResp.StatusCode >= 300 {
						log.Fatalf("Unable to initiate download of %s: %s", isoUrl, isoDlResp.Status)
					}
					byteCount, err := io.Copy(isoTarget, isoDlResp.Body)
					if err != nil {
						log.Fatalf("Download of %s aborted: %v", isoUrl, err)
					}
					log.Printf("Downloaded %d bytes", byteCount)
				}()
			}
			// We have the ISO now.
			log.Printf("Uploading %s to RocketSkates", isoPath)
			isoTarget, err := os.Open(isoPath)
			if err != nil {
				log.Fatalf("Unable to open %s for upload: %v", isoPath, err)
			}
			defer isoTarget.Close()
			params := isos.NewUploadIsoParams()
			params.Path = bootEnv.OS.IsoFile
			params.Body = isoTarget
			if _, err := session.Isos.UploadIso(params); err != nil {
				log.Fatalf("Error uploading %s: %v", isoPath, err)
			}
			if resp, err := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(*bootEnv.Name)); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, *bootEnv.Name, err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	}
	installCmd.Flags().BoolVar(&installDownloadIsos, "download", true, "Whether to try to download ISOs from their upstream")
	commands = append(commands, installCmd)

	res.AddCommand(commands...)
	return res
}
