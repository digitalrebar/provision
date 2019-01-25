package cli

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
)

var encodeJsonPtr = strings.NewReplacer("~", "~0", "/", "~1")

// String translates a pointerSegment into a regular string, encoding it as we go.
func makeJsonPtr(s string) string {
	return encodeJsonPtr.Replace(string(s))
}

func bufOrFileDecode(ref string, data interface{}) (err error) {
	if buf, terr := bufOrStdin(ref); terr != nil {
		err = fmt.Errorf("Unable to process reference object: %v\n", terr)
		return
	} else {
		err = api.DecodeYaml(buf, &data)
		if err != nil {
			err = fmt.Errorf("Unable to unmarshal reference object: %v\n", err)
			return
		}
	}
	return
}

func urlOrFileAsReadCloser(src string) (io.ReadCloser, error) {
	if s, err := os.Lstat(src); err == nil && s.Mode().IsRegular() {
		fi, err := os.Open(src)
		if err != nil {
			return nil, fmt.Errorf("Error opening %s: %v", src, err)
		}
		return fi, nil
	}
	if u, err := url.Parse(src); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		if res, err := client.Get(src); err != nil {
			return nil, err
		} else {
			return res.Body, nil
		}
	} else if err == nil && u.Scheme == "file" {
		return nil, fmt.Errorf("file:// scheme not supported")
	}
	return nil, fmt.Errorf("Must specify a file or url")
}

func bufOrFile(src string) ([]byte, error) {
	if s, err := os.Lstat(src); err == nil && s.Mode().IsRegular() {
		return ioutil.ReadFile(src)
	}
	if u, err := url.Parse(src); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		if res, err := client.Get(src); err != nil {
			return nil, err
		} else {
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			return []byte(body), err
		}
	} else if err == nil && u.Scheme == "file" {
		return nil, fmt.Errorf("file:// scheme not supported")
	}
	return []byte(src), nil
}

func bufOrStdin(src string) ([]byte, error) {
	if src == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return bufOrFile(src)
}

func into(src string, res interface{}) error {
	buf, err := bufOrStdin(src)
	if err != nil {
		return fmt.Errorf("Error reading from stdin: %v", err)
	}
	return api.DecodeYaml(buf, &res)
}

func safeMergeJSON(src interface{}, toMerge []byte) ([]byte, error) {
	sv := reflect.ValueOf(src)
	for sv.Kind() == reflect.Ptr || sv.Kind() == reflect.Interface {
		sv = sv.Elem()
	}
	toMergeObj := make(map[string]interface{})
	if err := json.Unmarshal(toMerge, &toMergeObj); err != nil {
		return nil, err
	}
	targetObj := map[string]interface{}{}
	if err := utils.Remarshal(src, &targetObj); err != nil {
		return nil, err
	}
	outObj, ok := utils.Merge(targetObj, toMergeObj).(map[string]interface{})
	if !ok {
		return nil, errors.New("Cannot happen in safeMergeJSON")
	}
	if sv.Kind() == reflect.Struct {
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
	// For Raw!!
	return json.Marshal(outObj)
}

func mergeInto(src models.Model, changes []byte) (models.Model, error) {
	buf, err := safeMergeJSON(src, changes)
	if err != nil {
		return nil, err
	}
	dest := models.Clone(src)
	return dest, json.Unmarshal(buf, &dest)
}

func mergeFromArgs(src models.Model, changes string) (models.Model, error) {
	// We have to load this and then convert to json to merge safely.
	data := map[string]interface{}{}
	if err := bufOrFileDecode(changes, &data); err != nil {
		return nil, err
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return mergeInto(src, buf)
}

func d(msg string, args ...interface{}) {
	if debug {
		log.Printf(msg, args...)
	}
}

func prettyPrintBuf(o interface{}) (buf []byte, err error) {
	var v interface{}
	if err := utils.Remarshal(o, &v); err != nil {
		return nil, err
	}
	return api.Pretty(format, v)
}

func prettyPrint(o interface{}) (err error) {
	if noPretty {
		fmt.Printf("%v", o)
		return nil
	}
	if buf, err := prettyPrintBuf(o); err != nil {
		return err
	} else {
		fmt.Println(string(buf))
	}
	return nil
}
