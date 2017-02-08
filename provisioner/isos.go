package provisioner

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/galthaus/swagger-test/models"
	"github.com/galthaus/swagger-test/restapi/operations/isos"
)

func ListIsos(params isos.ListIsosParams) middleware.Responder {
	res := []string{}
	ents, err := ioutil.ReadDir(path.Join(ProvOpts.FileRoot, "isos"))
	if err != nil {
		r := &models.Result{Code: int64(http.StatusNotFound),
			Messages: []string{fmt.Sprintf("list: error listing isos: %v", err)}}
		return isos.NewListIsosNotFound().WithPayload(r)
	}
	for _, ent := range ents {
		if !ent.Mode().IsRegular() {
			continue
		}
		res = append(res, ent.Name())
	}
	return isos.NewListIsosOK().WithPayload(res)
}

func GetIso(params isos.GetIsoParams) middleware.Responder {
	fileName := path.Join(ProvOpts.FileRoot, `isos`, path.Base(params.Name))
	f, err := os.Open(fileName)
	// GREG: errors and data are different types need to figure that out.
	if err != nil {
		r := &models.Result{Code: int64(404), Messages: []string{err.Error()}}
		return isos.NewGetIsoNotFound().WithPayload(r)
	}
	return isos.NewGetIsoOK().WithPayload(f)
}

func reloadBootenvsForIso(name string) {
	env := &BootEnv{}
	newEnv := &BootEnv{}
	for _, blob := range backend.list(env) {
		if err := json.Unmarshal(blob, env); err != nil {
			continue
		}
		if env.Available || env.OS.IsoFile != name {
			continue
		}
		json.Unmarshal(blob, newEnv)
		newEnv.Available = true
		backend.save(newEnv, env)
	}
}

func UploadIso(params isos.PostIsoParams) middleware.Responder {
	name := params.Name
	body := params.Body
	amount := params.HTTPRequest.ContentLength

	isoTmpName := path.Join(ProvOpts.FileRoot, `isos`, fmt.Sprintf(`.%s.part`, path.Base(name)))
	isoName := path.Join(ProvOpts.FileRoot, `isos`, path.Base(name))
	if _, err := os.Open(isoTmpName); err == nil {
		r := &models.Result{Code: int64(http.StatusConflict),
			Messages: []string{fmt.Sprintf("upload: iso %s already uploading", name)}}
		return isos.NewPostIsoConflict().WithPayload(r)
	}
	tgt, err := os.Create(isoTmpName)
	if err != nil {
		r := &models.Result{Code: int64(http.StatusConflict),
			Messages: []string{fmt.Sprintf("upload: Unable to upload %s: %v", name, err)}}
		return isos.NewPostIsoConflict().WithPayload(r)
	}

	copied, err := io.Copy(tgt, body)
	if err != nil {
		os.Remove(isoTmpName)
		r := &models.Result{Code: int64(http.StatusInsufficientStorage),
			Messages: []string{fmt.Sprintf("upload: Failed to upload %s: %v", name, err)}}
		return isos.NewPostIsoInsufficientStorage().WithPayload(r)
	}
	if amount != 0 && copied != amount {
		os.Remove(isoTmpName)
		r := &models.Result{Code: int64(http.StatusBadRequest),
			Messages: []string{fmt.Sprintf("upload: Failed to upload entire file %s: %d bytes expected, %d bytes recieved", name, amount, copied)}}
		return isos.NewPostIsoBadRequest().WithPayload(r)
	}
	os.Remove(isoName)
	os.Rename(isoTmpName, isoName)
	go reloadBootenvsForIso(name)
	return isos.NewPostIsoCreated().WithPayload(isos.PostIsoCreatedBody{Name: &name, Size: &copied})
}

func DeleteIso(params isos.DeleteIsoParams) middleware.Responder {
	isoName := path.Join(ProvOpts.FileRoot, `isos`, path.Base(params.Name))
	if err := os.Remove(isoName); err != nil {
		r := &models.Result{Code: int64(http.StatusNotFound),
			Messages: []string{fmt.Sprintf("delete: unable to delete %s: %v", params.Name, err)}}
		return isos.NewDeleteIsoNotFound().WithPayload(r)
	}
	return isos.NewDeleteIsoNoContent()
}
