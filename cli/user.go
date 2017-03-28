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
	d, e := session.Users.ListUsers(users.NewListUsersParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Get(id string) (interface{}, error) {
	d, e := session.Users.GetUser(users.NewGetUserParams().WithName(id))
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
	d, e := session.Users.CreateUser(users.NewCreateUserParams().WithBody(user))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Put(id string, obj interface{}) (interface{}, error) {
	user, ok := obj.(*models.User)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to user put")
	}
	d, e := session.Users.PutUser(users.NewPutUserParams().WithName(id).WithBody(user))
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
	d, e := session.Users.PatchUser(users.NewPatchUserParams().WithName(id).WithBody(data))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be UserOps) Delete(id string) (interface{}, error) {
	d, e := session.Users.DeleteUser(users.NewDeleteUserParams().WithName(id))
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
	res.AddCommand(commands...)
	return res
}
