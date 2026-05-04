# atlas.ps1 - โหลด .env แล้วรัน atlas command

# โหลด .env.dev หรือ .env
$envFile = if (Test-Path ".env.dev") { ".env.dev" } else { ".env" }

if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#=]+)=(.*)$') {
            [Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim(), "Process")
        }
    }
    Write-Host "Loaded $envFile"
}

# รัน atlas command
if (Test-Path ".\\atlas.exe") {
    & ".\\atlas.exe" @args
} elseif (Test-Path ".\\atlas") {
    & ".\\atlas" @args
} else {
    atlas @args
}