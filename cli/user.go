package cli

import (
	"fmt"
	"strconv"

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
		Use:   "token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]]",
		Short: "Get a login token for this user with optional parameters",
		Long:  "Creates a time-bound token for the specified user.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v needs at least 1 arg", c.UseLine())
			}

			scope := ""
			action := ""
			specific := ""
			ttl := ""
			index := 2
			for index < len(args) {
				if args[index-1] == "scope" {
					scope = args[index]
					index += 2
					continue
				}
				if args[index-1] == "action" {
					action = args[index]
					index += 2
					continue
				}
				if args[index-1] == "specific" {
					specific = args[index]
					index += 2
					continue
				}
				if args[index-1] == "ttl" {
					ttl = args[index]
					index += 2
					continue
				}
				return fmt.Errorf("%v does not support %s", c.UseLine(), args[index-1])
			}
			if index-1 != len(args) {
				return fmt.Errorf("%v needs at least 1 and pairs arg", c.UseLine())
			}

			dumpUsage = false

			p := users.NewGetUserTokenParams().WithName(args[0])
			if scope != "" {
				p = p.WithScope(&scope)
			}
			if action != "" {
				p = p.WithAction(&action)
			}
			if specific != "" {
				p = p.WithSpecific(&specific)
			}
			if ttl != "" {
				ttl64, e := strconv.ParseInt(ttl, 10, 64)
				if e != nil {
					return fmt.Errorf("ttl should be a number: %v", e)
				}
				p = p.WithTTL(&ttl64)
			}

			if d, e := session.Users.GetUserToken(p, basicAuth); e != nil {
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
