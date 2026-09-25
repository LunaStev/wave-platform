# Wave for Android

A native Kotlin and Jetpack Compose application. Start with this project to
contribute Android features; it does not embed the Wave website in a WebView.
See the [architecture decision](../docs/decisions/android-native.md) and
[Android roadmap](https://github.com/wavefnd/wave-platform/issues/76).

## Requirements

- JDK 17.
- An Android Studio version compatible with AGP 9.4, or Android command-line tools.
- SDK platform `platforms;android-37.0` and build tools `36.0.0`.
- An emulator or device running API 26 or newer for installation/UI tests.

Gradle 9.6.0 is provided by the checked-in Wrapper; a separate Gradle installation
is not required. The version catalog pins AGP 9.4.0, Kotlin/Compose compiler 2.2.10,
Compose BOM 2026.02.01, Activity 1.13.0, Lifecycle 2.10.0, and Coroutines 1.10.2.
These are the verified foundation versions, not floating latest-version selectors.

Open the `android` directory in Android Studio. Select JDK 17 as the Gradle JDK,
install the listed SDK packages, and let Studio create an ignored `local.properties`.
If Studio creates `gradle/gradle-daemon-jvm.properties`, check its JVM criteria too:
they take precedence over `JAVA_HOME`. Keep machine-generated toolchain changes
out of feature commits unless the supported build baseline is intentionally updated.
For command-line use, set `JAVA_HOME` and `ANDROID_HOME` to your installations:

```sh
sdkmanager 'platforms;android-37.0' 'build-tools;36.0.0'
cd android
./gradlew assembleDebug testDebugUnitTest lintDebug
```

On Windows, use `gradlew.bat`. The APK is produced at
`app/build/outputs/apk/debug/app-debug.apk`.

## Run against the local platform

Start the platform using the repository's [Docker instructions](../README.md#run-with-docker).
The debug app defaults to `http://10.0.2.2:8080/`, which reaches the host from the
standard Android emulator. Start an emulator in Android Studio, then run:

```sh
./gradlew installDebug
adb shell am start -n dev.wavelang.platform.debug/dev.wavelang.platform.MainActivity
```

For a USB-connected device, forward its localhost to the host and build with:

```sh
adb reverse tcp:8080 tcp:8080
./gradlew installDebug -Pwave.debugApiBaseUrl=http://127.0.0.1:8080/
```

To use an HTTPS development server:

```sh
./gradlew installDebug -Pwave.debugApiBaseUrl=https://your-development-host.example/
```

Only debug builds accept this override. HTTP is allowed only for the emulator
host, localhost, and `127.0.0.1`. Release configuration uses
`https://wave-lang.dev/` and does not include the debug cleartext exceptions.
No `.env`, account secret, or production signing key is needed for this reader.

## What is implemented

- Native compact bottom navigation and expanded navigation rail.
- The website's existing logo, light/dark colors, underlined tabs, thin separators,
  and grouped document lists, rendered with Compose and Android touch targets.
- Light/dark/device theme, English/Korean UI, persisted reading language.
- Actual anonymous XML document catalogs with Wave/Whale separation, translation
  merging, metadata ordering, and English fallback labels.
- Loading, empty, failure/retry states and protection against stale requests.
- Explicit external-browser article links while the native reader is developed.
- A truthful News placeholder; native news reading is a follow-up.

The app contains no account integration, offline downloads, or push registration.
The initial package ID is a contributor target, not a finalized store identity.

## Tests

```sh
./gradlew testDebugUnitTest lintDebug assembleDebug assembleDebugAndroidTest
./gradlew connectedDebugAndroidTest
```

The last command needs a running emulator/device. UI tests inject synthetic state
and do not contact the platform. Use an English emulator for the initial UI suite;
the app itself follows the device's English/Korean interface language.

Unit tests cover XML parsing, project/translation ordering, language defaults,
error/retry states, restored selections, and a response arriving after cancellation.
UI tests cover the real launcher, project isolation, navigation restoration, language selection, and
retry presentation. Reports are under `app/build/reports/`.
The device CI artifact also includes screenshots rendered from the synthetic UI
fixtures so visual changes can be reviewed without running a local emulator.

On a busy development machine, build before starting an emulator and use
`--max-workers=1` to limit parallel work. Device tests can also run in PR CI.

The Android CI workflow runs build, unit tests, lint, and emulator UI tests without
signing secrets. Debug APKs and reports are retained for seven days. Existing Go
and frontend/editor workflows are independent.

## Where to contribute

- `app/src/main/java/dev/wavelang/platform/data/`: anonymous catalog contract and XML adapter.
- `app/src/main/java/dev/wavelang/platform/ui/`: state, native screens, and theme.
- `app/src/main/res/values*/`: interface strings; document language is separate.
- `app/src/test/`: deterministic contract/state tests.
- `app/src/androidTest/`: device UI tests.

Keep new screens native. Reuse the existing API/content contracts and add small,
reviewable features under their roadmap issues. Full native document and article
readers are tracked in #84 and #85; do not replace them with an embedded website.
