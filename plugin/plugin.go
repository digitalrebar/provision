// Package plugin is used to write plugin providers in Go.
// It provides the framework for the rest of the plugin provider code,
// along with a set of interfaces that you app can satisfy to
// implement whatever custom behaviour the plugin provider needs to implement.
//
// A plugin provider is an executable that provides extensions to the base
// functionality in dr-provision.  This can be done in several different ways,
// depending on the functionality the plugin provider needs to implement:
//
// 1. Injecting custom content bundles into dr-provision to provide
//    additional tasks, params, etc.
//
// 2. Implementing additional per-object Actions that can be used for
//    a wide variety of things.
//
// 3. Providing additional files in the files/ space of the static file server.
//
// 4. Listening to the event stream from dr-provision to take action
//    whenever any number of selected events happen.
//
// 5. Define new object types that dr-provision will store and manage.
//
// github.com/digitalrebar/provision/cmds/incrementer provides a
// fully functional implementation of a basic plugin provider that
// you can use as an example and as a base for implementing your own
// plugin providers.
//
// github.com/digitalrebar/provision-plugins contains several production
// ready plugin provider implementations that you can use as a reference
// for implementing more advanced behaviours.
//
// At a higher level, a plugin provider is an application that has 3 ways
// of being invoked:
//
// 1. plugin_provider define
//
//    When invoked with a single argument of define, the plugin provider must
//    print the models.PluginProvider definition for the plugin provider in
//    JSON format on stdout.
//
// 2. plugin_provider unpack /path/to/filespace/for/this/provider
//
//    When invoked with unpack /path, the plugin provider must unpack any
//    embedded assets (extra executables and other artifacts like that) into
//    the path passed in as the argument.  Note that this does not include
//    the embedded content pack, which is emitted as part of the define
//    command.
//
// 3. plugin_provider listen /path/to/client/socket /path/to/server/socket
//
//    When invoked with listen, the plugin client must open an HTTP client
//    connection on the client socket to post events and status updates back
//    to dr-provision, and listen with an HTTP server on the server socket
//    to receive action requests, stop requests, and events from dr-provision.
//    Once both sockets are opened up and the plugin provider is ready to
//    be configured, it should emit `READY!` followed by a newline
//    on stdout.
//
//    In all cases, the following environment variables will be set when
//    the plugin provider is executed:
//
//    RS_ENDPOINT will be a URL to the usual dr-provision API endpoint
//    RS_TOKEN will be a long-lived token with superuser access rights
//    RS_FILESERVER will be a URL to the static file server
//    RS_WEBROOT will be the filesystem path to static file server space
//
//    The plugin provider will be executed with its current directory set
//    to a scratch directory it can use to hold temporary files.
//
// Once the plugin provider is ready, its HTTP server should listen on
// the following paths:
//
//  POST /api-plugin/v4/config
//
//  When a JSON object containing the Params field from the Plugin object
//  this instance of the plugin provider is backing is POSTed to this API
//  endpoint, the plugin should configure itself accordingly.
//  This is the first call made into the plugin provider
//  when it starts, and it can be called any time afterwards.
//
//  POST /api-plugin/v4/stop
//
//  When this API endpoint is POSTed to, the plugin provider should cleanly
//  shut down.
//
//  POST /api-plugin/v4/action
//
//  When a JSON object containing a fully filled out models.Action is POSTed
//  to this API endpoint, the plugin provider should take the appropriate
//  action and return the results of the action.  This endpoint must be
//  able to handle all of the actions listed in the AvailableActions section
//  of the definition that the define command returned.
//
//  POST /api-plugin/v4/publish (DEPRECATED, use api.EventStream instead)
//
//  When a JSON object containing a fully filled out models.Event is POSTed
//  to this API endpoint, the plugin provider should handle the event as
//  appropriate.  Events will only be published to this endpoint if the
//  plugin provider definition HasPublish flag is true.
//
//  This endpoint is deprecated, as it is synchronous and can cause
//  performance bottlenecks and potentially deadlocks, along with not
//  being filterable on the server side.  Using an api.EventStream
//  is a better solution.
//
// The HTTP client can POST back into dr-provision using the following
// paths on the client socket:
//
//  POST /api-plugin-server/v4/publish
//
//  The body should be a JSON serialized models.Event, which will be broadcast
//  to all interested parties.
//
//  POST /api-plugin-server/v4/leaving
//
//  This will cause dr-provision to cleanly shut down the plugin provider.
//  The body does not matter.
//
//  POST /api-plugin-server/v4/log
//
//  The body should be a JSON serialized logger.Line structure, which will be
//  added to the global dr-provision log.
package plugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/digitalrebar/provision/v4/plugin/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

