#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT/public/install.sh"
fetch() {
    printf '%s' '[{"id":1,"tag_name":"v1.0.0","draft":false,"published_at":"2026-01-01T00:00:00Z"},{"id":2,"tag_name":"nightly","draft":false,"published_at":"2026-10-05T00:00:00Z"},{"id":3,"tag_name":"v2.0.0","draft":true,"published_at":"2026-10-05T00:00:00Z"},{"id":4,"tag_name":"v1.1.0-alpha","prerelease":true,"draft":false,"published_at":"2026-10-01T00:00:00Z"}]'
}
[[ "$(latest_release wavefnd/Wave | jq -r .tag_name)" == v1.1.0-alpha ]]
fetch() {
    case "${*: -1}" in
        *page=1) jq -cn '[range(0;100) | {id:.,tag_name:"nightly",draft:false,published_at:"2026-10-05T00:00:00Z"}]' ;;
        *page=2) printf '%s' '[{"id":101,"tag_name":"v1.0.0","draft":false,"published_at":"2026-10-01T00:00:00Z"}]' ;;
        *) return 22 ;;
    esac
}
[[ "$(latest_release wavefnd/Wave | jq -r .tag_name)" == v1.0.0 ]]
for option in '--version 1.0.0' '--vex-version 1.0.0' '--version=1.0.0' nightly v1.0.0; do
    if message="$(bash "$ROOT/public/install.sh" $option 2>&1)"; then echo "accepted $option"; exit 1; fi
    [[ "$message" == *'manually'* ]]
done
for pair in 'Linux x86_64' 'Linux aarch64' 'Linux riscv64' 'Linux loongarch64' 'Darwin x86_64' 'Darwin arm64' 'FreeBSD amd64'; do
    resolve_platform $pair
    [[ -n "$WAVE_TARGET" && -n "$WAVE_SUFFIX" && -n "$VEX_SUFFIX" ]]
done
if (resolve_platform Linux i686 >/dev/null 2>&1); then exit 1; fi
if (resolve_platform Windows ARM64 >/dev/null 2>&1); then exit 1; fi
printf 'Latest-only policy and platform tests passed\n'
