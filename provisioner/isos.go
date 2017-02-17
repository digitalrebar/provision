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

	"github.com/rackn/rocket-skates/models"
	"github.com/rackn/rocket-skates/restapi/operations/isos"
)

func ListIsos(params isos.ListIsosParams, p *models.Principal) middleware.Responder {
	res := []string{}
	ents, err := ioutil.ReadDir(path.Join(ProvOpts.FileRoot, "isos"))
	if err != nil {
		e := NewError(http.StatusNotFound,
			fmt.Sprintf("list: error listing isos: %v", err))
		return isos.NewListIsosNotFound().WithPayload(e)
	}
	for _, ent := range ents {
		if !ent.Mode().IsRegular() {
			continue
		}
		res = append(res, ent.Name())
	}
	return isos.NewListIsosOK().WithPayload(res)
}

func reloadBootenvsForIso(name string) {
	env := &BootEnv{}
	newEnv := &BootEnv{}
	for _, blob := range list(env) {
		if err := json.Unmarshal(blob, env); err != nil {
			continue
		}
		if env.Available || env.OS.IsoFile != name {
			continue
		}
		json.Unmarshal(blob, newEnv)
		newEnv.Available = true
		save(newEnv, env)
	}
}

func UploadIso(params isos.PostIsoParams, p *models.Principal) middleware.Responder {
	name := params.Name
	body := params.Body
	amount := params.HTTPRequest.ContentLength

	isoTmpName := path.Join(ProvOpts.FileRoot, `isos`, fmt.Sprintf(`.%s.part`, path.Base(name)))
	isoName := path.Join(ProvOpts.FileRoot, `isos`, path.Base(name))
	if _, err := os.Open(isoTmpName); err == nil {
		e := NewError(http.StatusConflict,
			fmt.Sprintf("upload: iso %s already uploading", name))
		return isos.NewPostIsoConflict().WithPayload(e)
	}
	tgt, err := os.Create(isoTmpName)
	if err != nil {
		e := NewError(http.StatusConflict,
			fmt.Sprintf("upload: Unable to upload %s: %v", name, err))
		return isos.NewPostIsoConflict().WithPayload(e)
	}

	copied, err := io.Copy(tgt, body)
	if err != nil {
		os.Remove(isoTmpName)
		e := NewError(http.StatusInsufficientStorage,
			fmt.Sprintf("upload: Failed to upload %s: %v", name, err))
		return isos.NewPostIsoInsufficientStorage().WithPayload(e)
	}
	if amount != 0 && copied != amount {
		os.Remove(isoTmpName)
		e := NewError(http.StatusBadRequest,
			fmt.Sprintf("upload: Failed to upload entire file %s: %d bytes expected, %d bytes recieved", name, amount, copied))
		return isos.NewPostIsoBadRequest().WithPayload(e)
	}
	os.Remove(isoName)
	os.Rename(isoTmpName, isoName)
	go reloadBootenvsForIso(name)
	return isos.NewPostIsoCreated().WithPayload(isos.PostIsoCreatedBody{Name: &name, Size: &copied})
}

func GetIso(params isos.GetIsoParams, p *models.Principal) middleware.Responder {
	fileName := path.Join(ProvOpts.FileRoot, `isos`, path.Base(params.Name))
	f, err := os.Open(fileName)
	// GREG: errors and data are different types need to figure that out.
	if err != nil {
		e := NewError(http.StatusNotFound, err.Error())
		return isos.NewGetIsoNotFound().WithPayload(e)
	}
	return isos.NewGetIsoOK().WithPayload(f)
}

func DeleteIso(params isos.DeleteIsoParams, p *models.Principal) middleware.Responder {
	isoName := path.Join(ProvOpts.FileRoot, `isos`, path.Base(params.Name))
	if err := os.Remove(isoName); err != nil {
		e := NewError(http.StatusNotFound,
			fmt.Sprintf("delete: unable to delete %s: %v", params.Name, err))
		return isos.NewDeleteIsoNotFound().WithPayload(e)
	}
	return isos.NewDeleteIsoNoContent()
}
