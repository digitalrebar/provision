package api

import (
	"encoding/json"
	"fmt"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/models"
	yaml "github.com/ghodss/yaml"
)

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
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
	return models.DecodeYaml(buf, ref)
}

// Pretty marshals object acciording to the the fmt, in whatever
// passed for "pretty" according to fmt.
func Pretty(f string, obj interface{}) ([]byte, error) {
	switch f {
	case "json":
		return json.MarshalIndent(obj, "", "  ")
	case "yaml":
		return yaml.Marshal(obj)
	case "go":
		buf, err := Pretty("yaml", obj)
		if err == nil {
			return []byte(fmt.Sprintf("package main\n\nvar contentYamlString = `\n%s\n`\n", string(buf))), nil
		}
		return nil, err
	default:
		return nil, fmt.Errorf("Unknown pretty format %s", f)
	}
}

// GenPatch generates a JSON patch that will transform source into target.
// The generated patch will have all the applicable test clauses.
func GenPatch(source, target interface{}, paranoid bool) (jsonpatch2.Patch, error) {
	return models.GenPatch(source, target, paranoid)
}
