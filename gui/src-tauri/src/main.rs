// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use uuid::Uuid;
use std::process::{Command, Stdio};
use std::path::PathBuf;

// Download status structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DownloadStatus {
    pub id: String,
    pub name: String,
    pub path: String,
    pub progress: f64,
    pub downloaded: u64,
    pub uploaded: u64,
    pub download_speed: u64,
    pub upload_speed: u64,
    pub peers: u32,
    pub status: String,
    pub total_size: u64,
}

// Shared state for downloads
pub struct AppState {
    pub downloads: Arc<Mutex<HashMap<String, DownloadStatus>>>,
}

// Command: Add a new torrent
#[tauri::command]
fn add_torrent(
    state: tauri::State<AppState>,
    path: String,
) -> Result<String, String> {
    let id = Uuid::new_v4().to_string();
    let name = path.split('\\').last().unwrap_or(&path).split('/').last().unwrap_or("Unknown").to_string();
    
    let download = DownloadStatus {
        id: id.clone(),
        name: name.clone(),
        path: path.clone(),
        progress: 0.0,
        downloaded: 0,
        uploaded: 0,
        download_speed: 0,
        upload_speed: 0,
        peers: 0,
        status: "downloading".to_string(),
        total_size: 0,
    };
    
    let mut downloads = state.downloads.lock().map_err(|e| e.to_string())?;
    downloads.insert(id.clone(), download);
    
    println!("Added torrent: {} - {}", id, name);
    
    // В реальной реализации здесь будет вызов Go движка
    // Для демо просто создаём запись
    Ok(id)
}

// Command: Pause download
#[tauri::command]
fn pause_download(
    state: tauri::State<AppState>,
    id: String,
) -> Result<(), String> {
    let mut downloads = state.downloads.lock().map_err(|e| e.to_string())?;
    
    if let Some(download) = downloads.get_mut(&id) {
        download.status = "paused".to_string();
        println!("Paused download: {}", id);
        Ok(())
    } else {
        Err(format!("Download {} not found", id))
    }
}

// Command: Resume download
#[tauri::command]
fn resume_download(
    state: tauri::State<AppState>,
    id: String,
) -> Result<(), String> {
    let mut downloads = state.downloads.lock().map_err(|e| e.to_string())?;
    
    if let Some(download) = downloads.get_mut(&id) {
        download.status = "downloading".to_string();
        println!("Resumed download: {}", id);
        Ok(())
    } else {
        Err(format!("Download {} not found", id))
    }
}

// Command: Remove download
#[tauri::command]
fn remove_download(
    state: tauri::State<AppState>,
    id: String,
) -> Result<(), String> {
    let mut downloads = state.downloads.lock().map_err(|e| e.to_string())?;
    downloads.remove(&id);
    println!("Removed download: {}", id);
    Ok(())
}

// Command: Get all downloads
#[tauri::command]
fn get_downloads(
    state: tauri::State<AppState>,
) -> Result<Vec<DownloadStatus>, String> {
    let downloads = state.downloads.lock().map_err(|e| e.to_string())?;
    Ok(downloads.values().cloned().collect())
}

// Command: Get downloads status (for auto-update)
#[tauri::command]
fn get_downloads_status(
    state: tauri::State<AppState>,
) -> Result<Vec<DownloadStatus>, String> {
    let downloads = state.downloads.lock().map_err(|e| e.to_string())?;
    
    // Симуляция прогресса загрузки (в реальной версии - опрос Go движка)
    let mut status_list: Vec<DownloadStatus> = downloads.values().cloned().collect();
    
    for download in &mut status_list {
        if download.status == "downloading" {
            // Симуляция прогресса
            download.progress = (download.progress + 0.1).min(100.0);
            download.downloaded = (download.total_size as f64 * download.progress / 100.0) as u64;
            download.download_speed = 1024 * 1024; // 1 MB/s
            download.upload_speed = 512 * 1024; // 512 KB/s
            download.peers = 15;
            download.total_size = 1024 * 1024 * 100; // 100 MB для демо
        }
    }
    
    Ok(status_list)
}

// Command: Get app version
#[tauri::command]
fn get_version() -> String {
    "1.0.0".to_string()
}

// Command: Open download folder
#[tauri::command]
fn open_download_folder() -> Result<(), String> {
    let downloads_dir = dirs::download_dir()
        .ok_or_else(|| "Could not find downloads directory".to_string())?;
    
    #[cfg(target_os = "windows")]
    Command::new("explorer")
        .arg(downloads_dir)
        .spawn()
        .map_err(|e| e.to_string())?;
    
    #[cfg(target_os = "macos")]
    Command::new("open")
        .arg(downloads_dir)
        .spawn()
        .map_err(|e| e.to_string())?;
    
    #[cfg(target_os = "linux")]
    Command::new("xdg-open")
        .arg(downloads_dir)
        .spawn()
        .map_err(|e| e.to_string())?;
    
    Ok(())
}

fn main() {
    tauri::Builder::default()
        .manage(AppState {
            downloads: Arc::new(Mutex::new(HashMap::new())),
        })
        .invoke_handler(tauri::generate_handler![
            add_torrent,
            pause_download,
            resume_download,
            remove_download,
            get_downloads,
            get_downloads_status,
            get_version,
            open_download_folder
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
