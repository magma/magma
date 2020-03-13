/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package blobstore

import "magma/orc8r/lib/go/definitions"

var (
	SQLDriver      = definitions.GetEnvWithDefault("SQL_DRIVER", "sqlite3")
	DatabaseSource = definitions.GetEnvWithDefault("DATABASE_SOURCE", ":memory:")
)
