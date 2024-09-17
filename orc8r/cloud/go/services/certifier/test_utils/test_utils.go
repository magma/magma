package test_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/test_utils"
)

const (
	TestRootUsername        = "root"
	TestUsername            = "bob"
	TestQueryUsername       = "queryUser"
	TestTenantUsername      = "tenantUser"
	TestPassword            = "password"
	WriteTestNetworkId      = "N6789"
	TestTenantId            = int64(0)
	TestTenantNetworkId     = "N6780"
	TestDenyTenantId        = int64(1)
	TestDenyTenantNetworkId = "N7780"
)

func GetCertifierBlobstore(t *testing.T) storage.CertifierStorage {
	fact := test_utils.NewSQLBlobstore(t, storage.CertifierTableBlobstore)
	return storage.NewCertifierBlobstore(fact)
}

func CreateTestUser(t *testing.T, store storage.CertifierStorage, username string, password string, policies []*certprotos.Policy) string {
	user, token := createTestUser(t, username, password)
	err := store.PutUser(username, user)
	assert.NoError(t, err)
	policyList := certprotos.PolicyList{
		Token:    token,
		Policies: policies,
	}
	err = store.PutPolicy(token, &policyList)
	assert.NoError(t, err)
	return token
}

func createTestUser(t *testing.T, username string, password string) (*certprotos.User, string) {
	token, err := certifier.GenerateToken(certifier.Personal)
	assert.NoError(t, err)
	user := certprotos.User{
		Username: username,
		Password: []byte(password),
		Tokens:   &certprotos.TokenList{Tokens: []string{token}},
	}
	return &user, token
}
