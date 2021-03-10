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
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// NetworkModel describes models that represent a certain type of network.
// For example, an LTE network, that can be read/updated/deleted
type NetworkModel interface {
	serde.ValidatableModel
	// GetEmptyNetwork creates a new instance of the typed NetworkModel.
	// It should be empty
	GetEmptyNetwork() NetworkModel
	// ToConfiguratorNetwork should convert the Network model to
	// a configurator.network
	ToConfiguratorNetwork() configurator.Network
	// ToUpdateCriteria takes in the existing network and applies the change
	// from the model to create a NetworkUpdateCriteria
	ToUpdateCriteria() configurator.NetworkUpdateCriteria
	// FromConfiguratorNetwork should return a copy of the network
	FromConfiguratorNetwork(n configurator.Network) interface{}
}

// PartialNetworkModel describe models that represents a portion of network
// that can be read, updated, and deleted.
type PartialNetworkModel interface {
	serde.ValidatableModel
	// GetFromNetwork grabs the desired model from the configurator network.
	// Returns nil if it is not there.
	GetFromNetwork(network configurator.Network) interface{}
	// ToUpdateCriteria takes in the existing network and applies the change
	// from the model to create a NetworkUpdateCriteria
	ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error)
}

// GetPartialNetworkHandlers returns a set of GET/PUT/DELETE handlers according to the parameters.
// If the configKey is not "", it will add a delete handler for the network config for that key.
func GetPartialNetworkHandlers(path string, model PartialNetworkModel, configKey string, serdes serde.Registry) []obsidian.Handler {
	ret := []obsidian.Handler{
		GetPartialReadNetworkHandler(path, model, serdes),
		GetPartialUpdateNetworkHandler(path, model, serdes),
	}
	if configKey != "" {
		ret = append(ret, GetPartialDeleteNetworkHandler(path, configKey, serdes))
	}
	return ret
}

// GetPartialReadNetworkHandler returns a GET obsidian handler at the specified path.
// This function loads a network specified by the networkID and returns the
// part of the network that corresponds to the given model.
// Example:
//      (m *NetworkName) GetFromNetwork(network configurator.Network) interface{} {
// 			return string(network.Name)
// 		}
// 		getNameHandler := handlers.GetPartialReadNetworkHandler(URL, &models.NetworkName{})
//
//      would return a GET handler that can read the network name of a network with the specified ID.
func GetPartialReadNetworkHandler(path string, model PartialNetworkModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			network, err := configurator.LoadNetwork(networkID, true, true, serdes)
			if err == merrors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			ret := model.GetFromNetwork(network)
			if ret == nil {
				return obsidian.HttpError(fmt.Errorf("Not found"), http.StatusNotFound)
			}
			return c.JSON(http.StatusOK, ret)
		},
	}
}

// GetPartialUpdateNetworkHandler returns a PUT obsidian handler at the specified path.
// The handler will fetch the payload into the configModel and perform validations according to the swagger spec.
// updater will take the model and apply the change into an existing network.
// Example:
//      (m *NetworkName) ToUpdateCriteria(network configurator.Network) interface{} {
// 			return configurator.NetworkUpdateCriteria{
//				ID:   network.ID,
// 				Name: *m,
//			}
//      }
// 		putNameHandler := handlers.GetPartialUpdateNetworkHandler(URL, &models.NetworkName{})
//
//      would return a PUT handler that will intake a NetworkName model and update the corresponding network
func GetPartialUpdateNetworkHandler(path string, model PartialNetworkModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			requestedUpdate, nerr := GetAndValidatePayload(c, model)
			if nerr != nil {
				return nerr
			}

			network, err := configurator.LoadNetwork(networkID, true, true, serdes)
			if err == merrors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}

			updateCriteria, err := requestedUpdate.(PartialNetworkModel).ToUpdateCriteria(network)
			if err != nil {
				return obsidian.HttpError(err, http.StatusBadRequest)
			}
			err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria}, serdes)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

