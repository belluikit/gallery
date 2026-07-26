#!/usr/bin/env bash
# Package the gallery for Android, fix up the APK for Android 30+ (uncompressed
# + aligned resources.arsc), sign it with the debug key, and optionally install
# and launch on a running emulator/device.
#
# Why this script exists:
#   * The host shell exports CGO_CFLAGS/CGO_LDFLAGS pointing at Homebrew with a
#     macOS-only linker flag (-Wl,-no_warn_duplicate_libraries) that ld.lld
#     rejects when cross-linking for Android — so we clear them for the build.
#   * The (deprecated) `fyne` packager stores resources.arsc compressed, which
#     Android 14 (API 34) refuses to install; we repack it stored + aligned.
#
# Usage: scripts/android-build.sh [--install]
set -euo pipefail

cd "$(dirname "$0")/.."

: "${ANDROID_HOME:=$HOME/Library/Android/sdk}"
: "${ANDROID_NDK_HOME:=$ANDROID_HOME/ndk/30.0.14904198}"
export ANDROID_HOME ANDROID_NDK_HOME

BT="$ANDROID_HOME/build-tools/34.0.0"
ADB="$ANDROID_HOME/platform-tools/adb"
APK="MD3_Gallery.apk"
APP_ID="io.fyne.md3gallery"
ACTIVITY="org.golang.app.GoNativeActivity"

echo ">> packaging $APK for android/arm64"
env CGO_CFLAGS="" CGO_LDFLAGS="" fyne package -os android/arm64 --tags migrated_fynedo

echo ">> repacking (resources.arsc stored + aligned) and signing"
STAGE="$(mktemp -d)"
unzip -oq "$APK" -d "$STAGE"
rm -rf "$STAGE/META-INF"
rm -f /tmp/gallery-repack.apk /tmp/gallery-aligned.apk
( cd "$STAGE"
  zip -0 -X -q /tmp/gallery-repack.apk resources.arsc
  zip -rg -X -q /tmp/gallery-repack.apk . -x resources.arsc )
"$BT/zipalign" -p -f 4 /tmp/gallery-repack.apk /tmp/gallery-aligned.apk
"$BT/apksigner" sign \
  --ks "$HOME/.android/debug.keystore" --ks-pass pass:android \
  --ks-key-alias androiddebugkey --key-pass pass:android \
  /tmp/gallery-aligned.apk
cp /tmp/gallery-aligned.apk "$APK"
rm -rf "$STAGE"
echo ">> built signed, installable $APK"

if [[ "${1:-}" == "--install" ]]; then
  echo ">> installing and launching on device"
  "$ADB" install -r "$APK"
  "$ADB" shell am start -n "$APP_ID/$ACTIVITY"
fi
