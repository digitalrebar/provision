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
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type PluginStop interface {
	Stop(logger.Logger)
}

type PluginConfig interface {
	Config(logger.Logger, *api.Client, map[string]interface{}) *models.Error
}

type PluginPublisher interface {
	Publish(logger.Logger, *models.Event) *models.Error
}

type PluginActor interface {
	Action(logger.Logger, *models.Action) (interface{}, *models.Error)
}

var (
	thelog logger.Logger
	App    = &cobra.Command{
		Use:   "replaceme",
		Short: "Replace ME!",
	}
	debug   = false
	client  *http.Client
	session *api.Client
)

func Publish(t, a, k string, o interface{}) {
	e := &models.Event{Time: time.Now(), Type: t, Action: a, Key: k, Object: o}
	_, err := post("/publish", e)
	if err != nil {
		thelog.Errorf("Failed to publish event! %v %v", e, err)
	}
}

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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1  argument", cmd.UseLine())
			}
			App.GenBashCompletionFile(args[0])
			return nil
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "define",
		Short: "Digital Rebar Provision CLI Command Define",
		RunE: func(cmd *cobra.Command, args []string) error {
			if buf, err := json.MarshalIndent(def, "", "  "); err == nil {
				fmt.Println(string(buf))
				return nil
			} else {
				return err
			}
		},
	})
	App.AddCommand(&cobra.Command{
		Use:   "listen <socket path to plugin> <socket path from plugin>",
		Short: "Digital Rebar Provision CLI Command Listen",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				fmt.Printf("Failed\n")
				return fmt.Errorf("%v requires 2 argument", cmd.UseLine())
			}

			return Run(args[0], args[1], pc)
		},
	})
}

func Run(toPath, fromPath string, pc PluginConfig) error {
	// Get HTTP2 client on our socket.
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", fromPath)
			},
		},
	}

	gc := newGinServer(thelog.NoPublish())
	gc.NoMethod(func(c *gin.Context) { thelog.Warnf("Unknown method: %v\n", c) })
	gc.NoRoute(func(c *gin.Context) { thelog.Warnf("Unknown route: %v\n", c) })
	apiGroup := gc.Group("/api-plugin/v3")

	// Required Pieces.
	apiGroup.POST("/config", func(c *gin.Context) { configHandler(c, pc) })
	if ps, ok := pc.(PluginStop); ok {
		apiGroup.POST("/stop", func(c *gin.Context) { stopHandler(c, ps) })
	} else {
		apiGroup.POST("/stop", func(c *gin.Context) { stopHandler(c, nil) })
	}

	// Optional Pieces
	if pp, ok := pc.(PluginPublisher); ok {
		apiGroup.POST("/publish", func(c *gin.Context) { publishHandler(c, pp) })
	}
	if pa, ok := pc.(PluginActor); ok {
		apiGroup.POST("/action", func(c *gin.Context) { actionHandler(c, pa) })
	}

	go func() {
		for {
			if _, err := os.Stat(toPath); os.IsNotExist(err) {
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}
		fmt.Printf("READY!\n")
	}()
	return gc.RunUnix(toPath)
}

func logToDRP(l *logger.Line) {
	_, err := post("/log", l)
	if err != nil {
		thelog.NoPublish().Errorf("Failed to log line! %v %v", l, err)
	}
}

func stopHandler(c *gin.Context, ps PluginStop) {
	// GREG: Do better?
	if ps != nil {
		ps.Stop(thelog)
	}
	resp := models.Error{Code: http.StatusOK}
	c.JSON(resp.Code, resp)
	thelog.Infof("STOPPING\n")
	os.Exit(0)
}

func configHandler(c *gin.Context, pc PluginConfig) {
	var params map[string]interface{}
	if !assureDecode(c, &params) {
		return
	}

	thelog.Infof("Setting API session\n")

	default_endpoint := "https://127.0.0.1:8092"
	if ep := os.Getenv("RS_ENDPOINT"); ep != "" {
		default_endpoint = ep
	}
	default_token := ""
	if tk := os.Getenv("RS_TOKEN"); tk != "" {
		default_token = tk
	}

	var err2 error
	if default_token != "" {
		thelog.Infof("Starting session with endpoint and token: %s\n", default_endpoint)
		session, err2 = api.TokenSession(default_endpoint, default_token)
	} else {
		err2 = fmt.Errorf("Must have a token specified\n")
	}

	if err2 != nil {
		err := &models.Error{Code: 400, Model: "plugin", Key: "incrementer", Type: "plugin", Messages: []string{}}
		err.AddError(err2)
		c.JSON(err.Code, err)
		return
	}

	thelog.Infof("Received Config request: %v\n", params)
	resp := models.Error{Code: http.StatusOK}
	if err := pc.Config(thelog, session, params); err != nil {
		resp.Code = err.Code
		b, _ := json.Marshal(err)
		resp.Messages = append(resp.Messages, string(b))
	}
	c.JSON(resp.Code, resp)
}

func actionHandler(c *gin.Context, pa PluginActor) {
	var actionInfo models.Action
	if !assureDecode(c, &actionInfo) {
		return
	}
	if ret, err := pa.Action(thelog, &actionInfo); err != nil {
		c.JSON(err.Code, err)
	} else {
		c.JSON(200, ret)
	}
}

func publishHandler(c *gin.Context, pp PluginPublisher) {
	var event models.Event
	if !assureDecode(c, &event) {
		return
	}
	resp := models.Error{Code: http.StatusOK}
	if err := pp.Publish(thelog.NoPublish(), &event); err != nil {
		resp.Code = err.Code
		b, _ := json.Marshal(err)
		resp.Messages = append(resp.Messages, string(b))
	}
	c.JSON(resp.Code, resp)
}
