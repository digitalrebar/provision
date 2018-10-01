package cli

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/mail"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerCerts)
}

func registerCerts(app *cobra.Command) {
	app.AddCommand(certsCommands("certs"))
}

func certsCommands(bt string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   bt,
		Short: fmt.Sprintf("Access CLI commands relating to %v", bt),
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "csr [root] [cn] [hosts...]",
		Short: "Create a CSR and private key",
		Long:  "You must pass the root CA name, the common name, and may add additional hosts",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) >= 2 {
				return nil
			}
			return fmt.Errorf("%v: Expected more than 1 arguments", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			root := args[0]
			cn := args[1]
			hosts := []string{}
			if len(args) > 2 {
				hosts = args[2:]
			}

			csr, key, err := createCSR(root, cn, hosts)
			if err != nil {
				return generateError(err, "building csr")
			}

			answer := struct {
				CSR string
				Key string
			}{
				CSR: string(csr),
				Key: string(key),
			}
			return prettyPrint(answer)
		},
	})
	return cmd
}

func createCSR(label, CN string, hosts []string) (csrPem, key []byte, err error) {
	names := []csrName{}
	cname := csrName{
		C:  "US",
		ST: "Texas",
		L:  "Austin",
		O:  "RackN",
		OU: "CA Services",
	}
	names = append(names, cname)
	req := certReq{
		KeyRequest: &basicKeyReq{"ecdsa", 256},
		CN:         CN,
		Names:      names,
		Hosts:      hosts,
	}

	csrPem, key, err = parseCSRReq(&req)
	return
}

// Everything below here was taken from
// github.com/cloudflare/cfssl/csr/csr.go
// and lightly modified to unexport things
// and to not rely on some helper libraries.

