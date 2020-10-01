package calculations

import (
	"fmt"
	"github.com/golang/glog"
	"magma/cwf/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"
)

// AuthenticationsCalculation holds the parameters needed to run an authentication
// query and the registered prometheus gauge that the resulting value should be stored in
type AuthenticationsCalculation struct {
	CalculationParams
}

// Calculate returns the number of authentications over the past X days segmented
// by result code and networkID
func (x *AuthenticationsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]Result, error) {
	glog.Infof("Calculating Authentications. Days: %d", x.Days)

	query := fmt.Sprintf(`sum(increase(eap_auth[%dd])) by (code, %s)`, x.Days, metrics.NetworkLabelName)

	vec, err := query_api.QueryPrometheusVector(prometheusClient, query)
	if err != nil {
		return nil, fmt.Errorf("user Consumption query error: %s", err)
	}

	results := makeVectorResults(vec, x.Labels, x.Name)
	registerResults(x.CalculationParams, results)

	return results, nil
}
