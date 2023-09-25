#!/bin/bash

# $1: should be containerWorkspaceFolder from https://code.visualstudio.com/docs/remote/devcontainerjson-reference

sudo ln -s "$1"/lte/gateway/configs /etc/magma
echo "alias magtivate='source /home/vscode/build/python/bin/activate'" >> ~/.bashrc

echo "Generating compile_commands.json for C/C++ code navigation"
"$1"/dev_tools/gen_compilation_database.py

echo "Setting up Bazel Bash completion"
"$1"/bazel/scripts/setup_bazel_bash_completion.sh $(cat "$1"/.bazelversion)
