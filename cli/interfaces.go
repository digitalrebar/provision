package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/client/interfaces"
	"github.com/spf13/cobra"
)

type InterfaceOps struct{ CommonOps }

func (be InterfaceOps) GetIndexes() map[string]string {
	return map[string]string{}
}

func (be InterfaceOps) List(parms map[string]string) (interface{}, error) {
	d, e := session.Interfaces.ListInterfaces(interfaces.NewListInterfacesParams(), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be InterfaceOps) Get(id string) (interface{}, error) {
	d, e := session.Interfaces.GetInterface(interfaces.NewGetInterfaceParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addInterfaceCommands()
	App.AddCommand(tree)
}

func addInterfaceCommands() (res *cobra.Command) {
	singularName := "interface"
	name := "interfaces"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := commonOps(&InterfaceOps{CommonOps{Name: name, SingularName: singularName}})
	res.AddCommand(commands...)
	return res
}
