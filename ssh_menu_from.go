package sshmenu

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// FromReader imports host entries from a reader
func (s *SSHMenu) FromReader(data io.Reader) error {
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Host ") {
			entries := strings.Split(line, " ")
			if err := s.addHosts(entries[1:]); err != nil {
				return err
			}
		}
	}
	return nil
}

// addHosts splits host line and passes on
func (s *SSHMenu) addHosts(entries []string) error {
	if 0 == len(entries) {
		return errors.New("empty list of entries passed")
	}
	for _, entry := range entries {
		for _, host := range strings.Split(entry, " ") {
			if result, err := s.addHost(host, s.Data, host); err != nil {
				return err
			} else {
				s.Data = result
			}
		}
	}
	return nil
}

// addHost splits host and adds to data structure
func (s *SSHMenu) addHost(host string, currentDirectory DirectoryEntry, connect string) (DirectoryEntry, error) {
	for _, separator := range s.separators {
		splitted := strings.Split(host, string(separator))
		if len(splitted) == 1 {
			continue
		}
		var workingDirectory DirectoryEntry
		if value, ok := currentDirectory.Directories[splitted[0]]; !ok {
			workingDirectory = DirectoryEntry{
				Hosts:       make(map[string]string, 0),
				Directories: make(map[string]DirectoryEntry, 0),
			}
		} else {
			workingDirectory = value
		}
		result, err := s.addHost(strings.Join(splitted[1:], string(separator)), workingDirectory, connect)
		currentDirectory.Directories[splitted[0]] = result
		return currentDirectory, err
	}
	if host != "*" {
		currentDirectory.Hosts[host] = connect
	}

	return currentDirectory, nil
}
