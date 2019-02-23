/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package hello implements default (root) path handler which returns JSON
// encoded version
package hello

import (
	"net/http"

	"magma/orc8r/cloud/go/obsidian/config"
	"magma/orc8r/cloud/go/obsidian/handlers"

	"github.com/labstack/echo"
)

// GetObsidianHandlers returns all obsidian handlers for hello
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{
			Path:    "/",
			Methods: handlers.GET,
			HandlerFunc: func(c echo.Context) error {
				return c.JSON(
					http.StatusOK,
					map[string]string{
						"service": config.Product,
						"version": config.Version,
					},
				)
			},
		},
	}
}
