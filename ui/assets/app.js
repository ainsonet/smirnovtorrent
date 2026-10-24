// SmirnovTorrent UI Controller (v1.1.0)

const i18n = {
  "ru": {
    "drag_title": "Добавить торрент",
    "drag_desc": "Отпустите .torrent файл или Magnet-ссылку для добавления",
    "btn_resume": "Продолжить",
    "btn_resume_title": "Возобновить выбранный торрент",
    "btn_pause": "Пауза",
    "btn_pause_title": "Приостановить выбранный торрент",
    "btn_delete": "Удалить",
    "btn_delete_title": "Удалить торрент (Delete)",
    "btn_folder": "Папка",
    "btn_folder_title": "Открыть папку загрузки в проводнике",
    "btn_sequential": "Последовательно",
    "btn_sequential_title": "Последовательная загрузка (для потокового просмотра)",
    "btn_verify": "Проверить",
    "btn_verify_title": "Перепроверить целостность данных",
    "btn_redownload": "Заново",
    "btn_redownload_title": "Скачать заново / Исправить отсутствующие файлы",
    "btn_settings": "Настройки",
    "btn_settings_title": "Настройки программы",
    "search_placeholder": "Поиск торрентов...",
    "cat_title": "Категории",
    "cat_all": "Все торренты",
    "cat_downloading": "Загружаются",
    "cat_seeding": "Раздаются",
    "cat_completed": "Завершены",
    "cat_paused": "На паузе",
    "cat_missing": "Отсутствуют",
    "engine_active": "Клиент активен",
    "col_name": "Имя файла",
    "col_size": "Размер",
    "col_progress": "Прогресс",
    "col_status": "Статус",
    "col_dlspeed": "Загрузка",
    "col_ulspeed": "Отдача",
    "col_eta": "Осталось",
    "col_peers": "Пиры",
    "col_seeds": "Сиды",
    "col_ratio": "Рейтинг",
    "col_added": "Добавлен",
    "empty_title": "Здесь пока нет торрентов",
    "empty_desc": "Перетащите .torrent файл из проводника или воспользуйтесь кнопками ниже",
    "empty_select_file": "Выбрать .torrent файл",
    "empty_paste_magnet": "Вставить Magnet-ссылку",
    "tab_general": "Общие",
    "tab_files": "Файлы",
    "tab_peers": "Пиры",
    "tab_trackers": "Трекеры",
    "tab_pieces": "Карта частей",
    "sec_transfer": "Передача данных",
    "lbl_downloaded": "Загружено:",
    "lbl_uploaded": "Отдано:",
    "lbl_dl_speed": "Скорость загрузки:",
    "lbl_ul_speed": "Скорость отдачи:",
    "lbl_eta": "Оставшееся время:",
    "lbl_ratio": "Рейтинг (Ratio):",
    "sec_torrent_details": "Сведения о торренте",
    "lbl_total_size": "Общий размер:",
    "lbl_piece_size": "Размер части:",
    "lbl_piece_count": "Количество частей:",
    "lbl_info_hash": "Хэш (InfoHash):",
    "lbl_save_path": "Папка сохранения:",
    "lbl_created_date": "Дата создания:",
    "lbl_created_by": "Создано программой:",
    "col_file_name": "Имя файла",
    "col_file_size": "Размер",
    "col_file_prog": "Прогресс",
    "col_file_prio": "Приоритет",
    "col_peer_ip": "IP-адрес / Порт",
    "col_peer_client": "Клиент",
    "col_peer_dl": "Скорость приема",
    "col_peer_ul": "Скорость отдачи",
    "col_peer_prog": "Прогресс пира",
    "col_tracker_url": "Адрес трекера",
    "col_tracker_status": "Статус",
    "legend_done": "Скачано",
    "legend_part": "Частично",
    "legend_wait": "Ожидает",
    "status_ready": "Готов",
    "ctx_resume": "Возобновить",
    "ctx_pause": "Приостановить",
    "ctx_open_folder": "Открыть папку",
    "ctx_sequential": "Последовательно",
    "ctx_verify": "Перепроверить",
    "ctx_redownload": "Скачать заново",
    "ctx_copy_magnet": "Скопировать Magnet-ссылку",
    "ctx_copy_hash": "Скопировать хэш",
    "ctx_delete": "Удалить...",
    "ctx_file_open": "Открыть файл",
    "ctx_file_show_folder": "Показать в проводнике",
    "ctx_prio_header": "Приоритет",
    "ctx_prio_high": "Высокий приоритет",
    "ctx_prio_normal": "Обычный приоритет",
    "ctx_prio_none": "Не скачивать",
    "ctx_file_copy_name": "Скопировать имя",
    "ctx_file_copy_path": "Скопировать полный путь",
    "modal_magnet_title": "Добавить Magnet-ссылку",
    "modal_magnet_label": "Введите или вставьте Magnet URI:",
    "btn_cancel": "Отмена",
    "btn_submit_magnet": "Добавить в загрузки",
    "modal_delete_title": "Удаление торрента",
    "modal_delete_confirm": "Вы действительно хотите удалить",
    "modal_delete_files": "Также безвозвратно удалить скачанные файлы с диска",
    "modal_settings_title": "Настройки программы",
    "settings_section_appearance": "Внешний вид и язык",
    "setting_theme_label": "Тема оформления:",
    "theme_dark": "Тёмная",
    "theme_light": "Светлая",
    "setting_language_label": "Язык интерфейса:",
    "lbl_default_download_dir": "Папка загрузок по умолчанию:",
    "btn_browse": "Обзор...",
    "btn_reset_downloads": "В Загрузки",
    "btn_reset_downloads_title": "Вернуть стандартную папку 'Загрузки'",
    "hint_default_download_dir": "Все новые торренты будут сохраняться в эту папку",
    "lbl_listen_port": "Порт для входящих соединений (0 = авто):",
    "lbl_max_conns": "Макс. пиров на один торрент:",
    "lbl_dl_limit": "Ограничение загрузки (КБ/с, 0 = Без ограничений):",
    "lbl_ul_limit": "Ограничение отдачи (КБ/с, 0 = Без ограничений):",
    "sec_p2p": "Сетевые протоколы P2P",
    "lbl_enable_dht": "Включить DHT (Distributed Hash Table — поиск пиров без трекера)",
    "lbl_enable_pex": "Включить PEX (Peer Exchange — обмен списками участников)",
    "lbl_enable_upnp": "Включить UPnP / NAT-PMP (автоматическая переадресация портов роутера)",
    "sec_win_integration": "Интеграция с Windows",
    "lbl_assoc_torrent": "Ассоциировать со всеми .torrent файлами",
    "lbl_assoc_magnet": "Ассоциировать с magnet: ссылками в браузерах",
    "lbl_autostart": "Автоматически запускать SmirnovTorrent при входе в Windows",
    "sec_tray": "Системный трей и фоновый режим",
    "lbl_minimize_tray": "Сворачивать окно в системный трей при минимизации",
    "lbl_close_tray": "Кнопка закрытия [X] сворачивает программу в трей вместо выхода",
    "lbl_start_minimized": "Запускать программу сразу свёрнутой в трей",
    "btn_save_settings": "Сохранить настройки",
    "toast_settings_saved": "Настройки успешно сохранены",
    "toast_settings_error": "Ошибка сохранения: ",
    "toast_load_error": "Не удалось загрузить настройки: ",
    "prio_high": "Высокий",
    "prio_normal": "Обычный",
    "prio_skip": "Пропущен",
    "status_downloading": "Загрузка",
    "status_seeding": "Раздача",
    "status_paused": "На паузе",
    "status_checking": "Проверка",
    "status_missing": "Отсутствует",
    "status_metadata": "Получение инфо",
    "status_completed": "Завершено",
    "status_unknown": "Неизвестно",
    "done": "Готово",
    "day_short": "д",
    "hour_short": "ч",
    "min_short": "м",
    "sec_short": "с",
    "total_dl_prefix": "Всего скачано: ",
    "total_ul_prefix": "Всего отдано: ",
    "torrents_count_prefix": "Торрентов: ",
    "dht_nodes_suffix": "узлов",
    "pieces_suffix": "частей",
    "byte_b": "Б",
    "byte_kb": "КБ",
    "byte_mb": "МБ",
    "byte_gb": "ГБ",
    "byte_tb": "ТБ",
    "speed_unit": "/с",
    "toast_torrent_added": "Торрент добавлен: ",
    "toast_torrent_add_error": "Ошибка добавления торрента: ",
    "toast_error": "Ошибка: ",
    "toast_upload_error": "Ошибка загрузки: ",
    "toast_upload_fail": "Не удалось загрузить файл",
    "toast_magnet_copied": "Magnet-ссылка скопирована в буфер обмена",
    "toast_hash_copied": "Хэш торрента скопирован в буфер обмена",
    "toast_file_name_copied": "Имя файла скопировано",
    "toast_file_path_copied": "Путь к файлу скопирован",
    "toast_drop_torrent_only": "Пожалуйста, перетащите файл с расширением .torrent",
    "toast_drop_valid": "Пожалуйста, перетащите .torrent файл или Magnet-ссылку",
    "alert_invalid_magnet": "Пожалуйста, введите корректную Magnet-ссылку (начинается с magnet:?)",
    "toast_magnet_added": "Magnet-ссылка успешно добавлена!",
    "toast_magnet_add_error": "Ошибка добавления Magnet: ",
    "toast_action_fail": "Не удалось выполнить действие",
    "toast_seq_enabled": "Последовательная загрузка включена",
    "toast_seq_disabled": "Последовательная загрузка выключена",
    "toast_verify_started": "Запущена проверка целостности данных",
    "toast_redownload_started": "Торрент перезапущен для загрузки заново",
    "toast_open_folder_error": "Не удалось открыть папку: ",
    "toast_folder_not_found": "Папка не найдена",
    "toast_open_file_error": "Не удалось открыть файл: ",
    "toast_file_not_found": "Файл не найден на диске",
    "toast_explorer_error": "Ошибка проводника",
    "toast_prio_updated": "Приоритет файла обновлен",
    "toast_torrent_deleted": "Торрент удален",
    "toast_delete_error": "Ошибка удаления: ",
    "toast_browse_dir_error": "Ошибка выбора папки: ",
    "not_specified": "Не указано",
    "working": "Работает",
    "unknown": "Неизвестно",
    "untitled": "Без названия",
    "seq_on": " (Вкл)",
    "shortcut_space": "Пробел",
    "dht_active": "Активен"
  },
  "en": {
    "drag_title": "Add Torrent",
    "drag_desc": "Drop .torrent file or Magnet link here to add",
    "btn_resume": "Resume",
    "btn_resume_title": "Resume selected torrent",
    "btn_pause": "Pause",
    "btn_pause_title": "Pause selected torrent",
    "btn_delete": "Delete",
    "btn_delete_title": "Delete torrent (Delete)",
    "btn_folder": "Folder",
    "btn_folder_title": "Open download folder in Explorer",
    "btn_sequential": "Sequential",
    "btn_sequential_title": "Sequential download (for streaming playback)",
    "btn_verify": "Verify",
    "btn_verify_title": "Recheck data integrity",
    "btn_redownload": "Redownload",
    "btn_redownload_title": "Redownload / Fix missing files",
    "btn_settings": "Settings",
    "btn_settings_title": "Application Settings",
    "search_placeholder": "Search torrents...",
    "cat_title": "Categories",
    "cat_all": "All Torrents",
    "cat_downloading": "Downloading",
    "cat_seeding": "Seeding",
    "cat_completed": "Completed",
    "cat_paused": "Paused",
    "cat_missing": "Missing",
    "engine_active": "Client Active",
    "col_name": "Name",
    "col_size": "Size",
    "col_progress": "Progress",
    "col_status": "Status",
    "col_dlspeed": "Down Speed",
    "col_ulspeed": "Up Speed",
    "col_eta": "ETA",
    "col_peers": "Peers",
    "col_seeds": "Seeds",
    "col_ratio": "Ratio",
    "col_added": "Added",
    "empty_title": "No torrents yet",
    "empty_desc": "Drag and drop a .torrent file from Explorer or use the buttons below",
    "empty_select_file": "Select .torrent file",
    "empty_paste_magnet": "Paste Magnet link",
    "tab_general": "General",
    "tab_files": "Files",
    "tab_peers": "Peers",
    "tab_trackers": "Trackers",
    "tab_pieces": "Piece Map",
    "sec_transfer": "Transfer",
    "lbl_downloaded": "Downloaded:",
    "lbl_uploaded": "Uploaded:",
    "lbl_dl_speed": "Download Speed:",
    "lbl_ul_speed": "Upload Speed:",
    "lbl_eta": "ETA:",
    "lbl_ratio": "Ratio:",
    "sec_torrent_details": "Torrent Information",
    "lbl_total_size": "Total Size:",
    "lbl_piece_size": "Piece Size:",
    "lbl_piece_count": "Piece Count:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Save Path:",
    "lbl_created_date": "Created On:",
    "lbl_created_by": "Created By:",
    "col_file_name": "File Name",
    "col_file_size": "Size",
    "col_file_prog": "Progress",
    "col_file_prio": "Priority",
    "col_peer_ip": "IP Address / Port",
    "col_peer_client": "Client",
    "col_peer_dl": "Down Speed",
    "col_peer_ul": "Up Speed",
    "col_peer_prog": "Peer Progress",
    "col_tracker_url": "Tracker Address",
    "col_tracker_status": "Status",
    "legend_done": "Done",
    "legend_part": "Partial",
    "legend_wait": "Pending",
    "status_ready": "Ready",
    "ctx_resume": "Resume",
    "ctx_pause": "Pause",
    "ctx_open_folder": "Open Folder",
    "ctx_sequential": "Sequential",
    "ctx_verify": "Verify",
    "ctx_redownload": "Redownload",
    "ctx_copy_magnet": "Copy Magnet Link",
    "ctx_copy_hash": "Copy Hash",
    "ctx_delete": "Delete...",
    "ctx_file_open": "Open File",
    "ctx_file_show_folder": "Show in Folder",
    "ctx_prio_header": "Priority",
    "ctx_prio_high": "High Priority",
    "ctx_prio_normal": "Normal Priority",
    "ctx_prio_none": "Do Not Download",
    "ctx_file_copy_name": "Copy Name",
    "ctx_file_copy_path": "Copy Full Path",
    "modal_magnet_title": "Add Magnet Link",
    "modal_magnet_label": "Enter or paste Magnet URI:",
    "btn_cancel": "Cancel",
    "btn_submit_magnet": "Add to Downloads",
    "modal_delete_title": "Delete Torrent",
    "modal_delete_confirm": "Are you sure you want to delete",
    "modal_delete_files": "Also permanently delete downloaded files from disk",
    "modal_settings_title": "Application Settings",
    "settings_section_appearance": "Appearance & Language",
    "setting_theme_label": "Theme:",
    "theme_dark": "Dark",
    "theme_light": "Light",
    "setting_language_label": "Interface Language:",
    "lbl_default_download_dir": "Default Download Folder:",
    "btn_browse": "Browse...",
    "btn_reset_downloads": "To Downloads",
    "btn_reset_downloads_title": "Reset to default 'Downloads' folder",
    "hint_default_download_dir": "All new torrents will be saved to this folder",
    "lbl_listen_port": "Incoming Port (0 = auto):",
    "lbl_max_conns": "Max Peers per Torrent:",
    "lbl_dl_limit": "Download Limit (KB/s, 0 = unlimited):",
    "lbl_ul_limit": "Upload Limit (KB/s, 0 = unlimited):",
    "sec_p2p": "P2P Network Protocols",
    "lbl_enable_dht": "Enable DHT (Distributed Hash Table)",
    "lbl_enable_pex": "Enable PEX (Peer Exchange)",
    "lbl_enable_upnp": "Enable UPnP / NAT-PMP (Router port forwarding)",
    "sec_win_integration": "Windows Integration",
    "lbl_assoc_torrent": "Associate with .torrent files",
    "lbl_assoc_magnet": "Associate with magnet: links",
    "lbl_autostart": "Start SmirnovTorrent automatically on Windows login",
    "sec_tray": "System Tray & Background",
    "lbl_minimize_tray": "Minimize window to system tray",
    "lbl_close_tray": "Close button [X] minimizes to tray instead of quitting",
    "lbl_start_minimized": "Start application minimized in tray",
    "btn_save_settings": "Save Settings",
    "toast_settings_saved": "Settings saved successfully",
    "toast_settings_error": "Error saving settings: ",
    "toast_load_error": "Failed to load settings: ",
    "prio_high": "High",
    "prio_normal": "Normal",
    "prio_skip": "Skipped",
    "status_downloading": "Downloading",
    "status_seeding": "Seeding",
    "status_paused": "Paused",
    "status_checking": "Checking",
    "status_missing": "Missing",
    "status_metadata": "Fetching Info",
    "status_completed": "Completed",
    "status_unknown": "Unknown",
    "done": "Done",
    "day_short": "d",
    "hour_short": "h",
    "min_short": "m",
    "sec_short": "s",
    "total_dl_prefix": "Total Downloaded: ",
    "total_ul_prefix": "Total Uploaded: ",
    "torrents_count_prefix": "Torrents: ",
    "dht_nodes_suffix": "nodes",
    "pieces_suffix": "pieces",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent added: ",
    "toast_torrent_add_error": "Error adding torrent: ",
    "toast_error": "Error: ",
    "toast_upload_error": "Upload error: ",
    "toast_upload_fail": "Failed to upload file",
    "toast_magnet_copied": "Magnet link copied to clipboard",
    "toast_hash_copied": "Torrent hash copied to clipboard",
    "toast_file_name_copied": "File name copied",
    "toast_file_path_copied": "File path copied",
    "toast_drop_torrent_only": "Please drop a .torrent file",
    "toast_drop_valid": "Please drop a .torrent file or Magnet link",
    "alert_invalid_magnet": "Please enter a valid Magnet link (starts with magnet:?)",
    "toast_magnet_added": "Magnet link added successfully!",
    "toast_magnet_add_error": "Error adding Magnet: ",
    "toast_action_fail": "Failed to perform action",
    "toast_seq_enabled": "Sequential download enabled",
    "toast_seq_disabled": "Sequential download disabled",
    "toast_verify_started": "Data integrity verification started",
    "toast_redownload_started": "Torrent restarted for redownload",
    "toast_open_folder_error": "Failed to open folder: ",
    "toast_folder_not_found": "Folder not found",
    "toast_open_file_error": "Failed to open file: ",
    "toast_file_not_found": "File not found on disk",
    "toast_explorer_error": "Explorer error",
    "toast_prio_updated": "File priority updated",
    "toast_torrent_deleted": "Torrent deleted",
    "toast_delete_error": "Error deleting: ",
    "toast_browse_dir_error": "Error selecting folder: ",
    "not_specified": "Not specified",
    "working": "Working",
    "unknown": "Unknown",
    "untitled": "Untitled",
    "seq_on": " (On)",
    "shortcut_space": "Space",
    "dht_active": "Active"
  },
  "de": {
    "drag_title": "Torrent hinzufügen",
    "drag_desc": ".torrent-Datei oder Magnet-Link hier ablegen",
    "btn_resume": "Fortsetzen",
    "btn_resume_title": "Ausgewählten Torrent fortsetzen",
    "btn_pause": "Pause",
    "btn_pause_title": "Ausgewählten Torrent anhalten",
    "btn_delete": "Löschen",
    "btn_delete_title": "Torrent löschen (Entf)",
    "btn_folder": "Ordner",
    "btn_folder_title": "Download-Ordner im Explorer öffnen",
    "btn_sequential": "Sequenziell",
    "btn_sequential_title": "Sequenzieller Download (für Streaming)",
    "btn_verify": "Prüfen",
    "btn_verify_title": "Datenintegrität überprüfen",
    "btn_redownload": "Erneut laden",
    "btn_redownload_title": "Erneut laden / Fehlende Dateien reparieren",
    "btn_settings": "Einstellungen",
    "btn_settings_title": "Programmeinstellungen",
    "search_placeholder": "Torrents suchen...",
    "cat_title": "Kategorien",
    "cat_all": "Alle Torrents",
    "cat_downloading": "Herunterladen",
    "cat_seeding": "Verteilen",
    "cat_completed": "Abgeschlossen",
    "cat_paused": "Pausiert",
    "cat_missing": "Fehlend",
    "engine_active": "Client aktiv",
    "col_name": "Dateiname",
    "col_size": "Größe",
    "col_progress": "Fortschritt",
    "col_status": "Status",
    "col_dlspeed": "Download",
    "col_ulspeed": "Upload",
    "col_eta": "Restzeit",
    "col_peers": "Peers",
    "col_seeds": "Seeds",
    "col_ratio": "Verhältnis",
    "col_added": "Hinzugefügt",
    "empty_title": "Noch keine Torrents vorhanden",
    "empty_desc": "Ziehen Sie eine .torrent-Datei hierher oder nutzen Sie die Schaltflächen",
    "empty_select_file": ".torrent-Datei auswählen",
    "empty_paste_magnet": "Magnet-Link einfügen",
    "tab_general": "Allgemein",
    "tab_files": "Dateien",
    "tab_peers": "Peers",
    "tab_trackers": "Tracker",
    "tab_pieces": "Teilekarte",
    "sec_transfer": "Übertragung",
    "lbl_downloaded": "Geladen:",
    "lbl_uploaded": "Gesendet:",
    "lbl_dl_speed": "Download-Rate:",
    "lbl_ul_speed": "Upload-Rate:",
    "lbl_eta": "Restzeit:",
    "lbl_ratio": "Verhältnis:",
    "sec_torrent_details": "Torrent-Details",
    "lbl_total_size": "Gesamtgröße:",
    "lbl_piece_size": "Teilgröße:",
    "lbl_piece_count": "Teileanzahl:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Speicherpfad:",
    "lbl_created_date": "Erstelldatum:",
    "lbl_created_by": "Erstellt von:",
    "col_file_name": "Dateiname",
    "col_file_size": "Größe",
    "col_file_prog": "Fortschritt",
    "col_file_prio": "Priorität",
    "col_peer_ip": "IP-Adresse / Port",
    "col_peer_client": "Client",
    "col_peer_dl": "Download",
    "col_peer_ul": "Upload",
    "col_peer_prog": "Peer-Fortschritt",
    "col_tracker_url": "Tracker-Adresse",
    "col_tracker_status": "Status",
    "legend_done": "Fertig",
    "legend_part": "Teilweise",
    "legend_wait": "Wartend",
    "status_ready": "Bereit",
    "ctx_resume": "Fortsetzen",
    "ctx_pause": "Pause",
    "ctx_open_folder": "Ordner öffnen",
    "ctx_sequential": "Sequenziell",
    "ctx_verify": "Prüfen",
    "ctx_redownload": "Erneut laden",
    "ctx_copy_magnet": "Magnet-Link kopieren",
    "ctx_copy_hash": "Hash kopieren",
    "ctx_delete": "Löschen...",
    "ctx_file_open": "Datei öffnen",
    "ctx_file_show_folder": "Im Ordner anzeigen",
    "ctx_prio_header": "Priorität",
    "ctx_prio_high": "Hohe Priorität",
    "ctx_prio_normal": "Normale Priorität",
    "ctx_prio_none": "Nicht herunterladen",
    "ctx_file_copy_name": "Dateiname kopieren",
    "ctx_file_copy_path": "Vollständigen Pfad kopieren",
    "modal_magnet_title": "Magnet-Link hinzufügen",
    "modal_magnet_label": "Magnet-URI eingeben oder einfügen:",
    "btn_cancel": "Abbrechen",
    "btn_submit_magnet": "Zu Downloads hinzufügen",
    "modal_delete_title": "Torrent löschen",
    "modal_delete_confirm": "Möchten Sie diesen Torrent wirklich löschen",
    "modal_delete_files": "Heruntergeladene Dateien ebenfalls von der Festplatte löschen",
    "modal_settings_title": "Programmeinstellungen",
    "settings_section_appearance": "Erscheinungsbild & Sprache",
    "setting_theme_label": "Design / Thema:",
    "theme_dark": "Dunkel",
    "theme_light": "Hell",
    "setting_language_label": "Sprache:",
    "lbl_default_download_dir": "Standard-Downloadordner:",
    "btn_browse": "Durchsuchen...",
    "btn_reset_downloads": "Zu Downloads",
    "btn_reset_downloads_title": "Auf Standardordner zurücksetzen",
    "hint_default_download_dir": "Alle neuen Torrents werden hier gespeichert",
    "lbl_listen_port": "Eingehender Port (0 = auto):",
    "lbl_max_conns": "Max. Peers pro Torrent:",
    "lbl_dl_limit": "Download-Limit (KB/s, 0 = unbegrenzt):",
    "lbl_ul_limit": "Upload-Limit (KB/s, 0 = unbegrenzt):",
    "sec_p2p": "P2P-Protokolle",
    "lbl_enable_dht": "DHT aktivieren (Distributed Hash Table)",
    "lbl_enable_pex": "PEX aktivieren (Peer Exchange)",
    "lbl_enable_upnp": "UPnP / NAT-PMP aktivieren (Portweiterleitung)",
    "sec_win_integration": "Windows-Integration",
    "lbl_assoc_torrent": "Mit .torrent-Dateien verknüpfen",
    "lbl_assoc_magnet": "Mit magnet:-Links verknüpfen",
    "lbl_autostart": "Mit Windows automatisch starten",
    "sec_tray": "System-Tray & Hintergrund",
    "lbl_minimize_tray": "In System-Tray minimieren",
    "lbl_close_tray": "Schließen-Schaltfläche [X] minimiert in den Tray",
    "lbl_start_minimized": "Minimiert im Tray starten",
    "btn_save_settings": "Einstellungen speichern",
    "toast_settings_saved": "Einstellungen erfolgreich gespeichert",
    "toast_settings_error": "Fehler beim Speichern: ",
    "toast_load_error": "Einstellungen konnten nicht geladen werden: ",
    "prio_high": "Hoch",
    "prio_normal": "Normal",
    "prio_skip": "Übersprungen",
    "status_downloading": "Laden",
    "status_seeding": "Verteilen",
    "status_paused": "Pausiert",
    "status_checking": "Prüfen",
    "status_missing": "Fehlend",
    "status_metadata": "Info abrufen",
    "status_completed": "Fertig",
    "status_unknown": "Unbekannt",
    "done": "Fertig",
    "day_short": "T",
    "hour_short": "Std",
    "min_short": "Min",
    "sec_short": "Sek",
    "total_dl_prefix": "Gesamt geladen: ",
    "total_ul_prefix": "Gesamt gesendet: ",
    "torrents_count_prefix": "Torrents: ",
    "dht_nodes_suffix": "Knoten",
    "pieces_suffix": "Teile",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent hinzugefügt: ",
    "toast_torrent_add_error": "Fehler beim Hinzufügen des Torrents: ",
    "toast_error": "Fehler: ",
    "toast_upload_error": "Upload-Fehler: ",
    "toast_upload_fail": "Datei konnte nicht hochgeladen werden",
    "toast_magnet_copied": "Magnet-Link in die Zwischenablage kopiert",
    "toast_hash_copied": "Torrent-Hash in die Zwischenablage kopiert",
    "toast_file_name_copied": "Dateiname kopiert",
    "toast_file_path_copied": "Dateipfad kopiert",
    "toast_drop_torrent_only": "Bitte eine .torrent-Datei ablegen",
    "toast_drop_valid": "Bitte eine .torrent-Datei oder einen Magnet-Link ablegen",
    "alert_invalid_magnet": "Bitte einen gültigen Magnet-Link eingeben (beginnt mit magnet:?)",
    "toast_magnet_added": "Magnet-Link erfolgreich hinzugefügt!",
    "toast_magnet_add_error": "Fehler beim Hinzufügen von Magnet: ",
    "toast_action_fail": "Aktion konnte nicht ausgeführt werden",
    "toast_seq_enabled": "Sequentieller Download aktiviert",
    "toast_seq_disabled": "Sequentieller Download deaktiviert",
    "toast_verify_started": "Überprüfung der Datenintegrität gestartet",
    "toast_redownload_started": "Torrent zum erneuten Herunterladen neu gestartet",
    "toast_open_folder_error": "Ordner konnte nicht geöffnet werden: ",
    "toast_folder_not_found": "Ordner nicht gefunden",
    "toast_open_file_error": "Datei konnte nicht geöffnet werden: ",
    "toast_file_not_found": "Datei auf Datenträger nicht gefunden",
    "toast_explorer_error": "Explorer-Fehler",
    "toast_prio_updated": "Dateipriorität aktualisiert",
    "toast_torrent_deleted": "Torrent gelöscht",
    "toast_delete_error": "Fehler beim Löschen: ",
    "toast_browse_dir_error": "Fehler bei der Ordnerauswahl: ",
    "not_specified": "Nicht angegeben",
    "working": "Aktiv",
    "unknown": "Unbekannt",
    "untitled": "Ohne Titel",
    "seq_on": " (Ein)",
    "shortcut_space": "Leertaste",
    "dht_active": "Aktiv"
  },
  "fr": {
    "drag_title": "Ajouter un torrent",
    "drag_desc": "Déposez un fichier .torrent ou un lien Magnet ici",
    "btn_resume": "Reprendre",
    "btn_resume_title": "Reprendre le torrent sélectionné",
    "btn_pause": "Pause",
    "btn_pause_title": "Mettre en pause le torrent sélectionné",
    "btn_delete": "Supprimer",
    "btn_delete_title": "Supprimer le torrent (Suppr)",
    "btn_folder": "Dossier",
    "btn_folder_title": "Ouvrir le dossier dans l'Explorateur",
    "btn_sequential": "Séquentiel",
    "btn_sequential_title": "Téléchargement séquentiel (pour streaming)",
    "btn_verify": "Vérifier",
    "btn_verify_title": "Revérifier l'intégrité des données",
    "btn_redownload": "Retélécharger",
    "btn_redownload_title": "Retélécharger / Corriger fichiers manquants",
    "btn_settings": "Paramètres",
    "btn_settings_title": "Paramètres de l'application",
    "search_placeholder": "Rechercher des torrents...",
    "cat_title": "Catégories",
    "cat_all": "Tous les torrents",
    "cat_downloading": "Téléchargement",
    "cat_seeding": "Partage",
    "cat_completed": "Terminés",
    "cat_paused": "En pause",
    "cat_missing": "Manquants",
    "engine_active": "Client actif",
    "col_name": "Nom",
    "col_size": "Taille",
    "col_progress": "Progression",
    "col_status": "État",
    "col_dlspeed": "Réception",
    "col_ulspeed": "Envoi",
    "col_eta": "Temps restant",
    "col_peers": "Pairs",
    "col_seeds": "Sources",
    "col_ratio": "Ratio",
    "col_added": "Ajouté",
    "empty_title": "Aucun torrent",
    "empty_desc": "Glissez-déposez un fichier .torrent ou utilisez les boutons ci-dessous",
    "empty_select_file": "Choisir fichier .torrent",
    "empty_paste_magnet": "Coller un lien Magnet",
    "tab_general": "Général",
    "tab_files": "Fichiers",
    "tab_peers": "Pairs",
    "tab_trackers": "Trackers",
    "tab_pieces": "Pièces",
    "sec_transfer": "Transfert",
    "lbl_downloaded": "Téléchargé :",
    "lbl_uploaded": "Envoyé :",
    "lbl_dl_speed": "Vitesse de réception :",
    "lbl_ul_speed": "Vitesse d'envoi :",
    "lbl_eta": "Temps restant :",
    "lbl_ratio": "Ratio :",
    "sec_torrent_details": "Informations du torrent",
    "lbl_total_size": "Taille totale :",
    "lbl_piece_size": "Taille de pièce :",
    "lbl_piece_count": "Nombre de pièces :",
    "lbl_info_hash": "InfoHash :",
    "lbl_save_path": "Dossier de destination :",
    "lbl_created_date": "Date de création :",
    "lbl_created_by": "Créé par :",
    "col_file_name": "Nom du fichier",
    "col_file_size": "Taille",
    "col_file_prog": "Progression",
    "col_file_prio": "Priorité",
    "col_peer_ip": "Adresse IP / Port",
    "col_peer_client": "Client",
    "col_peer_dl": "Réception",
    "col_peer_ul": "Envoi",
    "col_peer_prog": "Progression",
    "col_tracker_url": "Adresse du tracker",
    "col_tracker_status": "État",
    "legend_done": "Reçu",
    "legend_part": "Partiel",
    "legend_wait": "En attente",
    "status_ready": "Prêt",
    "ctx_resume": "Reprendre",
    "ctx_pause": "Pause",
    "ctx_open_folder": "Ouvrir le dossier",
    "ctx_sequential": "Séquentiel",
    "ctx_verify": "Vérifier",
    "ctx_redownload": "Retélécharger",
    "ctx_copy_magnet": "Copier lien Magnet",
    "ctx_copy_hash": "Copier hash",
    "ctx_delete": "Supprimer...",
    "ctx_file_open": "Ouvrir le fichier",
    "ctx_file_show_folder": "Afficher dans le dossier",
    "ctx_prio_header": "Priorité",
    "ctx_prio_high": "Priorité haute",
    "ctx_prio_normal": "Priorité normale",
    "ctx_prio_none": "Ne pas télécharger",
    "ctx_file_copy_name": "Copier le nom",
    "ctx_file_copy_path": "Copier le chemin complet",
    "modal_magnet_title": "Ajouter un lien Magnet",
    "modal_magnet_label": "Entrez ou collez l'URI Magnet :",
    "btn_cancel": "Annuler",
    "btn_submit_magnet": "Ajouter aux téléchargements",
    "modal_delete_title": "Supprimer le torrent",
    "modal_delete_confirm": "Voulez-vous vraiment supprimer",
    "modal_delete_files": "Supprimer également les fichiers du disque dur",
    "modal_settings_title": "Paramètres",
    "settings_section_appearance": "Apparence & Langue",
    "setting_theme_label": "Thème :",
    "theme_dark": "Sombre",
    "theme_light": "Clair",
    "setting_language_label": "Langue :",
    "lbl_default_download_dir": "Dossier de téléchargement par défaut :",
    "btn_browse": "Parcourir...",
    "btn_reset_downloads": "Vers Téléchargements",
    "btn_reset_downloads_title": "Restaurer le dossier standard 'Téléchargements'",
    "hint_default_download_dir": "Tous les nouveaux torrents seront sauvegardés ici",
    "lbl_listen_port": "Port d'écoute (0 = auto) :",
    "lbl_max_conns": "Max pairs par torrent :",
    "lbl_dl_limit": "Limite de réception (Ko/s, 0 = illimité) :",
    "lbl_ul_limit": "Limite d'envoi (Ko/s, 0 = illimité) :",
    "sec_p2p": "Protocoles réseau P2P",
    "lbl_enable_dht": "Activer DHT (Distributed Hash Table)",
    "lbl_enable_pex": "Activer PEX (Peer Exchange)",
    "lbl_enable_upnp": "Activer UPnP / NAT-PMP (redirection de ports)",
    "sec_win_integration": "Intégration Windows",
    "lbl_assoc_torrent": "Associer aux fichiers .torrent",
    "lbl_assoc_magnet": "Associer aux liens magnet:",
    "lbl_autostart": "Lancer SmirnovTorrent au démarrage de Windows",
    "sec_tray": "Zone de notification et arrière-plan",
    "lbl_minimize_tray": "Réduire dans la zone de notification",
    "lbl_close_tray": "Fermer [X] minimise au lieu de quitter",
    "lbl_start_minimized": "Démarrer directement minimisé dans la zone de notification",
    "btn_save_settings": "Enregistrer les paramètres",
    "toast_settings_saved": "Paramètres enregistrés avec succès",
    "toast_settings_error": "Erreur d'enregistrement : ",
    "toast_load_error": "Impossible de charger les paramètres : ",
    "prio_high": "Haute",
    "prio_normal": "Normale",
    "prio_skip": "Ignoré",
    "status_downloading": "Réception",
    "status_seeding": "Partage",
    "status_paused": "En pause",
    "status_checking": "Vérification",
    "status_missing": "Manquant",
    "status_metadata": "Récupération info",
    "status_completed": "Terminé",
    "status_unknown": "Inconnu",
    "done": "Terminé",
    "day_short": "j",
    "hour_short": "h",
    "min_short": "m",
    "sec_short": "s",
    "total_dl_prefix": "Total téléchargé : ",
    "total_ul_prefix": "Total envoyé : ",
    "torrents_count_prefix": "Torrents : ",
    "dht_nodes_suffix": "nœuds",
    "pieces_suffix": "pièces",
    "byte_b": "B",
    "byte_kb": "Ko",
    "byte_mb": "Mo",
    "byte_gb": "Go",
    "byte_tb": "To",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent ajouté : ",
    "toast_torrent_add_error": "Erreur lors de l'ajout du torrent : ",
    "toast_error": "Erreur : ",
    "toast_upload_error": "Erreur de téléversement : ",
    "toast_upload_fail": "Échec du téléversement du fichier",
    "toast_magnet_copied": "Lien Magnet copié dans le presse-papiers",
    "toast_hash_copied": "Hash du torrent copié dans le presse-papiers",
    "toast_file_name_copied": "Nom de fichier copié",
    "toast_file_path_copied": "Chemin du fichier copié",
    "toast_drop_torrent_only": "Veuillez déposer un fichier .torrent",
    "toast_drop_valid": "Veuillez déposer un fichier .torrent ou un lien Magnet",
    "alert_invalid_magnet": "Veuillez saisir un lien Magnet valide (commence par magnet:?)",
    "toast_magnet_added": "Lien Magnet ajouté avec succès !",
    "toast_magnet_add_error": "Erreur lors de l'ajout du Magnet : ",
    "toast_action_fail": "Échec de l'exécution de l'action",
    "toast_seq_enabled": "Téléchargement séquentiel activé",
    "toast_seq_disabled": "Téléchargement séquentiel désactivé",
    "toast_verify_started": "Vérification de l'intégrité des données lancée",
    "toast_redownload_started": "Torrent redémarré pour retéléchargement",
    "toast_open_folder_error": "Impossible d'ouvrir le dossier : ",
    "toast_folder_not_found": "Dossier introuvable",
    "toast_open_file_error": "Impossible d'ouvrir le fichier : ",
    "toast_file_not_found": "Fichier introuvable sur le disque",
    "toast_explorer_error": "Erreur de l'explorateur",
    "toast_prio_updated": "Priorité du fichier mise à jour",
    "toast_torrent_deleted": "Torrent supprimé",
    "toast_delete_error": "Erreur de suppression : ",
    "toast_browse_dir_error": "Erreur de sélection de dossier : ",
    "not_specified": "Non spécifié",
    "working": "Actif",
    "unknown": "Inconnu",
    "untitled": "Sans titre",
    "seq_on": " (Activé)",
    "shortcut_space": "Espace",
    "dht_active": "Actif"
  },
  "es": {
    "drag_title": "Añadir torrent",
    "drag_desc": "Suelte el archivo .torrent o enlace Magnet aquí",
    "btn_resume": "Continuar",
    "btn_resume_title": "Reanudar torrent seleccionado",
    "btn_pause": "Pausar",
    "btn_pause_title": "Pausar torrent seleccionado",
    "btn_delete": "Eliminar",
    "btn_delete_title": "Eliminar torrent (Supr)",
    "btn_folder": "Carpeta",
    "btn_folder_title": "Abrir carpeta de descarga en el Explorador",
    "btn_sequential": "Secuencial",
    "btn_sequential_title": "Descarga secuencial (para streaming)",
    "btn_verify": "Verificar",
    "btn_verify_title": "Comprobar integridad de los datos",
    "btn_redownload": "Redescargar",
    "btn_redownload_title": "Descargar de nuevo / Reparar archivos",
    "btn_settings": "Ajustes",
    "btn_settings_title": "Ajustes del programa",
    "search_placeholder": "Buscar torrents...",
    "cat_title": "Categorías",
    "cat_all": "Todos",
    "cat_downloading": "Descargando",
    "cat_seeding": "Sembrando",
    "cat_completed": "Completados",
    "cat_paused": "En pausa",
    "cat_missing": "Faltantes",
    "engine_active": "Cliente activo",
    "col_name": "Nombre",
    "col_size": "Tamaño",
    "col_progress": "Progreso",
    "col_status": "Estado",
    "col_dlspeed": "Bajada",
    "col_ulspeed": "Subida",
    "col_eta": "Restante",
    "col_peers": "Pares",
    "col_seeds": "Semillas",
    "col_ratio": "Ratio",
    "col_added": "Añadido",
    "empty_title": "No hay torrents aún",
    "empty_desc": "Arrastre un archivo .torrent o use los botones de abajo",
    "empty_select_file": "Seleccionar archivo .torrent",
    "empty_paste_magnet": "Pegar enlace Magnet",
    "tab_general": "General",
    "tab_files": "Archivos",
    "tab_peers": "Pares",
    "tab_trackers": "Trackers",
    "tab_pieces": "Piezas",
    "sec_transfer": "Transferencia",
    "lbl_downloaded": "Descargado:",
    "lbl_uploaded": "Subido:",
    "lbl_dl_speed": "Velocidad de bajada:",
    "lbl_ul_speed": "Velocidad de subida:",
    "lbl_eta": "Tiempo restante:",
    "lbl_ratio": "Ratio:",
    "sec_torrent_details": "Detalles del torrent",
    "lbl_total_size": "Tamaño total:",
    "lbl_piece_size": "Tamaño de pieza:",
    "lbl_piece_count": "Cantidad de piezas:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Ruta de guardado:",
    "lbl_created_date": "Fecha de creación:",
    "lbl_created_by": "Creado por:",
    "col_file_name": "Nombre de archivo",
    "col_file_size": "Tamaño",
    "col_file_prog": "Progreso",
    "col_file_prio": "Prioridad",
    "col_peer_ip": "Dirección IP / Puerto",
    "col_peer_client": "Cliente",
    "col_peer_dl": "Bajada",
    "col_peer_ul": "Subida",
    "col_peer_prog": "Progreso",
    "col_tracker_url": "Dirección del Tracker",
    "col_tracker_status": "Estado",
    "legend_done": "Descargado",
    "legend_part": "Parcial",
    "legend_wait": "Pendiente",
    "status_ready": "Listo",
    "ctx_resume": "Continuar",
    "ctx_pause": "Pausar",
    "ctx_open_folder": "Abrir carpeta",
    "ctx_sequential": "Secuencial",
    "ctx_verify": "Verificar",
    "ctx_redownload": "Redescargar",
    "ctx_copy_magnet": "Copiar enlace Magnet",
    "ctx_copy_hash": "Copiar hash",
    "ctx_delete": "Eliminar...",
    "ctx_file_open": "Abrir archivo",
    "ctx_file_show_folder": "Mostrar en carpeta",
    "ctx_prio_header": "Prioridad",
    "ctx_prio_high": "Prioridad alta",
    "ctx_prio_normal": "Prioridad normal",
    "ctx_prio_none": "No descargar",
    "ctx_file_copy_name": "Copiar nombre",
    "ctx_file_copy_path": "Copiar ruta completa",
    "modal_magnet_title": "Añadir enlace Magnet",
    "modal_magnet_label": "Introduzca o pegue el enlace Magnet URI:",
    "btn_cancel": "Cancelar",
    "btn_submit_magnet": "Añadir a descargas",
    "modal_delete_title": "Eliminar torrent",
    "modal_delete_confirm": "¿Realmente desea eliminar",
    "modal_delete_files": "Eliminar también los archivos descargados del disco",
    "modal_settings_title": "Ajustes del programa",
    "settings_section_appearance": "Apariencia e Idioma",
    "setting_theme_label": "Tema:",
    "theme_dark": "Oscuro",
    "theme_light": "Claro",
    "setting_language_label": "Idioma de interfaz:",
    "lbl_default_download_dir": "Carpeta de descarga predeterminada:",
    "btn_browse": "Examinar...",
    "btn_reset_downloads": "A Descargas",
    "btn_reset_downloads_title": "Restaurar carpeta de Descargas predeterminada",
    "hint_default_download_dir": "Todos los nuevos torrents se guardarán aquí",
    "lbl_listen_port": "Puerto de escucha (0 = auto):",
    "lbl_max_conns": "Máx. pares por torrent:",
    "lbl_dl_limit": "Límite de descarga (KB/s, 0 = ilimitado):",
    "lbl_ul_limit": "Límite de subida (KB/s, 0 = ilimitado):",
    "sec_p2p": "Protocolos P2P",
    "lbl_enable_dht": "Activar DHT (Distributed Hash Table)",
    "lbl_enable_pex": "Activar PEX (Peer Exchange)",
    "lbl_enable_upnp": "Activar UPnP / NAT-PMP (redirección de puertos)",
    "sec_win_integration": "Integración con Windows",
    "lbl_assoc_torrent": "Asociar con archivos .torrent",
    "lbl_assoc_magnet": "Asociar con enlaces magnet:",
    "lbl_autostart": "Iniciar automáticamente con Windows",
    "sec_tray": "Bandeja del sistema",
    "lbl_minimize_tray": "Minimizar a la bandeja del sistema",
    "lbl_close_tray": "El botón cerrar [X] minimiza a la bandeja",
    "lbl_start_minimized": "Iniciar minimizado en la bandeja",
    "btn_save_settings": "Guardar ajustes",
    "toast_settings_saved": "Ajustes guardados correctamente",
    "toast_settings_error": "Error al guardar: ",
    "toast_load_error": "No se pudieron cargar los ajustes: ",
    "prio_high": "Alta",
    "prio_normal": "Normal",
    "prio_skip": "Omitido",
    "status_downloading": "Descargando",
    "status_seeding": "Sembrando",
    "status_paused": "En pausa",
    "status_checking": "Verificando",
    "status_missing": "Faltante",
    "status_metadata": "Obteniendo info",
    "status_completed": "Completado",
    "status_unknown": "Desconocido",
    "done": "Listo",
    "day_short": "d",
    "hour_short": "h",
    "min_short": "m",
    "sec_short": "s",
    "total_dl_prefix": "Total descargado: ",
    "total_ul_prefix": "Total subido: ",
    "torrents_count_prefix": "Torrents: ",
    "dht_nodes_suffix": "nodos",
    "pieces_suffix": "piezas",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent añadido: ",
    "toast_torrent_add_error": "Error al añadir el torrent: ",
    "toast_error": "Error: ",
    "toast_upload_error": "Error de subida: ",
    "toast_upload_fail": "No se pudo subir el archivo",
    "toast_magnet_copied": "Enlace Magnet copiado al portapapeles",
    "toast_hash_copied": "Hash del torrent copiado al portapapeles",
    "toast_file_name_copied": "Nombre de archivo copiado",
    "toast_file_path_copied": "Ruta del archivo copiada",
    "toast_drop_torrent_only": "Por favor, suelta un archivo .torrent",
    "toast_drop_valid": "Por favor, suelta un archivo .torrent o enlace Magnet",
    "alert_invalid_magnet": "Por favor, introduce un enlace Magnet válido (comienza con magnet:?)",
    "toast_magnet_added": "¡Enlace Magnet añadido con éxito!",
    "toast_magnet_add_error": "Error al añadir Magnet: ",
    "toast_action_fail": "No se pudo realizar la acción",
    "toast_seq_enabled": "Descarga secuencial activada",
    "toast_seq_disabled": "Descarga secuencial desactivada",
    "toast_verify_started": "Verificación de integridad iniciada",
    "toast_redownload_started": "Torrent reiniciado para volver a descargar",
    "toast_open_folder_error": "No se pudo abrir la carpeta: ",
    "toast_folder_not_found": "Carpeta no encontrada",
    "toast_open_file_error": "No se pudo abrir el archivo: ",
    "toast_file_not_found": "Archivo no encontrado en el disco",
    "toast_explorer_error": "Error del explorador",
    "toast_prio_updated": "Prioridad del archivo actualizada",
    "toast_torrent_deleted": "Torrent eliminado",
    "toast_delete_error": "Error al eliminar: ",
    "toast_browse_dir_error": "Error al seleccionar carpeta: ",
    "not_specified": "No especificado",
    "working": "Activo",
    "unknown": "Desconocido",
    "untitled": "Sin título",
    "seq_on": " (Activado)",
    "shortcut_space": "Espacio",
    "dht_active": "Activo"
  },
  "it": {
    "drag_title": "Aggiungi Torrent",
    "drag_desc": "Rilascia il file .torrent o il link Magnet qui",
    "btn_resume": "Riprendi",
    "btn_resume_title": "Riprendi torrent selezionato",
    "btn_pause": "Pausa",
    "btn_pause_title": "Metti in pausa torrent selezionato",
    "btn_delete": "Elimina",
    "btn_delete_title": "Elimina torrent (Canc)",
    "btn_folder": "Cartella",
    "btn_folder_title": "Apri cartella di download in Explorer",
    "btn_sequential": "Sequenziale",
    "btn_sequential_title": "Download sequenziale (per streaming)",
    "btn_verify": "Verifica",
    "btn_verify_title": "Verifica integrità dati",
    "btn_redownload": "Riscarica",
    "btn_redownload_title": "Riscarica / Ripara file mancanti",
    "btn_settings": "Impostazioni",
    "btn_settings_title": "Impostazioni applicazione",
    "search_placeholder": "Cerca torrent...",
    "cat_title": "Categorie",
    "cat_all": "Tutti i torrent",
    "cat_downloading": "In download",
    "cat_seeding": "In seeding",
    "cat_completed": "Completati",
    "cat_paused": "In pausa",
    "cat_missing": "Mancanti",
    "engine_active": "Client attivo",
    "col_name": "Nome",
    "col_size": "Dimensione",
    "col_progress": "Progresso",
    "col_status": "Stato",
    "col_dlspeed": "Download",
    "col_ulspeed": "Upload",
    "col_eta": "Tempo stimato",
    "col_peers": "Peer",
    "col_seeds": "Seed",
    "col_ratio": "Rapporto",
    "col_added": "Aggiunto",
    "empty_title": "Nessun torrent presente",
    "empty_desc": "Trascina un file .torrent o usa i pulsanti sottostanti",
    "empty_select_file": "Seleziona file .torrent",
    "empty_paste_magnet": "Incolla link Magnet",
    "tab_general": "Generale",
    "tab_files": "File",
    "tab_peers": "Peer",
    "tab_trackers": "Tracker",
    "tab_pieces": "Mappa parti",
    "sec_transfer": "Trasferimento",
    "lbl_downloaded": "Scaricato:",
    "lbl_uploaded": "Caricato:",
    "lbl_dl_speed": "Velocità download:",
    "lbl_ul_speed": "Velocità upload:",
    "lbl_eta": "Tempo residuo:",
    "lbl_ratio": "Rapporto:",
    "sec_torrent_details": "Dettagli torrent",
    "lbl_total_size": "Dimensione totale:",
    "lbl_piece_size": "Dimensione parte:",
    "lbl_piece_count": "Numero di parti:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Percorso di salvataggio:",
    "lbl_created_date": "Data di creazione:",
    "lbl_created_by": "Creato da:",
    "col_file_name": "Nome file",
    "col_file_size": "Dimensione",
    "col_file_prog": "Progresso",
    "col_file_prio": "Priorità",
    "col_peer_ip": "Indirizzo IP / Porta",
    "col_peer_client": "Client",
    "col_peer_dl": "Download",
    "col_peer_ul": "Upload",
    "col_peer_prog": "Progresso",
    "col_tracker_url": "Indirizzo tracker",
    "col_tracker_status": "Stato",
    "legend_done": "Completato",
    "legend_part": "Parziale",
    "legend_wait": "In attesa",
    "status_ready": "Pronto",
    "ctx_resume": "Riprendi",
    "ctx_pause": "Pausa",
    "ctx_open_folder": "Apri cartella",
    "ctx_sequential": "Sequenziale",
    "ctx_verify": "Verifica",
    "ctx_redownload": "Riscarica",
    "ctx_copy_magnet": "Copia link Magnet",
    "ctx_copy_hash": "Copia hash",
    "ctx_delete": "Elimina...",
    "ctx_file_open": "Apri file",
    "ctx_file_show_folder": "Mostra nella cartella",
    "ctx_prio_header": "Priorità",
    "ctx_prio_high": "Priorità alta",
    "ctx_prio_normal": "Priorità normale",
    "ctx_prio_none": "Non scaricare",
    "ctx_file_copy_name": "Copia nome",
    "ctx_file_copy_path": "Copia percorso completo",
    "modal_magnet_title": "Aggiungi link Magnet",
    "modal_magnet_label": "Inserisci o incolla l'URI Magnet:",
    "btn_cancel": "Annulla",
    "btn_submit_magnet": "Aggiungi ai download",
    "modal_delete_title": "Elimina torrent",
    "modal_delete_confirm": "Sei sicuro di voler eliminare",
    "modal_delete_files": "Elimina definitivamente anche i file dal disco",
    "modal_settings_title": "Impostazioni programma",
    "settings_section_appearance": "Aspetto e Lingua",
    "setting_theme_label": "Tema:",
    "theme_dark": "Scuro",
    "theme_light": "Chiaro",
    "setting_language_label": "Lingua interfaccia:",
    "lbl_default_download_dir": "Cartella di download predefinita:",
    "btn_browse": "Sfoglia...",
    "btn_reset_downloads": "In Download",
    "btn_reset_downloads_title": "Reimposta cartella Download predefinita",
    "hint_default_download_dir": "Tutti i nuovi torrent saranno salvati qui",
    "lbl_listen_port": "Porta in ingresso (0 = auto):",
    "lbl_max_conns": "Max peer per torrent:",
    "lbl_dl_limit": "Limite download (KB/s, 0 = illimitato):",
    "lbl_ul_limit": "Limite upload (KB/s, 0 = illimitato):",
    "sec_p2p": "Protocolli di rete P2P",
    "lbl_enable_dht": "Abilita DHT (Distributed Hash Table)",
    "lbl_enable_pex": "Abilita PEX (Peer Exchange)",
    "lbl_enable_upnp": "Abilita UPnP / NAT-PMP (inoltro porte)",
    "sec_win_integration": "Integrazione Windows",
    "lbl_assoc_torrent": "Associa a file .torrent",
    "lbl_assoc_magnet": "Associa a link magnet:",
    "lbl_autostart": "Avvia SmirnovTorrent all'accesso a Windows",
    "sec_tray": "Barra delle applicazioni e sfondo",
    "lbl_minimize_tray": "Riduci nella barra delle applicazioni",
    "lbl_close_tray": "Il tasto chiudi [X] riduce invece di uscire",
    "lbl_start_minimized": "Avvia ridotto nella barra",
    "btn_save_settings": "Salva impostazioni",
    "toast_settings_saved": "Impostazioni salvate con successo",
    "toast_settings_error": "Errore durante il salvataggio: ",
    "toast_load_error": "Impossibile caricare le impostazioni: ",
    "prio_high": "Alta",
    "prio_normal": "Normale",
    "prio_skip": "Saltato",
    "status_downloading": "Download",
    "status_seeding": "Seeding",
    "status_paused": "In pausa",
    "status_checking": "Verifica",
    "status_missing": "Mancante",
    "status_metadata": "Recupero info",
    "status_completed": "Completato",
    "status_unknown": "Sconosciuto",
    "done": "Fatto",
    "day_short": "g",
    "hour_short": "o",
    "min_short": "m",
    "sec_short": "s",
    "total_dl_prefix": "Totale scaricato: ",
    "total_ul_prefix": "Totale inviato: ",
    "torrents_count_prefix": "Torrent: ",
    "dht_nodes_suffix": "nodi",
    "pieces_suffix": "parti",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent aggiunto: ",
    "toast_torrent_add_error": "Errore durante l'aggiunta del torrent: ",
    "toast_error": "Errore: ",
    "toast_upload_error": "Errore di caricamento: ",
    "toast_upload_fail": "Impossibile caricare il file",
    "toast_magnet_copied": "Collegamento Magnet copiato negli appunti",
    "toast_hash_copied": "Hash del torrent copiato negli appunti",
    "toast_file_name_copied": "Nome del file copiato",
    "toast_file_path_copied": "Percorso del file copiato",
    "toast_drop_torrent_only": "Rilascia un file .torrent",
    "toast_drop_valid": "Rilascia un file .torrent o collegamento Magnet",
    "alert_invalid_magnet": "Inserisci un collegamento Magnet valido (inizia con magnet:?)",
    "toast_magnet_added": "Collegamento Magnet aggiunto con successo!",
    "toast_magnet_add_error": "Errore durante l'aggiunta di Magnet: ",
    "toast_action_fail": "Impossibile eseguire l'azione",
    "toast_seq_enabled": "Download sequenziale attivato",
    "toast_seq_disabled": "Download sequenziale disattivato",
    "toast_verify_started": "Verifica dell'integrità dei dati avviata",
    "toast_redownload_started": "Torrent riavviato per il download",
    "toast_open_folder_error": "Impossibile aprire la cartella: ",
    "toast_folder_not_found": "Cartella non trovata",
    "toast_open_file_error": "Impossibile aprire il file: ",
    "toast_file_not_found": "File non trovato sul disco",
    "toast_explorer_error": "Errore di Explorer",
    "toast_prio_updated": "Priorità del file aggiornata",
    "toast_torrent_deleted": "Torrent eliminato",
    "toast_delete_error": "Errore durante l'eliminazione: ",
    "toast_browse_dir_error": "Errore nella selezione della cartella: ",
    "not_specified": "Non specificato",
    "working": "Funzionante",
    "unknown": "Sconosciuto",
    "untitled": "Senza título",
    "seq_on": " (Attivo)",
    "shortcut_space": "Spazio",
    "dht_active": "Attivo"
  },
  "pt": {
    "drag_title": "Adicionar Torrent",
    "drag_desc": "Arraste o arquivo .torrent ou link Magnet para aqui",
    "btn_resume": "Continuar",
    "btn_resume_title": "Retomar torrent selecionado",
    "btn_pause": "Pausa",
    "btn_pause_title": "Pausar torrent selecionado",
    "btn_delete": "Excluir",
    "btn_delete_title": "Excluir torrent (Delete)",
    "btn_folder": "Pasta",
    "btn_folder_title": "Abrir pasta no Explorador de Arquivos",
    "btn_sequential": "Sequencial",
    "btn_sequential_title": "Download sequencial (para reprodução contínua)",
    "btn_verify": "Verificar",
    "btn_verify_title": "Verificar integridade dos dados",
    "btn_redownload": "Recarregar",
    "btn_redownload_title": "Baixar novamente / Corrigir arquivos em falta",
    "btn_settings": "Configurações",
    "btn_settings_title": "Configurações do aplicativo",
    "search_placeholder": "Pesquisar torrents...",
    "cat_title": "Categorias",
    "cat_all": "Todos",
    "cat_downloading": "Baixando",
    "cat_seeding": "Enviando",
    "cat_completed": "Concluídos",
    "cat_paused": "Em pausa",
    "cat_missing": "Em falta",
    "engine_active": "Cliente ativo",
    "col_name": "Nome",
    "col_size": "Tamanho",
    "col_progress": "Progresso",
    "col_status": "Status",
    "col_dlspeed": "Download",
    "col_ulspeed": "Upload",
    "col_eta": "Restante",
    "col_peers": "Pares",
    "col_seeds": "Seeds",
    "col_ratio": "Taxa",
    "col_added": "Adicionado",
    "empty_title": "Ainda não há torrents",
    "empty_desc": "Arraste um arquivo .torrent ou use os botões abaixo",
    "empty_select_file": "Selecionar arquivo .torrent",
    "empty_paste_magnet": "Colar link Magnet",
    "tab_general": "Geral",
    "tab_files": "Arquivos",
    "tab_peers": "Pares",
    "tab_trackers": "Rastreadores",
    "tab_pieces": "Mapa de partes",
    "sec_transfer": "Transferência",
    "lbl_downloaded": "Baixado:",
    "lbl_uploaded": "Enviado:",
    "lbl_dl_speed": "Velocidade de download:",
    "lbl_ul_speed": "Velocidade de upload:",
    "lbl_eta": "Tempo restante:",
    "lbl_ratio": "Taxa (Ratio):",
    "sec_torrent_details": "Detalhes do torrent",
    "lbl_total_size": "Tamanho total:",
    "lbl_piece_size": "Tamanho do bloco:",
    "lbl_piece_count": "Número de blocos:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Pasta de destino:",
    "lbl_created_date": "Data de criação:",
    "lbl_created_by": "Criado por:",
    "col_file_name": "Nome do arquivo",
    "col_file_size": "Tamanho",
    "col_file_prog": "Progresso",
    "col_file_prio": "Prioridade",
    "col_peer_ip": "Endereço IP / Porta",
    "col_peer_client": "Cliente",
    "col_peer_dl": "Download",
    "col_peer_ul": "Upload",
    "col_peer_prog": "Progresso do par",
    "col_tracker_url": "Endereço do rastreador",
    "col_tracker_status": "Status",
    "legend_done": "Baixado",
    "legend_part": "Parcial",
    "legend_wait": "Aguardando",
    "status_ready": "Pronto",
    "ctx_resume": "Continuar",
    "ctx_pause": "Pausa",
    "ctx_open_folder": "Abrir pasta",
    "ctx_sequential": "Sequencial",
    "ctx_verify": "Verificar",
    "ctx_redownload": "Baixar novamente",
    "ctx_copy_magnet": "Copiar link Magnet",
    "ctx_copy_hash": "Copiar hash",
    "ctx_delete": "Excluir...",
    "ctx_file_open": "Abrir arquivo",
    "ctx_file_show_folder": "Mostrar na pasta",
    "ctx_prio_header": "Prioridade",
    "ctx_prio_high": "Alta prioridade",
    "ctx_prio_normal": "Prioridade normal",
    "ctx_prio_none": "Não baixar",
    "ctx_file_copy_name": "Copiar nome",
    "ctx_file_copy_path": "Copiar caminho completo",
    "modal_magnet_title": "Adicionar link Magnet",
    "modal_magnet_label": "Insira ou cole a URI Magnet:",
    "btn_cancel": "Cancelar",
    "btn_submit_magnet": "Adicionar aos downloads",
    "modal_delete_title": "Excluir torrent",
    "modal_delete_confirm": "Tem certeza que deseja excluir",
    "modal_delete_files": "Também excluir permanentemente os arquivos baixados do disco",
    "modal_settings_title": "Configurações do aplicativo",
    "settings_section_appearance": "Aparência e Idioma",
    "setting_theme_label": "Tema:",
    "theme_dark": "Escuro",
    "theme_light": "Claro",
    "setting_language_label": "Idioma da interface:",
    "lbl_default_download_dir": "Pasta de download padrão:",
    "btn_browse": "Procurar...",
    "btn_reset_downloads": "Para Downloads",
    "btn_reset_downloads_title": "Restaurar pasta padrão de Downloads",
    "hint_default_download_dir": "Todos os novos torrents serão salvos nesta pasta",
    "lbl_listen_port": "Porta de entrada (0 = auto):",
    "lbl_max_conns": "Máx. pares por torrent:",
    "lbl_dl_limit": "Limite de download (KB/s, 0 = ilimitado):",
    "lbl_ul_limit": "Limite de upload (KB/s, 0 = ilimitado):",
    "sec_p2p": "Protocolos P2P",
    "lbl_enable_dht": "Ativar DHT (Distributed Hash Table)",
    "lbl_enable_pex": "Ativar PEX (Peer Exchange)",
    "lbl_enable_upnp": "Ativar UPnP / NAT-PMP (redirecionamento de portas)",
    "sec_win_integration": "Integração com o Windows",
    "lbl_assoc_torrent": "Associar com arquivos .torrent",
    "lbl_assoc_magnet": "Associar com links magnet:",
    "lbl_autostart": "Iniciar automaticamente com o Windows",
    "sec_tray": "Bandeja do sistema",
    "lbl_minimize_tray": "Minimizar para a bandeja do sistema",
    "lbl_close_tray": "O botão fechar [X] minimiza para a bandeja",
    "lbl_start_minimized": "Iniciar minimizado na bandeja",
    "btn_save_settings": "Salvar configurações",
    "toast_settings_saved": "Configurações salvas com sucesso",
    "toast_settings_error": "Erro ao salvar configurações: ",
    "toast_load_error": "Não foi possível carregar as configurações: ",
    "prio_high": "Alta",
    "prio_normal": "Normal",
    "prio_skip": "Ignorado",
    "status_downloading": "Baixando",
    "status_seeding": "Enviando",
    "status_paused": "Em pausa",
    "status_checking": "Verificando",
    "status_missing": "Em falta",
    "status_metadata": "Obtendo info",
    "status_completed": "Concluído",
    "status_unknown": "Desconhecido",
    "done": "Concluído",
    "day_short": "d",
    "hour_short": "h",
    "min_short": "m",
    "sec_short": "s",
    "total_dl_prefix": "Total baixado: ",
    "total_ul_prefix": "Total enviado: ",
    "torrents_count_prefix": "Torrents: ",
    "dht_nodes_suffix": "nós",
    "pieces_suffix": "partes",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent adicionado: ",
    "toast_torrent_add_error": "Erro ao adicionar torrent: ",
    "toast_error": "Erro: ",
    "toast_upload_error": "Erro de envio: ",
    "toast_upload_fail": "Falha ao enviar arquivo",
    "toast_magnet_copied": "Link Magnet copiado para a área de transferência",
    "toast_hash_copied": "Hash do torrent copiado para a área de transferência",
    "toast_file_name_copied": "Nome do arquivo copiado",
    "toast_file_path_copied": "Caminho do arquivo copiado",
    "toast_drop_torrent_only": "Por favor, solte um arquivo .torrent",
    "toast_drop_valid": "Por favor, solte um arquivo .torrent ou link Magnet",
    "alert_invalid_magnet": "Por favor, insira um link Magnet válido (começa com magnet:?)",
    "toast_magnet_added": "Link Magnet adicionado com sucesso!",
    "toast_magnet_add_error": "Erro ao adicionar Magnet: ",
    "toast_action_fail": "Falha ao executar a ação",
    "toast_seq_enabled": "Download sequencial ativado",
    "toast_seq_disabled": "Download sequencial desativado",
    "toast_verify_started": "Verificação de integridade iniciada",
    "toast_redownload_started": "Torrent reiniciado para novo download",
    "toast_open_folder_error": "Falha ao abrir a pasta: ",
    "toast_folder_not_found": "Pasta não encontrada",
    "toast_open_file_error": "Falha ao abrir o arquivo: ",
    "toast_file_not_found": "Arquivo não encontrado no disco",
    "toast_explorer_error": "Erro do Explorer",
    "toast_prio_updated": "Prioridade do arquivo atualizada",
    "toast_torrent_deleted": "Torrent excluído",
    "toast_delete_error": "Erro ao excluir: ",
    "toast_browse_dir_error": "Erro ao selecionar pasta: ",
    "not_specified": "Não especificado",
    "working": "Funcionando",
    "unknown": "Desconhecido",
    "untitled": "Sem título",
    "seq_on": " (Ligado)",
    "shortcut_space": "Espaço",
    "dht_active": "Ativo"
  },
  "zh": {
    "drag_title": "添加种子",
    "drag_desc": "将 .torrent 文件或磁力链接拖放至此处",
    "btn_resume": "继续",
    "btn_resume_title": "继续选中的种子",
    "btn_pause": "暂停",
    "btn_pause_title": "暂停选中的种子",
    "btn_delete": "删除",
    "btn_delete_title": "删除种子 (Delete)",
    "btn_folder": "文件夹",
    "btn_folder_title": "在资源管理器中打开下载文件夹",
    "btn_sequential": "顺序下载",
    "btn_sequential_title": "顺序下载（适用于流媒体边下边播）",
    "btn_verify": "校验",
    "btn_verify_title": "重新校验数据完整性",
    "btn_redownload": "重新下载",
    "btn_redownload_title": "重新下载 / 修复缺失文件",
    "btn_settings": "设置",
    "btn_settings_title": "程序设置",
    "search_placeholder": "搜索种子...",
    "cat_title": "分类",
    "cat_all": "全部种子",
    "cat_downloading": "正在下载",
    "cat_seeding": "正在做种",
    "cat_completed": "已完成",
    "cat_paused": "已暂停",
    "cat_missing": "文件缺失",
    "engine_active": "客户端运行中",
    "col_name": "名称",
    "col_size": "大小",
    "col_progress": "进度",
    "col_status": "状态",
    "col_dlspeed": "下载速度",
    "col_ulspeed": "上传速度",
    "col_eta": "剩余时间",
    "col_peers": "节点",
    "col_seeds": "做种数",
    "col_ratio": "分享率",
    "col_added": "添加时间",
    "empty_title": "当前暂无种子",
    "empty_desc": "从资源管理器拖入 .torrent 文件或点击下方按钮",
    "empty_select_file": "选择 .torrent 文件",
    "empty_paste_magnet": "粘贴磁力链接",
    "tab_general": "常规",
    "tab_files": "文件",
    "tab_peers": "节点",
    "tab_trackers": "Tracker",
    "tab_pieces": "区块地图",
    "sec_transfer": "传输统计",
    "lbl_downloaded": "已下载：",
    "lbl_uploaded": "已上传：",
    "lbl_dl_speed": "下载速度：",
    "lbl_ul_speed": "上传速度：",
    "lbl_eta": "剩余时间：",
    "lbl_ratio": "分享率：",
    "sec_torrent_details": "种子详情",
    "lbl_total_size": "总大小：",
    "lbl_piece_size": "分块大小：",
    "lbl_piece_count": "分块总数：",
    "lbl_info_hash": "特征码 (InfoHash)：",
    "lbl_save_path": "保存路径：",
    "lbl_created_date": "创建时间：",
    "lbl_created_by": "创建工具：",
    "col_file_name": "文件名",
    "col_file_size": "大小",
    "col_file_prog": "进度",
    "col_file_prio": "优先级",
    "col_peer_ip": "IP 地址 / 端口",
    "col_peer_client": "客户端",
    "col_peer_dl": "下载速率",
    "col_peer_ul": "上传速率",
    "col_peer_prog": "节点进度",
    "col_tracker_url": "Tracker 地址",
    "col_tracker_status": "状态",
    "legend_done": "已下载",
    "legend_part": "部分完成",
    "legend_wait": "等待中",
    "status_ready": "就绪",
    "ctx_resume": "继续",
    "ctx_pause": "暂停",
    "ctx_open_folder": "打开文件夹",
    "ctx_sequential": "顺序下载",
    "ctx_verify": "校验数据",
    "ctx_redownload": "重新下载",
    "ctx_copy_magnet": "复制磁力链接",
    "ctx_copy_hash": "复制特征码",
    "ctx_delete": "删除...",
    "ctx_file_open": "打开文件",
    "ctx_file_show_folder": "在文件夹中显示",
    "ctx_prio_header": "优先级",
    "ctx_prio_high": "高优先级",
    "ctx_prio_normal": "普通优先级",
    "ctx_prio_none": "不下载",
    "ctx_file_copy_name": "复制名称",
    "ctx_file_copy_path": "复制完整路径",
    "modal_magnet_title": "添加磁力链接",
    "modal_magnet_label": "输入或粘贴磁力链接 URI：",
    "btn_cancel": "取消",
    "btn_submit_magnet": "添加到下载列表",
    "modal_delete_title": "删除种子",
    "modal_delete_confirm": "确定要删除",
    "modal_delete_files": "同时从磁盘彻底删除已下载的文件",
    "modal_settings_title": "程序设置",
    "settings_section_appearance": "界面外观与语言",
    "setting_theme_label": "外观主题：",
    "theme_dark": "深色",
    "theme_light": "浅色",
    "setting_language_label": "界面语言：",
    "lbl_default_download_dir": "默认下载保存文件夹：",
    "btn_browse": "浏览...",
    "btn_reset_downloads": "默认下载夹",
    "btn_reset_downloads_title": "重置为默认'下载'文件夹",
    "hint_default_download_dir": "所有新添加的种子都将默认保存到该目录",
    "lbl_listen_port": "入站连接监听端口 (0 = 自动)：",
    "lbl_max_conns": "每个种子最大连接节点数：",
    "lbl_dl_limit": "全局下载限速 (KB/s, 0 = 不限)：",
    "lbl_ul_limit": "全局上传限速 (KB/s, 0 = 不限)：",
    "sec_p2p": "P2P 网络协议",
    "lbl_enable_dht": "启用 DHT (分布式哈希表)",
    "lbl_enable_pex": "启用 PEX (节点交换机制)",
    "lbl_enable_upnp": "启用 UPnP / NAT-PMP (路由器自动端口映射)",
    "sec_win_integration": "Windows 系统集成",
    "lbl_assoc_torrent": "关联 .torrent 文件",
    "lbl_assoc_magnet": "关联浏览器 magnet: 磁力链接",
    "lbl_autostart": "开机自动启动 SmirnovTorrent",
    "sec_tray": "系统托盘与后台运行",
    "lbl_minimize_tray": "窗口最小化时缩小至系统托盘",
    "lbl_close_tray": "点击关闭按钮 [X] 时最小化到托盘而不退出",
    "lbl_start_minimized": "启动时默认最小化到系统托盘",
    "btn_save_settings": "保存设置",
    "toast_settings_saved": "设置已成功保存",
    "toast_settings_error": "保存设置失败：",
    "toast_load_error": "读取设置失败：",
    "prio_high": "高",
    "prio_normal": "普通",
    "prio_skip": "跳过",
    "status_downloading": "下载中",
    "status_seeding": "做种中",
    "status_paused": "已暂停",
    "status_checking": "校验中",
    "status_missing": "缺失",
    "status_metadata": "解析元数据",
    "status_completed": "已完成",
    "status_unknown": "未知",
    "done": "完成",
    "day_short": "天",
    "hour_short": "时",
    "min_short": "分",
    "sec_short": "秒",
    "total_dl_prefix": "总下载：",
    "total_ul_prefix": "总上传：",
    "torrents_count_prefix": "种子数：",
    "dht_nodes_suffix": "节点",
    "pieces_suffix": "块",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "已添加种子: ",
    "toast_torrent_add_error": "添加种子失败: ",
    "toast_error": "错误: ",
    "toast_upload_error": "上传错误: ",
    "toast_upload_fail": "无法上传文件",
    "toast_magnet_copied": "磁力链接已复制到剪贴板",
    "toast_hash_copied": "种子哈希已复制到剪贴板",
    "toast_file_name_copied": "文件名已复制",
    "toast_file_path_copied": "文件路径已复制",
    "toast_drop_torrent_only": "请拖放 .torrent 文件",
    "toast_drop_valid": "请拖放 .torrent 文件或磁力链接",
    "alert_invalid_magnet": "请输入有效的磁力链接 (以 magnet:? 开头)",
    "toast_magnet_added": "磁力链接已成功添加！",
    "toast_magnet_add_error": "添加磁力链接错误: ",
    "toast_action_fail": "无法执行操作",
    "toast_seq_enabled": "顺序下载已启用",
    "toast_seq_disabled": "顺序下载已禁用",
    "toast_verify_started": "已启动数据完整性校验",
    "toast_redownload_started": "种子已重新启动以重新下载",
    "toast_open_folder_error": "无法打开文件夹: ",
    "toast_folder_not_found": "未找到文件夹",
    "toast_open_file_error": "无法打开文件: ",
    "toast_file_not_found": "磁盘上未找到文件",
    "toast_explorer_error": "资源管理器错误",
    "toast_prio_updated": "文件优先级已更新",
    "toast_torrent_deleted": "种子已删除",
    "toast_delete_error": "删除错误: ",
    "toast_browse_dir_error": "选择文件夹错误: ",
    "not_specified": "未指定",
    "working": "工作中",
    "unknown": "未知",
    "untitled": "无标题",
    "seq_on": " (开启)",
    "shortcut_space": "空格",
    "dht_active": "活跃"
  },
  "ja": {
    "drag_title": "トレントを追加",
    "drag_desc": ".torrent ファイルまたは Magnet リンクをここにドロップ",
    "btn_resume": "再開",
    "btn_resume_title": "選択したトレントを再開",
    "btn_pause": "一時停止",
    "btn_pause_title": "選択したトレントを一時停止",
    "btn_delete": "削除",
    "btn_delete_title": "トレントを削除 (Delete)",
    "btn_folder": "フォルダ",
    "btn_folder_title": "ダウンロードフォルダをエクスプローラーで開く",
    "btn_sequential": "順次ダウンロード",
    "btn_sequential_title": "順次ダウンロード（ストリーミング再生向け）",
    "btn_verify": "再検査",
    "btn_verify_title": "データの整合性を再確認",
    "btn_redownload": "再ダウンロード",
    "btn_redownload_title": "再ダウンロード / 欠落ファイルを修復",
    "btn_settings": "設定",
    "btn_settings_title": "アプリケーション設定",
    "search_placeholder": "トレントを検索...",
    "cat_title": "カテゴリ",
    "cat_all": "すべてのトレント",
    "cat_downloading": "ダウンロード中",
    "cat_seeding": "シード中",
    "cat_completed": "完了",
    "cat_paused": "一時停止中",
    "cat_missing": "ファイル不明",
    "engine_active": "クライアント稼働中",
    "col_name": "ファイル名",
    "col_size": "サイズ",
    "col_progress": "進捗",
    "col_status": "状態",
    "col_dlspeed": "下り速度",
    "col_ulspeed": "上り速度",
    "col_eta": "残り時間",
    "col_peers": "ピア",
    "col_seeds": "シード",
    "col_ratio": "共有比",
    "col_added": "追加日時",
    "empty_title": "トレントがまだありません",
    "empty_desc": ".torrent ファイルをドロップするか下のボタンを使用してください",
    "empty_select_file": ".torrent ファイルを選択",
    "empty_paste_magnet": "Magnet リンクを貼り付け",
    "tab_general": "一般",
    "tab_files": "ファイル",
    "tab_peers": "ピア",
    "tab_trackers": "トラッカー",
    "tab_pieces": "ピースマップ",
    "sec_transfer": "データ転送",
    "lbl_downloaded": "ダウンロード済み:",
    "lbl_uploaded": "アップロード済み:",
    "lbl_dl_speed": "下り速度:",
    "lbl_ul_speed": "上り速度:",
    "lbl_eta": "残り時間:",
    "lbl_ratio": "共有比:",
    "sec_torrent_details": "トレント情報",
    "lbl_total_size": "合計サイズ:",
    "lbl_piece_size": "ピースサイズ:",
    "lbl_piece_count": "ピース数:",
    "lbl_info_hash": "ハッシュ (InfoHash):",
    "lbl_save_path": "保存先パス:",
    "lbl_created_date": "作成日時:",
    "lbl_created_by": "作成ツール:",
    "col_file_name": "ファイル名",
    "col_file_size": "サイズ",
    "col_file_prog": "進捗",
    "col_file_prio": "優先度",
    "col_peer_ip": "IP アドレス / ポート",
    "col_peer_client": "クライアント",
    "col_peer_dl": "受信速度",
    "col_peer_ul": "送信速度",
    "col_peer_prog": "ピア進捗",
    "col_tracker_url": "トラッカー URL",
    "col_tracker_status": "状態",
    "legend_done": "完了",
    "legend_part": "一部完了",
    "legend_wait": "待機中",
    "status_ready": "準備完了",
    "ctx_resume": "再開",
    "ctx_pause": "一時停止",
    "ctx_open_folder": "フォルダを開く",
    "ctx_sequential": "順次ダウンロード",
    "ctx_verify": "再検査",
    "ctx_redownload": "再ダウンロード",
    "ctx_copy_magnet": "Magnet リンクをコピー",
    "ctx_copy_hash": "ハッシュをコピー",
    "ctx_delete": "削除...",
    "ctx_file_open": "ファイルを開く",
    "ctx_file_show_folder": "フォルダで表示",
    "ctx_prio_header": "優先度",
    "ctx_prio_high": "高優先度",
    "ctx_prio_normal": "通常優先度",
    "ctx_prio_none": "ダウンロードしない",
    "ctx_file_copy_name": "ファイル名をコピー",
    "ctx_file_copy_path": "フルパスをコピー",
    "modal_magnet_title": "Magnet リンクを追加",
    "modal_magnet_label": "Magnet URI を入力または貼り付け:",
    "btn_cancel": "キャンセル",
    "btn_submit_magnet": "ダウンロードに追加",
    "modal_delete_title": "トレントの削除",
    "modal_delete_confirm": "本当に削除しますか",
    "modal_delete_files": "ダウンロード済みのファイルもディスクから完全に削除する",
    "modal_settings_title": "アプリケーション設定",
    "settings_section_appearance": "外観と表示言語",
    "setting_theme_label": "テーマ:",
    "theme_dark": "ダーク",
    "theme_light": "ライト",
    "setting_language_label": "表示言語:",
    "lbl_default_download_dir": "デフォルトのダウンロード先フォルダ:",
    "btn_browse": "参照...",
    "btn_reset_downloads": "ダウンロードへ",
    "btn_reset_downloads_title": "標準の「ダウンロード」フォルダに戻す",
    "hint_default_download_dir": "新しいトレントはすべてこのフォルダに保存されます",
    "lbl_listen_port": "受信接続ポート (0 = 自動):",
    "lbl_max_conns": "トレントあたりの最大ピア数:",
    "lbl_dl_limit": "ダウンロード制限 (KB/s, 0 = 無制限):",
    "lbl_ul_limit": "アップロード制限 (KB/s, 0 = 無制限):",
    "sec_p2p": "P2P ネットワークプロトコル",
    "lbl_enable_dht": "DHT を有効にする (分散ハッシュテーブル)",
    "lbl_enable_pex": "PEX を有効にする (ピア交換)",
    "lbl_enable_upnp": "UPnP / NAT-PMP を有効にする (ルーター自動ポート開放)",
    "sec_win_integration": "Windows 統合",
    "lbl_assoc_torrent": ".torrent ファイルに関連付ける",
    "lbl_assoc_magnet": "magnet: リンクに関連付ける",
    "lbl_autostart": "Windows 起動時に自動開始",
    "sec_tray": "システムトレイとバックグラウンド動作",
    "lbl_minimize_tray": "最小化時にシステムトレイに格納する",
    "lbl_close_tray": "閉じるボタン [X] で終了せずトレイに最小化する",
    "lbl_start_minimized": "トレイに最小化した状態で起動する",
    "btn_save_settings": "設定を保存",
    "toast_settings_saved": "設定を正常に保存しました",
    "toast_settings_error": "保存エラー: ",
    "toast_load_error": "設定の読み込みに失敗しました: ",
    "prio_high": "高",
    "prio_normal": "標準",
    "prio_skip": "除外",
    "status_downloading": "ダウンロード中",
    "status_seeding": "シード中",
    "status_paused": "一時停止",
    "status_checking": "検査中",
    "status_missing": "ファイル不明",
    "status_metadata": "メタデータ取得中",
    "status_completed": "完了",
    "status_unknown": "不明",
    "done": "完了",
    "day_short": "日",
    "hour_short": "時間",
    "min_short": "分",
    "sec_short": "秒",
    "total_dl_prefix": "累計ダウンロード: ",
    "total_ul_prefix": "累計アップロード: ",
    "torrents_count_prefix": "トレント数: ",
    "dht_nodes_suffix": "ノード",
    "pieces_suffix": "ピース",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "トレントが追加されました: ",
    "toast_torrent_add_error": "トレントの追加エラー: ",
    "toast_error": "エラー: ",
    "toast_upload_error": "アップロードエラー: ",
    "toast_upload_fail": "ファイルをアップロードできませんでした",
    "toast_magnet_copied": "マグネットリンクをクリップボードにコピーしました",
    "toast_hash_copied": "トレントハッシュをクリップボードにコピーしました",
    "toast_file_name_copied": "ファイル名をコピーしました",
    "toast_file_path_copied": "ファイルパスをコピーしました",
    "toast_drop_torrent_only": ".torrent ファイルをドロップしてください",
    "toast_drop_valid": ".torrent ファイルまたはマグネットリンクをドロップしてください",
    "alert_invalid_magnet": "有効なマグネットリンクを入力してください (magnet:? で始まります)",
    "toast_magnet_added": "マグネットリンクが正常に追加されました！",
    "toast_magnet_add_error": "マグネット追加エラー: ",
    "toast_action_fail": "アクションを実行できませんでした",
    "toast_seq_enabled": "シーケンシャルダウンロードが有効になりました",
    "toast_seq_disabled": "シーケンシャルダウンロードが無効になりました",
    "toast_verify_started": "データの整合性チェックが開始されました",
    "toast_redownload_started": "トレントが再ダウンロードのために再開されました",
    "toast_open_folder_error": "フォルダを開けませんでした: ",
    "toast_folder_not_found": "フォルダが見つかりません",
    "toast_open_file_error": "ファイルを開けませんでした: ",
    "toast_file_not_found": "ディスク上にファイルが見つかりません",
    "toast_explorer_error": "エクスプローラーエラー",
    "toast_prio_updated": "ファイルの優先順位が更新されました",
    "toast_torrent_deleted": "トレントが削除されました",
    "toast_delete_error": "削除エラー: ",
    "toast_browse_dir_error": "フォルダ選択エラー: ",
    "not_specified": "指定なし",
    "working": "動作中",
    "unknown": "不明",
    "untitled": "無題",
    "seq_on": " (オン)",
    "shortcut_space": "スペース",
    "dht_active": "アクティブ"
  },
  "pl": {
    "drag_title": "Dodaj torrent",
    "drag_desc": "Upuść plik .torrent lub link Magnet tutaj",
    "btn_resume": "Wznów",
    "btn_resume_title": "Wznów wybrany torrent",
    "btn_pause": "Pauza",
    "btn_pause_title": "Wstrzymaj wybrany torrent",
    "btn_delete": "Usuń",
    "btn_delete_title": "Usuń torrent (Delete)",
    "btn_folder": "Folder",
    "btn_folder_title": "Otwórz folder pobierania w Eksploratorze",
    "btn_sequential": "Sekwencyjnie",
    "btn_sequential_title": "Pobieranie sekwencyjne (do odtwarzania w locie)",
    "btn_verify": "Sprawdź",
    "btn_verify_title": "Sprawdź spójność danych",
    "btn_redownload": "Pobierz ponownie",
    "btn_redownload_title": "Pobierz ponownie / Napraw brakujące pliki",
    "btn_settings": "Ustawienia",
    "btn_settings_title": "Ustawienia programu",
    "search_placeholder": "Szukaj torrentów...",
    "cat_title": "Kategorie",
    "cat_all": "Wszystkie",
    "cat_downloading": "Pobieranie",
    "cat_seeding": "Wysyłanie",
    "cat_completed": "Zakończone",
    "cat_paused": "Wstrzymane",
    "cat_missing": "Brakujące",
    "engine_active": "Klient aktywny",
    "col_name": "Nazwa",
    "col_size": "Rozmiar",
    "col_progress": "Postęp",
    "col_status": "Stan",
    "col_dlspeed": "Pobieranie",
    "col_ulspeed": "Wysyłanie",
    "col_eta": "Pozostało",
    "col_peers": "Peerzy",
    "col_seeds": "Seedy",
    "col_ratio": "Ratio",
    "col_added": "Dodano",
    "empty_title": "Brak torrentów",
    "empty_desc": "Przeciągnij plik .torrent lub użyj przycisków poniżej",
    "empty_select_file": "Wybierz plik .torrent",
    "empty_paste_magnet": "Wklej link Magnet",
    "tab_general": "Ogólne",
    "tab_files": "Pliki",
    "tab_peers": "Peerzy",
    "tab_trackers": "Trackery",
    "tab_pieces": "Mapa części",
    "sec_transfer": "Transfer danych",
    "lbl_downloaded": "Pobrano:",
    "lbl_uploaded": "Wysłano:",
    "lbl_dl_speed": "Prędkość pobierania:",
    "lbl_ul_speed": "Prędkość wysyłania:",
    "lbl_eta": "Pozostały czas:",
    "lbl_ratio": "Ratio:",
    "sec_torrent_details": "Informacje o torrencie",
    "lbl_total_size": "Całkowity rozmiar:",
    "lbl_piece_size": "Rozmiar części:",
    "lbl_piece_count": "Liczba części:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Ścieżka zapisu:",
    "lbl_created_date": "Data utworzenia:",
    "lbl_created_by": "Utworzony przez:",
    "col_file_name": "Nazwa pliku",
    "col_file_size": "Rozmiar",
    "col_file_prog": "Postęp",
    "col_file_prio": "Priorytet",
    "col_peer_ip": "Adres IP / Port",
    "col_peer_client": "Klient",
    "col_peer_dl": "Pobieranie",
    "col_peer_ul": "Wysyłanie",
    "col_peer_prog": "Postęp peera",
    "col_tracker_url": "Adres trackera",
    "col_tracker_status": "Stan",
    "legend_done": "Pobrano",
    "legend_part": "Częściowo",
    "legend_wait": "Oczekuje",
    "status_ready": "Gotowy",
    "ctx_resume": "Wznów",
    "ctx_pause": "Pauza",
    "ctx_open_folder": "Otwórz folder",
    "ctx_sequential": "Sekwencyjnie",
    "ctx_verify": "Sprawdź",
    "ctx_redownload": "Pobierz ponownie",
    "ctx_copy_magnet": "Kopiuj link Magnet",
    "ctx_copy_hash": "Kopiuj hash",
    "ctx_delete": "Usuń...",
    "ctx_file_open": "Otwórz plik",
    "ctx_file_show_folder": "Pokaż w folderze",
    "ctx_prio_header": "Priorytet",
    "ctx_prio_high": "Wysoki priorytet",
    "ctx_prio_normal": "Zwykły priorytet",
    "ctx_prio_none": "Nie pobieraj",
    "ctx_file_copy_name": "Kopiuj nazwę",
    "ctx_file_copy_path": "Kopiuj pełną ścieżkę",
    "modal_magnet_title": "Dodaj link Magnet",
    "modal_magnet_label": "Wpisz lub wklej URI Magnet:",
    "btn_cancel": "Anuluj",
    "btn_submit_magnet": "Dodaj do pobierania",
    "modal_delete_title": "Usuwanie torrenta",
    "modal_delete_confirm": "Czy na pewno chcesz usunąć",
    "modal_delete_files": "Usuń także pobrane pliki trwale z dysku",
    "modal_settings_title": "Ustawienia programu",
    "settings_section_appearance": "Wygląd i język",
    "setting_theme_label": "Motyw:",
    "theme_dark": "Ciemny",
    "theme_light": "Jasny",
    "setting_language_label": "Język interfejsu:",
    "lbl_default_download_dir": "Domyślny folder pobierania:",
    "btn_browse": "Przeglądaj...",
    "btn_reset_downloads": "Do Pobranych",
    "btn_reset_downloads_title": "Przywróć standardowy folder 'Pobrane'",
    "hint_default_download_dir": "Wszystkie nowe torrenty będą zapisywane w tym folderze",
    "lbl_listen_port": "Port połączeń przychodzących (0 = auto):",
    "lbl_max_conns": "Maks. peerów na torrent:",
    "lbl_dl_limit": "Limit pobierania (KB/s, 0 = bez limitu):",
    "lbl_ul_limit": "Limit wysyłania (KB/s, 0 = bez limitu):",
    "sec_p2p": "Protokoły sieciowe P2P",
    "lbl_enable_dht": "Włącz DHT (Distributed Hash Table)",
    "lbl_enable_pex": "Włącz PEX (Peer Exchange)",
    "lbl_enable_upnp": "Włącz UPnP / NAT-PMP (przekierowanie portów)",
    "sec_win_integration": "Integracja z Windows",
    "lbl_assoc_torrent": "Skojarz z plikami .torrent",
    "lbl_assoc_magnet": "Skojarz z linkami magnet:",
    "lbl_autostart": "Uruchamiaj SmirnovTorrent przy starcie Windows",
    "sec_tray": "Zasobnik systemowy i tło",
    "lbl_minimize_tray": "Minimalizuj do zasobnika systemowego",
    "lbl_close_tray": "Przycisk zamknij [X] minimalizuje zamiast zamykać",
    "lbl_start_minimized": "Uruchamiaj zminimalizowany w zasobniku",
    "btn_save_settings": "Zapisz ustawienia",
    "toast_settings_saved": "Ustawienia zostały pomyślnie zapisane",
    "toast_settings_error": "Błąd zapisu: ",
    "toast_load_error": "Nie udało się wczytać ustawień: ",
    "prio_high": "Wysoki",
    "prio_normal": "Zwykły",
    "prio_skip": "Pomiń",
    "status_downloading": "Pobieranie",
    "status_seeding": "Wysyłanie",
    "status_paused": "Wstrzymano",
    "status_checking": "Sprawdzanie",
    "status_missing": "Brakujące",
    "status_metadata": "Pobieranie info",
    "status_completed": "Zakończono",
    "status_unknown": "Nieznany",
    "done": "Gotowe",
    "day_short": "d",
    "hour_short": "g",
    "min_short": "m",
    "sec_short": "s",
    "total_dl_prefix": "Łącznie pobrano: ",
    "total_ul_prefix": "Łącznie wysłano: ",
    "torrents_count_prefix": "Torrenty: ",
    "dht_nodes_suffix": "węzłów",
    "pieces_suffix": "części",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Dodano torrent: ",
    "toast_torrent_add_error": "Błąd podczas dodawania torrenta: ",
    "toast_error": "Błąd: ",
    "toast_upload_error": "Błąd przesyłania: ",
    "toast_upload_fail": "Nie udało się przesłać pliku",
    "toast_magnet_copied": "Link Magnet skopiowany do schowka",
    "toast_hash_copied": "Hash torrenta skopiowany do schowka",
    "toast_file_name_copied": "Nazwa pliku skopiowana",
    "toast_file_path_copied": "Ścieżka do pliku skopiowana",
    "toast_drop_torrent_only": "Upuść plik .torrent",
    "toast_drop_valid": "Upuść plik .torrent lub link Magnet",
    "alert_invalid_magnet": "Wprowadź prawidłowy link Magnet (zaczyna się od magnet:?)",
    "toast_magnet_added": "Link Magnet został pomyślnie dodany!",
    "toast_magnet_add_error": "Błąd podczas dodawania Magnet: ",
    "toast_action_fail": "Nie udało się wykonać akcji",
    "toast_seq_enabled": "Pobieranie sekwencyjne włączone",
    "toast_seq_disabled": "Pobieranie sekwencyjne wyłączone",
    "toast_verify_started": "Uruchomiono sprawdzanie spójności danych",
    "toast_redownload_started": "Torrent zrestartowany w celu ponownego pobrania",
    "toast_open_folder_error": "Nie udało się otworzyć folderu: ",
    "toast_folder_not_found": "Nie znaleziono folderu",
    "toast_open_file_error": "Nie udało się otworzyć pliku: ",
    "toast_file_not_found": "Nie znaleziono pliku na dysku",
    "toast_explorer_error": "Błąd Eksploratora",
    "toast_prio_updated": "Zaktualizowano priorytet pliku",
    "toast_torrent_deleted": "Torrent usunięty",
    "toast_delete_error": "Błąd podczas usuwania: ",
    "toast_browse_dir_error": "Błąd wyboru folderu: ",
    "not_specified": "Nie określono",
    "working": "Działa",
    "unknown": "Nieznany",
    "untitled": "Bez tytułu",
    "seq_on": " (Wł.)",
    "shortcut_space": "Spacja",
    "dht_active": "Aktywny"
  },
  "tr": {
    "drag_title": "Torrent Ekle",
    "drag_desc": ".torrent dosyasını veya Magnet bağlantısını buraya bırakın",
    "btn_resume": "Devam Et",
    "btn_resume_title": "Seçili torrenti devam ettir",
    "btn_pause": "Duraklat",
    "btn_pause_title": "Seçili torrenti duraklat",
    "btn_delete": "Sil",
    "btn_delete_title": "Torrenti sil (Delete)",
    "btn_folder": "Klasör",
    "btn_folder_title": "İndirme klasörünü Gezginde aç",
    "btn_sequential": "Sıralı İndir",
    "btn_sequential_title": "Sıralı indirme (akış için)",
    "btn_verify": "Doğrula",
    "btn_verify_title": "Veri bütünlüğünü kontrol et",
    "btn_redownload": "Yeniden İndir",
    "btn_redownload_title": "Yeniden indir / Eksik dosyaları onar",
    "btn_settings": "Ayarlar",
    "btn_settings_title": "Program Ayarları",
    "search_placeholder": "Torrent ara...",
    "cat_title": "Kategoriler",
    "cat_all": "Tüm Torrentler",
    "cat_downloading": "İndirilenler",
    "cat_seeding": "Gönderilenler",
    "cat_completed": "Tamamlananlar",
    "cat_paused": "Duraklatılanlar",
    "cat_missing": "Eksik Dosyalar",
    "engine_active": "İstemci aktif",
    "col_name": "Dosya Adı",
    "col_size": "Boyut",
    "col_progress": "İlerleme",
    "col_status": "Durum",
    "col_dlspeed": "İndirme",
    "col_ulspeed": "Gönderme",
    "col_eta": "Kalan Süre",
    "col_peers": "Eşler",
    "col_seeds": "Kaynaklar",
    "col_ratio": "Oran",
    "col_added": "Eklendi",
    "empty_title": "Henüz torrent yok",
    "empty_desc": "Bir .torrent dosyasını buraya sürükleyin veya aşağıdaki butonları kullanın",
    "empty_select_file": ".torrent Dosyası Seç",
    "empty_paste_magnet": "Magnet Bağlantısı Yapıştır",
    "tab_general": "Genel",
    "tab_files": "Dosyalar",
    "tab_peers": "Eşler",
    "tab_trackers": "İzleyiciler",
    "tab_pieces": "Parça Haritası",
    "sec_transfer": "Veri Transferi",
    "lbl_downloaded": "İndirilen:",
    "lbl_uploaded": "Gönderilen:",
    "lbl_dl_speed": "İndirme Hızı:",
    "lbl_ul_speed": "Gönderme Hızı:",
    "lbl_eta": "Kalan Süre:",
    "lbl_ratio": "Oran (Ratio):",
    "sec_torrent_details": "Torrent Bilgileri",
    "lbl_total_size": "Toplam Boyut:",
    "lbl_piece_size": "Parça Boyutu:",
    "lbl_piece_count": "Parça Sayısı:",
    "lbl_info_hash": "InfoHash:",
    "lbl_save_path": "Kayıt Klasörü:",
    "lbl_created_date": "Oluşturulma Tarihi:",
    "lbl_created_by": "Oluşturan:",
    "col_file_name": "Dosya Adı",
    "col_file_size": "Boyut",
    "col_file_prog": "İlerleme",
    "col_file_prio": "Öncelik",
    "col_peer_ip": "IP Adresi / Port",
    "col_peer_client": "İstemci",
    "col_peer_dl": "İndirme Hızı",
    "col_peer_ul": "Gönderme Hızı",
    "col_peer_prog": "Eş İlerlemesi",
    "col_tracker_url": "İzleyici Adresi",
    "col_tracker_status": "Durum",
    "legend_done": "Tamamlandı",
    "legend_part": "Kısmi",
    "legend_wait": "Bekliyor",
    "status_ready": "Hazır",
    "ctx_resume": "Devam Et",
    "ctx_pause": "Duraklat",
    "ctx_open_folder": "Klasörü Aç",
    "ctx_sequential": "Sıralı İndir",
    "ctx_verify": "Doğrula",
    "ctx_redownload": "Yeniden İndir",
    "ctx_copy_magnet": "Magnet Bağlantısını Kopyala",
    "ctx_copy_hash": "Hash Kopyala",
    "ctx_delete": "Sil...",
    "ctx_file_open": "Dosyayı Aç",
    "ctx_file_show_folder": "Klasörde Göster",
    "ctx_prio_header": "Öncelik",
    "ctx_prio_high": "Yüksek Öncelik",
    "ctx_prio_normal": "Normal Öncelik",
    "ctx_prio_none": "İndirme",
    "ctx_file_copy_name": "Dosya Adını Kopyala",
    "ctx_file_copy_path": "Tam Yolu Kopyala",
    "modal_magnet_title": "Magnet Bağlantısı Ekle",
    "modal_magnet_label": "Magnet URI girin veya yapıştırın:",
    "btn_cancel": "İptal",
    "btn_submit_magnet": "İndirmelere Ekle",
    "modal_delete_title": "Torrenti Sil",
    "modal_delete_confirm": "Silmek istediğinize emin misiniz",
    "modal_delete_files": "İndirilen dosyaları da diskten kalıcı olarak sil",
    "modal_settings_title": "Program Ayarları",
    "settings_section_appearance": "Görünüm ve Dil",
    "setting_theme_label": "Tema:",
    "theme_dark": "Karanlık",
    "theme_light": "Aydınlık",
    "setting_language_label": "Arayüz Dili:",
    "lbl_default_download_dir": "Varsayılan İndirme Klasörü:",
    "btn_browse": "Gözat...",
    "btn_reset_downloads": "İndirilenler'e",
    "btn_reset_downloads_title": "Varsayılan 'İndirilenler' klasörüne sıfırla",
    "hint_default_download_dir": "Tüm yeni torrentler bu klasöre kaydedilecektir",
    "lbl_listen_port": "Gelen Bağlantı Portu (0 = otomatik):",
    "lbl_max_conns": "Torrent Başına Maks. Eş:",
    "lbl_dl_limit": "İndirme Sınırı (KB/s, 0 = sınırsız):",
    "lbl_ul_limit": "Gönderme Sınırı (KB/s, 0 = sınırsız):",
    "sec_p2p": "P2P Ağ Protokolleri",
    "lbl_enable_dht": "DHT'yi Etkinleştir (Distributed Hash Table)",
    "lbl_enable_pex": "PEX'i Etkinleştir (Peer Exchange)",
    "lbl_enable_upnp": "UPnP / NAT-PMP'yi Etkinleştir (port yönlendirme)",
    "sec_win_integration": "Windows Entegrasyonu",
    "lbl_assoc_torrent": ".torrent dosyalarıyla ilişkilendir",
    "lbl_assoc_magnet": "magnet: bağlantılarıyla ilişkilendir",
    "lbl_autostart": "Windows açılışında otomatik başlat",
    "sec_tray": "Sistem Tepsisi ve Arka Plan",
    "lbl_minimize_tray": "Küçültüldüğünde sistem tepsisine gönder",
    "lbl_close_tray": "Kapat düğmesi [X] çıkmak yerine tepsiye küçültür",
    "lbl_start_minimized": "Tepside simge durumunda başlat",
    "btn_save_settings": "Ayarları Kaydet",
    "toast_settings_saved": "Ayarlar başarıyla kaydedildi",
    "toast_settings_error": "Kaydetme hatası: ",
    "toast_load_error": "Ayarlar yüklenemedi: ",
    "prio_high": "Yüksek",
    "prio_normal": "Normal",
    "prio_skip": "Atla",
    "status_downloading": "İndiriliyor",
    "status_seeding": "Gönderiliyor",
    "status_paused": "Duraklatıldı",
    "status_checking": "Doğrulanıyor",
    "status_missing": "Eksik",
    "status_metadata": "Bilgi alınıyor",
    "status_completed": "Tamamlandı",
    "status_unknown": "Bilinmiyor",
    "done": "Tamam",
    "day_short": "g",
    "hour_short": "s",
    "min_short": "d",
    "sec_short": "sn",
    "total_dl_prefix": "Toplam İndirilen: ",
    "total_ul_prefix": "Toplam Gönderilen: ",
    "torrents_count_prefix": "Torrentler: ",
    "dht_nodes_suffix": "düğüm",
    "pieces_suffix": "parça",
    "byte_b": "B",
    "byte_kb": "KB",
    "byte_mb": "MB",
    "byte_gb": "GB",
    "byte_tb": "TB",
    "speed_unit": "/s",
    "toast_torrent_added": "Torrent eklendi: ",
    "toast_torrent_add_error": "Torrent ekleme hatası: ",
    "toast_error": "Hata: ",
    "toast_upload_error": "Yükleme hatası: ",
    "toast_upload_fail": "Dosya yüklenemedi",
    "toast_magnet_copied": "Magnet bağlantısı panoya kopyalandı",
    "toast_hash_copied": "Torrent karması panoya kopyalandı",
    "toast_file_name_copied": "Dosya adı kopyalandı",
    "toast_file_path_copied": "Dosya yolu kopyalandı",
    "toast_drop_torrent_only": "Lütfen bir .torrent dosyası bırakın",
    "toast_drop_valid": "Lütfen bir .torrent dosyası veya Magnet bağlantısı bırakın",
    "alert_invalid_magnet": "Lütfen geçerli bir Magnet bağlantısı girin (magnet:? ile başlar)",
    "toast_magnet_added": "Magnet bağlantısı başarıyla eklendi!",
    "toast_magnet_add_error": "Magnet ekleme hatası: ",
    "toast_action_fail": "İşlem gerçekleştirilemedi",
    "toast_seq_enabled": "Sıralı indirme etkinleştirildi",
    "toast_seq_disabled": "Sıralı indirme devre dışı bırakıldı",
    "toast_verify_started": "Veri bütünlüğü doğrulaması başlatıldı",
    "toast_redownload_started": "Torrent yeniden indirilmek üzere başlatıldı",
    "toast_open_folder_error": "Klasör açılamadı: ",
    "toast_folder_not_found": "Klasör bulunamadı",
    "toast_open_file_error": "Dosya açılamadı: ",
    "toast_file_not_found": "Dosya diskte bulunamadı",
    "toast_explorer_error": "Gezgin hatası",
    "toast_prio_updated": "Dosya önceliği güncellendi",
    "toast_torrent_deleted": "Torrent silindi",
    "toast_delete_error": "Silme hatası: ",
    "toast_browse_dir_error": "Klasör seçme hatası: ",
    "not_specified": "Belirtilmemiş",
    "working": "Çalışıyor",
    "unknown": "Bilinmiyor",
    "untitled": "Başlıksız",
    "seq_on": " (Açık)",
    "shortcut_space": "Boşluk",
    "dht_active": "Etkin"
  }
};

