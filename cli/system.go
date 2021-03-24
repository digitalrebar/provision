package cli

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/digitalrebar/provision/v4/api"

	"github.com/digitalrebar/provision/v4/models"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(systemInfo)
}

func systemInfo(app *cobra.Command) {
	tree := addSystemCommands()
	app.AddCommand(tree)
}

type cert struct {
	Cert, Key []byte
}

func addSystemCommands() (res *cobra.Command) {
	singularName := "system"
	name := "system"
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	consensus := &cobra.Command{
		Use:   "ha",
		Short: "Access CLI commands to get the state of high availability",
	}

	consensus.AddCommand(&cobra.Command{
		Use:   "id",
		Short: "Get the machine ID of this endpoint in the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "id").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "leader",
		Short: "Get the machine ID of the leader in the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "leader").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "active",
		Short: "Get the machine ID of the current active node in the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "active").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "peers",
		Short: "Get basic info on all members of the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "peers").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "state",
		Short: "Get the HA state of the system.",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res models.CurrentHAState
			if err := Session.Req().UrlFor("system", "consensus", "state").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "dump",
		Short: "Dump the detailed state of all members of the consensus system.",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "fullstate").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "failOverSafe [timeout]",
		Short: "Check to see if at least one non-observer passive node is caught up",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			if len(args) > 1 {
				return fmt.Errorf("Only an optional timeout argument is accepted")
			}
			if dur, err := models.ParseDuration(args[0], "s"); err != nil {
				return err
			} else if dur < 50*time.Millisecond {
				return fmt.Errorf("Duration %s too short, try something larger", args[0])
			} else if dur > 5*time.Second {
				return fmt.Errorf("Duration %s too long, try something shorter", args[0])
			} else {
				return nil
			}
		},
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			req := Session.Req().Post(nil).UrlFor("system", "consensus", "failOverSafe")
			if len(args) == 1 {
				req = req.Params("waitFor", args[0])
			}
			if err := req.Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "introduction [file]",
		Short: "Get an introduction from an existing cluster, save it in [file]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must pass a the name of a file to save the introduction to")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			intro := models.GlobalHaState{}
			err := Session.Req().UrlFor("system", "consensus", "introduction").Do(&intro)
			if err != nil {
				return err
			}
			tgt, err := os.Create(args[0])
			if err != nil {
				return err
			}
			defer tgt.Close()
			enc := json.NewEncoder(tgt)
			enc.SetIndent("", "  ")
			return enc.Encode(intro)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "join [file]",
		Short: "Join a cluster using the introduction saved in [file]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must pass a the name of a file to load the introduction from")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			intro := models.CurrentHAState{}
			if err := Session.Req().UrlFor("system", "consensus", "state").Do(&intro); err != nil {
				return err
			}
			src, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer src.Close()
			dec := json.NewDecoder(src)
			if err = dec.Decode(&intro.GlobalHaState); err != nil {
				return err
			}
			if err = Session.Req().Post(intro).UrlFor("system", "consensus", "enroll").Do(&intro); err != nil {
				return err
			}
			return prettyPrint(intro)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "enroll [endpointUrl] [endpointUser] [endpointPass] extraArgs",
		Short: "Have the endpoint at [endpointUrl] join the cluster.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("Need 3 or more arguments")
			}
			if len(args) > 3 {
				if len(args[3:])%2 != 0 {
					return fmt.Errorf("Extra enroll args must be present in even numbers")
				}
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			intro := models.CurrentHAState{}
			epSess, err := api.UserSessionTokenProxyContext(context.Background(),
				args[0],
				args[1], args[2],
				false, false)
			if err != nil {
				return err
			}
			if err = epSess.Req().UrlFor("system", "consensus", "state").Do(&intro); err != nil {
				return err
			}
			if err = Session.Req().UrlFor("system", "consensus", "introduction").Do(&intro.GlobalHaState); err != nil {
				return err
			}
			if len(args) > 3 {
				for i := 3; i < len(args); i += 2 {
					k, v := args[i], args[i+1]
					switch k {
					case "VirtInterface":
						intro.VirtInterface = v
					case "VirtInterfaceScript":
						intro.VirtInterfaceScript = v
					case "ConsensusAddr":
						intro.ConsensusAddr = v
					case "Observer":
						intro.Observer = strings.ToLower(v) == "true"
					case "LoadBalanced":
						if intro.ConsensusEnabled {
							return fmt.Errorf("LoadBalanced can only be set during self-enrollment")
						}
						intro.LoadBalanced = strings.ToLower(v) == "true"
					case "VirtAddr":
						if intro.ConsensusEnabled {
							return fmt.Errorf("VirtAddr can only be set during self-enrollment")
						}
						intro.VirtAddr = v
					case "HaID":
						if intro.ConsensusEnabled {
							return fmt.Errorf("HaID can only be set during self-enrollment")
						}
						intro.HaID = v
					default:
						return fmt.Errorf("Unknown node HA setting %s", k)
					}
				}
			}
			if err = epSess.Req().Post(intro).UrlFor("system", "consensus", "enroll").Do(&intro); err != nil {
				return err
			}
			return prettyPrint(intro)
		},
	})

	op := &ops{
		name:       name,
		singleName: singularName,
	}
	op.actions()
	res.AddCommand(op.extraCommands...)

	res.AddCommand(&cobra.Command{
		Use:   "upgrade [zip file]",
		Short: "Upgrade DRP with the provided file",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			filePath := args[0]
			fi, err := urlOrFileAsReadCloser(filePath)
			if err != nil {
				return fmt.Errorf("Error opening %s: %v", filePath, err)
			}
			defer fi.Close()
			if info, err := Session.PostBlob(fi, "system", "upgrade"); err != nil {
				return generateError(err, "Failed to post upgrade: %v", filePath)
			} else {
				return prettyPrint(info)
			}
		},
	})
	res.AddCommand(&cobra.Command{
		Use:   "passive",
		Short: "Switch DRP to HA Passive State",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("%v requires 0 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			r := Session.Req().Post(nil).UrlFor(path.Join("/", "system", "passive"))
			var info interface{}
			if err := r.Do(&info); err != nil {
				return generateError(err, "Failed to set passive state")
			} else {
				return prettyPrint(info)
			}
		},
	})
	res.AddCommand(&cobra.Command{
		Use:   "active",
		Short: "Switch DRP to HA Active State",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("%v requires 0 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			r := Session.Req().Post(nil).UrlFor(path.Join("/", "system", "active"))
			var info interface{}
			if err := r.Do(&info); err != nil {
				return generateError(err, "Failed to set active state")
			} else {
				return prettyPrint(info)
			}
		},
	})
	res.AddCommand(&cobra.Command{
		Use:   "signurl [URL]",
		Short: "Generate a RackN Signed URL for download",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			if newurl, err := signRackNUrl(args[0]); err != nil {
				return generateError(err, "Failed to sign url")
			} else {
				fmt.Println(newurl)
				return nil
			}
		},
	})

	certs := &cobra.Command{
		Use:   "certs",
		Short: "Access CLI commands to get and set the TLS cert the API uses",
	}

	certs.AddCommand(&cobra.Command{
		Use:   "get [certFile] [keyFile]",
		Short: "Get the current operating TLS certificate and private key, and save them in PEM format.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			runCert := cert{}
			err := Session.Req().UrlFor("system", "cert").Do(&runCert)
			if err != nil {
				return generateError(err, "Failed to fetch cert")
			}
			if err = ioutil.WriteFile(args[0], runCert.Cert, 0444); err != nil {
				return generateError(err, "Failed to save server certificate")
			}
			if err = ioutil.WriteFile(args[1], runCert.Key, 0400); err != nil {
				return generateError(err, "Failed to save server private key")
			}
			return nil
		},
	})
	certs.AddCommand(&cobra.Command{
		Use:   "set [certFile] [keyFile]",
		Short: "Set the current operating TLS certificate and private key using passed-in PEM encoded files",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			var err error
			runCert := &cert{}
			if runCert.Cert, err = ioutil.ReadFile(args[0]); err != nil {
				return generateError(err, "Failed to read certificate")
			}
			if runCert.Key, err = ioutil.ReadFile(args[1]); err != nil {
				return generateError(err, "Failed to read private key")
			}
			if _, err = tls.X509KeyPair(runCert.Cert, runCert.Key); err != nil {
				return generateError(err, "Invalid TLS certificate/key combination")
			}
			if err = Session.Req().Post(runCert).UrlFor("system", "cert").Do(&runCert); err != nil {
				return generateError(err, "Failed to update running TLS certificate")
			}
			return nil
		},
	})
	res.AddCommand(certs, consensus)
	return res
}
