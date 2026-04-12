# SmirnovTorrent

🌊 Легковесный BitTorrent клиент на Go

## 🚀 Особенности

- Парсинг .torrent файлов
- Подключение к трекерам
- Загрузка из нескольких источников (пиры)
- CLI интерфейс
- Проверка целостности данных (SHA-1)

## 📦 Установка

```bash
go build -o smirnovtorrent.exe ./cmd/smirnovtorrent
```

## 💻 Использование

```bash
# Загрузка торрента
smirnovtorrent download example.torrent

# Загрузка с magnet ссылки
smirnovtorrent download "magnet:?xt=urn:btih:..."
```

## 🏗️ Структура проекта

```
smirnovtorrent/
├── cmd/
│   └── smirnovtorrent/     # Точка входа (main.go)
├── internal/
│   ├── parser/             # Парсинг .torrent файлов
│   ├── tracker/            # Работа с трекерами
│   ├── peer/               # Протокол общения с пирами
│   └── engine/             # Основной движок загрузки
└── pkg/
    └── bencode/            # Bencode сериализация/десериализация
```

## 📝 План разработки

- [x] Структура проекта
- [ ] Парсер Bencode
- [ ] Парсер .torrent файлов
- [ ] Работа с трекерами (HTTP)
- [ ] Peer protocol (handshake, messages)
- [ ] Piece manager
- [ ] CLI интерфейс
- [ ] Поддержка нескольких пиров
- [ ] Seed-режим

## 📄 License

MIT
