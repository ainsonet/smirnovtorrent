package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"smirnovtorrent/internal/config"
	"smirnovtorrent/internal/dht"
	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/logger"
	"smirnovtorrent/internal/magnet"
	"smirnovtorrent/internal/parser"
)

//go:embed webui.html
var webUIAssets embed.FS

const version = "1.0.0"

var appLog *logger.Logger

func main() {
	// Инициализация логгера
	logConfig := logger.DefaultConfig()
	logConfig.Prefix = "[CLI]"
	appLog = logger.New(logConfig)
	defer appLog.Close()

	// Загрузка конфигурации
	cfgPath, err := config.GetConfigPath()
	if err != nil {
		appLog.Warn("Failed to get config path: %v", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		appLog.Warn("Failed to load config: %v", err)
		cfg = config.DefaultConfig()
	}

	appLog.Debug("Config loaded from: %s", cfgPath)
	appLog.Debug("DHT: %v, PEX: %v, Encryption: %v", cfg.EnableDHT, cfg.EnablePEX, cfg.EnableEncryption)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("SmirnovTorrent v%s\n", version)
	case "download":
		downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
		outputDir := downloadCmd.String("o", "", "Output directory")
		downloadLimit := downloadCmd.Int64("download-limit", cfg.DownloadRateLimit, "Download speed limit in bytes/sec (0 = unlimited)")
		uploadLimit := downloadCmd.Int64("upload-limit", cfg.UploadRateLimit, "Upload speed limit in bytes/sec (0 = unlimited)")
		enableDHT := downloadCmd.Bool("dht", cfg.EnableDHT, "Enable DHT peer discovery")
		enablePEX := downloadCmd.Bool("pex", cfg.EnablePEX, "Enable Peer Exchange")
		enableEncryption := downloadCmd.Bool("encrypt", cfg.EnableEncryption, "Enable MSE encryption")
		
		downloadCmd.Usage = func() {
			fmt.Println("Usage: smirnovtorrent download [options] <file.torrent|magnet>")
			downloadCmd.PrintDefaults()
		}
		
		downloadCmd.Parse(os.Args[2:])
		
		if downloadCmd.NArg() < 1 {
			fmt.Println("Error: torrent file or magnet link required")
			downloadCmd.Usage()
			os.Exit(1)
		}
		
		torrentSource := downloadCmd.Arg(0)
		
		// Применяем флаги к конфигурации
		if *downloadLimit > 0 {
			cfg.DownloadRateLimit = *downloadLimit
		}
		if *uploadLimit > 0 {
			cfg.UploadRateLimit = *uploadLimit
		}
		cfg.EnableDHT = *enableDHT
		cfg.EnablePEX = *enablePEX
		cfg.EnableEncryption = *enableEncryption
		
		download(torrentSource, cfg, *outputDir)
		
	case "webui":
		port := cfg.WebUIPort
		if len(os.Args) >= 3 {
			fmt.Sscanf(os.Args[2], "%d", &port)
		}
		startWebUI(port, cfg)
	case "info":
		if len(os.Args) < 3 {
			fmt.Println("Error: torrent file required")
			fmt.Println("Usage: smirnovtorrent info <file.torrent>")
			os.Exit(1)
		}
		showInfo(os.Args[2])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("SmirnovTorrent - Lightweight BitTorrent Client")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  smirnovtorrent <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  download <file.torrent|magnet>  Download a torrent")
	fmt.Println("  webui [port]                    Start web interface (default: 8080)")
	fmt.Println("  info <file.torrent>             Show torrent information")
	fmt.Println("  version                         Show version")
	fmt.Println("  help                            Show this help message")
	fmt.Println()
	fmt.Println("Download Options:")
	fmt.Println("  -o string           Output directory")
	fmt.Println("  -download-limit int Download speed limit in bytes/sec (0 = unlimited)")
	fmt.Println("  -upload-limit int   Upload speed limit in bytes/sec (0 = unlimited)")
	fmt.Println("  -dht                Enable DHT peer discovery (default: true)")
	fmt.Println("  -pex                Enable Peer Exchange (default: true)")
	fmt.Println("  -encrypt            Enable MSE encryption (default: true)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  smirnovtorrent download example.torrent")
	fmt.Println("  smirnovtorrent download \"magnet:?xt=urn:btih:...\"")
	fmt.Println("  smirnovtorrent download file.torrent -o ~/Downloads")
	fmt.Println("  smirnovtorrent download file.torrent -download-limit 1048576 -upload-limit 524288")
	fmt.Println("  smirnovtorrent download file.torrent -dht -pex -encrypt")
	fmt.Println("  smirnovtorrent webui 8080")
	fmt.Println("  smirnovtorrent info example.torrent")
}

func showInfo(path string) {
	torrent, err := parser.ParseFile(path)
	if err != nil {
		fmt.Printf("Error parsing torrent: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Torrent Information")
	fmt.Println("===================")
	fmt.Printf("Name:        %s\n", torrent.Info.Name)
	fmt.Printf("Info Hash:   %s\n", torrent.Info.InfoHash)
	fmt.Printf("Size:        %d bytes\n", torrent.TotalSize())
	fmt.Printf("Piece Size:  %s\n", torrent.PieceSize())
	fmt.Printf("Pieces:      %d\n", len(torrent.Info.Pieces)/20)
	fmt.Printf("Tracker:     %s\n", torrent.Announce)
	fmt.Println()
	fmt.Println("Files:")
	for _, file := range torrent.Info.Files {
		fmt.Printf("  - %s (%d bytes)\n", file.Path, file.Size)
	}
}

func download(source string, cfg *config.Config, outputDir string) {
	appLog.Info("Starting download from: %s", source)

	// Проверяем это magnet ссылка
	if magnet.IsMagnetLink(source) {
		downloadFromMagnet(source, cfg)
		return
	}

	torrent, err := parser.ParseFile(source)
	if err != nil {
		appLog.Error("Error parsing torrent: %v", err)
		fmt.Printf("Error parsing torrent: %v\n", err)
		os.Exit(1)
	}

	// Проверяем валидность торрента
	if err := torrent.IsValid(); err != nil {
		appLog.Error("Invalid torrent: %v", err)
		fmt.Printf("Invalid torrent: %v\n", err)
		os.Exit(1)
	}

	appLog.Info("Torrent: %s, Size: %d bytes", torrent.Info.Name, torrent.TotalSize())
	
	fmt.Printf("Starting download: %s\n", torrent.Info.Name)
	fmt.Printf("Size: %d bytes (%s)\n", torrent.TotalSize(), formatBytesFloat(float64(torrent.TotalSize())))
	fmt.Printf("Pieces: %d\n", len(torrent.Info.Pieces)/20)
	fmt.Printf("Piece size: %s\n", torrent.PieceSize())
	fmt.Printf("Tracker: %s\n", torrent.Announce)
	
	// Показываем лимиты скорости
	if cfg.DownloadRateLimit > 0 {
		fmt.Printf("Download limit: %s/s\n", engine.FormatRate(cfg.DownloadRateLimit))
	} else {
		fmt.Println("Download limit: unlimited")
	}

	if cfg.UploadRateLimit > 0 {
		fmt.Printf("Upload limit: %s/s\n", engine.FormatRate(cfg.UploadRateLimit))
	} else {
		fmt.Println("Upload limit: unlimited")
	}
	
	fmt.Println()

	// Создаём и запускаем движок загрузки
	eng := engine.NewDownloadEngine(torrent, outputDir)
	
	// Применяем конфигурацию
	if cfg.EnableDHT {
		appLog.Info("DHT enabled")
		eng.EnableDHT()
	}
	
	if cfg.EnablePEX {
		appLog.Info("PEX enabled")
		// PEX включается автоматически в PeerPool
	}
	
	if cfg.EnableEncryption {
		appLog.Info("Encryption enabled")
		// Encryption is enabled by default in PeerPool
	}
	
	// Rate limiting
	if cfg.DownloadRateLimit > 0 || cfg.UploadRateLimit > 0 {
		appLog.Info("Setting rate limits: DL=%d, UL=%d bytes/sec", cfg.DownloadRateLimit, cfg.UploadRateLimit)
		eng.SetRateLimits(cfg.DownloadRateLimit, cfg.UploadRateLimit)
	}
	
	// Resume support
	if cfg.EnableResume {
		appLog.Info("Resume support enabled")
		eng.EnableResume()
	}
	
	// Устанавливаем callback для обновления прогресса
	prog := NewProgressBar(40)
	eng.SetProgressCallback(func(progress float64, current, total, peers int, speed float64) {
		prog.Show(progress, current, total, peers, speed)
	})

	appLog.Info("Download started")
	
	if err := eng.Start(); err != nil {
		prog.Finish()
		appLog.Error("Download error: %v", err)
		fmt.Printf("Download error: %v\n", err)
		os.Exit(1)
	}

	prog.Finish()
	appLog.Info("Download completed successfully")
	fmt.Println()
	fmt.Println("Download completed successfully!")
}

// downloadFromMagnet загружает из magnet ссылки
func downloadFromMagnet(magnetLink string, cfg *config.Config) {
	appLog.Info("Magnet download: %s", magnetLink)
	
	fmt.Println("Magnet link download")
	fmt.Println("====================")

	link, err := magnet.Parse(magnetLink)
	if err != nil {
		appLog.Error("Error parsing magnet: %v", err)
		fmt.Printf("Error parsing magnet link: %v\n", err)
		os.Exit(1)
	}

	appLog.Info("Magnet info hash: %s", link.InfoHash)
	
	fmt.Printf("Info Hash: %s\n", link.InfoHash)
	if link.DisplayName != "" {
		fmt.Printf("Name: %s\n", link.DisplayName)
	}
	fmt.Printf("Trackers: %d\n", len(link.Trackers))
	for i, tracker := range link.Trackers {
		fmt.Printf("  %d. %s\n", i+1, tracker)
	}
	if link.DHT {
		fmt.Println("DHT: Enabled")
	}
	if link.PEX {
		fmt.Println("PEX: Enabled")
	}

	fmt.Println()
	
	if !cfg.EnableDHT && !link.DHT {
		appLog.Warn("DHT is disabled, magnet download may fail")
		fmt.Println("Warning: DHT is disabled. Enable with --dht flag or config.")
	}
	
	fmt.Println("Starting DHT client...")

	// Создаём DHT клиент
	dhtClient, err := dht.NewDHTClient(nil, 6882)
	if err != nil {
		appLog.Error("DHT client error: %v", err)
		fmt.Printf("Error creating DHT client: %v\n", err)
		os.Exit(1)
	}

	if err := dhtClient.Start(); err != nil {
		appLog.Error("DHT start error: %v", err)
		fmt.Printf("Error starting DHT: %v\n", err)
		os.Exit(1)
	}

	appLog.Info("DHT started, searching for peers...")
	fmt.Println("DHT started, searching for peers...")

	// Пытаемся получить метаданные через DHT (BEP 9)
	fmt.Println("Attempting metadata download via BEP 9...")
	
	// Получаем пиры через DHT
	peers, err := dhtClient.FindPeer(link.InfoHash)
	if err != nil {
		appLog.Warn("DHT peer discovery: %v", err)
		fmt.Printf("Warning: DHT peer discovery: %v\n", err)
	} else {
		appLog.Info("Found %d peers via DHT", len(peers))
		fmt.Printf("Found %d peers via DHT\n", len(peers))
	}

	// Если есть трекеры, получаем пиры от них
	if len(link.Trackers) > 0 {
		appLog.Info("Contacting %d trackers", len(link.Trackers))
		fmt.Println("Contacting trackers...")
		// В полной реализации здесь нужно получить пиры от трекеров
	}

	fmt.Println()
	fmt.Println("Note: Full magnet download with metadata is in progress.")
	fmt.Println("For now, you need to provide a .torrent file for complete download.")
	
	// Очищаем DHT клиент
	dhtClient.Stop()
	appLog.Info("DHT client stopped")
}

// startWebUI запускает веб-интерфейс
func startWebUI(port int, cfg *config.Config) {
	appLog.Info("Starting Web UI on port %d", port)
	
	fmt.Println("Starting Web UI...")
	fmt.Printf("Open http://localhost:%d in your browser\n", port)
	fmt.Println()

	// Создаём тестовый движок (для демонстрации)
	torrent := &parser.Torrent{
		Info: parser.TorrentInfo{
			Name:        "Demo Torrent",
			InfoHash:    "0000000000000000000000000000000000000000",
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
			Files:       []parser.FileInfo{{Path: "demo.txt", Size: 1024}},
		},
	}

	eng := engine.NewDownloadEngine(torrent, "")
	webui := NewWebUI(eng, port)

	// Запускаем веб-сервер
	appLog.Info("Web UI server starting")
	if err := webui.Start(); err != nil {
		appLog.Error("Web UI error: %v", err)
		fmt.Printf("Web UI error: %v\n", err)
		os.Exit(1)
	}
	appLog.Info("Web UI server stopped")
}

// formatBytesFloat форматирует размер в байтах для вывода (для float64 скорости)
func formatBytesFloat(bytes float64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%.1f B", bytes)
	}
	div, exp := float64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", bytes/div, suffixes[exp])
}

// ProgressBar прогресс бар для CLI
type ProgressBar struct {
	width int
}

// NewProgressBar создаёт новый прогресс бар
func NewProgressBar(width int) *ProgressBar {
	return &ProgressBar{width: width}
}

// Show показывает прогресс
func (pb *ProgressBar) Show(progress float64, current, total, peers int, speed float64) {
	// Простой вывод в консоль
	fmt.Printf("\r[%s] %6.2f%% | Peers: %3d | Speed: %8s/s", 
		strings.Repeat("=", int(progress/100*float64(pb.width))) + ">",
		progress,
		peers,
		formatBytesFloat(speed))
}

// Finish завершает прогресс бар
func (pb *ProgressBar) Finish() {
	fmt.Println()
}

