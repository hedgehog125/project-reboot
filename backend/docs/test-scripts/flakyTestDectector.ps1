param (
    [string]$Directory = ".",
    [int]$Interval = 5
)

Set-Location $PSScriptRoot
Set-Location ..\..

Write-Host "Starting continuous test runner..." -ForegroundColor Green
Write-Host "Target directory: $Directory"
Write-Host "Interval: $Interval seconds"
Write-Host "Press Ctrl+C to stop.`n"

while ($true) {
    $Timestamp = Get-Date -Format "HH:mm:ss"
    
    if (Test-Path $Directory) {
        Push-Location $Directory
        
        Write-Host "[$Timestamp] Cleaning test cache..." -ForegroundColor Cyan
        go clean -testcache
        
        Write-Host "[$Timestamp] Running tests..." -ForegroundColor Magenta
        go test ./...
        
        Pop-Location
    } else {
        Write-Host "Error: Directory '$Directory' not found." -ForegroundColor Red
        break
    }

    Write-Host "`n--- Waiting $Interval seconds ---" -ForegroundColor Gray
    Start-Sleep -Seconds $Interval
}
