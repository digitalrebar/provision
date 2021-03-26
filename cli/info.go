package cli

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	tftp "github.com/digitalrebar/tftp/v3"
	dhcp "github.com/krolaw/dhcp4"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerInfo)
}

func registerInfo(app *cobra.Command) {
	tree := addInfoCommands()
	app.AddCommand(tree)
}

func addInfoCommands() (res *cobra.Command) {
	singularName := "info"
	name := "info"
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	res.AddCommand(&cobra.Command{
		Use:   "check",
		Short: fmt.Sprintf("Fast API check that returns DRP Version"),
		Long:  `A helper function to return API response with version of DRP`,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Printf("{ \"active\": true }\n")
			return nil
		},
	})

	res.AddCommand(&cobra.Command{
		Use:   "get",
		Short: fmt.Sprintf("Get info about DRP"),
		Long:  `A helper function to return information about DRP`,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {

			d, err := Session.Info()
			if err != nil {
				return generateError(err, "Failed to fetch info %v", singularName)
			}
			return prettyPrint(d)
		},
	})

	res.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Get aliveness status of the various DRP ports",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			d, err := Session.Info()
			if err != nil {
				return generateError(err, "Failed to fetch port information")
			}
			type status struct {
				Enabled, Alive bool
				Port           int
			}
			results := map[string]*status{}
			chaddr, _ := net.ParseMAC("de:ad:be:ef:f0:01")
			for _, service := range []string{"API", "Static", "TFTP", "DHCP", "BINL"} {
				alive := false
				switch service {
				case "API":
					results[service] = &status{true, true, d.ApiPort}
				case "Static":
					if d.ProvisionerEnabled {
						host := net.JoinHostPort(Session.Host(), strconv.Itoa(d.FilePort))
						res, err := http.Get("http://" + host + "/")
						if err == nil {
							defer res.Body.Close()
							alive = true
						}
					}
					results[service] = &status{d.ProvisionerEnabled, alive, d.FilePort}
				case "TFTP":
					if d.ProvisionerEnabled && d.TftpEnabled {
						c, err := tftp.NewClient(net.JoinHostPort(Session.Host(), strconv.Itoa(d.TftpPort)))
						if err == nil {
							if src, err := c.Receive("lpxelinux.0", ""); err == nil {
								alive = true
								src.WriteTo(ioutil.Discard)
							}
						}
					}
					results[service] = &status{d.ProvisionerEnabled && d.TftpEnabled, alive, d.TftpPort}
				case "DHCP":
					if d.DhcpEnabled {
						var hosts []string
						hosts, err = net.LookupHost(Session.Host())
						if err != nil {
							return err
						}
						if len(hosts) == 0 {
							return fmt.Errorf("%s does not resolve to anything!", Session.Host())
						}
						xid := make([]byte, 4)
						rand.Read(xid)
						dest := &net.UDPAddr{IP: net.ParseIP(hosts[0]), Port: d.DhcpPort, Zone: ""}
						packet := dhcp.RequestPacket(dhcp.Request, chaddr, net.IPv4(0, 0, 0, 0), xid, false, nil)
						conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
						if err == nil {
							defer conn.Close()
							conn.WriteToUDP(packet, dest)
							conn.SetReadDeadline(time.Now().Add(10 * time.Second))
							reply := make([]byte, 1500)
							sz, err := conn.Read(reply)
							if err == nil {
								packet = dhcp.Packet(reply[:sz])
								options := packet.ParseOptions()
								t := options[dhcp.OptionDHCPMessageType]
								if t != nil && len(t) == 1 && dhcp.MessageType(t[0]) == dhcp.NAK {
									alive = true
								}
							}
						}
					}
					results[service] = &status{d.DhcpEnabled, alive, d.DhcpPort}
				case "BINL":
					if d.BinlEnabled {
						var hosts []string
						hosts, err = net.LookupHost(Session.Host())
						if err != nil {
							return err
						}
						if len(hosts) == 0 {
							return fmt.Errorf("%s does not resolve to anything!", Session.Host())
						}
						xid := make([]byte, 4)
						rand.Read(xid)
						dest := &net.UDPAddr{IP: net.ParseIP(hosts[0]), Port: d.BinlPort, Zone: ""}
						packet := dhcp.RequestPacket(dhcp.Request, chaddr, net.IPv4(0, 0, 0, 0), xid, false, nil)
						conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
						if err == nil {
							defer conn.Close()
							conn.WriteToUDP(packet, dest)
							conn.SetReadDeadline(time.Now().Add(10 * time.Second))
							reply := make([]byte, 1500)
							sz, err := conn.Read(reply)
							if err == nil {
								packet = dhcp.Packet(reply[:sz])
								options := packet.ParseOptions()
								t := options[dhcp.OptionDHCPMessageType]
								if t != nil && len(t) == 1 && dhcp.MessageType(t[0]) == dhcp.NAK {
									alive = true
								}
							}
						}
					}
					results[service] = &status{d.BinlEnabled, alive, d.BinlPort}
				}
			}
			prettyPrint(results)
			return nil
		},
	})
	return res
}
