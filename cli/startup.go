package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

type registerSection func(*cobra.Command)

var (
	version          = provision.RSVersion
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

var ppr = func(c *cobra.Command, a []string) error {
	c.SilenceUsage = true
	if session == nil {
		var err error
		if token != "" {
			session, err = api.TokenSession(endpoint, token)
		} else {
			home := os.ExpandEnv("${HOME}")
			tPath := os.ExpandEnv("${RS_TOKEN_CACHE}")
			if tPath == "" && home != "" {
				tPath = path.Join(home, ".cache", "drpcli", "tokens")
			}
			tokenFile := path.Join(tPath, "."+username+".token")
			if tPath != "" {
				if err := os.MkdirAll(tPath, 0700); err == nil {
					if tokenStr, err := ioutil.ReadFile(tokenFile); err == nil {
						session, err = api.TokenSession(endpoint, string(tokenStr))
						if err == nil {
							if _, err := session.Info(); err == nil {
								session.Trace(trace)
								session.TraceToken(traceToken)
								return nil
							}
							session.Close()
							session = nil
						}
					}
				}
			}
			session, err = api.UserSession(endpoint, username, password)
			if tPath != "" && err == nil {
				if err := os.MkdirAll(tPath, 700); err == nil {
					tok := &models.UserToken{}
					if err := session.
						Req().UrlFor("users", username, "token").
						Params("ttl", "7200").Do(&tok); err == nil {
						ioutil.WriteFile(tokenFile, []byte(tok.Token), 0600)
					}
				}
			}
		}
		if err != nil {
			return fmt.Errorf("Error creating session: %v", err)
		}
	}
	session.Trace(trace)
	session.TraceToken(traceToken)
	return nil
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
	home := os.ExpandEnv("${HOME}")
	if data, err := ioutil.ReadFile(fmt.Sprintf("%s/.drpclirc", home)); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			parts := strings.SplitN(line, "=", 2)

			switch parts[0] {
			case "RS_ENDPOINT":
				default_endpoint = parts[1]
			case "RS_TOKEN":
				default_token = parts[1]
			case "RS_USERNAME":
				default_username = parts[1]
			case "RS_PASSWORD":
				default_password = parts[1]
			case "RS_KEY":
				key := strings.SplitN(parts[1], ":", 2)
				if len(key) < 2 {
					continue
				}
				if key[0] == "" || key[1] == "" {
					continue
				}
				default_username = key[0]
				default_password = key[1]
			}
		}
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
		// contents needs some help.
		if c.Use == "contents" {
			for _, sc := range c.Commands() {
				if !strings.HasPrefix(sc.Use, "bundle") &&
					!strings.HasPrefix(sc.Use, "unbundle") &&
					!strings.HasPrefix(sc.Use, "document") {
					sc.PersistentPreRunE = ppr
				}
			}
		} else if c.Use == "users" {
			for _, sc := range c.Commands() {
				if !strings.HasPrefix(sc.Use, "passwordhash") {
					sc.PersistentPreRunE = ppr
				}
			}
		} else {
			c.PersistentPreRunE = ppr
		}
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
		Use:   "autocomplete [filename]",
		Short: "Generate CLI Command Bash AutoCompletion File (may require 'bash-completion' pkg be installed)",
		Long:  "Generate a bash autocomplete file as *filename*.\nPlace the generated file in /etc/bash_completion.d or /usr/local/etc/bash_completion.d.\nMay require the 'bash-completion' package is installed to work correctly.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1  argument", cmd.UseLine())
			}
			app.GenBashCompletionFile(args[0])
			return nil
		},
	})

	app.AddCommand(&cobra.Command{
		Use:   "gohai",
		Short: "Get basic system information as a JSON blob",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return gohai()
		},
	})

	return app
}
