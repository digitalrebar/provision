package models

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/pborman/uuid"
)

type Cert struct {
	Data [][]byte
	Key  ed25519.PrivateKey
	leaf *tls.Certificate
}

func (c *Cert) setTls() error {
	leaf, err := x509.ParseCertificate(c.Data[0])
	if err != nil {
		return err
	}
	c.leaf = &tls.Certificate{
		Certificate: c.Data,
		PrivateKey:  c.Key,
		Leaf:        leaf,
	}
	return nil
}

func (c *Cert) TLS() *tls.Certificate {
	if c.leaf == nil {
		log.Panicf("cert TLS called with nil leaf!")
	}
	return c.leaf
}

type GlobalHaState struct {
	LoadBalanced     bool
	Enabled          bool
	ConsensusEnabled bool
	ConsensusJoin    string
	VirtAddr         string
	ActiveUri        string
	Token            string
	HaID             string
	Valid            bool
	Roots            []Cert
}

func (g *GlobalHaState) FillTls() error {
	for i := range g.Roots {
		if err := (&g.Roots[i]).setTls(); err != nil {
			return err
		}
	}
	return nil
}

type NodeHaState struct {
	ConsensusID         uuid.UUID
	VirtInterface       string
	VirtInterfaceScript string
	ConsensusAddr       string
	ApiUrl              string
	Passive             bool
	Observer            bool
}

type CurrentHAState struct {
	GlobalHaState
	NodeHaState
}

func makeCert(template *x509.Certificate, parentCert *tls.Certificate) (*tls.Certificate, error) {
	var err error
	var priv ed25519.PrivateKey
	var public ed25519.PublicKey
	public, priv, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	var parent *x509.Certificate
	var parentPriv ed25519.PrivateKey
	if parentCert == nil {
		parent = template
		parentPriv = priv
	} else {
		parent = parentCert.Leaf
		parentPriv = parentCert.PrivateKey.(ed25519.PrivateKey)
	}
	var derBytes []byte
	derBytes, err = x509.CreateCertificate(rand.Reader, template, parent, public, parentPriv)
	if err != nil {
		return nil, err
	}
	finalCert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, err
	}
	return &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
		Leaf:        finalCert,
	}, nil
}

func (g *GlobalHaState) RotateRoot(templateMaker func() (*x509.Certificate, error)) (err error) {
	// Generate an initial certificate root.
	var template *x509.Certificate
	template, err = templateMaker()
	if err != nil {
		return
	}
	var finalCert *tls.Certificate
	finalCert, err = makeCert(template, nil)
	if err != nil {
		return
	}
	res := Cert{Data: finalCert.Certificate, Key: finalCert.PrivateKey.(ed25519.PrivateKey), leaf: finalCert}
	if len(g.Roots) == 0 {
		g.Roots = []Cert{res}
	} else if g.Roots[len(g.Roots)-1].leaf.Leaf.NotAfter.After(time.Now()) {
		copy(g.Roots[1:], g.Roots)
		g.Roots[0] = res
	} else {
		g.Roots = append([]Cert{res}, g.Roots...)
	}
	return
}

func (c *CurrentHAState) EndpointCert(templateMaker func() (*x509.Certificate, error)) (*tls.Certificate, error) {
	tmpl, err := templateMaker()
	if err != nil {
		return nil, err
	}
	addr, _, err := net.SplitHostPort(c.ConsensusAddr)
	if err != nil {
		return nil, err
	}
	tmpl.IPAddresses = []net.IP{net.ParseIP(addr)}
	return makeCert(tmpl, c.Roots[0].TLS())
}

func (c *CurrentHAState) OurIp() (string, error) {
	if !c.Enabled {
		return "", errors.New("HA not enabled")
	}
	if c.ConsensusAddr != "" {
		return c.VirtAddr, nil
	}
	if c.LoadBalanced {
		return c.VirtAddr, nil
	}
	ip, _, err := net.ParseCIDR(c.VirtAddr)
	return ip.String(), err
}

