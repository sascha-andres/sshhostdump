package sshhostdump

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

	if data, ok := s.Data.Hosts["test"]; !ok {
		t.Log("no test host found")
		t.Fail()
	} else {
		if data != "test" {
			t.Logf("expected [test] got [%s]", data)
			t.Fail()
		}
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
