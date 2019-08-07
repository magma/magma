/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/obsidian/models"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"

	"github.com/labstack/echo"
)

func PostOperatorEntityHandler(c echo.Context) error {
	operator, httpErr := getOperatorForWrite(c)
	if httpErr != nil {
		return httpErr
	}
	aclEntity := &models.ACLEntity{}
	if err := c.Bind(aclEntity); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	accessControlEntity := models.ACLEntityToProto(aclEntity)
	err := accessd.UpdateOperator(operator, []*accessprotos.AccessControl_Entity{accessControlEntity})
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to add permissions for %s: %s",
			operator.String(), err.Error()))
	}
	return c.NoContent(http.StatusCreated)
}

func DeleteOperatorEntityPermissionHandler(c echo.Context) error {
	operator, httpErr := getOperatorForWrite(c)
	if httpErr != nil {
		return httpErr
	}
	aclmap, err := accessd.GetOperatorACL(operator)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to get ACL for %s: %s",
			operator.String(), err.Error()))
	}
	network, httpErr := getNetwork(c)
	if httpErr != nil {
		return httpErr
	}
	delete(aclmap, network.HashString())
	var acl []*accessprotos.AccessControl_Entity
	for _, aclEntity := range aclmap {
		acl = append(acl, aclEntity)
	}
	if err := accessd.SetOperator(operator, acl); err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to remove permissions for %s: %s",
			operator.String(), err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}

func GetOperatorPermissionsHandler(c echo.Context) error {
	operator, httpErr := getOperatorForRead(c)
	if httpErr != nil {
		return httpErr
	}
	network, httpErr := getNetwork(c)
	if httpErr != nil {
		return httpErr
	}
	permissions, err := accessd.GetPermissions(operator, network)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to get %s permissions for %s: %s",
			operator.String(), network.String(), err.Error()))
	}
	permissionsMask := models.PermissionsMaskFromProto(permissions)
	return c.JSON(http.StatusOK, permissionsMask)
}

func PutOperatorPermissionsHandler(c echo.Context) error {
	operator, httpErr := getOperatorForWrite(c)
	if httpErr != nil {
		return httpErr
	}
	network, httpErr := getNetwork(c)
	if httpErr != nil {
		return httpErr
	}
	permissionsMask := models.PermissionsMask{}
	if err := c.Bind(&permissionsMask); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	accessControlEntity := &accessprotos.AccessControl_Entity{
		Id:          network,
		Permissions: models.PermissionsMaskToProto(permissionsMask),
	}
	accessControlList := []*accessprotos.AccessControl_Entity{
		accessControlEntity,
	}
	if err := accessd.UpdateOperator(operator, accessControlList); err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to update permissions for %s: %s",
			operator.String(), err.Error()))
	}
	return c.NoContent(http.StatusOK)
}
