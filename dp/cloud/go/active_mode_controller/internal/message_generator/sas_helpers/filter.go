package sas_helpers

import (
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func Filter(pending []*active_mode.Request, requests []*sas.Request) []*sas.Request {
	set := map[string]bool{}
	for _, r := range pending {
		set[r.Payload] = true
	}
	filtered := make([]*sas.Request, 0, len(requests))
	for _, r := range requests {
		if !set[string(r.Data)] {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
