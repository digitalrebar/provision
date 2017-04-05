package cli

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/rackn/rocket-skates/client/machines"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type MachineOps struct{}

func (be MachineOps) GetType() interface{} {
	return &models.Machine{}
}

func (be MachineOps) GetId(obj interface{}) (string, error) {
	machine, ok := obj.(*models.Machine)
	if !ok || machine.UUID == nil {
		return "", fmt.Errorf("Invalid type passed to machine create")
	}
	return machine.UUID.String(), nil
}

func (be MachineOps) List() (interface{}, error) {
	d, e := session.Machines.ListMachines(machines.NewListMachinesParams(), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Get(id string) (interface{}, error) {
	d, e := session.Machines.GetMachine(machines.NewGetMachineParams().WithUUID(strfmt.UUID(id)), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Create(obj interface{}) (interface{}, error) {
	machine, ok := obj.(*models.Machine)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to machine create")
	}
	d, e := session.Machines.CreateMachine(machines.NewCreateMachineParams().WithBody(machine), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to machine patch")
	}
	d, e := session.Machines.PatchMachine(machines.NewPatchMachineParams().WithUUID(strfmt.UUID(id)).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Delete(id string) (interface{}, error) {
	d, e := session.Machines.DeleteMachine(machines.NewDeleteMachineParams().WithUUID(strfmt.UUID(id)), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addMachineCommands()
	App.AddCommand(tree)
}

func addMachineCommands() (res *cobra.Command) {
	singularName := "machine"
	name := "machines"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	commands := commonOps(singularName, name, &MachineOps{})
	res.AddCommand(commands...)
	return res
}
