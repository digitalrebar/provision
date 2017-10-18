package api

import (
	"log"
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestContentCrud(t *testing.T) {
	summary := `
- Counts:
    bootenvs: 2
    jobs: 0
    leases: 0
    machines: 0
    params: 3
    plugins: 0
    preferences: 0
    profiles: 1
    reservations: 0
    stages: 0
    subnets: 0
    tasks: 0
    templates: 0
    users: 1
  Warnings: []
  meta:
    Description: Writable backing store
    Meta: {}
    Name: BackingStore
    Overwritable: false
    Source: ""
    Type: writable
    Version: user
    Writable: true
- Counts:
    templates: 1
  Warnings: []
  meta:
    Description: Local Override Store
    Meta: {}
    Name: LocalStore
    Overwritable: false
    Source: ""
    Type: local
    Version: user
    Writable: false
- Counts:
    templates: 1
  Warnings: []
  meta:
    Description: Initial Default Content
    Meta: {}
    Name: DefaultStore
    Overwritable: true
    Source: Unspecified
    Type: default
    Version: user
    Writable: false
`
	cs := []models.ContentSummary{}
	if err := DecodeYaml([]byte(summary), &cs); err != nil {
		log.Panicf("Unable to decode reference content summary: %v", err)
	}
	backingStore := `
meta:
  Description: Writable backing store
  Meta: {}
  Name: BackingStore
  Overwritable: false
  Source: Unknown
  Type: writable
  Version: user
  Writable: true
sections:
  bootenvs:
    ignore:
      Available: false
      BootParams: ""
      Description: The boot environment you should use to have unknown machines boot
        off their local hard drive
      Errors: []
      Initrds: []
      Kernel: ""
      Meta: {}
      Name: ignore
      OS:
        Codename: ""
        Family: ""
        IsoFile: ""
        IsoSha256: ""
        IsoUrl: ""
        Name: ignore
        Version: ""
      OnlyUnknown: true
      OptionalParams: []
      ReadOnly: false
      RequiredParams: []
      Templates:
      - Contents: |
          DEFAULT local
          PROMPT 0
          TIMEOUT 10
          LABEL local
          localboot 0
        ID: ""
        Name: pxelinux
        Path: pxelinux.cfg/default
      - Contents: exit
        ID: ""
        Name: elilo
        Path: elilo.conf
      - Contents: |
          #!ipxe
          chain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit
        ID: ""
        Name: ipxe
        Path: default.ipxe
      Validated: false
    local:
      Available: false
      BootParams: ""
      Description: The boot environment you should use to have known machines boot
        off their local hard drive
      Errors: []
      Initrds: []
      Kernel: ""
      Meta: {}
      Name: local
      OS:
        Codename: ""
        Family: ""
        IsoFile: ""
        IsoSha256: ""
        IsoUrl: ""
        Name: local
        Version: ""
      OnlyUnknown: false
      OptionalParams: []
      ReadOnly: false
      RequiredParams: []
      Templates:
      - Contents: |
          DEFAULT local
          PROMPT 0
          TIMEOUT 10
          LABEL local
          localboot 0
        ID: ""
        Name: pxelinux
        Path: pxelinux.cfg/{{.Machine.HexAddress}}
      - Contents: exit
        ID: ""
        Name: elilo
        Path: '{{.Machine.HexAddress}}.conf'
      - Contents: |
          #!ipxe
          exit
        ID: ""
        Name: ipxe
        Path: '{{.Machine.Address}}.ipxe'
      Validated: false
  jobs: {}
  leases: {}
  machines: {}
  params:
    incrementer/parameter:
      Available: false
      Description: ""
      Documentation: ""
      Errors: []
      Meta: {}
      Name: incrementer/parameter
      ReadOnly: false
      Schema:
        type: string
      Validated: false
    incrementer/step:
      Available: false
      Description: ""
      Documentation: ""
      Errors: []
      Meta: {}
      Name: incrementer/step
      ReadOnly: false
      Schema:
        type: integer
      Validated: false
    incrementer/touched:
      Available: false
      Description: ""
      Documentation: ""
      Errors: []
      Meta: {}
      Name: incrementer/touched
      ReadOnly: false
      Schema:
        type: integer
      Validated: false
  plugins: {}
  preferences: {}
  profiles:
    global:
      Available: false
      Description: ""
      Errors: []
      Meta: {}
      Name: global
      Params: {}
      ReadOnly: false
      Validated: false
  reservations: {}
  stages: {}
  subnets: {}
  tasks: {}
  templates: {}
  users:
    rocketskates:
      Available: false
      Errors: []
      Meta: {}
      Name: rocketskates
      PasswordHash: MTYzODQkOCQxJDk0YTBlZDI3N2IxMzNmMGU2NmNjMDdhMzU2ZWNmMzkxJDQ5M2E4OGI0YTdhMTkxN2ZiMDBkNzg2ODk4NjJjYjg0OTgwOWVkODQ1YTc0OGI2YWMyOThjMzkwMjk3Njg4OTQ=
      ReadOnly: false
      Validated: false
`
	bs := &models.Content{}
	if err := DecodeYaml([]byte(backingStore), bs); err != nil {
		log.Panicf("Unable to unmarshal backingStore: %v", err)
	}
	bs.Sections["users"]["rocketskates"].(map[string]interface{})["PasswordHash"] = "elided"
	tests := []crudTest{
		{
			name:      "List all content",
			expectRes: cs,
			expectErr: nil,
			op: func() (interface{}, error) {
				return session.GetContentSummary()
			},
		},
		{
			name:      "Get BackingStore",
			expectRes: bs,
			expectErr: nil,
			op: func() (interface{}, error) {
				res, err := session.GetContentItem("BackingStore")
				if err != nil {
					return res, err
				}
				res.Sections["users"]["rocketskates"].(map[string]interface{})["PasswordHash"] = "elided"
				return res, err
			},
		},
		{
			name:      "Get BarkingStore (that does not exist)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Key:      "BarkingStore",
				Type:     "API_ERROR",
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
				Type:     "API_ERROR",
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
				Key:      "",
				Type:     "STORE_ERROR",
				Messages: []string{"Content stores must have a name"},
				Code:     422,
			},
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				return session.CreateContent(barking)
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
				return session.CreateContent(barking)
			},
		},
		{
			name:      "Create Duplicate BarkingStore",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Key:      "BarkingStore",
				Type:     "API_ERROR",
				Messages: []string{"Content BarkingStore already exists"},
				Code:     409,
			},
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				barking.Meta.Name = "BarkingStore"
				return session.CreateContent(barking)
			},
		},
		{
			name:      "Update BarkingStore (that would break layers)",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "contents",
				Key:      "BarkingStore",
				Type:     "API_ERROR",
				Messages: []string{"New layer violates key restrictions: keysCannotBeOverridden: ignore is already in layer 0\n\tkeysCannotOverride: ignore would be overridden by layer 0"},
				Code:     500,
			},
			op: func() (interface{}, error) {
				barking := &models.Content{}
				barking.Fill()
				barking.Meta.Name = "BarkingStore"
				env, err := session.GetModel("bootenvs", "ignore")
				if err != nil {
					return nil, err
				}
				barking.Sections["bootenvs"] = map[string]interface{}{env.Key(): env}
				return session.ReplaceContent(barking)
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
				return session.ReplaceContent(barking)
			},
		},
		{
			name: "Make sure we can get the ignoble boot env",
			expectRes: mustDecode(&models.BootEnv{}, `
Available: true
Description: The boot environment you should use to have unknown machines boot off
  their local hard drive
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
    localboot 0
  Name: pxelinux
  Path: pxelinux.cfg/default
- Contents: exit
  Name: elilo
  Path: elilo.conf
- Contents: |
    #!ipxe
    chain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit
  Name: ipxe
  Path: default.ipxe
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
				Type:     "API_ERROR",
				Messages: []string{"bootenvs GET: ignoble: Not Found"},
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
