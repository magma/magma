// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate sh -c "test -e ../../../tools/go/gentargets/main.go && go run ../../../tools/go/gentargets/main.go --pkgbase github.com/facebookincubator/symphony --replace github.com/facebookincubator/symphony=fbc/symphony . || true"
//go:generate go run github.com/google/addlicense -c Facebook -y 2004-present -l bsd ./
