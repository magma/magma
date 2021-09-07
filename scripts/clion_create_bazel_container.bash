#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob


export PATH_TO_DOCKERFILE="${MAGMA_ROOT}"/experimental/bazel-base/Dockerfile
export BAZEL_OUTPUT_ROOT="${HOME}"/.cache/bazel-magma
export ASPECT_PATCHED="${BAZEL_OUTPUT_ROOT}/clwb_aspect"
export ASPECT_ORIG="/Users/marwhal/Library/Application Support/JetBrains/CLion2021.1/plugins/clwb/aspect"

rm -rf "${ASPECT_PATCHED}"
mkdir -p "${ASPECT_PATCHED}"
cp -r "${ASPECT_ORIG}" "${ASPECT_PATCHED}"

docker stop magma-builder
docker rm magma-builder

docker run -d \
  --name magma-builder \
  -v "${MAGMA_ROOT}":/magma \
  -v "${MAGMA_ROOT}"/lte/gateway/configs:/etc/magma \
  -v "${BAZEL_OUTPUT_ROOT}":"${BAZEL_OUTPUT_ROOT}" \
  -v "${ASPECT_PATCHED}":"${ASPECT_ORIG}" \
  -t magma/bazel-build:latest
