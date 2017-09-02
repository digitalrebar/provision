package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/leases"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type LeaseOps struct{ CommonOps }

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

func (be LeaseOps) GetId(obj interface{}) (string, error) {
	lease, ok := obj.(*models.Lease)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to lease create")
	}
	return lease.Addr.String(), nil
}

func (be LeaseOps) GetIndexes() map[string]string {
	b := &backend.Lease{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be LeaseOps) List(parms map[string]string) (interface{}, error) {
	params := leases.NewListLeasesParams()
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
		case "Available":
			params = params.WithAvailable(&v)
		case "Valid":
			params = params.WithValid(&v)
		case "Addr":
			params = params.WithAddr(&v)
		case "Token":
			params = params.WithToken(&v)
		case "Strategy":
			params = params.WithStrategy(&v)
		case "ExpireTime":
			params = params.WithExpireTime(&v)
		}
	}
	d, e := session.Leases.ListLeases(params, basicAuth)
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
	d, e := session.Leases.GetLease(leases.NewGetLeaseParams().WithAddress(s), basicAuth)
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
	d, e := session.Leases.CreateLease(leases.NewCreateLeaseParams().WithBody(lease), basicAuth)
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
	d, e := session.Leases.PatchLease(leases.NewPatchLeaseParams().WithAddress(s).WithBody(data), basicAuth)
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
	d, e := session.Leases.DeleteLease(leases.NewDeleteLeaseParams().WithAddress(s), basicAuth)
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

	commands := commonOps(&LeaseOps{CommonOps{Name: name, SingularName: singularName}})
	res.AddCommand(commands...)
	return res
}
