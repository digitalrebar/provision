package api

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/digitalrebar/provision/models"
)

func TestBootEnvCrud(t *testing.T) {
	localBootEnv := mustDecode(&models.BootEnv{}, `
Description: "The boot environment you should use to have known machines boot off their local hard drive"
Name: "local"
Meta:
  color: green
  feature-flags: change-stage-v2
  icon: radio
  title: Digital Rebar Provision
OS:
  Name: "local"
ReadOnly: true
Templates:
  - Contents: |
      DEFAULT local
      PROMPT 0
      TIMEOUT 10
      LABEL local
      {{.Param "pxelinux-local-boot"}}
    ID: ""
    Meta: {}
    Name: pxelinux
    Path: pxelinux.cfg/{{.Machine.HexAddress}}
  - Contents: |
      #!ipxe
      exit
    ID: ""
    Meta: {}
    Name: ipxe
    Path: '{{.Machine.Address}}.ipxe'
  - Contents: |
      DEFAULT local
      PROMPT 0
      TIMEOUT 10
      LABEL local
      {{.Param "pxelinux-local-boot"}}
    ID: ""
    Meta: {}
    Name: pxelinux-mac
    Path: pxelinux.cfg/{{.Machine.MacAddr "pxelinux"}}
  - Contents: |
      #!ipxe
      exit
    ID: ""
    Meta: {}
    Name: ipxe-mac
    Path: '{{.Machine.MacAddr "ipxe"}}.ipxe'
`).(*models.BootEnv)
	ignoreBootEnv := mustDecode(&models.BootEnv{}, `
Description: "The boot environment you should use to have unknown machines boot off their local hard drive"
Name:        "ignore"
Meta:
  color: green
  feature-flags: change-stage-v2
  icon: circle thin
  title: Digital Rebar Provision
OS:
  Name: "ignore"
OnlyUnknown: true
ReadOnly: true
Templates:
- Contents: |
    DEFAULT local
    PROMPT 0
    TIMEOUT 10
    LABEL local
    {{.Param "pxelinux-local-boot"}}
  Meta: {}
  Name: "pxelinux"
  Path: "pxelinux.cfg/default"
- Contents: |
    #!ipxe
    chain {{.ProvisionerURL}}/${netX/mac}.ipxe && exit || goto chainip
    :chainip
    chain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit
  Meta: {}
  Name: "ipxe"
  Path: "default.ipxe"`).(*models.BootEnv)
	fred := &models.BootEnv{Name: "fred"}
	fred.SetValid()
	fred.SetAvailable()
	testFill(fred)

	phred := models.Clone(fred)
	phred.(*models.BootEnv).OS.Name = "phred"
	tests := []crudTest{
		{
			name:      "List all feebles",
			expectRes: nil,
			expectErr: errors.New("No such Model: feebles"),
			op:        func() (interface{}, error) { return session.ListModel("feebles") },
		},
		{
			name:      "List all bootenvs",
			expectRes: []models.Model{ignoreBootEnv, localBootEnv},
			expectErr: nil,
			op:        func() (interface{}, error) { return session.ListModel("bootenvs") },
		},
		{
			name:      "List all bootenvs in reverse order",
			expectRes: []models.Model{localBootEnv, ignoreBootEnv},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "reverse", "true")
			},
		},
		{
			name:      "List all bootenvs by OnlyUnknown",
			expectRes: []models.Model{localBootEnv, ignoreBootEnv},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "sort", "OnlyUnknown")
			},
		},
		{
			name:      "List all bootenvs by OnlyUnknown in reverse",
			expectRes: []models.Model{ignoreBootEnv, localBootEnv},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "sort", "OnlyUnknown", "reverse", "true")
			},
		},
		{
			name:      "List just the local bootenv",
			expectRes: []models.Model{localBootEnv},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "Name", "local")
			},
		},
		{
			name:      "List the first bootenv",
			expectRes: []models.Model{ignoreBootEnv},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "limit", "1")
			},
		},
		{
			name:      "List the second bootenv",
			expectRes: []models.Model{localBootEnv},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "limit", "1", "offset", "1")
			},
		},
		{
			name:      "List no bootenvs",
			expectRes: []models.Model{},
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "limit", "0")
			},
		},
		{
			name:      "List a negative number of bootenvs",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Type:     "GET",
				Code:     406,
				Messages: []string{"Limit cannot be negative"},
			},
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "limit", "-1")
			},
		},
		{
			name:      "List with a negative offset",
			expectRes: []models.Model{localBootEnv},
			expectErr: &models.Error{
				Model:    "bootenvs",
				Type:     "GET",
				Code:     406,
				Messages: []string{"Offset cannot be negative"},
			},
			op: func() (interface{}, error) {
				return session.ListModel("bootenvs", "offset", "-1")
			},
		},
		{
			name:      "Test to see if the dweezil bootenv exists",
			expectRes: false,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ExistsModel("bootenvs", "dweezil")
			},
		},
		{
			name:      "Test to see if the local bootenv exists",
			expectRes: true,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.ExistsModel("bootenvs", "local")
			},
		},
		{
			name:      "Get the robert feeble",
			expectRes: nil,
			expectErr: errors.New("No such Model: feebles"),
			op: func() (interface{}, error) {
				return session.GetModel("feebles", "robert")
			},
		},
		{
			name:      "Get the local bootenv",
			expectRes: localBootEnv,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "local")
			},
		},
		{
			name:      "Get the ignore bootenv",
			expectRes: ignoreBootEnv,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "ignore")
			},
		},
		{
			name:      "Get the frabjulous bootenv",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "frabjulous",
				Type:     "GET",
				Code:     404,
				Messages: []string{"Not Found"},
			},
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "frabjulous")
			},
		},
		{
			name:      "Get the local bootenv by name",
			expectRes: localBootEnv,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "Name:local")
			},
		},
		{
			name:      "Get the ignore bootenv by name",
			expectRes: ignoreBootEnv,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "Name:ignore")
			},
		},
		{
			name:      "Get the ignore bootenv by OnlyUnknown",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "OnlyUnknown:true",
				Type:     "GET",
				Messages: []string{"Not Found"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "OnlyUnknown:true")
			},
		},
		{
			name:      "Delete fred bootenv (409)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "fred",
				Type:     "DELETE",
				Code:     404,
				Messages: []string{"Not Found"},
			},
			op: func() (interface{}, error) {
				return session.DeleteModel("bootenvs", "fred")
			},
		},
		{
			name:      "Create a fred bootenv",
			expectRes: fred,
			expectErr: nil,
			op: func() (interface{}, error) {
				m := &models.BootEnv{Name: "fred"}
				return m, session.CreateModel(m)
			},
		},
		{
			name:      "Create another fred bootenv",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "fred",
				Type:     "CREATE",
				Messages: []string{"already exists"},
				Code:     409,
			},
			op: func() (interface{}, error) {
				m := &models.BootEnv{Name: "fred"}
				return m, session.CreateModel(m)
			},
		},
		{
			name:      "List all bootenvs (with fred)",
			expectRes: []models.Model{fred, ignoreBootEnv, localBootEnv},
			expectErr: nil,
			op:        func() (interface{}, error) { return session.ListModel("bootenvs") },
		},
		{
			name:      "PUT Update bootenv fred OS name ->phred",
			expectRes: phred,
			expectErr: nil,
			op: func() (interface{}, error) {
				m := models.Clone(fred)
				m.(*models.BootEnv).OS.Name = "phred"
				return m, session.PutModel(m)
			},
		},
		{
			name:      "PATCH Update bootenv phred OS name (success)",
			expectRes: phred,
			expectErr: nil,
			op: func() (interface{}, error) {
				m := models.Clone(phred)
				m.(*models.BootEnv).OS.Name = "ffred"
				patch, _ := GenPatch(phred, m, false)
				phred.(*models.BootEnv).OS.Name = "ffred"
				return session.PatchModel("bootenvs", "fred", patch)
			},
		},
		{
			name:      "PATCH Update bootenv phred OS name (conflict)",
			expectRes: nil,
			expectErr: &models.Error{
				Model: "bootenvs",
				Key:   "fred",
				Type:  "PATCH",
				Messages: []string{
					"Patch error at line 0: Test op failed.",
					"Patch line: {\"op\":\"test\",\"path\":\"/OS/Name\",\"from\":\"\",\"value\":\"ffred\"}",
				},
				Code: 409,
			},
			op: func() (interface{}, error) {
				m := models.Clone(phred)
				m.(*models.BootEnv).OS.Name = "zfred"
				session.PutModel(m)
				m.(*models.BootEnv).OS.Name = "qfred"
				patch, _ := GenPatch(phred, m, false)
				phred.(*models.BootEnv).OS.Name = "qfred"
				return session.PatchModel("bootenvs", "fred", patch)
			},
		},
		{
			name:      "Fill fred with phred",
			expectRes: phred,
			expectErr: nil,
			op: func() (interface{}, error) {
				phred.(*models.BootEnv).OS.Name = "zfred"
				return fred, session.FillModel(fred, "fred")
			},
		},
		{
			name:      "Delete fred bootenv (success)",
			expectRes: fred,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.DeleteModel("bootenvs", "fred")
			},
		},
		{
			name:      "Delete fred bootenv (fail)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "fred",
				Type:     "DELETE",
				Code:     404,
				Messages: []string{"Not Found"},
			},
			op: func() (interface{}, error) {
				return session.DeleteModel("bootenvs", "fred")
			},
		},
	}
	for _, test := range tests {
		test.run(t)
	}
}

