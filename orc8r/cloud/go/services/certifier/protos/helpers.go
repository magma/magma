package protos

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/certifier/obsidian/models"
)

const (
	// UserType is the type of CertInfo used in blobstore type fields.
	UserType = "user"
	// PolicyType is the type of policy used in blobstore type fileds
	PolicyType = "policy"
)

func UserFromBlob(blob blobstore.Blob) (*User, error) {
	user := &User{}
	err := proto.Unmarshal(blob.Value, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) UserToBlob(username string) (blobstore.Blob, error) {
	marshalledUser, err := proto.Marshal(u)
	if err != nil {
		return blobstore.Blob{}, err
	}
	userBlob := blobstore.Blob{Type: UserType, Key: username, Value: marshalledUser}
	return userBlob, nil
}

func PolicyFromBlob(blob blobstore.Blob) (*PolicyList, error) {
	policy := &PolicyList{}
	err := proto.Unmarshal(blob.Value, policy)
	if err != nil {
		return policy, err
	}
	return policy, nil
}

func (p *PolicyList) PolicyToBlob(username string) (blobstore.Blob, error) {
	marshalledPolicy, err := proto.Marshal(p)
	if err != nil {
		return blobstore.Blob{}, err
	}
	policyBlob := blobstore.Blob{Type: PolicyType, Key: username, Value: marshalledPolicy}
	return policyBlob, nil
}

func PolicyListProtoToModel(policyLists []*PolicyList) []models.PolicyList {
	var policyListsModels []models.PolicyList
	for _, pl := range policyLists {
		policiesModel := policiesProtoToModel(pl.Policies)
		policyListsModel := models.PolicyList{
			Token:    &pl.Token,
			Policies: policiesModel,
		}
		policyListsModels = append(policyListsModels, policyListsModel)
	}
	return policyListsModels
}

func PoliciesModelToProto(policies *models.Policies) ([]*Policy, error) {
	policyProtos := make([]*Policy, len(*policies))
	for i, policyModel := range *policies {
		policyProto := &Policy{
			Effect: matchEffect(policyModel.Effect),
			Action: matchAction(policyModel.Action),
		}
		if err := setResource(policyModel, policyProto); err != nil {
			return nil, err
		}
		policyProtos[i] = policyProto
	}
	return policyProtos, nil
}

func policiesProtoToModel(policies []*Policy) models.Policies {
	var policiesModel models.Policies
	for _, p := range policies {
		policyModel := &models.Policy{
			Action: p.Action.String(),
			Effect: p.Effect.String(),
		}
		if path := p.GetPath(); path != nil {
			policyModel.Path = path.Path
		}
		if nid := p.GetNetwork(); nid != nil {
			policyModel.ResourceIDs = nid.Networks
		}
		if tid := p.GetTenant(); tid != nil {
			var tidStr []string
			for _, i := range tid.Tenants {
				tidStr = append(tidStr, strconv.FormatInt(i, 10))
			}
			policyModel.ResourceIDs = tidStr
		}
		policiesModel = append(policiesModel, policyModel)
	}
	return policiesModel
}

func convertTenantResourceIDs(ids []string) ([]int64, error) {
	var intIDs []int64
	for _, i := range ids {
		j, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			return []int64{}, fmt.Errorf("failed to convert tenant IDs to integers: %w", err)
		}
		intIDs = append(intIDs, j)
	}
	return intIDs, nil
}

func matchEffect(rawEffect string) Effect {
	switch rawEffect {
	case Effect_DENY.String():
		return Effect_DENY
	case Effect_ALLOW.String():
		return Effect_ALLOW
	default:
		return Effect_UNKNOWN
	}
}

func matchAction(rawAction string) Action {
	switch rawAction {
	case Action_READ.String():
		return Action_READ
	case Action_WRITE.String():
		return Action_WRITE
	default:
		return Action_NONE
	}
}

// setResource uses the resource value in the policy model and sets the
// resource based on its type in the policy proto
func setResource(policyModel *models.Policy, policyProto *Policy) error {
	switch policyModel.ResourceType {
	case models.PolicyResourceTypeNETWORKID:
		policyProto.Resource = &Policy_Network{Network: &NetworkResource{Networks: policyModel.ResourceIDs}}
	case models.PolicyResourceTypeTENANTID:
		tenantIDs, err := convertTenantResourceIDs(policyModel.ResourceIDs)
		if err != nil {
			return err
		}
		policyProto.Resource = &Policy_Tenant{Tenant: &TenantResource{Tenants: tenantIDs}}
	default:
		policyProto.Resource = &Policy_Path{Path: &PathResource{Path: policyModel.Path}}
	}
	return nil
}
