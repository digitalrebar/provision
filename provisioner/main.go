package provisioner

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/digitalrebar/digitalrebar/go/common/client"
	"github.com/digitalrebar/digitalrebar/go/common/service"
	"github.com/digitalrebar/digitalrebar/go/common/store"
	"github.com/digitalrebar/digitalrebar/go/common/version"
	"github.com/digitalrebar/digitalrebar/go/rebar-api/api"
	consul "github.com/hashicorp/consul/api"
)

var ProvOpts struct {
	VersionFlag    bool   `long:"version" description:"Print Version and exit"`
	BackEndType    string `long:"backend" description:"Storage backend to use. Can be either 'consul' or 'directory'" default:"consul"`
	DataRoot       string `long:"data-root" description:"Location we should store runtime information in" default:"digitalrebar/provisioner/boot-info"`
	StaticPort     int    `long:"static-port" description:"Port the static HTTP file server should listen on" default:"8091"`
	TftpPort       int    `long:"tftp-port" description:"Port for the TFTP server to listen on" default:"69"`
	FileRoot       string `long:"file-root" description:"Root of filesystem we should manage" default:"/tftpboot"`
	OurAddress     string `long:"static-ip" description:"IP address to advertise for the static HTTP file server" default:"192.168.124.11"`
	CommandURL     string `long:"endpoint" description:"DigitalRebar Endpoint" env:"EXTERNAL_REBAR_ENDPOINT"`
	RegisterConsul bool   `long:"register-consul" description:"Register services with Consul"`
}

var Logger *log.Logger

var ProvisionerURL string
var backends = map[string]store.SimpleStore{}
var backendMux = sync.Mutex{}
var rebarClient *api.Client

func InitializeProvisioner(apiPort int) {
	var err error

	// Set your custom logger if needed. Default one is log.Printf
	Logger = log.New(os.Stderr, "rocket-skates", log.LstdFlags|log.Lmicroseconds|log.LUTC)

	if ProvOpts.VersionFlag {
		Logger.Fatalf("Version: %s", version.REBAR_VERSION)
	}
	if ProvOpts.CommandURL != "" {
		rebarClient, err = api.TrustedSession("system", true)
		if err != nil {
			Logger.Fatalf("Error creating trusted Rebar API client: %v", err)
		}
	} else {
		Logger.Printf("Running without a rebar client endpoint - no updates for DR\n")
	}

	ProvisionerURL = fmt.Sprintf("http://%s:%d",
		ProvOpts.OurAddress,
		ProvOpts.StaticPort)
	Logger.Printf("Version: %s\n", version.REBAR_VERSION)
	var consulClient *consul.Client

	if ProvOpts.RegisterConsul {
		consulClient, err = client.Consul(true)
		if err != nil {
			Logger.Fatalf("Error talking to Consul: %v", err)
		}

		// Register service with Consul before continuing
		if err = service.Register(consulClient,
			&consul.AgentServiceRegistration{
				Name: "provisioner-service",
				Tags: []string{"deployment:system"},
				Port: ProvOpts.StaticPort,
				Check: &consul.AgentServiceCheck{
					HTTP:     fmt.Sprintf("http://[::]:%d/", ProvOpts.StaticPort),
					Interval: "10s",
				},
			},
			true); err != nil {
			log.Fatalf("Failed to register provisioner-service with Consul: %v", err)
		}

		if err = service.Register(consulClient,
			&consul.AgentServiceRegistration{
				Name: "provisioner-mgmt-service",
				Tags: []string{"revproxy"}, // We want to be exposed through the revproxy
				Port: apiPort,
				Check: &consul.AgentServiceCheck{
					HTTP:     fmt.Sprintf("http://[::]:%d/", ProvOpts.StaticPort),
					Interval: "10s",
				},
			},
			false); err != nil {
			log.Fatalf("Failed to register provisioner-mgmt-service with Consul: %v", err)
		}
		if err = service.Register(consulClient,
			&consul.AgentServiceRegistration{
				Name: "provisioner-tftp-service",
				Port: ProvOpts.TftpPort,
				Check: &consul.AgentServiceCheck{
					HTTP:     fmt.Sprintf("http://[::]:%d/", ProvOpts.TftpPort),
					Interval: "10s",
				},
			},
			true); err != nil {
			log.Fatalf("Failed to register provisioner-tftp-service with Consul: %v", err)
		}
	}
	var backend store.SimpleStore
	switch ProvOpts.BackEndType {
	case "consul":
		if consulClient == nil {
			consulClient, err = client.Consul(true)
			if err != nil {
				Logger.Fatalf("Error talking to Consul: %v", err)
			}
		}
		backend, err = store.NewSimpleConsulStore(consulClient, ProvOpts.DataRoot)
	case "directory":
		backend, err = store.NewFileBackend(ProvOpts.DataRoot)
	case "memory":
		backend = store.NewSimpleMemoryStore()
		err = nil
	case "bolt", "local":
		backend, err = store.NewSimpleLocalStore(ProvOpts.DataRoot)
	default:
		Logger.Fatalf("Unknown storage backend type %v\n", ProvOpts.BackEndType)
	}
	if err != nil {
		Logger.Fatalf("Error using backing store %s: %v", ProvOpts.BackEndType, err)
	}

	registerBackends(backend)

	go func() {
		if err := ServeTftp(fmt.Sprintf(":%d", ProvOpts.TftpPort)); err != nil {
			Logger.Fatalf("Error starting TFTP server: %v", err)
		}
	}()
	go func() {
		// Static file server must always be last, as all our health checks key off of it.
		if err := ServeStatic(fmt.Sprintf(":%d", ProvOpts.StaticPort), ProvOpts.FileRoot); err != nil {
			Logger.Fatalf("Error starting static file server: %v", err)
		}
	}()

}
