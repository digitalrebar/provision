package provisioner

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"

	strfmt "github.com/go-openapi/strfmt"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/rackn/rocket-skates/models"
	"github.com/rackn/rocket-skates/restapi/operations/machines"
	"github.com/satori/go.uuid"
)

type Machine struct {
	models.MachineOutput
}

func PopMachine(param string) *Machine {
	if _, err := uuid.FromString(param); err == nil {
		return &Machine{models.MachineOutput{models.MachineInput{UUID: strfmt.UUID(param)}}}
	} else {
		return &Machine{models.MachineOutput{models.MachineInput{Name: strfmt.Hostname(param)}}}
	}
}

func CastMachine(m1 *models.MachineInput) *Machine {
	return &Machine{models.MachineOutput{*m1}}
}

func MachinePatch(params machines.PatchMachineParams) middleware.Responder {
	newThing := PopMachine(params.UUID)
	patch, _ := json.Marshal(params.Body)
	item, code, err := updateThing(newThing, patch)
	if err != nil {
		r := &models.Result{Code: int64(code), Messages: []string{err.Message}}
		return machines.NewPatchMachineExpectationFailed().WithPayload(r)
	}
	original, ok := item.(models.MachineOutput)
	if !ok {
		r := &models.Result{Code: http.StatusInternalServerError,
			Messages: []string{"Failed to convert template"}}
		return machines.NewPatchMachineInternalServerError().WithPayload(r)
	}
	r := &models.Result{Code: int64(http.StatusOK), Messages: []string{}}
	m := machines.PatchMachineAcceptedBody{Result: r, Data: &original}
	return machines.NewPatchMachineAccepted().WithPayload(m)
}

func MachineDelete(params machines.DeleteMachineParams) middleware.Responder {
	code, err := deleteThing(PopMachine(params.UUID))
	if err != nil {
		r := &models.Result{Code: int64(code), Messages: []string{err.Message}}
		return machines.NewDeleteMachineConflict().WithPayload(r)
	}
	return machines.NewDeleteMachineNoContent()
}

func MachineGet(params machines.GetMachineParams) middleware.Responder {
	item, err := getThing(PopMachine(params.UUID))
	if err != nil {
		r := &models.Result{Code: http.StatusNotFound, Messages: []string{err.Message}}
		return machines.NewGetMachineNotFound().WithPayload(r)
	}
	r := &models.Result{Code: http.StatusOK, Messages: []string{}}
	original, ok := item.(models.MachineOutput)
	if !ok {
		r := &models.Result{Code: http.StatusInternalServerError, Messages: []string{err.Message}}
		return machines.NewGetMachineInternalServerError().WithPayload(r)
	}
	m := machines.GetMachineOKBody{Result: r, Data: &original}
	return machines.NewGetMachineOK().WithPayload(m)
}

func MachineList(params machines.ListMachinesParams) middleware.Responder {
	allthem, err := listThings(&Machine{})
	if err != nil {
		r := &models.Result{Code: http.StatusInternalServerError, Messages: []string{err.Message}}
		return machines.NewListMachinesInternalServerError().WithPayload(r)
	}
	r := &models.Result{Code: http.StatusOK, Messages: []string{}}
	data := make([]*models.MachineOutput, 0, 0)
	for _, j := range allthem {
		original, ok := j.(models.MachineOutput)
		if ok {
			data = append(data, &original)
		}
	}
	return machines.NewListMachinesOK().WithPayload(machines.ListMachinesOKBody{Result: r, Data: data})
}

func MachinePost(params machines.PostMachineParams) middleware.Responder {
	item, code, err := createThing(CastMachine(params.Body))
	if err != nil {
		r := &models.Result{Code: int64(code), Messages: []string{err.Message}}
		return machines.NewPostMachineConflict().WithPayload(r)
	}
	r := &models.Result{Code: http.StatusOK, Messages: []string{}}
	original, ok := item.(models.MachineOutput)
	if !ok {
		r := &models.Result{Code: http.StatusInternalServerError, Messages: []string{err.Message}}
		return machines.NewPostMachineInternalServerError().WithPayload(r)
	}
	m := machines.PostMachineCreatedBody{Result: r, Data: &original}
	return machines.NewPostMachineCreated().WithPayload(m)
}

