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

gen::
	operator-sdk generate k8s
	operator-sdk generate crds

test::
	go test ./...

vet::
	go vet -composites=false ./...

build::
	go install ./...
