package api

import (
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
Bundle: BasicStore
Name: "local"
Meta:
  color: green
  feature-flags: change-stage-v2
  icon: radio
  title: Digital Rebar Provision
OS:
  Name: "local"
Endpoint: ""
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
  - Contents: |
      if test $grub_platform == pc; then
          chainloader (hd0)
      else
          bpx=/efi/boot
          root='' prefix=''
          search --file --set=root $bpx/bootx64.efi || search --file --set=root $bpx/bootaa64.efi
          if test x$root == x; then
              echo "No EFI boot partiton found."
              echo "Rebooting in 120 seconds"
              sleep 120
              reboot
          fi
          if test -f ($root)/efi/microsoft/boot/bootmgfw.efi; then
              echo "Microsoft Windows found, chainloading into it"
              chainloader ($root)/efi/microsoft/boot/bootmgfw.efi
          fi
          for f in ($root)/efi/*; do
              if test -f $f/grub.cfg; then
                  prefix=$f
                  break
              fi
          done
          if test x$prefix == x; then
              echo "Unable to find grub.cfg"
              echo "Rebooting in 120 seconds"
              sleep 120
              reboot
          fi
          configfile $prefix/grub.cfg
      fi
    ID: ""
    Meta: {}
    Name: grub
    Path: grub/{{.Machine.Address}}.cfg
  - Contents: |
      if test $grub_platform == pc; then
          chainloader (hd0)
      else
          bpx=/efi/boot
          root='' prefix=''
          search --file --set=root $bpx/bootx64.efi || search --file --set=root $bpx/bootaa64.efi
          if test x$root == x; then
              echo "No EFI boot partiton found."
              echo "Rebooting in 120 seconds"
              sleep 120
              reboot
          fi
          if test -f ($root)/efi/microsoft/boot/bootmgfw.efi; then
              echo "Microsoft Windows found, chainloading into it"
              chainloader ($root)/efi/microsoft/boot/bootmgfw.efi
          fi
          for f in ($root)/efi/*; do
              if test -f $f/grub.cfg; then
                  prefix=$f
                  break
              fi
          done
          if test x$prefix == x; then
              echo "Unable to find grub.cfg"
              echo "Rebooting in 120 seconds"
              sleep 120
              reboot
          fi
          configfile $prefix/grub.cfg
      fi
    ID: ""
    Meta: {}
    Name: grub-mac
    Path: grub/{{.Machine.MacAddr "grub"}}.cfg
`).(*models.BootEnv)
	ignoreBootEnv := mustDecode(&models.BootEnv{}, `
Description: "The boot environment you should use to have unknown machines boot off their local hard drive"
Bundle: BasicStore
Name:        "ignore"
Meta:
  color: green
  feature-flags: change-stage-v2
  icon: circle thin
  title: Digital Rebar Provision
OS:
  Name: "ignore"
OnlyUnknown: true
Endpoint: ""
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
  Path: "default.ipxe"
- Contents: |
    set _kernel=linux
    set _module=initrd
    $_kernel
    if test $? != 18; then
        set _kernel=linuxefi
        set _module=initrdefi
    fi
    function kernel { $_kernel "$@"; }
    function module { $_module "$@"; }
    if test -s (tftp)/grub/${net_default_mac}.cfg; then
        echo "Booting via MAC"
        source (tftp)/grub/${net_default_mac}.cfg
        boot
    elif test -s (tftp)/grub/${net_default_ip}.cfg; then
        echo "Booting via IP"
        source (tftp)/grub/${net_default_ip}.cfg
        boot
    elif test $grub_platform == pc; then
        chainloader (hd0)
    else
        bpx=/efi/boot
        root='' prefix=''
        search --file --set=root $bpx/bootx64.efi || search --file --set=root $bpx/bootaa64.efi
        if test x$root == x; then
            echo "No EFI boot partiton found."
            echo "Rebooting in 120 seconds"
            sleep 120
            reboot
        fi
        if test -f ($root)/efi/microsoft/boot/bootmgfw.efi; then
            echo "Microsoft Windows found, chainloading into it"
            chainloader ($root)/efi/microsoft/boot/bootmgfw.efi
        fi
        for f in ($root)/efi/*; do
            if test -f $f/grub.cfg; then
                prefix=$f
                break
            fi
        done
        if test x$prefix == x; then
            echo "Unable to find grub.cfg"
            echo "Rebooting in 120 seconds"
            sleep 120
            reboot
        fi
        configfile $prefix/grub.cfg
    fi
  ID: ""
  Meta: {}
  Name: grub
  Path: grub/grub.cfg`).(*models.BootEnv)
	fred := &models.BootEnv{Name: "fred"}
	fred.SetValid()
	fred.SetAvailable()
	testFill(fred)

	phred := models.Clone(fred)
	phred.(*models.BootEnv).OS.Name = "phred"

	rt(t,
		"List all bootenvs",
		[]models.Model{ignoreBootEnv, localBootEnv},
		nil,
		func() (interface{}, error) { return session.ListModel("bootenvs") },
		nil)
	rt(t,
		"List all bootenvs in reverse order",
		[]models.Model{localBootEnv, ignoreBootEnv},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "reverse", "true")
		},
		nil)
	rt(t,
		"List all bootenvs by OnlyUnknown",
		[]models.Model{ignoreBootEnv},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "OnlyUnknown", "true")
		},
		nil)
	rt(t,
		"List all bootenvs by OnlyUnknown in reverse",
		[]models.Model{localBootEnv},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "OnlyUnknown", "false")
		},
		nil)
	rt(t,
		"List just the local bootenv",
		[]models.Model{localBootEnv},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "Name", "local")
		},
		nil)
	rt(t,
		"List the first bootenv",
		[]models.Model{ignoreBootEnv},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "limit", "1")
		},
		nil)
	rt(t,
		"List the second bootenv",
		[]models.Model{localBootEnv},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "limit", "1", "offset", "1")
		},
		nil)
	rt(t,
		"List no bootenvs",
		[]models.Model{},
		nil,
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "limit", "0")
		},
		nil)
	rt(t,
		"List a negative number of bootenvs",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Type:     "GET",
			Code:     406,
			Messages: []string{"Limit cannot be negative"},
		},
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "limit", "-1")
		},
		nil)
	rt(t,
		"List with a negative offset",
		[]models.Model{localBootEnv},
		&models.Error{
			Model:    "bootenvs",
			Type:     "GET",
			Code:     406,
			Messages: []string{"Offset cannot be negative"},
		},
		func() (interface{}, error) {
			return session.ListModel("bootenvs", "offset", "-1")
		},
		nil)
	rt(t,
		"Test to see if the dweezil bootenv exists",
		false,
		nil,
		func() (interface{}, error) {
			return session.ExistsModel("bootenvs", "dweezil")
		},
		nil)
	rt(t,
		"Test to see if the local bootenv exists",
		true,
		nil,
		func() (interface{}, error) {
			return session.ExistsModel("bootenvs", "local")
		},
		nil)
	rt(t,
		"Get the local bootenv",
		localBootEnv,
		nil,
		func() (interface{}, error) {
			return session.GetModel("bootenvs", "local")
		},
		nil)
	rt(t,
		"Get the ignore bootenv",
		ignoreBootEnv,
		nil,
		func() (interface{}, error) {
			return session.GetModel("bootenvs", "ignore")
		},
		nil)
	rt(t,
		"Get the frabjulous bootenv",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "frabjulous",
			Type:     "GET",
			Code:     404,
			Messages: []string{"Not Found"},
		},
		func() (interface{}, error) {
			return session.GetModel("bootenvs", "frabjulous")
		},
		nil)
	rt(t,
		"Get the local bootenv by name",
		localBootEnv,
		nil,
		func() (interface{}, error) {
			return session.GetModel("bootenvs", "Name:local")
		},
		nil)
	rt(t,
		"Get the ignore bootenv by name",
		ignoreBootEnv,
		nil,
		func() (interface{}, error) {
			return session.GetModel("bootenvs", "Name:ignore")
		},
		nil)
	rt(t,
		"Get the ignore bootenv by OnlyUnknown",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "OnlyUnknown:true",
			Type:     "GET",
			Messages: []string{"Not Found"},
			Code:     404,
		},
		func() (interface{}, error) {
			return session.GetModel("bootenvs", "OnlyUnknown:true")
		},
		nil)
	rt(t,
		"Delete fred bootenv (409)",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "fred",
			Type:     "DELETE",
			Code:     404,
			Messages: []string{"Not Found"},
		},
		func() (interface{}, error) {
			return session.DeleteModel("bootenvs", "fred")
		},
		nil)
	rt(t,
		"Create a fred bootenv",
		fred,
		nil,
		func() (interface{}, error) {
			m := &models.BootEnv{Name: "fred"}
			return m, session.CreateModel(m)
		},
		nil)
	rt(t,
		"Create another fred bootenv",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "fred",
			Type:     "CREATE",
			Messages: []string{"already exists"},
			Code:     409,
		},
		func() (interface{}, error) {
			m := &models.BootEnv{Name: "fred"}
			return m, session.CreateModel(m)
		},
		nil)
	rt(t,
		"List all bootenvs (with fred)",
		[]models.Model{fred, ignoreBootEnv, localBootEnv},
		nil,
		func() (interface{}, error) { return session.ListModel("bootenvs") },
		nil)
	rt(t,
		"PUT Update bootenv fred OS name ->phred",
		phred,
		nil,
		func() (interface{}, error) {
			m := models.Clone(fred)
			m.(*models.BootEnv).OS.Name = "phred"
			return m, session.PutModel(m)
		},
		nil)
	rt(t,
		"PATCH Update bootenv phred OS name (success)",
		phred,
		nil,
		func() (interface{}, error) {
			m := models.Clone(phred)
			m.(*models.BootEnv).OS.Name = "ffred"
			patch, _ := GenPatch(phred, m, false)
			phred.(*models.BootEnv).OS.Name = "ffred"
			return session.PatchModel("bootenvs", "fred", patch)
		},
		nil)
	rt(t,
		"PATCH Update bootenv phred OS name (conflict)",
		nil,
		&models.Error{
			Model: "bootenvs",
			Key:   "fred",
			Type:  "PATCH",
			Messages: []string{
				"Patch error at line 0: Test op failed.",
				"Patch line: {\"op\":\"test\",\"path\":\"/OS/Name\",\"from\":\"\",\"value\":\"ffred\"}",
			},
			Code: 409,
		},
		func() (interface{}, error) {
			m := models.Clone(phred)
			m.(*models.BootEnv).OS.Name = "zfred"
			session.PutModel(m)
			m.(*models.BootEnv).OS.Name = "qfred"
			patch, _ := GenPatch(phred, m, false)
			phred.(*models.BootEnv).OS.Name = "qfred"
			return session.PatchModel("bootenvs", "fred", patch)
		},
		nil)
	rt(t,
		"Fill fred with phred",
		phred,
		nil,
		func() (interface{}, error) {
			phred.(*models.BootEnv).OS.Name = "zfred"
			return fred, session.FillModel(fred, "fred")
		},
		nil)
	rt(t,
		"Delete fred bootenv (success)",
		fred,
		nil,
		func() (interface{}, error) {
			return session.DeleteModel("bootenvs", "fred")
		},
		nil)
	rt(t,
		"Delete fred bootenv (fail)",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "fred",
			Type:     "DELETE",
			Code:     404,
			Messages: []string{"Not Found"},
		},
		func() (interface{}, error) {
			return session.DeleteModel("bootenvs", "fred")
		},
		nil)

}

func TestBootEnvImport(t *testing.T) {
	fredstage := mustDecode(&models.Stage{}, `
Name: fred
BootEnv: fredhammer
`).(*models.Stage)
	fredhammer := mustDecode(&models.BootEnv{}, `
Available: true
BootParams: Acounted for
Endpoint: ""
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
	rt(t,
		"Create fred stage",
		mustDecode(&models.Stage{}, `
Available: false
BootEnv: fredhammer
Description: ""
Endpoint: ""
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
		nil,
		func() (interface{}, error) {
			return fredstage, session.CreateModel(fredstage)
		},
		nil)
	rt(t,
		"Test bootenv install from nonsenical location",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "",
			Type:     "CLIENT_ERROR",
			Messages: []string{"stat does/not/exist: no such file or directory"},
		},
		func() (interface{}, error) {
			return session.InstallBootEnvFromFile("does/not/exist")
		},
		nil)
	rt(t,
		"Test install from a known to exist directory",
		nil,
		&models.Error{
			Model:    "bootenvs",
			Key:      "",
			Type:     "CLIENT_ERROR",
			Messages: []string{". is a directory.  It needs to be a file."},
		},
		func() (interface{}, error) {
			return session.InstallBootEnvFromFile(".")
		},
		nil)
	rt(t,
		"Test install from a valid path that contains invalid content",
		nil,
		&models.Error{
			Model: "bootenvs",
			Key:   "",
			Type:  "CLIENT_ERROR",
			Messages: []string{
				"error converting YAML to JSON: yaml: line 1: did not find expected node content",
			},
			Code: 0,
		},
		func() (interface{}, error) {
			return session.InstallBootEnvFromFile("test-data/badhammer.yml")
		},
		nil)
	rt(t,
		"Test install from a valid flat install (no ISO)",
		fredhammer,
		nil,
		func() (interface{}, error) {
			res, err := session.InstallBootEnvFromFile("test-data/fredhammer.yml")
			if err != nil {
				return res, err
			}
			res.Errors = []string{"Fake error"}
			return res, err
		},
		nil)
	rt(t,
		"Test ISO upload for valid install (manual with invalid path)",
		nil,
		&models.Error{
			Model:    "isos",
			Key:      "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
			Type:     "DOWNLOAD_NOT_ALLOWED",
			Messages: []string{"Iso not present at server, not present locally, and automatic download forbidden"},
		},
		func() (interface{}, error) {
			env, err := session.GetModel("bootenvs", "fredhammer")
			if err != nil {
				return env, err
			}
			return env, session.InstallISOForBootenv(env.(*models.BootEnv), "/no/iso/here", false)
		},
		nil)
	rt(t,
		"Test ISO upload for valid install (automatic with invalid path and bad source)",
		nil,
		&models.Error{
			Model:    "isos",
			Key:      "http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
			Type:     "DOWNLOAD_FAILED",
			Messages: []string{"open /no/iso/here/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar: no such file or directory"},
		},
		func() (interface{}, error) {
			env, err := session.GetModel("bootenvs", "fredhammer")
			if err != nil {
				return env, err
			}
			return env, session.InstallISOForBootenv(env.(*models.BootEnv), "/no/iso/here", true)
		},
		nil)
	rt(t,
		"Test ISO upload for valid install (valid path, invalid source)",
		nil,
		&models.Error{
			Model:    "isos",
			Key:      "http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
			Type:     "DOWNLOAD_FAILED",
			Messages: []string{"Unable to start download of http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar: 404 Not Found"},
		},
		func() (interface{}, error) {
			return nil, session.InstallISOForBootenv(fredhammer, tmpDir, true)
		},
		nil)
	rt(t,
		"Test ISO upload for valid install (valid path, valid source)",
		fredhammer,
		nil,
		func() (interface{}, error) {
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
			err = session.InstallISOForBootenv(fredhammer, tmpDir, true)
			if err != nil {
				return nil, err
			}
			fredhammer.Available = true
			fredhammer.Errors = []string{}
			time.Sleep(15 * time.Second)
			return session.GetModel("bootenvs", fredhammer.Key())
		},
		nil)
	rt(t,
		"Test to make sure fredhammer bootenv is available",
		mustDecode(&models.BootEnv{}, `---
Available: true
BootParams: Acounted for
Description: ""
Documentation: ""
Endpoint: ""
Errors: []
Initrds:
- stage1.img
Kernel: vmlinuz0
Meta: {}
Name: fredhammer
OS:
  Codename: ""
  Family: ""
  IsoFile: sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar
  IsoSha256: ""
  IsoUrl: http://127.0.0.1:10012/files/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar
  Name: sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b
  SupportedArchitectures: {}
  Version: ""
OnlyUnknown: false
OptionalParams:
- ntp_servers
- access_keys
ReadOnly: false
RequiredParams: []
Templates:
- Contents: 'Attention all '
  ID: ""
  Meta: {}
  Name: pxelinux
  Path: pxelinux.cfg/{{.Machine.HexAddress}}
- Contents: planets of the
  ID: ""
  Meta: {}
  Name: elilo
  Path: '{{.Machine.HexAddress}}.conf'
- Contents: Solar Federation
  ID: ""
  Meta: {}
  Name: ipxe
  Path: '{{.Machine.Address}}.ipxe'
- Contents: We have assumed control
  ID: ""
  Meta: {}
  Name: control.sh
  Path: '{{.Machine.Path}}/control.sh'
  Validated: true
`),
		nil,
		func() (interface{}, error) {
			err := session.Req().UrlForM(fredhammer).Do(fredhammer)
			return fredhammer, err
		},
		nil)
	rt(t,
		"Test to make sure fred stage is available",
		mustDecode(&models.Stage{}, `
Available: true
BootEnv: fredhammer
Description: ""
Endpoint: ""
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
		nil,
		func() (interface{}, error) {
			err := session.Req().UrlForM(fredstage).Do(fredstage)
			return fredstage, err
		},
		nil)
	rt(t,
		"Clean up after fredhammer (flat install)",
		nil,
		nil,
		func() (interface{}, error) {
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
		nil)
	rt(t,
		"Install local3 bootenv (flat install)",
		mustDecode(&models.BootEnv{}, `
Available: true
Endpoint: ""
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
		nil,
		func() (interface{}, error) {
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
		nil)
	rt(t,
		"Install local3 bootenv (/bootenvs without /templates)",
		nil,
		&models.Error{
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
		func() (interface{}, error) {
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
		nil)
	rt(t,
		"Install local3 bootenv (/bootenvs and /templates)",
		mustDecode(&models.BootEnv{}, `
Available: true
Name: local3
Endpoint: ""
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
		nil,
		func() (interface{}, error) {
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
		nil)
}
