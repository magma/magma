/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian/handlers"
)

const (
	operatorsRootPath       = handlers.OPERATORS_ROOT
	operatorsDetailPath     = operatorsRootPath + "/:operator_id"
	operatorEntitiesPath    = operatorsDetailPath + "/entities"
	operatorNetworkPath     = operatorEntitiesPath + "/network/:network_id"
	operatorPermissionsPath = operatorNetworkPath + "/permissions"
	operatorCertificatePath = operatorsDetailPath + "/certificate"
)

// GetObsidianHandlers returns all the handlers for accessd
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		// operator_handlers.go
		{
			Path:        operatorsRootPath,
			Methods:     handlers.GET,
			HandlerFunc: GetOperatorsRootHandler,
		},
		{
			Path:        operatorsRootPath,
			Methods:     handlers.POST,
			HandlerFunc: PostOperatorsRootHandler,
		},
		{
			Path:        operatorsDetailPath,
			Methods:     handlers.GET,
			HandlerFunc: GetOperatorsDetailHandler,
		},
		{
			Path:        operatorsDetailPath,
			Methods:     handlers.DELETE,
			HandlerFunc: DeleteOperatorsDetailHandler,
		},

		// permissions_handlers.go
		{
			Path:        operatorEntitiesPath,
			Methods:     handlers.POST,
			HandlerFunc: PostOperatorEntityHandler,
		},
		{
			Path:        operatorNetworkPath,
			Methods:     handlers.DELETE,
			HandlerFunc: DeleteOperatorEntityPermissionHandler,
		},
		{
			Path:        operatorPermissionsPath,
			Methods:     handlers.GET,
			HandlerFunc: GetOperatorPermissionsHandler,
		},
		{
			Path:        operatorPermissionsPath,
			Methods:     handlers.PUT,
			HandlerFunc: PutOperatorPermissionsHandler,
		},

		// certificate_handlers.go
		{
			Path:        operatorCertificatePath,
			Methods:     handlers.GET,
			HandlerFunc: GetOperatorCertificateHandler,
		},
		{
			Path:        operatorCertificatePath,
			Methods:     handlers.POST,
			HandlerFunc: PostOperatorCertificateHandler,
		},
		{
			Path:        operatorCertificatePath,
			Methods:     handlers.DELETE,
			HandlerFunc: DeleteOperatorCertificateHandler,
		},
	}
}
