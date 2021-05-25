package models

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mohae/deepcopy"

	"github.com/VictorLowther/jsonpatch2"
	yaml "github.com/ghodss/yaml"
)

var (
	baseModels = func() map[string]reflect.Type {
		res := map[string]reflect.Type{}
		for _, m := range All() {
			vv := reflect.ValueOf(m)
			for vv.Kind() == reflect.Interface || vv.Kind() == reflect.Ptr {
				vv = vv.Elem()
			}
			res[m.Prefix()] = vv.Type()
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

// BlobInfo contains information on an uploaded file or ISO.
// swagger:model
type BlobInfo struct {
	Path string
	Size int64
}

// Model is the interface that pretty much all non-Error objects
// returned by the API satisfy.
type Model interface {
	Prefix() string
	Key() string
	KeyName() string
}

// Filler interface defines if a model can be filled/initialized
type Filler interface {
	Model
	Fill()
}

// Slicer interface defines is a model can be operated on as slices
type Slicer interface {
	Filler
	SliceOf() interface{}
	ToModels(interface{}) []Model
}

// All returns a slice containing a single blank instance of all the
// Models.
func All() []Model {
	return []Model{
		&BootEnv{},
		&Interface{},
		&Job{},
		&Lease{},
		&Context{},
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
		&Pool{},
		&Endpoint{},
		&VersionSet{},
		&CatalogItem{},
		&AsyncAction{},
		&AsyncActionTemplate{},
		&AsyncActionCron{},
	}
}

// AllPrefixes returns a slice containing the prefix names of all the
// Models.
func AllPrefixes() []string {
	all := All()
	res := make([]string, len(all))
	for i := range all {
		res[i] = all[i].Prefix()
	}
	return res
}

// New returns a new blank instance of the Model with the passed-in
// prefix.
func New(prefix string) (Slicer, error) {
	var res Slicer
	if v, ok := baseModels[prefix]; ok {
		res = reflect.New(v).Interface().(Slicer)
	} else if v, ok = baseModels[strings.TrimSuffix(prefix, "s")]; ok {
		res = reflect.New(v).Interface().(Slicer)
	} else {
		res = &RawModel{"Type": prefix}
	}
	res.Fill()
	return res, nil
}

// Clone returns a deep copy of the passed-in Model
func Clone(m Model) Model {
	if m == nil {
		return nil
	}
	return deepcopy.Copy(m).(Model)
}

var (
	validMachineName  = regexp.MustCompile(`^(\pL|\pN)+([- _.]+|\pN+|\pL+)+$`)
	validEndpointName = regexp.MustCompile(`^(\pL|\pN)+([- _.:]+|\pN+|\pL+)+$`)
	validName         = regexp.MustCompile(`^\pL+([- _.]+|\pN+|\pL+)+$`)
	validUserName     = regexp.MustCompile(`^\pL+([- @_.]+|\pN+|\pL+)+$`)
	validParamName    = regexp.MustCompile(`^\pL+([- _./]+|\pN+|\pL+)+$`)
	validNumberName   = validMachineName
)

func validMatch(msg, s string, re *regexp.Regexp) error {
	if re.MatchString(s) {
		return nil
	}
	return fmt.Errorf("%s `%s`", msg, s)
}

// ValidMachineName validates that the string is a valid Machine Name
func ValidMachineName(msg, s string) error {
	return validMatch(msg, s, validMachineName)
}

// ValidEndpointName validates that the string is a valid Endpoint Name
func ValidEndpointName(msg, s string) error {
	return validMatch(msg, s, validEndpointName)
}

func ValidNumberName(msg, s string) error {
	return validMatch(msg, s, validNumberName)
}

// ValidName validates that the string is a valid Name
func ValidName(msg, s string) error {
	return validMatch(msg, s, validName)
}

// ValidUserName validates that the string is a valid Username
func ValidUserName(msg, s string) error {
	return validMatch(msg, s, validUserName)
}

// ValidParamName validates that the string is a valid Param Name
func ValidParamName(msg, s string) error {
	return validMatch(msg, s, validParamName)
}

// NameSetter interface if the model can change names
type NameSetter interface {
	Model
	SetName(string)
}

// Paramer interface defines if the model has parameters
type Paramer interface {
	Model
	GetParams() map[string]interface{}
	SetParams(map[string]interface{})
}

// Profiler interface defines if the model has profiles
type Profiler interface {
	Model
	GetProfiles() []string
	SetProfiles([]string)
}

// BootEnver interface defines if the model has a boot env
type BootEnver interface {
	Model
	GetBootEnv() string
	SetBootEnv(string)
}

// Tasker interface defines if the model has a task list
type Tasker interface {
	Model
	GetTasks() []string
	SetTasks([]string)
}

// TaskRunner interface defines if the object can run tasks
type TaskRunner interface {
	Tasker
	RunningTask() int
}

// Docer interface defines if the object has a documentation field
type Docer interface {
	Model
	GetDocumentation() string
}

// Descer interface defines if the object has a description field
type Descer interface {
	Model
	GetDescription() string
}

// Actor interface should be implemented this if you want actions
type Actor interface {
	Model
	CanHaveActions() bool
}

// FibBackoff takes function and retries it in a fibonacci backup sequence
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

// Remarshal remarshals src onto dest.
func Remarshal(src, dest interface{}) error {
	buf, err := json.Marshal(src)
	if err == nil {
		err = json.Unmarshal(buf, dest)
	}
	return err
}

// RandString returns a random string of n characters
// The range of characters is limited to the base64 set
func RandString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("Failed to read random\n")
		return "ARGH!"
	}
	base64 := base64.URLEncoding.EncodeToString(b)
	return base64[:n]
}

var (
	unitMap = map[string]time.Duration{
		"ns": time.Nanosecond,
		"us": time.Microsecond,
		"µs": time.Microsecond, // U+00B5 = micro symbol
		"μs": time.Microsecond, // U+03BC = Greek letter mu
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		`d`:  time.Hour * 24,
		`w`:  time.Hour * 24 * 7,
		`mo`: time.Hour * 24 * 30,
		`y`:  time.Hour * 24 * 365,
	}
	lengthRE   = regexp.MustCompile(`^([.0-9]+)`)
	durationRE = regexp.MustCompile(`[.0-9]+\s?((mo|(n|u|m|µ|μ)?s|m|h|d|w|y)?[a-z]*)$`)
)

func ParseDuration(s, unit string) (time.Duration, error) {
	if s == "never" {
		// never = 100 years, more or less
		return time.Hour * 24 * 365 * 100, nil
	}
	parts := lengthRE.FindStringSubmatch(s)
	if parts == nil || len(parts) != 2 || parts[1] == "" {
		return 0, fmt.Errorf("Invalid duration '%s': number not valid", s)
	}
	length, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid duration '%s': number not valid", s)
	}
	parts = durationRE.FindStringSubmatch(s)
	switch len(parts) {
	case 3:
		if parts[2] != "" {
			unit = parts[2]
			break
		}
		fallthrough
	case 2:
		if parts[1] != "" {
			unit = parts[1]
			break
		}
	}
	duration, ok := unitMap[unit]
	if !ok {
		return 0, fmt.Errorf("Invalid duration '%s': unit must be [ns] nanoseconds, [us] microseconds, [ms] milliseconds, [s]seconds, [m]inutes, [h]ours, [d]ays, [w]eeks, [mo]nths, [y]ears, or never", s)
	}
	return time.Duration(int64(length * float64(duration))), nil
}
