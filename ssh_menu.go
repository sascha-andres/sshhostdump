package sshmenu

// NewSSHMenu creates a new instance of SSH menu
func NewSSHMenu(separators string) (*SSHMenu, error) {
	entry := DirectoryEntry{
		Directories: make(map[string]DirectoryEntry),
		Hosts:       make([]string, 0),
	}
	return &SSHMenu{
		Data:       entry,
		separators: separators,
	}, nil
}
