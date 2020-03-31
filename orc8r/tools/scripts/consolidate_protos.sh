#!/usr/bin/env bash

# consolidate_protos consolidates magma and imported proto files
# to /tmp/magma_protos for easy consumption by e.g. Wireshark.

set -e

outdir=/tmp/magma_protos

magma=~/fbsource/fbcode/magma
include=/usr/local/include

ignore=( -not -path '*/migrations/*' )

pushd "$magma"
find -L . -name '*.proto' "${ignore[@]}" | cpio -pdm --insecure "$outdir"
popd

pushd "$include"
find -L . -name '*.proto' "${ignore[@]}" | cpio -pdm --insecure "$outdir"
popd
