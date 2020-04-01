# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
clean::
	go clean ./...

fmt::
	go fmt ./...

# operator-sdk enforces a strict directory
# structure and magma ignores build
# directories. Copy to allow generation,
# then remove
gen::
	cp -r docker build
	operator-sdk generate k8s
	operator-sdk generate crds
	rm -rf build

test::
	go test ./...

vet::
	go vet -composites=false ./...

build::
	go install ./...
