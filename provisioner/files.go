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

func GetFile(params files.GetFileParams) middleware.Responder {
	fileName := path.Join(ProvOpts.FileRoot, `files`, path.Clean(params.Path))
	f, err := os.Open(fileName)
	// GREG: errors and data are different types need to figure that out.
	if err != nil {
		r := &models.Result{Code: int64(404), Messages: []string{err.Error()}}
		return files.NewGetFileNotFound().WithPayload(r)
	}
	return files.NewGetFileOK().WithPayload(f)
}

func ListFiles(params files.ListFilesParams) middleware.Responder {
	pathPart := "/" // params.Path
	res := []string{}
	ents, err := ioutil.ReadDir(path.Join(ProvOpts.FileRoot, "files", path.Clean(pathPart)))
	if err != nil {
		r := &models.Result{Code: int64(http.StatusNotFound),
			Messages: []string{fmt.Sprintf("list: error listing files: %v", err)}}
		return files.NewListFilesNotFound().WithPayload(r)
	}
	for _, ent := range ents {
		if ent.Mode().IsRegular() {
			res = append(res, ent.Name())
		} else if ent.Mode().IsDir() {
			res = append(res, ent.Name()+"/")
		}
	}

	if err != nil {
		r := &models.Result{Code: int64(http.StatusNotFound),
			Messages: []string{fmt.Sprintf("list: error listing files: %v", err)}}
		return files.NewListFilesNotFound().WithPayload(r)
	}
	return files.NewListFilesOK().WithPayload(res)
}

func UploadFile(params files.PostFileParams) middleware.Responder {
	name := params.Path
	body := params.Body
	amount := params.HTTPRequest.ContentLength

	fileTmpName := path.Join(ProvOpts.FileRoot, `files`, fmt.Sprintf(`.%s.part`, path.Clean(name)))
	fileName := path.Join(ProvOpts.FileRoot, `files`, path.Clean(name))
	if strings.HasSuffix(fileName, "/") {
		r := &models.Result{Code: int64(http.StatusForbidden),
			Messages: []string{fmt.Sprintf("upload: Cannot upload a directory")}}
		return files.NewPostFileBadRequest().WithPayload(r)
	}
	if err := os.MkdirAll(path.Dir(fileName), 0755); err != nil {
		r := &models.Result{Code: int64(http.StatusConflict),
			Messages: []string{fmt.Sprintf("upload: unable to create directory %s", path.Clean(path.Dir(name)))}}
		return files.NewPostFileConflict().WithPayload(r)
	}
	if _, err := os.Open(fileTmpName); err == nil {
		r := &models.Result{Code: int64(http.StatusConflict),
			Messages: []string{fmt.Sprintf("upload: file %s already uploading", name)}}
		return files.NewPostFileConflict().WithPayload(r)
	}
	tgt, err := os.Create(fileTmpName)
	if err != nil {
		r := &models.Result{Code: int64(http.StatusConflict),
			Messages: []string{fmt.Sprintf("upload: Unable to upload %s: %v", name, err)}}
		return files.NewPostFileConflict().WithPayload(r)
	}

	copied, err := io.Copy(tgt, body)
	if err != nil {
		os.Remove(fileTmpName)
		r := &models.Result{Code: int64(http.StatusInsufficientStorage),
			Messages: []string{fmt.Sprintf("upload: Failed to upload %s: %v", name, err)}}
		return files.NewPostFileInsufficientStorage().WithPayload(r)

	}
	if amount != 0 && copied != amount {
		os.Remove(fileTmpName)
		r := &models.Result{Code: int64(http.StatusForbidden),
			Messages: []string{fmt.Sprintf("upload: Failed to upload entire file %s: %d bytes expected, %d bytes recieved", name, amount, copied)}}
		return files.NewPostFileBadRequest().WithPayload(r)
	}
	os.Remove(fileName)
	os.Rename(fileTmpName, fileName)
	return files.NewPostFileCreated().WithPayload(files.PostFileCreatedBody{Name: &name, Size: &copied})
}

func DeleteFile(params files.DeleteFileParams) middleware.Responder {
	fileName := path.Join(ProvOpts.FileRoot, `files`, path.Clean(params.Path))
	if fileName == path.Join(ProvOpts.FileRoot, `files`) {
		r := &models.Result{Code: int64(http.StatusForbidden),
			Messages: []string{"delete: Not allowed to remove files dir"}}
		return files.NewDeleteFileForbidden().WithPayload(r)
	}
	if err := os.Remove(fileName); err != nil {
		r := &models.Result{Code: int64(http.StatusNotFound),
			Messages: []string{fmt.Sprintf("delete: unable to delete %s: %v", params.Path, err)}}
		return files.NewDeleteFileNotFound().WithPayload(r)
	}
	return files.NewDeleteFileNoContent()
}
