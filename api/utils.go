package api

import (
	"encoding/json"

	"github.com/VictorLowther/jsonpatch2"
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
	return yaml.Unmarshal(buf, ref)
}

// GenPatch generates a JSON patch that will transform source into target.
// The generated patch will have all the applicable test clauses.

func GenPatch(source, target interface{}) (jsonpatch2.Patch, error) {
	srcBuf, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	tgtBuf, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	return jsonpatch2.Generate(srcBuf, tgtBuf, true)
}
