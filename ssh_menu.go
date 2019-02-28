package sshhostdump

// NewSSHMenu creates a new instance of SSH menu
func NewSSHMenu(separators string) (*SSHMenu, error) {
	entry := DirectoryEntry{
		Directories: make(map[string]DirectoryEntry),
		Hosts:       make(map[string]string),
	}
	return &SSHMenu{
		Data:       entry,
		separators: separators,
	}, nil
}
