package cli

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/rackn/rocket-skates/client/leases"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type LeaseOps struct{}

func convertStringToAddress(id string) (strfmt.IPv4, error) {
	var s strfmt.IPv4
	err := s.Scan(id)
	if err != nil {
		return "", fmt.Errorf("%v is not a valid IPv4: %v", id, err)
	}
	return s, nil
}

func (be LeaseOps) GetType() interface{} {
	return &models.Lease{}
}

func (be LeaseOps) List() (interface{}, error) {
	d, e := Session.Leases.ListLeases(leases.NewListLeasesParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be LeaseOps) Get(id string) (interface{}, error) {
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := Session.Leases.GetLease(leases.NewGetLeaseParams().WithAddress(s))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be LeaseOps) Create(obj interface{}) (interface{}, error) {
	lease, ok := obj.(*models.Lease)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to lease create")
	}
	d, e := Session.Leases.CreateLease(leases.NewCreateLeaseParams().WithBody(lease))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be LeaseOps) Put(id string, obj interface{}) (interface{}, error) {
	lease, ok := obj.(*models.Lease)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to lease put")
	}
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := Session.Leases.PutLease(leases.NewPutLeaseParams().WithAddress(s).WithBody(lease))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be LeaseOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to lease patch")
	}
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := Session.Leases.PatchLease(leases.NewPatchLeaseParams().WithAddress(s).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be LeaseOps) Delete(id string) (interface{}, error) {
	s, e := convertStringToAddress(id)
	if e != nil {
		return nil, e
	}
	d, e := Session.Leases.DeleteLease(leases.NewDeleteLeaseParams().WithAddress(s))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addLeaseCommands()
	App.AddCommand(tree)
}

func addLeaseCommands() (res *cobra.Command) {
	singularName := "lease"
	name := "leases"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	commands := commonOps(singularName, name, &LeaseOps{})
	res.AddCommand(commands...)
	return res
}
