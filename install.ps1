param(
    [ValidateSet("install", "update", "binary", "extension")]
    [string]$Command = "install"
)

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BinaryName = "butter.exe"
$BinaryPath = Join-Path $ScriptDir $BinaryName
$ExtensionDir = Join-Path $ScriptDir "butter-extension"
$VsixOutput = Join-Path $ScriptDir "butter-extension.vsix"

function Write-Info  { Write-Host "[INFO] $args" -ForegroundColor Green }
function Write-Warn  { Write-Host "[WARN] $args" -ForegroundColor Yellow }
function Write-Error { Write-Host "[ERROR] $args" -ForegroundColor Red }

function Install-Binary {
    Write-Info "Building Butter compiler binary..."
    
    $go = Get-Command "go" -ErrorAction SilentlyContinue
    if (-not $go) {
        Write-Error "Go is not installed. Please install Go from https://go.dev/dl/"
        exit 1
    }
    
    Push-Location $ScriptDir
    go build -o $BinaryName main.go
    Pop-Location
    
    if (-not (Test-Path $BinaryName)) {
        Write-Error "Build failed — no binary produced"
        exit 1
    }
    
    $installDirs = @(
        "$env:USERPROFILE\.local\bin",
        "$env:APPDATA\butter"
    )
    
    $installDir = $installDirs[1]
    foreach ($dir in $installDirs) {
        if (Test-Path $dir -PathType Container) {
            $installDir = $dir
            break
        }
    }
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    
    Copy-Item -Path $BinaryName -Destination (Join-Path $installDir $BinaryName) -Force
    Write-Info "Butter compiler installed to $installDir\$BinaryName"
    $installedPath = Join-Path $installDir $BinaryName
    if (Test-Path $installedPath) {
        $version = & $installedPath --version
        Write-Info "$version"
    }
    
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($userPath -notlike "*$installDir*") {
        $newPath = "$installDir;$userPath"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        $env:PATH = "$installDir;$env:PATH"
        Write-Warn "Added $installDir to user PATH. Restart your terminal or run:"
        Write-Warn "  `$env:PATH = `"$installDir;`$env:PATH`""
    }
}

function Install-Extension {
    Write-Info "Installing VS Code Butter extension..."
    
    $code = Get-Command "code" -ErrorAction SilentlyContinue
    if (-not $code) {
        Write-Warn "VS Code 'code' CLI not found. Skipping extension installation."
        Write-Warn "To install manually, open butter-extension/ in VS Code and press F5."
        return
    }
    
    $vsce = Get-Command "vsce" -ErrorAction SilentlyContinue
    if ($vsce) {
        Push-Location $ExtensionDir
        & vsce package --out $VsixOutput 2>$null
        Pop-Location
        
        if (Test-Path $VsixOutput) {
            & code --install-extension $VsixOutput --force
            Remove-Item $VsixOutput -Force
            Write-Info "VS Code Butter extension installed successfully."
            return
        }
    }
    
    $extensionsDir = "$env:USERPROFILE\.vscode\extensions\butter-extension"
    New-Item -ItemType Directory -Force -Path $extensionsDir | Out-Null
    Copy-Item -Path "$ExtensionDir\*" -Destination $extensionsDir -Recurse -Force
    Write-Info "VS Code Butter extension copied to $extensionsDir"
    Write-Warn "Restart VS Code for the extension to take effect."
}

switch ($Command) {
    "install" {
        Install-Binary
        Install-Extension
        Write-Info "Butter installation complete."
    }
    "update" {
        Install-Binary
        Install-Extension
        Write-Info "Butter update complete."
    }
    "binary" {
        Install-Binary
    }
    "extension" {
        Install-Extension
    }
}
