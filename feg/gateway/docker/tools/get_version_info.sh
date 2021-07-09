#!/bin/bash

branch=$(git rev-parse --abbrev-ref HEAD)
tag=$(git tag --points-at HEAD)
chash=$(git rev-parse --short HEAD)
commit_date=$(git show -s --format=%ci)

printf '{\n "MAGMA_BUILD_BRANCH":"%s",\n"MAGMA_BUILD_TAG":"%s",\n"MAGMA_BUILD_COMMIT_HASH":"%s",\n "MAGMA_BUILD_COMMIT_DATE":"%s"\n}' "$branch" "$tag" "$chash" "$commit_date"
