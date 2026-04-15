package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"smirnovtorrent/internal/config"
	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/logger"
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

	// Проверяем что файл существует
	if _, err := os.Stat(source); os.IsNotExist(err) {
		appLog.Error("Torrent file not found: %s", source)
		fmt.Printf("Error: Torrent file not found: %s\n", source)
		os.Exit(1)
	}

	// Проверяем расширение
	if !strings.HasSuffix(strings.ToLower(source), ".torrent") {
		appLog.Error("Not a .torrent file: %s", source)
		fmt.Printf("Error: Not a .torrent file: %s\n", source)
		os.Exit(1)
	}

	// Если outputDir не указан, используем текущую директорию
	if outputDir == "" {
		outputDir = "."
	}

	fmt.Printf("Starting download: %s\n", filepath.Base(source))
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println()

	// Создаём и запускаем движок загрузки
	eng, err := engine.NewAnacrolixEngine(source, outputDir)
	if err != nil {
		appLog.Error("Engine creation error: %v", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer eng.Stop()

	// Загружаем торрент
	if err := eng.LoadTorrent(source); err != nil {
		appLog.Error("Torrent load error: %v", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
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
	fmt.Printf("Files saved to: %s\n", outputDir)
}

// startWebUI запускает веб-интерфейс
func startWebUI(port int, cfg *config.Config) {
	appLog.Info("Starting Web UI on port %d", port)

	fmt.Println("Starting Web UI...")
	fmt.Printf("Open http://localhost:%d in your browser\n", port)
	fmt.Println("Add .torrent files via the interface to start downloading.")
	fmt.Println()

	webui := NewWebUI(port)

	// Запускаем веб-сервер (он откроет браузер автоматически)
	appLog.Info("Web UI server starting")
	if err := webui.Start(); err != nil {
		appLog.Error("Web UI error: %v", err)
		fmt.Printf("Web UI error: %v\n", err)
		os.Exit(1)
	}
	appLog.Info("Web UI server stopped")
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
	speedStr := formatSpeed(speed)
	fmt.Printf("\r[%s] %6.2f%% | Peers: %3d | Speed: %8s/s", 
		strings.Repeat("=", int(progress/100*float64(pb.width))) + ">",
		progress,
		peers,
		speedStr)
}

// formatSpeed форматирует скорость
func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec < 1024 {
		return fmt.Sprintf("%.1f B", bytesPerSec)
	}
	div, exp := float64(1024), 0
	for n := bytesPerSec / 1024; n >= 1024; n /= 1024 {
		div *= 1024
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", bytesPerSec/div, suffixes[exp])
}

// Finish завершает прогресс бар
func (pb *ProgressBar) Finish() {
	fmt.Println()
}

