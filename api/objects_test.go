package api

import (
	"testing"
)

func TestObject(t *testing.T) {
	test := &crudTest{
		name: "get objects",
		expectRes: []string{
			"async_action_crons",
			"async_action_templates",
			"async_actions",
			"bootenvs",
			"catalog_items",
			"contexts",
			"endpoints",
			"jobs",
			"kk",
			"leases",
			"machines",
			"params",
			"plugins",
			"pools",
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
			"version_sets",
			"workflows",
		},
		expectErr: nil,
		op: func() (interface{}, error) {
			return session.Objects()
		},
	}
	test.run(t)

}
