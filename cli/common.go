package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
)

func generateError(err error, sfmt string, args ...interface{}) error {
	s := fmt.Sprintf(sfmt, args...)

	ee, ok := err.(*models.Error)
	if !ok {
		return fmt.Errorf(s+": %v", err)
	}
	return ee
}

var listLimit = -1
var listOffset = -1

func PatchWithString(key, js string, op *ops) error {
	data, clone, err := session.GetModelForPatch(op.name, key)
	if err != nil {
		return generateError(err, "Failed to fetch %v: %v", op.singleName, key)
	}
	if err := api.DecodeYaml([]byte(js), &clone); err != nil {
		return fmt.Errorf("Unable to merge objects: %v\n", err)
	}
	res, err := session.PatchTo(data, clone)
	if err != nil {
		return err
	}
	return prettyPrint(res)
}

// The input function takes the object and returns the modified object and if the object changed.
func PatchWithFunction(key string, op *ops, fn func(models.Model) (models.Model, bool)) error {
	data, clone, err := session.GetModelForPatch(op.name, key)
	if err != nil {
		return generateError(err, "Failed to fetch %v: %v", op.singleName, key)
	}
	newobj, changed := fn(clone)
	if !changed {
		return nil
	}
	res, err := session.PatchTo(data, newobj)
	if err != nil {
		return err
	}
	return prettyPrint(res)
}
