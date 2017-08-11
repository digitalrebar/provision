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
	"sort"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision"
	apiclient "github.com/digitalrebar/provision/client"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

var (
	App = &cobra.Command{
		Use:   "drpcli",
		Short: "A CLI application for interacting with the DigitalRebar Provision API",
	}

	version   = provision.RS_VERSION
	debug     = false
	endpoint  = "https://127.0.0.1:8092"
	token     = ""
	username  = "rocketskates"
	password  = "r0cketsk8ts"
	format    = "json"
	session   *apiclient.DigitalRebarProvision
	basicAuth runtime.ClientAuthInfoWriter
	uf        func(*cobra.Command) error
	dumpUsage = true
	force     = false
	noPretty  = false
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
		"The Digital Rebar Provision API endpoint to talk to")
	App.PersistentFlags().StringVarP(&username,
		"username", "U", username,
		"Name of the Digital Rebar Provision user to talk to")
	App.PersistentFlags().StringVarP(&password,
		"password", "P", password,
		"password of the Digital Rebar Provision user")
	App.PersistentFlags().StringVarP(&token,
		"token", "T", token,
		"token of the Digital Rebar Provision access")
	App.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")
	App.PersistentFlags().StringVarP(&format,
		"format", "F", "json",
		`The serialzation we expect for output.  Can be "json" or "yaml"`)
	App.PersistentFlags().BoolVarP(&force,
		"force", "f", false,
		"When needed, attempt to force the operation - used on some update/patch calls")

	uf = App.UsageFunc()
	App.SetUsageFunc(MyUsage)

	App.PersistentPreRun = func(c *cobra.Command, a []string) {
		if session == nil {
			var err error
			d("Talking to Digital Rebar Provision with %v (%v:%v)", endpoint, username, password)
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
		Short: "Digital Rebar Provision CLI Command Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %v\n", version)
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "autocomplete <filename>",
		Short: "Digital Rebar Provision CLI Command Bash AutoCompletion File",
		Long:  "Generate a bash autocomplete file as <filename>.\nPlace the generated file in /etc/bash_completion.d or /usr/local/etc/bash_completion.d.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1  argument", cmd.UseLine())
			}
			App.GenBashCompletionFile(args[0])
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
	if noPretty {
		fmt.Printf("%v", o)
		return nil
	}
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

type CommonOps struct {
	Name         string
	SingularName string
}

func (co CommonOps) GetName() string {
	return co.Name
}

func (co CommonOps) GetSingularName() string {
	return co.SingularName
}

type Payloader interface {
	GetPayload() interface{}
}

type ICommonOps interface {
	GetName() string
	GetSingularName() string
}

type CommonTypeOps interface {
	GetType() interface{}
	GetId(interface{}) (string, error)
}

type ListOp interface {
	ICommonOps
	List(params map[string]string) (interface{}, error)
	GetIndexes() map[string]string
}

type GetOp interface {
	ICommonOps
	Get(string) (interface{}, error)
	GetIndexes() map[string]string
}

type CreateOps interface {
	ICommonOps
	CommonTypeOps
	Create(interface{}) (interface{}, error)
}

type ModOps interface {
	ICommonOps
	CommonTypeOps
}

type PatchOps interface {
	ModOps
	Patch(string, interface{}) (interface{}, error)
}

type UpdateOps interface {
	ModOps
	Update(string, interface{}) (interface{}, error)
}

type DeleteOps interface {
	ICommonOps
	Delete(string) (interface{}, error)
}

type UploadOp interface {
	ICommonOps
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

var listLimit = -1
var listOffset = -1

func Get(id string, gptrs GetOp) (interface{}, error) {
	return gptrs.Get(id)
}

func Create(input string, ptrs CreateOps) (interface{}, error) {
	var buf []byte
	var err error
	if input == "-" {
		buf, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("Error reading from stdin: %v", err)
		}
	} else {
		buf = []byte(input)
	}
	var obj interface{}
	obj = ptrs.GetType()
	err = yaml.Unmarshal(buf, obj)
	if err != nil {
		obj = ""
		err2 := yaml.Unmarshal(buf, &obj)
		if err2 != nil {
			return nil, fmt.Errorf("Invalid %v object: %v and %v", ptrs.GetSingularName(), err, err2)
		}
	}
	if data, err := ptrs.Create(obj); err != nil {
		return nil, generateError(err, "Unable to create new %v", ptrs.GetSingularName())
	} else {
		return data, nil
	}
}

