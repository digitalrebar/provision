package cli

import "testing"

func TestCertsCli(t *testing.T) {
	cliTest(true, false, "certs").run(t)
	cliTest(true, true, "certs", "csr").run(t)
	cliTest(true, true, "certs", "csr", "root").run(t)
	cliTest(false, false, "certs", "csr", "root", "cn1").run(t)
	cliTest(false, false, "certs", "csr", "root", "cn1", "an1").run(t)
}
