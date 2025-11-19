# Port the application is running on
$port = 9000

# Output file
$outputFile = "resource_usage_1.csv"

# Find the process listening on the port
Write-Host "Looking for process on port $port..."
$tcpConnection = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue
if (-not $tcpConnection) {
    Write-Host "No process found listening on port $port. Start the application and re-run this script."
    exit
}
$processId = $tcpConnection.OwningProcess
$process = Get-Process -Id $processId

# Write header to CSV
"Timestamp,CPU_TotalSeconds,Memory_MB" | Out-File -FilePath $outputFile

Write-Host "Monitoring process '$($process.ProcessName)' (PID: $($process.Id)) on port $port"
Write-Host "Press Ctrl+C to stop monitoring."

# Monitor and log
try {
    while ($true) {
        $process.Refresh()
        $timestamp = (Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
        # Note: $process.CPU is not available on all systems/versions of PowerShell
        # TotalProcessorTime is a cumulative value. We can calculate the differential, but for simplicity, we log the total.
        $cpu = $process.TotalProcessorTime.TotalSeconds
        $mem = [math]::Round($process.WorkingSet / 1MB, 2)
        "$timestamp,$cpu,$mem" | Out-File -FilePath $outputFile -Append
        Start-Sleep -Seconds 1
    }
}
catch {
    Write-Host "Monitoring stopped."
}
