package models

import (
	"fmt"

	"github.com/digitalrebar/provision/v4/store"
)

// Plugin Provider describes the available functions that could be
// instantiated by a plugin.
// swagger:model
type PluginProvider struct {
	Meta

	// Name is the unique name of the PluginProvider.
	// Each Plugin provider must have a unique Name.
	Name string `index:",key"`

	// The version of the PluginProvider.  This is a semver compatible string.
	Version string

	// This is used to indicate what version the plugin is built for
	// This is effectively the API version of the protocol that
	// plugin providers use to communicate with dr-provision.
	// Right now, all plugin providers must set this to version 4,
	// which is the only supported protocol version.
	PluginVersion int

	// If AutoStart is true, a Plugin will be created for this
	// Provider at provider definition time, if one is not already present.
	AutoStart bool

	// HasPlugin is deprecated, plugin provider binaries should use a websocket
	// event stream instead.
	HasPublish bool

	// AvailableActions lists the actions that this PluginProvider
	// can take.
	AvailableActions []AvailableAction

	// RequiredParams and OptionalParams
	// are Params that must be present on a Plugin for the Provider
	// to operate.
	RequiredParams []string
	OptionalParams []string

	// Object prefixes that can be accessed by this plugin.
	// The interface can be empty struct{} or a JSONSchema draft v4
	// This allows PluginProviders to define custom Object types that dr-provision will
	// store and check the validity of.
	StoreObjects map[string]interface{}

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
	if p.StoreObjects == nil {
		p.StoreObjects = map[string]interface{}{}
	}
	for _, a := range p.AvailableActions {
		a.Fill()
	}
}

// Store extracts the content bundle in the Content field of the
// PluginProvider into a Store.
func (p *PluginProvider) Store() (store.Store, error) {
	content := &Content{}
	content.Fill()

	if p.Content != "" {
		codec := store.YamlCodec
		if err := codec.Decode([]byte(p.Content), content); err != nil {
			return nil, err
		}
	}
	cName := p.Name
	content.Meta.Name = cName
	content.Meta.Version = p.Version
	if content.Meta.Description == "" || content.Meta.Description == "Unspecified" {
		content.Meta.Description = fmt.Sprintf("Content layer for %s plugin provider", p.Name)
	}
	if content.Meta.Source == "" || content.Meta.Source == "Unspecified" {
		content.Meta.Source = "FromPluginProvider"
	}
	content.Meta.Type = "plugin"
	meta := content.GenerateMetaMap()
	if v, ok := meta["Color"]; ok {
		meta["color"] = v
	}
	if v, ok := meta["Icon"]; ok {
		meta["icon"] = v
	}
	p.SetMeta(meta)
	s, _ := store.Open("memory:///")
	return s, content.ToStore(s)
}

// AutoPlugin - builds a plugin model if auto start is true, otherwise nil
func (p *PluginProvider) AutoPlugin() *Plugin {
	if p.AutoStart {
		pl := &Plugin{Name: p.Name, Provider: p.Name}
		pl.Fill()
		pl.SetMeta(p.GetMeta())
		return pl
	}
	return nil
}

// swagger:model
type PluginProviderUploadInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// AvailableAction is an Action that a Plugin instantiated by a
// PluginProvider.  Assumes that there are parameters on the call in
// addition to the machine.
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

// Action is an additional command that can be added to other Models
// by a Plugin.
type Action struct {
	Model      interface{}
	Plugin     string
	Command    string
	CommandSet string
	Params     map[string]interface{}
}

func (m *Action) Fill() {
	if m.Params == nil {
		m.Params = map[string]interface{}{}
	}
}
