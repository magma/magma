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
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"

	"github.com/facebookincubator/prometheus-configmanager/alertmanager/config"
	"github.com/labstack/echo/v4"

	"magma/orc8r/cloud/go/services/obsidian"
)

const (
	ReceiverNamePathParam  = "receiver"
	ReceiverNameQueryParam = "receiver"
)

func GetConfigureAlertReceiverHandler(configManagerURL string, client HttpClient) func(echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, configureAlertReceiver, client)
}

func GetRetrieveAlertReceiverHandler(configManagerURL string, client HttpClient) func(echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, retrieveAlertReceivers, client)
}

func GetUpdateAlertReceiverHandler(configManagerURL string, client HttpClient) func(echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, updateAlertReceiver, client)
}

func GetDeleteAlertReceiverHandler(configManagerURL string, client HttpClient) func(echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, deleteAlertReceiver, client)
}

// getHandlerWithReceiverFunc returns an echo HandlerFunc that checks the
// networkID and runs the given handlerImplFunc that communicates with the
// alertmanager config service
func getHandlerWithReceiverFunc(configManagerURL string, handlerImplFunc func(echo.Context, string, HttpClient) error, client HttpClient) func(echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := makeNetworkReceiverPath(configManagerURL, networkID)
		return handlerImplFunc(c, url, client)
	}
}

func GetRetrieveAlertRouteHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return getHandlerWithRouteFunc(configManagerURL, retrieveAlertRoute, client)
}

func GetUpdateAlertRouteHandler(configManagerURL string, client HttpClient) func(c echo.Context) error {
	return getHandlerWithRouteFunc(configManagerURL, updateAlertRoute, client)
}

// getHandlerWithRouteFunc returns an echo HandlerFunc that checks the
// networkID and runs the given handlerImplFunc that communicates with the
// alertmanager config service for routing trees
func getHandlerWithRouteFunc(configManagerURL string, handlerImplFunc func(echo.Context, string, HttpClient) error, client HttpClient) func(echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := makeNetworkRoutePath(configManagerURL, networkID)
		return handlerImplFunc(c, url, client)
	}
}

func configureAlertReceiver(c echo.Context, url string, client HttpClient) error {
	receiver, err := buildReceiverFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sendErr := sendConfig(receiver, url, http.MethodPost, client)
	if sendErr != nil {
		return sendErr
	}
	return c.NoContent(http.StatusOK)
}

func retrieveAlertReceivers(c echo.Context, url string, client HttpClient) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error reading receivers: %v", body.Message))
	}
	var recs []config.Receiver
	err = json.NewDecoder(resp.Body).Decode(&recs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error decoding server response %v", err))
	}
	return c.JSON(http.StatusOK, recs)
}

func updateAlertReceiver(c echo.Context, url string, client HttpClient) error {
	receiver, err := buildReceiverFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	receiverName := c.Param(ReceiverNamePathParam)
	if receiverName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "receiver name not provided")
	}
	if receiverName != receiver.Name {
		return echo.NewHTTPError(http.StatusBadRequest, "new receiver configuration must have same name")
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(receiverName))

	sendErr := sendConfig(receiver, url, http.MethodPut, client)
	if sendErr != nil {
		return sendErr
	}
	return c.NoContent(http.StatusOK)
}

func deleteAlertReceiver(c echo.Context, url string, client HttpClient) error {
	receiverName := c.QueryParam(ReceiverNameQueryParam)
	if receiverName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "receiver name not provided")
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(receiverName))

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error deleting receiver: %v", body.Message))
	}
	return c.NoContent(http.StatusOK)
}

func retrieveAlertRoute(c echo.Context, url string, client HttpClient) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return echo.NewHTTPError(resp.StatusCode, fmt.Sprintf("error reading alerting route: %v", body.Message))
	}
	var route config.Route
	err = json.NewDecoder(resp.Body).Decode(&route)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error decoding server response %v", err))
	}
	return c.JSON(http.StatusOK, route)
}

func updateAlertRoute(c echo.Context, url string, client HttpClient) error {
	route, err := buildRouteFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid route specification: %v", err))
	}

	sendErr := sendConfig(route, url, http.MethodPost, client)
	if sendErr != nil {
		return echo.NewHTTPError(sendErr.Code, fmt.Sprintf("error updating alert route: %v", sendErr.Message))
	}
	return c.NoContent(http.StatusOK)
}

func buildReceiverFromContext(c echo.Context) (config.Receiver, error) {
	wrapper := config.Receiver{}
	err := json.NewDecoder(c.Request().Body).Decode(&wrapper)
	if err != nil {
		return config.Receiver{}, err
	}
	return wrapper, nil
}

func buildRouteFromContext(c echo.Context) (config.Route, error) {
	route := config.Route{}
	err := json.NewDecoder(c.Request().Body).Decode(&route)
	if err != nil {
		return config.Route{}, err
	}
	return route, nil
}

func makeNetworkReceiverPath(configManagerURL, networkID string) string {
	return configManagerURL + "/" + networkID + "/receiver"
}

func makeNetworkRoutePath(configManagerURL, networkID string) string {
	return configManagerURL + "/" + networkID + "/route"
}
