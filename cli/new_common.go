package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type ops struct {
	name          string
	singleName    string
	example       func() models.Model
	mustPut       bool
	noCreate      bool
	noUpdate      bool
	noDestroy     bool
	noWait        bool
	extraCommands []*cobra.Command
	actionName    string
}

func (o *ops) refOrFill(key string) (data models.Model, err error) {
	data = o.example()

	if ref == "" {
		if err = session.FillModel(data, key); err != nil {
			return
		}
	} else {
		err = bufOrFileDecode(ref, &data)
	}
	return
}

func (o *ops) addCommand(c *cobra.Command) {
	o.extraCommands = append(o.extraCommands, c)
}

func (o *ops) commands() []*cobra.Command {
	cmds := []*cobra.Command{}
	listCmd := &cobra.Command{
		Use:   "list [filters...]",
		Short: fmt.Sprintf("List all %v", o.name),
		Long: fmt.Sprintf(`This will list all %v by default.
You can narrow down the items returned using index filters.
Use the "indexes" command to get the indexes available for %v.

To filter by indexes, you can use the following stanzas:

* *index* Eq *value*
  This will return items Equal to *value* according to *index*
* *index* Ne *value*
  This will return items Not Equal to *value* according to *index*
* *index* Lt *value*
  This will return items Less Than *value* according to *index*
* *index* Lte *value*
  This will return items Less Than Or Equal to *value* according to *index*
* *index* Gt *value*
  This will return items Greater Than *value* according to *index*
* *index* Gte *value*
  This will return items Greater Than Or Equal to *value* according to *index*
* *index* Between *lower* *upper*
  This will return items Greater Than Or Equal to *lower*
  and Less Than Or Equal to *upper* according to *index*
* *index* Except *lower* *upper*
  This will return items Less Than *lower* or
  Greater Than *upper* according to *index*

You can chain any number of filters together, and they will pipeline into
each other as appropriate.  After the above filters have been applied, you can
further tweak how the results are returned using the following meta-filters:

* 'reverse' to return items in reverse order
* 'limit' *number* to only return the first *number* items
* 'offset' *number* to skip *number* items
* 'sort' *index* to sort items according to *index*
`, o.name, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			if strings.Contains(args[0], "=") {
				for _, a := range args {
					ar := strings.SplitN(a, "=", 2)
					if len(ar) != 2 {
						return fmt.Errorf("Filter argument requires an '=' separator: %s", a)
					}
				}
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			req := session.Req().List(o.name)
			if len(args) > 0 && strings.Contains(args[0], "=") {
				// Old-style structured args
				if listLimit != -1 {
					args = append(args, fmt.Sprintf("limit=%d", listLimit))
				}
				if listOffset != -1 {
					args = append(args, fmt.Sprintf("offset=%d", listOffset))
				}
				pargs := []string{}
				for _, arg := range args {
					a := strings.SplitN(arg, "=", 2)
					pargs = append(pargs, a...)
				}
				req.Params(pargs...)
			} else {
				// New-style freeform args
				if listLimit != -1 {
					args = append(args, "limit", fmt.Sprintf("%d", listLimit))
				}
				if listOffset != -1 {
					args = append(args, "offset", fmt.Sprintf("%d", listOffset))
				}
				if len(args) > 0 {
					req = session.Req().Filter(o.name, args...)
				}
			}
			data := []interface{}{}
			err := req.Do(&data)
			if err != nil {
				return generateError(err, "listing %v", o.name)
			} else {
				return prettyPrint(data)
			}
		},
	}
	cmds = append(cmds, listCmd)
	listCmd.Flags().IntVar(&listLimit, "limit", -1, "Maximum number of items to return")
	listCmd.Flags().IntVar(&listOffset, "offset", -1, "Number of items to skip before starting to return data")
	cmds = append(cmds, &cobra.Command{
		Use:   "indexes",
		Short: fmt.Sprintf("Get indexes for %s", o.name),
		Long:  fmt.Sprintf("Different object types can have indexes on various fields."),
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			indexes, err := session.Indexes(o.name)
			if err != nil {
				return generateError(err, "Error fetching indexes")
			}
			return prettyPrint(indexes)
		},
	})
	cmds = append(cmds, &cobra.Command{
		Use:   "show [id]",
		Short: fmt.Sprintf("Show a single %v by id", o.name),
		Long: fmt.Sprintf(`This will show a %v by ID.
You may also show a single item using a unique index.  In that case,
format id as *index*:*value*
`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			data := o.example()
			if err := session.FillModel(data, args[0]); err != nil {
				return generateError(err, "Failed to fetch %v: %v", o.singleName, args[0])
			} else {
				return prettyPrint(data)
			}
		},
	})
	cmds = append(cmds, &cobra.Command{
		Use:   "exists [id]",
		Short: fmt.Sprintf("See if a %v exists by id", o.name),
		Long:  fmt.Sprintf("This will detect if a %v exists.", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			exists, err := session.ExistsModel(o.name, args[0])
			if err != nil {
				return generateError(err, "Failed to test %v: %v", o.name, args[0])
			}
			if exists {
				return nil
			}
			return fmt.Errorf("%s:%s does not exist", o.name, args[0])
		},
	})
	if !o.noCreate {
		cmds = append(cmds, &cobra.Command{
			Use:   "create [json]",
			Short: fmt.Sprintf("Create a new %v with the passed-in JSON or string key", o.singleName),
			Long: `
As a useful shortcut, '-' can be passed to indicate that the JSON should
be read from stdin.

In either case, for the Machine, BootEnv, User, and Profile objects, a string may be provided to create a new
empty object of that type.  For User, BootEnv, Machine, and Profile, it will be the object's name.
`,
			Args: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				ref := o.example()
				if err := into(args[0], ref); err != nil {
					if args[0] != "-" {
						if tgt, ok := ref.(models.NameSetter); ok {
							tgt.SetName(args[0])
						} else {
							return fmt.Errorf("Unable to create a new %s by name", o.singleName)
						}
					} else {
						return fmt.Errorf("Unable to create a new %s: %v", o.singleName, err)
					}
				}
				if err := session.CreateModel(ref); err != nil {
					return generateError(err, "Unable to create new %v", o.singleName)
				}
				return prettyPrint(ref)
			},
		})
	}
	if !o.noUpdate {
		if o.mustPut {
			cmds = append(cmds, &cobra.Command{
				Use:   "update [id] [json]",
				Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", o.singleName),
				Long:  `As a useful shortcut, '-' can be passed to indicate that the JSON should be read from stdin`,
				Args: func(c *cobra.Command, args []string) error {
					if len(args) != 2 {
						return fmt.Errorf("%v requires 2 arguments", c.UseLine())
					}
					return nil
				},
				RunE: func(c *cobra.Command, args []string) error {
					ref, err := o.refOrFill(args[0])
					if err != nil {
						return err
					}
					toPut, err := mergeFromArgs(ref, args[1])
					if err != nil {
						return generateError(err, "Failed to generate changed %s:%s object", o.name, args[0])
					}
					if err := session.PutModel(toPut); err != nil {
						return generateError(err, "Unable to update %v", args[0])
					}
					return prettyPrint(toPut)
				},
			})
		} else {
			cmds = append(cmds, &cobra.Command{
				Use:   "update [id] [json]",
				Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", o.singleName),
				Long:  `As a useful shortcut, '-' can be passed to indicate that the JSON should be read from stdin`,
				Args: func(c *cobra.Command, args []string) error {
					if len(args) != 2 {
						return fmt.Errorf("%v requires 2 arguments", c.UseLine())
					}
					return nil
				},
				RunE: func(c *cobra.Command, args []string) error {
					refObj, err := o.refOrFill(args[0])
					if err != nil {
						return err
					}
					toPut, err := mergeFromArgs(refObj, args[1])
					if err != nil {
						return generateError(err, "Failed to generate changed %s:%s object", o.name, args[0])
					}
					if res, err := session.PatchToFull(refObj, toPut, ref != ""); err != nil {
						return generateError(err, "Unable to update %v", args[0])
					} else {
						return prettyPrint(res)
					}
				},
			})
		}
	}
	if !o.noDestroy {
		cmds = append(cmds, &cobra.Command{
			Use:   "destroy [id]",
			Short: fmt.Sprintf("Destroy %v by id", o.singleName),
			Long:  fmt.Sprintf("This will destroy the %v.", o.singleName),
			Args: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				_, err := session.DeleteModel(o.name, args[0])
				if err != nil {
					return generateError(err, "Unable to destroy %v %v", o.singleName, args[0])
				}
				fmt.Printf("Deleted %v %v\n", o.singleName, args[0])
				return nil
			},
		})
	}
	if !o.noWait && o.example != nil {
		cmds = append(cmds, &cobra.Command{
			Use:   "wait [id] [field] [value] [timeout]",
			Short: fmt.Sprintf("Wait for a %s's field to become a value within a number of seconds", o.singleName),
			Long: `
This function waits for the value to become the new value.

Timeout is optional, defaults to 1 day, and is measured in seconds.

Returns the following strings:
  complete - field is equal to value
  interrupt - user interrupted the command
  timeout - timeout has exceeded`,
			Args: func(c *cobra.Command, args []string) error {
				if len(args) < 3 {
					return fmt.Errorf("%v requires at least 3 arguments", c.UseLine())
				}
				if len(args) > 4 {
					return fmt.Errorf("%v requires at most 4 arguments", c.UseLine())
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				id := args[0]
				field := args[1]
				value := args[2]
				timeout := time.Hour * 24
				if len(args) == 4 {
					t, e := strconv.ParseInt(args[3], 10, 64)
					if e != nil {
						return e
					}
					timeout = time.Second * time.Duration(t)
				}
				testfn := api.EqualItem(field, value)
				item, err := o.refOrFill(id)
				if err != nil {
					return err
				}
				es, err := session.Events()
				if err != nil {
					return err
				}
				defer es.Close()
				res, err := es.WaitFor(item, testfn, timeout)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		})
	}
	cmds = append(cmds, o.extraCommands...)
	return cmds
}

func (o *ops) params() {
	aggregate := false
	getParams := &cobra.Command{
		Use:   "params [id] [json]",
		Short: fmt.Sprintf("Gets/sets all parameters for the %s", o.singleName),
		Long:  fmt.Sprintf(`A helper function to return all or set all the parameters on the %s`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			if len(args) == 1 {
				req := session.Req().UrlFor(o.name, args[0], "params")
				if aggregate {
					req.Params("aggregate", "true")
				}
				res := map[string]interface{}{}
				if err := req.Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
				return prettyPrint(res)
			}
			val := map[string]interface{}{}
			if err := into(args[1], &val); err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			res := map[string]interface{}{}
			if ref == "" {
				if err := session.Req().Post(val).UrlFor(o.name, args[0], "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				var data map[string]interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				if err := session.Req().ParanoidPatch().PatchObj(data, val).UrlFor(o.name, args[0], "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(res)
		},
	}
	getParams.Flags().BoolVar(&aggregate, "aggregate", false, "Should machine return aggregated view")
	o.addCommand(getParams)
	getParam := &cobra.Command{
		Use:   "get [id] param [key]",
		Short: fmt.Sprintf("Get a parameter from the %s", o.singleName),
		Long:  fmt.Sprintf(`A helper function to return the value of the parameter on the %s`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			var res interface{}
			req := session.Req().UrlFor(o.name, uuid, "params", key)
			if aggregate {
				req.Params("aggregate", "true")
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
			}
			return prettyPrint(res)
		},
	}
	getParam.Flags().BoolVar(&aggregate, "aggregate", false, "Should machine return aggregated view")
	o.addCommand(getParam)
	o.addCommand(&cobra.Command{
		Use:   "add [id] param [key] to [json blob]",
		Short: fmt.Sprintf("Add the %s param *key* to *blob*", o.name),
		Long:  fmt.Sprintf(`Helper function to add parameters to the %s. Fails is already present.`, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			newValue := args[4]
			var value interface{}
			if err := into(newValue, &value); err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}

			res := map[string]interface{}{}
			if ref == "" {
				if err := session.Req().UrlFor(o.name, uuid, "params").Do(&res); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				if err := bufOrFileDecode(ref, &res); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
			}

			if _, ok := res[key]; ok {
				return fmt.Errorf("Key, %s, already present on %s %s", key, o.singleName, uuid)
			}

			var params interface{}
			path := fmt.Sprintf("/%s", makeJsonPtr(key))
			patch := jsonpatch2.Patch{
				jsonpatch2.Operation{
					Op:    "test",
					Path:  "",
					Value: res,
				},
				jsonpatch2.Operation{
					Op:    "add",
					Path:  path,
					Value: value,
				},
			}
			if err := session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&params); err != nil {
				return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
			}
			return prettyPrint(value)
		},
	})
	o.addCommand(&cobra.Command{
		Use:   "set [id] param [key] to [json blob]",
		Short: fmt.Sprintf("Set the %s param *key* to *blob*", o.name),
		Long:  fmt.Sprintf(`Helper function to update the %s parameters.`, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			newValue := args[4]
			var value interface{}
			if err := into(newValue, &value); err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			var params interface{}
			if ref == "" {
				if err := session.Req().Post(value).UrlFor(o.name, uuid, "params", key).Do(&params); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			} else {
				var data interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				path := fmt.Sprintf("/%s", makeJsonPtr(key))
				patch := jsonpatch2.Patch{
					jsonpatch2.Operation{
						Op:    "test",
						Path:  path,
						Value: data,
					},
					jsonpatch2.Operation{
						Op:    "replace",
						Path:  path,
						Value: value,
					},
				}
				if err := session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&params); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(value)
		},
	})
	o.addCommand(&cobra.Command{
		Use:   "remove [id] param [key]",
		Short: fmt.Sprintf("Remove the param *key* from %s", o.name),
		Long:  fmt.Sprintf(`Helper function to update the %s parameters.`, o.name),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			key := args[2]
			var param interface{}
			if ref == "" {
				err := session.Req().Del().UrlFor(o.name, uuid, "params", key).Do(&param)
				if err != nil {
					return generateError(err, "Failed to delete param %v: %v", key, uuid)
				}
			} else {
				var data interface{}
				if err := bufOrFileDecode(ref, &data); err != nil {
					return generateError(err, "Failed to parse ref %s: %v", o.singleName, err)
				}
				path := fmt.Sprintf("/%s", makeJsonPtr(key))
				patch := jsonpatch2.Patch{
					jsonpatch2.Operation{
						Op:    "test",
						Path:  path,
						Value: data,
					},
					jsonpatch2.Operation{
						Op:   "remove",
						Path: path,
					},
				}
				if err := session.Req().Patch(patch).UrlFor(o.name, uuid, "params").Do(&param); err != nil {
					return generateError(err, "Failed to fetch params %v: %v", o.singleName, uuid)
				}
			}
			return prettyPrint(param)
		},
	})
}

