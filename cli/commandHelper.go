package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func (o *ops) commands() []*cobra.Command {
	canSlim := false
	canDecode := false
	canParam := false
	if _, ok := o.example().(models.MetaHaver); ok {
		canSlim = true
	}
	if _, ok := o.example().(models.Paramer); ok {
		canSlim = true
		canParam = true
		canDecode = true
	}
	slim := ""
	params := ""
	decode := false
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
* *index* Re *re2 compatible regular expression*
  This will return items in *index* that match the passed-in regular expression
  We use the regular expression syntax described at
  https://github.com/google/re2/wiki/Syntax
* *index* Between *lower* *upper*
  This will return items Greater Than Or Equal to *lower*
  and Less Than Or Equal to *upper* according to *index*
* *index* Except *lower* *upper*
  This will return items Less Than *lower* or
  Greater Than *upper* according to *index*
* *index* In *comma,separated,list,of,values*
  This will return any items In the set passed for the
  comma-separated list of values.
* *index* Nin *comma,separated,list,of,values*
  This will return any items Not In the set passed for the
  comma-separated list of values.

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
			req := Session.Req().List(o.name)
			if len(args) > 0 && strings.Contains(args[0], "=") {
				// Old-style structured args
				if listLimit != -1 {
					args = append(args, fmt.Sprintf("limit=%d", listLimit))
				}
				if listOffset != -1 {
					args = append(args, fmt.Sprintf("offset=%d", listOffset))
				}
				if slim != "" {
					args = append(args, fmt.Sprintf("slim=%s", slim))
				}
				if params != "" {
					args = append(args, fmt.Sprintf("params=%s", params))
				}
				if decode {
					args = append(args, "decode=true")
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
				if slim != "" {
					args = append(args, "slim", slim)
				}
				if params != "" {
					args = append(args, "params", params)
				}
				if decode {
					args = append(args, "decode")
				}
				if len(args) > 0 {
					req = Session.Req().Filter(o.name, args...)
				}
			}
			data := []interface{}{}
			err := req.Do(&data)
			if err != nil {
				return generateError(err, "listing %v", o.name)
			}
			return prettyPrint(data)
		},
	}
	cmds = append(cmds, listCmd)
	listCmd.Flags().IntVar(&listLimit, "limit", -1, "Maximum number of items to return")
	listCmd.Flags().IntVar(&listOffset, "offset", -1, "Number of items to skip before starting to return data")
	if canSlim {
		listCmd.Flags().StringVar(&slim,
			"slim",
			"",
			"Should elide certain fields.  Can be 'Params', 'Meta', or a comma-separated list of both.")
	}
	if canParam {
		listCmd.Flags().StringVar(&params,
			"params",
			"",
			"Should return only the parameters specified as a comma-separated list of parameter names.")
	}
	if canDecode {
		listCmd.Flags().BoolVar(&decode,
			"decode",
			false,
			"Should decode any secure params")
	}
	cmds = append(cmds, &cobra.Command{
		Use:   "indexes",
		Short: fmt.Sprintf("Get indexes for %s", o.name),
		Long:  fmt.Sprintf("Different object types can have indexes on various fields."),
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			indexes, err := Session.Indexes(o.name)
			if err != nil {
				return generateError(err, "Error fetching indexes")
			}
			return prettyPrint(indexes)
		},
	})
	showCmd := &cobra.Command{
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
			req := Session.Req().UrlFor(o.name, args[0])
			if slim != "" {
				req = req.Params("slim", slim)
			}
			if decode {
				req = req.Params("decode", "true")
			}
			if params != "" {
				req = req.Params("params", params)
			}
			if err := req.Do(&data); err != nil {
				return generateError(err, "Failed to fetch %v: %v", o.singleName, args[0])
			}
			return prettyPrint(data)
		},
	}
	if canSlim {
		showCmd.Flags().StringVar(&slim,
			"slim",
			"",
			"Should elide certain fields.  Can be 'Params', 'Meta', or a comma-separated list of both.")
	}
	if canParam {
		showCmd.Flags().StringVar(&params,
			"params",
			"",
			"Should return only the parameters specified as a comma-separated list of parameter names.")
	}
	if canDecode {
		showCmd.Flags().BoolVar(&decode,
			"decode",
			false,
			"Should decode any secure params.")
	}

	cmds = append(cmds, showCmd)
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
			exists, err := Session.ExistsModel(o.name, args[0])
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
				req := Session.Req().Post(ref).UrlFor(ref.Prefix())
				if force {
					req.Params("force", "true")
				}
				if err := req.Do(&ref); err != nil {
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
					req := Session.Req().Put(toPut).UrlForM(toPut)
					if force {
						req.Params("force", "true")
					}
					if err := req.Do(&toPut); err != nil {
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
					res := models.Clone(refObj)
					req := Session.Req()
					if ref != "" {
						req = req.ParanoidPatch()
					}
					req = req.PatchTo(refObj, toPut)
					if force {
						req = req.Params("force", "true")
					}
					if err := req.Do(&res); err != nil {
						return generateError(err, "Unable to update %v", args[0])
					}
					return prettyPrint(res)
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
				_, err := Session.DeleteModel(o.name, args[0])
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
				svalue := args[2]
				timeout := time.Hour * 24
				if len(args) == 4 {
					t, e := strconv.ParseInt(args[3], 10, 64)
					if e != nil {
						return e
					}
					timeout = time.Second * time.Duration(t)
				}

				var value interface{}
				if err := json.Unmarshal([]byte(svalue), &value); err != nil {
					value = svalue
				}

				testfn := api.EqualItem(field, value)
				item, err := o.refOrFill(id)
				if err != nil {
					return err
				}
				es, err := Session.Events()
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
