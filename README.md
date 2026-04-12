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
# Показать информацию о торренте
smirnovtorrent info example.torrent

# Загрузка торрента (в разработке)
smirnovtorrent download example.torrent

# Загрузка с magnet ссылки (в разработке)
smirnovtorrent download "magnet:?xt=urn:btih:..."

# Показать версию
smirnovtorrent version
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
- [x] Парсер Bencode
- [x] Парсер .torrent файлов
- [x] Работа с трекерами (HTTP)
- [x] Peer protocol (handshake, messages)
- [x] Piece manager
- [x] Download Engine (базовый)
- [x] Multi-file torrent support
- [ ] CLI интерфейс (полный)
- [ ] Поддержка нескольких пиров (параллельная загрузка)
- [ ] Seed-режим

## 📄 License

MIT
