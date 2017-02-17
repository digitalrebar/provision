package provisioner

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/digitalrebar/digitalrebar/go/rebar-api/api"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/rackn/rocket-skates/models"
	"github.com/rackn/rocket-skates/restapi/operations/bootenvs"
)

// RenderData is the struct that is passed to templates as a source of
// parameters and useful methods.
type RenderData struct {
	Machine        *Machine // The Machine that the template is being rendered for.
	Env            *BootEnv // The boot environment that provided the template.
	ProvisionerURL string   // The URL to the provisioner that all files should be fetched from
	CommandURL     string   // The URL of the API endpoint that this machine should talk to for command and control
}

func (r *RenderData) ProvisionerAddress() string {
	return ProvOpts.OurAddress
}

// BootParams is a helper function that expands the BootParams
// template from the boot environment.
func (r *RenderData) BootParams() (string, error) {
	res := &bytes.Buffer{}
	if r.Env.bootParamsTmpl == nil {
		return "", nil
	}
	if err := r.Env.bootParamsTmpl.Execute(res, r); err != nil {
		return "", err
	}
	return res.String(), nil
}

func (r *RenderData) ParseUrl(segment, rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	switch segment {
	case "scheme":
		return parsedUrl.Scheme, nil
	case "host":
		return parsedUrl.Host, nil
	case "path":
		return parsedUrl.Path, nil
	}
	return "", fmt.Errorf("No idea how to get URL part %s from %s", segment, rawUrl)
}

// Param is a helper function for extracting a parameter from Machine.Params
func (r *RenderData) Param(key string) (interface{}, error) {
	res, ok := r.Machine.Params[key]
	if !ok {
		return nil, fmt.Errorf("No such machine parameter %s", key)
	}
	return res, nil
}

// TemplateInfo holds information on the templates in the boot
// environment that will be expanded into files.
type TemplateInfo struct {
	pathTmpl  *template.Template
	finalPath string
	contents  *Template
}

// BootEnv encapsulates the machine-agnostic information needed by the
// provisioner to set up a boot environment.
type BootEnv struct {
	models.BootenvOutput

	bootParamsTmpl *template.Template
	TemplateInfo   []*TemplateInfo // The templates that should be expanded into files for the bot environment.
}

func CastBootenv(m1 *models.BootenvInput) *BootEnv {
	return &BootEnv{models.BootenvOutput{*m1, make([]string, 0, 0)}, nil, nil}
}

func NewBootenv(name string) *BootEnv {
	return &BootEnv{models.BootenvOutput{models.BootenvInput{Name: name}, make([]string, 0, 0)}, nil, nil}
}

func BootenvList(params bootenvs.ListBootenvsParams, p *models.Principal) middleware.Responder {
	allthem, err := listThings(&BootEnv{})
	if err != nil {
		return bootenvs.NewListBootenvsInternalServerError().WithPayload(err)
	}
	data := make([]*models.BootenvOutput, 0, 0)
	for _, j := range allthem {
		original, ok := j.(models.BootenvOutput)
		if ok {
			data = append(data, &original)
		}
	}
	return bootenvs.NewListBootenvsOK().WithPayload(data)
}

func BootenvPost(params bootenvs.PostBootenvParams, p *models.Principal) middleware.Responder {
	item, code, err := createThing(CastBootenv(params.Body))
	if err != nil {
		return bootenvs.NewPostBootenvConflict().WithPayload(err)
	}
	original, ok := item.(models.BootenvOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Could not marshal bootenv")
		return bootenvs.NewPostBootenvInternalServerError().WithPayload(e)
	}
	if code == http.StatusOK {
		return bootenvs.NewPostBootenvOK().WithPayload(&original)
	}
	return bootenvs.NewPostBootenvCreated().WithPayload(&original)
}

func BootenvGet(params bootenvs.GetBootenvParams, p *models.Principal) middleware.Responder {
	item, err := getThing(NewBootenv(params.Name))
	if err != nil {
		return bootenvs.NewGetBootenvNotFound().WithPayload(err)
	}
	original, ok := item.(models.BootenvOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Could not marshal bootenv")
		return bootenvs.NewGetBootenvInternalServerError().WithPayload(e)
	}
	return bootenvs.NewGetBootenvOK().WithPayload(&original)
}

