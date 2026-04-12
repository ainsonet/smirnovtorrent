package parser

import (
	"errors"
	"os"

	"smirnovtorrent/pkg/bencode"
)

// ParseFile читает и парсит .torrent файл по пути
func ParseFile(path string) (*Torrent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse парсит bencode данные .torrent файла
func Parse(data []byte) (*Torrent, error) {
	val, err := bencode.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	dict, ok := val.(bencode.Dict)
	if !ok {
		return nil, errors.New("root must be a dictionary")
	}

	torrent := &Torrent{}

	// Announce URL
	if announce, exists := dict["announce"]; exists {
		if str, ok := announce.(bencode.String); ok {
			torrent.Announce = string(str)
		}
	}

	// Announce list (backup trackers)
	if announceList, exists := dict["announce-list"]; exists {
		if list, ok := announceList.(bencode.List); ok {
			for _, tier := range list {
				if tierList, ok := tier.(bencode.List); ok {
					var urls []string
					for _, url := range tierList {
						if str, ok := url.(bencode.String); ok {
							urls = append(urls, string(str))
						}
					}
					if len(urls) > 0 {
						torrent.AnnounceList = append(torrent.AnnounceList, urls)
					}
				}
			}
		}
	}

	// Info dictionary
	if info, exists := dict["info"]; exists {
		if infoDict, ok := info.(bencode.Dict); ok {
			infoBytes, err := bencode.Marshal(infoDict)
			if err != nil {
				return nil, err
			}

			torrent.Info = parseInfo(infoDict)
			torrent.Info.InfoHash = CalculateInfoHash(infoBytes)
		}
	}

	// Optional fields
	if createdBy, exists := dict["created by"]; exists {
		if str, ok := createdBy.(bencode.String); ok {
			torrent.CreatedBy = string(str)
		}
	}

	if creationDate, exists := dict["creation date"]; exists {
		if date, ok := creationDate.(bencode.Int); ok {
			torrent.CreationDate = int64(date)
		}
	}

	if comment, exists := dict["comment"]; exists {
		if str, ok := comment.(bencode.String); ok {
			torrent.Comment = string(str)
		}
	}

	if encoding, exists := dict["encoding"]; exists {
		if str, ok := encoding.(bencode.String); ok {
			torrent.Encoding = string(str)
		}
	}

	return torrent, nil
}

func parseInfo(infoDict bencode.Dict) TorrentInfo {
	info := TorrentInfo{}

	// Name
	if name, exists := infoDict["name"]; exists {
		if str, ok := name.(bencode.String); ok {
			info.Name = string(str)
		}
	}

	// Piece length
	if pieceLength, exists := infoDict["piece length"]; exists {
		if length, ok := pieceLength.(bencode.Int); ok {
			info.PieceLength = int(length)
		}
	}

	// Pieces
	if pieces, exists := infoDict["pieces"]; exists {
		if p, ok := pieces.(bencode.String); ok {
			info.Pieces = []byte(p)
		}
	}

	// Single file mode
	if length, exists := infoDict["length"]; exists {
		if len(info.Files) == 0 {
			// Это single-file torrent
			if l, ok := length.(bencode.Int); ok {
				info.Files = append(info.Files, FileInfo{
					Path: info.Name,
					Size: int64(l),
				})
			}
		}
	}

	// Multi-file mode
	if files, exists := infoDict["files"]; exists {
		if fileList, ok := files.(bencode.List); ok {
			for _, fileVal := range fileList {
				if fileDict, ok := fileVal.(bencode.Dict); ok {
					file := parseFileInfo(fileDict, info.Name)
					info.Files = append(info.Files, file)
				}
			}
		}
	}

	return info
}

func parseFileInfo(fileDict bencode.Dict, torrentName string) FileInfo {
	file := FileInfo{}

	// Length
	if length, exists := fileDict["length"]; exists {
		if l, ok := length.(bencode.Int); ok {
			file.Size = int64(l)
		}
	}

	// Path
	if path, exists := fileDict["path"]; exists {
		if pathList, ok := path.(bencode.List); ok {
			var pathStr string
			for i, part := range pathList {
				if str, ok := part.(bencode.String); ok {
					if i > 0 {
						pathStr += "/"
					}
					pathStr += string(str)
				}
			}
			file.Path = pathStr
		}
	}

	// MD5 sum (optional)
	if md5sum, exists := fileDict["md5sum"]; exists {
		if str, ok := md5sum.(bencode.String); ok {
			file.MD5Sum = string(str)
		}
	}

	// Если путь пустой, используем имя торрента
	if file.Path == "" {
		file.Path = torrentName
	}

	return file
}

// IsValid проверяет базовую валидность торрента
func (t *Torrent) IsValid() error {
	if t.Announce == "" && len(t.AnnounceList) == 0 {
		return errors.New("no tracker URL found")
	}

	if t.Info.PieceLength == 0 {
		return errors.New("piece length is required")
	}

	if len(t.Info.Pieces)%20 != 0 {
		return errors.New("pieces length must be multiple of 20 (SHA-1 hash size)")
	}

	if len(t.Info.Files) == 0 {
		return errors.New("no files in torrent")
	}

	return nil
}
