package api

import (
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestInterfaces(t *testing.T) {
	tests := []*crudTest{
		{
			name:      "list interfaces",
			expectRes: nil,
			expectErr: nil,
			op: func() (interface{}, error) {
				_, err := session.ListModel("interfaces", nil)
				return nil, err
			},
		},
		{
			name:      "show missing interfaces",
			expectRes: nil,
			expectErr: &models.Error{
				Model:    "interfaces",
				Key:      "missing",
				Type:     "GET",
				Messages: []string{"No interface"},
				Code:     404,
			},
			op: func() (interface{}, error) {
				return session.GetModel("interfaces", "missing")
			},
		},
	}
	for _, test := range tests {
		test.run(t)
	}

}
