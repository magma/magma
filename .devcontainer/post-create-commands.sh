#!/bin/bash

# $1: should be containerWorkspaceFolder from https://code.visualstudio.com/docs/remote/devcontainerjson-reference

sudo ln -s "$1"/lte/gateway/configs /etc/magma
echo "alias magtivate='source /home/vscode/build/python/bin/activate'" >> ~/.bashrc

# Pull in cache to speed up Bazel build. The cache is populated by a GitHub Action job periodically (.github/workflows/bazel-cache-push.yml)
# Fetch repository cache
wget -qO- https://magma-cache.s3.amazonaws.com/bazel-cache-repo-devcontainer.tar.gz | tar xvfz - 
# Fetch build cache
wget -qO- https://magma-cache.s3.amazonaws.com/bazel-cache-devcontainer.tar.gz | tar xvfz - 
