/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package migration

import (
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func RunCustomPluginMigrations(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, migratedGatewayMetasByNetwork map[string]map[string]MigratedGatewayMeta) error {
	for _, plug := range allPlugins {
		err := plug.RunCustomMigrations(sc, builder, migratedGatewayMetasByNetwork)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
