package provisioner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VictorLowther/jsonpatch"
)

func listThings(thing keySaver) ([]interface{}, *pkgError) {
	things := backend.list(thing)
	res := make([]interface{}, 0, len(things))
	for _, obj := range things {
		var buf interface{}
		if err := json.Unmarshal(obj, &buf); err != nil {
			return nil, NewError(fmt.Sprintf("list: error unmarshalling %v: %v", string(obj), err))
		}
		res = append(res, buf)
	}
	return res, nil
}

func createThing(newThing keySaver) (interface{}, int, *pkgError) {
	finalStatus := http.StatusCreated
	oldThing := newThing.newIsh()
	if err := backend.load(oldThing); err == nil {
		Logger.Printf("backend: Updating %v\n", oldThing.key())
		Logger.Printf("backend: Updating new %v\n", newThing.key())
		finalStatus = http.StatusAccepted
	} else {
		Logger.Printf("backend: Creating %v\n", newThing.key())
		oldThing = nil
	}
	if err := backend.save(newThing, oldThing); err != nil {
		Logger.Printf("backend: Save failed: %v\n", err)
		return nil, http.StatusConflict, NewError(err.Error())
	}
	return newThing, finalStatus, nil
}

func getThing(thing keySaver) (interface{}, *pkgError) {
	if err := backend.load(thing); err != nil {
		return nil, NewError(err.Error())
	}
	return thing, nil
}

func putThing(newThing keySaver) (interface{}, int, *pkgError) {
	oldThing := newThing.newIsh()
	if err := backend.load(oldThing); err != nil {
		return nil, http.StatusNotFound, NewError(err.Error())
	}
	if err := backend.save(newThing, oldThing); err != nil {
		Logger.Printf("backend: Save failed: %v\n", err)
		return nil, http.StatusConflict, NewError(err.Error())
	}
	return newThing, http.StatusOK, nil
}

func updateThing(oldThing keySaver, patch []byte) (interface{}, int, *pkgError) {
	if err := backend.load(oldThing); err != nil {
		return nil, http.StatusNotFound, NewError(err.Error())
	}
	var err error
	newThing := &Template{}
	oldThingBuf, _ := json.Marshal(oldThing)
	newThingBuf, err, loc := jsonpatch.ApplyJSON(oldThingBuf, patch)
	if err != nil {
		return nil, http.StatusConflict, NewError(fmt.Sprintf("Failed to apply patch at %d: %v\n", loc, err))
	}
	if err := json.Unmarshal(newThingBuf, &newThing); err != nil {
		return nil, http.StatusExpectationFailed, NewError(err.Error())
	}
	if err := backend.save(newThing, oldThing); err != nil {
		return nil, http.StatusConflict, NewError(err.Error())
	}

	return newThing, http.StatusAccepted, nil
}

func deleteThing(thing keySaver) (int, *pkgError) {
	if err := backend.load(thing); err != nil {
		return http.StatusConflict, NewError(err.Error())
	}
	if err := backend.remove(thing); err != nil {
		return http.StatusConflict, NewError(fmt.Sprintf("Failed to delete %s: %v", thing.key(), err))
	}
	return http.StatusOK, nil
}