func (o *ops) bootenv() {
	o.addCommand(&cobra.Command{
		Use:   "bootenv [id] [bootenv]",
		Short: fmt.Sprintf("Set the %s's bootenv", o.singleName),
		Long:  fmt.Sprintf(`Helper function to update the %s's bootenv.`, o.singleName),
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
			ex := models.Clone(data).(models.BootEnver)
			ex.SetBootEnv(args[1])
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})
}

func (o *ops) profiles() {
	o.addCommand(&cobra.Command{
		Use:   "addprofile [id] [profile]",
		Short: fmt.Sprintf("Add profile to the machine's profile list"),
		Long:  `Helper function to add a profile to the machine's profile list.`,
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
			ex := models.Clone(data).(models.Profiler)
			ex.SetProfiles(append(ex.GetProfiles(), args[1]))
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})
	o.addCommand(&cobra.Command{
		Use:   "removeprofile [id] [profile]",
		Short: fmt.Sprintf("Remove a profile from the machine's list"),
		Long:  `Helper function to update the machine's profile list by removing one.`,
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
			ex := models.Clone(data).(models.Profiler)
			newProfiles := []string{}
			for _, s := range ex.GetProfiles() {
				if s == args[1] {
					continue
				}
				newProfiles = append(newProfiles, s)
			}
			ex.SetProfiles(newProfiles)
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})

}

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

