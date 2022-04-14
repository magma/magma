#!/bin/bash

# $1: should be containerWorkspaceFolder from https://code.visualstudio.com/docs/remote/devcontainerjson-reference

sudo ln -s "$1"/lte/gateway/configs /etc/magma
echo "alias magtivate='source /home/vscode/build/python/bin/activate'" >> ~/.bashrc

# Only pull in cache for devcontainer opened inside GitHub Codespaces.
# Locally opened Devcontainer has access to persistent cache directory in .bazel-cache and .bazel-cache-repo.
# This is a little hacky, but GITHUB_CODESPACE_TOKEN should only be defined for the Codespace case
if [[ -z $GITHUB_CODESPACE_TOKEN ]]; then
    echo "Assuming the devcontainer is opened locally, not pulling in Bazel cache..."
else
    echo "Assuming the devcontainer is opened in GitHub Codespaces, pulling in Bazel cache..."
    # Pull in cache to speed up Bazel build. The cache is populated by a GitHub Action job periodically (.github/workflows/bazel-cache-push.yml)
    # Fetch repository cache
    wget -qO- https://magma-cache.s3.amazonaws.com/bazel-cache-repo-devcontainer.tar.gz | tar xvfz -
    # Fetch build cache
    wget -qO- https://magma-cache.s3.amazonaws.com/bazel-cache-devcontainer.tar.gz | tar xvfz -
fi
