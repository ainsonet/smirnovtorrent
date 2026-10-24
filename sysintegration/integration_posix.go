//go:build !windows

package sysintegration

// GetStatus returns dummy integration status on non-Windows platforms.
func GetStatus() (IntegrationStatus, error) {
	return IntegrationStatus{
		TorrentAssociated: false,
		MagnetAssociated:  false,
		AutostartEnabled:  false,
	}, nil
}

func IsTorrentAssociated() bool {
	return false
}

func SetTorrentAssociation(enable bool) error {
	return nil
}

func IsMagnetAssociated() bool {
	return false
}

func SetMagnetAssociation(enable bool) error {
	return nil
}

func IsAutostartEnabled() bool {
	return false
}

func SetAutostart(enable bool) error {
	return nil
}

func SetTestRegistryPaths(classes, run string) func() {
	return func() {}
}
