package models

import (
	"fmt"
	"strings"

	"github.com/VictorLowther/jsonpatch2/utils"
)

// RawModel is a raw model that Plugins can specialize to save custom
// data in the dr-provision backing store.
type RawModel map[string]interface{}

func (r *RawModel) String() string {
	return fmt.Sprintf("%v:%v %v", (*r)["Type"], (*r)["Id"], (*r)["Params"])
}

// Access Interface
func (r *RawModel) IsReadOnly() bool {
	b, ok := (*r)["ReadOnly"]
	if !ok {
		return false
	}
	return b.(bool)
}

func (r *RawModel) SetReadOnly(v bool) {
	(*r)["ReadOnly"] = v
}

func (r *RawModel) SetBundle(name string) {
	(*r)["Bundle"] = name
}

// Owner Interface
func (r *RawModel) GetEndpoint() string {
	sobj, _ := r.GetStringField("Endpoint")
	return sobj
}

func (r *RawModel) SetEndpoint(n string) {
	(*r)["Endpoint"] = n
}

// Helpers to get fields
func (r *RawModel) GetStringField(field string) (string, bool) {
	if val, ok := (*r)[field]; ok {
		if sval, ok := val.(string); ok {
			return sval, true
		}
	}
	return "", false
}

func (r *RawModel) getValidated() bool {
	b, ok := (*r)["Validated"]
	if !ok {
		return false
	}
	return b.(bool)
}
func (r *RawModel) setValidated(v bool) {
	if !v {
		delete((*r), "Validated")
		return
	}
	(*r)["Validated"] = true
}
func (r *RawModel) getAvailable() bool {
	b, ok := (*r)["Available"]
	if !ok {
		return false
	}
	return b.(bool)
}
func (r *RawModel) setAvailable(v bool) {
	if !v {
		delete((*r), "Available")
		return
	}
	(*r)["Available"] = true
}

func (r *RawModel) getErrors() []string {
	e, ok := (*r)["Errors"]
	if !ok {
		return []string{}
	}
	switch v := e.(type) {
	case []string:
		return v
	case []interface{}:
		res := make([]string, len(v))
		var ok bool
		for i, val := range v {
			res[i], ok = val.(string)
			if !ok {
				res[i] = fmt.Sprintf("%v", val)
			}
		}
		return res
	default:
		return []string{"Errors field exists and cannot be stringified!"}
	}
}

// Validator Interface
func (r *RawModel) SaveValidation() *Validation {
	return &Validation{
		Validated: r.getValidated(),
		Available: r.getAvailable(),
		Errors:    r.getErrors(),
	}
}

func (r *RawModel) RestoreValidation(or *RawModel) {
	r.setValidated(or.getValidated())
	r.setAvailable(or.getAvailable())
	(*r)["Errors"] = or.getErrors()
}

func (r *RawModel) ClearValidation() {
	r.setValidated(false)
	r.setAvailable(false)
	(*r)["Errors"] = []string{}
}

func (r *RawModel) ForceChange() {
	(*r)["forceChange"] = true
}

func (r *RawModel) ChangeForced() bool {
	return r != nil && (*r)["forceChange"] != nil && (*r)["forceChange"].(bool)
}

func (r *RawModel) Errorf(fmtStr string, args ...interface{}) {
	r.setAvailable(false)
	e := r.getErrors()
	(*r)["Errors"] = append(e, r.Key()+": "+fmt.Sprintf(fmtStr, args...))
}

func (r *RawModel) AddError(err error) {
	if err != nil {
		e := r.getErrors()
		switch o := err.(type) {
		case *Validation:
			e = append(e, o.Errors...)
		case *Error:
			e = append(e, o.Messages...)
		default:
			e = append(e, r.Key()+": "+err.Error())
		}
		(*r)["Errors"] = e
	}
}

func (r *RawModel) HasError() error {
	if len(r.getErrors()) == 0 {
		return nil
	}
	return r
}

func (r *RawModel) Useable() bool {
	return r.getValidated()
}

func (r *RawModel) IsAvailable() bool {
	return r.getAvailable()
}

func (r *RawModel) SetInvalid() bool {
	r.setValidated(false)
	return false
}

func (r *RawModel) SetValid() bool {
	b := len(r.getErrors()) == 0
	r.setValidated(b)
	return b
}

func (r *RawModel) SetAvailable() bool {
	b := len(r.getErrors()) == 0
	r.setAvailable(b)
	return b
}

func (r *RawModel) Error() string {
	return strings.Join(r.getErrors(), "\n")
}

func (r *RawModel) MakeError(code int, errType string, obj Model) error {
	if len(r.getErrors()) == 0 {
		return nil
	}
	return &Error{
		Model:    obj.Prefix(),
		Key:      obj.Key(),
		Code:     code,
		Type:     errType,
		Messages: r.getErrors(),
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

func (r *RawModel) IsPartial() bool {
	b, ok := (*r)["Partial"]
	if !ok {
		return false
	}
	return b.(bool)
}

func (r *RawModel) SetPartial() {
	(*r)["Partial"] = true
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
	// Validation
	if (*r)["Validated"] == nil {
		(*r)["Validated"] = false
	}
	if (*r)["Available"] == nil {
		(*r)["Available"] = false
	}
	if (*r)["Errors"] == nil {
		(*r)["Errors"] = []string{}
	}
	// Meta
	if (*r)["Meta"] == nil {
		(*r)["Meta"] = Meta{}
	}
	// Params
	if (*r)["Params"] == nil {
		(*r)["Params"] = map[string]interface{}{}
	}
	// Bundled
	if (*r)["Bundle"] == nil {
		(*r)["Bundle"] = ""
	}
	// Owner
	if (*r)["Endpoint"] == nil {
		(*r)["Endpoint"] = ""
	}
	// Accessible
	if (*r)["ReadOnly"] == nil {
		(*r)["ReadOnly"] = false
	}
	// Other common fields
	if (*r)["Documentation"] == nil {
		(*r)["Documentation"] = ""
	}
	if (*r)["Description"] == nil {
		(*r)["Description"] = ""
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
