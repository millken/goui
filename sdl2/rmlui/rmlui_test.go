package rmlui

import "testing"

func TestCoreVersion(t *testing.T) {

	version := RmlGetVersion()
	if version != "6.0-dev" {
		t.Errorf("Expected version , got %s", version)
	}
}