// GetPartialDeleteNetworkHandler returns a DELETE obsidian handler at the specified path.
// The handler will delete a network config specified by the key.
// Example:
// 		deleteNetworkFeaturesHandler := handlers.GetPartialDeleteNetworkHandler(URL, "orc8r_features")
//
//      would return a DELETE handler that will remove the network features config from the corresponding network
func GetPartialDeleteNetworkHandler(path string, key string, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.DELETE,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			update := configurator.NetworkUpdateCriteria{
				ID:              networkID,
				ConfigsToDelete: []string{key},
			}
			err := configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update}, serdes)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

func GetTypedNetworkCRUDHandlers(listCreatePath string, getUpdateDeletePath string, networkType string, network NetworkModel, serdes serde.Registry) []obsidian.Handler {
	return []obsidian.Handler{
		getListTypedNetworksHandler(listCreatePath, networkType),
		getCreateTypedNetworkHandler(listCreatePath, networkType, network, serdes),
		getGetTypedNetworkHandler(getUpdateDeletePath, networkType, network, serdes),
		getUpdateTypedNetworkHandler(getUpdateDeletePath, networkType, network, serdes),
		getDeleteTypedNetworkHandler(getUpdateDeletePath, networkType, serdes),
	}
}

func getListTypedNetworksHandler(path string, networkType string) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			ids, err := configurator.ListNetworksOfType(networkType)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			sort.Strings(ids)
			return c.JSON(http.StatusOK, ids)
		},
	}
}

func getCreateTypedNetworkHandler(path string, networkType string, network NetworkModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.POST,
		HandlerFunc: func(c echo.Context) error {
			payload, err := getAndValidateNetwork(c, network)
			if err != nil {
				return err
			}
			err = configurator.CreateNetwork(payload.ToConfiguratorNetwork(), serdes)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusCreated)

		},
	}
}

func getGetTypedNetworkHandler(path string, networkType string, networkModel NetworkModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			nid, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}

			network, err := configurator.LoadNetwork(nid, true, true, serdes)
			if err == merrors.ErrNotFound {
				return c.NoContent(http.StatusNotFound)
			}
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			if network.Type != networkType {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a <%s> network", nid, networkType))
			}

			ret := (networkModel.GetEmptyNetwork()).FromConfiguratorNetwork(network)
			return c.JSON(http.StatusOK, ret)
		},
	}
}

func getUpdateTypedNetworkHandler(path string, networkType string, networkModel NetworkModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			nid, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}

			payload, err := getAndValidateNetwork(c, networkModel)
			if err != nil {
				return err
			}

			network, err := configurator.LoadNetwork(nid, false, false, serdes)
			if err == merrors.ErrNotFound {
				return c.NoContent(http.StatusNotFound)
			}
			if err != nil {
				return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
			}
			if network.Type != networkType {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a <%s> network", nid, networkType))
			}

			err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{payload.ToUpdateCriteria()}, serdes)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

func getDeleteTypedNetworkHandler(path string, networkType string, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.DELETE,
		HandlerFunc: func(c echo.Context) error {
			nid, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}

			network, err := configurator.LoadNetwork(nid, false, false, serdes)
			if err == merrors.ErrNotFound {
				return c.NoContent(http.StatusNotFound)
			}
			if err != nil {
				return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
			}
			if network.Type != networkType {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a <%s> network", nid, networkType))
			}

			err = configurator.DeleteNetwork(nid)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

// getAndValidateNetwork can be used by any model that implements NetworkModel
func getAndValidateNetwork(c echo.Context, network interface{}) (NetworkModel, error) {
	iModel := reflect.New(reflect.TypeOf(network).Elem()).Interface().(NetworkModel)
	if err := c.Bind(iModel); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	// Run validations specified by the swagger spec
	if err := iModel.ValidateModel(); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	return iModel, nil
}
