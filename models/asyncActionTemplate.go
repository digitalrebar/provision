package models

// AsyncActionTemplate is the base format of the AsyncAction
//
// swagger:model
type AsyncActionTemplate struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Name is the key of this particular AsyncActionTemplate.
	// required: true
	Name string `index:",key"`
	// Tasks is a list of strings that match the same as the machine's Task list.
	// stages, actions, and other things are allowed as well.
	Tasks []string
	// Description is a one-line description of the parameter.
	Description string
	// Documentation details what the parameter does, what values it can
	// take, what it is used for, etc.
	Documentation string
	// An array of profiles to apply to this machine in order when looking
	// for a parameter during rendering.
	Profiles []string
	// The Parameters that have been directly set on the AsyncActionTemplate.
	Params map[string]interface{}
}

func (aat *AsyncActionTemplate) GetMeta() Meta {
	return aat.Meta
}

func (aat *AsyncActionTemplate) SetMeta(d Meta) {
	aat.Meta = d
}

// GetDocumentaiton returns the object's Documentation
func (aat *AsyncActionTemplate) GetDocumentation() string {
	return aat.Documentation
}

// GetDescription returns the object's Description
func (aat *AsyncActionTemplate) GetDescription() string {
	return aat.Description
}

func (aat *AsyncActionTemplate) Validate() {
	aat.AddError(ValidName("Invalid AsyncActionTemplate Name", aat.Name))
}

func (aat *AsyncActionTemplate) Prefix() string {
	return "async_action_templates"
}

func (aat *AsyncActionTemplate) Key() string {
	return aat.Name
}

func (aat *AsyncActionTemplate) KeyName() string {
	return "Name"
}

func (aat *AsyncActionTemplate) Fill() {
	if aat.Meta == nil {
		aat.Meta = Meta{}
	}
	if aat.Profiles == nil {
		aat.Profiles = []string{}
	}
	if aat.Params == nil {
		aat.Params = map[string]interface{}{}
	}
	if aat.Tasks == nil {
		aat.Tasks = []string{}
	}
	aat.Validation.fill(aat)
}

func (aat *AsyncActionTemplate) AuthKey() string {
	return aat.Key()
}

func (aat *AsyncActionTemplate) SliceOf() interface{} {
	s := []*AsyncActionTemplate{}
	return &s
}

func (aat *AsyncActionTemplate) ToModels(obj interface{}) []Model {
	items := obj.(*[]*AsyncActionTemplate)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (aat *AsyncActionTemplate) CanHaveActions() bool {
	return true
}

// match Profiler interface

// GetProfiles gets the profiles on this stage
func (aat *AsyncActionTemplate) GetProfiles() []string {
	return aat.Profiles
}

// SetProfiles sets the profiles on this stage
func (aat *AsyncActionTemplate) SetProfiles(p []string) {
	aat.Profiles = p
}

// match Paramer interface

// GetParams gets the parameters on this stage
func (aat *AsyncActionTemplate) GetParams() map[string]interface{} {
	return copyMap(aat.Params)
}

// SetParams sets the parameters on this stage
func (aat *AsyncActionTemplate) SetParams(p map[string]interface{}) {
	aat.Params = copyMap(p)
}

// match TaskRunner interface

// GetTasks returns the tasks associated with this stage
func (aat *AsyncActionTemplate) GetTasks() []string {
	return aat.Tasks
}

// SetTasks sets the tasks in this stage
func (aat *AsyncActionTemplate) SetTasks(t []string) {
	aat.Tasks = t
}

// SetName sets the name of the object
func (aat *AsyncActionTemplate) SetName(n string) {
	aat.Name = n
}
