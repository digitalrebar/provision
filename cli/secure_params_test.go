package cli

import "testing"

func TestSecureParams(t *testing.T) {
	cliTest(false, false, "params", "create", "-").Stdin(`---
Name: secure
Secure: true
Schema:
  type: string`).run(t)
	cliTest(false, false, "machines", "create", "bob").run(t)
	cliTest(false, false, "profiles", "create", "bob").run(t)
	cliTest(false, false, "plugins", "create", "-").Stdin(`---
Name: bob
Provider: incrementer`).run(t)
	cliTest(false, false, "contents", "upload", "-").Stdin(licenseLayer).run(t)
	for _, tgt := range []string{"machines", "profiles", "plugins"} {
		cliTest(false, false, tgt, "set", "Name:bob", "param", "secure", "to", "Bob").run(t)
		cliTest(false, false, tgt, "get", "Name:bob", "param", "secure").run(t)
		cliTest(false, false, tgt, "get", "Name:bob", "param", "secure", "--decode").run(t)
	}
	cliTest(false, false, "roles", "create", "-").Stdin(`---
Name: secretSetter
Claims:
  - Scope: "params"
    Action: "get"
    Specific: "*"
  - Scope: "machines,profiles,plugins"
    Action: "get,update,updateSecure"
    Specific: "*"`).run(t)
	cliTest(false, false, "roles", "create", "-").Stdin(`---
Name: secretGetter
Claims:
  - Scope: "params"
    Action: "get"
    Specific: "*"
  - Scope: "machines,profiles,plugins"
    Action: "get,getSecure"
    Specific: "*"`).run(t)
	cliTest(false, false, "users", "create", "fred").run(t)
	cliTest(false, false, "users", "create", "fred2").run(t)
	cliTest(false, false, "users", "password", "fred", "fred").run(t)
	cliTest(false, false, "users", "password", "fred2", "fred").run(t)
	cliTest(false, false, "users", "update", "fred", `{"Roles":["secretSetter"]}`).run(t)
	for _, tgt := range []string{"machines", "profiles", "plugins"} {
		cliTest(false, false, "-T", "", "-U", "fred", "-P", "fred", tgt, "set", "Name:bob", "param", "secure", "to", "Fred").run(t)
		cliTest(false, false, "-T", "", "-U", "fred", "-P", "fred", tgt, "get", "Name:bob", "param", "secure").run(t)
		cliTest(false, true, "-T", "", "-U", "fred", "-P", "fred", tgt, "get", "Name:bob", "param", "secure", "--decode").run(t)
	}
	cliTest(false, false, "users", "update", "fred2", `{"Roles":["secretGetter"]}`).run(t)
	for _, tgt := range []string{"machines", "profiles", "plugins"} {
		cliTest(false, true, "-T", "", "-U", "fred2", "-P", "fred", tgt, "set", "Name:bob", "param", "secure", "to", "Freddy").run(t)
		cliTest(false, false, "-T", "", "-U", "fred2", "-P", "fred", tgt, "get", "Name:bob", "param", "secure").run(t)
		cliTest(false, false, "-T", "", "-U", "fred2", "-P", "fred", tgt, "get", "Name:bob", "param", "secure", "--decode").run(t)
	}
	cliTest(false, false, "users", "destroy", "fred2").run(t)
	cliTest(false, false, "users", "destroy", "fred").run(t)
	cliTest(false, false, "roles", "destroy", "secretSetter").run(t)
	cliTest(false, false, "roles", "destroy", "secretGetter").run(t)
	cliTest(false, false, "contents", "destroy", "rackn-license").run(t)
	for _, tgt := range []string{"machines", "profiles", "plugins"} {
		cliTest(false, false, tgt, "destroy", "Name:bob").run(t)
	}
	cliTest(false, false, "params", "destroy", "secure").run(t)
	verifyClean(t)
}

func TestSecureParamUpgrade(t *testing.T) {
	insecureContent := `---
meta:
  Name: insecure
sections:
  params:
    test:
      name: test
      secure: false
      schema:
        type: string
`
	secureContent := `---
meta:
  Name: insecure
sections:
  params:
    test:
      name: test
      secure: true
      schema:
        type: string
`
	cliTest(false, false, "contents", "create", "-").Stdin(insecureContent).run(t)
	cliTest(false, false, "profiles", "create", "bob").run(t)
	cliTest(false, false, "profiles", "set", "bob", "param", "test", "to", `"foo"`).run(t)
	cliTest(false, false, "profiles", "get", "bob", "param", "test").run(t)
	cliTest(false, false, "contents", "update", "insecure", "-").Stdin(secureContent).run(t)
	cliTest(false, false, "profiles", "get", "bob", "param", "test").run(t)
	cliTest(false, false, "profiles", "get", "bob", "param", "test", "--decode").run(t)
	cliTest(false, false, "profiles", "destroy", "bob").run(t)
	cliTest(false, false, "contents", "destroy", "insecure").run(t)
	verifyClean(t)
}