func BootenvPut(params bootenvs.PutBootenvParams, p *models.Principal) middleware.Responder {
	item, err := putThing(CastBootenv(params.Body))
	if err != nil {
		if err.Code == http.StatusConflict {
			return bootenvs.NewPutBootenvConflict().WithPayload(err)
		}
		return bootenvs.NewPutBootenvNotFound().WithPayload(err)
	}
	original, ok := item.(models.BootenvOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Could not marshal bootenv")
		return bootenvs.NewPutBootenvInternalServerError().WithPayload(e)
	}
	return bootenvs.NewPutBootenvOK().WithPayload(&original)
}

func BootenvPatch(params bootenvs.PatchBootenvParams, p *models.Principal) middleware.Responder {
	newThing := NewBootenv(params.Name)
	patch, _ := json.Marshal(params.Body)
	item, err := patchThing(newThing, patch)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return bootenvs.NewPatchBootenvNotFound().WithPayload(err)
		}
		if err.Code == http.StatusConflict {
			return bootenvs.NewPatchBootenvConflict().WithPayload(err)
		}
		return bootenvs.NewPatchBootenvExpectationFailed().WithPayload(err)
	}
	original, ok := item.(models.BootenvOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "Could not marshal bootenv")
		return bootenvs.NewPatchBootenvInternalServerError().WithPayload(e)
	}
	return bootenvs.NewPatchBootenvOK().WithPayload(&original)
}

func BootenvDelete(params bootenvs.DeleteBootenvParams, p *models.Principal) middleware.Responder {
	err := deleteThing(NewBootenv(params.Name))
	if err != nil {
		if err.Code == http.StatusNotFound {
			return bootenvs.NewDeleteBootenvNotFound().WithPayload(err)
		}
		return bootenvs.NewDeleteBootenvConflict().WithPayload(err)
	}
	return bootenvs.NewDeleteBootenvNoContent()
}

func (b *BootEnv) Error() string {
	return strings.Join(b.Errors, "\n")
}

func (b *BootEnv) errorOrNil() error {
	if len(b.Errors) == 0 {
		return nil
	}
	return b
}

func (b *BootEnv) Errorf(arg string, args ...interface{}) {
	b.Errors = append(b.Errors, fmt.Sprintf(arg, args...))
}

// PathFor expands the partial paths for kernels and initrds into full
// paths appropriate for specific protocols.
//
// proto can be one of 3 choices:
//    http: Will expand to the URL the file can be accessed over.
//    tftp: Will expand to the path the file can be accessed at via TFTP.
//    disk: Will expand to the path of the file inside the provisioner container.
func (b *BootEnv) PathFor(proto, f string) string {
	res := b.OS.Name
	if res != "discovery" {
		res = path.Join(res, "install")
	}
	switch proto {
	case "disk":
		return path.Join(ProvOpts.FileRoot, res, f)
	case "tftp":
		return path.Join(res, f)
	case "http":
		return ProvisionerURL + "/" + path.Join(res, f)
	default:
		Logger.Fatalf("Unknown protocol %v", proto)
	}
	return ""
}

func (b *BootEnv) parseTemplates() {
	if b.TemplateInfo == nil {
		b.TemplateInfo = make([]*TemplateInfo, 0, len(b.Templates))
		for i := 0; i < len(b.Templates); i++ {
			b.TemplateInfo[i] = &TemplateInfo{}
		}
	}
	for ii, templateParams := range b.Templates {
		templateInfo := b.TemplateInfo[ii]
		pathTmpl, err := template.New(templateParams.Name).Parse(templateParams.Path)
		if err != nil {
			b.Errorf("bootenv: Error compiling path template %s (%s): %v",
				templateParams.Name,
				templateParams.Path,
				err)
			continue
		}
		templateInfo.pathTmpl = pathTmpl.Option("missingkey=error")
		if templateInfo.contents == nil {
			tmpl := NewTemplate(templateParams.UUID)
			if err := load(tmpl); err != nil {
				b.Errorf("bootenv: Error loading template %s for %s: %v",
					templateParams.UUID,
					templateParams.Name,
					err)
				continue
			}
			if err := tmpl.Parse(); err != nil {
				b.Errorf("bootenv: Error compiling template %s: %v\n---template---\n %s",
					templateParams.Name,
					err,
					tmpl.Contents)
				continue
			}
			templateInfo.contents = tmpl
		}

	}
	if b.BootParams != "" {
		tmpl, err := template.New("machine").Parse(b.BootParams)
		if err != nil {
			b.Errorf("bootenv: Error compiling boot parameter template: %v\n----TEMPLATE---\n%s",
				err,
				b.BootParams)
		}
		b.bootParamsTmpl = tmpl.Option("missingkey=error")
	}
	return
}

