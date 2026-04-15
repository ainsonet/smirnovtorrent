# Setup Logo

To display the logo on your website, copy your logo file to this directory:

```
Copy from: C:\Users\user\Documents\Visual Studio Code\SmirnovTorrent\logo2.png
Copy to:   website\logo2.png
```

Or use this command in PowerShell:

```powershell
Copy-Item "C:\Users\user\Documents\Visual Studio Code\SmirnovTorrent\logo2.png" website\logo2.png
```

The logo will appear:
- In the navigation bar (40x40px)
- In the hero section (120x120px)
- In the footer (40x40px)

If the logo file is missing, it will be hidden automatically (no broken image icon).
