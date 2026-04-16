# [SmirnovTorrent](https://ainsonet.github.io/smirnovtorrent/)

[![Версия](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/ainsonet/smirnovtorrent/releases)
[![Лицензия](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24.0-00ADD8.svg?logo=go)](https://go.dev/)

Лёгкий, современный BitTorrent клиент с веб-интерфейсом, написанный на Go.

<img src="https://raw.githubusercontent.com/ainsonet/smirnovtorrent/master/logo.png" alt="SmirnovTorrent Logo" width="200"/>

## Особенности

- Высокая производительность - Оптимизирован для быстрых загрузок
- Современный веб-интерфейс - Красивый UI с поддержкой тёмной/светлой темы
- Безопасность - Поддержка MSE/PE шифрования
- DHT и PEX - Децентрализованное обнаружение пиров
- Продолжение загрузок - Сохранение прогресса при остановке
- Статистика в реальном времени - Скорости, пиры, прогресс
- Удобная организация - Автоматическая сортировка файлов
- Портативность - Не требует установки

## Содержание

- [SmirnovTorrent](#smirnovtorrent)
  - [Особенности](#особенности)
  - [Содержание](#содержание)
  - [Быстрый старт](#быстрый-старт)
    - [Самый простой способ (Windows)](#самый-простой-способ-windows)
    - [Windows (готовый исполняемый файл)](#windows-готовый-исполняемый-файл)
    - [Из командной строки](#из-командной-строки)
  - [Установка](#установка)
    - [Windows](#windows)
    - [Из исходников](#из-исходников)
  - [Использование](#использование)
    - [Веб-интерфейс (рекомендуется)](#веб-интерфейс-рекомендуется)
    - [Командная строка](#командная-строка)
  - [Конфигурация](#конфигурация)
    - [Файл конфигурации](#файл-конфигурации)
    - [Параметры](#параметры)
  - [Сборка из исходников](#сборка-из-исходников)
    - [Требования](#требования)
    - [Сборка](#сборка)
    - [Проверка сборки](#проверка-сборки)
  - [Структура проекта](#структура-проекта)
  - [Документация API](#документация-api)
    - [Эндпоинты веб-интерфейса](#эндпоинты-веб-интерфейса)
    - [Пример ответа API](#пример-ответа-api)
  - [Вклад](#вклад)
    - [Как внести вклад](#как-внести-вклад)
    - [Правила разработки](#правила-разработки)
  - [Лицензия](#лицензия)
  - [Благодарности](#благодарности)
  - [Контакты](#контакты)

## Быстрый старт

### Самый простой способ (Windows)

Просто запустите **start.bat** в папке с программой - веб-интерфейс откроется автоматически!

```
cmd/smirnovtorrent/start.bat
```

### Windows (готовый исполняемый файл)

1. [Скачайте последнюю версию](https://github.com/ainsonet/smirnovtorrent/releases)
2. Распакуйте архив
3. Запустите `smirnovtorrent.exe` или `start.bat`
4. Откройте браузер по адресу `http://localhost:8080`

### Из командной строки

```bash
# Загрузить торрент файл
smirnovtorrent download example.torrent

# Запустить веб-интерфейс
smirnovtorrent webui

# Показать информацию о торренте
smirnovtorrent info example.torrent
```

## Установка

### Windows

Скачайте последнюю версию с [GitHub Releases](https://github.com/ainsonet/smirnovtorrent/releases):

```powershell
# Используя PowerShell
Invoke-WebRequest -Uri "https://github.com/ainsonet/smirnovtorrent/releases/download/v1.0.0/smirnovtorrent-windows.zip" -OutFile "smirnovtorrent.zip"
Expand-Archive smirnovtorrent.zip
.\smirnovtorrent\smirnovtorrent.exe webui
```

### Из исходников

Требуется **Go 1.24.0** или новее:

```bash
# Клонируйте репозиторий
git clone https://github.com/ainsonet/smirnovtorrent.git
cd smirnovtorrent

# Скомпилируйте
go build -o smirnovtorrent.exe ./cmd/smirnovtorrent

# Запустите
.\smirnovtorrent.exe webui
```

## Использование

### Веб-интерфейс (рекомендуется)

**Самый удобный способ** - запустить `start.bat` или выполнить:

```bash
smirnovtorrent webui
```

Веб-интерфейс предоставляет удобный GUI для управления загрузками:

1. **Откройте браузер:**
   - Перейдите по адресу `http://localhost:8080`
   - Браузер откроется автоматически (Windows)

2. **Добавление торрента:**
   - Нажмите кнопку **"Обзор"** (folder icon)
   - Выберите `.torrent` файл
   - Нажмите **"Добавить"**
   - Загрузка начнётся автоматически

3. **Управление:**
   - **Pause** - Приостановить загрузку
   - **Open Folder** - Открыть папку с файлами (после завершения)
   - **Remove** - Удалить торрент

### Командная строка

Полный список команд:

```bash
# Показать справку
smirnovtorrent help

# Скачать торрент
smirnovtorrent download <файл.torrent|magnet> [опции]

# Опции для загрузки:
#   -o string         Выходная директория
#   -download-limit   Лимит загрузки (байт/сек, 0 = без ограничений)
#   -upload-limit     Лимит отдачи (байт/сек, 0 = без ограничений)
#   -dht              Включить DHT (по умолчанию: true)
#   -pex              Включить PEX (по умолчанию: true)
#   -encrypt          Включить шифрование (по умолчанию: true)

# Примеры:
smirnovtorrent download example.torrent
smirnovtorrent download file.torrent -o ~/Downloads
smirnovtorrent download file.torrent -download-limit 1048576 -upload-limit 524288
smirnovtorrent download "magnet:?xt=urn:btih:..."

# Запустить веб-интерфейс
smirnovtorrent webui          # Порт 8080 (по умолчанию)
smirnovtorrent webui 9000     # Порт 9000

# Информация о торренте
smirnovtorrent info example.torrent

# Показать версию
smirnovtorrent version
```

## Конфигурация

### Файл конфигурации

Создайте файл `smirnovtorrent.json` в директории приложения:

```json
{
  "WebUIPort": 8080,
  "DownloadRateLimit": 0,
  "UploadRateLimit": 0,
  "EnableDHT": true,
  "EnablePEX": true,
  "EnableEncryption": true,
  "DefaultDownloadDir": "C:\\Users\\user\\Downloads"
}
```

### Параметры

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|--------------|----------|
| `WebUIPort` | int | `8080` | Порт веб-интерфейса |
| `DownloadRateLimit` | int | `0` | Лимит загрузки (байт/сек) |
| `UploadRateLimit` | int | `0` | Лимит отдачи (байт/сек) |
| `EnableDHT` | bool | `true` | Включить DHT |
| `EnablePEX` | bool | `true` | Включить Peer Exchange |
| `EnableEncryption` | bool | `true` | Включить MSE шифрование |
| `DefaultDownloadDir` | string | `.` | Директория по умолчанию |

## Сборка из исходников

### Требования

- **Go** 1.24.0 или новее
- **Git** (для клонирования репозитория)

### Сборка

```bash
# Клонируйте репозиторий
git clone https://github.com/ainsonet/smirnovtorrent.git
cd smirnovtorrent

# Обычная сборка
go build -o smirnovtorrent.exe ./cmd/smirnovtorrent

# Сборка с оптимизациями
go build -ldflags="-s -w" -o smirnovtorrent.exe ./cmd/smirnovtorrent

# Сборка инсталлятора
powershell -ExecutionPolicy Bypass -File build-installer.ps1
```

### Проверка сборки

```bash
# Запустите тесты
go test ./...

# Проверьте форматирование
go fmt ./...

# Проверьте линтером
golangci-lint run
```

## Структура проекта

```
SmirnovTorrent/
├── cmd/
│   └── smirnovtorrent/
│       ├── main.go           # Точка входа CLI
│       ├── webui.go          # Веб-сервер
│       ├── webui.html        # Веб-интерфейс
│       └── start.bat         # Быстрый запуск GUI
│
├── internal/
│   ├── config/               # Управление конфигурацией
│   ├── engine/
│   │   ├── anacrolix_engine.go  # Движок загрузки
│   │   └── resume.go         # Сохранение прогресса
│   ├── logger/               # Логирование
│   ├── parser/               # Парсер .torrent файлов
│   └── tracker/              # Tracker клиент
│
├── pkg/
│   └── bencode/              # Bencode кодирование
│
├── website/                  # Официальный сайт
│   ├── index.html
│   ├── styles.css
│   ├── script.js
│   └── logo.png
│
├── build-installer.ps1       # Создание инсталлятора
├── go.mod                    # Зависимости Go
├── LICENSE                   # Лицензия MIT
└── README.md                 # Документация
```

## Документация API

### Эндпоинты веб-интерфейса

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/` | Главная страница |
| `GET` | `/api/status` | Статус загрузки |
| `POST` | `/api/add` | Добавить торрент |
| `POST` | `/api/start` | Запустить загрузку |
| `POST` | `/api/stop` | Остановить загрузку |
| `POST` | `/api/pause` | Приостановить |
| `POST` | `/api/resume` | Возобновить |
| `POST` | `/api/remove` | Удалить торрент |
| `POST` | `/api/open-folder` | Открыть папку |
| `GET` | `/logo.png` | Логотип |

### Пример ответа API

```json
{
  "status": "completed",
  "progress": 100,
  "downloaded": 5368709120,
  "uploaded": 2147483648,
  "totalSize": 5368709120,
  "activePeers": 15,
  "downloadSpeed": 0,
  "uploadSpeed": 524288,
  "torrentName": "example-file.iso",
  "path": "C:\\Downloads\\example.torrent"
}
```

## Вклад

Мы приветствуем вклад! Пожалуйста, прочитайте наши правила перед созданием pull request'ов.

### Как внести вклад

1. **Форкните** репозиторий
2. **Создайте ветку** для вашей фичи:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Сделайте изменения** и запишите тесты
4. **Закоммитьте** изменения:
   ```bash
   git commit -m 'Добавлена amazing фича'
   ```
5. **Пушните** в ветку:
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Откройте Pull Request**

### Правила разработки

- Следуйте стилю кода Go
- Пишите понятные комментарии
- Добавляйте тесты для новых функций
- Обновляйте документацию при необходимости

## Лицензия

Этот проект распространяется под лицензией **MIT**. См. файл [LICENSE](LICENSE) для деталей.

## Благодарности

- [anacrolix/torrent](https://github.com/anacrolix/torrent) - Отличный BitTorrent клиент на Go
- [GoLang](https://golang.org/) - Великолепный язык программирования
- Все контрибьюторы проекта

## Контакты

- **Автор**: Дмитрий Смирнов
- **GitHub**: [ainsonet/smirnovtorrent](https://github.com/ainsonet/smirnovtorrent)

---

**Разработано в России**
