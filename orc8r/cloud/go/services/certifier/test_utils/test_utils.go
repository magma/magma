package test_utils

import (
	"testing"

	"magma/orc8r/cloud/go/services/certifier"
	certProtos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/test_utils"
)

func GetCertifierBlobstore(t *testing.T) storage.CertifierStorage {
	fact := test_utils.NewSQLBlobstore(t, storage.CertifierTableBlobstore)
	return storage.NewCertifierBlobstore(fact)
}

func CreateTestUser(t *testing.T, username string, password string) (certProtos.Operator, string) {
	token, err := certifier.GenerateToken(certifier.Personal)
	if err != nil {
		t.Errorf("failed to generate user token %v", token)
	}
	user := certProtos.Operator{
		Username: username,
		Password: []byte(password),
		Tokens:   &certProtos.Operator_TokenList{Token: []string{token}},
	}
	return user, token
}

func CreateUserPolicy(t *testing.T, token string) certProtos.Policy {
	policy := certProtos.Policy{
		Token:     token,
		Effect:    certProtos.Effect_DENY,
		Action:    certProtos.Action_READ,
		Resources: &certProtos.Policy_ResourceList{Resource: []string{"*"}},
	}
	return policy
}
func CreateAdminPolicy(t *testing.T, token string) certProtos.Policy {
	policy := certProtos.Policy{
		Token:     token,
		Effect:    certProtos.Effect_ALLOW,
		Action:    certProtos.Action_WRITE,
		Resources: &certProtos.Policy_ResourceList{Resource: []string{"*"}},
	}
	return policy
}