let state = {
  torrents: [],
  selectedHash: null,
  activeFilter: 'all',
  searchQuery: '',
  sortField: 'added_at',
  sortAsc: false,
  activeTab: 'general',
  pollInterval: null,
  defaultDownloadDir: '',
  theme: (function() {
    try { return localStorage.getItem('smirnovtorrent_theme') || 'dark'; } catch(e) { return 'dark'; }
  })(),
  language: (function() {
    try {
      const stored = localStorage.getItem('smirnovtorrent_language');
      if (stored === 'uk') {
        try { localStorage.setItem('smirnovtorrent_language', 'ru'); } catch(e) {}
        return 'ru';
      }
      return (stored && i18n[stored]) ? stored : 'ru';
    } catch(e) { return 'ru'; }
  })(),
  contextFile: null, // { hash, fileIndex, path, name }
  settingsInitialTheme: null,
  settingsInitialLanguage: null
};

function t(key, fallback = '') {
  const dict = i18n[state.language] || i18n.ru;
  if (dict && dict[key] !== undefined) return dict[key];
  if (i18n.ru && i18n.ru[key] !== undefined) return i18n.ru[key];
  return fallback;
}

function applyTheme(theme) {
  if (theme !== 'light' && theme !== 'dark') theme = 'dark';
  state.theme = theme;
  document.documentElement.setAttribute('data-theme', theme);
  const logo = document.querySelector('.brand-logo');
  if (logo) {
    logo.src = theme === 'light' ? 'smirnovtorrent-light.jpg' : 'smirnovtorrent.jpg';
  }
  const favicon = document.querySelector('link[rel="icon"]');
  if (favicon) {
    favicon.href = theme === 'light' ? 'smirnovtorrent-light.jpg' : 'smirnovtorrent.jpg';
  }
  try {
    localStorage.setItem('smirnovtorrent_theme', theme);
  } catch (e) {}
  try {
    fetch('/api/window/theme', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ theme: theme }),
    }).catch(() => {});
  } catch (e) {}
}

