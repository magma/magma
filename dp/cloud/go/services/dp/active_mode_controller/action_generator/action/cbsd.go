package action

import (
	"golang.org/x/exp/slices"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/eirp"
	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/grant"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func SetAvailableFrequences(cbsd *storage.DBCbsd) Action {
	calc := eirp.NewCalculator(cbsd)
	availableFrequencies := grant.CalcAvailableFrequencies(cbsd.Channels, calc)
	cbsd.AvailableFrequencies = availableFrequencies

	data := &storage.DBCbsd{
		Id:                   cbsd.Id,
		AvailableFrequencies: availableFrequencies,
	}
	mask := db.NewIncludeMask("available_frequencies")
	return &UpdateCbsd{Data: data, Mask: mask}
}

func UnsetGrantFrequency(cbsd *storage.DetailedCbsd, gt *storage.DetailedGrant) Action {
	newFrequencies := grant.UnsetGrantFrequency(cbsd.Cbsd, gt.Grant)
	if slices.Equal(newFrequencies, cbsd.Cbsd.AvailableFrequencies) {
		return nil
	}
	cbsd.Cbsd.AvailableFrequencies = newFrequencies

	data := &storage.DBCbsd{
		Id:                   cbsd.Cbsd.Id,
		AvailableFrequencies: newFrequencies,
	}
	mask := db.NewIncludeMask("available_frequencies")
	return &UpdateCbsd{Data: data, Mask: mask}
}

func RemoveIdleGrants(cbsd *storage.DetailedCbsd) []Action {
	var actions []Action
	var notIdleGrants []*storage.DetailedGrant

	for _, gt := range cbsd.Grants {
		if gt.GrantState.Name.String != "idle" {
			notIdleGrants = append(notIdleGrants, gt)
			continue
		}

		action := UnsetGrantFrequency(cbsd, gt)
		if action != nil {
			actions = append(actions, action)
		}

		actDelGrat := &DeleteGrant{Id: gt.Grant.Id.Int64}
		actions = append(actions, actDelGrat)
	}

	cbsd.Grants = notIdleGrants
	return actions
}
