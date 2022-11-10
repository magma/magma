#!/bin/bash

# $1: should be containerWorkspaceFolder from https://code.visualstudio.com/docs/remote/devcontainerjson-reference

sudo ln -s "$1"/lte/gateway/configs /etc/magma
echo "alias magtivate='source /home/vscode/build/python/bin/activate'" >> ~/.bashrc

# Only pull in cache for devcontainer opened inside GitHub Codespaces.
# Locally opened Devcontainer has access to persistent cache directory in .bazel-cache and .bazel-cache-repo.
# This is a little hacky, but GITHUB_CODESPACE_TOKEN should only be defined for the Codespace case
if [[ -z $GITHUB_CODESPACE_TOKEN ]]; then
    echo "Assuming the devcontainer is opened locally, not using Bazel remote cache..."
else
    echo "Assuming the devcontainer is opened in GitHub Codespaces, using read-only Bazel remote cache..."
    cache_key=bazel-base-image  #  the devcontainer is based on the bazel-base container
    (cd "$1" && bazel/scripts/remote_cache_bazelrc_setup.sh $cache_key)
fi

echo "Generating compile_commands.json for C/C++ code navigation"
"$1"/dev_tools/gen_compilation_database.py

echo "Setting up Bazel Bash completion"
"$1"/bazel/scripts/setup_bazel_bash_completion.sh $(cat "$1"/.bazelversion)
