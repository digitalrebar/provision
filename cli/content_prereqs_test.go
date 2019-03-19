package cli

import "testing"

func TestContentVersions(t *testing.T) {
	cliTest(false, false, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.3
  Prerequisites: "BasicStore: >=3.12.0"`).run(t)
	cliTest(false, true, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.3
  Prerequisites: "BasicStore: >=9.99.0"`).run(t)
	cliTest(false, true, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.foo
  Prerequisites: 'BasicStore: >=9.99.0'`).run(t)
	cliTest(false, false, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.4
  Prerequisites: 'BasicStore: >=3.12.0'`).run(t)
	cliTest(false, false, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.4
  Prerequisites: 'BasicStore: >=3.12.x'`).run(t)
	cliTest(false, true, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.4
  Prerequisites: 'BasicStore: >=3.12.foo'`).run(t)
	cliTest(false, false, "contents", "upload", "--format", "yaml", "-").Stdin(`
meta:
  Name: basic
  Version: 1.2.4
  Prerequisites: 'foo, BasicStore: >=3.12.x'`).run(t)
	cliTest(false, false, "contents", "show", "basic", "--format", "yaml").run(t)
	cliTest(false, false, "contents", "destroy", "basic").run(t)
	verifyClean(t)
}
