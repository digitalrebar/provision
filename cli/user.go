package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerUser)
}

func registerUser(app *cobra.Command) {
	op := &ops{
		name:       "users",
		singleName: "user",
		example:    func() models.Model { return &models.User{} },
	}
	op.addCommand(&cobra.Command{
		Use:   "password [id] [password]",
		Short: "Set the password for this id",
		Long:  "Set the password for this id",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v needs 2 args", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			pwd := &models.UserPassword{Password: args[1]}
			res := &models.User{}
			if err := session.Req().Put(pwd).UrlFor("users", args[0], "password").Do(res); err != nil {
				return generateError(err, "Error: putUserPassword: %v", err)
			}
			return prettyPrint(res)
		},
	})
	op.addCommand(&cobra.Command{
		Use:   "passwordhash [password]",
		Short: "Get a password hash for a password",
		Long:  "Get a password hash for a password.  This can be used in content packages.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v needs 1 arg", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			res := &models.User{}
			if err := res.ChangePassword(args[0]); err != nil {
				return generateError(err, "Error: generating password: %v", err)
			}
			fmt.Printf("%s\n", string(res.PasswordHash))
			return nil
		},
	})
	tokenArgs := []string{}
	op.addCommand(&cobra.Command{
		Use:   "token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]]",
		Short: "Get a login token for this user with optional parameters",
		Long:  "Creates a time-bound token for the specified user.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v needs at least 1 arg", c.UseLine())
			}
			tokenArgs = []string{}
			index := 2
			for index < len(args) {
				switch args[index-1] {
				case "scope", "action", "specific", "ttl":
					tokenArgs = append(tokenArgs, args[index-1], args[index])
					index += 2
				default:
					return fmt.Errorf("%v does not support %s", c.UseLine(), args[index-1])
				}
			}
			if index-1 != len(args) {
				return fmt.Errorf("%v needs at least 1 and pairs arg", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			res := &models.UserToken{}
			if err := session.Req().UrlFor("users", args[0], "token").Params(tokenArgs...).Do(res); err != nil {
				return generateError(err, "Error: getToken: %v", err)
			}
			return prettyPrint(res)
		},
	})
	op.command(app)
}
