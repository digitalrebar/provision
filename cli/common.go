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
	"github.com/digitalrebar/provision"
	apiclient "github.com/digitalrebar/provision/client"
	"github.com/digitalrebar/provision/models"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

var (
	App = &cobra.Command{
		Use:   "rscli",
		Short: "A CLI application for interacting with the DigitalRebar Provision API",
	}

	version            = provision.RS_VERSION
	debug              = false
	endpoint           = "https://127.0.0.1:8092"
	token              = ""
	username, password string
	format             = "json"
	session            *apiclient.DigitalRebarProvision
	basicAuth          runtime.ClientAuthInfoWriter
	uf                 func(*cobra.Command) error
	dumpUsage          = true
)

func MyUsage(c *cobra.Command) error {
	if dumpUsage {
		return uf(c)
	}
	return nil
}

func init() {
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		endpoint = ep
	}
	if tk := os.Getenv("RS_TOKEN"); tk != "" {
		token = tk
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
	App.PersistentFlags().StringVarP(&token,
		"token", "T", token,
		"token of the Rocket-Skates access")
	App.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")
	App.PersistentFlags().StringVarP(&format,
		"format", "F", "json",
		`The serialzation we expect for output.  Can be "json" or "yaml"`)

	uf = App.UsageFunc()
	App.SetUsageFunc(MyUsage)

	App.PersistentPreRun = func(c *cobra.Command, a []string) {
		if session == nil {
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
		}
		if token != "" {
			basicAuth = httptransport.BearerToken(token)
		} else {
			basicAuth = httptransport.BasicAuth(username, password)
		}
	}
	App.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Rocket-Skates CLI Command Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %v\n", version)
			return nil
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

func prettyPrint(o interface{}) (err error) {
	var buf []byte
	switch format {
	case "json":
		buf, err = json.MarshalIndent(o, "", "  ")
	case "yaml":
		buf, err = yaml.Marshal(o)
	default:
		return fmt.Errorf("Unknown pretty format %s", format)
	}
	if err != nil {
		return fmt.Errorf("Failed to unmarshal returned object! %s", err.Error())
	}
	fmt.Println(string(buf))
	return nil
}

type Payloader interface {
	GetPayload() interface{}
}

type ListOp interface {
	List() (interface{}, error)
}

type GetOp interface {
	Get(string) (interface{}, error)
}

type ModOps interface {
	GetType() interface{}
	GetId(interface{}) (string, error)
	Create(interface{}) (interface{}, error)
	Patch(string, interface{}) (interface{}, error)
}

type DeleteOps interface {
	Delete(string) (interface{}, error)
}

type UploadOp interface {
	Upload(string, *os.File) (interface{}, error)
}

func generateError(err error, sfmt string, args ...interface{}) error {
	s := fmt.Sprintf(sfmt, args...)

	obj, ok := err.(Payloader)
	if !ok {
		return fmt.Errorf(s+": %v", err)
	}

	e := obj.GetPayload()
	if e == nil {
		return fmt.Errorf(s+": %v", err)
	}

	ee, ok := e.(*models.Error)
	if !ok {
		return fmt.Errorf(s+": %v", err)
	}

	s = ""
	first := true
	for _, ns := range ee.Messages {
		if !first {
			s = s + "\n"
		}
		first = false
		s = s + ns
	}
	return fmt.Errorf(s)
}

// TODO: Consider adding Match someday

func commonOps(singularName, name string, pobj interface{}) (commands []*cobra.Command) {
	commands = make([]*cobra.Command, 0, 0)

	if ptrs, ok := pobj.(ListOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "list",
			Short: fmt.Sprintf("List all %v", name),
			RunE: func(c *cobra.Command, args []string) error {
				dumpUsage = false
				if data, err := ptrs.List(); err != nil {
					return generateError(err, "Error listing %v", name)
				} else {
					return prettyPrint(data)
				}
			},
		})
	}
	if gptrs, ok := pobj.(GetOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "show [id]",
			Short: fmt.Sprintf("Show a single %v by id", singularName),
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				dumpUsage = false
				if data, err := gptrs.Get(args[0]); err != nil {
					return generateError(err, "Failed to fetch %v: %v", singularName, args[0])
				} else {
					return prettyPrint(data)
				}
			},
		})
		commands = append(commands, &cobra.Command{
			Use:   "exists [id]",
			Short: fmt.Sprintf("See if a %v exists by id", singularName),
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				dumpUsage = false
				if _, err := gptrs.Get(args[0]); err != nil {
					return generateError(err, "Failed to fetch %v: %v", singularName, args[0])
				}
				return nil
			},
		})

		if ptrs, ok := pobj.(ModOps); ok {
			commands = append(commands, &cobra.Command{
				Use:   "create [json]",
				Short: fmt.Sprintf("Create a new %v with the passed-in JSON", singularName),
				Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
				RunE: func(c *cobra.Command, args []string) error {
					if len(args) != 1 {
						return fmt.Errorf("%v requires 1 argument", c.UseLine())
					}
					dumpUsage = false
					var buf []byte
					var err error
					if args[0] == "-" {
						buf, err = ioutil.ReadAll(os.Stdin)
						if err != nil {
							return fmt.Errorf("Error reading from stdin: %v", err)
						}
					} else {
						buf = []byte(args[0])
					}
					obj := ptrs.GetType()
					err = yaml.Unmarshal(buf, obj)
					if err != nil {
						return fmt.Errorf("Invalid %v object: %v", singularName, err)
					}
					if data, err := ptrs.Create(obj); err != nil {
						return generateError(err, "Unable to create new %v", singularName)
					} else {
						return prettyPrint(data)
					}
				},
			})

			commands = append(commands, &cobra.Command{
				Use:   "update [id] [json]",
				Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", singularName),
				Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
				RunE: func(c *cobra.Command, args []string) error {
					if len(args) != 2 {
						return fmt.Errorf("%v requires 2 arguments", c.UseLine())
					}
					dumpUsage = false
					data, err := gptrs.Get(args[0])
					if err != nil {
						return generateError(err, "Failed to fetch %v: %v", singularName, args[0])
					}
					var buf []byte

					baseObj, err := json.Marshal(data)
					if err != nil {
						return fmt.Errorf("Unable to marshal object: %v\n", err)
					}
					if args[1] == "-" {
						buf, err = ioutil.ReadAll(os.Stdin)
						if err != nil {
							return fmt.Errorf("Error reading from stdin: %v", err)
						}
					} else {
						buf = []byte(args[1])
					}
					var intermediate interface{}
					err = yaml.Unmarshal(buf, &intermediate)
					if err != nil {
						return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
					}
					updateObj, err := json.Marshal(intermediate)
					if err != nil {
						return fmt.Errorf("Unable to marshal input stream: %v\n", err)
					}

					merged, err := safeMergeJSON(data, updateObj)
					if err != nil {
						return fmt.Errorf("Unable to merge objects: %v\n", err)
					}
					patch, err := jsonpatch2.Generate(baseObj, merged, true)
					if err != nil {
						return fmt.Errorf("Error generating patch: %v", err)
					}
					p := models.Patch{}
					if err := utils.Remarshal(&patch, &p); err != nil {
						return fmt.Errorf("Error translating patch: %v", err)
					}

					if data, err := ptrs.Patch(args[0], p); err != nil {
						return generateError(err, "Unable to patch %v", args[0])
					} else {
						return prettyPrint(data)
					}
				},
			})

			commands = append(commands, &cobra.Command{
				Use:   "patch [objectJson] [changesJson]",
				Short: fmt.Sprintf("Patch %v with the passed-in JSON", singularName),
				RunE: func(c *cobra.Command, args []string) error {
					if len(args) != 2 {
						return fmt.Errorf("%v requires 2 arguments", c.UseLine())
					}
					dumpUsage = false
					obj := ptrs.GetType()
					if err := yaml.Unmarshal([]byte(args[0]), obj); err != nil {
						return fmt.Errorf("Unable to parse %v JSON %v\nError: %v", c.UseLine(), args[0], err)
					}
					newObj := ptrs.GetType()
					yaml.Unmarshal([]byte(args[0]), newObj)
					if err := yaml.Unmarshal([]byte(args[1]), newObj); err != nil {
						return fmt.Errorf("Unable to parse %v JSON %v\nError: %v", c.UseLine(), args[1], err)
					}
					newBuf, _ := json.Marshal(newObj)
					patch, err := jsonpatch2.Generate([]byte(args[0]), newBuf, true)
					if err != nil {
						return fmt.Errorf("Cannot generate JSON Patch\n%v", err)
					}
					p := models.Patch{}
					if err := utils.Remarshal(&patch, &p); err != nil {
						return fmt.Errorf("Error translating patch: %v", err)
					}

					id, err := ptrs.GetId(obj)
					if err != nil {
						return fmt.Errorf("Cannot get key for obj: %v", err)
					}
					if data, err := ptrs.Patch(id, p); err != nil {
						return generateError(err, "Unable to patch %v", args[0])
					} else {
						return prettyPrint(data)
					}
				},
			})
		}
	}

	if ptrs, ok := pobj.(DeleteOps); ok {
		commands = append(commands, &cobra.Command{
			Use:   "destroy [id]",
			Short: fmt.Sprintf("Destroy %v by id", singularName),
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				dumpUsage = false
				if _, err := ptrs.Delete(args[0]); err != nil {
					return generateError(err, "Unable to destroy %v %v", singularName, args[0])
				} else {
					fmt.Printf("Deleted %v %v\n", singularName, args[0])
					return nil
				}
			},
		})
	}

	if ptrs, ok := pobj.(UploadOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "upload [file] as [name]",
			Short: "Upload a local file to RocketSkates",
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 3 {
					return fmt.Errorf("Wrong number of args: expected 3, got %d", len(args))
				}
				dumpUsage = false
				f, err := os.Open(args[0])
				if err != nil {
					return fmt.Errorf("Failed to open %s: %v", args[0], err)
				}
				defer f.Close()
				if d, err := ptrs.Upload(args[2], f); err != nil {
					return generateError(err, "Error uploading: %v", args[0])
				} else {
					return prettyPrint(d)
				}
			},
		})
	}

	return
}
