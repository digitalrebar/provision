package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/tasks"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type TaskOps struct{}

func (be TaskOps) GetType() interface{} {
	return &models.Task{}
}

func (be TaskOps) GetId(obj interface{}) (string, error) {
	task, ok := obj.(*models.Task)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to task create")
	}
	return *task.Name, nil
}

func (be TaskOps) GetIndexes() map[string]string {
	b := &backend.Task{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be TaskOps) List(parms map[string]string) (interface{}, error) {
	params := tasks.NewListTasksParams()
	if listLimit != -1 {
		t1 := int64(listLimit)
		params = params.WithLimit(&t1)
	}
	if listOffset != -1 {
		t1 := int64(listOffset)
		params = params.WithOffset(&t1)
	}
	for k, v := range parms {
		switch k {
		case "Name":
			params = params.WithName(&v)
		}
	}
	d, e := session.Tasks.ListTasks(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TaskOps) Get(id string) (interface{}, error) {
	d, e := session.Tasks.GetTask(tasks.NewGetTaskParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TaskOps) Create(obj interface{}) (interface{}, error) {
	task, ok := obj.(*models.Task)
	if !ok {
		name, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to task create")
		}
		task = &models.Task{Name: &name}
	}
	d, e := session.Tasks.CreateTask(tasks.NewCreateTaskParams().WithBody(task), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TaskOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to task patch")
	}
	d, e := session.Tasks.PatchTask(tasks.NewPatchTaskParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be TaskOps) Delete(id string) (interface{}, error) {
	d, e := session.Tasks.DeleteTask(tasks.NewDeleteTaskParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addTaskCommands()
	App.AddCommand(tree)
}

func addTaskCommands() (res *cobra.Command) {
	singularName := "task"
	name := "tasks"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &TaskOps{}
	commands := commonOps(singularName, name, mo)
	res.AddCommand(commands...)
	return res
}
