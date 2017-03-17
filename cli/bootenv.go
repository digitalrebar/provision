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

	"github.com/VictorLowther/jsonpatch"
	bootenvs "github.com/rackn/rocket-skates/client/boot_envs"
	"github.com/rackn/rocket-skates/client/isos"
	"github.com/rackn/rocket-skates/client/templates"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

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
	commands := make([]*cobra.Command, 0, 0)
	commands = append(commands, &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("List all %v", name),
		Run: func(c *cobra.Command, args []string) {
			if resp, err := session.BootEnvs.ListBootEnvs(bootenvs.NewListBootEnvsParams()); err != nil {
				log.Fatalf("Error listing %v: %v", name, err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	/* Match not supported today
	commands = append(commands, &cobra.Command{
		Use:   "match [json]",
		Short: fmt.Sprintf("List all %v that match the template in [json]", name),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			objs := []interface{}{}
			vals := map[string]interface{}{}
			if err := json.Unmarshal([]byte(args[0]), &vals); err != nil {
				log.Fatalf("Matches not valid JSON\n%v", err)
			}
			if err := session.Match(session.UrlPath(maker()), vals, &objs); err != nil {
				log.Fatalf("Error getting matches for %v\nError:%v\n", singularName, err)
			}
			fmt.Println(prettyJSON(objs))
		},
	})
	*/
	commands = append(commands, &cobra.Command{
		Use:   "show [id]",
		Short: fmt.Sprintf("Show a single %v by id", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			if resp, err := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(args[0])); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	/* Sample not supported today
	commands = append(commands, &cobra.Command{
		Use:   "sample",
		Short: fmt.Sprintf("Get the default values for a %v", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 0 {
				log.Fatalf("%v takes no arguments", c.UseLine())
			}
			obj := maker()
			if err := session.Init(obj); err != nil {
				log.Fatalf("Unable to fetch defaults for %v: %v\n", singularName, err)
			}
			fmt.Println(prettyJSON(obj))
		},
	})
	*/
	commands = append(commands, &cobra.Command{
		Use:   "create [json]",
		Short: fmt.Sprintf("Create a new %v with the passed-in JSON", singularName),
		Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			var buf []byte
			var err error
			if args[0] == "-" {
				buf, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Fatalf("Error reading from stdin: %v", err)
				}
			} else {
				buf = []byte(args[0])
			}
			bootenv := &models.BootEnv{}
			err = yaml.Unmarshal(buf, bootenv)
			if err != nil {
				log.Fatalf("Invalid %v object: %v\n", singularName, err)
			}
			if resp, err := session.BootEnvs.CreateBootEnv(bootenvs.NewCreateBootEnvParams().WithBody(bootenv)); err != nil {
				log.Fatalf("Unable to create new %v: %v\n", singularName, err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
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

	commands = append(commands, &cobra.Command{
		Use:   "update [id] [json]",
		Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", singularName),
		Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatalf("%v requires 2 arguments\n", c.UseLine())
			}
			if resp, err := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(args[0])); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
			} else {
				var buf []byte
				var err error
				if args[1] == "-" {
					buf, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						log.Fatalf("Error reading from stdin: %v", err)
					}
				} else {
					buf = []byte(args[1])
				}
				bootenv := resp.Payload
				buf2, err := yaml.Marshal(bootenv)
				if err != nil {
					log.Fatalf("Unable to marshal object: %v\n", err)
				}

				merged, err := safeMergeJSON(buf2, buf)
				if err != nil {
					log.Fatalf("Unable to merge objects: %v\n", err)
				}

				bootenv = &models.BootEnv{}
				err = yaml.Unmarshal(merged, bootenv)
				if err != nil {
					log.Fatalf("Unable to unmarshal merged object: %v\n", err)
				}

				if resp, err := session.BootEnvs.PutBootEnv(bootenvs.NewPutBootEnvParams().WithName(args[0]).WithBody(bootenv)); err != nil {
					log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
				} else {
					fmt.Println(pretty(resp.Payload))
				}
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "patch [objectJson] [changesJson]",
		Short: fmt.Sprintf("Patch %v with the passed-in JSON", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatalf("%v requires 2 arguments\n", c.UseLine())
			}
			obj := &models.BootEnv{}
			if err := yaml.Unmarshal([]byte(args[0]), obj); err != nil {
				log.Fatalf("Unable to parse %v JSON %v\nError: %v\n", c.UseLine(), args[0], err)
			}
			newObj := &models.BootEnv{}
			yaml.Unmarshal([]byte(args[0]), newObj)
			if err := yaml.Unmarshal([]byte(args[1]), newObj); err != nil {
				log.Fatalf("Unable to parse %v JSON %v\nError: %v\n", c.UseLine(), args[1], err)
			}
			newBuf, _ := yaml.Marshal(newObj)
			patch, err := jsonpatch.GenerateJSON([]byte(args[0]), newBuf, true)
			if err != nil {
				log.Fatalf("Cannot generate JSON Patch\n%v\n", err)
			}
			p := []*models.JSONPatchOperation{}
			err = yaml.Unmarshal(patch, p)
			if err != nil {
				log.Fatalf("Cannot generate JSON Patch Object\n%v\n", err)
			}
			if resp, err := session.BootEnvs.PatchBootEnv(bootenvs.NewPatchBootEnvParams().WithName(*obj.Name).WithBody(p)); err != nil {
				log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "destroy [id]",
		Short: fmt.Sprintf("Destroy %v by id", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			if _, err := session.BootEnvs.DeleteBootEnv(bootenvs.NewDeleteBootEnvParams().WithName(args[0])); err != nil {
				log.Fatalf("Unable to destroy %v %v\nError: %v\n", singularName, args[0], err)
			} else {
				fmt.Printf("Deleted %v %v\n", singularName, args[0])
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "exists [id]",
		Short: fmt.Sprintf("See if a %v exists by id", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			if _, err := session.BootEnvs.GetBootEnv(bootenvs.NewGetBootEnvParams().WithName(args[0])); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
			} else {
				os.Exit(0)
			}
		},
	})

	res.AddCommand(commands...)
	return res
}
