$ErrorActionPreference = "Stop"

$Version = if ($env:VERSION) { $env:VERSION } else { "dev" }
$LDFlags = "-s -w -X main.version=$Version"
$OutDir = "bin"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

Write-Host "Building bull $Version ..."

go build -ldflags="$LDFlags" -o "$OutDir/bull.exe" ./cmd/bull/
Write-Host "  -> $OutDir/bull.exe"

# Cross compile (set $env:CROSS = "1" to enable)
if ($env:CROSS -eq "1") {
    $targets = @(
        @{ os = "linux";   arch = "amd64"; ext = "" },
        @{ os = "linux";   arch = "arm64"; ext = "" },
        @{ os = "darwin";  arch = "amd64"; ext = "" },
        @{ os = "darwin";  arch = "arm64"; ext = "" },
        @{ os = "windows"; arch = "amd64"; ext = ".exe" }
    )
    foreach ($t in $targets) {
        $output = "$OutDir/bull-$($t.os)-$($t.arch)$($t.ext)"
        Write-Host "  building $($t.os)/$($t.arch) ..."
        $env:GOOS = $t.os
        $env:GOARCH = $t.arch
        go build -ldflags="$LDFlags" -o $output ./cmd/bull/
        Write-Host "  -> $output"
    }
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
}

Write-Host "Done."
