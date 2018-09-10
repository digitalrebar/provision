package cli

import (
	"testing"
)

func TestMultiArch(t *testing.T) {
	cliTest(false, false, "machines", "create", "-").Stdin(`{"Name":"amd64","Arch":"amd64"}`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`{"Name":"arm64","Arch":"arm64"}`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`{"Name":"arm","Arch":"arm"}`).run(t)
	cliTest(false, false, "bootenvs", "create", "-").Stdin(`---
Name: march-discover
OS:
  Name: march-discover
  SupportedArchitectures:
    x86_64:
      BootParams: "I am amd64, AKA x86_64"
      Kernel: vmlinuz0
      IsoFile: march-amd64.tar
    aarch64:
      BootParams: "I am aarch64, AKA arm64"
      Kernel: vmlinuz0
      IsoFile: march-arm64.tar
OnlyUnknown: true
Templates:
  - Contents: |
      chain {{.ProvisionerURL}}/${netX/mac}.ipxe && exit || goto chainip
      :chainip
      chain {{.ProvisionerURL}}/${netX/ip}.ipxe && exit || goto sledgehammer
      :sledgehammer
      kernel {{.Env.PathFor "http" .Env.Kernel}} {{.BootParams}} BOOTIF=01-${netX/mac:hexhyp}
    Name: ipxe
    Path: default.ipxe
`).run(t)
	cliTest(false, false, "prefs", "set", "unknownBootEnv", "march-discover").run(t)
	cliTest(false, false, "files", "static", "default.ipxe").run(t)
	cliTest(false, false, "bootenvs", "create", "-").Stdin(`---
Name: march-install
OS:
  Name: march
  SupportedArchitectures:
    x86_64:
      BootParams: "I am amx64, AKA x86_64"
      Kernel: vmlinuz0
      IsoFile: march-amd64.tar
    aarch64:
      BootParams: "I am arm64, AKA aarch64"
      Kernel: vmlinuz0
      IsoFile: march-arm64.tar
Templates:
- Contents: |
    {{.Env.PathFor "tftp" .Env.Kernel }}
    {{.BootParams}}
    {{.Env.InstallUrl}}
  Name: ipxe
  Path: '{{.Machine.Name}}/{{.Machine.Arch}}/kernel'
`).run(t)
	cliTest(false, false, "isos", "upload", "test-data/march-amd64.tar", "as", "march-amd64.tar").run(t)
	cliTest(false, false, "files", "static", "march/install").run(t)
	cliTest(false, false, "machines", "bootenv", "Name:amd64", "march-install").run(t)
	cliTest(false, false, "isos", "upload", "test-data/march-arm64.tar", "as", "march-arm64.tar").run(t)
	cliTest(false, false, "files", "static", "march/arm64/install").run(t)
	cliTest(false, false, "machines", "bootenv", "Name:arm64", "march-install").run(t)
	cliTest(false, true, "machines", "bootenv", "Name:arm", "march-install").run(t)
	cliTest(false, false, "machines", "destroy", "Name:arm").run(t)
	cliTest(false, false, "files", "static", "amd64/amd64/kernel").run(t)
	cliTest(false, false, "files", "static", "arm64/arm64/kernel").run(t)
	cliTest(false, false, "files", "static", "march/install/vmlinuz0").run(t)
	cliTest(false, false, "files", "static", "march/arm64/install/vmlinuz0").run(t)
	cliTest(false, false, "machines", "destroy", "Name:arm64").run(t)
	cliTest(false, false, "machines", "destroy", "Name:amd64").run(t)
	cliTest(false, false, "isos", "destroy", "march-amd64.tar").run(t)
	cliTest(false, false, "isos", "destroy", "march-arm64.tar").run(t)
	cliTest(false, false, "bootenvs", "destroy", "march-install").run(t)
	verifyClean(t)
}
