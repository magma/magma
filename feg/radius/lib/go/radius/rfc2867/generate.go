/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate go run ../cmd/radius-dict-gen/main.go -package rfc2867 -output generated.go -ref Acct-Status-Type:fbc/lib/go/radius/rfc2866 /usr/share/freeradius/dictionary.rfc2867

package rfc2867