func Update(id, input string, ptrs ModOps) (interface{}, error) {
	var buf []byte
	var err error
	if input == "-" {
		buf, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("Error reading from stdin: %v", err)
		}
	} else {
		buf = []byte(input)
	}
	var intermediate interface{}
	err = yaml.Unmarshal(buf, &intermediate)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
	}

	updateObj, err := json.Marshal(intermediate)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal input stream: %v\n", err)
	}
	data, err := Get(id, ptrs.(GetOp))
	if err != nil {
		return nil, generateError(err, "Failed to fetch %v: %v", ptrs.GetSingularName(), id)
	}
	baseObj, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal object: %v\n", err)
	}

	merged, err := safeMergeJSON(data, updateObj)
	if err != nil {
		return nil, fmt.Errorf("Unable to merge objects: %v\n", err)
	}

	// if the caller provides update, use it because we have Patch issues.
	if uptrs, ok := ptrs.(UpdateOps); ok {
		obj := ptrs.GetType()
		err = yaml.Unmarshal(merged, obj)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshal merged input stream: %v\n", err)
		}

		return uptrs.Update(id, obj)
	}
	// Else use Patch
	patch, err := jsonpatch2.Generate(baseObj, merged, true)
	if err != nil {
		return nil, fmt.Errorf("Error generating patch: %v", err)
	}
	p := models.Patch{}
	if err := utils.Remarshal(&patch, &p); err != nil {
		return nil, fmt.Errorf("Error translating patch: %v", err)
	}

	pptrs, _ := ptrs.(PatchOps)
	if data, err := pptrs.Patch(id, p); err != nil {
		return nil, generateError(err, "Unable to patch %v", id)
	} else {
		return data, nil
	}
}

func Delete(id string, ptrs DeleteOps) (interface{}, error) {
	return ptrs.Delete(id)
}

