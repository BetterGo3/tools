# Rebuild gopls against the fork GOROOT so it picks up go/printer changes.
$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent $PSScriptRoot
$GOROOT = Join-Path $ProjectRoot "go"
$env:GOROOT = $GOROOT
$env:GOEXPERIMENT = "genericmethods"
Set-Location (Join-Path $PSScriptRoot "gopls")
& "$GOROOT\bin\go.exe" build -o (Join-Path $PSScriptRoot "gopls.exe") .
Write-Host "Built $PSScriptRoot\gopls.exe"
Get-Item (Join-Path $PSScriptRoot "gopls.exe") | Select-Object FullName, LastWriteTime, Length