function applyLanguage(lang) {
  if (!i18n[lang] || lang === 'uk') lang = 'ru';
  state.language = lang;
  try {
    localStorage.setItem('smirnovtorrent_language', lang);
  } catch (e) {}
  document.documentElement.setAttribute('lang', lang);
  syncCustomLanguageSelect(lang);

  const dict = i18n[lang];
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    if (dict && dict[key] !== undefined) {
      el.textContent = dict[key];
    }
  });

  document.querySelectorAll('[data-i18n-title]').forEach(el => {
    const key = el.getAttribute('data-i18n-title');
    if (dict && dict[key] !== undefined) {
      el.title = dict[key];
    }
  });

  document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
    const key = el.getAttribute('data-i18n-placeholder');
    if (dict && dict[key] !== undefined) {
      el.placeholder = dict[key];
    }
  });

  if (!lastGlobalStats) {
    const countEl = document.getElementById('status-torrents-count');
    if (countEl) countEl.textContent = t('torrents_count_prefix', 'Торрентов: ') + (state.torrents ? state.torrents.length : 0);
    const dhtEl = document.getElementById('status-dht-nodes');
    if (dhtEl) dhtEl.textContent = 'DHT: ' + t('dht_active', 'Активен');
    const totalDlEl = document.getElementById('global-total-dl');
    if (totalDlEl) totalDlEl.textContent = t('total_dl_prefix', 'Всего скачано: ') + '0 ' + t('byte_b', 'Б');
    const totalUlEl = document.getElementById('global-total-ul');
    if (totalUlEl) totalUlEl.textContent = t('total_ul_prefix', 'Всего отдано: ') + '0 ' + t('byte_b', 'Б');
  }

  if (state.torrents && state.torrents.length > 0) {
    renderTorrents();
  }
  if (state.selectedHash) {
    fetchTorrentDetails(state.selectedHash);
  }
  updateGlobalStatsDisplay();
}

