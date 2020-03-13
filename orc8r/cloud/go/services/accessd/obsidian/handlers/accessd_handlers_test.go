/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers_test

import (
	"crypto/x509"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/obsidian/access"
	accessTests "magma/orc8r/cloud/go/obsidian/access/tests"
	"magma/orc8r/cloud/go/services/accessd/obsidian/handlers"
	"magma/orc8r/cloud/go/services/accessd/obsidian/models"
	"magma/orc8r/cloud/go/services/accessd/test_init"
	"magma/orc8r/cloud/go/test_utils"
	securityCert "magma/orc8r/lib/go/security/cert"
	certifierTestUtils "magma/orc8r/lib/go/security/csr"
)

const (
	operator1ID         = models.OperatorID("op1")
	operator2ID         = models.OperatorID("op2")
	adminID             = models.OperatorID("admin")
	network1ID          = models.NetworkID("net1")
	network2ID          = models.NetworkID("net2")
	network3ID          = models.NetworkID("net3")
	testAdminOperatorID = "Obsidian_Unit_Test_Admin_Operator"
)

func testInit(t *testing.T) (string, map[models.OperatorID]models.Certificate, map[models.OperatorID]models.ACLType) {
	test_init.StartTestService(t)
	testOperatorSerialNumber := accessTests.StartMockAccessControl(t, testAdminOperatorID)
	certificates := make(map[models.OperatorID]models.Certificate)
	acls := make(map[models.OperatorID]models.ACLType)
	initializeOperators(t, certificates, acls)
	return testOperatorSerialNumber, certificates, acls
}

func cleanup(t *testing.T) {
	test_utils.DropTableFromSharedTestDB(t, "access_control_blobstore")
	test_utils.DropTableFromSharedTestDB(t, "certificate_info_blobstore")
}

func TestListOperators(t *testing.T) {
	defer cleanup(t)
	testInit(t)
	operatorsListExpected := []string{
		string(operator1ID),
		string(operator2ID),
		string(adminID),
	}
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := handlers.GetOperatorsRootHandler(c)
	assert.NoError(t, err)
	actualSet, err := getOperatorSetFromJSON(rec.Body.Bytes())
	assert.NoError(t, err)
	// We just want to check that our operators are registered,
	// if there are more that's fine
	for _, opid := range operatorsListExpected {
		_, ok := actualSet[opid]
		assert.True(t, ok, "Operator %s not found", opid)
	}
	// Assert that there are no empty strings
	_, ok := actualSet[""]
	assert.False(t, ok, "Empty string found in operator set")
}

func TestGetOperatorsDetail(t *testing.T) {
	defer cleanup(t)
	testOperatorSN, certificates, acls := testInit(t)
	expectedSNs := []models.CertificateSn{
		certToSerialNumber(t, certificates[operator1ID]),
	}
	expectedACL := acls[operator1ID]
	expectedRecord := models.OperatorRecord{
		CertificateSns: expectedSNs,
		Entities:       expectedACL,
	}
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err := handlers.GetOperatorsDetailHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assertOperatorRecordResponse(t, &expectedRecord, rec.Body.String())
}

func TestGetOperatorsDetailSpecificPermission(t *testing.T) {
	defer cleanup(t)
	_, certificates, acls := testInit(t)
	expectedSNs := []models.CertificateSn{
		certToSerialNumber(t, certificates[operator1ID]),
	}
	expectedACL := acls[operator1ID]
	expectedRecord := models.OperatorRecord{
		CertificateSns: expectedSNs,
		Entities:       expectedACL,
	}
	serialNum := certToSerialNumber(t, certificates[operator2ID])
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, string(serialNum))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err := handlers.GetOperatorsDetailHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assertOperatorRecordResponse(t, &expectedRecord, rec.Body.String())
}

func TestGetOperatorsDetailRestrictedPermission(t *testing.T) {
	defer cleanup(t)
	_, certificates, _ := testInit(t)
	serialNum := certToSerialNumber(t, certificates[operator1ID])
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, string(serialNum))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator2ID))
	err := handlers.GetOperatorsDetailHandler(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code=403")
}

func TestDeleteOperators(t *testing.T) {
	defer cleanup(t)
	testOperatorSN, _, _ := testInit(t)
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err := handlers.DeleteOperatorsDetailHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.GetOperatorsDetailHandler(c)
	assert.NoError(t, err)
	expected := models.OperatorRecord{
		CertificateSns: []models.CertificateSn{},
	}
	actual := models.OperatorRecord{}
	err = json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)
}

func TestDeleteOperatorRestrictedPermission(t *testing.T) {
	defer cleanup(t)
	_, certificates, _ := testInit(t)
	serialNum := certToSerialNumber(t, certificates[operator1ID])
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, string(serialNum))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator2ID))
	err := handlers.DeleteOperatorsDetailHandler(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code=403")
}

