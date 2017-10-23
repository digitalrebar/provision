package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"

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
	return yaml.Unmarshal(buf, ref)
}

func makeParams(params ...string) map[string]string {
	if len(params)%2 != 0 {
		log.Panicf("makeParams requires an even number of arguments")
	}
	res := map[string]string{}
	for i := 1; i <= len(params); i += 2 {
		res[params[i-1]] = params[i]
	}
	return res
}

// Unmarshal is a helper for decoding the body of a response from the server.
// It should be called in one of two ways:
//
// The first is when you expect the response body to contain a blob
// of data that needs to be streamed somewhere.  In that case, ref
// should be an io.Writer, and the Content-Type header will be ignored.
//
// The second is when you expect the response body to contain a
// serialized object to be unmarshalled.  In that case, the response's
// Content-Type will be used as a hint to decide how to unmarshall the
// recieved data into ref.
//
// In either case, if there are any errors in the unmarshalling
// process or the response StatusCode indicates non-success, an error
// will be returned and you should not expect ref to contain vaild
// data.
func Unmarshal(resp *http.Response, ref interface{}) error {
	if resp != nil {
		defer resp.Body.Close()
	}
	if wr, ok := ref.(io.Writer); ok && resp.StatusCode < 300 {
		_, err := io.Copy(wr, resp.Body)
		return err
	}
	var dec Decoder
	ct := resp.Header.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)
	switch mt {
	case "application/json":
		dec = json.NewDecoder(resp.Body)
	default:
		buf := &bytes.Buffer{}
		io.Copy(buf, resp.Body)
		log.Printf("Got %v: %v", ct, string(buf.Bytes()))
		log.Printf("%v", resp.Request.URL)
		return fmt.Errorf("Cannot handle content-type %s", ct)
	}
	if dec == nil {
		return fmt.Errorf("No decoder for content-type %s", ct)
	}
	if resp.StatusCode >= 400 {
		res := &models.Error{}
		if err := dec.Decode(res); err != nil {
			return err
		}
		return res
	}
	var err error
	if ref != nil {
		err = dec.Decode(ref)
	}
	return err
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
