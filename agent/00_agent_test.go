package agent

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var err error
	actuallyPowerThings = false

	tmpDir, err = ioutil.TempDir("", "cli-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	err = fakeServer()
	if err != nil {
		log.Fatalf("Failed with error: %v", err)
	}

	ret := m.Run()

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Printf("Removing temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	os.Exit(ret)
}