// JoinInitrds joins the fully expanded initrd paths into a comma-separated string.
func (b *BootEnv) JoinInitrds(proto string) string {
	fullInitrds := make([]string, len(b.Initrds))
	for i, initrd := range b.Initrds {
		fullInitrds[i] = b.PathFor(proto, initrd)
	}
	return strings.Join(fullInitrds, " ")
}

func (b *BootEnv) prefix() string {
	return "bootenvs"
}

func (b *BootEnv) key() string {
	return path.Join(b.prefix(), b.Name)
}

func (b *BootEnv) typeName() string {
	return "BOOTENV"
}

func (b *BootEnv) newIsh() keySaver {
	res := NewBootenv(b.Name)
	return keySaver(res)
}

// RenderPaths renders the paths of the templates for this machine.
func (b *BootEnv) RenderPaths(machine *Machine) error {
	vars := &RenderData{
		Machine:        machine,
		Env:            b,
		ProvisionerURL: ProvisionerURL,
		CommandURL:     ProvOpts.CommandURL,
	}
	for ii, templateParams := range b.Templates {
		templateInfos := b.TemplateInfo[ii]
		pathBuf := &bytes.Buffer{}
		if err := templateInfos.pathTmpl.Execute(pathBuf, vars); err != nil {
			b.Errorf("template: Error rendering path %s (%s): %v",
				templateParams.Name,
				templateParams.Path,
				err)
			continue
		}
		templateInfos.finalPath = filepath.Join(ProvOpts.FileRoot, pathBuf.String())
	}
	return b.errorOrNil()
}

// RenderTemplates renders the templates in the bootenv with the data from the machine.
func (b *BootEnv) RenderTemplates(machine *Machine) error {
	vars := &RenderData{
		Machine:        machine,
		Env:            b,
		ProvisionerURL: ProvisionerURL,
		CommandURL:     ProvOpts.CommandURL,
	}
	b.parseTemplates()
	b.RenderPaths(machine)
	var missingParams []string
	for _, param := range b.RequiredParams {
		if _, ok := machine.Params[param]; !ok {
			missingParams = append(missingParams, param)
		}
	}
	if len(missingParams) > 0 {
		b.Errorf("bootenv: %s missing required machine params for %s:\n %v", b.Name, machine.Name, missingParams)
	}
	for ii, templateParams := range b.Templates {
		templateInfos := b.TemplateInfo[ii]
		tmplPath := templateInfos.finalPath
		if err := os.MkdirAll(path.Dir(tmplPath), 0755); err != nil {
			b.Errorf("template: Unable to create dir for %s: %v", tmplPath, err)
			continue
		}

		tmplDest, err := os.Create(tmplPath)
		if err != nil {
			b.Errorf("template: Unable to create file %s: %v", tmplPath, err)
			continue
		}
		defer tmplDest.Close()
		if err := templateInfos.contents.Render(tmplDest, vars); err != nil {
			os.Remove(tmplPath)
			b.Errorf("template: Error rendering template %s: %v\n---template---\n %s",
				templateParams.Name,
				err,
				templateInfos.contents.Contents)
			continue
		}
		tmplDest.Sync()
	}
	return b.errorOrNil()
}

// DeleteRenderedTemplates deletes the templates that were rendered
// for this bootenv/machine combination.
func (b *BootEnv) DeleteRenderedTemplates(machine *Machine) {
	b.parseTemplates()
	b.RenderPaths(machine)
	for ii, _ := range b.Templates {
		tmpli := b.TemplateInfo[ii]
		if tmpli.finalPath != "" {
			os.Remove(tmpli.finalPath)
		}
	}
}

