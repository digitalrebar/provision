package models

import (
	"fmt"
)

type RawModel struct {
	Validation
	Access
	Meta
	Type   string
	Id     string
	Params map[string]interface{}
}

func (r *RawModel) String() string {
	return fmt.Sprintf("%s:%s %v", r.Type, r.Id, r.Params)
}

func (r *RawModel) GetMeta() Meta {
	return r.Meta
}

func (r *RawModel) SetMeta(d Meta) {
	r.Meta = d
}

// match Paramer interface
func (r *RawModel) GetParams() map[string]interface{} {
	return copyMap(r.Params)
}

func (r *RawModel) SetParams(p map[string]interface{}) {
	r.Params = copyMap(p)
}

func (r *RawModel) Prefix() string {
	return r.Type
}

func (r *RawModel) Key() string {
	return r.Id
}

func (r *RawModel) KeyName() string {
	return "Id"
}

func (r *RawModel) Fill() {
	r.Validation.fill()
	if r.Meta == nil {
		r.Meta = Meta{}
	}
	if r.Params == nil {
		r.Params = map[string]interface{}{}
	}
	return
}

func (r *RawModel) AuthKey() string {
	return r.Key()
}

func (r *RawModel) SliceOf() interface{} {
	return &[]RawModel{}
}

func (r *RawModel) ToModels(obj interface{}) []Model {
	items := obj.(*[]*RawModel)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (r *RawModel) CanHaveActions() bool {
	return true
}