func (cOpts *CurrentHAState) Validate() error {
	// Validate HA args.
	if !cOpts.Enabled {
		return nil
	}
	ourAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	consensusAddr := ""
	consensusPort := ""

	if cOpts.ConsensusAddr != "" {
		consensusAddr, consensusPort, err = net.SplitHostPort(cOpts.ConsensusAddr)
		if err != nil {
			return err
		}
		cAddrOk := false
		if net.ParseIP(consensusAddr) == nil {
			return fmt.Errorf("Must specify an IP address for the consensus address")
		}
		for _, ourAddr := range ourAddrs {
			if ourAddr.(*net.IPNet).IP.String() == consensusAddr {
				cAddrOk = true
				break
			}
		}
		if !cAddrOk {
			return fmt.Errorf("Consensus address %s is not present on the system", consensusAddr)
		}
		portNo, _ := strconv.ParseInt(consensusPort, 10, 32)
		if portNo < 0 || portNo > 65536 {
			return fmt.Errorf("Consensus port %d is out of range", portNo)
		}
	}
	if cOpts.LoadBalanced {
		if cOpts.VirtAddr == "" {
			return fmt.Errorf("Error: HA must specify an address that eternal systems will see this system as")
		}
		if net.ParseIP(cOpts.VirtAddr) == nil {
			return fmt.Errorf("Error: Invalid HA address %s", cOpts.VirtAddr)
		}
		lbAddrOk := true
		for _, ourAddr := range ourAddrs {
			if ourAddr.String() == cOpts.VirtAddr {
				lbAddrOk = false
				break
			}
		}
		if !lbAddrOk {
			return fmt.Errorf("Virt address %s is present on the system, not permitted when load balanced", cOpts.VirtAddr)
		}
	} else {
		if cOpts.VirtAddr == "" {
			return fmt.Errorf("Error: HA must specify a VIP in CIDR format that DRP will move around")
		}
		// In HA mode with a VIP, force everything to talk to the VIP address.
		ip, cidr, err := net.ParseCIDR(cOpts.VirtAddr)
		if err != nil {
			return fmt.Errorf("Error: HA IP address %s not valid: %v", cOpts.VirtAddr, err)
		}
		if consensusAddr != "" && consensusAddr == ip.String() {
			return fmt.Errorf("Error: Consensus address %s cannot be the same as the HA virtual IP %s", consensusAddr, cOpts.VirtAddr)
		}
		cidr.IP = ip

		if cOpts.VirtInterface == "" {
			return fmt.Errorf("Error: HA must specify an interface for the VIP that DRP will move around")
		}

		if _, err = net.InterfaceByName(cOpts.VirtInterface); err != nil {
			return fmt.Errorf("Error: HA interface %s not found: %v", cOpts.VirtInterface, err)
		}
	}
	return nil
}

func GetHaState(base string) (*CurrentHAState, error) {
	haStateFile := path.Join(base, "ha-state.json")
	stateFi, err := os.OpenFile(haStateFile, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer stateFi.Close()
	dec := json.NewDecoder(stateFi)
	st := &CurrentHAState{}
	if err = dec.Decode(st); err != nil || !st.Valid {
		st.ConsensusID = uuid.NewRandom()
		return st, SetHaState(base, st)
	}
	if err = st.FillTls(); err != nil {
		return nil, err
	}
	return st, nil
}

func SetHaState(base string, state *CurrentHAState) error {
	haStateFile := path.Join(base, "ha-state.json")
	stateFi, err := ioutil.TempFile(base, ".ha-state-")
	if err != nil {
		return err
	}
	state.Valid = true
	defer os.Remove(stateFi.Name())
	defer stateFi.Close()
	stateFi.Truncate(0)
	enc := json.NewEncoder(stateFi)
	if err = enc.Encode(state); err != nil {
		return err
	}
	if err = stateFi.Sync(); err != nil {
		return err
	}
	return os.Rename(stateFi.Name(), haStateFile)
}
