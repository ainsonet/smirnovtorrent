# SmirnovTorrent - Быстрая установка

## 🚀 Установка (Windows)

### Вариант 1: Готовый установщик (Рекомендуется)

1. Скачайте установщик: `SmirnovTorrent_1.0.0_x64-setup.exe`
2. Запустите установщик
3. Следуйте мастеру установки
4. Готово! Приложение доступно в меню Пуск

### Вариант 2: Сборка из исходников

Если хотите собрать последнюю версию:

```powershell
# Запустите скрипт сборки
.\build-installer.ps1
```

Скрипт автоматически:
- Проверит наличие Node.js, npm и Rust
- Установит зависимости
- Соберёт frontend
- Создаст .exe установщик

**Требования для сборки:**
- Node.js >= 18.0.0
- Rust >= 1.70.0
- Windows 10/11

## 💻 Использование

### После установки:

1. **Запуск приложения**
   - Через меню Пуск → SmirnovTorrent
   - Или ярлык на рабочем столе

2. **Добавление торрента**
   - Нажмите "Browse" и выберите .torrent файл
   - Или вставьте magnet ссылку
   - Нажмите "Add"

3. **Управление загрузкой**
   - Pause - приостановить
   - Resume - продолжить
   - Remove - удалить

4. **Папка загрузок**
   - Файлы сохраняются в `C:\Users\[Имя]\Downloads\SmirnovTorrent`
   - Можно изменить в настройках

## ⚙️ Настройки

Приложение создаёт конфиг в:
```
C:\Users\[Имя]\AppData\Roaming\smirnovtorrent\config.json
```

### Пример конфигурации:
```json
{
  "download_rate_limit": 1048576,
  "upload_rate_limit": 524288,
  "enable_dht": true,
  "enable_pex": true,
  "enable_encryption": true,
  "output_dir": "C:/Downloads"
}
```

## 🔧 Troubleshooting

### Установщик не запускается
- Убедитесь, что у вас Windows 10/11
- Проверьте наличие WebView2: https://developer.microsoft.com/en-us/microsoft-edge/webview2/

### Ошибка при сборке
```powershell
# Обновите Rust
rustup update

# Очистите кэш
cargo clean

# Попробуйте снова
.\build-installer.ps1
```

### Антивирус блокирует
- Добавьте приложение в исключения
- Или используйте подписанный установщик (в разработке)

## 📊 Системные требования

**Минимальные:**
- Windows 10
- 2 GB RAM
- 100 MB свободного места

**Рекомендуемые:**
- Windows 11
- 4 GB RAM
- SSD для быстрых загрузок

## 🆘 Поддержка

- Документация: https://github.com/ainsonet/smirnovtorrent
- Issues: https://github.com/ainsonet/smirnovtorrent/issues
- Email: support@smirnovtorrent.com

---

**Приятного использования!** 🎉
