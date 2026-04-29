#!/usr/bin/env bash
# Publish an existing release artifact from rel/ to a GitHub Release.
# Does NOT build — run `npm run package` first to produce zips.
#
# Usage:
#   npm run publish               # auto-pick if exactly one unpublished zip exists
#   npm run publish -- <version>  # publish a specific version, e.g. 1.0.0
set -euo pipefail

PLUGIN_ID="novant-datasource"
REL_DIR="rel"

# gh CLI must be available.
if ! command -v gh >/dev/null 2>&1; then
  echo "Error: gh CLI not found. Install it from https://cli.github.com/"
  exit 1
fi

# Find all zips in rel/.
ZIPS=()
shopt -s nullglob
for f in "${REL_DIR}/${PLUGIN_ID}"-*.zip; do
  ZIPS+=("$f")
done
shopt -u nullglob

if [ ${#ZIPS[@]} -eq 0 ]; then
  echo "No artifacts in ${REL_DIR}/. Run 'npm run package' first."
  exit 1
fi

# Fetch existing GitHub release tags.
EXISTING_TAGS=$(gh release list --limit 100 --json tagName --jq '.[].tagName' 2>/dev/null || echo "")

# Sort zips into published / unpublished by checking against existing release tags.
UNPUBLISHED=()
PUBLISHED=()
for zip in "${ZIPS[@]}"; do
  filename=$(basename "$zip")
  version=${filename#${PLUGIN_ID}-}
  version=${version%.zip}
  tag="v${version}"
  if echo "${EXISTING_TAGS}" | grep -qFx "${tag}"; then
    PUBLISHED+=("${version}")
  else
    UNPUBLISHED+=("${version}")
  fi
done

# Show inventory.
echo "Artifacts in ${REL_DIR}/:"
if [ ${#PUBLISHED[@]} -gt 0 ]; then
  for v in "${PUBLISHED[@]}"; do
    echo "   v${v}  (already published)"
  done
fi
if [ ${#UNPUBLISHED[@]} -gt 0 ]; then
  for v in "${UNPUBLISHED[@]}"; do
    echo "   v${v}  (unpublished)"
  done
fi
echo

# Decide which version to publish.
TARGET="${1:-}"

if [ -n "${TARGET}" ]; then
  # Explicit target via arg.
  found=""
  if [ ${#UNPUBLISHED[@]} -gt 0 ] && printf '%s\n' "${UNPUBLISHED[@]}" | grep -qFx "${TARGET}"; then
    VERSION="${TARGET}"
    found="unpublished"
  elif [ ${#PUBLISHED[@]} -gt 0 ] && printf '%s\n' "${PUBLISHED[@]}" | grep -qFx "${TARGET}"; then
    echo "Error: v${TARGET} is already published."
    exit 1
  fi
  if [ -z "${found}" ]; then
    echo "Error: no artifact for version ${TARGET} in ${REL_DIR}/."
    exit 1
  fi
elif [ ${#UNPUBLISHED[@]} -eq 0 ]; then
  echo "All artifacts are already published. Bump version, run 'npm run package',"
  echo "then 'npm run publish' to release a new one."
  exit 1
elif [ ${#UNPUBLISHED[@]} -eq 1 ]; then
  VERSION="${UNPUBLISHED[0]}"
else
  # Multiple unpublished — prompt for a numbered selection.
  echo "Multiple unpublished artifacts. Select one:"
  i=1
  for v in "${UNPUBLISHED[@]}"; do
    echo "   ${i}) v${v}"
    i=$((i + 1))
  done
  echo
  read -r -p "Select [1-${#UNPUBLISHED[@]}, q to abort]: " choice < /dev/tty
  if [[ "${choice}" == "q" || -z "${choice}" ]]; then
    echo "Aborted."
    exit 1
  fi
  if ! [[ "${choice}" =~ ^[0-9]+$ ]] || [ "${choice}" -lt 1 ] || [ "${choice}" -gt ${#UNPUBLISHED[@]} ]; then
    echo "Invalid selection."
    exit 1
  fi
  VERSION="${UNPUBLISHED[$((choice - 1))]}"
fi

TAG="v${VERSION}"
ZIP_PATH="${REL_DIR}/${PLUGIN_ID}-${VERSION}.zip"

# Working tree must be clean since the tag will be placed on HEAD.
if ! git diff-index --quiet HEAD --; then
  echo "Error: working tree is dirty. Commit or stash changes first."
  exit 1
fi

# Local tag must not already exist (would mean a half-finished previous publish).
if git rev-parse "${TAG}" >/dev/null 2>&1; then
  echo "Error: local tag ${TAG} exists but no GitHub release was found for it."
  echo "Either delete it (git tag -d ${TAG}) or finish manually:"
  echo "   git push origin ${TAG} && gh release create ${TAG} ${ZIP_PATH} --generate-notes"
  exit 1
fi

# Confirm with the user before any destructive action.
echo "About to publish:"
echo "   Tag:    ${TAG}  (at $(git rev-parse --short HEAD))"
echo "   Asset:  ${ZIP_PATH}"
echo "   Remote: $(git remote get-url origin 2>/dev/null || echo 'origin')"
echo
read -r -p "Proceed? [y/N]: " confirm < /dev/tty
if [[ "${confirm}" != "y" && "${confirm}" != "Y" ]]; then
  echo "Aborted."
  exit 1
fi

echo "==> Tagging ${TAG}"
git tag "${TAG}"

echo "==> Pushing tag to origin"
git push origin "${TAG}"

echo "==> Creating GitHub release"
gh release create "${TAG}" "${ZIP_PATH}" --title "${TAG}" --generate-notes

echo
echo "Done: released ${TAG}"
