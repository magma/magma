package test_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/certifier"
	certProtos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/test_utils"
)

const TestRootUsername = "root"
const TestUsername = "bob"

func GetCertifierBlobstore(t *testing.T) storage.CertifierStorage {
	fact := test_utils.NewSQLBlobstore(t, storage.CertifierTableBlobstore)
	return storage.NewCertifierBlobstore(fact)
}

func CreateTestAdmin(t *testing.T, store storage.CertifierStorage) string {
	user, token := createTestUser(t, TestRootUsername, "password")
	err := store.PutUser(TestRootUsername, &user)
	assert.NoError(t, err)
	policy := createTestAdminPolicy(token)
	err = store.PutPolicy(token, &policy)
	assert.NoError(t, err)
	return token
}

func CreateTestUser(t *testing.T, store storage.CertifierStorage) string {
	user, token := createTestUser(t, TestUsername, "password")
	err := store.PutUser(TestUsername, &user)
	assert.NoError(t, err)
	policy := createTestUserPolicy(token)
	err = store.PutPolicy(token, &policy)
	assert.NoError(t, err)
	return token
}

func createTestUser(t *testing.T, username string, password string) (certProtos.User, string) {
	token, err := certifier.GenerateToken(certifier.Personal)
	assert.NoError(t, err)
	user := certProtos.User{
		Username: username,
		Password: []byte(password),
		Tokens:   &certProtos.TokenList{Token: []string{token}},
	}
	return user, token
}

func createTestUserPolicy(token string) certProtos.Policy {
	policy := certProtos.Policy{
		Token:     token,
		Effect:    certProtos.Effect_ALLOW,
		Action:    certProtos.Action_READ,
		Resources: &certProtos.ResourceList{Resource: []string{"/**"}},
	}
	return policy
}

func createTestAdminPolicy(token string) certProtos.Policy {
	policy := certProtos.Policy{
		Token:     token,
		Effect:    certProtos.Effect_ALLOW,
		Action:    certProtos.Action_WRITE,
		Resources: &certProtos.ResourceList{Resource: []string{"/**"}},
	}
	return policy
}
