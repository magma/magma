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

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/obsidian/models"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/certifier"

	"github.com/labstack/echo"
)

func GetOperatorsRootHandler(c echo.Context) error {
	operators, err := accessd.ListOperators()
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Permission Denied"), http.StatusForbidden)
	}
	operatorIDs := make([]string, len(operators))
	for i, operator := range operators {
		operatorIDs[i] = operator.GetOperator()
	}
	return c.JSON(http.StatusOK, operatorIDs)
}

func PostOperatorsRootHandler(c echo.Context) error {
	createOpRecord := new(models.CreateOperatorRecord)
	if err := c.Bind(createOpRecord); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	operator := identity.NewOperator(string(createOpRecord.Operator))
	allOperators, err := accessd.ListOperators()
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Permission Denied"), http.StatusForbidden)
	}
	for _, listedOperator := range allOperators {
		if operator.HashString() == listedOperator.HashString() {
			return obsidian.HttpError(fmt.Errorf("Operator already exists"), http.StatusBadRequest)
		}
	}
	accessControlEntities := models.ACLToProto(createOpRecord.Entities)
	if err := accessd.SetOperator(operator, accessControlEntities); err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to create operator %s: %s",
			operator.String(), err.Error()))
	}
	csr := models.CSRToProto(createOpRecord.Csr, operator)
	certificate, err := certifier.SignCSR(csr)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to sign CSR: %s", err.Error()))
	}
	return c.Blob(http.StatusCreated, "string", certificate.CertDer)
}

func GetOperatorsDetailHandler(c echo.Context) error {
	operator, httpErr := getOperatorForRead(c)
	if httpErr != nil {
		return httpErr
	}
	certificateSNs, err := getCertificateSNs(operator)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to get certificates for operator %s: %s",
			operator.String(), err.Error()))
	}
	accessControlEntities, err := accessd.GetOperatorACL(operator)
	var accessControlList []*accessprotos.AccessControl_Entity
	for _, entity := range accessControlEntities {
		accessControlList = append(accessControlList, entity)
	}
	acl := models.ACLFromProto(accessControlList)
	operatorRecord := &models.OperatorRecord{
		CertificateSns: certificateSNs,
		Entities:       acl,
	}
	return c.JSON(http.StatusOK, operatorRecord)
}

func DeleteOperatorsDetailHandler(c echo.Context) error {
	operator, httpErr := getOperatorForWrite(c)
	if httpErr != nil {
		return httpErr
	}
	certificateSNs, err := certifier.FindCertificates(operator)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to find certificates for %s: %s",
			operator.String(), err.Error()))
	}
	for _, certificate := range certificateSNs {
		certifier.RevokeCertificateSN(certificate)
	}
	err = accessd.DeleteOperator(operator)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to delete operator %s: %s",
			operator.String(), err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}
