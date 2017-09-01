package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/params"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/spf13/cobra"
)

type ParamOps struct{ CommonOps }

func (be ParamOps) GetType() interface{} {
	return &models.Param{}
}

func (be ParamOps) GetId(obj interface{}) (string, error) {
	param, ok := obj.(*models.Param)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to param create")
	}
	return *param.Name, nil
}

func (be ParamOps) GetIndexes() map[string]string {
	b := &backend.Param{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be ParamOps) List(parms map[string]string) (interface{}, error) {
	params := params.NewListParamsParams()
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
		case "Name":
			params = params.WithName(&v)
		}
	}
	d, e := session.Params.ListParams(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Get(id string) (interface{}, error) {
	d, e := session.Params.GetParam(params.NewGetParamParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Create(obj interface{}) (interface{}, error) {
	param, ok := obj.(*models.Param)
	if !ok {
		name, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to param create")
		}
		param = &models.Param{Name: &name}
	}
	d, e := session.Params.CreateParam(params.NewCreateParamParams().WithBody(param), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to param patch")
	}
	d, e := session.Params.PatchParam(params.NewPatchParamParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Delete(id string) (interface{}, error) {
	d, e := session.Params.DeleteParam(params.NewDeleteParamParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addParamCommands()
	App.AddCommand(tree)
}

func addParamCommands() (res *cobra.Command) {
	singularName := "param"
	name := "params"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &ParamOps{CommonOps{Name: name, SingularName: singularName}}
	commands := commonOps(mo)
	res.AddCommand(commands...)
	return res
}
