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

func CastMachine(m1 *models.MachineInput) *Machine {
	return &Machine{models.MachineOutput{*m1, make([]string, 0, 0)}}
}

func PopMachine(param string) *Machine {
	if _, err := uuid.FromString(param); err == nil {
		return &Machine{models.MachineOutput{models.MachineInput{UUID: strfmt.UUID(param)},
			make([]string, 0, 0)}}
	} else {
		return &Machine{models.MachineOutput{models.MachineInput{Name: strfmt.Hostname(param)},
			make([]string, 0, 0)}}
	}
}

func MachineList(params machines.ListMachinesParams, p *models.Principal) middleware.Responder {
	allthem, err := listThings(&Machine{})
	if err != nil {
		return machines.NewListMachinesInternalServerError().WithPayload(err)
	}
	data := make([]*models.MachineOutput, 0, 0)
	for _, j := range allthem {
		original, ok := j.(models.MachineOutput)
		if ok {
			data = append(data, &original)
		}
	}
	return machines.NewListMachinesOK().WithPayload(data)
}

func MachinePost(params machines.PostMachineParams, p *models.Principal) middleware.Responder {
	item, code, err := createThing(CastMachine(params.Body))
	if err != nil {
		return machines.NewPostMachineConflict().WithPayload(err)
	}
	original, ok := item.(models.MachineOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "failed to marshall machine")
		return machines.NewPostMachineInternalServerError().WithPayload(e)
	}
	if code == http.StatusOK {
		return machines.NewPostMachineOK().WithPayload(&original)
	}
	return machines.NewPostMachineCreated().WithPayload(&original)
}

func MachineGet(params machines.GetMachineParams, p *models.Principal) middleware.Responder {
	item, err := getThing(PopMachine(params.UUID))
	if err != nil {
		return machines.NewGetMachineNotFound().WithPayload(err)
	}
	original, ok := item.(models.MachineOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "failed to marshall machine")
		return machines.NewGetMachineInternalServerError().WithPayload(e)
	}
	return machines.NewGetMachineOK().WithPayload(&original)
}

func MachinePut(params machines.PutMachineParams, p *models.Principal) middleware.Responder {
	item, err := putThing(CastMachine(params.Body))
	if err != nil {
		if err.Code == http.StatusNotFound {
			return machines.NewPutMachineNotFound().WithPayload(err)
		}
		return machines.NewPutMachineConflict().WithPayload(err)
	}
	original, ok := item.(models.MachineOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "failed to marshall machine")
		return machines.NewPutMachineInternalServerError().WithPayload(e)
	}
	return machines.NewPutMachineOK().WithPayload(&original)
}

func MachinePatch(params machines.PatchMachineParams, p *models.Principal) middleware.Responder {
	newThing := PopMachine(params.UUID)
	patch, _ := json.Marshal(params.Body)
	item, err := patchThing(newThing, patch)
	if err != nil {
		if err.Code == http.StatusNotFound {
			return machines.NewPatchMachineNotFound().WithPayload(err)
		}
		if err.Code == http.StatusConflict {
			return machines.NewPatchMachineConflict().WithPayload(err)
		}
		return machines.NewPatchMachineExpectationFailed().WithPayload(err)
	}
	original, ok := item.(models.MachineOutput)
	if !ok {
		e := NewError(http.StatusInternalServerError, "failed to marshall machine")
		return machines.NewPatchMachineInternalServerError().WithPayload(e)
	}
	return machines.NewPatchMachineOK().WithPayload(&original)
}

func MachineDelete(params machines.DeleteMachineParams, p *models.Principal) middleware.Responder {
	err := deleteThing(PopMachine(params.UUID))
	if err != nil {
		if err.Code == http.StatusNotFound {
			return machines.NewDeleteMachineNotFound().WithPayload(err)
		}
		return machines.NewDeleteMachineConflict().WithPayload(err)
	}
	return machines.NewDeleteMachineNoContent()
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
	res := &Machine{models.MachineOutput{models.MachineInput{Name: n.Name, UUID: strfmt.UUID(n.MachineOutput.MachineInput.UUID)},
		make([]string, 0, 0)}}
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
		oldBootEnv := NewBootenv(old.BootEnv)
		if err := load(oldBootEnv); err != nil {
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
	bootEnv := NewBootenv(n.BootEnv)
	if err := load(bootEnv); err != nil {
		return err
	}
	if err := bootEnv.RenderTemplates(n); err != nil {
		return err
	}
	return nil
}

func (n *Machine) onDelete() error {
	bootEnv := NewBootenv(n.BootEnv)
	if err := load(bootEnv); err != nil {
		return err
	}
	bootEnv.DeleteRenderedTemplates(n)
	return nil
}

func (b *Machine) List() ([]*Machine, error) {
	things := list(b)
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
