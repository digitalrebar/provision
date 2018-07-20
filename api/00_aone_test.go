package api

import (
	"testing"
)

func TestFirst(t *testing.T) {
	err := fakeServer()
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}
}