func (o *ops) actions() {
	actionName := o.actionName
	if actionName == "" {
		actionName = "action"
	}
	actionsName := fmt.Sprintf("%ss", actionName)
	prefix := "system"
	idStr := ""
	argCount := 0
	evenCount := 1
	if o.example != nil {
		prefix = o.example().Prefix()
		idStr = " [id]"
		argCount = 1
		evenCount = 0
	}
	plugin := ""
	actions := &cobra.Command{
		Use:   fmt.Sprintf("%s%s", actionsName, idStr),
		Short: fmt.Sprintf("Display actions for this %s", o.singleName),
		Long:  fmt.Sprintf("Helper function to display the %s's actions.", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != argCount {
				return fmt.Errorf("%v requires %d argument", c.UseLine(), argCount)
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			res := []models.AvailableAction{}
			var req *api.R
			id := "system"
			if argCount == 1 {
				id = args[0]
				req = session.Req().UrlFor(prefix, id, actionsName)
			} else {
				req = session.Req().UrlFor(prefix, actionsName)
			}
			if plugin != "" {
				req = req.Params("plugin", plugin)
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch actions %v: %v", o.singleName, id)
			}
			return prettyPrint(res)
		},
	}
	actions.Flags().StringVar(&plugin, "plugin", "", "Plugin to filter action search")
	o.addCommand(actions)
	action := &cobra.Command{
		Use:   fmt.Sprintf("%s%s [action]", actionName, idStr),
		Short: fmt.Sprintf("Display the action for this %s", o.singleName),
		Long:  fmt.Sprintf("Helper function to display the %s's action.", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != argCount+1 {
				return fmt.Errorf("%v requires %d arguments", c.UseLine(), argCount+1)
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			action := args[argCount]
			res := &models.AvailableAction{}
			var req *api.R
			id := "system"
			if argCount == 1 {
				id = args[0]
				req = session.Req().UrlFor(prefix, id, actionsName, action)
			} else {
				req = session.Req().UrlFor(prefix, actionsName, action)
			}
			if plugin != "" {
				req = req.Params("plugin", plugin)
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch action %v: %v", o.singleName, id)
			}
			return prettyPrint(res)
		},
	}
	action.Flags().StringVar(&plugin, "plugin", "", "Plugin to filter action search")
	o.addCommand(action)
	actionParams := map[string]interface{}{}
	runaction := &cobra.Command{
		Use:   fmt.Sprintf("run%s%s [command] [- | JSON or YAML Map of objects | pairs of string objects]", actionName, idStr),
		Short: "Run action on object from plugin",
		Args: func(c *cobra.Command, args []string) error {
			actionParams = map[string]interface{}{}
			if len(args) == argCount+2 {
				if err := into(args[argCount+1], &actionParams); err != nil {
					return err
				}
				return nil
			}
			if len(args) >= argCount+1 && len(args)%2 == evenCount {
				for i := argCount + 1; i < len(args); i += 2 {
					var obj interface{}
					if err := api.DecodeYaml([]byte(args[i+1]), &obj); err != nil {
						return fmt.Errorf("Invalid parameters: %s %v\n", args[i+1], err)
					}
					actionParams[args[i]] = obj
				}
				return nil
			}
			if argCount == 1 {
				return fmt.Errorf("runaction either takes three arguments or a multiple of two, not %d", len(args))
			}
			return fmt.Errorf("runaction either takes two arguments or one plus a multiple of two, not %d", len(args))
		},

		RunE: func(c *cobra.Command, args []string) error {
			command := args[argCount]
			var resp interface{}
			var req *api.R
			if argCount == 1 {
				id := args[0]
				req = session.Req().Post(actionParams).UrlFor(prefix, id, actionsName, command)
			} else {
				req = session.Req().Post(actionParams).UrlFor(prefix, actionsName, command)
			}
			if plugin != "" {
				req = req.Params("plugin", plugin)
			}
			if err := req.Do(&resp); err != nil {
				return generateError(err, "Error running action")
			}
			return prettyPrint(resp)
		},
	}
	runaction.Flags().StringVar(&plugin, "plugin", "", "Plugin to filter action search")
	o.addCommand(runaction)
}

func (o *ops) command(app *cobra.Command) {
	res := &cobra.Command{
		Use:   o.name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", o.name),
	}
	if o.example != nil {
		ref := o.example()
		if _, ok := ref.(models.BootEnver); ok {
			o.bootenv()
		}
		if _, ok := ref.(models.Paramer); ok {
			o.params()
		}
		if _, ok := ref.(models.Profiler); ok {
			o.profiles()
		}
		if _, ok := ref.(models.Tasker); ok {
			o.tasks()
		}
		if _, ok := ref.(models.Actor); ok {
			o.actions()
		}
		res.AddCommand(o.commands()...)
	}
	app.AddCommand(res)
}
