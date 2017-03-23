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

func (be MachineOps) List() (interface{}, error) {
	d, e := Session.Machines.ListMachines(machines.NewListMachinesParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Get(id string) (interface{}, error) {
	d, e := Session.Machines.GetMachine(machines.NewGetMachineParams().WithUUID(strfmt.UUID(id)))
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
	d, e := Session.Machines.CreateMachine(machines.NewCreateMachineParams().WithBody(machine))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Put(id string, obj interface{}) (interface{}, error) {
	machine, ok := obj.(*models.Machine)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to machine put")
	}
	d, e := Session.Machines.PutMachine(machines.NewPutMachineParams().WithUUID(strfmt.UUID(id)).WithBody(machine))
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
	d, e := Session.Machines.PatchMachine(machines.NewPatchMachineParams().WithUUID(strfmt.UUID(id)).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Delete(id string) (interface{}, error) {
	d, e := Session.Machines.DeleteMachine(machines.NewDeleteMachineParams().WithUUID(strfmt.UUID(id)))
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
