package cli

import "testing"

func TestMeta(t *testing.T) {
	cliTest(false, false, "profiles", "create", "bob").run(t)
	cliTest(false, false, "profiles", "meta", "bob").run(t)
	cliTest(false, true, "profiles", "meta", "remove", "bob", "key", "foo").run(t)
	cliTest(false, false, "profiles", "meta", "add", "bob", "key", "foo", "val", "bar").run(t)
	cliTest(false, true, "profiles", "meta", "add", "bob", "key", "foo", "val", "bar").run(t)
	cliTest(false, false, "profiles", "meta", "set", "bob", "key", "foo", "val", "baz").run(t)
	cliTest(false, false, "profiles", "meta", "remove", "bob", "key", "foo").run(t)
	cliTest(false, false, "profiles", "meta", "set", "bob", "key", "foo", "val", "baz").run(t)
	cliTest(false, false, "profiles", "meta", "bob").run(t)
	cliTest(false, false, "profiles", "destroy", "bob").run(t)
	cliTest(false, true, "profiles", "meta", "bob").run(t)
	verifyClean(t)
}
