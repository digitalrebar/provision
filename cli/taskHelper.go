package cli

import (
	"fmt"
	"strconv"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func (o *ops) tasks() {
	o.addCommand(&cobra.Command{
		Use:   "addtask [id] [task]",
		Short: fmt.Sprintf("Add task to the %s's task list", o.singleName),
		Long:  fmt.Sprintf(`Helper function to add a task to the %s's task list.`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			data, err := o.refOrFill(args[0])
			if err != nil {
				return err
			}
			ex := models.Clone(data).(models.Tasker)
			ex.SetTasks(append(ex.GetTasks(), args[1]))
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})

	o.addCommand(&cobra.Command{
		Use:   "removetask [id] [task]",
		Short: fmt.Sprintf("Remove a task from the %s's list", o.singleName),
		Long:  fmt.Sprintf(`Helper function to update the %s's task list by removing one.`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			data, err := o.refOrFill(args[0])
			if err != nil {
				return err
			}
			ex := models.Clone(data).(models.Tasker)
			newTasks := []string{}
			for _, s := range ex.GetTasks() {
				if s == args[1] {
					continue
				}
				newTasks = append(newTasks, s)
			}
			ex.SetTasks(newTasks)
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})
	if _, ok := o.example().(models.TaskRunner); ok {
		o.addCommand(&cobra.Command{
			Use:   "inserttask [id] [task] [offset]",
			Short: fmt.Sprintf("Insert a task at [offset] from %s's running task", o.singleName),
			Long:  `If [offset] is not present, it is assumed to be just after the current task`,
			Args: func(c *cobra.Command, args []string) error {
				if len(args) < 2 || len(args) > 3 {
					return fmt.Errorf("%v expects 2 or 3 arguments", c.UseLine())
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				offset := 0
				var err error
				if len(args) == 3 {
					offset, err = strconv.Atoi(args[2])
					if err != nil {
						return err
					}
				}
				data, err := o.refOrFill(args[0])
				if err != nil {
					return err
				}
				ex := models.Clone(data).(models.TaskRunner)
				rt := ex.RunningTask()
				taskList := ex.GetTasks()
				var immutable, mutable []string
				if rt == -1 {
					immutable = []string{}
					mutable = taskList
				} else if rt == len(taskList) {
					immutable = taskList
					mutable = []string{}
				} else {
					immutable = taskList[:rt+1]
					mutable = taskList[rt+1:]
				}
				insertAt := offset
				if offset < 0 {
					insertAt = len(mutable) - offset + 1
				}
				if insertAt < 0 || insertAt > len(mutable) {
					return fmt.Errorf("Invalid offset %v", offset)
				}
				if insertAt == len(mutable) {
					mutable = append(mutable, args[1])
				} else if insertAt == 0 {
					mutable = append([]string{args[1]}, mutable...)
				} else {
					head := make([]string, insertAt)
					copy(head, mutable)
					head = append(head, args[1])
					tail := mutable[insertAt:]
					mutable = append(head, tail...)
				}
				ex.SetTasks(append(immutable, mutable...))
				res, err := session.PatchToFull(data, ex, ref != "")
				if err != nil {
					return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
				}
				return prettyPrint(res)
			},
		})
	}

}
