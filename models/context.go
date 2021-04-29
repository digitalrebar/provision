package models

// Context defines an alternate task execution environment for a machine.
// This allows Digital Rebar to manage and run tasks against machines that
// mey not be able to run the Agent.  See https://provision.readthedocs.io/en/latest/doc/arch/provision.html#context
// for more detailed information on how to make an environment for a Context.
//
// swagger:model
type Context struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Name is the name of this Context.  It must be unique.
	Name string `index:",key"`
	// Image is the name of the prebuilt execution environment that the Engine should use to create
	// specific execution environments for this Context when Tasks should run on behalf
	// of a Machine.  Images must contain all the tools needed to run the Tasks
	// that are designed to run in them, as well as a version of drpcli
	// with a context-aware `machines processjobs` command.
	Image string
	// Engine is the name of the Plugin that provides the functionality
	// needed to manage the execution environment that Tasks run in on
	// behalf of a given Machine in the Context.  An Engine could be a
	// Plugin that interfaces with Docker or Podman locally, Kubernetes,
	// Rancher, vSphere, AWS, or any number of other things.
	Engine string

	// Description is a one-line summary of the purpose of this Context
	Description string
	// Documentation should contain any special notes or caveats to keep in mind
	// when using this Context.
	Documentation string
}

func (c *Context) Prefix() string {
	return "contexts"
}

func (c *Context) Key() string {
	return c.Name
}

func (c *Context) KeyName() string {
	return "Name"
}

func (c *Context) AuthKey() string {
	return c.Key()
}

func (c *Context) Fill() {
	c.Validation.fill(c)
	if c.Meta == nil {
		c.Meta = Meta{}
	}
}

func (c *Context) SliceOf() interface{} {
	s := []*Context{}
	return &s
}

func (c *Context) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Context)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (c *Context) GetMeta() Meta {
	return c.Meta
}

func (c *Context) SetMeta(d Meta) {
	c.Meta = d
}
