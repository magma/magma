package akatataipx

import "encoding/json"

// AkaState an internal state of the AKA protocol
type AkaState struct {
	MSK      []byte `json:"msk"`
	Identity string `json:"identity"`
}

// Serialize serializes the given AkaState to string
func (s AkaState) Serialize() string {
	b, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// DeserializeState deserializes the given string to AkaState
// If the state is invalid, a new empty state is returned
func DeserializeState(s string) AkaState {
	var state AkaState
	err := json.Unmarshal([]byte(s), &state)
	if err != nil {
		return AkaState{}
	}
	return state
}

func getAkaState(identity string, msk []byte) AkaState {
	return AkaState{
		Identity: identity,
		MSK:      msk,
	}
}
