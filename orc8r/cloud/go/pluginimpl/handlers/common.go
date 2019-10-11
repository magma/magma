/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"net/http"
	"reflect"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"

	"github.com/labstack/echo"
)

// GetAndValidatePayload can be used by any model that implements ValidateModel
// Example:
// 	payload, nerr := GetAndValidatePayload(c, &models.DNSConfigRecord{})
//	if nerr != nil {
//		return nil, nerr
//	}
//	record := payload.(*models.DNSConfigRecord)
func GetAndValidatePayload(c echo.Context, model interface{}) (serde.ValidatableModel, *echo.HTTPError) {
	iModel := reflect.New(reflect.TypeOf(model).Elem()).Interface().(serde.ValidatableModel)
	if err := c.Bind(iModel); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	// Run validations specified by the swagger spec
	if err := iModel.ValidateModel(); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	return iModel, nil
}
