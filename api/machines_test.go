package api

import (
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestMachineCrud(t *testing.T) {
	rt(t, "Create initial helpers", nil, nil,
		func() (interface{}, error) {
			for _, n := range []string{"stage-prof", "jill", "jean"} {
				p := &models.Profile{Name: n}
				if err := session.CreateModel(p); err != nil {
					return nil, err
				}
			}
			for _, t := range []string{"jamie", "justine"} {
				t := &models.Task{Name: t}
				if err := session.CreateModel(t); err != nil {
					return nil, err
				}
			}
			st1 := &models.Stage{Name: "stage1", BootEnv: "local", Tasks: []string{"jamie", "justine"}}
			st2 := &models.Stage{
				Name:    "stage2",
				BootEnv: "local",
				Templates: []models.TemplateInfo{
					{
						Contents: `{{.Param "sp-param"}}`,
						Name:     "test",
						Path:     `{{.Machine.Path}}/file`,
					},
				}}
			b := &models.BootEnv{Name: "foo"}
			if err := session.CreateModel(st1); err != nil {
				return nil, err
			}
			if err := session.CreateModel(st2); err != nil {
				return nil, err
			}
			if err := session.CreateModel(b); err != nil {
				return nil, err
			}
			return nil, nil
		}, nil)
	rt(t, "Get nonexistent machine", nil,
		&models.Error{
			Model:    "machines",
			Key:      "foo",
			Type:     "GET",
			Messages: []string{"Not Found"},
			Code:     404,
		},
		func() (interface{}, error) {
			return session.GetModel("machines", "foo")
		}, nil)
	rt(t, "Get machine list (no machines)", []models.Model{}, nil, func() (interface{}, error) {
		return session.ListModel("machines")
	}, nil)
	baseM := mustDecode(&models.Machine{},
		`{"Uuid":"24679e38-53a2-4a82-99dd-5280139de00c"}`).(*models.Machine)
	rt(t, "Create machine (no name)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Missing Name"},
			Code:     422,
		},
		func() (interface{}, error) {
			return baseM, session.CreateModel(baseM)
		}, nil)
	rt(t, "Create machine (bad bootenv and stage)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Stage baz does not exist", "Bootenv bar does not exist"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = "foo"
			r.BootEnv = "bar"
			r.Stage = "baz"
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (invalid \\name)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Name must not contain a '/' or '\\'"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo\`
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (invalid /name)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Name must not contain a '/' or '\\'"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo/`
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (missing profile)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Profile foo (at 0) does not exist"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo`
			r.Profiles = []string{"foo"}
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (duplicate profile)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Duplicate profile jill: at 0 and 1"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo`
			r.Profiles = []string{"jill", "jill"}
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (bootenv not available)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"BootEnv foo is not available"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo`
			r.BootEnv = "foo"
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (bootenv only for unknown machines)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"BootEnv ignore does not allow Machine assignments, it has the OnlyUnknown flag."},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo`
			r.BootEnv = "ignore"
			return r, session.CreateModel(r)
		}, nil)
	rt(t, "Create machine (wants nonexistent task)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "24679e38-53a2-4a82-99dd-5280139de00c",
			Type:     "ValidationError",
			Messages: []string{"Task foo (at 0) does not exist"},
			Code:     422,
		},
		func() (interface{}, error) {
			r := models.Clone(baseM).(*models.Machine)
			r.Name = `foo`
			r.Tasks = []string{"foo"}
			return r, session.CreateModel(r)
		}, nil)
	baseM.Name = "foo"
	rt(t, "Create machine", baseM, nil,
		func() (interface{}, error) {
			return baseM, session.CreateModel(baseM)
		}, nil)
}
