package plugin

/*
 * This is used by plugins to define their base App.
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/provision/plugin/mux"
	"github.com/spf13/cobra"
)

// PluginStop defines the Stop routine used to inform a plugin
// that it should stop.
type PluginStop interface {
	Stop(logger.Logger)
}

// PluginConfig defines the Config routine used to configure a
// specific instance of a plugin.
type PluginConfig interface {
	Config(logger.Logger, *api.Client, map[string]interface{}) *models.Error
}

// PluginPublisher defines the Publish routine used to send events
// to a plugin.
type PluginPublisher interface {
	Publish(logger.Logger, *models.Event) *models.Error
}

// PluginActor defines the Action routine used to invoke actions
// by the plugin.
type PluginActor interface {
	Action(logger.Logger, *models.Action) (interface{}, *models.Error)
}

// PluginValidator defines the Validate routine used to ensure that
// the environment is valid around the define timeframe.
type PluginValidator interface {
	Validate(logger.Logger, *api.Client) (interface{}, *models.Error)
}

// PluginUnpacker defines the Unpack routine used to unpack embedded
// assets to the specified path.
type PluginUnpacker interface {
	Unpack(logger.Logger, string) error
}

var (
	thelog logger.Logger
	// App is the global cobra command structure.
	App = &cobra.Command{
		Use:   "replaceme",
		Short: "Replace ME!",
	}
	debug   = false
	client  *http.Client
	session *api.Client
)

// Publish allows the plugin to generate events back to DRP.
func Publish(t, a, k string, o interface{}) {
	if client == nil {
		return
	}
	e := &models.Event{Time: time.Now(), Type: t, Action: a, Key: k, Object: o}
	_, err := mux.Post(client, "/publish", e)
	if err != nil {
		thelog.Errorf("Failed to publish event! %v %v", e, err)
	}
}

// Leaving allows the plugin to inform DRP that it is about to exit.
func Leaving(e *models.Error) {
	if client == nil {
		return
	}
	_, err := mux.Post(client, "/leaving", e)
	if err != nil {
		thelog.Errorf("Failed to send leaving event! %v %v", e, err)
	}
}

func ListObjects(prefix string) ([]*models.RawModel, *models.Error) {
	if client == nil {
		return nil, models.NewError("plugin-mux", 400, fmt.Sprintf("No client to look up %s", prefix))
	}
	data, err := mux.Get(client, fmt.Sprintf("/objects/%s", prefix))
	if err != nil {
		return nil, models.NewError("plugin-mux", 400, fmt.Sprintf("Failed to list %s: %v", prefix, err))
	}
	var m []*models.RawModel
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, models.NewError("plugin-mux", 400, fmt.Sprintf("Failed to marshal list %s: %v", prefix, err))
	}
	return m, nil
}

// InitApp initializes the plugin system and makes the base actions
// available in cobra CLI.
func InitApp(use, short, version string, def *models.PluginProvider, pc PluginConfig) {
	App.Use = use
	App.Short = short

	localLogger := log.New(ioutil.Discard, App.Use, log.LstdFlags|log.Lmicroseconds|log.LUTC)
	thelog = logger.New(localLogger).SetDefaultLevel(logger.Debug).SetPublisher(logToDRP).Log(App.Use)

	App.PersistentFlags().BoolVarP(&debug,
		"debug", "d", false,
		"Whether the CLI should run in debug mode")

	App.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Digital Rebar Provision CLI Command Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %v\n", version)
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "autocomplete <filename>",
		Short: "Digital Rebar Provision CLI Command Bash AutoCompletion File",
		Long:  "Generate a bash autocomplete file as <filename>.\nPlace the generated file in /etc/bash_completion.d or /usr/local/etc/bash_completion.d.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1  argument", cmd.UseLine())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			App.GenBashCompletionFile(args[0])
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		SilenceUsage: true,
		Use:          "define",
		Short:        "Digital Rebar Provision CLI Command Define",
		Args: func(c *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var theDef interface{}
			defaultToken := ""
			if tk := os.Getenv("RS_TOKEN"); tk != "" {
				defaultToken = tk
			}
			if pv, ok := pc.(PluginValidator); ok && defaultToken != "catalog" {
				session, err2 := buildSession()
				if err2 != nil {
					return err2
				}
				ndef, err := pv.Validate(thelog, session)
				if err != nil {
					return err
				}
				theDef = ndef
			} else {
				theDef = def
			}
			buf, err := json.MarshalIndent(theDef, "", "  ")
			if err == nil {
				fmt.Println(string(buf))
				return nil
			}
			return err
		},
	})
	App.AddCommand(&cobra.Command{
		SilenceUsage: true,
		Use:          "listen <socket path to plugin> <socket path from plugin>",
		Short:        "Digital Rebar Provision CLI Command Listen",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				fmt.Printf("Failed\n")
				return fmt.Errorf("%v requires 2 argument", cmd.UseLine())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], args[1], pc)
		},
	})
	App.AddCommand(&cobra.Command{
		SilenceUsage: true,
		Use:          "unpack [loc]",
		Short:        "Unpack embedded static content to [loc]",
		Args: func(c *cobra.Command, args []string) error {
			if args[0] == `` {
				return fmt.Errorf("Not a valid location: `%s`", args[0])
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if pu, ok := pc.(PluginUnpacker); ok {
				if err := os.MkdirAll(args[0], 0755); err != nil {
					return err
				}
				return pu.Unpack(thelog, args[0])
			}
			return nil
		},
	})
}

// run implements the listen part of the CLI.
func run(toPath, fromPath string, pc PluginConfig) error {
	// Get HTTP2 client on our socket.
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", fromPath)
			},
		},
	}
	pmux := mux.New(thelog)
	pmux.Handle("/api-plugin/v3/config",
		func(w http.ResponseWriter, r *http.Request) { configHandler(w, r, pc) })
	if ps, ok := pc.(PluginStop); ok {
		pmux.Handle("/api-plugin/v3/stop",
			func(w http.ResponseWriter, r *http.Request) { stopHandler(w, r, ps) })
	} else {
		pmux.Handle("/api-plugin/v3/stop",
			func(w http.ResponseWriter, r *http.Request) { stopHandler(w, r, nil) })
	}

	// Optional Pieces
	if pp, ok := pc.(PluginPublisher); ok {
		pmux.Handle("/api-plugin/v3/publish",
			func(w http.ResponseWriter, r *http.Request) { publishHandler(w, r, pp) })
	}
	if pa, ok := pc.(PluginActor); ok {
		pmux.Handle("/api-plugin/v3/action",
			func(w http.ResponseWriter, r *http.Request) { actionHandler(w, r, pa) })
	}
	os.Remove(toPath)
	sock, err := net.Listen("unix", toPath)
	if err != nil {
		return err
	}
	defer sock.Close()
	go func() {
		fmt.Printf("READY!\n")
	}()
	return http.Serve(sock, pmux)
}

func logToDRP(l *logger.Line) {
	if client == nil {
		fmt.Fprintf(os.Stderr, "local log: %v\n", l)
		return
	}
	_, err := mux.Post(client, "/log", l)
	if err != nil {
		thelog.NoRepublish().Errorf("Failed to log line! %v %v", l, err)
	}
}

func stopHandler(w http.ResponseWriter, r *http.Request, ps PluginStop) {
	l := w.(logger.Logger)
	if ps != nil {
		ps.Stop(l)
	}
	resp := models.Error{Code: http.StatusOK}
	mux.JsonResponse(w, resp.Code, resp)
	l.Infof("STOPPING\n")
	os.Exit(0)
}

func buildSession() (*api.Client, error) {
	defaultEndpoint := "https://127.0.0.1:8092"
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		defaultEndpoint = ep
	}
	defaultToken := ""
	if tk := os.Getenv("RS_TOKEN"); tk != "" {
		defaultToken = tk
	}

	var session *api.Client
	var err2 error
	if defaultToken != "" {
		session, err2 = api.TokenSession(defaultEndpoint, defaultToken)
	} else {
		err2 = fmt.Errorf("Must have a token specified")
	}

	return session, err2
}

func configHandler(w http.ResponseWriter, r *http.Request, pc PluginConfig) {
	var params map[string]interface{}
	if !mux.AssureDecode(w, r, &params) {
		return
	}
	l := w.(logger.Logger)

	l.Infof("Setting API session\n")
	session, err2 := buildSession()
	if err2 != nil {
		err := &models.Error{Code: 400, Model: "plugin", Key: "incrementer", Type: "plugin", Messages: []string{}}
		err.AddError(err2)
		mux.JsonResponse(w, err.Code, err)
		return
	}

	l.Infof("Received Config request: %v\n", params)
	resp := models.Error{Code: http.StatusOK}
	if err := pc.Config(l, session, params); err != nil {
		resp.Code = err.Code
		b, _ := json.Marshal(err)
		resp.Messages = append(resp.Messages, string(b))
	}
	mux.JsonResponse(w, resp.Code, resp)
}

func actionHandler(w http.ResponseWriter, r *http.Request, pa PluginActor) {
	var actionInfo models.Action
	if !mux.AssureDecode(w, r, &actionInfo) {
		return
	}
	l := w.(logger.Logger)
	if ret, err := pa.Action(l, &actionInfo); err != nil {
		mux.JsonResponse(w, err.Code, err)
	} else {
		mux.JsonResponse(w, http.StatusOK, ret)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request, pp PluginPublisher) {
	var event models.Event
	if !mux.AssureDecode(w, r, &event) {
		return
	}
	l := w.(logger.Logger)
	resp := models.Error{Code: http.StatusOK}
	if err := pp.Publish(l.NoRepublish(), &event); err != nil {
		resp.Code = err.Code
		b, _ := json.Marshal(err)
		resp.Messages = append(resp.Messages, string(b))
	}
	mux.JsonResponse(w, resp.Code, resp)
}
