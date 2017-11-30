package models

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
	"text/template"

	dhcp "github.com/krolaw/dhcp4"
)

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

func (o *DhcpOption) ConvertOptionValueToByte(value string) ([]byte, error) {
	code := dhcp.OptionCode(o.Code)
	switch code {
	// Single IP-like address
	case dhcp.OptionSubnetMask,
		dhcp.OptionBroadcastAddress,
		dhcp.OptionSwapServer,
		dhcp.OptionRouterSolicitationAddress,
		dhcp.OptionRequestedIPAddress,
		dhcp.OptionServerIdentifier:
		return []byte(net.ParseIP(value).To4()), nil

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

		addrs := make([]net.IP, 0)
		alist := strings.Split(value, ",")
		for i := range alist {
			addrs = append(addrs, net.ParseIP(alist[i]).To4())
		}
		return dhcp.JoinIPs(addrs), nil

	// String like value
	case dhcp.OptionHostName,
		dhcp.OptionMeritDumpFile,
		dhcp.OptionDomainName,
		dhcp.OptionRootPath,
		dhcp.OptionExtensionsPath,
		dhcp.OptionNetworkInformationServiceDomain,
		dhcp.OptionVendorSpecificInformation, // This is wrong, but ...
		dhcp.OptionNetBIOSOverTCPIPScope,
		dhcp.OptionNetworkInformationServicePlusDomain,
		dhcp.OptionTFTPServerName,
		dhcp.OptionBootFileName,
		dhcp.OptionMessage,
		dhcp.OptionVendorClassIdentifier,
		dhcp.OptionClientIdentifier,
		dhcp.OptionUserClass,
		dhcp.OptionTZPOSIXString,
		dhcp.OptionTZDatabaseString:
		return []byte(value), nil

	// 4 byte integer value
	case dhcp.OptionTimeOffset,
		dhcp.OptionPathMTUAgingTimeout,
		dhcp.OptionARPCacheTimeout,
		dhcp.OptionTCPKeepaliveInterval,
		dhcp.OptionIPAddressLeaseTime,
		dhcp.OptionRenewalTimeValue,
		dhcp.OptionRebindingTimeValue:
		answer := make([]byte, 4)
		ival, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		binary.BigEndian.PutUint32(answer, uint32(ival))
		return answer, nil

	// 2 byte integer value
	case dhcp.OptionBootFileSize,
		dhcp.OptionMaximumDatagramReassemblySize,
		dhcp.OptionInterfaceMTU,
		dhcp.OptionMaximumDHCPMessageSize:
		answer := make([]byte, 2)
		ival, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		binary.BigEndian.PutUint16(answer, uint16(ival))
		return answer, nil

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
		dhcp.OptionOverload,
		dhcp.OptionDHCPMessageType:
		answer := make([]byte, 1)
		ival, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		answer[0] = byte(ival)
		return answer, nil

		// Empty
	case dhcp.Pad, dhcp.End:
		return make([]byte, 0), nil
	}

	return nil, errors.New("Invalid Option: " + code.String() + " " + value)
}

func (o *DhcpOption) RenderToDHCP(srcOpts map[int]string) (code byte, val []byte, err error) {
	tmpl, err := template.New("dhcp_option").Parse(o.Value)
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
