package models

// AsyncActionCron is time system to apply an async action template to something
// It copies the params and profiles into the async action, applies it to the filtered
// machines at the specified cron-style time.
//
// swagger:model
type AsyncActionCron struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Name is the key of this particular AsyncActionCron.
	// required: true
	Name string `index:",key"`
	// Cron is cron string to indicate how often the entry should be applied
	Cron string
	// Filter is a "list"-style filter string to find machines to apply the cron too
	// Filter is already assumed to have AsyncActionMode == true && Runnable == true
	Filter string
	// AsyncActionTemplate is template to apply
	AsyncActionTemplate string
	// Description is a one-line description of the parameter.
	Description string
	// Documentation details what the parameter does, what values it can
	// take, what it is used for, etc.
	Documentation string
	// Profiles to apply to this machine in order when looking
	// for a parameter during rendering.
	Profiles []string
	// Params that have been directly set on the AsyncActionCron.
	Params map[string]interface{}
}

func (aac *AsyncActionCron) GetMeta() Meta {
	return aac.Meta
}

func (aac *AsyncActionCron) SetMeta(d Meta) {
	aac.Meta = d
}

// GetDocumentaiton returns the object's Documentation
func (aac *AsyncActionCron) GetDocumentation() string {
	return aac.Documentation
}

// GetDescription returns the object's Description
func (aac *AsyncActionCron) GetDescription() string {
	return aac.Description
}

func (aac *AsyncActionCron) Validate() {
	aac.AddError(ValidName("Invalid AsyncActionCron Name", aac.Name))
}

func (aac *AsyncActionCron) Prefix() string {
	return "async_action_crons"
}

func (aac *AsyncActionCron) Key() string {
	return aac.Name
}

func (aac *AsyncActionCron) KeyName() string {
	return "Name"
}

func (aac *AsyncActionCron) Fill() {
	if aac.Meta == nil {
		aac.Meta = Meta{}
	}
	if aac.Profiles == nil {
		aac.Profiles = []string{}
	}
	if aac.Params == nil {
		aac.Params = map[string]interface{}{}
	}
	aac.Validation.fill(aac)
}

func (aac *AsyncActionCron) AuthKey() string {
	return aac.Key()
}

func (aac *AsyncActionCron) SliceOf() interface{} {
	s := []*AsyncActionCron{}
	return &s
}

func (aac *AsyncActionCron) ToModels(obj interface{}) []Model {
	items := obj.(*[]*AsyncActionCron)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (aac *AsyncActionCron) CanHaveActions() bool {
	return true
}

// match Profiler interface

// GetProfiles gets the profiles on this stage
func (aac *AsyncActionCron) GetProfiles() []string {
	return aac.Profiles
}

// SetProfiles sets the profiles on this stage
func (aac *AsyncActionCron) SetProfiles(p []string) {
	aac.Profiles = p
}

// match Paramer interface

// GetParams gets the parameters on this stage
func (aac *AsyncActionCron) GetParams() map[string]interface{} {
	return copyMap(aac.Params)
}

// SetParams sets the parameters on this stage
func (aac *AsyncActionCron) SetParams(p map[string]interface{}) {
	aac.Params = copyMap(p)
}

// SetName sets the name of the object
func (aac *AsyncActionCron) SetName(n string) {
	aac.Name = n
}
