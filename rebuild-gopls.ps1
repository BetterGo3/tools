# Rebuild gopls against the fork GOROOT so it picks up go/printer changes.
$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent $PSScriptRoot
$GOROOT = Join-Path $ProjectRoot "go"
$env:GOROOT = $GOROOT
$env:GOEXPERIMENT = "genericmethods"
$outDir = Join-Path $PSScriptRoot "bin"
New-Item -ItemType Directory -Force -Path $outDir | Out-Null
$outPath = Join-Path $outDir "gopls.exe"
Set-Location (Join-Path $PSScriptRoot "gopls")
& "$GOROOT\bin\go.exe" build -o $outPath .
Write-Host "Built $outPath"
Get-Item $outPath | Select-Object FullName, LastWriteTime, Length
