package rspace

import (
	"testing"
)

func TestAbs(t *testing.T) {
	got := GetStatus()
	if got.Message != "OK" {
		t.Errorf("Expected 'OK' but was %v", got.Message)
	}
	if len(got.RSpaceVersion) == 0 {
		t.Errorf("RSpaceVersion must be non-empty")

	}
}
