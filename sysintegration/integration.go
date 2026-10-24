package sysintegration

import (
	"smirnovtorrent/config"
)

// IntegrationStatus represents the current state of Windows shell and startup integration.
type IntegrationStatus struct {
	TorrentAssociated bool `json:"torrent_associated"`
	MagnetAssociated  bool `json:"magnet_associated"`
	AutostartEnabled  bool `json:"autostart_enabled"`
}

// ApplyIntegration syncs the operating system configuration to match AppConfig settings.
func ApplyIntegration(cfg config.AppConfig) error {
	if err := SetTorrentAssociation(cfg.AssociateTorrents); err != nil {
		return err
	}
	if err := SetMagnetAssociation(cfg.AssociateMagnets); err != nil {
		return err
	}
	if err := SetAutostart(cfg.Autostart); err != nil {
		return err
	}
	return nil
}
