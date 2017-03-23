package main

import (
	"fmt"

	"github.com/rackn/rocket-skates/client/subnets"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type SubnetOps struct{}

func (be SubnetOps) GetType() interface{} {
	return &models.Subnet{}
}

func (be SubnetOps) List() (interface{}, error) {
	d, e := session.Subnets.ListSubnets(subnets.NewListSubnetsParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Get(id string) (interface{}, error) {
	d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(id))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Create(obj interface{}) (interface{}, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet create")
	}
	d, e := session.Subnets.CreateSubnet(subnets.NewCreateSubnetParams().WithBody(subnet))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Put(id string, obj interface{}) (interface{}, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet put")
	}
	d, e := session.Subnets.PutSubnet(subnets.NewPutSubnetParams().WithName(id).WithBody(subnet))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.([]*models.JSONPatchOperation)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet patch")
	}
	d, e := session.Subnets.PatchSubnet(subnets.NewPatchSubnetParams().WithName(id).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Delete(id string) (interface{}, error) {
	d, e := session.Subnets.DeleteSubnet(subnets.NewDeleteSubnetParams().WithName(id))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addSubnetCommands()
	app.AddCommand(tree)
}

func addSubnetCommands() (res *cobra.Command) {
	singularName := "subnet"
	name := "subnets"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := commonOps(singularName, name, &SubnetOps{})
	res.AddCommand(commands...)
	return res
}
