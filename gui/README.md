# SmirnovTorrent Desktop GUI

Modern desktop application for SmirnovTorrent built with [Tauri](https://tauri.app/).

## 🚀 Features

- 🎨 **Modern UI** - Clean, dark theme interface
- 📊 **Real-time Stats** - Live download/upload speeds, peer count
- 🎯 **Progress Tracking** - Visual progress bars for each download
- ⏯️ **Control Downloads** - Pause, resume, and remove downloads
- 📝 **Activity Log** - Real-time logging of all actions
- 🔍 **File Browser** - Easy torrent file selection
- 🧲 **Magnet Support** - Add torrents via magnet links
- ⚡ **Fast & Lightweight** - Native performance with minimal resource usage

## 📋 Prerequisites

### Node.js & npm
- Node.js >= 18.0.0
- npm >= 9.0.0

### Rust & Cargo
- Rust >= 1.70.0
- Cargo >= 1.70.0

### System Dependencies

#### Windows
- [WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) (usually pre-installed)
- Visual Studio C++ Build Tools

#### macOS
- Xcode Command Line Tools
- Xcode 14 or later

#### Linux
```bash
# Debian/Ubuntu
sudo apt update
sudo apt install -y libwebkit2gtk-4.0-dev \
    build-essential \
    curl \
    wget \
    libssl-dev \
    libgtk-3-dev \
    libayatana-appindicator3-dev \
    librsvg2-dev

# Fedora
sudo dnf install -y webkit2gtk3-devel \
    openssl-devel \
    curl \
    wget \
    libappindicator-gtk3-devel \
    librsvg2-devel
```

## 🛠️ Installation

### 1. Navigate to GUI directory
```bash
cd gui
```

### 2. Install Node.js dependencies
```bash
npm install
```

### 3. Install Rust dependencies
```bash
cargo install
```

## 🚀 Development

### Start development server
```bash
npm run tauri dev
```

This will:
1. Start the Vite dev server on port 3000
2. Launch the Tauri application window
3. Enable hot-reload for frontend changes

## 📦 Building

### Build for current platform
```bash
npm run tauri build
```

### Build specific targets
```bash
# Windows
npm run tauri build -- --target x86_64-pc-windows-msvc

# macOS
npm run tauri build -- --target x86_64-apple-darwin

# Linux
npm run tauri build -- --target x86_64-unknown-linux-gnu
```

Build artifacts will be in `gui/src-tauri/target/release/bundle/`

## 📁 Project Structure

```
gui/
├── src-tauri/              # Rust backend
│   ├── src/
│   │   └── main.rs        # Tauri commands & app logic
│   ├── icons/             # Application icons
│   ├── Cargo.toml         # Rust dependencies
│   ├── build.rs           # Build script
│   └── tauri.conf.json    # Tauri configuration
├── public/
│   ├── index.html         # Main HTML (injected from root)
│   ├── style.css          # Application styles
│   └── main.js            # Frontend JavaScript
├── index.html             # Main HTML file
├── package.json           # Node.js dependencies
└── vite.config.js         # Vite configuration
```

## 🎮 Usage

### Adding a Torrent

1. **Via File Browser:**
   - Click "Browse" button
   - Select a `.torrent` file
   - Click "Add"

2. **Via Path:**
   - Enter path to `.torrent` file
   - Click "Add"

3. **Via Magnet Link:**
   - Paste magnet link
   - Click "Add"

### Managing Downloads

- **Pause:** Click "Pause" button to temporarily stop download
- **Resume:** Click "Resume" button to continue paused download
- **Remove:** Click "Remove" button to delete download

### Monitoring Progress

The dashboard shows:
- Download speed (total & per-download)
- Upload speed (total & per-download)
- Active peer count
- Total downloaded data
- Individual progress bars

## 🔧 Configuration

Edit `gui/src-tauri/tauri.conf.json` to customize:

- Window size and title
- Application bundle settings
- Security options
- Build configuration

## 🐛 Troubleshooting

### Build fails on Linux
```bash
# Install missing dependencies
sudo apt install libwebkit2gtk-4.0-dev libgtk-3-dev
```

### WebView2 error on Windows
Download and install: https://developer.microsoft.com/en-us/microsoft-edge/webview2/

### Rust compilation errors
```bash
# Update Rust
rustup update

# Clear cargo cache
cargo clean
```

## 📊 Screenshots

### Main Dashboard
- Modern dark theme
- Real-time statistics
- Active downloads list
- Activity log

## 🔗 Links

- [SmirnovTorrent Main Project](https://github.com/ainsonet/smirnovtorrent)
- [Tauri Documentation](https://tauri.app/docs)
- [Vite Documentation](https://vitejs.dev/)

## 📄 License

MIT License - same as main SmirnovTorrent project

---

**Built with ❤️ using Tauri + Vite + Rust**
