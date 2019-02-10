package models

// Plugin represents a single instance of a running plugin.
// This contains the configuration need to start this plugin instance.
// swagger:model
type Plugin struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// The name of the plugin instance.  THis must be unique across all
	// plugins.
	//
	// required: true
	Name string
	// A description of this plugin.  This can contain any reference
	// information for humans you want associated with the plugin.
	Description string
	// Documentation of this plugin.  This should tell what
	// the plugin is for, any special considerations that
	// should be taken into account when using it, etc. in rich structured text (rst).
	Documentation string
	// Any additional parameters that may be needed to configure
	// the plugin.
	Params map[string]interface{}
	// The plugin provider for this plugin
	//
	// required: true
	Provider string
	// Error unrelated to the object validity, but the execution
	// of the plugin.
	PluginErrors []string
}

func (p *Plugin) GetMeta() Meta {
	return p.Meta
}

func (p *Plugin) SetMeta(d Meta) {
	p.Meta = d
}

func (p *Plugin) GetDocumentation() string {
	return p.Documentation
}

func (p *Plugin) Validate() {
	p.AddError(ValidName("Invalid Name", p.Name))
	p.AddError(ValidName("Invalid Provider", p.Provider))
	for k := range p.Params {
		p.AddError(ValidParamName("Invalid Param Name", k))
	}
}

func (p *Plugin) SetName(s string) {
	p.Name = s
}

func (n *Plugin) Prefix() string {
	return "plugins"
}

func (n *Plugin) Key() string {
	return n.Name
}

func (n *Plugin) KeyName() string {
	return "Name"
}

func (n *Plugin) Fill() {
	if n.Meta == nil {
		n.Meta = Meta{}
	}
	n.Validation.fill()
	if n.Params == nil {
		n.Params = map[string]interface{}{}
	}
	if n.PluginErrors == nil {
		n.PluginErrors = []string{}
	}
}

func (n *Plugin) AuthKey() string {
	return n.Key()
}

func (p *Plugin) SliceOf() interface{} {
	s := []*Plugin{}
	return &s
}

func (p *Plugin) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Plugin)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// match Paramer interface
func (p *Plugin) GetParams() map[string]interface{} {
	return copyMap(p.Params)
}

func (p *Plugin) SetParams(pl map[string]interface{}) {
	p.Params = copyMap(pl)
}

func (p *Plugin) CanHaveActions() bool {
	return true
}
