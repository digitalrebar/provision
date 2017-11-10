package api

import (
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestLeaseCrud(t *testing.T) {
	rt(t, "Get empty lease list", []models.Model{}, nil, func() (interface{}, error) {
		return session.ListModel("leases")
	}, nil)
	rt(t, "Get a malformed lease address", nil,
		&models.Error{
			Model:    "leases",
			Key:      "foo",
			Type:     "GET",
			Messages: []string{"address not valid"},
			Code:     400,
		},
		func() (interface{}, error) {
			return session.GetModel("leases", "foo")
		}, nil)
	rt(t, "Get a nonexistent lease", nil,
		&models.Error{
			Model:    "leases",
			Key:      "00000000",
			Type:     "GET",
			Messages: []string{"Not Found"},
			Code:     404,
		},
		func() (interface{}, error) {
			return session.GetModel("leases", "0.0.0.0")
		}, nil)
}
