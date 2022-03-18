package sas_helpers

import (
	"encoding/json"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
)

func Build(reqs []*sas.Request) []string {
	byType := [sas.RequestTypeCount][]json.RawMessage{}
	for _, r := range reqs {
		byType[r.Type] = append(byType[r.Type], r.Data)
	}
	payloads := make([]string, 0, len(byType))
	// TODO change this to be deterministic
	for k, v := range byType {
		if len(v) != 0 {
			payloads = append(payloads, toRequest(sas.RequestType(k), v))
		}
	}
	return payloads
}

func toRequest(requestType sas.RequestType, reqs []json.RawMessage) string {
	data := map[string][]json.RawMessage{
		requestType.String(): reqs,
	}
	payload, _ := json.Marshal(data)
	return string(payload)
}
