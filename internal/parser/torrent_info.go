package parser

import (
	"crypto/sha1"
	"encoding/hex"
)

// TorrentInfo содержит информацию из .torrent файла
type TorrentInfo struct {
	// Хеш инфы (info dictionary)
	InfoHash string

	// Имя торрента
	Name string

	// Размер каждого куска (в байтах)
	PieceLength int

	// Список кусков (каждый 20 байт = SHA-1 хеш)
	Pieces []byte

	// Файлы в торренте
	Files []FileInfo

	// Аннотация (опционально)
	Comment string

	// Создатель (опционально)
	Creator string
}

// FileInfo информация о файле внутри торрента
type FileInfo struct {
	// Путь файла (может быть многоуровневым)
	Path string

	// Размер файла в байтах
	Size int64

	// MD5 хеш (опционально)
	MD5Sum string
}

// Torrent структура всего .torrent файла
type Torrent struct {
	// Announce URL трекера
	Announce string

	// Список URL трекеров (backup)
	AnnounceList [][]string

	// Meta info (основная информация)
	Info TorrentInfo

	// Создавшая программа
	CreatedBy string

	// Время создания (Unix timestamp)
	CreationDate int64

	// Кодировка имен файлов
	Encoding string

	// Комментарий (опционально)
	Comment string
}

// PieceSize возвращает размер кусков в более понятном формате
func (t *Torrent) PieceSize() string {
	return formatBytes(t.Info.PieceLength)
}

// TotalSize возвращает общий размер всех файлов
func (t *Torrent) TotalSize() int64 {
	var total int64
	for _, f := range t.Info.Files {
		total += f.Size
	}
	return total
}

// InfoHashBytes возвращает InfoHash в виде байтового массива
func (t *Torrent) InfoHashBytes() []byte {
	hash, _ := hex.DecodeString(t.Info.InfoHash)
	return hash
}

// Helper functions

func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string(rune('K' + exp)) + "." + string(rune(bytes/int(div)%unit)) + "B"
}

// CalculateInfoHash вычисляет хеш info словаря
func CalculateInfoHash(infoDict []byte) string {
	hash := sha1.Sum(infoDict)
	return hex.EncodeToString(hash[:])
}
