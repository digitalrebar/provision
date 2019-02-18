package cli

import (
	"os"
	"testing"
)

func TestContentCli(t *testing.T) {
	var (
		contentCreateBadJSONString         = "{asdgasdg"
		contentCreateBadJSON2String        = "[asdgasdg]"
		contentCreateInputString    string = `{
  "meta": {
    "Name": "john"
  }
}
`

		contentUpdateBadJSONString         = "asdgasdg"
		contentUpdateBadInputString string = `{
  "meta": {
    "Name": "john2"
  }
}
`
		contentUpdateInputString string = `{
  "meta": {
    "Description": "Fred Rules",
    "Name": "john"
  }
}
`
		contentWithProfileString string = `{
"meta": {
  "Name": "withProfile"
},
"sections": {
  "params": {
    "secureFoo": {
      "Name": "secureFoo",
      "Secure": true,
      "Schema": { "type":"string"},
    }
  },
  "profiles": {
    "englobal": {
      "Name": "englobal",
      "Params": {
        "foo": "bar",
        "secureFoo": "secureBar"
}
}
}
}
}
`
	)

	cliTest(true, false, "contents").run(t)
	cliTest(false, false, "contents", "list").run(t)

	cliTest(false, true, "contents", "create").run(t)
	cliTest(false, true, "contents", "create", "john", "john2").run(t)
	cliTest(false, true, "contents", "create", "https://gi333thub.com/digitalrebar/provision-content/releases/download/v1.3.0/drp-community-content.yamljj").run(t)
	cliTest(false, false, "contents", "create", "https://github.com/digitalrebar/provision-content/releases/download/v1.3.0/drp-community-content.yaml").run(t)
	cliTest(false, false, "contents", "destroy", "drp-community-content").run(t)
	cliTest(false, false, "contents", "create", "test-data/content.yaml").run(t)
	cliTest(false, false, "contents", "destroy", "drp-community-contrib").run(t)
	cliTest(false, true, "contents", "create", contentCreateBadJSONString).run(t)
	cliTest(false, true, "contents", "create", contentCreateBadJSON2String).run(t)
	cliTest(false, false, "contents", "create", contentCreateInputString).run(t)
	cliTest(false, true, "contents", "create", contentCreateInputString).run(t)
	cliTest(false, false, "contents", "destroy", "john").run(t)
	cliTest(false, false, "contents", "upload", contentCreateInputString).run(t)
	cliTest(false, false, "contents", "upload", contentCreateInputString).run(t)
	cliTest(false, false, "contents", "list").run(t)
	cliTest(false, true, "contents", "list", "--limit=-1", "--offset=-1").run(t)
	cliTest(false, true, "contents", "list", "Cow").run(t)
	cliTest(false, true, "contents", "list", "Cow=john").run(t)

	cliTest(true, true, "contents", "show").run(t)
	cliTest(true, true, "contents", "show", "john", "john2").run(t)
	cliTest(false, true, "contents", "show", "john2").run(t)
	cliTest(false, false, "contents", "show", "john").run(t)

	cliTest(false, true, "contents", "exists").run(t)
	cliTest(false, true, "contents", "exists", "john", "john2").run(t)
	cliTest(false, false, "contents", "exists", "john").run(t)
	cliTest(false, true, "contents", "exists", "john2").run(t)
	cliTest(true, true, "contents", "exists", "john", "john2").run(t)

	cliTest(false, true, "contents", "update").run(t)
	cliTest(false, true, "contents", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "contents", "update", "john", contentUpdateBadJSONString).run(t)
	cliTest(false, true, "contents", "update", "john", contentUpdateBadInputString).run(t)
	cliTest(false, false, "contents", "update", "john", contentUpdateInputString).run(t)
	cliTest(false, true, "contents", "update", "john2", contentUpdateInputString).run(t)
	cliTest(false, false, "contents", "show", "john").run(t)

	cliTest(false, true, "contents", "destroy").run(t)
	cliTest(false, true, "contents", "destroy", "john", "june").run(t)
	cliTest(false, false, "contents", "destroy", "john").run(t)
	cliTest(false, true, "contents", "destroy", "john").run(t)
	cliTest(false, false, "contents", "list").run(t)

	cliTest(false, false, "contents", "create", "-").Stdin(contentCreateInputString + "\n").run(t)
	cliTest(false, false, "contents", "list").run(t)
	cliTest(false, false, "contents", "update", "john", "-").Stdin(contentUpdateInputString + "\n").run(t)
	cliTest(false, false, "contents", "show", "john").run(t)

	cliTest(false, false, "contents", "destroy", "john").run(t)
	cliTest(false, false, "contents", "create", "-").Stdin(contentWithProfileString + "\n").run(t)
	cliTest(false, true, "profiles", "set", "englobal", "param", "foo", "to", "baz").run(t)
	cliTest(false, false, "profiles", "get", "englobal", "param", "foo").run(t)
	cliTest(false, false, "profiles", "get", "englobal", "param", "secureFoo").run(t)
	cliTest(false, false, "profiles", "get", "englobal", "param", "secureFoo", "--decode").run(t)
	cliTest(false, false, "profiles", "show", "englobal", "--decode").run(t)
	cliTest(false, false, "profiles", "show", "englobal").run(t)
	cliTest(false, true, "contents", "show", "withProfile").run(t)
	cliTest(false, false, "contents", "show", "withProfile", "--key", "/tmp/drpkey.json").run(t)
	cliTest(false, false, "contents", "destroy", "withProfile").run(t)
	cliTest(false, false, "contents", "list").run(t)

	cliTest(true, true, "contents", "bundlize").run(t)
	cliTest(false, true, "contents", "bundlize", "greg.yaml").run(t)
	cliTest(true, true, "contents", "bundlize", "greg.yaml", "greg").run(t)
	cliTest(false, true, "contents", "bundlize", "greg.yaml", "Name=greg").run(t)
	cliTest(false, true, "contents", "bundlize", "greg.yaml", "Name=greg", "profiles:test1").run(t)

	cliTest(false, false, "profiles", "create", "test1").run(t)
	cliTest(false, false, "contents", "bundlize", "greg.yaml", "Name=greg", "profiles:test1").run(t)
	cliTest(false, false, "profiles", "show", "test1").run(t)

	// Test if greg.yaml exists
	if _, err := os.Stat("greg.yaml"); os.IsNotExist(err) {
		t.Errorf("Failed to find greg.yaml\n")
	}
	os.Remove("greg.yaml")

	cliTest(false, false, "contents", "bundlize", "greg.yaml", "Name=greg", "profiles:test1", "--delete", "--reload").run(t)
	cliTest(false, false, "profiles", "show", "test1").run(t)
	cliTest(false, false, "contents", "destroy", "greg").run(t)
	cliTest(false, true, "profiles", "show", "test1").run(t)

	// Test if greg.yaml exists
	if _, err := os.Stat("greg.yaml"); os.IsNotExist(err) {
		t.Errorf("Failed to find greg.yaml\n")
	}
	os.Remove("greg.yaml")

	cliTest(false, false, "profiles", "create", "test1").run(t)
	cliTest(false, false, "params", "create", "-").Stdin(`---
Name: secureTest
Secure: true
Schema:
  type: string`).run(t)
	cliTest(false, false, "profiles", "set", "test1", "param", "secureTest", "to", "Fred").run(t)
	cliTest(false, false, "profiles", "show", "test1").run(t)
	cliTest(false, false, "profiles", "show", "test1", "--decode").run(t)
	cliTest(false, true, "contents", "bundlize", "greg.yaml", "Name=greg", "profiles:test1", "params:secureTest", "--delete", "--reload").run(t)
	cliTest(false, false, "contents", "bundlize", "greg.yaml", "Name=greg", "profiles:test1", "params:secureTest", "--delete", "--reload", "--key", "/tmp/drpkey.json").run(t)
	cliTest(false, false, "profiles", "show", "test1").run(t)
	cliTest(false, false, "profiles", "show", "test1", "--decode").run(t)
	cliTest(false, false, "params", "show", "secureTest").run(t)
	cliTest(false, false, "contents", "destroy", "greg").run(t)
	cliTest(false, true, "profiles", "show", "test1").run(t)
	os.Remove("greg.yaml")

	cliTest(false, false, "profiles", "create", "test1").run(t)
	cliTest(false, false, "contents", "bundlize", "greg.yaml", "Name=greg", "profiles:test1", "--delete").run(t)
	cliTest(false, true, "profiles", "show", "test1").run(t)
	// Test if greg.yaml exists
	if _, err := os.Stat("greg.yaml"); os.IsNotExist(err) {
		t.Errorf("Failed to find greg.yaml\n")
	}

	cliTest(true, true, "contents", "convert").run(t)
	cliTest(true, true, "contents", "convert", "gg", "ff").run(t)
	cliTest(false, true, "contents", "convert", "gg.yaml").run(t)
	cliTest(false, false, "contents", "convert", "greg.yaml").run(t)
	cliTest(false, false, "profiles", "show", "test1").run(t)
	cliTest(false, false, "profiles", "destroy", "test1").run(t)

	verifyClean(t)

	os.Remove("greg.yaml")
}

func TestContentRequiredFeatures(t *testing.T) {
	okContent := `---
meta:
  Name: insecure
  RequiredFeatures: 'secure-params'
sections:
  params:
    test:
      Name: test
      Secure: true
      Schema:
        type: string
`

	badContent := `---
meta:
  Name: insecure
  RequiredFeatures: 'missing-feature'
sections:
  params:
    test:
      Name: test
      Secure: true
      Schema:
        type: string
`
	cliTest(false, true, "contents", "create", "-").Stdin(badContent).run(t)
	cliTest(false, false, "contents", "create", "-").Stdin(okContent).run(t)
	cliTest(false, true, "contents", "update", "insecure", "-").Stdin(badContent).run(t)
	cliTest(false, false, "contents", "destroy", "insecure").run(t)
	verifyClean(t)
}