func TestBootEnvImport(t *testing.T) {
	fredstage := mustDecode(&models.Stage{}, `
Name: fred
BootEnv: fredhammer
`).(*models.Stage)
	fredhammer := mustDecode(&models.BootEnv{}, `
Available: true
BootParams: Acounted for
Errors:
- Fake error
Initrds:
- stage1.img
Kernel: vmlinuz0
Name: fredhammer
OS:
  IsoFile: sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar
  IsoUrl: http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar
  Name: sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b
OptionalParams:
- ntp_servers
- access_keys
Templates:
- Contents: 'Attention all '
  Meta: {}
  Name: pxelinux
  Path: pxelinux.cfg/{{.Machine.HexAddress}}
- Contents: planets of the
  Meta: {}
  Name: elilo
  Path: '{{.Machine.HexAddress}}.conf'
- Contents: Solar Federation
  Meta: {}
  Name: ipxe
  Path: '{{.Machine.Address}}.ipxe'
- Contents: We have assumed control
  Meta: {}
  Name: control.sh
  Path: '{{.Machine.Path}}/control.sh'
Validated: true
`).(*models.BootEnv)
	tests := []crudTest{
		{
			name: "Create fred stage",
			expectRes: mustDecode(&models.Stage{}, `
Available: false
BootEnv: fredhammer
Description: ""
Errors:
- BootEnv fredhammer does not exist
Meta: {}
Name: fred
OptionalParams: []
Profiles: []
ReadOnly: false
Reboot: false
RequiredParams: []
RunnerWait: true
Tasks: []
Templates: []
Validated: true`),
			expectErr: nil,
			op: func() (interface{}, error) {
				return fredstage, session.CreateModel(fredstage)
			},
		},
		{
			name:      "Test bootenv install from nonsenical location",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "",
				Type:     "CLIENT_ERROR",
				Messages: []string{"stat does/not/exist: no such file or directory"},
			},
			op: func() (interface{}, error) {
				return session.InstallBootEnvFromFile("does/not/exist")
			},
		},
		{
			name:      "Test install from a known to exist directory",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "",
				Type:     "CLIENT_ERROR",
				Messages: []string{". is a directory.  It needs to be a file."},
			},
			op: func() (interface{}, error) {
				return session.InstallBootEnvFromFile(".")
			},
		},
		{
			name:      "Test install from a valid path that contains invalid content",
			expectRes: nil,
			expectErr: &models.Error{
				Model: "bootenvs",
				Key:   "",
				Type:  "CLIENT_ERROR",
				Messages: []string{
					"error converting YAML to JSON: yaml: line 1: did not find expected node content",
				},
				Code: 0,
			},
			op: func() (interface{}, error) {
				return session.InstallBootEnvFromFile("test-data/badhammer.yml")
			},
		},
		{
			name:      "Test install from a valid flat install (no ISO)",
			expectRes: fredhammer,
			expectErr: nil,
			op: func() (interface{}, error) {
				res, err := session.InstallBootEnvFromFile("test-data/fredhammer.yml")
				if err != nil {
					return res, err
				}
				res.Errors = []string{"Fake error"}
				return res, err
			},
		},
		{
			name:      "Test ISO upload for valid install (manual with invalid path)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "isos",
				Key:      "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
				Type:     "DOWNLOAD_NOT_ALLOWED",
				Messages: []string{"Iso not present at server, not present locally, and automatic download forbidden"},
			},
			op: func() (interface{}, error) {
				env, err := session.GetModel("bootenvs", "fredhammer")
				if err != nil {
					return env, err
				}
				return env, session.InstallISOForBootenv(env.(*models.BootEnv), "/no/iso/here", false)
			},
		},
		{
			name:      "Test ISO upload for valid install (automatic with invalid path and bad source)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "isos",
				Key:      "http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
				Type:     "DOWNLOAD_FAILED",
				Messages: []string{"open /no/iso/here: no such file or directory"},
			},
			op: func() (interface{}, error) {
				env, err := session.GetModel("bootenvs", "fredhammer")
				if err != nil {
					return env, err
				}
				return env, session.InstallISOForBootenv(env.(*models.BootEnv), "/no/iso/here", true)
			},
		},
		{
			name:      "Test ISO upload for valid install (valid path, invalid source)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "isos",
				Key:      "http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
				Type:     "DOWNLOAD_FAILED",
				Messages: []string{"Unable to start download of http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar: 404 Not Found"},
			},
			op: func() (interface{}, error) {
				return nil, session.InstallISOForBootenv(fredhammer, path.Join(tmpDir, "fredhammer.iso"), true)
			},
		},
		{
			name:      "Test ISO upload for valid install (valid path, valid source)",
			expectRes: fredhammer,
			expectErr: nil,
			op: func() (interface{}, error) {
				dst, err := os.Create(path.Join(tmpDir, "tftpboot", "files", fredhammer.OS.IsoFile))
				if err != nil {
					return nil, err
				}
				src, err := os.Open(path.Join("../cli/test-data/", fredhammer.OS.IsoFile))
				if err != nil {
					dst.Close()
					return nil, err
				}
				_, err = io.Copy(dst, src)
				dst.Close()
				src.Close()
				if err != nil {
					return nil, err
				}
				err = session.InstallISOForBootenv(fredhammer, path.Join(tmpDir, "fredhammer.iso"), true)
				if err != nil {
					return nil, err
				}
				fredhammer.Available = true
				fredhammer.Errors = []string{}
				time.Sleep(15 * time.Second)
				return session.GetModel("bootenvs", fredhammer.Key())
			},
		},
		{
			name: "Test to make sure fred stage is available",
			expectRes: mustDecode(&models.Stage{}, `
Available: true
BootEnv: fredhammer
Description: ""
Errors: []
Meta: {}
Name: fred
OptionalParams: []
Profiles: []
ReadOnly: false
Reboot: false
RequiredParams: []
RunnerWait: true
Tasks: []
Templates: []
Validated: true`),
			expectErr: nil,
			op: func() (interface{}, error) {
				return fredstage, session.Req().Fill(fredstage)
			},
		},
		{
			name:      "Clean up after fredhammer (flat install)",
			expectRes: nil,
			expectErr: nil,
			op: func() (interface{}, error) {
				st, err := session.DeleteModel("stages", "fred")
				if err != nil {
					return st, err
				}
				env, err := session.DeleteModel("bootenvs", "fredhammer")
				if err != nil {
					return env, err
				}
				if err := session.DeleteBlob("isos", env.(*models.BootEnv).OS.IsoFile); err != nil {
					return env, err
				}
				if err := session.DeleteBlob("files", env.(*models.BootEnv).OS.IsoFile); err != nil {
					return env, err
				}
				return nil, nil
			},
		},
		{
			name: "Install local3 bootenv (flat install)",
			expectRes: mustDecode(&models.BootEnv{}, `
Available: true
Name: local3
OS:
  Name: local3
Templates:
- ID: local3-pxelinux.tmpl
  Meta: {}
  Name: pxelinux
  Path: pxelinux.cfg/{{.Machine.HexAddress}}
- ID: local3-elilo.tmpl
  Meta: {}
  Name: elilo
  Path: '{{.Machine.HexAddress}}.conf'
- ID: local3-ipxe.tmpl
  Meta: {}
  Name: ipxe
  Path: '{{.Machine.Address}}.ipxe'
Validated: true
`),
			expectErr: nil,
			op: func() (interface{}, error) {
				res, err := session.InstallBootEnvFromFile("../cli/test-data/local3.yml")
				if err == nil {
					_, err = session.DeleteModel("bootenvs", "local3")
				}
				if err == nil {
					_, err = session.DeleteModel("templates", "local3-pxelinux.tmpl")
				}
				if err == nil {
					_, err = session.DeleteModel("templates", "local3-elilo.tmpl")
				}
				if err == nil {
					_, err = session.DeleteModel("templates", "local3-ipxe.tmpl")
				}
				return res, err
			},
		},
		{
			name:      "Install local3 bootenv (/bootenvs without /templates)",
			expectRes: nil,
			expectErr: &models.Error{
				Model: "bootenvs",
				Key:   "local3",
				Type:  "CLIENT_ERROR",
				Messages: []string{
					"Unable to import template local3-pxelinux.tmpl",
					"Unable to import template local3-elilo.tmpl",
					"Unable to import template local3-ipxe.tmpl",
				},
				Code: 0,
			},
			op: func() (interface{}, error) {
				benv := path.Join(tmpDir, "bootenvs")
				tgt := path.Join(benv, "local3.yml")
				if err := os.MkdirAll(benv, 0755); err != nil {
					return nil, err
				}
				src, err := filepath.Abs("../cli/test-data/local3.yml")
				if err != nil {
					return nil, err
				}
				os.Symlink(src, tgt)
				res, err := session.InstallBootEnvFromFile(tgt)
				if err == nil {
					_, err = session.DeleteModel("bootenvs", "local3")
				}
				return res, err
			},
		},
		{
			name: "Install local3 bootenv (/bootenvs and /templates)",
			expectRes: mustDecode(&models.BootEnv{}, `
Available: true
Name: local3
OS:
  Name: local3
Templates:
- ID: local3-pxelinux.tmpl
  Meta: {}
  Name: pxelinux
  Path: pxelinux.cfg/{{.Machine.HexAddress}}
- ID: local3-elilo.tmpl
  Meta: {}
  Name: elilo
  Path: '{{.Machine.HexAddress}}.conf'
- ID: local3-ipxe.tmpl
  Meta: {}
  Name: ipxe
  Path: '{{.Machine.Address}}.ipxe'
Validated: true
`),
			expectErr: nil,
			op: func() (interface{}, error) {
				tmplts := path.Join(tmpDir, "templates")
				if err := os.MkdirAll(tmplts, 0755); err != nil {
					return nil, err
				}
				for _, name := range []string{"local3-pxelinux.tmpl", "local3-elilo.tmpl", "local3-ipxe.tmpl"} {
					tgt := path.Join(tmplts, name)
					src, err := filepath.Abs(path.Join("../cli/test-data", name))
					if err != nil {
						return nil, err
					}
					if err := os.Symlink(src, tgt); err != nil {
						return nil, err
					}
				}
				res, err := session.InstallBootEnvFromFile(path.Join(tmpDir, "bootenvs", "local3.yml"))
				if err == nil {
					_, err = session.DeleteModel("bootenvs", "local3")
				}
				if err == nil {
					_, err = session.DeleteModel("templates", "local3-pxelinux.tmpl")
				}
				if err == nil {
					_, err = session.DeleteModel("templates", "local3-elilo.tmpl")
				}
				if err == nil {
					_, err = session.DeleteModel("templates", "local3-ipxe.tmpl")
				}
				return res, err
			},
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}