/*
Copyright (c) 2014 CloudFlare Inc.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.

Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

const (
	curveP256 = 256
	curveP384 = 384
	curveP521 = 521
)

// A csrName contains the SubjectInfo fields.
type csrName struct {
	C            string // Country
	ST           string // State
	L            string // Locality
	O            string // OrganisationName
	OU           string // OrganisationalUnitName
	SerialNumber string
}

// A keyReq is a generic request for a new key.
type keyReq interface {
	Algo() string
	Size() int
	Generate() (crypto.PrivateKey, error)
	SigAlgo() x509.SignatureAlgorithm
}

// A basicKeyReq contains the algorithm and key size for a new private key.
type basicKeyReq struct {
	A string `json:"algo" yaml:"algo"`
	S int    `json:"size" yaml:"size"`
}

// NewBasicKeyRequest returns a default BasicKeyRequest.
func NewBasicKeyRequest() *basicKeyReq {
	return &basicKeyReq{"ecdsa", curveP256}
}

// Algo returns the requested key algorithm represented as a string.
func (kr *basicKeyReq) Algo() string {
	return kr.A
}

// Size returns the requested key size.
func (kr *basicKeyReq) Size() int {
	return kr.S
}

// Generate generates a key as specified in the request. Currently,
// only ECDSA and RSA are supported.
func (kr *basicKeyReq) Generate() (crypto.PrivateKey, error) {
	switch kr.Algo() {
	case "rsa":
		if kr.Size() < 2048 {
			return nil, errors.New("RSA key is too weak")
		}
		if kr.Size() > 8192 {
			return nil, errors.New("RSA key size too large")
		}
		return rsa.GenerateKey(rand.Reader, kr.Size())
	case "ecdsa":
		var curve elliptic.Curve
		switch kr.Size() {
		case curveP256:
			curve = elliptic.P256()
		case curveP384:
			curve = elliptic.P384()
		case curveP521:
			curve = elliptic.P521()
		default:
			return nil, errors.New("invalid curve")
		}
		return ecdsa.GenerateKey(curve, rand.Reader)
	default:
		return nil, errors.New("invalid algorithm")
	}
}

// SigAlgo returns an appropriate X.509 signature algorithm given the
// key request's type and size.
func (kr *basicKeyReq) SigAlgo() x509.SignatureAlgorithm {
	switch kr.Algo() {
	case "rsa":
		switch {
		case kr.Size() >= 4096:
			return x509.SHA512WithRSA
		case kr.Size() >= 3072:
			return x509.SHA384WithRSA
		case kr.Size() >= 2048:
			return x509.SHA256WithRSA
		default:
			return x509.SHA1WithRSA
		}
	case "ecdsa":
		switch kr.Size() {
		case curveP521:
			return x509.ECDSAWithSHA512
		case curveP384:
			return x509.ECDSAWithSHA384
		case curveP256:
			return x509.ECDSAWithSHA256
		default:
			return x509.ECDSAWithSHA1
		}
	default:
		return x509.UnknownSignatureAlgorithm
	}
}

// caConfig is a section used in the requests initialising a new CA.
type caConfig struct {
	PathLength  int    `json:"pathlen" yaml:"pathlen"`
	PathLenZero bool   `json:"pathlenzero" yaml:"pathlenzero"`
	Expiry      string `json:"expiry" yaml:"expiry"`
	Backdate    string `json:"backdate" yaml:"backdate"`
}

// A certReq encapsulates the API interface to the
// certificate request functionality.
type certReq struct {
	CN           string
	Names        []csrName `json:"names" yaml:"names"`
	Hosts        []string  `json:"hosts" yaml:"hosts"`
	KeyRequest   keyReq    `json:"key,omitempty" yaml:"key,omitempty"`
	CA           *caConfig `json:"ca,omitempty" yaml:"ca,omitempty"`
	SerialNumber string    `json:"serialnumber,omitempty" yaml:"serialnumber,omitempty"`
}

// appendIf appends to a if s is not an empty string.
func appendIf(s string, a *[]string) {
	if s != "" {
		*a = append(*a, s)
	}
}

// Name returns the PKIX name for the request.
func (cr *certReq) Name() pkix.Name {
	var name pkix.Name
	name.CommonName = cr.CN

	for _, n := range cr.Names {
		appendIf(n.C, &name.Country)
		appendIf(n.ST, &name.Province)
		appendIf(n.L, &name.Locality)
		appendIf(n.O, &name.Organization)
		appendIf(n.OU, &name.OrganizationalUnit)
	}
	name.SerialNumber = cr.SerialNumber
	return name
}

// csrRestraints CSR information RFC 5280, 4.2.1.9
type csrRestraints struct {
	IsCA       bool `asn1:"optional"`
	MaxPathLen int  `asn1:"optional,default:-1"`
}

// appendCAInfoToCSR appends CAConfig BasicConstraint extension to a CSR
func appendCAInfoToCSR(reqConf *caConfig, csr *x509.CertificateRequest) error {
	pathlen := reqConf.PathLength
	if pathlen == 0 && !reqConf.PathLenZero {
		pathlen = -1
	}
	val, err := asn1.Marshal(csrRestraints{true, pathlen})

	if err != nil {
		return err
	}

	csr.ExtraExtensions = []pkix.Extension{
		{
			Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
			Value:    val,
			Critical: true,
		},
	}

	return nil
}

// genCSR creates a new CSR from a CertificateRequest structure and
// an existing key. The KeyRequest field is ignored.
func genCSR(priv crypto.Signer, req *certReq) (csr []byte, err error) {
	var tpl = x509.CertificateRequest{
		Subject:            req.Name(),
		SignatureAlgorithm: req.KeyRequest.SigAlgo(),
	}

	for i := range req.Hosts {
		if ip := net.ParseIP(req.Hosts[i]); ip != nil {
			tpl.IPAddresses = append(tpl.IPAddresses, ip)
		} else if email, eerr := mail.ParseAddress(req.Hosts[i]); eerr == nil && email != nil {
			tpl.EmailAddresses = append(tpl.EmailAddresses, email.Address)
		} else {
			tpl.DNSNames = append(tpl.DNSNames, req.Hosts[i])
		}
	}

	if req.CA != nil {
		err = appendCAInfoToCSR(req.CA, &tpl)
		if err != nil {
			return
		}
	}

	csr, err = x509.CreateCertificateRequest(rand.Reader, &tpl, priv)
	if err != nil {
		return
	}
	block := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	}

	csr = pem.EncodeToMemory(&block)
	return
}

// parseCSRReq takes a certificate request and generates a key and
// CSR from it. It does no validation -- caveat emptor. It will,
// however, fail if the key request is not valid (i.e., an unsupported
// curve or RSA key size). The lack of validation was specifically
// chosen to allow the end user to define a policy and validate the
// request appropriately before calling this function.
func parseCSRReq(req *certReq) (csr, key []byte, err error) {
	if req.KeyRequest == nil {
		req.KeyRequest = NewBasicKeyRequest()
	}

	priv, err := req.KeyRequest.Generate()
	if err != nil {
		return
	}

	switch priv := priv.(type) {
	case *rsa.PrivateKey:
		key = x509.MarshalPKCS1PrivateKey(priv)
		block := pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: key,
		}
		key = pem.EncodeToMemory(&block)
	case *ecdsa.PrivateKey:
		key, err = x509.MarshalECPrivateKey(priv)
		if err != nil {
			return
		}
		block := pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: key,
		}
		key = pem.EncodeToMemory(&block)
	default:
		panic("Generate should have failed to produce a valid key.")
	}

	csr, err = genCSR(priv.(crypto.Signer), req)
	if err != nil {
	}
	return
}