func commonOps(pobj interface{}) (commands []*cobra.Command) {
	commands = make([]*cobra.Command, 0, 0)
	updateAdded := false
	if ptrs, ok := pobj.(ListOp); ok {
		idxs := ptrs.GetIndexes()
		bigidxstr := ""
		if len(idxs) > 0 {
			keys := []string{}
			for k := range idxs {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			idxstr := ""
			idxsingle := "notallowed"
			for _, k := range keys {
				if k == "Key" {
					continue
				}
				idxsingle = k
				idxstr += fmt.Sprintf("*  %s = %s\n", k, idxs[k])
			}
			bigidxstr = fmt.Sprintf(`
You may specify:

*  Offset = integer, 0-based inclusive starting point in filter data.
*  Limit = integer, number of items to return

Functional Indexs:

%s

Functions:

*  Eq(value) = Return items that are equal to value
*  Lt(value) = Return items that are less than value
*  Lte(value) = Return items that less than or equal to value
*  Gt(value) = Return items that are greater than value
*  Gte(value) = Return items that greater than or equal to value
*  Between(lower,upper) = Return items that are inclusively between lower and upper
*  Except(lower,upper) = Return items that are not inclusively between lower and upper

Example:

*  %v=fred - returns items named fred
*  %v=Lt(fred) - returns items that alphabetically less than fred.
*  %v=Lt(fred)&Available=true - returns items with Name less than fred and Available is true

`, idxstr, idxsingle, idxsingle, idxsingle)
		}
		listCmd := &cobra.Command{
			Use:   "list [key=value] ...",
			Short: fmt.Sprintf("List all %v", ptrs.GetName()),
			Long:  fmt.Sprintf("This will list all %v by default.\n%s\n", ptrs.GetName(), bigidxstr),
			RunE: func(c *cobra.Command, args []string) error {
				dumpUsage = false

				parms := map[string]string{}

				for _, a := range args {
					ar := strings.SplitN(a, "=", 2)
					if len(ar) != 2 {
						return fmt.Errorf("Filter argument requires an '=' separator: %s", a)
					}
					parms[ar[0]] = ar[1]
				}

				if data, err := ptrs.List(parms); err != nil {
					return generateError(err, "listing %v", ptrs.GetName())
				} else {
					return prettyPrint(data)
				}
			},
		}
		if len(idxs) > 0 {
			listCmd.Flags().IntVar(&listLimit, "limit", -1, "Maximum number of items to return")
			listCmd.Flags().IntVar(&listOffset, "offset", -1, "Number of items to skip before starting to return data")
		}

		commands = append(commands, listCmd)

	}
	if gptrs, ok := pobj.(GetOp); ok {
		idxs := gptrs.GetIndexes()
		bigidxstr := ""
		if len(idxs) > 0 {
			keys := []string{}
			for k := range idxs {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			idxstr := ""
			idxsingle := "notallowed"
			for _, k := range keys {
				if k == "Key" {
					continue
				}
				idxsingle = k
				idxstr += fmt.Sprintf("*  %s = %s\n", k, idxs[k])
			}
			bigidxstr = fmt.Sprintf(`
You may specify the id in the request by the using normal key or by index.

Functional Indexs:

%s

When using the index name, use the following form:

* Index:Value

Example:

* e.g: %s:fred

`, idxstr, idxsingle)
		}
		commands = append(commands, &cobra.Command{
			Use:   "show [id]",
			Short: fmt.Sprintf("Show a single %v by id", gptrs.GetSingularName()),
			Long:  fmt.Sprintf("This will show a %v.\n%s\n", gptrs.GetName(), bigidxstr),
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				dumpUsage = false
				if data, err := Get(args[0], gptrs); err != nil {
					return generateError(err, "Failed to fetch %v: %v", gptrs.GetSingularName(), args[0])
				} else {
					return prettyPrint(data)
				}
			},
		})
		commands = append(commands, &cobra.Command{
			Use:   "exists [id]",
			Short: fmt.Sprintf("See if a %v exists by id", gptrs.GetSingularName()),
			Long:  fmt.Sprintf("This will detect if a %v exists.\n%s\n", gptrs.GetName(), bigidxstr),
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				dumpUsage = false
				if _, err := Get(args[0], gptrs); err != nil {
					return generateError(err, "Failed to fetch %v: %v", gptrs.GetSingularName(), args[0])
				}
				return nil
			},
		})

		if ptrs, ok := pobj.(CreateOps); ok {
			commands = append(commands, &cobra.Command{
				Use:   "create [json]",
				Short: fmt.Sprintf("Create a new %v with the passed-in JSON or string key", ptrs.GetSingularName()),
				Long: `
As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin.

In either case, for the Machine, BootEnv, User, and Profile objects, a string may be provided to create a new
empty object of that type.  For User, BootEnv, Machine, and Profile, it will be the object's name.
`,
				RunE: func(c *cobra.Command, args []string) error {
					if len(args) != 1 {
						return fmt.Errorf("%v requires 1 argument", c.UseLine())
					}
					dumpUsage = false

					if data, err := Create(args[0], ptrs); err != nil {
						return err
					} else {
						return prettyPrint(data)
					}
				},
			})
		}

		if ptrs, ok := pobj.(PatchOps); ok {
			updateAdded = true
			commands = append(commands, &cobra.Command{
				Use:   "update [id] [json]",
				Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", ptrs.GetSingularName()),
				Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
				RunE: func(c *cobra.Command, args []string) error {
					if len(args) != 2 {
						return fmt.Errorf("%v requires 2 arguments", c.UseLine())
					}
					dumpUsage = false

					if data, err := Update(args[0], args[1], ptrs); err != nil {
						return err
					} else {
						return prettyPrint(data)
					}
				},
			})

			commands = append(commands, &cobra.Command{
				Use:   "patch [objectJson] [changesJson]",
				Short: fmt.Sprintf("Patch %v with the passed-in JSON", ptrs.GetSingularName()),
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

		if !updateAdded {
			if ptrs, ok := pobj.(UpdateOps); ok {
				commands = append(commands, &cobra.Command{
					Use:   "update [id] [json]",
					Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", ptrs.GetSingularName()),
					Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
					RunE: func(c *cobra.Command, args []string) error {
						if len(args) != 2 {
							return fmt.Errorf("%v requires 2 arguments", c.UseLine())
						}
						dumpUsage = false

						if data, err := Update(args[0], args[1], ptrs); err != nil {
							return err
						} else {
							return prettyPrint(data)
						}
					},
				})
			}
		}
	}

	if ptrs, ok := pobj.(DeleteOps); ok {
		commands = append(commands, &cobra.Command{
			Use:   "destroy [id]",
			Short: fmt.Sprintf("Destroy %v by id", ptrs.GetSingularName()),
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("%v requires 1 argument", c.UseLine())
				}
				dumpUsage = false
				if _, err := Delete(args[0], ptrs); err != nil {
					return generateError(err, "Unable to destroy %v %v", ptrs.GetSingularName(), args[0])
				} else {
					fmt.Printf("Deleted %v %v\n", ptrs.GetSingularName(), args[0])
					return nil
				}
			},
		})
	}

	if ptrs, ok := pobj.(UploadOp); ok {
		commands = append(commands, &cobra.Command{
			Use:   "upload [file] as [name]",
			Short: "Upload a local file to Digital Rebar Provision",
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
