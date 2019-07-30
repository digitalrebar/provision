package cli

import (
	"path"
	"runtime"
	"testing"
)

func TestLoadIncrementer(t *testing.T) {
	cliTest(false, false,
		"plugin_providers", "upload", "incrementer", "from", path.Join("../bin", runtime.GOOS, runtime.GOARCH, "incrementer")).run(t)
}
