package provisioner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VictorLowther/jsonpatch"
	"github.com/rackn/rocket-skates/models"
)

func listThings(thing keySaver) ([]interface{}, *models.Error) {
	things := list(thing)
	res := make([]interface{}, 0, len(things))
	for _, obj := range things {
		var buf interface{}
		if err := json.Unmarshal(obj, &buf); err != nil {
			return nil, NewError(http.StatusInternalServerError,
				fmt.Sprintf("list: error unmarshalling %v: %v", string(obj), err))
		}
		res = append(res, buf)
	}
	return res, nil
}

func createThing(newThing keySaver) (interface{}, int, *models.Error) {
	finalStatus := http.StatusCreated
	oldThing := newThing.newIsh()
	if err := load(oldThing); err == nil {
		Logger.Printf("backend: Updating %v\n", oldThing.key())
		Logger.Printf("backend: Updating new %v\n", newThing.key())
		finalStatus = http.StatusOK
	} else {
		Logger.Printf("backend: Creating %v\n", newThing.key())
		oldThing = nil
	}
	if err := save(newThing, oldThing); err != nil {
		Logger.Printf("backend: Save failed: %v\n", err)
		return nil, http.StatusConflict, NewError(http.StatusConflict, err.Error())
	}
	return newThing, finalStatus, nil
}

func getThing(thing keySaver) (interface{}, *models.Error) {
	if err := load(thing); err != nil {
		return nil, NewError(http.StatusNotFound, err.Error())
	}
	return thing, nil
}

func putThing(newThing keySaver) (interface{}, *models.Error) {
	oldThing := newThing.newIsh()
	if err := load(oldThing); err != nil {
		return nil, NewError(http.StatusNotFound, err.Error())
	}
	if err := save(newThing, oldThing); err != nil {
		Logger.Printf("backend: Save failed: %v\n", err)
		return nil, NewError(http.StatusConflict, err.Error())
	}
	return newThing, nil
}

func patchThing(oldThing keySaver, patch []byte) (interface{}, *models.Error) {
	if err := load(oldThing); err != nil {
		return nil, NewError(http.StatusNotFound, err.Error())
	}
	var err error
	newThing := &Template{}
	oldThingBuf, _ := json.Marshal(oldThing)
	newThingBuf, err, loc := jsonpatch.ApplyJSON(oldThingBuf, patch)
	if err != nil {
		return nil, NewError(http.StatusConflict, fmt.Sprintf("Failed to apply patch at %d: %v\n", loc, err))
	}
	if err := json.Unmarshal(newThingBuf, &newThing); err != nil {
		return nil, NewError(http.StatusExpectationFailed, err.Error())
	}
	if err := save(newThing, oldThing); err != nil {
		return nil, NewError(http.StatusConflict, err.Error())
	}

	return newThing, nil
}

func deleteThing(thing keySaver) *models.Error {
	if err := load(thing); err != nil {
		return NewError(http.StatusNotFound, err.Error())
	}
	if err := remove(thing); err != nil {
		return NewError(http.StatusConflict, fmt.Sprintf("Failed to delete %s: %v", thing.key(), err))
	}
	return nil
}
