/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/lib/go/merrors"
)

func listNetworks(c echo.Context) error {
	networks, err := configurator.ListNetworkIDs(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if networks == nil {
		networks = []string{}
	}
	return c.JSON(http.StatusOK, networks)
}

func registerNetwork(c echo.Context) error {
	payload, nerr := GetAndValidatePayload(c, &models.Network{})
	if nerr != nil {
		return nerr
	}
	network := payload.(*models.Network).ToConfiguratorNetwork()
	createdNetworks, err := configurator.CreateNetworks(
		c.Request().Context(),
		[]configurator.Network{network},
		serdes.Network,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, createdNetworks[0].ID)
}

func getNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	network, err := configurator.LoadNetwork(c.Request().Context(), networkID, true, true, serdes.Network)
	if err == merrors.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	ret := (&models.Network{}).FromConfiguratorNetwork(network)
	return c.JSON(http.StatusOK, ret)
}

func updateNetwork(c echo.Context) error {
	network, nerr := GetAndValidatePayload(c, &models.Network{})
	if nerr != nil {
		return nerr
	}
	update := network.(*models.Network).ToUpdateCriteria()
	err := configurator.UpdateNetworks(c.Request().Context(), []configurator.NetworkUpdateCriteria{update}, serdes.Network)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteNetwork(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	err := configurator.DeleteNetwork(c.Request().Context(), networkID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func CreateDNSRecord(c echo.Context) error {
	networkID, domain, nerr := getNetworkIDAndDomain(c)
	if nerr != nil {
		return nerr
	}

	record, nerr := getRecordAndValidate(c, domain)
	if nerr != nil {
		return nerr
	}

	reqCtx := c.Request().Context()
	dnsConfig, nerr := getExistingDNSConfig(reqCtx, networkID)
	if nerr != nil {
		return nerr
	}

	// check the domain is not already registered
	for _, record := range dnsConfig.Records {
		if record.Domain == domain {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("A record with domain:%s already exists", domain))
		}
	}

	dnsConfig.Records = append(dnsConfig.Records, record)
	nerr = updateDNSConfig(reqCtx, networkID, dnsConfig)
	if nerr != nil {
		return nerr
	}
	return c.JSON(http.StatusCreated, domain)
}

func UpdateDNSRecord(c echo.Context) error {
	networkID, domain, nerr := getNetworkIDAndDomain(c)
	if nerr != nil {
		return nerr
	}

	record, nerr := getRecordAndValidate(c, domain)
	if nerr != nil {
		return nerr
	}

	reqCtx := c.Request().Context()
	dnsConfig, nerr := getExistingDNSConfig(reqCtx, networkID)
	if nerr != nil {
		return nerr
	}

	for i, existingRecord := range dnsConfig.Records {
		if existingRecord.Domain == domain {
			dnsConfig.Records[i] = record
			nerr = updateDNSConfig(reqCtx, networkID, dnsConfig)
			if nerr != nil {
				return nerr
			}
			return c.NoContent(http.StatusNoContent)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound)
}

func ReadDNSRecord(c echo.Context) error {
	networkID, domain, nerr := getNetworkIDAndDomain(c)
	if nerr != nil {
		return nerr
	}

	dnsConfig, nerr := getExistingDNSConfig(c.Request().Context(), networkID)
	if nerr != nil {
		return nerr
	}
	for _, record := range dnsConfig.Records {
		if record.Domain == domain {
			return c.JSON(http.StatusOK, record)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound)
}

func DeleteDNSRecord(c echo.Context) error {
	networkID, domain, nerr := getNetworkIDAndDomain(c)
	if nerr != nil {
		return nerr
	}

	reqCtx := c.Request().Context()
	dnsConfig, nerr := getExistingDNSConfig(reqCtx, networkID)
	if nerr != nil {
		return nerr
	}

	for i, record := range dnsConfig.Records {
		if record.Domain == domain {
			if i == len(dnsConfig.Records)-1 {
				dnsConfig.Records = dnsConfig.Records[:i]
			} else {
				dnsConfig.Records = append(dnsConfig.Records[:i], dnsConfig.Records[i+1:]...)
			}
			nerr = updateDNSConfig(reqCtx, networkID, dnsConfig)
			if nerr != nil {
				return nerr
			}
			return c.NoContent(http.StatusNoContent)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound)
}

func updateDNSConfig(ctx context.Context, networkID string, dnsConfig *models.NetworkDNSConfig) *echo.HTTPError {
	err := configurator.UpdateNetworks(
		ctx,
		[]configurator.NetworkUpdateCriteria{
			{
				ID:                   networkID,
				ConfigsToAddOrUpdate: map[string]interface{}{orc8r.DnsdNetworkType: dnsConfig},
			},
		},
		serdes.Network,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func getNetworkIDAndDomain(c echo.Context) (string, string, *echo.HTTPError) {
	vals, nerr := obsidian.GetParamValues(c, "network_id", "domain")
	if nerr != nil {
		return "", "", nerr
	}
	return vals[0], vals[1], nil
}

func getExistingDNSConfig(ctx context.Context, networkID string) (*models.NetworkDNSConfig, *echo.HTTPError) {
	iDNSConfig, err := configurator.LoadNetworkConfig(ctx, networkID, orc8r.DnsdNetworkType, serdes.Network)
	if err == merrors.ErrNotFound {
		return nil, echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return iDNSConfig.(*models.NetworkDNSConfig), nil
}

func getRecordAndValidate(c echo.Context, domain string) (*models.DNSConfigRecord, *echo.HTTPError) {
	payload, nerr := GetAndValidatePayload(c, &models.DNSConfigRecord{})
	if nerr != nil {
		return nil, nerr
	}
	record := payload.(*models.DNSConfigRecord)

	if record.Domain != domain {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Domain name in param and record don't match"))
	}
	return record, nil
}