// Formatting Utilities
function formatBytes(bytes, decimals = 1) {
  const b = t('byte_b', 'Б');
  if (typeof bytes !== 'number' || isNaN(bytes) || bytes <= 0) return '0 ' + b;
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = [b, t('byte_kb', 'КБ'), t('byte_mb', 'МБ'), t('byte_gb', 'ГБ'), t('byte_tb', 'ТБ')];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const idx = Math.min(Math.max(0, i), sizes.length - 1);
  return parseFloat((bytes / Math.pow(k, idx)).toFixed(dm)) + ' ' + sizes[idx];
}

function formatSpeed(bytesPerSec) {
  const su = t('speed_unit', '/с');
  const b = t('byte_b', 'Б');
  if (!bytesPerSec || bytesPerSec <= 0) return '0 ' + b + su;
  return formatBytes(bytesPerSec) + su;
}

function formatETA(seconds) {
  if (seconds == null || seconds < 0 || seconds > 86400 * 365) return '∞';
  if (seconds === 0) return t('done', 'Готово');
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  const dS = t('day_short', 'д');
  const hS = t('hour_short', 'ч');
  const mS = t('min_short', 'м');
  const sS = t('sec_short', 'с');
  if (d > 0) return `${d}${dS} ${h}${hS}`;
  if (h > 0) return `${h}${hS} ${m}${mS}`;
  if (m > 0) return `${m}${mS} ${s}${sS}`;
  return `${s}${sS}`;
}

