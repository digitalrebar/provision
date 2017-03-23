package main

import (
	"fmt"

	"github.com/rackn/rocket-skates/client/params"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type ParamOps struct{}

func (be ParamOps) GetType() interface{} {
	return &models.Param{}
}

func (be ParamOps) List() (interface{}, error) {
	d, e := session.Params.ListParams(params.NewListParamsParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Get(id string) (interface{}, error) {
	d, e := session.Params.GetParam(params.NewGetParamParams().WithName(id))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Create(obj interface{}) (interface{}, error) {
	param, ok := obj.(*models.Param)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to param create")
	}
	d, e := session.Params.CreateParam(params.NewCreateParamParams().WithBody(param))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Put(id string, obj interface{}) (interface{}, error) {
	param, ok := obj.(*models.Param)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to param put")
	}
	d, e := session.Params.PutParam(params.NewPutParamParams().WithName(id).WithBody(param))
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
	d, e := session.Params.PatchParam(params.NewPatchParamParams().WithName(id).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ParamOps) Delete(id string) (interface{}, error) {
	d, e := session.Params.DeleteParam(params.NewDeleteParamParams().WithName(id))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addParamCommands()
	app.AddCommand(tree)
}

func addParamCommands() (res *cobra.Command) {
	singularName := "param"
	name := "params"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := commonOps(singularName, name, &ParamOps{})
	res.AddCommand(commands...)
	return res
}
