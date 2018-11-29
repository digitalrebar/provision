package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"text/template"

	dhcp "github.com/krolaw/dhcp4"
)

func DHCPOptionParser(code dhcp.OptionCode) (func(string) ([]byte, error), func([]byte) string) {
	switch code {
	case dhcp.OptionDHCPMessageType:
		return func(s string) ([]byte, error) {
				switch s {
				case "dis":
					return []byte{1}, nil
				case "ofr":
					return []byte{2}, nil
				case "req":
					return []byte{3}, nil
				case "dec":
					return []byte{4}, nil
				case "ack":
					return []byte{5}, nil
				case "nak":
					return []byte{6}, nil
				case "rel":
					return []byte{7}, nil
				case "inf":
					return []byte{8}, nil
				default:
					return nil, fmt.Errorf("Invalid message type %s", s)
				}
			}, func(buf []byte) string {
				switch buf[0] {
				case 1:
					return "dis"
				case 2:
					return "ofr"
				case 3:
					return "req"
				case 4:
					return "dec"
				case 5:
					return "ack"
				case 6:
					return "nak"
				case 7:
					return "rel"
				case 8:
					return "inf"
				default:
					return "unk"
				}
			}
	// Single IP-like address
	case dhcp.OptionSubnetMask,
		dhcp.OptionBroadcastAddress,
		dhcp.OptionSwapServer,
		dhcp.OptionRouterSolicitationAddress,
		dhcp.OptionRequestedIPAddress,
		dhcp.OptionServerIdentifier:
		return func(s string) ([]byte, error) {
				return []byte(net.ParseIP(s).To4()), nil
			}, func(buf []byte) string {
				return net.IP(buf).To4().String()
			}
		// Multiple IP-like address
	case dhcp.OptionRouter,
		dhcp.OptionTimeServer,
		dhcp.OptionNameServer,
		dhcp.OptionDomainNameServer,
		dhcp.OptionLogServer,
		dhcp.OptionCookieServer,
		dhcp.OptionLPRServer,
		dhcp.OptionImpressServer,
		dhcp.OptionResourceLocationServer,
		dhcp.OptionPolicyFilter, // This is special and could validate more (2Ips per)
		dhcp.OptionStaticRoute,  // This is special and could validate more (2IPs per)
		dhcp.OptionNetworkInformationServers,
		dhcp.OptionNetworkTimeProtocolServers,
		dhcp.OptionNetBIOSOverTCPIPNameServer,
		dhcp.OptionNetBIOSOverTCPIPDatagramDistributionServer,
		dhcp.OptionXWindowSystemFontServer,
		dhcp.OptionXWindowSystemDisplayManager,
		dhcp.OptionNetworkInformationServicePlusServers,
		dhcp.OptionMobileIPHomeAgent,
		dhcp.OptionSimpleMailTransportProtocol,
		dhcp.OptionPostOfficeProtocolServer,
		dhcp.OptionNetworkNewsTransportProtocol,
		dhcp.OptionDefaultWorldWideWebServer,
		dhcp.OptionDefaultFingerServer,
		dhcp.OptionDefaultInternetRelayChatServer,
		dhcp.OptionStreetTalkServer,
		dhcp.OptionStreetTalkDirectoryAssistance:
		return func(s string) ([]byte, error) {
				addrs := make([]net.IP, 0)
				alist := strings.Split(s, ",")
				for i := range alist {
					addrs = append(addrs, net.ParseIP(alist[i]).To4())
				}
				return dhcp.JoinIPs(addrs), nil
			}, func(buf []byte) string {
				ips := []string{}
				for len(buf) >= 4 {
					ips = append(ips, net.IP(buf[:4]).To4().String())
					buf = buf[4:]
				}
				return strings.Join(ips, ",")
			}
		// String like value
	case dhcp.OptionHostName,
		dhcp.OptionMeritDumpFile,
		dhcp.OptionDomainName,
		dhcp.OptionRootPath,
		dhcp.OptionExtensionsPath,
		dhcp.OptionNetworkInformationServiceDomain,
		dhcp.OptionNetBIOSOverTCPIPScope,
		dhcp.OptionNetworkInformationServicePlusDomain,
		dhcp.OptionTFTPServerName,
		dhcp.OptionBootFileName,
		dhcp.OptionMessage,
		dhcp.OptionVendorClassIdentifier,
		dhcp.OptionUserClass,
		dhcp.OptionTZPOSIXString,
		dhcp.OptionTZDatabaseString:
		return func(s string) ([]byte, error) {
				return []byte(s), nil
			}, func(buf []byte) string {
				return string(buf)
			}
		// 4 byte integer value
	case dhcp.OptionTimeOffset,
		dhcp.OptionPathMTUAgingTimeout,
		dhcp.OptionARPCacheTimeout,
		dhcp.OptionTCPKeepaliveInterval,
		dhcp.OptionIPAddressLeaseTime,
		dhcp.OptionRenewalTimeValue,
		dhcp.OptionRebindingTimeValue:
		return func(s string) ([]byte, error) {
				answer := make([]byte, 4)
				ival, err := strconv.Atoi(s)
				if err != nil {
					return nil, err
				}
				binary.BigEndian.PutUint32(answer, uint32(ival))
				return answer, nil
			}, func(buf []byte) string {
				return fmt.Sprintf("%d", binary.BigEndian.Uint32(buf))
			}
		// 2 byte integer value
	case dhcp.OptionBootFileSize,
		dhcp.OptionMaximumDatagramReassemblySize,
		dhcp.OptionInterfaceMTU,
		dhcp.OptionMaximumDHCPMessageSize,
		dhcp.OptionClientArchitecture:
		return func(s string) ([]byte, error) {
				answer := make([]byte, 2)
				ival, err := strconv.Atoi(s)
				if err != nil {
					return nil, err
				}
				binary.BigEndian.PutUint16(answer, uint16(ival))
				return answer, nil
			}, func(buf []byte) string {
				return fmt.Sprintf("%d", binary.BigEndian.Uint16(buf))
			}
		// 1 byte integer value
	case dhcp.OptionIPForwardingEnableDisable,
		dhcp.OptionNonLocalSourceRoutingEnableDisable,
		dhcp.OptionDefaultIPTimeToLive,
		dhcp.OptionAllSubnetsAreLocal,
		dhcp.OptionPerformMaskDiscovery,
		dhcp.OptionMaskSupplier,
		dhcp.OptionPerformRouterDiscovery,
		dhcp.OptionTrailerEncapsulation,
		dhcp.OptionEthernetEncapsulation,
		dhcp.OptionTCPDefaultTTL,
		dhcp.OptionTCPKeepaliveGarbage,
		dhcp.OptionNetBIOSOverTCPIPNodeType,
		dhcp.OptionOverload:
		return func(s string) ([]byte, error) {
				answer := make([]byte, 1)
				ival, err := strconv.Atoi(s)
				if err != nil {
					return nil, err
				}
				answer[0] = byte(ival)
				return answer, nil
			}, func(buf []byte) string {
				return fmt.Sprintf("%d", buf[0])
			}
		// Empty
	case dhcp.Pad, dhcp.End:
		return func(s string) ([]byte, error) {
				return []byte{}, nil
			}, func(buf []byte) string {
				return ""
			}
		// Untyped array of bytes
	default:
		return func(s string) ([]byte, error) {
				if strings.HasPrefix(s, "string:") {
					return []byte(strings.TrimPrefix(s, "string:")), nil
				}
				res := []byte{}
				for _, b := range strings.Split(s, ",") {
					ival, err := strconv.Atoi(b)
					if err != nil {
						return nil, err
					}
					res = append(res, byte(ival))
				}
				return res, nil
			}, func(buf []byte) string {
				vals := make([]string, len(buf))
				for i := range buf {
					vals[i] = fmt.Sprintf("%d", buf[i])
				}
				return strings.Join(vals, ",")
			}
	}
}

