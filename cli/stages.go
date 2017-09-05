package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/backend"
	stages "github.com/digitalrebar/provision/client/stages"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/spf13/cobra"
)

type StageOps struct{ CommonOps }

func (be StageOps) GetType() interface{} {
	return &models.Stage{}
}

func (be StageOps) GetId(obj interface{}) (string, error) {
	stage, ok := obj.(*models.Stage)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to stage create")
	}
	return *stage.Name, nil
}

func (be StageOps) GetIndexes() map[string]string {
	b := &backend.Stage{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be StageOps) List(parms map[string]string) (interface{}, error) {
	params := stages.NewListStagesParams()
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
		case "Available":
			params = params.WithAvailable(&v)
		case "Valid":
			params = params.WithValid(&v)
		case "Name":
			params = params.WithName(&v)
		case "BootEnv":
			params = params.WithBootEnv(&v)
		case "Reboot":
			params = params.WithReboot(&v)
		}
	}

	d, e := session.Stages.ListStages(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be StageOps) Get(id string) (interface{}, error) {
	d, e := session.Stages.GetStage(stages.NewGetStageParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be StageOps) Create(obj interface{}) (interface{}, error) {
	stage, ok := obj.(*models.Stage)
	if !ok {
		name, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to stage create")
		}
		stage = &models.Stage{Name: &name}
	}
	d, e := session.Stages.CreateStage(stages.NewCreateStageParams().WithBody(stage), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be StageOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to stage patch")
	}
	d, e := session.Stages.PatchStage(stages.NewPatchStageParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be StageOps) Delete(id string) (interface{}, error) {
	d, e := session.Stages.DeleteStage(stages.NewDeleteStageParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addStageCommands()
	App.AddCommand(tree)
}

func addStageCommands() (res *cobra.Command) {
	singularName := "stage"
	name := "stages"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	beo := &StageOps{CommonOps{Name: name, SingularName: singularName}}
	commands := commonOps(beo)

	commands = append(commands, &cobra.Command{
		Use:   "bootenv [id] [bootenv]",
		Short: fmt.Sprintf("Set the stage's bootenv"),
		Long:  `Helper function to update the stage's bootenv.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithString(args[0], "{ \"BootEnv\": \""+args[1]+"\" }", beo)
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "addprofile [id] [profile]",
		Short: fmt.Sprintf("Add profile to the stage's profile list"),
		Long:  `Helper function to add a profile to the stage's profile list.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithFunction(args[0], beo, func(data interface{}) (interface{}, bool) {
				stage, _ := data.(*models.Stage)
				stage.Profiles = append(stage.Profiles, args[1])
				return stage, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "removeprofile [id] [profile]",
		Short: fmt.Sprintf("Remove a profile from the stage's list"),
		Long:  `Helper function to update the stage's profile list by removing one.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithFunction(args[0], beo, func(data interface{}) (interface{}, bool) {
				changed := false
				stage, _ := data.(*models.Stage)
				newProfiles := []string{}
				for _, s := range stage.Profiles {
					if s == args[1] {
						changed = true
						continue
					}
					newProfiles = append(newProfiles, s)
				}
				stage.Profiles = newProfiles
				if len(newProfiles) == 0 {
					stage.Profiles = nil
				}
				return stage, changed
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "addtask [id] [task]",
		Short: fmt.Sprintf("Add task to the stage's task list"),
		Long:  `Helper function to add a task to the stage's task list.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithFunction(args[0], beo, func(data interface{}) (interface{}, bool) {
				stage, _ := data.(*models.Stage)
				stage.Tasks = append(stage.Tasks, args[1])
				return stage, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "removetask [id] [task]",
		Short: fmt.Sprintf("Remove a task from the stage's list"),
		Long:  `Helper function to update the stage's task list by removing one.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithFunction(args[0], beo, func(data interface{}) (interface{}, bool) {
				changed := false
				stage, _ := data.(*models.Stage)
				newTasks := []string{}
				for _, s := range stage.Tasks {
					if s == args[1] {
						changed = true
						continue
					}
					newTasks = append(newTasks, s)
				}
				stage.Tasks = newTasks
				if len(newTasks) == 0 {
					stage.Tasks = nil
				}
				return stage, changed
			})
		},
	})

	res.AddCommand(commands...)
	return res
}
