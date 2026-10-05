#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT/public/install.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
printf 'fixture\n' > "$TMP/archive"
hash="$(sha256_file "$TMP/archive")"
verify_hash "$TMP/archive" "$hash"
verify_hash "$TMP/archive" "$(tr '[:lower:]' '[:upper:]' <<< "$hash")"
for bad in invalid '' "$(printf '%064d' 0)"; do
    if (verify_hash "$TMP/archive" "$bad" >/dev/null 2>&1); then exit 1; fi
done
valid="$(jq -cn --arg hash "$hash" '{assets:[{name:"wave.tar.gz",state:"uploaded",digest:("sha256:"+$hash)}]}')"
asset "$valid" wave.tar.gz >/dev/null
[[ "$(asset "$valid" absent true)" == null ]]
for filter in '.assets=[]' '.assets[0].digest=null' '.assets[0].state="new"' '.assets += [.assets[0]]'; do
    if (asset "$(jq "$filter" <<< "$valid")" wave.tar.gz >/dev/null 2>&1); then exit 1; fi
done
printf 'Checksum tests passed\n'
