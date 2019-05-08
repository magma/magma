package serde

import (
	"encoding/json"

	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
)

type GatewayStatusSerde struct{}

func (*GatewayStatusSerde) GetDomain() string {
	return state.SerdeDomain
}

func (s *GatewayStatusSerde) GetType() string {
	return "gw_state"
}

func (s *GatewayStatusSerde) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (s *GatewayStatusSerde) Deserialize(in []byte) (interface{}, error) {
	response := checkind_models.GatewayStatus{}
	err := json.Unmarshal(in, &response)
	return response, err
}
