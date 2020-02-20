package models

type Context struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Name is the name of this Context.
	Name string
	// Image the OS image that jobs will execute in when running in this Context.
	// This is usually a Docker container, a VM image, or something similar.
	Image string
	// Engine is the system that runs the Image.  This is something like
	// docker, kubernetes, AWS, or something similar.
	Engine        string
	Description   string
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
