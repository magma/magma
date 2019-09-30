/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate go run ../cmd/radius-dict-gen/main.go -package rfc5176 -output generated.go -ref Error-Cause:fbc/lib/go/radius/rfc3576 /usr/share/freeradius/dictionary.rfc5176

package rfc5176
