package test_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/test_utils"
)

const TestRootUsername = "root"
const TestUsername = "bob"
const TestPassword = "password"
const WriteTestNetworkId = "N6789"

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

func createTestUserPolicy(token string) certprotos.Policy {
	resources := []*certprotos.Resource{
		{
			Effect:       certprotos.Effect_ALLOW,
			Action:       certprotos.Action_READ,
			ResourceType: certprotos.ResourceType_URI,
			Resource:     "**",
		},
		{
			Effect:       certprotos.Effect_ALLOW,
			Action:       certprotos.Action_READ,
			ResourceType: certprotos.ResourceType_NETWORK_ID,
			Resource:     "**",
		},
		{
			Effect:       certprotos.Effect_ALLOW,
			Action:       certprotos.Action_WRITE,
			ResourceType: certprotos.ResourceType_NETWORK_ID,
			Resource:     WriteTestNetworkId,
		},
	}

	policy := certprotos.Policy{
		Token:     token,
		Resources: &certprotos.ResourceList{Resources: resources},
	}
	return policy
}

func createTestAdminPolicy(token string) certprotos.Policy {
	resources := []*certprotos.Resource{
		{
			Effect:       certprotos.Effect_ALLOW,
			Action:       certprotos.Action_WRITE,
			ResourceType: certprotos.ResourceType_URI,
			Resource:     "**",
		},
		{
			Effect:       certprotos.Effect_ALLOW,
			Action:       certprotos.Action_WRITE,
			ResourceType: certprotos.ResourceType_NETWORK_ID,
			Resource:     "**",
		},
		{
			Effect:       certprotos.Effect_ALLOW,
			Action:       certprotos.Action_WRITE,
			ResourceType: certprotos.ResourceType_TENANT_ID,
			Resource:     "**",
		},
	}
	policy := certprotos.Policy{
		Token:     token,
		Resources: &certprotos.ResourceList{Resources: resources},
	}
	return policy
}