const localeMap = {
  ru: 'ru-RU',
  en: 'en-US',
  de: 'de-DE',
  fr: 'fr-FR',
  es: 'es-ES',
  it: 'it-IT',
  pt: 'pt-PT',
  zh: 'zh-CN',
  ja: 'ja-JP',
  pl: 'pl-PL',
  tr: 'tr-TR'
};

function formatDate(isoStr) {
  if (!isoStr) return '-';
  const d = new Date(isoStr);
  if (isNaN(d.getTime()) || d.getFullYear() <= 1970) return '-';
  const loc = localeMap[state.language] || 'ru-RU';
  return d.toLocaleDateString(loc) + ' ' + d.toLocaleTimeString(loc, { hour: '2-digit', minute: '2-digit' });
}

function escapeHtml(str) {
  if (!str) return '';
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

function getStatusInfo(status) {
  switch (status) {
    case 'downloading':
      return { label: t('status_downloading', 'Загрузка'), className: 'status-downloading', barClass: '' };
    case 'seeding':
      return { label: t('status_seeding', 'Раздача'), className: 'status-seeding', barClass: 'bar-seeding' };
    case 'paused':
      return { label: t('status_paused', 'На паузе'), className: 'status-paused', barClass: 'bar-paused' };
    case 'checking':
      return { label: t('status_checking', 'Проверка'), className: 'status-checking', barClass: 'bar-checking' };
    case 'missing':
      return { label: t('status_missing', 'Отсутствует'), className: 'status-missing', barClass: 'bar-missing' };
    case 'metadata':
      return { label: t('status_metadata', 'Получение инфо'), className: 'status-metadata', barClass: 'bar-checking' };
    default:
      return { label: status || t('status_unknown', 'Неизвестно'), className: 'status-paused', barClass: '' };
  }
}

// Toast Notifications
function showToast(message) {
  const container = document.getElementById('toast-container');
  if (!container) return;
  const toast = document.createElement('div');
  toast.className = 'toast';
  toast.innerHTML = `
    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="var(--accent-blue)" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
    <span>${escapeHtml(message)}</span>
  `;
  container.appendChild(toast);
  setTimeout(() => {
    if (toast.parentNode) toast.parentNode.removeChild(toast);
  }, 3000);
}

// DOM Elements
const tbody = document.getElementById('torrent-tbody');
const emptyState = document.getElementById('empty-state');
const filterInput = document.getElementById('filter-input');
const btnResume = document.getElementById('btn-resume');
const btnPause = document.getElementById('btn-pause');
const btnDelete = document.getElementById('btn-delete');
const btnOpenFolder = document.getElementById('btn-open-folder');
const btnSeq = document.getElementById('btn-sequential');
const btnVerify = document.getElementById('btn-verify');
const btnRedownload = document.getElementById('btn-redownload');

// Context Menus
const ctxTorrent = document.getElementById('context-menu-torrent');
const ctxFile = document.getElementById('context-menu-file');

// Initialization
document.addEventListener('DOMContentLoaded', () => {
  if (typeof document !== 'undefined' && 'fonts' in document) {
    try {
      document.fonts.load('16px "Twemoji Country Flags"');
    } catch (e) {}
  }

  setupEventListeners();
  setupResizer();
  setupDragAndDrop();
  setupKeyboardShortcuts();
  setupContextMenus();

  applyTheme(state.theme);
  applyLanguage(state.language);

  // Sync settings with server
  fetch('/api/settings')
    .then(r => r.json())
    .then(cfg => {
      if (cfg.theme && cfg.theme !== state.theme) {
        applyTheme(cfg.theme);
      }
      if (cfg.language && cfg.language !== state.language) {
        applyLanguage(cfg.language);
      }
    })
    .catch(() => {});

  fetchTorrents();
  fetchGlobalStats();

  state.pollInterval = setInterval(() => {
    fetchTorrents();
    fetchGlobalStats();
    if (state.selectedHash) {
      fetchTorrentDetails(state.selectedHash);
    }
  }, 1000);
});

// Event Listeners
function setupEventListeners() {
  // Navigation categories
  document.querySelectorAll('.nav-item').forEach(item => {
    item.addEventListener('click', () => {
      document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
      item.classList.add('active');
      state.activeFilter = item.dataset.filter;
      renderTorrents();
    });
  });

  // Search filter
  filterInput.addEventListener('input', (e) => {
    state.searchQuery = e.target.value.toLowerCase().trim();
    renderTorrents();
  });

  // Sort headers
  document.querySelectorAll('.torrent-table th').forEach(th => {
    th.addEventListener('click', () => {
      const field = th.dataset.sort;
      if (!field) return;
      if (state.sortField === field) {
        state.sortAsc = !state.sortAsc;
      } else {
        state.sortField = field;
        state.sortAsc = true;
      }
      renderTorrents();
    });
  });

  // Inspector Tabs
  document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
      document.querySelectorAll('.tab-pane').forEach(p => p.classList.remove('active'));
      btn.classList.add('active');
      const tabId = btn.dataset.tab;
      const targetPane = document.getElementById('pane-' + tabId);
      if (targetPane) targetPane.classList.add('active');
      state.activeTab = tabId;
      if (state.selectedHash) {
        fetchTorrentDetails(state.selectedHash);
      }
    });
  });

  // Top action buttons
  const btnAddFile = document.getElementById('btn-add-file');
  if (btnAddFile) btnAddFile.addEventListener('click', addTorrentFile);
  const btnAddMagnet = document.getElementById('btn-add-magnet');
  if (btnAddMagnet) btnAddMagnet.addEventListener('click', () => openModal('modal-magnet'));
  btnResume.addEventListener('click', () => performAction('resume'));
  btnPause.addEventListener('click', () => performAction('pause'));
  btnDelete.addEventListener('click', openDeleteModal);
  btnOpenFolder.addEventListener('click', openFolder);
  btnSeq.addEventListener('click', toggleSequential);
  btnVerify.addEventListener('click', verifyTorrent);
  if (btnRedownload) btnRedownload.addEventListener('click', redownloadTorrent);
  document.getElementById('btn-settings').addEventListener('click', openSettings);

  // Modal triggers
  document.querySelectorAll('[data-close]').forEach(el => {
    el.addEventListener('click', () => closeModal(el.dataset.close));
  });

  document.querySelectorAll('.modal-backdrop').forEach(backdrop => {
    backdrop.addEventListener('click', (e) => {
      if (e.target === backdrop) {
        backdrop.classList.add('hidden');
      }
    });
  });

  document.getElementById('btn-submit-magnet').addEventListener('click', submitMagnet);
  document.getElementById('btn-confirm-delete').addEventListener('click', confirmDelete);
  document.getElementById('btn-browse-dir').addEventListener('click', browseDownloadDir);
  document.getElementById('btn-save-settings').addEventListener('click', saveSettings);

  const btnResetDownloads = document.getElementById('btn-reset-downloads');
  if (btnResetDownloads) {
    btnResetDownloads.addEventListener('click', () => {
      if (state.defaultDownloadDir) {
        document.getElementById('setting-download-dir').value = state.defaultDownloadDir;
        showToast(t('toast_dir_reset', "Сброшено на стандартную папку 'Загрузки'"));
      }
    });
  }

  const themeSelect = document.getElementById('setting-theme');
  if (themeSelect) {
    themeSelect.addEventListener('change', (e) => {
      applyTheme(e.target.value);
    });
  }

  const langSelect = document.getElementById('setting-language');
  if (langSelect) {
    langSelect.addEventListener('change', (e) => {
      applyLanguage(e.target.value);
    });
  }

  setupCustomLanguageSelect();
}

