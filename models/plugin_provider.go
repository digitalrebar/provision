package models

// Plugin Provider describes the available functions that could be
// instantiated by a plugin.
// swagger:model
type PluginProvider struct {
	Meta

	Name    string
	Version string

	// This is used to indicate what version the plugin is built for
	PluginVersion int

	HasPublish       bool
	AvailableActions []AvailableAction

	RequiredParams []string
	OptionalParams []string

	// Object prefixes that can be accessed by this plugin.
	StoreObjects []string

	// Documentation of this plugin provider.  This should tell what
	// the plugin provider is for, any special considerations that
	// should be taken into account when using it, etc. in rich structured text (rst).
	Documentation string

	// Content Bundle Yaml string - can be optional or empty
	Content string
}

func (p *PluginProvider) GetMeta() Meta {
	return p.Meta
}

func (p *PluginProvider) SetMeta(d Meta) {
	p.Meta = d
}

func (p *PluginProvider) GetDocumentation() string {
	return p.Documentation
}

func (p *PluginProvider) Prefix() string  { return "plugin_providers" }
func (p *PluginProvider) Key() string     { return p.Name }
func (p *PluginProvider) KeyName() string { return "Name" }

func (p *PluginProvider) SliceOf() interface{} {
	s := []*PluginProvider{}
	return &s
}

func (p *PluginProvider) ToModels(obj interface{}) []Model {
	items := obj.(*[]*PluginProvider)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (p *PluginProvider) Fill() {
	if p.Meta == nil {
		p.Meta = Meta{}
	}
	if p.RequiredParams == nil {
		p.RequiredParams = []string{}
	}
	if p.OptionalParams == nil {
		p.OptionalParams = []string{}
	}
	if p.AvailableActions == nil {
		p.AvailableActions = []AvailableAction{}
	}
	for _, a := range p.AvailableActions {
		a.Fill()
	}
}

// swagger:model
type PluginProviderUploadInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// Plugins can provide actions for machines
// Assumes that there are parameters on the
// call in addition to the machine.
//
// swagger:model
type AvailableAction struct {
	Provider       string
	Model          string
	Command        string
	RequiredParams []string
	OptionalParams []string
}

func (a *AvailableAction) Fill() {
	if a.RequiredParams == nil {
		a.RequiredParams = []string{}
	}
	if a.OptionalParams == nil {
		a.OptionalParams = []string{}
	}
}

//
// Params is built from the caller, plus
// the machine, plus profiles, plus global.
//
// This is used by the frontend to talk to
// the plugin.
//
type Action struct {
	Model   interface{}
	Plugin  string
	Command string
	Params  map[string]interface{}
}

func (m *Action) Fill() {
	if m.Params == nil {
		m.Params = map[string]interface{}{}
	}
}
