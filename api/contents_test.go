package api

import (
	"testing"

	"github.com/digitalrebar/provision/v4/models"
)

func TestContentCrud(t *testing.T) {
	tests := []crudTest{
		{
			name:      "Get BarkingStore (that does not exist)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Key:      "BarkingStore",
				Type:     "GET",
				Messages: []string{"No such content store"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return session.GetContentItem("BarkingStore")
			},
		},
		{
			name:      "Delete BarkingStore (that does not exist)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Key:      "BarkingStore",
				Type:     "DELETE",
				Messages: []string{"No such content store"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return nil, session.DeleteContent("BarkingStore")
			},
		},
		{
			name:      "Create Bad BarkingStore (no name)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Type:     "STORE_ERROR",
				Messages: []string{"Store at content- has no Name metadata"},
				Code:     422,
			},
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				return session.CreateContent(barking, false)
			},
		},
		{
			name: "Create BarkingStore",
			expectRes: mustDecode(&models.ContentSummary{}, `
Counts: {}
Warnings: []
meta:
  Description: ""
  Meta: {}
  Name: BarkingStore
  Overwritable: false
  Source: ""
  Type: dynamic
  Version: ""
  Writable: false
`),
			expectErr: nil,
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				barking.Meta.Name = "BarkingStore"
				return session.CreateContent(barking, false)
			},
		},
		{
			name:      "Update BarkingStore (that would break layers)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Key:      "BarkingStore",
				Type:     "PUT",
				Messages: []string{"profiles:global in layer writable would override layer content-BarkingStore"},
				Code:     500,
			},
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				barking.Meta.Name = "BarkingStore"
				env, err := session.GetModel("profiles", "global")
				if err != nil {
					return nil, err
				}
				barking.Sections["profiles"] = map[string]interface{}{env.Key(): env}
				return session.ReplaceContent(barking, false)
			},
		},
		{
			name: "Update BarkingStore",
			expectRes: mustDecode(&models.ContentSummary{}, `
Counts:
  bootenvs: 1
Warnings: []
meta:
  Description: ""
  Meta: {}
  Name: BarkingStore
  Overwritable: false
  Source: ""
  Type: dynamic
  Version: ""
  Writable: false
`),
			expectErr: nil,
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				barking.Meta.Name = "BarkingStore"
				env, err := session.GetModel("bootenvs", "ignore")
				if err != nil {
					return nil, err
				}
				env.(*models.BootEnv).Name = "ignoble"
				barking.Sections["bootenvs"] = map[string]interface{}{env.Key(): env}
				return session.ReplaceContent(barking, false)
			},
		},
		{
			name: "Make sure we can get the ignoble boot env",
			expectRes: mustDecode(&models.BootEnv{}, `
Available: true
Endpoint: ""
Bundle: BarkingStore
Description: The boot environment you should use to have unknown machines boot off
  their local hard drive
Meta:
  color: green
  feature-flags: change-stage-v2
  icon: circle thin
  title: Digital Rebar Provision
Name: ignoble
OS:
  Name: ignore
OnlyUnknown: true
ReadOnly: true
Templates:
- Contents: |
    DEFAULT local
    PROMPT 0
    TIMEOUT 10
    LABEL local
    {{.Param "pxelinux-local-boot"}}
  Name: pxelinux
  Meta: {}
  Path: pxelinux.cfg/default
- Contents: |
    #!ipxe
    chain {{.ProvisionerURL}}/${netX/mac}.ipxe && exit || goto chainip
    :chainip
    chain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit
  Name: ipxe
  Meta: {}
  Path: default.ipxe
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
  Path: grub/grub.cfg
Validated: true
`),
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "ignoble")
			},
		},
		{
			name:      "Delete BarkingStore",
			expectRes: nil,
			expectErr: nil,
			op: func() (interface{}, error) {
				return nil, session.DeleteContent("BarkingStore")
			},
		},
		{
			name:      "Make sure the ignoble boot env is gone",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "bootenvs",
				Key:      "ignoble",
				Type:     "GET",
				Messages: []string{"Not Found"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return session.GetModel("bootenvs", "ignoble")
			},
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}
