@echo off
echo ===========================================
echo    Building SmirnovTorrent for Windows
echo ===========================================

where rsrc >nul 2>&1
if %ERRORLEVEL% equ 0 (
    if exist smirnovtorrent.ico (
        echo [INFO] Updating Windows executable icon resources...
        rsrc -ico smirnovtorrent.ico -arch amd64 -o rsrc_windows_amd64.syso
        rsrc -ico smirnovtorrent.ico -arch 386 -o rsrc_windows_386.syso
        rsrc -ico smirnovtorrent.ico -arch arm64 -o rsrc_windows_arm64.syso
    )
)

go build -ldflags "-H windowsgui -s -w" -o SmirnovTorrent.exe .
if %ERRORLEVEL% equ 0 (
    echo [OK] Build successful! Created SmirnovTorrent.exe
) else (
    echo [ERROR] Build failed.
)
pause
