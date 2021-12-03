package protos

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/certifier/storage"

	"github.com/golang/protobuf/proto"
)

func UserFromBlob(blob blobstore.Blob) (User, error) {
	user := User{}
	err := proto.Unmarshal(blob.Value, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func UserToBlob(username string, user *User) (blobstore.Blob, error) {
	marshalledUser, err := proto.Marshal(user)
	if err != nil {
		return blobstore.Blob{}, err
	}
	userBlob := blobstore.Blob{Type: storage.UserType, Key: username, Value: marshalledUser}
	return userBlob, nil
}

func PolicyFromBlob(blob blobstore.Blob) (Policy, error) {
	policy := Policy{}
	err := proto.Unmarshal(blob.Value, &policy)
	if err != nil {
		return policy, err
	}
	return policy, nil

}

func PolicyToBlob(username string, policy *Policy) (blobstore.Blob, error) {
	marshalledPolicy, err := proto.Marshal(policy)
	if err != nil {
		return blobstore.Blob{}, err
	}
	policyBlob := blobstore.Blob{Type: storage.PolicyType, Key: username, Value: marshalledPolicy}
	return policyBlob, nil
}
