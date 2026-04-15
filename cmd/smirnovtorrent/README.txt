SmirnovTorrent - Инструкция
=============================

Быстрый старт:
--------------
1. Просто запустите start.bat - откроется веб-интерфейс
2. Или выполните в командной строке: smirnovtorrent.exe webui
3. Откройте браузер по адресу: http://localhost:8080

Использование:
--------------
- Нажмите "Обзор" и выберите .torrent файл
- Нажмите "Добавить" - загрузка начнётся автоматически
- После завершения можно открыть папку с файлами

Команды:
--------
smirnovtorrent.exe webui           - Запустить веб-интерфейс
smirnovtorrent.exe download file.torrent - Скачать торрент
smirnovtorrent.exe info file.torrent - Показать информацию
smirnovtorrent.exe version         - Версия программы
smirnovtorrent.exe help            - Справка

Настройки:
----------
Создайте файл smirnovtorrent.json в папке с программой:

{
  "WebUIPort": 8080,
  "DownloadRateLimit": 0,
  "UploadRateLimit": 0,
  "EnableDHT": true,
  "EnablePEX": true,
  "EnableEncryption": true
}

Параметры:
- WebUIPort: порт веб-интерфейса (по умолчанию 8080)
- DownloadRateLimit: лимит загрузки байт/сек (0 = без ограничений)
- UploadRateLimit: лимит отдачи байт/сек (0 = без ограничений)
- EnableDHT: включить DHT сеть
- EnablePEX: включить обмен пирами
- EnableEncryption: включить шифрование

Версия: 1.0.0
Сайт: https://github.com/ainsonet/smirnovtorrent