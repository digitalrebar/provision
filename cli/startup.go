package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/api"
	"github.com/spf13/cobra"
)

type registerSection func(*cobra.Command)

var (
	version          = provision.RS_VERSION
	debug            = false
	endpoint         = "https://127.0.0.1:8092"
	default_endpoint = "https://127.0.0.1:8092"
	token            = ""
	default_token    = ""
	username         = "rocketskates"
	default_username = "rocketskates"
	password         = "r0cketsk8ts"
	default_password = "r0cketsk8ts"
	format           = "json"
	session          *api.Client
	force            = false
	noPretty         = false
	ref              = ""
	default_ref      = ""
	trace            = ""
	traceToken       = ""
	registrations    = []registerSection{}
)

func addRegistrar(rs registerSection) {
	registrations = append(registrations, rs)
}

var ppr = func(c *cobra.Command, a []string) {
	c.SilenceUsage = true
	if session == nil {
		var err error
		if token != "" {
			session, err = api.TokenSession(endpoint, token)
		} else {
			session, err = api.UserSession(endpoint, username, password)
		}
		if err != nil {
			log.Printf("Error creating session: %v", err)
			os.Exit(1)
		}
	}
	session.Trace(trace)
	session.TraceToken(traceToken)
}

func NewApp() *cobra.Command {
	app := &cobra.Command{
		Use:   "drpcli",
		Short: "A CLI application for interacting with the DigitalRebar Provision API",
	}
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		default_endpoint = ep
	}
	if tk := os.Getenv("RS_TOKEN"); tk != "" {
		default_token = tk
	}
	if kv := os.Getenv("RS_KEY"); kv != "" {
		key := strings.SplitN(kv, ":", 2)
		if len(key) < 2 {
			log.Fatal("RS_KEY does not contain a username:password pair!")
		}
		if key[0] == "" || key[1] == "" {
			log.Fatal("RS_KEY contains an invalid username:password pair!")
		}
		default_username = key[0]
		default_password = key[1]
	}
	app.PersistentFlags().StringVarP(&endpoint,
		"endpoint", "E", default_endpoint,
		"The Digital Rebar Provision API endpoint to talk to")
	app.PersistentFlags().StringVarP(&username,
		"username", "U", default_username,
		"Name of the Digital Rebar Provision user to talk to")
	app.PersistentFlags().StringVarP(&password,
		"password", "P", default_password,
		"password of the Digital Rebar Provision user")
	app.PersistentFlags().StringVarP(&token,
		"token", "T", default_token,
		"token of the Digital Rebar Provision access")
	app.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")
	app.PersistentFlags().StringVarP(&format,
		"format", "F", "json",
		`The serialzation we expect for output.  Can be "json" or "yaml"`)
	app.PersistentFlags().BoolVarP(&force,
		"force", "f", false,
		"When needed, attempt to force the operation - used on some update/patch calls")
	app.PersistentFlags().StringVarP(&ref,
		"ref", "r", default_ref,
		"A reference object for update commands that can be a file name, yaml, or json blob")
	app.PersistentFlags().StringVarP(&trace,
		"trace", "t", "",
		"The log level API requests should be logged at on the server side")
	app.PersistentFlags().StringVarP(&traceToken,
		"traceToken", "Z", "",
		"A token that individual traced requests should report in the server logs")

	for _, rs := range registrations {
		rs(app)
	}

	for _, c := range app.Commands() {
		c.PersistentPreRun = ppr
	}
	// top-level commands that do not need PersistentPreRun go here.
	app.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Digital Rebar Provision CLI Command Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %v\n", version)
			return nil
		},
	})
	app.AddCommand(&cobra.Command{
		Use:   "autocomplete <filename>",
		Short: "Digital Rebar Provision CLI Command Bash AutoCompletion File",
		Long:  "Generate a bash autocomplete file as <filename>.\nPlace the generated file in /etc/bash_completion.d or /usr/local/etc/bash_completion.d.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1  argument", cmd.UseLine())
			}
			app.GenBashCompletionFile(args[0])
			return nil
		},
	})

	return app
}
