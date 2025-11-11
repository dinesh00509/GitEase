Write-Host "Installing GitEase..."

$version = "latest"
$repo = "dinesh00509/gitease"

if ($version -eq "latest") {
    $version = (Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest").tag_name
}

$url = "https://github.com/$repo/releases/download/$version/gitease_Windows_x86_64.zip"
$temp = "$env:TEMP\gitease.zip"

Write-Host " Downloading $url..."
Invoke-WebRequest -Uri $url -OutFile $temp

Write-Host " Extracting..."
Expand-Archive -Path $temp -DestinationPath $env:TEMP -Force

Write-Host " Installing..."
Move-Item -Path "$env:TEMP\gitease.exe" -Destination "C:\Windows\System32\gitease.exe" -Force

Write-Host "GitEase installed successfully!"
gitease --version

