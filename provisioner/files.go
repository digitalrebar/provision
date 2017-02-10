package provisioner

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/rackn/rocket-skates/models"
	"github.com/rackn/rocket-skates/restapi/operations/files"
)

func ListFiles(params files.ListFilesParams, p *models.Principal) middleware.Responder {
	pathPart := "/" // params.Path
	res := []string{}
	ents, err := ioutil.ReadDir(path.Join(ProvOpts.FileRoot, "files", path.Clean(pathPart)))
	if err != nil {
		e := NewError(http.StatusNotFound,
			fmt.Sprintf("list: error listing files: %v", err))
		return files.NewListFilesNotFound().WithPayload(e)
	}
	for _, ent := range ents {
		if ent.Mode().IsRegular() {
			res = append(res, ent.Name())
		} else if ent.Mode().IsDir() {
			res = append(res, ent.Name()+"/")
		}
	}

	if err != nil {
		e := NewError(http.StatusNotFound,
			fmt.Sprintf("list: error listing files: %v", err))
		return files.NewListFilesNotFound().WithPayload(e)
	}
	return files.NewListFilesOK().WithPayload(res)
}

func UploadFile(params files.PostFileParams, p *models.Principal) middleware.Responder {
	name := params.Path
	body := params.Body
	amount := params.HTTPRequest.ContentLength

	fileTmpName := path.Join(ProvOpts.FileRoot, `files`, fmt.Sprintf(`.%s.part`, path.Clean(name)))
	fileName := path.Join(ProvOpts.FileRoot, `files`, path.Clean(name))
	if strings.HasSuffix(fileName, "/") {
		e := NewError(http.StatusBadRequest,
			fmt.Sprintf("upload: Cannot upload a directory"))
		return files.NewPostFileBadRequest().WithPayload(e)
	}
	if err := os.MkdirAll(path.Dir(fileName), 0755); err != nil {
		e := NewError(http.StatusConflict,
			fmt.Sprintf("upload: unable to create directory %s", path.Clean(path.Dir(name))))
		return files.NewPostFileConflict().WithPayload(e)
	}
	if _, err := os.Open(fileTmpName); err == nil {
		e := NewError(http.StatusConflict,
			fmt.Sprintf("upload: file %s already uploading", name))
		return files.NewPostFileConflict().WithPayload(e)
	}
	tgt, err := os.Create(fileTmpName)
	if err != nil {
		e := NewError(http.StatusConflict,
			fmt.Sprintf("upload: Unable to upload %s: %v", name, err))
		return files.NewPostFileConflict().WithPayload(e)
	}

	copied, err := io.Copy(tgt, body)
	if err != nil {
		os.Remove(fileTmpName)
		e := NewError(http.StatusInsufficientStorage,
			fmt.Sprintf("upload: Failed to upload %s: %v", name, err))
		return files.NewPostFileInsufficientStorage().WithPayload(e)

	}
	if amount != 0 && copied != amount {
		os.Remove(fileTmpName)
		e := NewError(http.StatusBadRequest,
			fmt.Sprintf("upload: Failed to upload entire file %s: %d bytes expected, %d bytes recieved", name, amount, copied))
		return files.NewPostFileBadRequest().WithPayload(e)
	}
	os.Remove(fileName)
	os.Rename(fileTmpName, fileName)
	return files.NewPostFileCreated().WithPayload(files.PostFileCreatedBody{Name: &name, Size: &copied})
}

func GetFile(params files.GetFileParams, p *models.Principal) middleware.Responder {
	fileName := path.Join(ProvOpts.FileRoot, `files`, path.Clean(params.Path))
	f, err := os.Open(fileName)
	// GREG: errors and data are different types need to figure that out.
	if err != nil {
		e := NewError(http.StatusNotFound, err.Error())
		return files.NewGetFileNotFound().WithPayload(e)
	}
	return files.NewGetFileOK().WithPayload(f)
}

func DeleteFile(params files.DeleteFileParams, p *models.Principal) middleware.Responder {
	fileName := path.Join(ProvOpts.FileRoot, `files`, path.Clean(params.Path))
	if fileName == path.Join(ProvOpts.FileRoot, `files`) {
		e := NewError(http.StatusForbidden, "delete: Not allowed to remove files dir")
		return files.NewDeleteFileForbidden().WithPayload(e)
	}
	if err := os.Remove(fileName); err != nil {
		e := NewError(http.StatusNotFound,
			fmt.Sprintf("delete: unable to delete %s: %v", params.Path, err))
		return files.NewDeleteFileNotFound().WithPayload(e)
	}
	return files.NewDeleteFileNoContent()
}
