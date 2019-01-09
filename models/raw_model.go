package models

import (
	"fmt"
	"strings"

	"github.com/VictorLowther/jsonpatch2/utils"
)

type RawModel map[string]interface{}

func (r *RawModel) String() string {
	return fmt.Sprintf("%v:%v %v", (*r)["Type"], (*r)["Id"], (*r)["Params"])
}

// Access Interface
func (r *RawModel) IsReadOnly() bool {
	return (*r)["ReadOnly"].(bool)
}

// Owner Interface
func (r *RawModel) GetEndpoint() string {
	return (*r)["Endpoint"].(string)
}

// Validator Interface
func (r *RawModel) SaveValidation() *Validation {
	return &Validation{
		Validated: (*r)["Validated"].(bool),
		Available: (*r)["Available"].(bool),
		Errors:    (*r)["Errors"].([]string),
	}
}

func (r *RawModel) RestoreValidation(or *RawModel) {
	(*r)["Validated"] = (*or)["Validated"]
	(*r)["Available"] = (*or)["Available"]
	(*r)["Errors"] = (*or)["Errors"]
}

func (r *RawModel) ClearValidation() {
	(*r)["Validated"] = false
	(*r)["Available"] = false
	(*r)["Errors"] = []string{}
}

func (r *RawModel) ForceChange() {
	(*r)["forceChange"] = true
}

func (r *RawModel) ChangeForced() bool {
	return r != nil && (*r)["forceChange"] != nil && (*r)["forceChange"].(bool)
}

func (r *RawModel) Errorf(fmtStr string, args ...interface{}) {
	(*r)["Available"] = false
	if (*r)["Errors"] == nil {
		(*r)["Errors"] = []string{}
	}
	(*r)["Errors"] = append((*r)["Errors"].([]string), fmt.Sprintf(fmtStr, args...))
}

func (r *RawModel) AddError(err error) {
	if err != nil {
		if (*r)["Errors"] == nil {
			(*r)["Errors"] = []string{}
		}
		switch o := err.(type) {
		case *Validation:
			(*r)["Errors"] = append((*r)["Errors"].([]string), o.Errors...)
		case *Error:
			(*r)["Errors"] = append((*r)["Errors"].([]string), o.Messages...)
		default:
			(*r)["Errors"] = append((*r)["Errors"].([]string), err.Error())
		}
	}
}

func (r *RawModel) HasError() error {
	if len((*r)["Errors"].([]string)) == 0 {
		return nil
	}
	return r
}

func (r *RawModel) Useable() bool {
	return (*r)["Validated"].(bool)
}

func (r *RawModel) IsAvailable() bool {
	return (*r)["Available"].(bool)
}

func (r *RawModel) SetInvalid() bool {
	(*r)["Validated"] = false
	return false
}

func (r *RawModel) SetValid() bool {
	(*r)["Validated"] = (*r)["Validated"].(bool) || len((*r)["Errors"].([]string)) == 0
	return (*r)["Validated"].(bool)
}

func (r *RawModel) SetAvailable() bool {
	(*r)["Available"] = (*r)["Available"].(bool) || len((*r)["Errors"].([]string)) == 0
	return (*r)["Available"].(bool)
}

func (r *RawModel) Error() string {
	return strings.Join((*r)["Errors"].([]string), "\n")
}

func (r *RawModel) MakeError(code int, errType string, obj Model) error {
	if len((*r)["Errors"].([]string)) == 0 {
		return nil
	}
	return &Error{
		Model:    obj.Prefix(),
		Key:      obj.Key(),
		Code:     code,
		Type:     errType,
		Messages: (*r)["Errors"].([]string),
	}
}

// MetaHaver Interface
func (r *RawModel) GetMeta() Meta {
	if m, ok := (*r)["Meta"].(Meta); ok {
		return m
	}
	m := Meta{}
	if utils.Remarshal((*r)["Meta"], &m) != nil {
		return Meta{}
	}
	return m
}

func (r *RawModel) SetMeta(d Meta) {
	(*r)["Meta"] = d
}

// match Paramer interface
func (r *RawModel) GetParams() map[string]interface{} {
	return copyMap((*r)["Params"].(map[string]interface{}))
}

func (r *RawModel) SetParams(p map[string]interface{}) {
	(*r)["Params"] = copyMap(p)
}

func (r *RawModel) Prefix() string {
	return (*r)["Type"].(string)
}

func (r *RawModel) Key() string {
	s, ok := (*r)["Id"]
	if !ok {
		return ""
	}
	return s.(string)
}

func (r *RawModel) KeyName() string {
	return "Id"
}

func (r *RawModel) Fill() {
	boolFields := []string{"Available", "Validated", "forceChange", "ReadOnly"}
	for _, f := range boolFields {
		if (*r)[f] == nil {
			(*r)[f] = false
		}
	}
	if (*r)["Errors"] == nil {
		(*r)["Errors"] = []string{}
	}
	if (*r)["Meta"] == nil {
		(*r)["Meta"] = Meta{}
	}
	if (*r)["Params"] == nil {
		(*r)["Params"] = map[string]interface{}{}
	}
	if (*r)["Documentation"] == nil {
		(*r)["Documentation"] = ""
	}
	if (*r)["Description"] == nil {
		(*r)["Description"] = ""
	}
	if (*r)["Endpoint"] == nil {
		(*r)["Endpoint"] = ""
	}
	return
}

func (r *RawModel) AuthKey() string {
	return r.Key()
}

func (r *RawModel) SliceOf() interface{} {
	return &[]*RawModel{}
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
