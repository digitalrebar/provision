package api

import (
	"testing"
)

func TestObject(t *testing.T) {
	test := &crudTest{
		name: "get objects",
		expectRes: []string{
			"bootenvs",
			"jobs",
			"leases",
			"machines",
			"params",
			"plugins",
			"preferences",
			"profiles",
			"reservations",
			"roles",
			"stages",
			"subnets",
			"tasks",
			"templates",
			"tenants",
			"users",
			"workflows",
		},
		expectErr: nil,
		op: func() (interface{}, error) {
			return session.Objects()
		},
	}
	test.run(t)

}