// Tests adding and deleting an entity as well as getting permissions
func TestAddDeleteEntityPermissionsCheck(t *testing.T) {
	defer cleanup(t)
	testOperatorSN, _, _ := testInit(t)
	// Add network 4 to operator 1's ACL
	newNetworkID := models.NetworkID("net4")
	newEntity := &models.ACLEntity{
		EntityType: models.ACLEntityEntityTypeNETWORK,
		NetworkID:  newNetworkID,
		Permissions: models.PermissionsMask{
			models.PermissionTypeREAD,
			models.PermissionTypeWRITE,
		},
	}
	newEntityBinary, err := newEntity.MarshalBinary()
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(string(newEntityBinary)))
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.PostOperatorEntityHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Check that the permissions were set properly
	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id", "network_id")
	c.SetParamValues(string(operator1ID), string(newNetworkID))
	err = handlers.GetOperatorPermissionsHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var permissions models.PermissionsMask
	err = json.Unmarshal(rec.Body.Bytes(), &permissions)
	assert.NoError(t, err)
	expected := models.PermissionsMask{models.PermissionTypeREAD, models.PermissionTypeWRITE}
	assert.ElementsMatch(t, expected, permissions)

	// Delete entity
	req = httptest.NewRequest(echo.DELETE, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id", "network_id")
	c.SetParamValues(string(operator1ID), string(newNetworkID))
	err = handlers.DeleteOperatorEntityPermissionHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Check that permissions were removed
	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id", "network_id")
	c.SetParamValues(string(operator1ID), string(newNetworkID))
	err = handlers.GetOperatorPermissionsHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &permissions)
	assert.NoError(t, err)
	assert.True(t, permissions[0] == models.PermissionTypeNONE && permissions[1] == models.PermissionTypeNONE)
}

func TestPutPermissions(t *testing.T) {
	defer cleanup(t)
	testOperatorSN, _, _ := testInit(t)
	// Update permissions of operator 1 for network 2
	newPermissions := models.PermissionsMask{
		models.PermissionTypeREAD,
		models.PermissionTypeWRITE,
	}
	newPermissionsBytes, err := json.Marshal(newPermissions)
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(string(newPermissionsBytes)))
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id", "network_id")
	c.SetParamValues(string(operator1ID), string(network2ID))
	err = handlers.PutOperatorPermissionsHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check that permissions were updated properly
	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id", "network_id")
	c.SetParamValues(string(operator1ID), string(network2ID))
	err = handlers.GetOperatorPermissionsHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var permissions models.PermissionsMask
	err = json.Unmarshal(rec.Body.Bytes(), &permissions)
	assert.NoError(t, err)
	expected := models.PermissionsMask{models.PermissionTypeREAD, models.PermissionTypeWRITE}
	assert.ElementsMatch(t, expected, permissions)
}

func TestGetCertificate(t *testing.T) {
	defer cleanup(t)
	testOperatorSN, certificates, _ := testInit(t)
	op1CertSN := certToSerialNumber(t, certificates[operator1ID])
	op1GetCertExpectedResponse, err := json.Marshal([]string{
		string(op1CertSN),
	})
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.GetOperatorCertificateHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, string(op1GetCertExpectedResponse), rec.Body.String())
}

func TestPostDeleteCertificate(t *testing.T) {
	defer cleanup(t)
	testOperatorSN, certificates, _ := testInit(t)
	csr, err := certifierTestUtils.CreateCSR(
		time.Duration(int64(time.Hour*24*365)),
		string(operator1ID),
		string(operator1ID),
	)
	assert.NoError(t, err)
	modelCSR := models.CSRFromProto(csr)
	modelCSRBytes, err := modelCSR.MarshalBinary()
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(string(modelCSRBytes)))
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.PostOperatorCertificateHandler(c)
	assert.NoError(t, err)
	certificate := models.Certificate(rec.Body.Bytes())
	certificateSN := certToSerialNumber(t, certificate)
	oldCertSN := certToSerialNumber(t, certificates[operator1ID])
	expectedSNs := []models.CertificateSn{
		oldCertSN,
		certificateSN,
	}

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.GetOperatorCertificateHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var actualSNs []models.CertificateSn
	err = json.Unmarshal(rec.Body.Bytes(), &actualSNs)
	assert.ElementsMatch(t, expectedSNs, actualSNs)

	req = httptest.NewRequest(echo.DELETE, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.DeleteOperatorCertificateHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(access.CLIENT_CERT_SN_KEY, testOperatorSN)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("operator_id")
	c.SetParamValues(string(operator1ID))
	err = handlers.GetOperatorCertificateHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "[]", rec.Body.String())
}

// Helpers

func assertOperatorRecordResponse(t *testing.T, expectedRecord *models.OperatorRecord, response string) {
	var actualRecord models.OperatorRecord
	err := actualRecord.UnmarshalBinary([]byte(response))
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedRecord.CertificateSns, actualRecord.CertificateSns)
	assert.ElementsMatch(t, expectedRecord.Entities, actualRecord.Entities)
}

