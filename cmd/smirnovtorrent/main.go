package main

import (
	"fmt"
	"os"

	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/magnet"
	"smirnovtorrent/internal/parser"
)

const version = "0.7.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("SmirnovTorrent v%s\n", version)
	case "download":
		if len(os.Args) < 3 {
			fmt.Println("Error: torrent file or magnet link required")
			fmt.Println("Usage: smirnovtorrent download <file.torrent|magnet>")
			os.Exit(1)
		}
		torrentSource := os.Args[2]
		download(torrentSource)
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
	fmt.Println("  info <file.torrent>             Show torrent information")
	fmt.Println("  version                         Show version")
	fmt.Println("  help                            Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  smirnovtorrent download example.torrent")
	fmt.Println("  smirnovtorrent download \"magnet:?xt=urn:btih:...\"")
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

func download(source string) {
	// Проверяем это magnet ссылка
	if magnet.IsMagnetLink(source) {
		downloadFromMagnet(source)
		return
	}

	torrent, err := parser.ParseFile(source)
	if err != nil {
		fmt.Printf("Error parsing torrent: %v\n", err)
		os.Exit(1)
	}

	// Проверяем валидность торрента
	if err := torrent.IsValid(); err != nil {
		fmt.Printf("Invalid torrent: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting download: %s\n", torrent.Info.Name)
	fmt.Printf("Size: %d bytes (%s)\n", torrent.TotalSize(), formatBytes(float64(torrent.TotalSize())))
	fmt.Printf("Pieces: %d\n", len(torrent.Info.Pieces)/20)
	fmt.Printf("Piece size: %s\n", torrent.PieceSize())
	fmt.Printf("Tracker: %s\n", torrent.Announce)
	fmt.Println()

	// Создаём и запускаем движок загрузки
	eng := engine.NewDownloadEngine(torrent, "")
	
	// Устанавливаем callback для обновления прогресса
	prog := NewProgressBar(40)
	eng.SetProgressCallback(func(progress float64, current, total, peers int, speed float64) {
		prog.Show(progress, current, total, peers, speed)
	})

	if err := eng.Start(); err != nil {
		prog.Finish()
		fmt.Printf("Download error: %v\n", err)
		os.Exit(1)
	}

	prog.Finish()
	fmt.Println()
	fmt.Println("Download completed successfully!")
}

// downloadFromMagnet загружает из magnet ссылки
func downloadFromMagnet(magnetLink string) {
	fmt.Println("Magnet link download")
	fmt.Println("====================")

	link, err := magnet.Parse(magnetLink)
	if err != nil {
		fmt.Printf("Error parsing magnet link: %v\n", err)
		os.Exit(1)
	}

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
	fmt.Println("Note: Full magnet link support with DHT is coming soon!")
	fmt.Println("For now, you need to provide a .torrent file.")
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
