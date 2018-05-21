package cli

import (
	"fmt"

	"github.com/cloudflare/cfssl/csr"
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

// validator does nothing and will never return an error. It exists because creating a
// csr.Generator requires a validator.
func validator(req *csr.CertificateRequest) error {
	return nil
}

func createCSR(label, CN string, hosts []string) (csrPem, key []byte, err error) {
	csrPem = nil
	key = nil

	// Make CSR for this host
	names := make([]csr.Name, 0, 0)
	name := csr.Name{
		C:  "US",
		ST: "Texas",
		L:  "Austin",
		O:  "RackN",
		OU: "CA Services",
	}
	names = append(names, name)
	req := csr.CertificateRequest{
		KeyRequest: &csr.BasicKeyRequest{"ecdsa", 256},
		CN:         CN,
		Names:      names,
		Hosts:      hosts,
	}

	g := &csr.Generator{Validator: validator}
	csrPem, key, err = g.ProcessRequest(&req)
	return
}
