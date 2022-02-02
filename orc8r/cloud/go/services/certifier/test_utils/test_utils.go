package test_utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

const (
	TestRootUsername        = "root"
	TestUsername            = "bob"
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

func CreateTestAdmin(t *testing.T, store storage.CertifierStorage) string {
	user, token := createTestUser(t, TestRootUsername, TestPassword)
	err := store.PutUser(TestRootUsername, &user)
	assert.NoError(t, err)
	policy := createTestAdminPolicy(token)
	err = store.PutPolicy(token, &policy)
	assert.NoError(t, err)
	return token
}

func CreateTestUser(t *testing.T, store storage.CertifierStorage) string {
	user, token := createTestUser(t, TestUsername, TestPassword)
	err := store.PutUser(TestUsername, &user)
	assert.NoError(t, err)
	policy := createTestUserPolicy(token)
	err = store.PutPolicy(token, &policy)
	assert.NoError(t, err)
	return token
}

func CreateTestTenantUser(t *testing.T, store storage.CertifierStorage) string {
	user, token := createTestTenantUser(t, TestTenantUsername, TestPassword)
	err := store.PutUser(TestTenantUsername, &user)
	assert.NoError(t, err)
	tenants.CreateTenant(context.Background(), TestTenantId, &protos.Tenant{
		Name:     string(TestTenantId),
		Networks: []string{TestTenantNetworkId},
	})
	policy := createTestTenantUserPolicy(token)
	err = store.PutPolicy(token, &policy)
	assert.NoError(t, err)
	return token
}

func createTestUser(t *testing.T, username string, password string) (certprotos.User, string) {
	token, err := certifier.GenerateToken(certifier.Personal)
	assert.NoError(t, err)
	user := certprotos.User{
		Username: username,
		Password: []byte(password),
		Tokens:   &certprotos.TokenList{Tokens: []string{token}},
	}
	return user, token
}

func createTestTenantUser(t *testing.T, username string, password string) (certprotos.User, string) {
	token, err := certifier.GenerateToken(certifier.Personal)
	assert.NoError(t, err)
	user := certprotos.User{
		Username: username,
		Password: []byte(password),
		Tokens:   &certprotos.TokenList{Tokens: []string{token}},
	}
	return user, token
}

func createTestUserPolicy(token string) certprotos.PolicyList {
	policies := []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_READ,
			Resource: &certprotos.Policy_Path{Path: &certprotos.PathResource{Path: "**"}},
		},
		{
			Effect:   certprotos.Effect_DENY,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Network{Network: &certprotos.NetworkResource{Networks: []string{WriteTestNetworkId}}},
		},
	}

	policy := certprotos.PolicyList{
		Token:    token,
		Policies: policies,
	}
	return policy
}

func createTestAdminPolicy(token string) certprotos.PolicyList {
	policies := []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Path{Path: &certprotos.PathResource{Path: "**"}},
		},
	}
	policy := certprotos.PolicyList{
		Token:    token,
		Policies: policies,
	}
	return policy
}

func createTestTenantUserPolicy(token string) certprotos.PolicyList {
	policies := []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Tenant{Tenant: &certprotos.TenantResource{Tenants: []int64{TestTenantId}}},
		},
	}
	policy := certprotos.PolicyList{
		Token:    token,
		Policies: policies,
	}
	return policy
}
