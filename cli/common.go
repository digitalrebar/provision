package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"

	"github.com/VictorLowther/jsonpatch"
	"github.com/VictorLowther/jsonpatch/utils"
	"github.com/go-openapi/runtime"
	apiclient "github.com/rackn/rocket-skates/client"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

var (
	Version            = "1.1.1"
	Debug              = false
	Endpoint           = "https://127.0.0.1:8092"
	Username, Password string
	Format             = "json"
	App                = &cobra.Command{
		Use:   "rscli",
		Short: "A CLI application for interacting with the Rocket-Skates API",
	}
	Session   *apiclient.RocketSkates
	BasicAuth runtime.ClientAuthInfoWriter
)

func safeMergeJSON(target, toMerge []byte) ([]byte, error) {
	targetObj := make(map[string]interface{})
	toMergeObj := make(map[string]interface{})
	if err := json.Unmarshal(target, &targetObj); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(toMerge, &toMergeObj); err != nil {
		return nil, err
	}
	outObj, ok := utils.Merge(targetObj, toMergeObj).(map[string]interface{})
	if !ok {
		return nil, errors.New("Cannot happen in safeMergeJSON")
	}
	keys := make([]string, 0)
	for k := range outObj {
		if _, ok := targetObj[k]; !ok {
			keys = append(keys, k)
		}
	}
	for _, k := range keys {
		delete(outObj, k)
	}
	return json.Marshal(outObj)
}

func D(msg string, args ...interface{}) {
	d(msg, args)
}

func d(msg string, args ...interface{}) {
	if Debug {
		log.Printf(msg, args...)
	}
}

func pretty(o interface{}) (res string) {
	var buf []byte
	var err error
	switch Format {
	case "json":
		buf, err = json.MarshalIndent(o, "", "  ")
	case "yaml":
		buf, err = yaml.Marshal(o)
	default:
		log.Fatalf("Unknown pretty format %s", Format)
	}
	if err != nil {
		log.Fatalf("Failed to unmarshal returned object!")
	}
	return string(buf)
}

type ListOp interface {
	List() (interface{}, error)
}

type GetOp interface {
	Get(string) (interface{}, error)
}

type ModOps interface {
	GetType() interface{}
	Create(interface{}) (interface{}, error)
	Put(string, interface{}) (interface{}, error)
	Patch(string, interface{}) (interface{}, error)
	Delete(string) (interface{}, error)
}

type UploadOp interface {
	Upload(string, *os.File) (interface{}, error)
}

// TODO: Consider adding Match someday
// TODO: Make PATCH WORK!!