// DhcpOption is a representation of a specific DHCP option.
// swagger:model
type DhcpOption struct {
	// Code is a DHCP Option Code.
	//
	// required: true
	Code byte
	// Value is a text/template that will be expanded
	// and then converted into the proper format
	// for the option code
	//
	// required: true
	Value string
}

func (o *DhcpOption) String() string {
	return fmt.Sprintf("code:%03d val:%q", o.Code, o.Value)
}

func (o *DhcpOption) Fill(s string) error {
	buf, err := o.ConvertOptionValueToByte(s)
	if err != nil {
		return err
	}
	o.FillFromPacketOpt(buf)
	return nil
}

func (o *DhcpOption) AddToPacket(p *dhcp.Packet) error {
	val, err := o.ConvertOptionValueToByte(o.Value)
	if err != nil {
		return err
	}
	p.AddOption(dhcp.OptionCode(o.Code), val)
	return nil
}

func (o *DhcpOption) FillFromPacketOpt(buf []byte) {
	_, fn := DHCPOptionParser(dhcp.OptionCode(o.Code))
	o.Value = fn(buf)
}

func (o *DhcpOption) ConvertOptionValueToByte(value string) ([]byte, error) {
	fn, _ := DHCPOptionParser(dhcp.OptionCode(o.Code))
	return fn(value)
}

func (o DhcpOption) RenderToDHCP(srcOpts map[int]string) (code byte, val []byte, err error) {
	tmpl, err := template.New("dhcp_option").Funcs(DrpSafeFuncMap()).Parse(o.Value)
	if err != nil {
		return o.Code, nil, err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, srcOpts); err != nil {
		return o.Code, nil, err
	}
	val, err = o.ConvertOptionValueToByte(buf.String())
	return o.Code, val, err
}

func DHCPOptionsInOrder(p dhcp.Packet) []*DhcpOption {
	res := []*DhcpOption{}
	for opts := p.Options(); len(opts) > 2; opts = opts[2+opts[1]:] {
		switch dhcp.OptionCode(opts[0]) {
		case dhcp.End:
			return res
		case dhcp.Pad:
			opts = opts[1:]
		default:
			if len(opts) < int(opts[1])+2 {
				return res
			}
			opt, val := &DhcpOption{Code: opts[0]}, opts[2:2+opts[1]]
			opt.FillFromPacketOpt(val)
			res = append(res, opt)
		}
	}
	return res
}
