package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/VictorLowther/jsonpatch2"
	yaml "github.com/ghodss/yaml"
)

var (
	modelPrefixes = func() map[string]Model {
		res := map[string]Model{}
		for _, m := range All() {
			res[m.Prefix()] = m
		}
		return res
	}()
)

func copyMap(m map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		res[k] = v
	}
	return res
}

// swagger:model
type BlobInfo struct {
	Path string
	Size int64
}

type Model interface {
	Prefix() string
	Key() string
	KeyName() string
}

type Filler interface {
	Model
	Fill()
}

type Slicer interface {
	Filler
	SliceOf() interface{}
	ToModels(interface{}) []Model
}

func All() []Model {
	return []Model{
		&BootEnv{},
		&Interface{},
		&Job{},
		&Lease{},
		&Machine{},
		&Param{},
		&PluginProvider{},
		&Plugin{},
		&Pref{},
		&Profile{},
		&Reservation{},
		&Role{},
		&Stage{},
		&Subnet{},
		&Task{},
		&Template{},
		&User{},
		&Workflow{},
		&Tenant{},
	}
}

func AllPrefixes() []string {
	all := All()
	res := make([]string, len(all))
	for i := range all {
		res[i] = all[i].Prefix()
	}
	return res
}

func New(kind string) (Slicer, error) {
	for _, i := range All() {
		key := i.Prefix()
		if key == kind || kind == strings.TrimSuffix(key, "s") {
			res := i.(Slicer)
			res.Fill()
			return res, nil
		}
	}
	return nil, fmt.Errorf("No such Model: %s", kind)
}

func Clone(m Model) Model {
	if m == nil {
		return nil
	}
	res, err := New(m.Prefix())
	if err != nil {
		log.Panicf("Failed to make a new %s: %v", m.Prefix(), err)
	}
	buf := bytes.Buffer{}
	enc, dec := json.NewEncoder(&buf), json.NewDecoder(&buf)
	if err := enc.Encode(m); err != nil {
		log.Panicf("Failed to encode %s:%s: %v", m.Prefix(), m.Key(), err)
	}
	if err := dec.Decode(res); err != nil {
		log.Panicf("Failed to decode %s:%s: %v", m.Prefix(), m.Key(), err)
	}
	return res
}

var (
	validName      = regexp.MustCompile(`^\pL+([- _.]+|\pN+|\pL+)+$`)
	validParamName = regexp.MustCompile(`^\pL+([- _./]+|\pN+|\pL+)+$`)
)

func ValidName(msg, s string) error {
	if validName.MatchString(s) {
		return nil
	}
	return fmt.Errorf("%s `%s`", msg, s)
}

func ValidParamName(msg, s string) error {
	if validParamName.MatchString(s) {
		return nil
	}
	return fmt.Errorf("%s `%s`", msg, s)
}

type NameSetter interface {
	Model
	SetName(string)
}

type Paramer interface {
	Model
	GetParams() map[string]interface{}
	SetParams(map[string]interface{})
}

type Profiler interface {
	Model
	GetProfiles() []string
	SetProfiles([]string)
}

type BootEnver interface {
	Model
	GetBootEnv() string
	SetBootEnv(string)
}

type Tasker interface {
	Model
	GetTasks() []string
	SetTasks([]string)
}

type TaskRunner interface {
	Tasker
	RunningTask() int
}

// Only implement this if you want actions
type Actor interface {
	Model
	CanHaveActions() bool
}

func FibBackoff(thunk func() error) {
	timeouts := []time.Duration{
		time.Second,
		time.Second,
		2 * time.Second,
		3 * time.Second,
		5 * time.Second,
		8 * time.Second,
	}
	for _, d := range timeouts {
		if thunk() == nil {
			return
		}
		time.Sleep(d)
	}
}

// GenPatch generates a JSON patch that will transform source into target.
// The generated patch will have all the applicable test clauses.
func GenPatch(source, target interface{}, paranoid bool) (jsonpatch2.Patch, error) {
	srcBuf, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	tgtBuf, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}
	return jsonpatch2.GenerateFull(srcBuf, tgtBuf, true, paranoid)
}

// DecodeYaml is a helper function for dealing with user input -- when
// accepting input from the user, we want to treat both YAML and JSON
// as first-class citizens.  The YAML library we use makes that easier
// by using the json struct tags for all marshalling and unmarshalling
// purposes.
//
// Note that the REST API does not use YAML as a wire protocol, so
// this function should never be used to decode data coming from the
// provision service.
func DecodeYaml(buf []byte, ref interface{}) error {
	return yaml.Unmarshal(buf, ref)
}
