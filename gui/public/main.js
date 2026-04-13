// SmirnovTorrent Desktop GUI - Main JavaScript

let downloads = new Map();
let updateInterval = null;

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
  log('SmirnovTorrent Desktop GUI initialized', 'info');
  loadDownloads();
  startAutoUpdate();
});

// Browse for torrent file
async function browseFile() {
  try {
    const { open } = await import('@tauri-apps/api/dialog');
    const selected = await open({
      multiple: false,
      filters: [{
        name: 'Torrent',
        extensions: ['torrent']
      }]
    });
    
    if (selected) {
      document.getElementById('torrentPath').value = selected;
      log(`Selected file: ${selected}`, 'info');
    }
  } catch (error) {
    log(`Error browsing file: ${error}`, 'error');
  }
}

// Add torrent
async function addTorrent() {
  const path = document.getElementById('torrentPath').value.trim();
  
  if (!path) {
    log('Please enter a torrent path or magnet link', 'warning');
    return;
  }

  try {
    const { invoke } = await import('@tauri-apps/api/tauri');
    
    log(`Adding torrent: ${path}`, 'info');
    
    // Invoke Rust backend to add torrent
    const downloadId = await invoke('add_torrent', { path });
    
    downloads.set(downloadId, {
      id: downloadId,
      name: path.split('/').pop() || 'Unknown',
      path: path,
      progress: 0,
      downloaded: 0,
      uploaded: 0,
      downloadSpeed: 0,
      uploadSpeed: 0,
      peers: 0,
      status: 'downloading',
      totalSize: 0
    });
    
    log(`Torrent added successfully: ${downloadId}`, 'success');
    renderDownloads();
    document.getElementById('torrentPath').value = '';
    
  } catch (error) {
    log(`Error adding torrent: ${error}`, 'error');
  }
}

// Pause download
async function pauseDownload(id) {
  try {
    const { invoke } = await import('@tauri-apps/api/tauri');
    await invoke('pause_download', { id });
    
    const download = downloads.get(id);
    if (download) {
      download.status = 'paused';
      renderDownloads();
      log(`Download paused: ${id}`, 'info');
    }
  } catch (error) {
    log(`Error pausing download: ${error}`, 'error');
  }
}

// Resume download
async function resumeDownload(id) {
  try {
    const { invoke } = await import('@tauri-apps/api/tauri');
    await invoke('resume_download', { id });
    
    const download = downloads.get(id);
    if (download) {
      download.status = 'downloading';
      renderDownloads();
      log(`Download resumed: ${id}`, 'success');
    }
  } catch (error) {
    log(`Error resuming download: ${error}`, 'error');
  }
}

// Remove download
async function removeDownload(id) {
  try {
    const { invoke } = await import('@tauri-apps/api/tauri');
    await invoke('remove_download', { id });
    
    downloads.delete(id);
    renderDownloads();
    log(`Download removed: ${id}`, 'info');
  } catch (error) {
    log(`Error removing download: ${error}`, 'error');
  }
}

// Render downloads list
function renderDownloads() {
  const container = document.getElementById('downloadsList');
  
  if (downloads.size === 0) {
    container.innerHTML = '<div class="empty-state">No active downloads</div>';
    updateStats();
    return;
  }
  
  container.innerHTML = Array.from(downloads.values()).map(dl => `
    <div class="download-item" data-id="${dl.id}">
      <div class="download-header">
        <div class="download-name">${escapeHtml(dl.name)}</div>
        <div class="download-status status-${dl.status}">${dl.status.toUpperCase()}</div>
      </div>
      
      <div class="progress-container">
        <div class="progress-bar">
          <div class="progress-fill" style="width: ${dl.progress}%"></div>
        </div>
        <div class="progress-text">${dl.progress.toFixed(1)}% complete</div>
      </div>
      
      <div class="download-info">
        <div class="info-item">
          <div class="info-label">Download Speed</div>
          <div class="info-value">${formatSpeed(dl.downloadSpeed)}</div>
        </div>
        <div class="info-item">
          <div class="info-label">Upload Speed</div>
          <div class="info-value">${formatSpeed(dl.uploadSpeed)}</div>
        </div>
        <div class="info-item">
          <div class="info-label">Peers</div>
          <div class="info-value">${dl.peers}</div>
        </div>
        <div class="info-item">
          <div class="info-label">Downloaded</div>
          <div class="info-value">${formatBytes(dl.downloaded)}</div>
        </div>
        <div class="info-item">
          <div class="info-label">Uploaded</div>
          <div class="info-value">${formatBytes(dl.uploaded)}</div>
        </div>
        <div class="info-item">
          <div class="info-label">Size</div>
          <div class="info-value">${formatBytes(dl.totalSize)}</div>
        </div>
      </div>
      
      <div class="download-actions">
        ${dl.status === 'downloading' 
          ? `<button class="btn-pause" onclick="pauseDownload('${dl.id}')">Pause</button>`
          : `<button class="btn-resume" onclick="resumeDownload('${dl.id}')">Resume</button>`
        }
        <button class="btn-remove" onclick="removeDownload('${dl.id}')">Remove</button>
      </div>
    </div>
  `).join('');
  
  updateStats();
}

