package main

import (
	"fmt"
	"os"

	"smirnovtorrent/internal/engine"
	"smirnovtorrent/internal/parser"
)

const version = "0.4.0"

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
	if len(source) > 8 && source[:8] == "magnet:?" {
		fmt.Println("Magnet link support coming soon!")
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
	fmt.Printf("Size: %d bytes\n", torrent.TotalSize())
	fmt.Printf("Pieces: %d\n", len(torrent.Info.Pieces)/20)
	fmt.Printf("Tracker: %s\n", torrent.Announce)
	fmt.Println()

	// Создаём и запускаем движок загрузки
	eng := engine.NewDownloadEngine(torrent, "")
	
	if err := eng.Start(); err != nil {
		fmt.Printf("Download error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Download completed successfully!")
}