func certToSerialNumber(t *testing.T, certificate models.Certificate) models.CertificateSn {
	cert, err := x509.ParseCertificate(certificate)
	assert.NoError(t, err)
	return models.CertificateSn(securityCert.SerialToString(cert.SerialNumber))
}

func getOperatorSetFromJSON(jsonstr []byte) (map[string]bool, error) {
	var operatorList []string
	operatorSet := make(map[string]bool)
	err := json.Unmarshal(jsonstr, &operatorList)
	if err != nil {
		return operatorSet, err
	}
	for _, item := range operatorList {
		operatorSet[item] = true
	}
	return operatorSet, nil
}

// Initializers

func initializeOperators(t *testing.T, certificates map[models.OperatorID]models.Certificate, acls map[models.OperatorID]models.ACLType) {
	initializeOperator1(t, certificates, acls)
	initializeOperator2(t, certificates, acls)
	initializeAdminOperator(t, certificates, acls)
}

func initializeOperator1(t *testing.T, certificates map[models.OperatorID]models.Certificate, acls map[models.OperatorID]models.ACLType) {
	operator1CSR, err := certifierTestUtils.CreateCSR(
		time.Duration(int64(time.Hour*24*365)),
		string(operator1ID),
		string(operator1ID),
	)
	assert.NoError(t, err)
	operator1ACL := models.ACLType{
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeNETWORK,
			NetworkID:  network1ID,
			Permissions: models.PermissionsMask{
				models.PermissionTypeREAD,
				models.PermissionTypeWRITE,
			},
		},
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeNETWORK,
			NetworkID:  network2ID,
			Permissions: models.PermissionsMask{
				models.PermissionTypeREAD,
				models.PermissionTypeNONE,
			},
		},
	}
	acls[operator1ID] = operator1ACL
	createOperator1 := models.CreateOperatorRecord{
		Csr:      models.CSRFromProto(operator1CSR),
		Entities: operator1ACL,
		Operator: operator1ID,
	}
	certificates[operator1ID] = createOperator(t, createOperator1)
}

func initializeOperator2(t *testing.T, certificates map[models.OperatorID]models.Certificate, acls map[models.OperatorID]models.ACLType) {
	operator2CSR, err := certifierTestUtils.CreateCSR(
		time.Duration(int64(time.Hour*24*365)),
		string(operator2ID),
		string(operator2ID),
	)
	assert.NoError(t, err)
	operator2ACL := models.ACLType{
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeNETWORK,
			NetworkID:  network2ID,
			Permissions: models.PermissionsMask{
				models.PermissionTypeREAD,
				models.PermissionTypeWRITE,
			},
		},
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeNETWORK,
			NetworkID:  network3ID,
			Permissions: models.PermissionsMask{
				models.PermissionTypeWRITE,
				models.PermissionTypeNONE,
			},
		},
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeOPERATOR,
			OperatorID: operator1ID,
			Permissions: models.PermissionsMask{
				models.PermissionTypeREAD,
				models.PermissionTypeNONE,
			},
		},
	}
	acls[operator2ID] = operator2ACL
	createOperator2 := models.CreateOperatorRecord{
		Csr:      models.CSRFromProto(operator2CSR),
		Entities: operator2ACL,
		Operator: operator2ID,
	}
	certificates[operator2ID] = createOperator(t, createOperator2)
}

func initializeAdminOperator(t *testing.T, certificates map[models.OperatorID]models.Certificate, acls map[models.OperatorID]models.ACLType) {
	adminCSR, err := certifierTestUtils.CreateCSR(
		time.Duration(int64(time.Hour*24*365)),
		string(adminID),
		string(adminID),
	)
	assert.NoError(t, err)
	adminACL := models.ACLType{
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeNETWORKWILDCARD,
			Permissions: models.PermissionsMask{
				models.PermissionTypeREAD,
				models.PermissionTypeWRITE,
			},
		},
		&models.ACLEntity{
			EntityType: models.ACLEntityEntityTypeOPERATORWILDCARD,
			Permissions: models.PermissionsMask{
				models.PermissionTypeREAD,
				models.PermissionTypeWRITE,
			},
		},
	}
	acls[adminID] = adminACL
	createAdmin := models.CreateOperatorRecord{
		Csr:      models.CSRFromProto(adminCSR),
		Entities: adminACL,
		Operator: adminID,
	}
	certificates[adminID] = createOperator(t, createAdmin)
}

func createOperator(t *testing.T, createRecord models.CreateOperatorRecord) models.Certificate {
	createRecordBinary, err := createRecord.MarshalBinary()
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(string(createRecordBinary)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = handlers.PostOperatorsRootHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	response := rec.Body.String()
	return []byte(response)
}
