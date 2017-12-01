package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
)

func bufOrFile(src string) ([]byte, error) {
	if s, err := os.Lstat(src); err == nil && s.Mode().IsRegular() {
		return ioutil.ReadFile(src)
	}
	return []byte(src), nil
}

func bufOrStdin(src string) ([]byte, error) {
	if src == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return []byte(src), nil
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
	if sv.Kind() != reflect.Struct {
		log.Panicf("first arg to safeMergeJSON is not a struct! %#v", src)
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

func mergeInto(src models.Model, changes []byte) (models.Model, error) {
	buf, err := safeMergeJSON(src, changes)
	if err != nil {
		return nil, err
	}
	dest := models.Clone(src)
	return dest, json.Unmarshal(buf, &dest)
}

func mergeFromArgs(src models.Model, changes string) (models.Model, error) {
	buf, err := bufOrStdin(changes)
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
