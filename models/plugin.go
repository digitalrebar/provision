package models

// Plugin represents a single instance of a running plugin.
// This contains the configuration need to start this plugin instance.
// swagger:model
type Plugin struct {
	Validation
	Access
	Meta
	// The name of the plugin instance.  THis must be unique across all
	// plugins.
	//
	// required: true
	Name string
	// A description of this plugin.  This can contain any reference
	// information for humans you want associated with the plugin.
	Description string
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

func (b *Plugin) SliceOf() interface{} {
	s := []*Plugin{}
	return &s
}

func (b *Plugin) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Plugin)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// match Paramer interface
func (b *Plugin) GetParams() map[string]interface{} {
	return copyMap(b.Params)
}

func (b *Plugin) SetParams(p map[string]interface{}) {
	b.Params = copyMap(p)
}
