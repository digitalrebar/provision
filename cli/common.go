package cli

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	apiclient "github.com/rackn/rocket-skates/client"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

var (
	App = &cobra.Command{
		Use:   "rscli",
		Short: "A CLI application for interacting with the Rocket-Skates API",
	}

	version            = "1.1.1"
	debug              = false
	endpoint           = "https://127.0.0.1:8092"
	username, password string
	format             = "json"
	session            *apiclient.RocketSkates
	basicAuth          runtime.ClientAuthInfoWriter
)

func init() {
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		endpoint = ep
	}
	if kv := os.Getenv("RS_KEY"); kv != "" {
		key := strings.SplitN(kv, ":", 2)
		if len(key) < 2 {
			log.Fatal("RS_KEY does not contain a username:password pair!")
		}
		if key[0] == "" || key[1] == "" {
			log.Fatal("RS_KEY contains an invalid username:password pair!")
		}
		username = key[0]
		password = key[1]
	}
	App.PersistentFlags().StringVarP(&endpoint,
		"endpoint", "E", endpoint,
		"The Rocket-Skates API endpoint to talk to")
	App.PersistentFlags().StringVarP(&username,
		"username", "U", username,
		"Name of the Rocket-Skates user to talk to")
	App.PersistentFlags().StringVarP(&password,
		"password", "P", password,
		"password of the Rocket-Skates user")
	App.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")
	App.PersistentFlags().StringVarP(&format,
		"format", "F", "json",
		`The serialzation we expect for output.  Can be "json" or "yaml"`)

	App.PersistentPreRun = func(c *cobra.Command, a []string) {
		var err error
		d("Talking to Rocket-Skates with %v (%v:%v)", endpoint, username, password)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc := &http.Client{Transport: tr}
		epURL, err := url.Parse(endpoint)
		if err != nil {
			log.Fatalf("Error handling endpoint %s: %v", endpoint, err)
		}
		transport := httptransport.NewWithClient(epURL.Host, "/api/v3", []string{epURL.Scheme}, hc)
		session = apiclient.New(transport, strfmt.Default)
		basicAuth = httptransport.BasicAuth(username, password)

		if err != nil {
			if c.Use != "version" {
				log.Fatalf("Could not connect to Rocket-Skates: %v\n", err.Error())
			}
		}
	}
	App.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Rocket-Skates CLI Command Version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Version: %v", version)
		},
	})

}

func safeMergeJSON(src interface{}, toMerge []byte) ([]byte, error) {
	toMergeObj := make(map[string]interface{})
	if err := json.Unmarshal(toMerge, &toMergeObj); err != nil {
		return nil, err
	}
	buf, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	var targetObj map[string]interface{}
	if err := json.Unmarshal(buf, &targetObj); err != nil {
		return nil, err
	}
	outObj, ok := utils.Merge(targetObj, toMergeObj).(map[string]interface{})
	if !ok {
		return nil, errors.New("Cannot happen in safeMergeJSON")
	}
	sv := reflect.ValueOf(src)
	for sv.Kind() == reflect.Ptr || sv.Kind() == reflect.Interface {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		log.Panicf("first arg to safeMergeJSON is not a struct! %#v", src)
	}
	finalObj := map[string]interface{}{}
	for i := 0; i < sv.NumField(); i++ {
		vf := sv.Field(i)
		if !vf.CanSet() {
			continue
		}
		tf := sv.Type().Field(i)
		mapField := tf.Name
		if tag, ok := tf.Tag.Lookup(`json`); ok {
			tagVals := strings.Split(tag, `,`)
			if tagVals[0] == "-" {
				continue
			}
			if tagVals[0] != "" {
				mapField = tagVals[0]
			}
		}
		if v, ok := outObj[mapField]; ok {
			finalObj[mapField] = v
		}
	}
	return json.Marshal(finalObj)
}

func d(msg string, args ...interface{}) {
	if debug {
		log.Printf(msg, args...)
	}
}

func pretty(o interface{}) (res string) {
	var buf []byte
	var err error
	switch format {
	case "json":
		buf, err = json.MarshalIndent(o, "", "  ")
	case "yaml":
		buf, err = yaml.Marshal(o)
	default:
		log.Fatalf("Unknown pretty format %s", format)
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
					log.Println(pretty(data))
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
					c.Printf("%v requires 1 argument\n", c.UseLine())
					return
				}
				if data, err := gptrs.Get(args[0]); err != nil {
					c.Printf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
					return
				} else {
					log.Println(pretty(data))
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
						log.Println(pretty(data))
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
					data, err := gptrs.Get(args[0])
					if err != nil {
						log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
					}
					var buf []byte

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

					merged, err := safeMergeJSON(data, updateObj)
					if err != nil {
						log.Fatalf("Unable to merge objects: %v\n", err)
					}
					patch, err := jsonpatch2.Generate(baseObj, merged, true)
					if err != nil {
						log.Fatalf("Error generating patch: %v", err)
					}
					p := models.Patch{}
					if err := utils.Remarshal(&patch, &p); err != nil {
						log.Fatalf("Error translating patch: %v", err)
					}

					if data, err := ptrs.Patch(args[0], p); err != nil {
						log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
					} else {
						fmt.Println(pretty(data))
					}
				},
			})

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
					patch, err := jsonpatch2.Generate([]byte(args[0]), newBuf, true)
					if err != nil {
						log.Fatalf("Cannot generate JSON Patch\n%v\n", err)
					}
					p := models.Patch{}
					if err := utils.Remarshal(&patch, &p); err != nil {
						log.Fatalf("Error translating patch: %v", err)
					}

					if data, err := ptrs.Patch("id", p); err != nil {
						log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
					} else {
						log.Println(pretty(data))
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
						log.Printf("Deleted %v %v\n", singularName, args[0])
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
					log.Println(pretty(d))
				}
			},
		})
	}

	return
}
