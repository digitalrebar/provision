package api

import (
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestJobs(t *testing.T) {
	machine1 := mustDecode(&models.Machine, `
Address: 192.168.100.110
BootEnv: local
Name: john
Uuid: 3e7031fe-3062-45f1-835c-92541bc9cbd3
Validated: true
`)
}
