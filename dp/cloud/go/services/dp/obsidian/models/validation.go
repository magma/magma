package models

import (
	"context"
	"fmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
)

func (m *MutableCbsd) ValidateModel(_ context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	var res []error
	if !*m.GrantRedundancy && *m.CarrierAggregationEnabled {
		err := fmt.Errorf("grant_redundancy cannot be set to false when carrier_aggregation_enabled is enabled")
		res = append(res, err)
	}
	if m.Capabilities.MaxIbwMhz < m.FrequencyPreferences.BandwidthMhz {
		err := fmt.Errorf("max_ibw_mhz cannot be less than bandwidth_mhz")
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