var json = jsoniter.ConfigFastest

var (
	thelog logger.Logger
	// App is the global cobra command structure.
	App = &cobra.Command{
		Use:   "replaceme",
		Short: "Replace ME!",
	}
	debug    = false
	client   *http.Client
	session  *api.Client
	es       *api.EventStream
	esHandle int64
	events   <-chan api.RecievedEvent
)

// Publish allows the plugin provider to generate events back to DRP.
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
// available in cobra CLI.  It provides default implementations of the
// define, unpack, and listen commands, which will be backed by all the interfaces
// that whatever is passed in as pc satisfy.
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
			return run(args[0], args[1], pc, def)
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
func run(toPath, fromPath string, pc PluginConfig, def *models.PluginProvider) error {
	// Get HTTP2 client on our socket.
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", fromPath)
			},
		},
	}
	pmux := mux.New(thelog)
	pmux.Handle("/api-plugin/v4/config",
		func(w http.ResponseWriter, r *http.Request) { configHandler(w, r, def, pc) })
	if ps, ok := pc.(PluginStop); ok {
		pmux.Handle("/api-plugin/v4/stop",
			func(w http.ResponseWriter, r *http.Request) { stopHandler(w, r, ps) })
	} else {
		pmux.Handle("/api-plugin/v4/stop",
			func(w http.ResponseWriter, r *http.Request) { stopHandler(w, r, nil) })
	}

	// Optional Pieces
	_, hasPSE := pc.(PluginEventSelecter)
	if pp, ok := pc.(PluginPublisher); ok && def.HasPublish && !hasPSE {
		pmux.Handle("/api-plugin/v4/publish",
			func(w http.ResponseWriter, r *http.Request) { publishHandler(w, r, pp) })
	}
	if pa, ok := pc.(PluginActor); ok {
		pmux.Handle("/api-plugin/v4/action",
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
	if es != nil {
		es.Deregister(esHandle)
		es.Close()
	}
	if r.Body != nil {
		r.Body.Close()
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

func configHandler(w http.ResponseWriter, r *http.Request, def *models.PluginProvider, pc PluginConfig) {
	var params map[string]interface{}
	if !mux.AssureDecode(w, r, &params) {
		return
	}
	l := w.(logger.Logger)

	l.Infof("Setting API session\n")
	session, err2 := buildSession()
	if err2 != nil {
		err := &models.Error{Code: 400, Model: "plugin", Key: "unknown", Type: "plugin", Messages: []string{}}
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
		mux.JsonResponse(w, resp.Code, resp)
		return
	}

	// We need to handle a few cases.
	// Does the plugin have the following:
	//    eventSelector
	//    publish Function
	//    marked HasPublish.
	//
	// The plugin should either HasPublish or use eventSelector
	// if both selector HasPublish and EventSelector - panic
	pse, hasPSE := pc.(PluginEventSelecter)
	pp, hasPublishFunct := pc.(PluginPublisher)
	hasPublish := def.HasPublish
	if hasPublish && hasPSE {
		err := &models.Error{Code: 400, Model: "plugin", Key: "unknown", Type: "plugin", Messages: []string{}}
		err.Errorf("Plugin can NOT have both HasPublish and EventSelector: %s", def.Name)
		mux.JsonResponse(w, err.Code, err)
		return
	}

	if hasPSE && !hasPublishFunct {
		err := &models.Error{Code: 400, Model: "plugin", Key: "unknown", Type: "plugin", Messages: []string{}}
		err.Errorf("Plugin can has an EventSelector, but no Publish function: %s", def.Name)
		mux.JsonResponse(w, err.Code, err)
		return
	}

	if hasPSE && hasPublishFunct {
		var esErr error
		if es != nil {
			es.Deregister(esHandle)
			es.Close()
		}
		es, esErr = session.Events()
		if esErr != nil {
			err := models.NewError("plugins", 500, fmt.Sprintf("Unable to create event stream: %v", esErr))
			mux.JsonResponse(w, err.Code, err)
			return
		}
		esHandle, events, esErr = es.Register(pse.SelectEvents()...)
		if esErr != nil {
			es.Close()
			err := models.NewError("plugins", 500, fmt.Sprintf("Unable to register for machine events: %v", esErr))
			mux.JsonResponse(w, err.Code, err)
			return
		}
		go func(l logger.Logger, eventStream <-chan api.RecievedEvent) {
			for {
				evt, ok := <-eventStream
				if !ok {
					return
				}
				if err := pp.Publish(l, &evt.E); err != nil {
					l.Errorf("Error processing event: %v", err)
				}
			}
		}(w.(logger.Logger).NoRepublish(), events)

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
