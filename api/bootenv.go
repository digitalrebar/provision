package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/digitalrebar/provision/v4/models"
)

func (c *Client) InstallRawTemplateFromFile(src string) (*models.Template, error) {
	tid := path.Base(src)
	return c.InstallRawTemplateFromFileWithId(src, tid)
}

func (c *Client) InstallRawTemplateFromFileWithId(src, tid string) (*models.Template, error) {
	buf, readErr := ioutil.ReadFile(src)
	if readErr != nil {
		err := &models.Error{
			Model: "templates",
			Type:  "CLIENT_ERROR",
			Key:   tid,
		}
		err.Errorf("Unable to import template %s", tid)
		return nil, err
	}
	tmpl := &models.Template{ID: tid}
	if fillErr := c.Req().Fill(tmpl); fillErr == nil {
		tmpl.Contents = string(buf)
		return tmpl, c.PutModel(tmpl)
	}
	tmpl.Contents = string(buf)
	return tmpl, c.CreateModel(tmpl)
}

func (c *Client) UploadISOForBootEnv(env *models.BootEnv, src io.Reader, dest string) (models.BlobInfo, error) {
	return c.PostBlob(src, "isos", dest)
}

func (c *Client) InstallISOForBootenv(env *models.BootEnv, src string, downloadOK bool) error {
	isoFiles := map[string]string{}
	if env.OS.IsoFile != "" {
		isoFiles[env.OS.IsoFile] = env.OS.IsoUrl
	}
	for _, archInfo := range env.OS.SupportedArchitectures {
		if archInfo.IsoFile != "" {
			isoFiles[archInfo.IsoFile] = archInfo.IsoUrl
		}
	}
	if len(isoFiles) == 0 {
		return nil
	}
	isos, err := c.ListBlobs("isos")
	if err != nil {
		return err
	}
	for _, iso := range isos {
		if _, ok := isoFiles[iso]; ok {
			delete(isoFiles, iso)
		}
	}
	if len(isoFiles) == 0 {
		return nil
	}
	for isoFile, isoUrl := range isoFiles {
		if st, err := os.Stat(path.Join(src, isoFile)); err != nil {
			if !downloadOK {
				res := &models.Error{
					Model: "isos",
					Type:  "DOWNLOAD_NOT_ALLOWED",
					Key:   isoFile,
				}
				res.Errorf("Iso not present at server, not present locally, and automatic download forbidden")
				return res
			}
			if isoUrl == "" {
				res := &models.Error{
					Model: "isos",
					Type:  "DOWNLOAD_NOT_POSSIBLE",
					Key:   isoFile,
				}
				res.Errorf("Bootenv %s does not have a valid upstream source for the ISO it needs", env.Key())
				return res
			}
			err = func() error {
				isoTarget, err2 := os.Create(path.Join(src, isoFile))
				if err2 != nil {
					return err2
				}
				defer isoTarget.Close()
				resp, err2 := http.Get(isoUrl)
				if err2 != nil {
					os.Remove(src)
					return err2
				}
				defer resp.Body.Close()
				if resp.StatusCode >= 300 {
					os.Remove(src)
					return fmt.Errorf("Unable to start download of %s: %s", isoUrl, resp.Status)
				}
				_, err2 = io.Copy(isoTarget, resp.Body)
				if err2 != nil {
					os.Remove(src)
				}
				return err2
			}()
			if err != nil {
				res := &models.Error{
					Model: "isos",
					Type:  "DOWNLOAD_FAILED",
					Key:   isoUrl,
				}
				res.AddError(err)
				return res
			}
		} else if st.IsDir() {
			return &models.Error{Model: "isos", Type: "ISO_SRC_IS_A_DIR", Key: src}
		}
		var isoSrc *os.File
		isoSrc, err = os.Open(path.Join(src, isoFile))
		if err != nil {
			res := &models.Error{Model: "isos", Type: "UPLOAD_NOT_POSSIBLE", Key: src}
			res.AddError(err)
			return res
		}
		func() {
			defer isoSrc.Close()
			_, err = c.UploadISOForBootEnv(env, isoSrc, isoFile)
		}()
		if err != nil {
			break
		}
	}
	return err
}

func (c *Client) InstallBootEnvFromFile(src string) (*models.BootEnv, error) {
	env := &models.BootEnv{}
	err := &models.Error{
		Model: "bootenvs",
		Type:  "CLIENT_ERROR",
	}
	if st, statErr := os.Stat(src); statErr != nil {
		err.AddError(statErr)
		return nil, err
	} else if st.IsDir() {
		err.Errorf("%s is a directory.  It needs to be a file.", src)
		return nil, err
	}
	buf, readErr := ioutil.ReadFile(src)
	if readErr != nil {
		err.AddError(readErr)
		return nil, err
	}
	if decodeErr := DecodeYaml(buf, env); decodeErr != nil {
		err.AddError(decodeErr)
		return nil, err
	}
	err.Key = env.Key()
	if found, _ := c.ExistsModel("bootenvs", env.Key()); found {
		return env, c.Req().Fill(env)
	}
	srcDir := path.Dir(src)
	tmplDir := srcDir
	if path.Base(srcDir) == "bootenvs" {
		tmplDir = path.Join(path.Dir(srcDir), "templates")
	}
	// Upload all directly-referenced templates.
	for _, ti := range env.Templates {
		if ti.ID == "" {
			continue
		}
		if _, tmplErr := c.InstallRawTemplateFromFile(path.Join(tmplDir, ti.ID)); tmplErr != nil {
			err.AddError(tmplErr)
		}
	}
	files, dirErr := ioutil.ReadDir(tmplDir)
	if dirErr == nil {
		treatAllAsTemplates := path.Base(tmplDir) == "templates"
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			fname := f.Name()
			if !treatAllAsTemplates && !strings.HasSuffix(fname, ".tmpl") {
				continue
			}
			if _, tmplErr := c.InstallRawTemplateFromFile(path.Join(tmplDir, f.Name())); tmplErr != nil {
				err.AddError(tmplErr)
			}
		}
	} else if !os.IsNotExist(dirErr) {
		err.Errorf("Cannot import extra templates: %v", dirErr)
	}

	if err.ContainsError() {
		return nil, err
	}
	return env, c.CreateModel(env)
}
