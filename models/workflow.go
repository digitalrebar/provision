package models

type Workflow struct {
	Validation
	Access
	Meta
	Name        string
	Description string
	Stages      []string
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
