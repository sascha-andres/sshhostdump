package sshmenu

type (
	// SSHMenu is the root type
	SSHMenu struct {
		Data DirectoryEntry `yaml:"root"`

		separators string `yaml:"-"`
	}

	// DirectoryEntry are the entries at one level
	DirectoryEntry struct {
		Hosts       map[string]string         `yaml:"hosts"`
		Directories map[string]DirectoryEntry `yaml:"directories"`
	}
)
