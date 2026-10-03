#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

# Load the installer helpers without running its download/install entry point.
source <(sed -n '1,/^resolve_shell_rc()/p' "$ROOT/public/install.sh" | sed '$d')

archive="$TMP/wave.tar.gz"
sums="$TMP/SHA256SUMS"
name="wave.tar.gz"
printf 'checksum regression fixture\n' > "$archive"
digest="$(sha256_file "$archive")"

check_ok() {
    local entry="$1"
    printf '%s\n' "$entry" > "$sums"
    verify_checksum "$archive" "$sums" "$name" >/dev/null
}

check_fail() {
    local entry="$1"
    printf '%s\n' "$entry" > "$sums"
    if (verify_checksum "$archive" "$sums" "$name" >/dev/null 2>&1); then
        echo "expected checksum verification to fail: $entry" >&2
        exit 1
    fi
}

upper="$(printf '%s' "$digest" | LC_ALL=C tr '[:lower:]' '[:upper:]')"
mixed="$(printf '%s' "$digest" | sed 's/a/A/g; s/b/B/g; s/c/C/g')"

check_ok "$digest  $name"
check_ok "$upper  $name"
check_ok "$mixed *$name"
check_fail "$(printf '%064d' 0)  $name"
check_fail "not-a-digest  $name"
check_fail "$digest  another-file.tar.gz"

echo "install checksum tests passed"


# GitHub digests replace Wave's public checksum attachment.
asset_response="$(jq -n --arg name "$name" --arg hash "$upper" \
    '{assets: [{name: "unrelated.zip", state: "uploaded", digest: null},
               {name: $name, state: "uploaded", digest: ("sha256:" + $hash)}]}')"
curl() { printf '%s' "$asset_response"; }
[[ "$(published_asset_digest wavefnd/Wave v0.2.1-pre-beta "$name")" == "$upper" ]]
verify_hash "$archive" "$(published_asset_digest wavefnd/Wave v0.2.1-pre-beta "$name")" "$name" >/dev/null
valid="$asset_response"
for filter in '.assets = []' '.assets[1].name = "wrong.zip"' \
    '.assets[1].digest = null' '.assets[1].digest = "sha256:bad"' \
    '.assets[1].digest |= sub("sha256:"; "sha512:")' \
    '.assets[1].state = "new"' '.assets += [.assets[1]]'; do
    asset_response="$(jq "$filter" <<< "$valid")"
    if (published_asset_digest wavefnd/Wave v0.2.1-pre-beta "$name" >/dev/null 2>&1); then
        echo "accepted invalid asset metadata: $filter" >&2; exit 1
    fi
done
asset_response='not JSON'
if (published_asset_digest wavefnd/Wave v0.2.1-pre-beta "$name" >/dev/null 2>&1); then exit 1; fi
curl() { return 22; }
if (published_asset_digest wavefnd/Wave v0.2.1-pre-beta "$name" >/dev/null 2>&1); then exit 1; fi
if (verify_hash "$archive" "$(printf '%064d' 0)" "$name" >/dev/null 2>&1); then exit 1; fi
echo 'installer GitHub digest tests passed'
