package engine

import "time"

// TorrentSummary is the lightweight view for the main torrent list table
type TorrentSummary struct {
	InfoHash       string    `json:"info_hash"`
	Name           string    `json:"name"`
	Size           int64     `json:"size"`
	Downloaded     int64     `json:"downloaded"`
	Uploaded       int64     `json:"uploaded"`
	Progress       float64   `json:"progress"`       // 0.0 to 100.0
	Status         string    `json:"status"`         // downloading, seeding, paused, metadata, checking, complete
	DownloadSpeed  int64     `json:"download_speed"` // bytes/sec
	UploadSpeed    int64     `json:"upload_speed"`   // bytes/sec
	ETA            int64     `json:"eta"`            // seconds remaining (-1 = unknown / inf)
	PeersConnected int       `json:"peers_connected"`
	PeersTotal     int       `json:"peers_total"`
	SeedsConnected int       `json:"seeds_connected"`
	Ratio          float64   `json:"ratio"`
	AddedAt        time.Time `json:"added_at"`
	SavePath       string    `json:"save_path"`
	Sequential     bool      `json:"sequential"`
	NumFiles       int       `json:"num_files"`
}

// FileEntry represents a single file inside the torrent
type FileEntry struct {
	Index          int     `json:"index"`
	Path           string  `json:"path"`
	Length         int64   `json:"length"`
	BytesCompleted int64   `json:"bytes_completed"`
	Progress       float64 `json:"progress"`
	Priority       int     `json:"priority"` // 0 = Do Not Download, 1 = Normal, 2 = High
}

// PeerEntry represents a connected peer
type PeerEntry struct {
	Address      string  `json:"address"`
	Client       string  `json:"client"`
	DownloadRate float64 `json:"download_rate"`
	UploadRate   float64 `json:"upload_rate"`
	Progress     float64 `json:"progress"`
}

// TrackerEntry represents a tracker URL and status
type TrackerEntry struct {
	URL     string `json:"url"`
	Status  string `json:"status"`
	Tier    int    `json:"tier"`
	Message string `json:"message"`
}

// TorrentDetail contains all detailed data for the inspector tabs
type TorrentDetail struct {
	Summary      TorrentSummary `json:"summary"`
	PieceLength  int64          `json:"piece_length"`
	PieceCount   int            `json:"piece_count"`
	PieceMap     []byte         `json:"piece_map"` // 0 = not downloaded, 1 = partial/downloading, 2 = complete
	Comment      string         `json:"comment"`
	CreatedBy    string         `json:"created_by"`
	CreationDate int64          `json:"creation_date"`
	Files        []FileEntry    `json:"files"`
	Peers        []PeerEntry    `json:"peers"`
	Trackers     []TrackerEntry `json:"trackers"`
}

// GlobalStats holds aggregated metrics for the status bar
type GlobalStats struct {
	TotalDownloadSpeed int64 `json:"total_download_speed"`
	TotalUploadSpeed   int64 `json:"total_upload_speed"`
	TotalDownloaded    int64 `json:"total_downloaded"`
	TotalUploaded      int64 `json:"total_uploaded"`
	ActiveTorrents     int   `json:"active_torrents"`
	TotalTorrents      int   `json:"total_torrents"`
	DHTNodes           int   `json:"dht_nodes"`
}
