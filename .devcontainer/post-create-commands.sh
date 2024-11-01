#!/bin/bash
# $1: should be containerWorkspaceFolder from https://code.visualstudio.com/docs/remote/devcontainerjson-reference
if [ "$1" != "$MAGMA_ROOT" ]
then
    sudo rm -rf "$MAGMA_ROOT"
    sudo ln -s "$1" "$MAGMA_ROOT"
fi

sudo ln -s "$MAGMA_ROOT"/lte/gateway/configs /etc/magma
echo "alias magtivate='source /home/vscode/build/python/bin/activate'" >> ~/.bashrc

echo "Generating compile_commands.json for C/C++ code navigation"
"$MAGMA_ROOT"/dev_tools/gen_compilation_database.py

echo "Setting up Bazel Bash completion"
"$MAGMA_ROOT"/bazel/scripts/setup_bazel_bash_completion.sh $(cat "$MAGMA_ROOT"/.bazelversion)
