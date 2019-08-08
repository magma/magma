/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian"
)

const (
	operatorsRootPath       = obsidian.OperatorsRoot
	operatorsDetailPath     = operatorsRootPath + "/:operator_id"
	operatorEntitiesPath    = operatorsDetailPath + "/entities"
	operatorNetworkPath     = operatorEntitiesPath + "/network/:network_id"
	operatorPermissionsPath = operatorNetworkPath + "/permissions"
	operatorCertificatePath = operatorsDetailPath + "/certificate"
)

// GetObsidianHandlers returns all the handlers for accessd
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// operator_handlers.go
		{
			Path:        operatorsRootPath,
			Methods:     obsidian.GET,
			HandlerFunc: GetOperatorsRootHandler,
		},
		{
			Path:        operatorsRootPath,
			Methods:     obsidian.POST,
			HandlerFunc: PostOperatorsRootHandler,
		},
		{
			Path:        operatorsDetailPath,
			Methods:     obsidian.GET,
			HandlerFunc: GetOperatorsDetailHandler,
		},
		{
			Path:        operatorsDetailPath,
			Methods:     obsidian.DELETE,
			HandlerFunc: DeleteOperatorsDetailHandler,
		},

		// permissions_handlers.go
		{
			Path:        operatorEntitiesPath,
			Methods:     obsidian.POST,
			HandlerFunc: PostOperatorEntityHandler,
		},
		{
			Path:        operatorNetworkPath,
			Methods:     obsidian.DELETE,
			HandlerFunc: DeleteOperatorEntityPermissionHandler,
		},
		{
			Path:        operatorPermissionsPath,
			Methods:     obsidian.GET,
			HandlerFunc: GetOperatorPermissionsHandler,
		},
		{
			Path:        operatorPermissionsPath,
			Methods:     obsidian.PUT,
			HandlerFunc: PutOperatorPermissionsHandler,
		},

		// certificate_handlers.go
		{
			Path:        operatorCertificatePath,
			Methods:     obsidian.GET,
			HandlerFunc: GetOperatorCertificateHandler,
		},
		{
			Path:        operatorCertificatePath,
			Methods:     obsidian.POST,
			HandlerFunc: PostOperatorCertificateHandler,
		},
		{
			Path:        operatorCertificatePath,
			Methods:     obsidian.DELETE,
			HandlerFunc: DeleteOperatorCertificateHandler,
		},
	}
}
