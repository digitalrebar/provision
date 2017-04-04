package cli

import (
	"fmt"

	"github.com/rackn/rocket-skates/client/users"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

type UserOps struct{}

func (be UserOps) GetType() interface{} {
	return &models.User{}
}

func (be UserOps) GetId(obj interface{}) (string, error) {
	user, ok := obj.(*models.User)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to user create")
	}
	return *user.Name, nil
}

func (be UserOps) List() (interface{}, error) {
	d, e := session.Users.ListUsers(users.NewListUsersParams(), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Get(id string) (interface{}, error) {
	d, e := session.Users.GetUser(users.NewGetUserParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Create(obj interface{}) (interface{}, error) {
	user, ok := obj.(*models.User)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to user create")
	}
	d, e := session.Users.CreateUser(users.NewCreateUserParams().WithBody(user), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to user patch")
	}
	d, e := session.Users.PatchUser(users.NewPatchUserParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Delete(id string) (interface{}, error) {
	d, e := session.Users.DeleteUser(users.NewDeleteUserParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addUserCommands()
	App.AddCommand(tree)
}

func addUserCommands() (res *cobra.Command) {
	singularName := "user"
	name := "users"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := commonOps(singularName, name, &UserOps{})

	tokenCmd := &cobra.Command{
		Use:   "token [id]",
		Short: "Get a login token for this user",
		Long:  "Creates a time-bound token for the specified user.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v needs 1 arg", c.UseLine())
			}
			dumpUsage = false
			if d, e := session.Users.GetUserToken(users.NewGetUserTokenParams().WithName(args[0]), basicAuth); e != nil {
				return generateError(e, "Error: getToken: %v", e)
			} else {
				return prettyPrint(d.Payload)
			}
		},
	}
	commands = append(commands, tokenCmd)

	res.AddCommand(commands...)
	return res
}