// Custom Language Dropdown Helper
function syncCustomLanguageSelect(val) {
  const select = document.getElementById('setting-language');
  const targetVal = val || (select ? select.value : 'ru');
  const wrapper = document.getElementById('custom-language-select');
  if (!wrapper) return;
  const triggerText = wrapper.querySelector('.custom-select-text');
  const options = wrapper.querySelectorAll('.custom-select-option');
  options.forEach(opt => {
    if (opt.getAttribute('data-value') === targetVal) {
      opt.classList.add('selected');
      if (triggerText) triggerText.innerHTML = opt.innerHTML;
    } else {
      opt.classList.remove('selected');
    }
  });
}

function setupCustomLanguageSelect() {
  const wrapper = document.getElementById('custom-language-select');
  const select = document.getElementById('setting-language');
  if (!wrapper || !select) return;

  const trigger = wrapper.querySelector('.custom-select-trigger');
  const options = wrapper.querySelectorAll('.custom-select-option');

  syncCustomLanguageSelect(select.value);

  if (trigger) {
    trigger.addEventListener('click', (e) => {
      e.stopPropagation();
      const isOpen = wrapper.classList.toggle('open');
      trigger.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
    });

    trigger.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
        e.preventDefault();
        wrapper.classList.add('open');
        trigger.setAttribute('aria-expanded', 'true');
        const selectedOpt = wrapper.querySelector('.custom-select-option.selected') || options[0];
        if (selectedOpt) selectedOpt.focus();
      } else if (e.key === 'Escape') {
        wrapper.classList.remove('open');
        trigger.setAttribute('aria-expanded', 'false');
      }
    });
  }

  options.forEach(opt => {
    opt.setAttribute('tabindex', '0');
    opt.addEventListener('click', (e) => {
      e.stopPropagation();
      const val = opt.getAttribute('data-value');
      if (val && select.value !== val) {
        select.value = val;
        select.dispatchEvent(new Event('change'));
      }
      syncCustomLanguageSelect(val);
      wrapper.classList.remove('open');
      if (trigger) {
        trigger.setAttribute('aria-expanded', 'false');
        trigger.focus();
      }
    });

    opt.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        opt.click();
      } else if (e.key === 'ArrowDown') {
        e.preventDefault();
        const next = opt.nextElementSibling;
        if (next && next.classList.contains('custom-select-option')) next.focus();
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        const prev = opt.previousElementSibling;
        if (prev && prev.classList.contains('custom-select-option')) prev.focus();
      } else if (e.key === 'Escape') {
        wrapper.classList.remove('open');
        if (trigger) {
          trigger.setAttribute('aria-expanded', 'false');
          trigger.focus();
        }
      }
    });
  });

  document.addEventListener('click', (e) => {
    if (!wrapper.contains(e.target)) {
      wrapper.classList.remove('open');
      if (trigger) trigger.setAttribute('aria-expanded', 'false');
    }
  });

  select.addEventListener('change', () => {
    syncCustomLanguageSelect(select.value);
  });
}

// Drag & Drop Implementation
function setupDragAndDrop() {
  const overlay = document.getElementById('drag-drop-overlay');
  let dragCounter = 0;

  window.addEventListener('dragenter', (e) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounter++;
    if (overlay) overlay.classList.remove('hidden');
  });

  window.addEventListener('dragover', (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'copy';
    }
    if (overlay && overlay.classList.contains('hidden')) {
      overlay.classList.remove('hidden');
    }
  });

  window.addEventListener('dragleave', (e) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounter--;
    if (dragCounter <= 0) {
      dragCounter = 0;
      if (overlay) overlay.classList.add('hidden');
    }
  });

  window.addEventListener('drop', async (e) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounter = 0;
    if (overlay) overlay.classList.add('hidden');

    if (!e.dataTransfer) return;

    // 1. Files dropped
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      let count = 0;
      for (let i = 0; i < e.dataTransfer.files.length; i++) {
        const file = e.dataTransfer.files[i];
        if (file.name.toLowerCase().endsWith('.torrent')) {
          await uploadTorrentFile(file);
          count++;
        }
      }
      if (count === 0) {
        showToast(t("toast_drop_torrent_only", "Пожалуйста, перетащите файл с расширением .torrent"));
      }
      return;
    }

    // 2. Text / Magnet link dropped
    const text = (e.dataTransfer.getData('text/plain') || e.dataTransfer.getData('text/uri-list') || '').trim();
    if (text) {
      if (text.toLowerCase().startsWith('magnet:?')) {
        await addMagnetURI(text);
      } else {
        showToast(t("toast_drop_valid", "Пожалуйста, перетащите .torrent файл или Magnet-ссылку"));
      }
    }
  });
}

