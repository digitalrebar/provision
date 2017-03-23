package cli

import (
	"os"

	"github.com/rackn/rocket-skates/client/isos"
	"github.com/spf13/cobra"
)

type IsoOps struct{}

func (be IsoOps) List() (interface{}, error) {
	d, e := Session.Isos.ListIsos(isos.NewListIsosParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be IsoOps) Upload(path string, f *os.File) (interface{}, error) {
	d, e := Session.Isos.UploadIso(isos.NewUploadIsoParams().WithPath(path).WithBody(f))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addIsoCommands()
	App.AddCommand(tree)
}

func addIsoCommands() (res *cobra.Command) {
	singularName := "iso"
	name := "isos"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: "Commands to manage isos on the provisioner",
	}
	commands := commonOps(singularName, name, &IsoOps{})
	res.AddCommand(commands...)
	return res
}
