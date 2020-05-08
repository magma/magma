// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/facebookincubator/symphony/jobrunner"
)

func main() {
	jobs := os.Args[1:]
	jobrunner.RunJobOnAllTenants(jobs...)
}
