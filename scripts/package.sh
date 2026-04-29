#!/usr/bin/env bash
# Build and package the plugin into a distributable .zip for GitHub Releases.
# Output: rel/novant-datasource-<version>.zip
set -euo pipefail

VERSION=$(node -p "require('./package.json').version")
PLUGIN_ID="novant-datasource"
REL_DIR="rel"
ZIP_PATH="${REL_DIR}/${PLUGIN_ID}-${VERSION}.zip"

echo "==> Building ${PLUGIN_ID} v${VERSION}"
npm run build:all

echo "==> Staging ${PLUGIN_ID}/"
rm -rf .package "${ZIP_PATH}"
mkdir -p .package "${REL_DIR}"
cp -R dist ".package/${PLUGIN_ID}"

echo "==> Creating ${ZIP_PATH}"
( cd .package && zip -rq "../${ZIP_PATH}" "${PLUGIN_ID}" )

echo "==> Cleanup"
rm -rf .package

echo
echo "Done: ${ZIP_PATH}"
ls -lh "${ZIP_PATH}"
