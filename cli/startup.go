package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	v4 "github.com/digitalrebar/provision/v4"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

type registerSection func(*cobra.Command)

var (
	version              = v4.RSVersion
	debug                = false
	catalog              = "https://repo.rackn.io"
	defaultCatalog       = "https://repo.rackn.io"
	endpoint             = "https://127.0.0.1:8092"
	defaultEndpoint      = "https://127.0.0.1:8092"
	defaultEndpoints     = []string{"https://127.0.0.1:8092"}
	token                = ""
	defaultToken         = ""
	username             = "rocketskates"
	defaultUsername      = "rocketskates"
	password             = "r0cketsk8ts"
	defaultPassword      = "r0cketsk8ts"
	downloadProxy        = ""
	defaultDownloadProxy = ""
	format               = ""
	defaultFormat        = "json"
	noColor              = false
	defaultNoColor       = false
	// Format is idx=val1,val2,val3;idx2=val1,val2, ...
	// idx = 0 = json string color
	// idx = 1 = json bool color
	// idx = 2 = json number color
	// idx = 3 = json null color
	// idx = 4 = json key color
	// idx = 5 = header color
	// idx = 6 = border color
	// idx = 7 = value1 color
	// idx = 8 = value2 color
	colorString           = "0=32;1=33;2=36;3=90;4=34,1;5=35,6=95;7=32;8=92"
	defaultColorString    = "0=32;1=33;2=36;3=90;4=34,1;5=35;6=95;7=32;8=92"
	printFields           = ""
	defaultPrintFields    = ""
	defaultTruncateLength = 40
	truncateLength        = 40
	noHeader              = false
	defaultNoHeader       = false
	// Session is the global client access session
	Session         *api.Client
	noToken         = false
	force           = false
	ref             = ""
	defaultRef      = ""
	trace           = ""
	traceToken      = ""
	defaultUrlProxy = ""
	urlProxy        = ""
	registrations   = []registerSection{}
	objectErrorsAreFatal = false
)

func addRegistrar(rs registerSection) {
	registrations = append(registrations, rs)
}

var ppr = func(c *cobra.Command, a []string) error {
	c.SilenceUsage = true
	if Session == nil {
		epInList := false
		for i := range defaultEndpoints {
			if defaultEndpoints[i] == endpoint {
				epInList = true
				break
			}
		}
		if !epInList {
			l := len(defaultEndpoints)
			defaultEndpoints = append(defaultEndpoints, endpoint)
			defaultEndpoints[0], defaultEndpoints[l] = defaultEndpoints[l], defaultEndpoints[0]
		}
		var sessErr error
		for _, endpoint = range defaultEndpoints {
			if token != "" {
				Session, sessErr = api.TokenSession(endpoint, token)
			} else {
				home := os.ExpandEnv("${HOME}")
				tPath := os.ExpandEnv("${RS_TOKEN_CACHE}")
				if tPath == "" && home != "" {
					tPath = path.Join(home, ".cache", "drpcli", "tokens")
				}
				tokenFile := path.Join(tPath, "."+username+".token")
				if !noToken && tPath != "" {
					if err := os.MkdirAll(tPath, 0700); err == nil {
						if tokenStr, err := ioutil.ReadFile(tokenFile); err == nil {
							Session, sessErr = api.TokenSession(endpoint, string(tokenStr))
							if sessErr == nil {
								if _, err := Session.Info(); err == nil {
									Session.Trace(trace)
									Session.TraceToken(traceToken)
									break
								}
							}
						}
					}
				}
				Session, sessErr = api.UserSessionToken(endpoint, username, password, !noToken)
				if !noToken && tPath != "" && sessErr == nil {
					if err := os.MkdirAll(tPath, 700); err == nil {
						tok := &models.UserToken{}
						if err := Session.
							Req().UrlFor("users", username, "token").
							Params("ttl", "7200").Do(&tok); err == nil {
							ioutil.WriteFile(tokenFile, []byte(tok.Token), 0600)
						}
					}
				}
			}
			if sessErr == nil {
				break
			}
		}
		if sessErr != nil {
			return fmt.Errorf("Error creating Session: %v", sessErr)
		}
	}
	// We have a session.
	Session.UrlProxy(urlProxy)
	Session.Trace(trace)
	Session.TraceToken(traceToken)
	return nil
}

