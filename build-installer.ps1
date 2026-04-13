# SmirnovTorrent Windows Installer Build Script
# Creates .exe installer using Tauri + NSIS

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "SmirnovTorrent Installer Builder" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Check prerequisites
Write-Host "[1/5] Checking prerequisites..." -ForegroundColor Yellow

# Check Node.js
try {
    $nodeVersion = node --version
    Write-Host "  ✓ Node.js: $nodeVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Node.js not found. Please install from https://nodejs.org/" -ForegroundColor Red
    exit 1
}

# Check npm
try {
    $npmVersion = npm --version
    Write-Host "  ✓ npm: v$npmVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ npm not found" -ForegroundColor Red
    exit 1
}

# Check Rust
try {
    $rustVersion = rustc --version
    Write-Host "  ✓ Rust: $rustVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Rust not found. Please install from https://rustup.rs/" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "[2/5] Installing dependencies..." -ForegroundColor Yellow

# Navigate to gui directory
$guiDir = Join-Path $PSScriptRoot "gui"
Set-Location $guiDir

# Install npm packages
Write-Host "  Installing npm packages..." -ForegroundColor Gray
npm install

Write-Host ""
Write-Host "[3/5] Building frontend..." -ForegroundColor Yellow
npm run build

Write-Host ""
Write-Host "[4/5] Building Tauri application..." -ForegroundColor Yellow
Write-Host "  This may take several minutes on first build..." -ForegroundColor Gray

# Build Tauri with installer
npm run tauri build

Write-Host ""
Write-Host "[5/5] Installer created!" -ForegroundColor Green
Write-Host ""

# Find the installer
$bundleDir = Join-Path $guiDir "src-tauri\target\release\bundle"

if (Test-Path (Join-Path $bundleDir "nsis")) {
    $nsisDir = Join-Path $bundleDir "nsis"
    $installer = Get-ChildItem -Path $nsisDir -Filter "*.exe" | Select-Object -First 1
    
    if ($installer) {
        Write-Host "============================================" -ForegroundColor Cyan
        Write-Host "SUCCESS!" -ForegroundColor Green
        Write-Host "============================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Installer location:" -ForegroundColor Yellow
        Write-Host "  $($installer.FullName)" -ForegroundColor White
        Write-Host ""
        Write-Host "File size: $([math]::Round($installer.Length / 1MB, 2)) MB" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "To install:" -ForegroundColor Yellow
        Write-Host "  1. Double-click the installer" -ForegroundColor White
        Write-Host "  2. Follow the installation wizard" -ForegroundColor White
        Write-Host "  3. Launch SmirnovTorrent from Start Menu or Desktop" -ForegroundColor White
        Write-Host ""
        Write-Host "============================================" -ForegroundColor Cyan
    }
}

if (Test-Path (Join-Path $bundleDir "msi")) {
    $msiDir = Join-Path $bundleDir "msi"
    $msi = Get-ChildItem -Path $msiDir -Filter "*.msi" | Select-Object -First 1
    
    if ($msi) {
        Write-Host "MSI Installer:" -ForegroundColor Yellow
        Write-Host "  $($msi.FullName)" -ForegroundColor White
        Write-Host ""
    }
}

Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
