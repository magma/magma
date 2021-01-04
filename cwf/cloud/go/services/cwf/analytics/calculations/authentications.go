package calculations

import (
	"fmt"

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
)

// AuthenticationsCalculation holds the parameters needed to run an authentication
// query and the registered prometheus gauge that the resulting value should be stored in
type AuthenticationsCalculation struct {
	calculations.BaseCalculation
}

// Calculate returns the number of authentications over the past X days segmented
// by result code and networkID
func (x *AuthenticationsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.Infof("Calculating Authentications. Days: %d", x.Days)

	query := fmt.Sprintf(`sum(increase(eap_auth[%dd])) by (code, %s)`, x.Days, metrics.NetworkLabelName)

	vec, err := query_api.QueryPrometheusVector(prometheusClient, query)
	if err != nil {
		return nil, fmt.Errorf("user Consumption query error: %s", err)
	}

	results := calculations.MakeVectorResults(vec, x.Labels, x.Name)
	return results, nil
}