// NewApp is the app start function
func NewApp() *cobra.Command {
	// Don't use color on windows
	if runtime.GOOS == "windows" {
		defaultNoColor = true
	}
	app := &cobra.Command{
		Use:   "drpcli",
		Short: "A CLI application for interacting with the DigitalRebar Provision API",
	}
	if oeaf := os.Getenv("RS_OBJECT_ERRORS_ARE_FATAL"); oeaf == "true" {
		objectErrorsAreFatal = true
	}
	if dep := os.Getenv("RS_ENDPOINTS"); dep != "" {
		defaultEndpoints = strings.Split(dep, " ")
	}
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		defaultEndpoints = []string{ep}
	}
	if ep := os.Getenv("RS_URL_PROXY"); ep != "" {
		defaultUrlProxy = ep
	}
	if tk := os.Getenv("RS_TOKEN"); tk != "" {
		defaultToken = tk
	}
	if tk := os.Getenv("RS_CATALOG"); tk != "" {
		defaultCatalog = tk
	}
	if tk := os.Getenv("RS_FORMAT"); tk != "" {
		defaultFormat = tk
	}
	if tk := os.Getenv("RS_PRINT_FIELDS"); tk != "" {
		defaultPrintFields = tk
	}
	if tk := os.Getenv("RS_DOWNLOAD_PROXY"); tk != "" {
		defaultDownloadProxy = tk
	}
	if tk := os.Getenv("RS_NO_HEADER"); tk != "" {
		var e error
		defaultNoHeader, e = strconv.ParseBool(tk)
		if e != nil {
			log.Fatal("RS_NO_HEADER should be a boolean value")
		}
	}
	if tk := os.Getenv("RS_NO_COLOR"); tk != "" {
		var e error
		defaultNoColor, e = strconv.ParseBool(tk)
		if e != nil {
			log.Fatal("RS_NO_COLOR should be a boolean value")
		}
	}
	if tk := os.Getenv("RS_COLORS"); tk != "" {
		defaultColorString = tk
	}
	if tk := os.Getenv("RS_TRUNCATE_LENGTH"); tk != "" {
		var e error
		defaultTruncateLength, e = strconv.Atoi(tk)
		if e != nil {
			log.Fatal("RS_TRUNCATE_LENGTH should be an integer value")
		}
	}
	if kv := os.Getenv("RS_KEY"); kv != "" {
		key := strings.SplitN(kv, ":", 2)
		if len(key) < 2 {
			log.Fatal("RS_KEY does not contain a username:password pair!")
		}
		if key[0] == "" || key[1] == "" {
			log.Fatal("RS_KEY contains an invalid username:password pair!")
		}
		defaultUsername = key[0]
		defaultPassword = key[1]
	}
	home := os.ExpandEnv("${HOME}")
	if data, err := ioutil.ReadFile(fmt.Sprintf("%s/.drpclirc", home)); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			parts := strings.SplitN(line, "=", 2)

			switch parts[0] {
			case "RS_NO_HEADER":
				var e error
				defaultNoHeader, e = strconv.ParseBool(parts[1])
				if e != nil {
					log.Fatal("RS_NO_HEADER should be a boolean value in drpclirc")
				}
			case "RS_NO_COLOR":
				var e error
				defaultNoColor, e = strconv.ParseBool(parts[1])
				if e != nil {
					log.Fatal("RS_NO_HEADER should be a boolean value in drpclirc")
				}
			case "RS_COLORS":
				defaultColorString = parts[1]
			case "RS_ENDPOINT":
				defaultEndpoints = []string{parts[1]}
			case "RS_TOKEN":
				defaultToken = parts[1]
			case "RS_USERNAME":
				defaultUsername = parts[1]
			case "RS_PASSWORD":
				defaultPassword = parts[1]
			case "RS_URL_PROXY":
				defaultUrlProxy = parts[1]
			case "RS_DOWNLOAD_PROXY":
				defaultDownloadProxy = parts[1]
			case "RS_FORMAT":
				defaultFormat = parts[1]
			case "RS_PRINT_FIELDS":
				defaultPrintFields = parts[1]
			case "RS_TRUNCATE_LENGTH":
				var e error
				defaultTruncateLength, e = strconv.Atoi(parts[1])
				if e != nil {
					log.Fatal("RS_TRUNCATE_LENGTH should be an integer value in drpclirc")
				}
			case "RS_KEY":
				key := strings.SplitN(parts[1], ":", 2)
				if len(key) < 2 {
					continue
				}
				if key[0] == "" || key[1] == "" {
					continue
				}
				defaultUsername = key[0]
				defaultPassword = key[1]
			}
		}
	}
	app.PersistentFlags().StringVarP(&endpoint,
		"endpoint", "E", defaultEndpoints[0],
		"The Digital Rebar Provision API endpoint to talk to")
	app.PersistentFlags().StringVarP(&username,
		"username", "U", defaultUsername,
		"Name of the Digital Rebar Provision user to talk to")
	app.PersistentFlags().StringVarP(&password,
		"password", "P", defaultPassword,
		"password of the Digital Rebar Provision user")
	app.PersistentFlags().StringVarP(&token,
		"token", "T", defaultToken,
		"token of the Digital Rebar Provision access")
	app.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")
	app.PersistentFlags().BoolVarP(&noColor,
		"no-color", "N", defaultNoColor,
		"Whether the CLI should output colorized strings")
	app.PersistentFlags().StringVarP(&colorString,
		"colors", "C", defaultColorString,
		`The colors for JSON and Table/Text colorization.  8 values in the for 0=val,val;1=val,val2...`)
	app.PersistentFlags().StringVarP(&format,
		"format", "F", defaultFormat,
		`The serialization we expect for output.  Can be "json" or "yaml" or "text" or "table"`)
	app.PersistentFlags().StringVarP(&printFields,
		"print-fields", "J", defaultPrintFields,
		`The fields of the object to display in "text" or "table" mode. Comma separated`)
	app.PersistentFlags().BoolVarP(&noHeader,
		"no-header", "H", defaultNoHeader,
		`Should header be shown in "text" or "table" mode`)
	app.PersistentFlags().IntVarP(&truncateLength,
		"truncate-length", "j", defaultTruncateLength,
		`Truncate columns at this length`)
	app.PersistentFlags().BoolVarP(&force,
		"force", "f", false,
		"When needed, attempt to force the operation - used on some update/patch calls")
	app.PersistentFlags().StringVarP(&ref,
		"ref", "r", defaultRef,
		"A reference object for update commands that can be a file name, yaml, or json blob")
	app.PersistentFlags().StringVarP(&trace,
		"trace", "t", "",
		"The log level API requests should be logged at on the server side")
	app.PersistentFlags().StringVarP(&traceToken,
		"traceToken", "Z", "",
		"A token that individual traced requests should report in the server logs")
	app.PersistentFlags().StringVarP(&catalog,
		"catalog", "c", defaultCatalog,
		"The catalog file to use to get product information")
	app.PersistentFlags().StringVarP(&downloadProxy,
		"download-proxy", "D", defaultDownloadProxy,
		"HTTP Proxy to use for downloading catalog and content")
	app.PersistentFlags().BoolVarP(&noToken,
		"noToken", "x", noToken,
		"Do not use token auth or token cache")
	app.PersistentFlags().BoolVarP(&objectErrorsAreFatal,
		"exitEarly", "X", false,
		"Cause drpcli to exit if a command results in an object that has errors")
	app.PersistentFlags().StringVarP(&urlProxy,
		"url-proxy", "u", defaultUrlProxy,
		"URL Proxy for passing actions through another DRP")
	if runtime.GOOS != "windows" {
		app.AddCommand(&cobra.Command{
			Use:   "proxy [socket]",
			Short: "Run a local UNIX socket proxy for further drpcli commands.  Requires RS_LOCAL_PROXY to be set in the env.",
			RunE: func(c *cobra.Command, args []string) error {
				if len(args) != 1 {
					return fmt.Errorf("No location for the local proxy socket")
				}
				if pl := os.Getenv("RS_LOCAL_PROXY"); pl != "" {
					return fmt.Errorf("Local proxy already running at %s", pl)
				}
				return Session.RunProxy(args[0])
			},
		})
	}

	for _, rs := range registrations {
		rs(app)
	}

	for _, c := range app.Commands() {
		// contents needs some help.
		switch c.Use {
		case "catalog":
			for _, sc := range c.Commands() {
				if strings.HasPrefix(sc.Use, "updateLocal") {
					sc.PersistentPreRunE = ppr
				}
			}
		case "contents":
			for _, sc := range c.Commands() {
				if !strings.HasPrefix(sc.Use, "bundle") &&
					!strings.HasPrefix(sc.Use, "unbundle") &&
					!strings.HasPrefix(sc.Use, "document") {
					sc.PersistentPreRunE = ppr
				}
			}
		case "users":
			for _, sc := range c.Commands() {
				if !strings.HasPrefix(sc.Use, "passwordhash") {
					sc.PersistentPreRunE = ppr
				}
			}
		case "support":
			for _, sc := range c.Commands() {
				if !strings.HasPrefix(sc.Use, "bundle") {
					sc.PersistentPreRunE = ppr
				}
			}

		default:
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

	app.AddCommand(&cobra.Command{
		Use:   "fingerprint",
		Short: "Get the machine fingerprint used to determine what machine we are running on",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			w := &models.Whoami{}
			if err := w.Fill(); err != nil {
				return err
			}
			return prettyPrint(w)
		},
	})
	app.AddCommand(agentHandler)

	return app
}

func ResetDefaults() {
	defaultEndpoint = "https://127.0.0.1:8092"
	defaultEndpoints = []string{"https://127.0.0.1:8092"}
}
