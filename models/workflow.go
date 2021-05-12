package models

// Workflow contains a list of Stages. When it is applied to a Machine,
// that machine's Tasks list is populated with the contents of the Stages in the Workflow.
//
// swagger:model
type Workflow struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	Name          string `index:",key"`
	Description   string
	Documentation string
	Stages        []string
}

func (w *Workflow) GetMeta() Meta {
	return w.Meta
}

func (w *Workflow) SetMeta(d Meta) {
	w.Meta = d
}

// GetDocumentaiton returns the object's Documentation
func (w *Workflow) GetDocumentation() string {
	return w.Documentation
}

// GetDescription returns the object's Description
func (w *Workflow) GetDescription() string {
	return w.Description
}

func (w *Workflow) Prefix() string {
	return "workflows"
}

func (w *Workflow) Key() string {
	return w.Name
}

func (w *Workflow) KeyName() string {
	return "Name"
}

func (w *Workflow) Fill() {
	w.Validation.fill(w)
	if w.Meta == nil {
		w.Meta = Meta{}
	}
	if w.Stages == nil {
		w.Stages = []string{}
	}
}

func (w *Workflow) AuthKey() string {
	return w.Key()
}

func (w *Workflow) SliceOf() interface{} {
	ws := []*Workflow{}
	return &ws
}

func (w *Workflow) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Workflow)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (w *Workflow) Validate() {
	w.AddError(ValidName("Invalid Name", w.Name))
	for _, stageName := range w.Stages {
		w.AddError(ValidName("Invalid Stage Name", stageName))
	}
}

func (w *Workflow) CanHaveActions() bool {
	return true
}
