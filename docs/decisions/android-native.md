# Native Android application

Status: accepted project direction (September 25, 2026).

Related work: [architecture #77](https://github.com/wavefnd/wave-platform/issues/77),
[application shell #79](https://github.com/wavefnd/wave-platform/issues/79),
and [Android roadmap #76](https://github.com/wavefnd/wave-platform/issues/76).

## Decision

Build a native Kotlin application with Jetpack Compose and Material 3 in `android/`.
Do not embed the website in a WebView or ship a Trusted Web Activity. The maintainer
explicitly selected native Android with an Android-specific UI/UX.

Reuse the Go host, public XML API contracts, content, and brand colors. Implement
Android presentation and lifecycle handling independently of Vue. A native client
costs additional UI implementation and testing, but gives the planned offline
storage, App Links, notifications, accessibility, and account lifecycle work a
platform-owned implementation surface.

A web shell or TWA would reuse more existing Vue presentation and browser sessions.
Those approaches do not meet the selected native UI requirement, and browser
session reuse must not be assumed for subsequent authentication work.

## Initial structure and support

Use a single `:app` module with `data` and `ui` packages. Separate the catalog
repository from its ViewModel and stateless screen components; introduce more
modules only when a feature establishes a useful boundary. Dependencies and build
tools are pinned in the version catalog and checksum-verified Gradle Wrapper.

The initial baseline is Android 8.0/API 26 and newer, compiled and targeted at
API 37. API 26 keeps the first contributor environment and platform behavior
manageable while retaining support for older devices. This is a baseline to test,
not a claim that every supported OS/device has already been verified.

The app preserves the website's identity rather than introducing a separate visual
direction. Colors follow `frontend/src/ui/tokens.css`; the logo is the unchanged PNG
embedded in `frontend/public/img/wave-logo.ico`. Underlined project tabs, grouped
document rows, thin separators, compact typography, and 2–8 dp corner radii follow
`frontend/src/ui/service-pages.css` and `frontend/src/styles.css`. Android uses its
platform sans-serif fallback and at least 48 dp touch targets. Native rendering
does not imply replacing the existing brand or content hierarchy.

The app uses a bottom navigation bar on compact windows and a navigation rail at
600 dp and wider. Content remains width-bounded. Docs, News, and Settings form the
initial navigation. Project selection and reading language are independent;
interface text follows the device's English/Korean locale, and documentation keeps
all nine existing language choices. Themes support light, dark, and device settings.

## Executable feasibility slice

The foundation fetches anonymous document catalogs from the existing XML endpoint,
merges translations over English by path, and displays native Wave/Whale lists.
It has explicit loading, retry, and empty states. Catalog responses cannot replace
a more recently selected language. Tests cover the real contract using synthetic
fixtures, including namespaced XML and an unavailable service.

Article actions are explicitly labeled as opening an external browser. They are a
temporary catalog-preview action, not an embedded browser or the final native
reader. The News destination likewise describes its current availability without
invented articles. The first product release should include native docs and
blog/release reading; this scaffold alone is not that release.

## Boundaries for subsequent contributions

The full XML transport/error contract remains #83; this prototype only reads public
catalogs. Native article reading remains #84, and native news reading remains #85.
Offline delivery follows #89/#91 and #90. Sign-in and storage of account material
follow #88; the catalog does not receive website session cookies or bypass
TOTP/Turnstile. Push transport remains #107, and verified App Links remain #86.

The current application namespace is `dev.wavelang.platform`; debug installs use
`dev.wavelang.platform.debug`. These establish a buildable contributor target and
do not reserve a store identity or select production signing certificates. Review
store identity, distribution, signing, and device coverage before public release.

## References

- [Compose](https://developer.android.com/develop/ui/compose/documentation)
- [Android architecture](https://developer.android.com/topic/architecture/recommendations)
- [AGP 9.4 compatibility](https://developer.android.com/build/releases/agp-9-4-0-release-notes)
