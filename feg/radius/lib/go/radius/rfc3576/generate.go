/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate go run ../cmd/radius-dict-gen/main.go -package rfc3576 -output generated.go -ref Service-Type:fbc/lib/go/radius/rfc2865 /usr/share/freeradius/dictionary.rfc3576

package rfc3576
