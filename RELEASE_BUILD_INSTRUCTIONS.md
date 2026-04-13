# 📦 SmirnovTorrent v1.0.0 - Release Build Instructions

## ✅ Готовые установщики

### Windows
- **NSIS Installer**: `SmirnovTorrent_1.0.0_x64-setup.exe` (рекомендуется)
- **MSI Installer**: `SmirnovTorrent_1.0.0_x64_en-US.msi`

---

## 🚀 Сборка установщика (Windows)

### Быстрая сборка

```powershell
# Запустите скрипт сборки
.\build-installer.ps1
```

Скрипт автоматически:
1. Проверит Node.js, npm, Rust
2. Установит зависимости
3. Соберёт frontend
4. Создаст .exe установщик

### Ручная сборка (по шагам)

#### 1. Проверка требований

```powershell
# Node.js (>= 18.0.0)
node --version

# npm (>= 9.0.0)
npm --version

# Rust (>= 1.70.0)
rustc --version
```

#### 2. Установка зависимостей

```powershell
cd gui
npm install
```

#### 3. Сборка frontend

```powershell
npm run build
```

#### 4. Сборка установщика

```powershell
# NSIS установщик (.exe)
npm run tauri build

# Или конкретный формат
npm run tauri build -- --bundles nsis

# MSI установщик
npm run tauri build -- --bundles msi
```

---

## 📁 Результаты сборки

После сборки файлы появятся в:

```
gui/src-tauri/target/release/bundle/
├── nsis/
│   └── SmirnovTorrent_1.0.0_x64-setup.exe    # Рекомендуемый
├── msi/
│   └── SmirnovTorrent_1.0.0_x64_en-US.msi
└── app/
    └── SmirnovTorrent.exe                     # Portable версия
```

---

## ⚙️ Настройка установщика

### tauri.conf.json

```json
{
  "tauri": {
    "bundle": {
      "identifier": "com.smirnovtorrent.desktop",
      "windows": {
        "nsis": {
          "installMode": "currentUser",
          "languages": ["English", "Russian"],
          "displayLanguageSelector": true
        }
      }
    }
  }
}
```

### Изменение иконки

1. Поместите `icon.ico` (256x256) в `gui/src-tauri/icons/`
2. Обновите `tauri.conf.json`:
```json
"icon": ["icons/icon.ico"]
```

### Добавление лицензии

Поместите `LICENSE` в корень проекта и укажите в `tauri.conf.json`:
```json
"nsis": {
  "license": "LICENSE"
}
```

---

## 🌐 Мультиязычность

Установщик поддерживает:
- English
- Russian

Язык выбирается автоматически по системе или вручную.

### Добавить язык

В `tauri.conf.json`:
```json
"nsis": {
  "languages": ["English", "Russian", "Spanish"]
}
```

---

## 📊 Размер установщика

Ожидаемый размер:
- **NSIS**: ~15-25 MB
- **MSI**: ~20-30 MB
- **Portable**: ~10-15 MB

Размер зависит от включённых зависимостей.

---

## 🔐 Подпись установщика (Production)

Для подписи кода в Windows:

```json
"windows": {
  "certificateThumbprint": "YOUR_CERT_THUMBPRINT",
  "digestAlgorithm": "sha256",
  "timestampUrl": "http://timestamp.digicert.com"
}
```

Получить сертификат можно в удостоверяющих центрах:
- DigiCert
- Sectigo
- GlobalSign

---

## 🧪 Тестирование установщика

### 1. Чистая установка

```powershell
# Удалите предыдущую версию
Uninstall-Package SmirnovTorrent

# Запустите установщик
.\SmirnovTorrent_1.0.0_x64-setup.exe

# Проверьте установку
Test-Path "$env:PROGRAMFILES\SmirnovTorrent"
```

### 2. Проверка функциональности

- [ ] Запуск приложения
- [ ] Добавление торрента
- [ ] Загрузка файлов
- [ ] Пауза/Возобновление
- [ ] Настройки
- [ ] Завершение работы

### 3. Проверка реестра

```powershell
# Проверка записи в реестре
Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* |
  Where-Object { $_.DisplayName -like "*SmirnovTorrent*" }
```

---

## 📝 Чеклист перед релизом

- [ ] Все тесты проходят (`go test ./...`)
- [ ] GUI собирается без ошибок
- [ ] Установщик создаётся корректно
- [ ] Иконки отображаются правильно
- [ ] Языки работают
- [ ] Лицензия включена
- [ ] Размер в норме (< 30 MB)
- [ ] Установка проходит успешно
- [ ] Удаление работает
- [ ] Документация обновлена

---

## 🎯 Публикация релиза

### GitHub Releases

1. Создайте тег:
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. Создайте релиз на GitHub:
   - Перейдите в Releases
   - Draft a new release
   - Tag: v1.0.0
   - Attach установщики:
     - `SmirnovTorrent_1.0.0_x64-setup.exe`
     - `SmirnovTorrent_1.0.0_x64_en-US.msi`
   - Добавьте описание
   - Publish

### Chocolatey (опционально)

```powershell
# Создайте пакет
choco pack smirnovtorrent.nuspec

# Опубликуйте
choco push smirnovtorrent.1.0.0.nupkg --api-key=YOUR_KEY
```

---

## 🐛 Troubleshooting

### Ошибка: "Rust not found"
```powershell
# Переустановите Rust
rustup self uninstall
rustup install stable
```

### Ошибка: "WebView2 not found"
Скачайте: https://developer.microsoft.com/en-us/microsoft-edge/webview2/

### Ошибка сборки Tauri
```powershell
cd gui
cargo clean
npm run tauri build
```

### Установщик не создаётся
Проверьте:
- Node.js >= 18.0.0
- Rust >= 1.70.0
- Все зависимости установлены (`npm install`)

---

## 📞 Поддержка

- Issues: https://github.com/ainsonet/smirnovtorrent/issues
- Docs: https://github.com/ainsonet/smirnovtorrent/tree/main/gui

---

**Удачной сборки!** 🚀
