package sshmenu

import (
	"strings"
	"testing"
)

const separators = "./"

func TestFromReader(t *testing.T) {
	s, err := NewSSHMenu(separators)
	if err != nil {
		t.Logf("constructor threw error: %s", err)
		t.Fail()
	}
	err = s.FromReader(strings.NewReader(`Host test
Host dir.test
Host two.entries for.one.host
Host two.server`))
	if err != nil {
		t.Logf("error: %s", err)
		t.Fail()
	}

	if s.Data.Hosts[0] != "test" {
		t.Log("expected test to be a root level host")
		t.Fail()
	}
}

func TestEmptyHostEntries(t *testing.T) {
	s, err := NewSSHMenu(separators)
	if err != nil {
		t.Logf("constructor threw error: %s", err)
		t.Fail()
	}
	err = s.addHosts(make([]string, 0))
	if err == nil {
		t.Log("expected error, received none")
		t.Fail()
	}
}