// Context Menus
function setupContextMenus() {
  // Hide menus on click anywhere or on resize/scroll
  window.addEventListener('click', () => hideContextMenus());
  window.addEventListener('resize', () => hideContextMenus());
  document.addEventListener('scroll', () => hideContextMenus(), true);

  // Torrent Menu Actions
  document.getElementById('ctx-torrent-toggle').addEventListener('click', () => {
    hideContextMenus();
    const t = getSelectedTorrent();
    if (!t) return;
    if (t.status === 'paused') {
      performAction('resume');
    } else {
      performAction('pause');
    }
  });

  document.getElementById('ctx-torrent-open-folder').addEventListener('click', () => {
    hideContextMenus();
    openFolder();
  });

  document.getElementById('ctx-torrent-sequential').addEventListener('click', () => {
    hideContextMenus();
    toggleSequential();
  });

  document.getElementById('ctx-torrent-verify').addEventListener('click', () => {
    hideContextMenus();
    verifyTorrent();
  });

  document.getElementById('ctx-torrent-redownload').addEventListener('click', () => {
    hideContextMenus();
    redownloadTorrent();
  });

  document.getElementById('ctx-torrent-copy-magnet').addEventListener('click', async () => {
    hideContextMenus();
    const t = getSelectedTorrent();
    if (!t) return;
    try {
      const res = await fetch(`/api/torrent?hash=${encodeURIComponent(t.info_hash)}`);
      const data = await res.json();
      const magnet = data.summary.magnet_uri || `magnet:?xt=urn:btih:${t.info_hash}&dn=${encodeURIComponent(t.name)}`;
      await navigator.clipboard.writeText(magnet);
      showToast(t("toast_magnet_copied", "Magnet-ссылка скопирована в буфер обмена"));
    } catch {
      const magnet = `magnet:?xt=urn:btih:${t.info_hash}&dn=${encodeURIComponent(t.name)}`;
      await navigator.clipboard.writeText(magnet);
      showToast(t("toast_magnet_copied", "Magnet-ссылка скопирована в буфер обмена"));
    }
  });

  document.getElementById('ctx-torrent-copy-hash').addEventListener('click', async () => {
    hideContextMenus();
    const t = getSelectedTorrent();
    if (!t) return;
    await navigator.clipboard.writeText(t.info_hash);
    showToast(t("toast_hash_copied", "Хэш торрента скопирован в буфер обмена"));
  });

  document.getElementById('ctx-torrent-delete').addEventListener('click', () => {
    hideContextMenus();
    openDeleteModal();
  });

  // File Menu Actions
  document.getElementById('ctx-file-open').addEventListener('click', () => {
    hideContextMenus();
    if (!state.contextFile) return;
    openTargetFile(state.contextFile.hash, state.contextFile.fileIndex);
  });

  document.getElementById('ctx-file-show-folder').addEventListener('click', () => {
    hideContextMenus();
    if (!state.contextFile) return;
    showFileInFolder(state.contextFile.hash, state.contextFile.fileIndex);
  });

  document.getElementById('ctx-file-prio-high').addEventListener('click', () => {
    hideContextMenus();
    if (!state.contextFile) return;
    setFilePriority(state.contextFile.hash, state.contextFile.fileIndex, 2);
  });

  document.getElementById('ctx-file-prio-normal').addEventListener('click', () => {
    hideContextMenus();
    if (!state.contextFile) return;
    setFilePriority(state.contextFile.hash, state.contextFile.fileIndex, 1);
  });

  document.getElementById('ctx-file-prio-none').addEventListener('click', () => {
    hideContextMenus();
    if (!state.contextFile) return;
    setFilePriority(state.contextFile.hash, state.contextFile.fileIndex, 0);
  });

  document.getElementById('ctx-file-copy-name').addEventListener('click', async () => {
    hideContextMenus();
    if (!state.contextFile) return;
    await navigator.clipboard.writeText(state.contextFile.name);
    showToast(t("toast_file_name_copied", "Имя файла скопировано"));
  });

  document.getElementById('ctx-file-copy-path').addEventListener('click', async () => {
    hideContextMenus();
    if (!state.contextFile) return;
    await navigator.clipboard.writeText(state.contextFile.path);
    showToast(t("toast_file_path_copied", "Путь к файлу скопирован"));
  });
}

function hideContextMenus() {
  if (ctxTorrent) ctxTorrent.classList.add('hidden');
  if (ctxFile) ctxFile.classList.add('hidden');
}

function showContextMenu(menu, x, y) {
  hideContextMenus();
  menu.classList.remove('hidden');
  const w = menu.offsetWidth || 230;
  const h = menu.offsetHeight || 280;
  const clampedX = Math.min(x, window.innerWidth - w - 10);
  const clampedY = Math.min(y, window.innerHeight - h - 10);
  menu.style.left = `${Math.max(10, clampedX)}px`;
  menu.style.top = `${Math.max(10, clampedY)}px`;
}

function getSelectedTorrent() {
  if (!state.selectedHash) return null;
  return state.torrents.find(t => t.info_hash === state.selectedHash) || null;
}

// Keyboard shortcuts
function setupKeyboardShortcuts() {
  document.addEventListener('keydown', (e) => {
    // Ctrl+F: Focus search
    if ((e.ctrlKey || e.metaKey) && (e.key === 'f' || e.key === 'F')) {
      e.preventDefault();
      filterInput.focus();
      filterInput.select();
    }
    // Ctrl+O: Add file
    if ((e.ctrlKey || e.metaKey) && (e.key === 'o' || e.key === 'O')) {
      e.preventDefault();
      addTorrentFile();
    }
    // Ctrl+M: Add magnet
    if ((e.ctrlKey || e.metaKey) && (e.key === 'm' || e.key === 'M')) {
      e.preventDefault();
      openModal('modal-magnet');
    }
    // Escape: Close context menus, modals, search
    if (e.key === 'Escape') {
      hideContextMenus();
      const openModals = document.querySelectorAll('.modal-backdrop:not(.hidden)');
      if (openModals.length > 0) {
        openModals.forEach(m => m.classList.add('hidden'));
      } else if (document.activeElement === filterInput) {
        filterInput.value = '';
        state.searchQuery = '';
        filterInput.blur();
        renderTorrents();
      }
    }
    // Delete: Prompt remove selected torrent
    if (e.key === 'Delete' && state.selectedHash && !document.querySelector('.modal-backdrop:not(.hidden)')) {
      openDeleteModal();
    }
    // Enter in magnet textarea
    if (e.key === 'Enter' && e.ctrlKey && !document.getElementById('modal-magnet').classList.contains('hidden')) {
      submitMagnet();
    }

    // Space / Enter shortcuts when not inside text inputs
    const inInput = e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA';
    if (!inInput && !document.querySelector('.modal-backdrop:not(.hidden)')) {
      if (e.key === ' ' && state.selectedHash) {
        e.preventDefault();
        const t = getSelectedTorrent();
        if (t) {
          if (t.status === 'paused') performAction('resume');
          else performAction('pause');
        }
      } else if (e.key === 'Enter') {
        e.preventDefault();
        if (state.activeTab === 'files' && state.contextFile) {
          openTargetFile(state.contextFile.hash, state.contextFile.fileIndex);
        } else if (state.selectedHash) {
          openFolder();
        }
      }
    }
  });
}

// Resizer logic
function setupResizer() {
  const resizer = document.getElementById('drag-resizer');
  const detailsView = document.getElementById('details-view');
  let isResizing = false;

  resizer.addEventListener('mousedown', (e) => {
    e.preventDefault();
    isResizing = true;
    resizer.classList.add('resizing');
    document.body.style.cursor = 'row-resize';
  });

  window.addEventListener('mousemove', (e) => {
    if (!isResizing) return;
    const containerHeight = document.querySelector('.content-split').offsetHeight;
    const newDetailsHeight = containerHeight - e.clientY + document.querySelector('.toolbar').offsetHeight;
    if (newDetailsHeight >= 100 && newDetailsHeight <= containerHeight - 100) {
      detailsView.style.height = `${newDetailsHeight}px`;
    }
  });

  window.addEventListener('mouseup', () => {
    if (isResizing) {
      isResizing = false;
      resizer.classList.remove('resizing');
      document.body.style.cursor = '';
    }
  });
}

// Fetch & Render Torrents
async function fetchTorrents() {
  try {
    const res = await fetch('/api/torrents');
    if (!res.ok) return;
    state.torrents = await res.json() || [];
    updateSidebarBadges();
    renderTorrents();
  } catch (err) {
    console.error("fetchTorrents error:", err);
  }
}

let lastGlobalStats = null;

function updateGlobalStatsDisplay() {
  if (!lastGlobalStats) return;
  const stats = lastGlobalStats;
  const dlEl = document.getElementById('global-dl-speed');
  const ulEl = document.getElementById('global-ul-speed');
  const totalDlEl = document.getElementById('global-total-dl');
  const totalUlEl = document.getElementById('global-total-ul');
  const countEl = document.getElementById('status-torrents-count');
  const dhtEl = document.getElementById('status-dht-nodes');

  if (dlEl) dlEl.textContent = formatSpeed(stats.total_download_speed);
  if (ulEl) ulEl.textContent = formatSpeed(stats.total_upload_speed);
  if (totalDlEl) totalDlEl.textContent = t('total_dl_prefix', 'Всего скачано: ') + formatBytes(stats.total_downloaded);
  if (totalUlEl) totalUlEl.textContent = t('total_ul_prefix', 'Всего отдано: ') + formatBytes(stats.total_uploaded);
  if (countEl) countEl.textContent = t('torrents_count_prefix', 'Торрентов: ') + stats.total_torrents;
  if (dhtEl) dhtEl.textContent = 'DHT: ' + stats.dht_nodes + ' ' + t('dht_nodes_suffix', 'узлов');
}

async function fetchGlobalStats() {
  try {
    const res = await fetch('/api/stats');
    if (!res.ok) return;
    const stats = await res.json();
    lastGlobalStats = stats;
    updateGlobalStatsDisplay();
  } catch (err) {
    console.error("fetchGlobalStats error:", err);
  }
}

function updateSidebarBadges() {
  const all = state.torrents.length;
  let dl = 0, seed = 0, comp = 0, paused = 0, missing = 0;
  state.torrents.forEach(t => {
    if (t.status === 'downloading') dl++;
    else if (t.status === 'seeding') seed++;
    else if (t.status === 'paused') paused++;
    else if (t.status === 'missing') missing++;
    if (t.status === 'seeding' || (t.progress >= 100 && t.status !== 'missing' && t.status !== 'error')) comp++;
  });
  document.getElementById('badge-all').textContent = all;
  document.getElementById('badge-dl').textContent = dl;
  document.getElementById('badge-seed').textContent = seed;
  document.getElementById('badge-comp').textContent = comp;
  document.getElementById('badge-paused').textContent = paused;
  const badgeMissing = document.getElementById('badge-missing');
  if (badgeMissing) badgeMissing.textContent = missing;
}

function updateSortHeaders() {
  document.querySelectorAll('.torrent-table th').forEach(th => {
    th.classList.remove('sort-asc', 'sort-desc');
    if (th.dataset.sort === state.sortField) {
      th.classList.add(state.sortAsc ? 'sort-asc' : 'sort-desc');
    }
  });
}

function renderTorrents() {
  updateSortHeaders();
  let list = [...state.torrents];

  // Category filter
  if (state.activeFilter === 'downloading') {
    list = list.filter(t => t.status === 'downloading');
  } else if (state.activeFilter === 'seeding') {
    list = list.filter(t => t.status === 'seeding');
  } else if (state.activeFilter === 'completed') {
    list = list.filter(t => t.status === 'seeding' || (t.progress >= 100 && t.status !== 'missing' && t.status !== 'error'));
  } else if (state.activeFilter === 'paused') {
    list = list.filter(t => t.status === 'paused');
  } else if (state.activeFilter === 'missing') {
    list = list.filter(t => t.status === 'missing');
  }

  // Search filter
  if (state.searchQuery) {
    list = list.filter(t => (t.name && t.name.toLowerCase().includes(state.searchQuery)) || (t.info_hash && t.info_hash.toLowerCase().includes(state.searchQuery)));
  }

  // Sorting
  list.sort((a, b) => {
    let vA = a[state.sortField];
    let vB = b[state.sortField];
    if (state.sortField === 'added_at') {
      const tA = new Date(vA).getTime() || 0;
      const tB = new Date(vB).getTime() || 0;
      return state.sortAsc ? (tA - tB) : (tB - tA);
    }
    if (state.sortField === 'eta') {
      const eA = (vA == null || vA < 0) ? Infinity : Number(vA);
      const eB = (vB == null || vB < 0) ? Infinity : Number(vB);
      if (eA === eB) return 0;
      return state.sortAsc ? (eA - eB) : (eB - eA);
    }
    if (typeof vA === 'string' || typeof vB === 'string') {
      const sA = (vA == null ? '' : String(vA));
      const sB = (vB == null ? '' : String(vB));
      return state.sortAsc ? sA.localeCompare(sB) : sB.localeCompare(sA);
    }
    const nA = (vA == null ? 0 : Number(vA)) || 0;
    const nB = (vB == null ? 0 : Number(vB)) || 0;
    return state.sortAsc ? (nA - nB) : (nB - nA);
  });

  if (list.length === 0) {
    tbody.innerHTML = '';
    emptyState.classList.add('visible');
    state.selectedHash = null;
    clearDetails();
    updateToolbarButtons(null);
    return;
  }
  emptyState.classList.remove('visible');

  // Single-pass reconciliation of table rows
  const currentHashes = new Set(list.map(t => t.info_hash));
  const existingRows = new Map();
  const children = Array.from(tbody.children);
  for (const child of children) {
    const hash = child.dataset ? child.dataset.hash : null;
    if (!hash || !currentHashes.has(hash)) {
      child.remove();
    } else {
      existingRows.set(hash, child);
    }
  }

  list.forEach((t, index) => {
    let tr = existingRows.get(t.info_hash);
    const isNew = !tr;

    const progVal = (typeof t.progress === 'number' && !isNaN(t.progress)) ? t.progress : 0;
    const progPercent = progVal.toFixed(1);
    const nameStr = t.name || t('untitled', 'Без названия');
    const formattedName = escapeHtml(nameStr);
    const formattedSize = formatBytes(t.size);
    const formattedDl = formatSpeed(t.download_speed);
    const formattedUl = formatSpeed(t.upload_speed);
    const formattedEta = formatETA(t.eta);
    const formattedPeers = `${t.peers_connected || 0} (${t.peers_total || 0})`;
    const formattedSeeds = `${t.seeds_connected || 0}`;
    const ratioVal = (typeof t.ratio === 'number' && !isNaN(t.ratio)) ? t.ratio : 0;
    const formattedRatio = ratioVal.toFixed(2);
    const formattedDate = formatDate(t.added_at);
    const sInfo = getStatusInfo(t.status);

    if (isNew) {
      tr = document.createElement('tr');
      tr.dataset.hash = t.info_hash;
      tr.innerHTML = `
        <td class="col-name" title="${formattedName}">
          <div class="torrent-title">
            <svg viewBox="0 0 24 24" class="torrent-icon" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            <span>${formattedName}</span>
          </div>
        </td>
        <td class="col-size">${formattedSize}</td>
        <td class="col-progress">
          <div class="progress-cell">
            <div class="progress-track">
              <div class="progress-bar ${sInfo.barClass}" style="width: ${progPercent}%"></div>
            </div>
            <span class="progress-text">${progPercent}%</span>
          </div>
        </td>
        <td class="col-status">
          <span class="status-pill ${sInfo.className}">${sInfo.label}</span>
        </td>
        <td class="col-dlspeed">${formattedDl}</td>
        <td class="col-ulspeed">${formattedUl}</td>
        <td class="col-eta">${formattedEta}</td>
        <td class="col-peers">${formattedPeers}</td>
        <td class="col-seeds">${formattedSeeds}</td>
        <td class="col-ratio">${formattedRatio}</td>
        <td class="col-added">${formattedDate}</td>
      `;

      tr.addEventListener('click', () => {
        selectTorrentRow(tr, t.info_hash);
      });

      tr.addEventListener('dblclick', () => {
        selectTorrentRow(tr, t.info_hash);
        openFolder();
      });

      tr.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        selectTorrentRow(tr, t.info_hash);
        populateTorrentContextMenu(t);
        showContextMenu(ctxTorrent, e.clientX, e.clientY);
      });
    } else {
      // In-place updates for existing rows
      const nameSpan = tr.querySelector('.torrent-title span');
      if (nameSpan && nameSpan.textContent !== nameStr) {
        nameSpan.textContent = nameStr;
        tr.querySelector('.col-name').title = nameStr;
      }
      const td1 = tr.children[1];
      if (td1 && td1.textContent !== formattedSize) td1.textContent = formattedSize;

      const td2 = tr.children[2];
      if (td2) {
        const bar = td2.querySelector('.progress-bar');
        const text = td2.querySelector('.progress-text');
        if (bar) {
          const expectedWidth = `${progPercent}%`;
          if (bar.style.width !== expectedWidth) bar.style.width = expectedWidth;
          const expectedClass = `progress-bar ${sInfo.barClass}`;
          if (bar.className !== expectedClass) bar.className = expectedClass;
        }
        if (text && text.textContent !== `${progPercent}%`) text.textContent = `${progPercent}%`;
      }

      const td3 = tr.children[3];
      if (td3) {
        const pill = td3.querySelector('.status-pill');
        if (pill) {
          if (pill.textContent !== sInfo.label) pill.textContent = sInfo.label;
          const expectedPillClass = `status-pill ${sInfo.className}`;
          if (pill.className !== expectedPillClass) pill.className = expectedPillClass;
        }
      }

      const td4 = tr.children[4];
      if (td4 && td4.textContent !== formattedDl) td4.textContent = formattedDl;
      const td5 = tr.children[5];
      if (td5 && td5.textContent !== formattedUl) td5.textContent = formattedUl;
      const td6 = tr.children[6];
      if (td6 && td6.textContent !== formattedEta) td6.textContent = formattedEta;
      const td7 = tr.children[7];
      if (td7 && td7.textContent !== formattedPeers) td7.textContent = formattedPeers;
      const td8 = tr.children[8];
      if (td8 && td8.textContent !== formattedSeeds) td8.textContent = formattedSeeds;
      const td9 = tr.children[9];
      if (td9 && td9.textContent !== formattedRatio) td9.textContent = formattedRatio;
      const td10 = tr.children[10];
      if (td10 && td10.textContent !== formattedDate) td10.textContent = formattedDate;
    }

    if (t.info_hash === state.selectedHash) {
      tr.classList.add('selected');
    } else {
      tr.classList.remove('selected');
    }

    if (tbody.children[index] !== tr) {
      tbody.insertBefore(tr, tbody.children[index] || null);
    }
  });

  const selectedTorrent = getSelectedTorrent();
  updateToolbarButtons(selectedTorrent);
}

function selectTorrentRow(tr, hash) {
  tbody.querySelectorAll('tr.selected').forEach(r => r.classList.remove('selected'));
  tr.classList.add('selected');
  state.selectedHash = hash;
  const current = state.torrents.find(item => item.info_hash === hash);
  updateToolbarButtons(current);
  fetchTorrentDetails(hash);
}

function updateToolbarButtons(selectedTorrent) {
  const hasSelection = !!selectedTorrent;
  btnResume.disabled = !hasSelection;
  btnPause.disabled = !hasSelection;
  btnDelete.disabled = !hasSelection;
  btnOpenFolder.disabled = !hasSelection;
  btnSeq.disabled = !hasSelection;
  btnVerify.disabled = !hasSelection;
  if (btnRedownload) btnRedownload.disabled = !hasSelection;

  if (hasSelection) {
    if (selectedTorrent.sequential) {
      btnSeq.classList.add('btn-primary');
    } else {
      btnSeq.classList.remove('btn-primary');
    }
  } else {
    btnSeq.classList.remove('btn-primary');
  }
}

function populateTorrentContextMenu(tItem) {
  const iconToggle = document.getElementById('ctx-icon-toggle');
  const labelToggle = document.getElementById('ctx-label-toggle');
  if (tItem.status === 'paused') {
    iconToggle.innerHTML = `<svg viewBox="0 0 24 24" class="ctx-svg"><polygon points="5 3 19 12 5 21 5 3"/></svg>`;
    labelToggle.textContent = t('ctx_resume', 'Возобновить');
  } else {
    iconToggle.innerHTML = `<svg viewBox="0 0 24 24" class="ctx-svg"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>`;
    labelToggle.textContent = t('ctx_pause', 'Приостановить');
  }

  const labelSeq = document.getElementById('ctx-label-sequential');
  if (tItem.sequential) {
    labelSeq.textContent = t('ctx_sequential', 'Последовательно') + t('seq_on', ' (Вкл)');
  } else {
    labelSeq.textContent = t('ctx_sequential', 'Последовательно');
  }
}