// Update statistics
function updateStats() {
  const totalDownloadSpeed = Array.from(downloads.values())
    .reduce((sum, dl) => sum + dl.downloadSpeed, 0);
  const totalUploadSpeed = Array.from(downloads.values())
    .reduce((sum, dl) => sum + dl.uploadSpeed, 0);
  const totalPeers = Array.from(downloads.values())
    .reduce((sum, dl) => sum + dl.peers, 0);
  const totalDownloaded = Array.from(downloads.values())
    .reduce((sum, dl) => sum + dl.downloaded, 0);
  
  document.getElementById('downloadSpeed').textContent = formatSpeed(totalDownloadSpeed);
  document.getElementById('uploadSpeed').textContent = formatSpeed(totalUploadSpeed);
  document.getElementById('activePeers').textContent = totalPeers;
  document.getElementById('totalDownloaded').textContent = formatBytes(totalDownloaded);
}

// Auto-update downloads
async function startAutoUpdate() {
  try {
    const { invoke } = await import('@tauri-apps/api/tauri');
    
    updateInterval = setInterval(async () => {
      try {
        const updates = await invoke('get_downloads_status');
        
        if (updates && Array.isArray(updates)) {
          updates.forEach(update => {
            if (downloads.has(update.id)) {
              const dl = downloads.get(update.id);
              Object.assign(dl, update);
            }
          });
          renderDownloads();
        }
      } catch (error) {
        // Silent error for auto-update
      }
    }, 2000); // Update every 2 seconds
  } catch (error) {
    log(`Error starting auto-update: ${error}`, 'error');
  }
}

// Load downloads from backend
async function loadDownloads() {
  try {
    const { invoke } = await import('@tauri-apps/api/tauri');
    const downloadsList = await invoke('get_downloads');
    
    if (downloadsList && Array.isArray(downloadsList)) {
      downloadsList.forEach(dl => {
        downloads.set(dl.id, dl);
      });
      renderDownloads();
      log(`Loaded ${downloads.size} downloads`, 'info');
    }
  } catch (error) {
    log(`Error loading downloads: ${error}`, 'error');
  }
}

// Logging
function log(message, type = 'info') {
  const logOutput = document.getElementById('logOutput');
  const timestamp = new Date().toLocaleTimeString();
  const entry = document.createElement('div');
  entry.className = `log-entry log-${type}`;
  entry.textContent = `[${timestamp}] ${message}`;
  logOutput.appendChild(entry);
  logOutput.scrollTop = logOutput.scrollHeight;
}

function clearLog() {
  document.getElementById('logOutput').innerHTML = '';
  log('Log cleared', 'info');
}

// Utility functions
function formatSpeed(bytesPerSecond) {
  if (bytesPerSecond < 1024) {
    return `${Math.round(bytesPerSecond)} B/s`;
  }
  if (bytesPerSecond < 1024 * 1024) {
    return `${(bytesPerSecond / 1024).toFixed(1)} KB/s`;
  }
  return `${(bytesPerSecond / (1024 * 1024)).toFixed(1)} MB/s`;
}

function formatBytes(bytes) {
  if (bytes < 1024) {
    return `${Math.round(bytes)} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  if (bytes < 1024 * 1024 * 1024) {
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// Cleanup
window.addEventListener('beforeunload', () => {
  if (updateInterval) {
    clearInterval(updateInterval);
  }
});
