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
	"magma/orc8r/cloud/go/services/certifier"

	"github.com/labstack/echo"
)

func GetOperatorCertificateHandler(c echo.Context) error {
	operator, httpErr := getOperatorForRead(c)
	if httpErr != nil {
		return httpErr
	}
	certificateSNs, err := getCertificateSNs(operator)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to get certificates for operator %s: %s",
			operator.String(), err.Error()))
	}
	return c.JSON(http.StatusOK, certificateSNs)
}

func PostOperatorCertificateHandler(c echo.Context) error {
	operator, httpErr := getOperatorForWrite(c)
	if httpErr != nil {
		return httpErr
	}
	modelCSR := &models.CsrType{}
	if err := c.Bind(modelCSR); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	csr := models.CSRToProto(modelCSR, operator)
	certificate, err := certifier.SignCSR(csr)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to sign CSR for operator %s: %s",
			operator.String(), err.Error()))
	}
	return c.Blob(http.StatusOK, "string", certificate.CertDer)
}

func DeleteOperatorCertificateHandler(c echo.Context) error {
	operator, httpErr := getOperatorForWrite(c)
	if httpErr != nil {
		return httpErr
	}
	certificates, err := certifier.FindCertificates(operator)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to find certificates for %s: %s",
			operator.String(), err.Error()))
	}
	for _, certificateSN := range certificates {
		if err := certifier.RevokeCertificateSN(certificateSN); err != nil {
			return obsidian.HttpError(fmt.Errorf("Failed to revoke certificate %s: %s",
				certificateSN, err.Error()))
		}
	}
	if err := accessd.DeleteOperator(operator); err != nil {
		return obsidian.HttpError(fmt.Errorf("Failed to delete operator %s: %s",
			operator.String(), err.Error()))
	}
	return c.NoContent(http.StatusNoContent)
}
