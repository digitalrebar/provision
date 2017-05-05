package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/client/subnets"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type SubnetOps struct{}

func (be SubnetOps) GetType() interface{} {
	return &models.Subnet{}
}

func (be SubnetOps) GetId(obj interface{}) (string, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to subnet create")
	}
	return *subnet.Name, nil
}

func (be SubnetOps) GetIndexes() map[string]string {
	return map[string]string{"Name": "string", "Subnet": "CIDR IP Address", "Strategy": "string", "NextServer": "IP Address"}
}

func (be SubnetOps) List(parms map[string]string) (interface{}, error) {
	params := subnets.NewListSubnetsParams()
	if listLimit != -1 {
		t1 := int64(listLimit)
		params = params.WithLimit(&t1)
	}
	if listOffset != -1 {
		t1 := int64(listOffset)
		params = params.WithOffset(&t1)
	}
	for k, v := range parms {
		switch k {
		case "Name":
			params = params.WithName(&v)
		case "Subnet":
			params = params.WithSubnet(&v)
		case "Strategy":
			params = params.WithStrategy(&v)
		case "NextServer":
			params = params.WithNextServer(&v)
		}
	}

	d, e := session.Subnets.ListSubnets(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Get(id string) (interface{}, error) {
	d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(id), basicAuth)
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
	d, e := session.Subnets.CreateSubnet(subnets.NewCreateSubnetParams().WithBody(subnet), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet patch")
	}
	d, e := session.Subnets.PatchSubnet(subnets.NewPatchSubnetParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Delete(id string) (interface{}, error) {
	d, e := session.Subnets.DeleteSubnet(subnets.NewDeleteSubnetParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addSubnetCommands()
	App.AddCommand(tree)
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
