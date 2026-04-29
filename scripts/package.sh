#!/usr/bin/env bash
# Build and package the plugin into a distributable .zip for GitHub Releases.
# Output: ./novant-datasource-<version>.zip
set -euo pipefail

VERSION=$(node -p "require('./package.json').version")
PLUGIN_ID="novant-datasource"
ZIP_NAME="${PLUGIN_ID}-${VERSION}.zip"

echo "==> Building ${PLUGIN_ID} v${VERSION}"
npm run build:all

echo "==> Staging ${PLUGIN_ID}/"
rm -rf .package "${ZIP_NAME}"
mkdir -p .package
cp -R dist ".package/${PLUGIN_ID}"

echo "==> Creating ${ZIP_NAME}"
( cd .package && zip -rq "../${ZIP_NAME}" "${PLUGIN_ID}" )

echo "==> Cleanup"
rm -rf .package

echo
echo "Done: ${ZIP_NAME}"
ls -lh "${ZIP_NAME}"