// Inspector Details
function clearDetails() {
  const fields = [
    'info-downloaded', 'info-uploaded', 'info-dlspeed', 'info-ulspeed',
    'info-eta', 'info-ratio', 'info-size', 'info-piecesize',
    'info-piececount', 'info-hash', 'info-savepath', 'info-date', 'info-creator'
  ];
  fields.forEach(id => {
    const el = document.getElementById(id);
    if (el) el.textContent = '-';
  });

  const filesCount = document.getElementById('files-count');
  if (filesCount) filesCount.textContent = '0';

  const filesTbody = document.getElementById('files-tbody');
  if (filesTbody) {
    filesTbody.innerHTML = '';
    delete filesTbody.dataset.renderedHash;
  }

  const trackersTbody = document.getElementById('trackers-tbody');
  if (trackersTbody) trackersTbody.innerHTML = '';

  const peersTbody = document.getElementById('peers-tbody');
  if (peersTbody) peersTbody.innerHTML = '';

  const piecesGrid = document.getElementById('pieces-grid');
  if (piecesGrid) piecesGrid.innerHTML = '';

  const piecesStats = document.getElementById('pieces-stats');
  if (piecesStats) piecesStats.textContent = '0 / 0 ' + t('pieces_suffix', 'частей');
}

async function fetchTorrentDetails(hash) {
  try {
    const res = await fetch(`/api/torrent?hash=${encodeURIComponent(hash)}`);
    if (!res.ok) {
      if (res.status === 404 && state.selectedHash === hash) {
        state.selectedHash = null;
        clearDetails();
        updateToolbarButtons(null);
      }
      return;
    }
    const data = await res.json();
    renderTorrentDetails(data);
  } catch (err) {
    console.error("fetchTorrentDetails error:", err);
  }
}

function renderTorrentDetails(detail) {
  if (!detail || !detail.summary) return;
  const s = detail.summary;

  document.getElementById('info-downloaded').textContent = formatBytes(s.downloaded);
  document.getElementById('info-uploaded').textContent = formatBytes(s.uploaded);
  document.getElementById('info-dlspeed').textContent = formatSpeed(s.download_speed);
  document.getElementById('info-ulspeed').textContent = formatSpeed(s.upload_speed);
  document.getElementById('info-eta').textContent = formatETA(s.eta);
  const ratioVal = (typeof s.ratio === 'number' && !isNaN(s.ratio)) ? s.ratio : 0;
  document.getElementById('info-ratio').textContent = ratioVal.toFixed(2);

  document.getElementById('info-size').textContent = formatBytes(s.size);
  document.getElementById('info-piecesize').textContent = formatBytes(detail.piece_length);
  document.getElementById('info-piececount').textContent = detail.piece_count || 0;
  document.getElementById('info-hash').textContent = s.info_hash || '-';
  document.getElementById('info-savepath').textContent = s.save_path || '-';
  document.getElementById('info-date').textContent = formatDate(s.added_at);
  document.getElementById('info-creator').textContent = detail.created_by || t('not_specified', 'Не указано');

  // Files tab
  const files = detail.files || [];
  document.getElementById('files-count').textContent = files.length;
  renderFilesTab(files, s.info_hash);

  // Peers tab
  renderPeersTab(detail.peers || []);

  // Trackers tab
  renderTrackersTab(detail.trackers || []);

  // Pieces map tab
  renderPiecesTab(detail.piece_map || '', detail.piece_count || 0);
}

function renderFilesTab(files, hash) {
  const tbody = document.getElementById('files-tbody');
  if (!tbody) return;

  const currentCount = tbody.children.length;
  const needsRebuild = currentCount !== files.length || tbody.dataset.renderedHash !== hash;

  if (needsRebuild) {
    tbody.innerHTML = '';
    tbody.dataset.renderedHash = hash;
    files.forEach((f, idx) => {
      const tr = document.createElement('tr');
      tr.dataset.index = idx;
      tr.dataset.hash = hash;
      const progVal = (typeof f.progress === 'number' && !isNaN(f.progress)) ? f.progress : 0;
      const progPercent = progVal.toFixed(1);
      const prioBadge = getPriorityBadge(f.priority);
      const fileName = f.path.split('/').pop() || f.path;

      tr.innerHTML = `
        <td class="col-file-name" title="${escapeHtml(f.path)}">
          <span style="display: flex; align-items: center; gap: 6px;">
            <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            <span>${escapeHtml(fileName)}</span>
          </span>
        </td>
        <td class="col-file-size font-mono">${formatBytes(f.length)}</td>
        <td class="col-file-prog font-mono">${progPercent}%</td>
        <td class="col-file-prio">${prioBadge}</td>
      `;

      tr.addEventListener('click', () => {
        tbody.querySelectorAll('tr.selected').forEach(r => r.classList.remove('selected'));
        tr.classList.add('selected');
      });

      tr.addEventListener('dblclick', () => {
        openTargetFile(hash, idx);
      });

      tr.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        tbody.querySelectorAll('tr.selected').forEach(r => r.classList.remove('selected'));
        tr.classList.add('selected');
        state.contextFile = {
          hash: hash,
          fileIndex: idx,
          path: f.path,
          name: fileName
        };
        showContextMenu(ctxFile, e.clientX, e.clientY);
      });

      tbody.appendChild(tr);
    });
  } else {
    // In-place updates for files
    files.forEach((f, idx) => {
      const tr = tbody.children[idx];
      if (!tr) return;
      const progVal = (typeof f.progress === 'number' && !isNaN(f.progress)) ? f.progress : 0;
      const progPercent = progVal.toFixed(1) + '%';
      const tdProg = tr.children[2];
      if (tdProg && tdProg.textContent !== progPercent) tdProg.textContent = progPercent;

      const prioBadge = getPriorityBadge(f.priority);
      const tdPrio = tr.children[3];
      if (tdPrio && tdPrio.innerHTML !== prioBadge) tdPrio.innerHTML = prioBadge;
    });
  }
}

function getPriorityBadge(prio) {
  switch (prio) {
    case 2:
      return `<span class="prio-badge prio-high">${escapeHtml(t('prio_high', 'Высокий'))}</span>`;
    case 0:
      return `<span class="prio-badge prio-none">${escapeHtml(t('prio_skip', 'Пропущен'))}</span>`;
    default:
      return `<span class="prio-badge prio-normal">${escapeHtml(t('prio_normal', 'Обычный'))}</span>`;
  }
}

function renderPeersTab(peers) {
  const tbody = document.getElementById('peers-tbody');
  if (!tbody) return;
  tbody.innerHTML = '';
  peers.forEach(p => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td class="font-mono">${escapeHtml(p.address)}</td>
      <td>${escapeHtml(p.client || t('unknown', 'Неизвестно'))}</td>
      <td class="font-mono">${formatSpeed(p.download_rate)}</td>
      <td class="font-mono">${formatSpeed(p.upload_rate)}</td>
      <td class="font-mono">${(p.progress || 0).toFixed(1)}%</td>
    `;
    tbody.appendChild(tr);
  });
}

function renderTrackersTab(trackers) {
  const tbody = document.getElementById('trackers-tbody');
  if (!tbody) return;
  tbody.innerHTML = '';
  trackers.forEach((t, i) => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td class="font-mono">${i + 1}</td>
      <td class="font-mono">${escapeHtml(t.url)}</td>
      <td><span class="status-pill status-ready">${escapeHtml(t.status || t('working', 'Работает'))}</span></td>
    `;
    tbody.appendChild(tr);
  });
}

function renderPiecesTab(pieceMap, totalPieces) {
  const grid = document.getElementById('pieces-grid');
  const stats = document.getElementById('pieces-stats');
  if (!grid || !stats) return;

  grid.innerHTML = '';
  if (!pieceMap || totalPieces <= 0) {
    stats.textContent = '0 / 0 ' + t('pieces_suffix', 'частей');
    return;
  }

  // Base64 decode piece states
  let rawBytes;
  try {
    rawBytes = Uint8Array.from(atob(pieceMap), c => c.charCodeAt(0));
  } catch {
    rawBytes = new Uint8Array(0);
  }

  let completed = 0;
  for (let i = 0; i < rawBytes.length; i++) {
    if (rawBytes[i] === 2) completed++;
  }

  stats.textContent = `${completed} / ${totalPieces} ${t('pieces_suffix', 'частей')} (${((completed / totalPieces) * 100).toFixed(1)}%)`;

  const maxDisplay = 1200; // Limit rendering for DOM performance with massive torrents
  const step = Math.max(1, Math.ceil(rawBytes.length / maxDisplay));

  for (let i = 0; i < rawBytes.length; i += step) {
    const div = document.createElement('div');
    div.className = 'piece-block';
    const val = rawBytes[i];
    if (val === 2) {
      div.classList.add('done');
    } else if (val === 1) {
      div.classList.add('part');
    }
    grid.appendChild(div);
  }
}

// API Action Handlers
async function addTorrentFile() {
  try {
    const res = await fetch('/api/add-file-dialog', { method: 'POST' });
    const data = await res.json();
    if (data.success && data.torrent) {
      state.selectedHash = data.torrent.info_hash;
      showToast(t('toast_torrent_added', 'Торрент добавлен: ') + data.torrent.name);
      await fetchTorrents();
    } else if (data.error && data.error !== 'cancelled') {
      showToast(t("toast_torrent_add_error", "Ошибка добавления торрента: ") + data.error);
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function uploadTorrentFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  try {
    const res = await fetch('/api/upload-torrent', {
      method: 'POST',
      body: formData
    });
    const data = await res.json();
    if (data.success && data.torrent) {
      state.selectedHash = data.torrent.info_hash;
      showToast(t('toast_torrent_added', 'Торрент добавлен: ') + data.torrent.name);
      await fetchTorrents();
    } else {
      showToast(t("toast_error", "Ошибка: ") + (data.error || t("toast_upload_fail", "Не удалось загрузить файл")));
    }
  } catch (err) {
    showToast(t("toast_upload_error", "Ошибка загрузки: ") + err.message);
  }
}

async function submitMagnet() {
  const input = document.getElementById('magnet-input');
  const uri = input.value.trim();
  if (!uri || !uri.toLowerCase().startsWith('magnet:?')) {
    alert(t("alert_invalid_magnet", "Пожалуйста, введите корректную Magnet-ссылку (начинается с magnet:?)"));
    return;
  }
  await addMagnetURI(uri);
  closeModal('modal-magnet');
  input.value = '';
}

async function addMagnetURI(uri) {
  try {
    const res = await fetch('/api/add-magnet', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ magnet_uri: uri })
    });
    const data = await res.json();
    if (data.success && data.torrent) {
      state.selectedHash = data.torrent.info_hash;
      showToast(t("toast_magnet_added", "Magnet-ссылка успешно добавлена!"));
      await fetchTorrents();
    } else {
      showToast(t("toast_magnet_add_error", "Ошибка добавления Magnet: ") + (data.error || t("unknown", "Неизвестная ошибка")));
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function performAction(action) {
  if (!state.selectedHash) return;
  try {
    const res = await fetch(`/api/torrent/${action}?hash=${encodeURIComponent(state.selectedHash)}`, {
      method: 'POST'
    });
    const data = await res.json();
    if (data.success) {
      await fetchTorrents();
      if (state.selectedHash) await fetchTorrentDetails(state.selectedHash);
    } else {
      showToast(t("toast_error", "Ошибка: ") + (data.error || t("toast_action_fail", "Не удалось выполнить действие")));
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function toggleSequential() {
  const t = getSelectedTorrent();
  if (!t) return;
  const newSeq = !t.sequential;
  try {
    const res = await fetch(`/api/torrent/sequential?hash=${encodeURIComponent(t.info_hash)}&enabled=${newSeq}`, {
      method: 'POST'
    });
    const data = await res.json();
    if (data.success) {
      t.sequential = newSeq;
      updateToolbarButtons(t);
      showToast(newSeq ? t("toast_seq_enabled", "Последовательная загрузка включена") : t("toast_seq_disabled", "Последовательная загрузка выключена"));
      await fetchTorrents();
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function verifyTorrent() {
  if (!state.selectedHash) return;
  try {
    const res = await fetch(`/api/torrent/verify?hash=${encodeURIComponent(state.selectedHash)}`, {
      method: 'POST'
    });
    const data = await res.json();
    if (data.success) {
      showToast(t("toast_verify_started", "Запущена проверка целостности данных"));
      await fetchTorrents();
    } else {
      showToast(t("toast_error", "Ошибка: ") + data.error);
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function redownloadTorrent() {
  if (!state.selectedHash) return;
  try {
    const res = await fetch(`/api/torrent/redownload?hash=${encodeURIComponent(state.selectedHash)}`, {
      method: 'POST'
    });
    const data = await res.json();
    if (data.success) {
      showToast(t("toast_redownload_started", "Торрент перезапущен для загрузки заново"));
      await fetchTorrents();
    } else {
      showToast(t("toast_error", "Ошибка: ") + data.error);
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function openFolder() {
  if (!state.selectedHash) return;
  try {
    const res = await fetch(`/api/open-folder?hash=${encodeURIComponent(state.selectedHash)}`);
    const data = await res.json();
    if (!data.success) {
      showToast(t("toast_open_folder_error", "Не удалось открыть папку: ") + (data.error || t("toast_folder_not_found", "Папка не найдена")));
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function openTargetFile(hash, fileIndex) {
  try {
    const res = await fetch(`/api/torrent/open-file?hash=${encodeURIComponent(hash)}&file_index=${fileIndex}`);
    const data = await res.json();
    if (!data.success) {
      showToast(t("toast_open_file_error", "Не удалось открыть файл: ") + (data.error || t("toast_file_not_found", "Файл не найден на диске")));
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function showFileInFolder(hash, fileIndex) {
  try {
    const res = await fetch(`/api/torrent/show-file-in-folder?hash=${encodeURIComponent(hash)}&file_index=${fileIndex}`);
    const data = await res.json();
    if (!data.success) {
      showToast(t("toast_open_folder_error", "Не удалось открыть папку: ") + (data.error || t("toast_explorer_error", "Ошибка проводника")));
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

async function setFilePriority(hash, fileIndex, priority) {
  try {
    const res = await fetch('/api/torrent/file-priority', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ hash, file_index: fileIndex, priority })
    });
    const data = await res.json();
    if (data.success) {
      showToast(t("toast_prio_updated", "Приоритет файла обновлен"));
      if (state.selectedHash) await fetchTorrentDetails(state.selectedHash);
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

// Delete Modal
function openDeleteModal() {
  const t = getSelectedTorrent();
  if (!t) return;
  document.getElementById('delete-torrent-name').textContent = t.name;
  document.getElementById('delete-files-checkbox').checked = false;
  openModal('modal-delete');
}

async function confirmDelete() {
  const t = getSelectedTorrent();
  if (!t) return;
  const deleteFiles = document.getElementById('delete-files-checkbox').checked;
  closeModal('modal-delete');

  try {
    const res = await fetch(`/api/torrent/delete?hash=${encodeURIComponent(t.info_hash)}&deleteFiles=${deleteFiles}`, {
      method: 'POST'
    });
    const data = await res.json();
    if (data.success) {
      state.selectedHash = null;
      clearDetails();
      updateToolbarButtons(null);
      showToast(t("toast_torrent_deleted", "Торрент удален"));
      await fetchTorrents();
    } else {
      showToast(t("toast_delete_error", "Ошибка удаления: ") + data.error);
    }
  } catch (err) {
    showToast(t("toast_error", "Ошибка: ") + err.message);
  }
}

// Settings Modal
async function openSettings() {
  try {
    const res = await fetch('/api/settings');
    if (!res.ok) return;
    const cfg = await res.json();

    state.defaultDownloadDir = cfg.default_download_dir || '';
    document.getElementById('setting-download-dir').value = cfg.download_dir || '';
    document.getElementById('setting-port').value = cfg.listen_port || 0;
    document.getElementById('setting-max-conns').value = cfg.max_established || 80;
    document.getElementById('setting-dl-limit').value = cfg.download_limit_bps > 0 ? Math.round(cfg.download_limit_bps / 1024) : 0;
    document.getElementById('setting-ul-limit').value = cfg.upload_limit_bps > 0 ? Math.round(cfg.upload_limit_bps / 1024) : 0;

    document.getElementById('setting-dht').checked = !!cfg.enable_dht;
    document.getElementById('setting-pex').checked = !!cfg.enable_pex;
    document.getElementById('setting-upnp').checked = !!cfg.enable_upnp;

    document.getElementById('setting-assoc-torrent').checked = !!cfg.associate_torrents;
    document.getElementById('setting-assoc-magnet').checked = !!cfg.associate_magnets;
    document.getElementById('setting-autostart').checked = !!cfg.autostart;

    document.getElementById('setting-minimize-tray').checked = !!cfg.minimize_to_tray;
    document.getElementById('setting-close-tray').checked = !!cfg.close_to_tray;
    document.getElementById('setting-start-minimized').checked = !!cfg.start_minimized;

    const themeEl = document.getElementById('setting-theme');
    if (themeEl) themeEl.value = cfg.theme || state.theme || 'dark';

    const langEl = document.getElementById('setting-language');
    if (langEl) langEl.value = (cfg.language && i18n[cfg.language]) ? cfg.language : (i18n[state.language] ? state.language : 'ru');
    syncCustomLanguageSelect(langEl ? langEl.value : state.language);

    state.settingsInitialTheme = themeEl ? themeEl.value : state.theme;
    state.settingsInitialLanguage = langEl ? langEl.value : state.language;

    openModal('modal-settings');
  } catch (err) {
    showToast(t("toast_load_error", "Не удалось загрузить настройки: ") + err.message);
  }
}

async function browseDownloadDir() {
  try {
    const res = await fetch('/api/browse-dir', { method: 'POST' });
    const data = await res.json();
    if (data.path) {
      document.getElementById('setting-download-dir').value = data.path;
    }
  } catch (err) {
    showToast(t("toast_browse_dir_error", "Ошибка выбора папки: ") + err.message);
  }
}

async function saveSettings() {
  const dlLimitKb = parseInt(document.getElementById('setting-dl-limit').value, 10) || 0;
  const ulLimitKb = parseInt(document.getElementById('setting-ul-limit').value, 10) || 0;

  const themeEl = document.getElementById('setting-theme');
  const langEl = document.getElementById('setting-language');
  const chosenTheme = themeEl ? themeEl.value : state.theme;
  const chosenLang = langEl ? langEl.value : state.language;

  const cfg = {
    download_dir: document.getElementById('setting-download-dir').value.trim(),
    listen_port: parseInt(document.getElementById('setting-port').value, 10) || 0,
    max_established: parseInt(document.getElementById('setting-max-conns').value, 10) || 80,
    download_limit_bps: dlLimitKb * 1024,
    upload_limit_bps: ulLimitKb * 1024,
    enable_dht: document.getElementById('setting-dht').checked,
    enable_pex: document.getElementById('setting-pex').checked,
    enable_upnp: document.getElementById('setting-upnp').checked,
    associate_torrents: document.getElementById('setting-assoc-torrent').checked,
    associate_magnets: document.getElementById('setting-assoc-magnet').checked,
    autostart: document.getElementById('setting-autostart').checked,
    minimize_to_tray: document.getElementById('setting-minimize-tray').checked,
    close_to_tray: document.getElementById('setting-close-tray').checked,
    start_minimized: document.getElementById('setting-start-minimized').checked,
    theme: chosenTheme,
    language: chosenLang
  };

  try {
    const res = await fetch('/api/settings', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(cfg)
    });
    const data = await res.json();
    if (data.success) {
      state.settingsInitialTheme = null;
      state.settingsInitialLanguage = null;
      applyTheme(chosenTheme);
      applyLanguage(chosenLang);
      closeModal('modal-settings');
      showToast(t('toast_settings_saved', "Настройки успешно сохранены"));
    } else {
      showToast(t('toast_settings_error', "Ошибка сохранения: ") + data.error);
    }
  } catch (err) {
    showToast(t('toast_settings_error', "Ошибка сохранения: ") + err.message);
  }
}

// Modal Helpers
function openModal(id) {
  const modal = document.getElementById(id);
  if (modal) modal.classList.remove('hidden');
}

function closeModal(id) {
  const modal = document.getElementById(id);
  if (modal) modal.classList.add('hidden');
  if (id === 'modal-settings') {
    if (state.settingsInitialTheme && state.settingsInitialTheme !== state.theme) {
      applyTheme(state.settingsInitialTheme);
    }
    if (state.settingsInitialLanguage && state.settingsInitialLanguage !== state.language) {
      applyLanguage(state.settingsInitialLanguage);
    }
    state.settingsInitialTheme = null;
    state.settingsInitialLanguage = null;
  }
}