func (b *BootEnv) explodeIso() error {
	// Only explode install things
	if !strings.HasSuffix(b.Name, "-install") {
		Logger.Printf("Explode ISO: Skipping %s becausing not -install\n", b.Name)
		return nil
	}
	// Only work on things that are requested.
	if b.OS.IsoFile == "" {
		Logger.Printf("Explode ISO: Skipping %s becausing no iso image specified\n", b.Name)
		return nil
	}
	// Have we already exploded this?  If file exists, then good!
	canaryPath := b.PathFor("disk", "."+b.OS.Name+".rebar_canary")
	buf, err := ioutil.ReadFile(canaryPath)
	if err == nil && len(buf) != 0 && string(bytes.TrimSpace(buf)) == b.OS.IsoSha256 {
		Logger.Printf("Explode ISO: Skipping %s becausing canary file, %s, in place and has proper SHA256\n", b.Name, canaryPath)
		return nil
	}

	isoPath := filepath.Join(ProvOpts.FileRoot, "isos", b.OS.IsoFile)
	if _, err := os.Stat(isoPath); os.IsNotExist(err) {
		Logger.Printf("Explode ISO: Skipping %s becausing iso doesn't exist: %s\n", b.Name, isoPath)
		return nil
	}

	f, err := os.Open(isoPath)
	if err != nil {
		return fmt.Errorf("Explode ISO: For %s, failed to open iso file %s: %v", b.Name, isoPath, err)
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return fmt.Errorf("Explode ISO: For %s, failed to read iso file %s: %v", b.Name, isoPath, err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	// This will wind up being saved along with the rest of the
	// hash because explodeIso is called by OnChange before the struct gets saved.
	if b.OS.IsoSha256 == "" {
		b.OS.IsoSha256 = hash
	}

	if hash != b.OS.IsoSha256 {
		return fmt.Errorf("iso: Iso checksum bad.  Re-download image: %s: actual: %v expected: %v", isoPath, hash, b.OS.IsoSha256)
	}

	// Call extract script
	// /explode_iso.sh b.OS.Name isoPath path.Dir(canaryPath)
	cmdName := "/explode_iso.sh"
	cmdArgs := []string{b.OS.Name, isoPath, path.Dir(canaryPath), b.OS.IsoSha256}
	if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return fmt.Errorf("Explode ISO: Exec command failed for %s: %s\n", b.Name, err)
	}
	return nil
}

func (b *BootEnv) onChange(oldThing interface{}) error {
	seenPxeLinux := false
	seenELilo := false
	seenIPXE := false
	b.Errors = []string{}
	for _, template := range b.Templates {
		if template.Name == "pxelinux" {
			seenPxeLinux = true
		}
		if template.Name == "elilo" {
			seenELilo = true
		}
		if template.Name == "ipxe" {
			seenIPXE = true
		}
		if template.Name == "" ||
			template.Path == "" ||
			template.UUID == "" {
			b.Errorf("bootenv: Illegal template: %+v", template)
		}
	}
	if !seenIPXE {
		if !(seenPxeLinux && seenELilo) {
			b.Errorf("bootenv: Missing elilo or pxelinux template")
		}
	}

	// Make sure the ISO is exploded
	if b.OS.IsoFile != "" {
		Logger.Printf("Exploding ISO for %s\n", b.OS.Name)
		if err := b.explodeIso(); err != nil {
			b.Errorf("bootenv: Unable to expand ISO %s: %v", b.OS.IsoFile, err)
		}
	}

	b.parseTemplates()
	if b.Kernel != "" {
		kPath := b.PathFor("disk", b.Kernel)
		kernelStat, err := os.Stat(kPath)
		if err != nil {
			b.Errorf("bootenv: %s: missing kernel %s (%s)",
				b.Name,
				b.Kernel,
				kPath)
		} else if !kernelStat.Mode().IsRegular() {
			b.Errorf("bootenv: %s: invalid kernel %s (%s)",
				b.Name,
				b.Kernel,
				kPath)
		}
	}
	if len(b.Initrds) > 0 {
		for _, initrd := range b.Initrds {
			iPath := b.PathFor("disk", initrd)
			initrdStat, err := os.Stat(iPath)
			if err != nil {
				b.Errorf("bootenv: %s: missing initrd %s (%s)",
					b.Name,
					initrd,
					iPath)
				continue
			}
			if !initrdStat.Mode().IsRegular() {
				b.Errorf("bootenv: %s: invalid initrd %s (%s)",
					b.Name,
					initrd,
					iPath)
			}
		}
	}

	if old, ok := oldThing.(*BootEnv); ok && old != nil {
		if old.Name != b.Name {
			b.Errorf("bootenv: Cannot change name of bootenv %s", old.Name)
		}
		machine := &Machine{}
		machines, err := machine.List()
		if err != nil {
			b.Errorf("bootenv: Failed to get list of current machines: %v", err)
		}

		for _, machine := range machines {
			if machine.BootEnv != old.Name {
				continue
			}
			if err := b.RenderTemplates(machine); err != nil {
				b.Errorf("bootenv: Failed to render templates for machine %s: %v", machine.Name, err)
			}
		}
	}
	b.Available = (len(b.Errors) == 0)
	return nil
}

func (b *BootEnv) onDelete() error {
	b.Errors = []string{}
	machine := &Machine{}
	machines, err := machine.List()
	if err == nil {
		for _, machine := range machines {
			if machine.BootEnv != b.Name {
				continue
			}
			b.Errorf("Bootenv %s in use by Machine %s", b.Name, machine.Name)
		}
	}
	return b.errorOrNil()
}

func (b *BootEnv) List() ([]*BootEnv, error) {
	things := list(b)
	res := make([]*BootEnv, len(things))
	for i, blob := range things {
		bootenv := &BootEnv{}
		if err := json.Unmarshal(blob, bootenv); err != nil {
			return nil, err
		}
		res[i] = bootenv
	}
	return res, nil
}

func (b *BootEnv) RebuildRebarData() error {
	// We aren't running with a rebar client endpoint - SUCCEED!
	if rebarClient == nil {
		return nil
	}

	preferredOses := map[string]int{
		"centos-7.3.1611": 0,
		"centos-7.2.1511": 1,
		"centos-7.1.1503": 2,
		"ubuntu-16.04":    3,
		"ubuntu-14.04":    4,
		"ubuntu-15.04":    5,
		"debian-8":        6,
		"centos-6.8":      7,
		"centos-6.6":      8,
		"debian-7":        9,
		"redhat-6.5":      10,
		"ubuntu-12.04":    11,
	}

	attrValOSes := make(map[string]bool)
	attrValOS := "STRING"
	attrPref := 1000

	if !b.Available {
		return b.errorOrNil()
	}

	bes, err := b.List()
	if err != nil {
		return err
	}

	for _, be := range bes {
		if !strings.HasSuffix(be.Name, "-install") {
			continue
		}
		if !be.Available {
			continue
		}
		attrValOSes[be.OS.Name] = true
		numPref, ok := preferredOses[be.OS.Name]
		if !ok {
			numPref = 999
		}
		if numPref < attrPref {
			attrValOS = be.OS.Name
			attrPref = numPref
		}
	}

	deployment := &api.Deployment{}
	if err := rebarClient.Fetch(deployment, "system"); err != nil {
		return err
	}

	role := &api.Role{}
	if err := rebarClient.Fetch(role, "provisioner-service"); err != nil {
		return err
	}

	var tgt api.Attriber
	for {
		drs := []*api.DeploymentRole{}
		matcher := make(map[string]interface{})
		matcher["role_id"] = role.ID
		matcher["deployment_id"] = deployment.ID
		dr := &api.DeploymentRole{}
		if err := rebarClient.Match(rebarClient.UrlPath(dr), matcher, &drs); err != nil {
			return err
		}
		if len(drs) != 0 {
			tgt = drs[0]
			break
		}
		log.Printf("Waiting for provisioner-service (%v) to show up in system(%v)", role.ID, deployment.ID)
		log.Printf("drs: %#v, err: %#v", drs, err)
		time.Sleep(5 * time.Second)
	}

	attrib := &api.Attrib{}
	attrib.SetId("provisioner-available-oses")
	attrib, err = rebarClient.GetAttrib(tgt, attrib, "")
	if err != nil {
		return err
	}
	attrib.Value = attrValOSes
	if err := rebarClient.SetAttrib(tgt, attrib, ""); err != nil {
		return err
	}

	attrib = &api.Attrib{}
	attrib.SetId("provisioner-default-os")
	attrib, err = rebarClient.GetAttrib(tgt, attrib, "")
	if err != nil {
		return err
	}
	attrib.Value = attrValOS
	if err := rebarClient.SetAttrib(tgt, attrib, ""); err != nil {
		return err
	}

	if err := rebarClient.Commit(tgt); err != nil {
		return err
	}

	return nil
}