func commonOps(singularName, name string, pobj interface{}) (commands []*cobra.Command) {
	commands = make([]*cobra.Command, 0, 0)

	if ptrs, ok := pobj.(ListOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "list",
			Short: fmt.Sprintf("List all %v", name),
			Run: func(c *cobra.Command, args []string) {
				if data, err := ptrs.List(); err != nil {
					log.Fatalf("Error listing %v: %v", name, err)
				} else {
					fmt.Println(pretty(data))
				}
			},
		})
	}
	if gptrs, ok := pobj.(GetOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "show [id]",
			Short: fmt.Sprintf("Show a single %v by id", singularName),
			Run: func(c *cobra.Command, args []string) {
				if len(args) != 1 {
					log.Fatalf("%v requires 1 argument\n", c.UseLine())
				}
				if data, err := gptrs.Get(args[0]); err != nil {
					log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
				} else {
					fmt.Println(pretty(data))
				}
			},
		})
		commands = append(commands, &cobra.Command{
			Use:   "exists [id]",
			Short: fmt.Sprintf("See if a %v exists by id", singularName),
			Run: func(c *cobra.Command, args []string) {
				if len(args) != 1 {
					log.Fatalf("%v requires 1 argument\n", c.UseLine())
				}
				if _, err := gptrs.Get(args[0]); err != nil {
					log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
				} else {
					os.Exit(0)
				}
			},
		})

		if ptrs, ok := pobj.(ModOps); ok {
			commands = append(commands, &cobra.Command{
				Use:   "create [json]",
				Short: fmt.Sprintf("Create a new %v with the passed-in JSON", singularName),
				Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
				Run: func(c *cobra.Command, args []string) {
					if len(args) != 1 {
						log.Fatalf("%v requires 1 argument\n", c.UseLine())
					}
					var buf []byte
					var err error
					if args[0] == "-" {
						buf, err = ioutil.ReadAll(os.Stdin)
						if err != nil {
							log.Fatalf("Error reading from stdin: %v", err)
						}
					} else {
						buf = []byte(args[0])
					}
					obj := ptrs.GetType()
					err = yaml.Unmarshal(buf, obj)
					if err != nil {
						log.Fatalf("Invalid %v object: %v\n", singularName, err)
					}
					if data, err := ptrs.Create(obj); err != nil {
						log.Fatalf("Unable to create new %v: %v\n", singularName, err)
					} else {
						fmt.Println(pretty(data))
					}
				},
			})

			commands = append(commands, &cobra.Command{
				Use:   "update [id] [json]",
				Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", singularName),
				Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
				Run: func(c *cobra.Command, args []string) {
					if len(args) != 2 {
						log.Fatalf("%v requires 2 arguments\n", c.UseLine())
					}
					if data, err := gptrs.Get(args[0]); err != nil {
						log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
					} else {
						var buf []byte
						var err error

						baseObj, err := json.Marshal(data)
						if err != nil {
							log.Fatalf("Unable to marshal object: %v\n", err)
						}

						if args[1] == "-" {
							buf, err = ioutil.ReadAll(os.Stdin)
							if err != nil {
								log.Fatalf("Error reading from stdin: %v", err)
							}
						} else {
							buf = []byte(args[1])
						}
						var intermediate interface{}
						err = yaml.Unmarshal(buf, &intermediate)
						if err != nil {
							log.Fatalf("Unable to unmarshal input stream: %v\n", err)
						}
						updateObj, err := json.Marshal(intermediate)
						if err != nil {
							log.Fatalf("Unable to marshal input stream: %v\n", err)
						}

						merged, err := safeMergeJSON(baseObj, updateObj)
						if err != nil {
							log.Fatalf("Unable to merge objects: %v\n", err)
						}

						obj := ptrs.GetType()
						err = yaml.Unmarshal(merged, obj)
						if err != nil {
							log.Fatalf("Unable to unmarshal merged object: %v\n", err)
						}

						if data, err := ptrs.Put(args[0], obj); err != nil {
							log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
						} else {
							fmt.Println(pretty(data))
						}
					}
				},
			})

			// GREG: This needs more help on arg processing.
			commands = append(commands, &cobra.Command{
				Use:   "patch [objectJson] [changesJson]",
				Short: fmt.Sprintf("Patch %v with the passed-in JSON", singularName),
				Run: func(c *cobra.Command, args []string) {
					if len(args) != 2 {
						log.Fatalf("%v requires 2 arguments\n", c.UseLine())
					}
					obj := &models.BootEnv{}
					if err := yaml.Unmarshal([]byte(args[0]), obj); err != nil {
						log.Fatalf("Unable to parse %v JSON %v\nError: %v\n", c.UseLine(), args[0], err)
					}
					newObj := &models.BootEnv{}
					yaml.Unmarshal([]byte(args[0]), newObj)
					if err := yaml.Unmarshal([]byte(args[1]), newObj); err != nil {
						log.Fatalf("Unable to parse %v JSON %v\nError: %v\n", c.UseLine(), args[1], err)
					}
					newBuf, _ := yaml.Marshal(newObj)
					patch, err := jsonpatch.GenerateJSON([]byte(args[0]), newBuf, true)
					if err != nil {
						log.Fatalf("Cannot generate JSON Patch\n%v\n", err)
					}
					p := models.Patch{}
					err = yaml.Unmarshal(patch, &p)
					if err != nil {
						log.Fatalf("Cannot generate JSON Patch Object\n%v\n", err)
					}
					if data, err := ptrs.Patch("id", p); err != nil {
						log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
					} else {
						fmt.Println(pretty(data))
					}
				},
			})

			commands = append(commands, &cobra.Command{
				Use:   "destroy [id]",
				Short: fmt.Sprintf("Destroy %v by id", singularName),
				Run: func(c *cobra.Command, args []string) {
					if len(args) != 1 {
						log.Fatalf("%v requires 1 argument\n", c.UseLine())
					}
					if _, err := ptrs.Delete(args[0]); err != nil {
						log.Fatalf("Unable to destroy %v %v\nError: %v\n", singularName, args[0], err)
					} else {
						fmt.Printf("Deleted %v %v\n", singularName, args[0])
					}
				},
			})
		}
	}

	if ptrs, ok := pobj.(UploadOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "upload [file] as [name]",
			Short: "Upload a local file to RocketSkates",
			Run: func(c *cobra.Command, args []string) {
				if len(args) != 3 {
					log.Fatalf("Wrong number of args: expected 3, got %d", len(args))
				}
				f, err := os.Open(args[0])
				if err != nil {
					log.Fatalf("Failed to open %s: %v", args[0], err)
				}
				defer f.Close()
				if d, err := ptrs.Upload(args[2], f); err != nil {
					log.Fatalf("Error uploading: %v", err)
				} else {
					fmt.Println(pretty(d))
				}
			},
		})
	}

	return
}
