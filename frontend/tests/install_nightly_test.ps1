$ErrorActionPreference = 'Stop'
$script = Join-Path $PSScriptRoot '../public/install.ps1'
$tokens = $null
$errors = $null
$ast = [System.Management.Automation.Language.Parser]::ParseFile((Resolve-Path $script), [ref]$tokens, [ref]$errors)
if ($errors.Count) { throw $errors }
# Load only helper definitions; do not run any installation entry point.
foreach ($function in $ast.FindAll({param($node) $node -is [System.Management.Automation.Language.FunctionDefinitionAst]}, $false)) {
    Invoke-Expression $function.Extent.Text
}
function Fail($Message) { throw $Message }
foreach ($tag in @('nightly', 'vnightly', 'NIGHTLY', 'vNightly')) {
    foreach ($helper in @('Normalize-Version', 'Assert-Version')) {
        $rejected = $false
        try { & $helper $tag } catch {
            $rejected = $_.Exception.Message -like '*manual download*'
        }
        if (-not $rejected) { throw "Accepted $tag through $helper" }
    }
}
if ((Normalize-Version '0.2.1-pre-beta') -ne 'v0.2.1-pre-beta') { throw 'version normalization changed' }
Assert-Version 'v0.2.1-pre-beta'
function Invoke-RestMethod {
    param($Headers, $Uri)
    return ,@(@{tag_name='nightly'; prerelease=$true}, @{tag_name='v0.2.1-pre-beta'; prerelease=$true})
}
if ((Resolve-LatestVersion 'wavefnd/Wave') -ne 'v0.2.1-pre-beta') { throw 'versioned prerelease not selected' }
function Invoke-RestMethod {
    param($Headers, $Uri)
    if ($Uri.EndsWith('page=1')) { return ,@(@{tag_name='nightly'; prerelease=$true}) }
    return ,@(@{tag_name='v0.2.0-pre-beta'; prerelease=$true})
}
if ((Resolve-LatestVersion 'wavefnd/Wave') -ne 'v0.2.0-pre-beta') { throw 'pagination failed' }
$script:emptyRequests = 0
function Invoke-RestMethod {
    param($Headers, $Uri)
    $script:emptyRequests++
    if ($script:emptyRequests -gt 1) { throw 'unexpected extra page after empty array' }
    return ,@()
}
$rejected = $false
try { Resolve-LatestVersion 'wavefnd/Wave' } catch { $rejected = $true }
if (-not $rejected -or $script:emptyRequests -ne 1) { throw 'empty release list accepted or paginated' }
Write-Host 'installer Nightly policy tests passed'


# Load real helper definitions above, mock only GitHub's JSON response.
$WaveRepo = 'wavefnd/Wave'
$temp = [System.IO.Path]::GetTempFileName()
try {
    [System.IO.File]::WriteAllText($temp, 'release digest fixture')
    $hash = (Get-FileHash -LiteralPath $temp -Algorithm SHA256).Hash
    $names = @('wave-v0.2.1-pre-beta-x86_64-pc-windows-msvc.zip',
               'wave-v0.2.1-pre-beta-x86_64-pc-windows-gnu.zip')
    $msvc = @{ name=$names[0]; state='uploaded'; digest="sha256:$hash" }
    $gnu = @{ name=$names[1]; state='uploaded'; digest="sha256:$hash" }
    $script:assetResponse = @{ assets=@($gnu, $msvc) }
    function Invoke-RestMethod { param($Headers, $Uri); return $script:assetResponse }
    if ((Get-WaveReleaseAsset 'v0.2.1-pre-beta' $names).name -ne $names[0]) { throw 'MSVC not preferred' }
    Assert-Hash $temp $hash $names[0]
    Assert-Hash $temp $hash.ToLowerInvariant() $names[0]
    $script:assetResponse = @{ assets=@($gnu) }
    if ((Get-WaveReleaseAsset 'v0.2.0-pre-beta' $names).name -ne $names[1]) { throw 'legacy GNU not retained' }
    foreach ($bad in @(
        @{}, @{name=$names[0]; state='uploaded'; digest=$null},
        @{name=$names[0]; state='uploaded'; digest='sha256:bad'},
        @{name=$names[0]; state='uploaded'; digest="sha512:$hash"},
        @{name=$names[0]; state='new'; digest="sha256:$hash"}
    )) {
        $script:assetResponse = @{ assets=@($bad) }
        $rejected = $false
        try { Get-WaveReleaseAsset 'v0.2.1-pre-beta' $names } catch { $rejected = $true }
        if (-not $rejected) { throw 'invalid asset accepted' }
    }
    $script:assetResponse = @{ assets=@($msvc, $msvc) }
    $rejected = $false
    try { Get-WaveReleaseAsset 'v0.2.1-pre-beta' $names } catch { $rejected = $true }
    if (-not $rejected) { throw 'duplicate asset accepted' }
    function Invoke-RestMethod { param($Headers, $Uri); throw 'network unavailable' }
    $rejected = $false
    try { Get-WaveReleaseAsset 'v0.2.1-pre-beta' $names } catch { $rejected = $true }
    if (-not $rejected) { throw 'API failure accepted' }
    foreach ($badHash in @($null, '', 'invalid', ('0' * 64))) {
        $rejected = $false
        try { Assert-Hash $temp $badHash $names[0] } catch { $rejected = $true }
        if (-not $rejected) { throw 'invalid checksum accepted' }
    }
} finally { Remove-Item -LiteralPath $temp }
Write-Host 'installer GitHub digest tests passed'