func MachinePut(params machines.PutMachineParams) middleware.Responder {
	item, code, err := putThing(CastMachine(params.Body))
	if err != nil {
		r := &models.Result{Code: int64(code), Messages: []string{err.Message}}
		return machines.NewPutMachineNotFound().WithPayload(r)
	}
	r := &models.Result{Code: http.StatusOK, Messages: []string{}}
	original, ok := item.(models.MachineOutput)
	if !ok {
		r := &models.Result{Code: http.StatusInternalServerError,
			Messages: []string{"failed to cast template"}}
		return machines.NewPutMachineInternalServerError().WithPayload(r)
	}
	m := machines.PutMachineOKBody{Result: r, Data: &original}
	return machines.NewPutMachineOK().WithPayload(m)
}

// HexAddress returns Address in raw hexadecimal format, suitable for
// pxelinux and elilo usage.
func (n *Machine) HexAddress() string {
	addr := net.ParseIP(n.Address.String()).To4()
	hexIP := []byte(addr)
	return fmt.Sprintf("%02X%02X%02X%02X", hexIP[0], hexIP[1], hexIP[2], hexIP[3])
}

func (n *Machine) ShortName() string {
	idx := strings.Index(n.Name.String(), ".")
	if idx == -1 {
		return n.Name.String()
	}
	return n.Name.String()[:idx]
}

func (n *Machine) UUID() string {
	if n.MachineOutput.MachineInput.UUID == "" {
		return n.Name.String()
	}
	return n.MachineOutput.MachineInput.UUID.String()
}

func (n *Machine) Url() string {
	return ProvisionerURL + "/" + n.key()
}

func (n *Machine) prefix() string {
	return "machines"
}

func (n *Machine) Path() string {
	return path.Join(n.prefix(), n.UUID())
}

func (n *Machine) key() string {
	return n.Path()
}

func (n *Machine) tenantId() int64 {
	return n.TenantID
}

func (n *Machine) setTenantId(tid int64) {
	n.TenantID = tid
}

func (n *Machine) typeName() string {
	return "MACHINE"
}

func (n *Machine) newIsh() keySaver {
	res := &Machine{models.MachineOutput{models.MachineInput{Name: n.Name, UUID: n.MachineOutput.MachineInput.UUID}}}
	return keySaver(res)
}

func (n *Machine) onChange(oldThing interface{}) error {
	if old, ok := oldThing.(*Machine); ok && old != nil {
		if old.MachineOutput.MachineInput.UUID != "" {
			if old.MachineOutput.MachineInput.UUID != n.MachineOutput.MachineInput.UUID {
				return fmt.Errorf("machine: Cannot change machine UUID %s", old.UUID)
			}
		} else if old.Name != n.Name {
			return fmt.Errorf("machine: Cannot change name of machine %s", old.Name)
		}
		oldBootEnv := &BootEnv{models.BootenvInput{Name: old.BootEnv}, nil, nil}
		if err := backend.load(oldBootEnv); err != nil {
			return err
		}
		oldBootEnv.DeleteRenderedTemplates(old)
	}
	addr := net.ParseIP(n.Address.String())
	if addr != nil {
		addr = addr.To4()
	}
	if addr == nil {
		return fmt.Errorf("machine: %s  is not a valid IPv4 address", n.Address)
	}
	bootEnv := &BootEnv{models.BootenvInput{Name: n.BootEnv}, nil, nil}
	if err := backend.load(bootEnv); err != nil {
		return err
	}
	if err := bootEnv.RenderTemplates(n); err != nil {
		return err
	}
	return nil
}

func (n *Machine) onDelete() error {
	bootEnv := &BootEnv{models.BootenvInput{Name: n.BootEnv}, nil, nil}
	if err := backend.load(bootEnv); err != nil {
		return err
	}
	bootEnv.DeleteRenderedTemplates(n)
	return nil
}

func (b *Machine) List() ([]*Machine, error) {
	things := backend.list(b)
	res := make([]*Machine, len(things))
	for i, blob := range things {
		machine := &Machine{}
		if err := json.Unmarshal(blob, machine); err != nil {
			return nil, err
		}
		res[i] = machine
	}
	return res, nil
}

func (b *Machine) RebuildRebarData() error {
	return nil
}
