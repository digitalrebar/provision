package models

// Plugin represents a single instance of a running plugin.
// This contains the configuration need to start this plugin instance.
// swagger:model
type Plugin struct {
	Validation
	Access
	MetaData
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

func (n *Plugin) Prefix() string {
	return "plugins"
}

func (n *Plugin) Key() string {
	return n.Name
}

func (n *Plugin) Fill() {
	n.MetaData.fill()
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
