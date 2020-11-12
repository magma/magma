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
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
)

// PartialEntityModel describe models that represents a portion of network
// entity that can be read and updated.
type PartialEntityModel interface {
	serde.ValidatableModel
	// FromBackendModels the same PartialEntityModel from the configurator
	// entities attached to the networkID and key.
	FromBackendModels(networkID string, key string) error
	// ToUpdateCriteria returns a EntityUpdateCriteria needed to apply
	// the change in the model.
	ToUpdateCriteria(networkID string, key string) ([]configurator.EntityUpdateCriteria, error)
}

// GetPartialEntityHandlers returns both GET and PUT handlers for modifying the portion of a
// network entity specified by the model.
// - path : 	the url at which the handler will be registered.
// - paramName: the parameter name in the url at which the entity key is stored
// - model: 	the input and output of the handler and it also provides FromBackendModels
//   and ToUpdateCriteria to go between the configurator model.
func GetPartialEntityHandlers(path string, paramName string, model PartialEntityModel, serdes serde.Registry) []obsidian.Handler {
	return []obsidian.Handler{
		GetPartialUpdateEntityHandler(path, paramName, model, serdes),
		GetPartialReadEntityHandler(path, paramName, model, serdes),
	}
}

// GetPartialReadEntityHandler returns a GET obsidian handler at the specified path.
// This function loads a portion of the gateway specified by the model's FromBackendModels function.
// Example:
// 		(m *TierName) FromBackendModels(networkID, tierID string) error {
// 			entity, err := configurator.LoadEntity(networkID, orc8r.UpgradeTierEntityType, key, configurator.EntityLoadCriteria{LoadMetadata: true})
//			if err != nil {
//				return err
//			}
//			*m = TierName(entity.Name)
//			return nil
//		}
// 		getTierNameHandler := handlers.GetPartialReadEntityHandler(URL, "tier_id", new(models.TierName))
//      would return a GET handler that can read the tier name of a tier with the specified ID.
func GetPartialReadEntityHandler(path string, paramName string, model PartialEntityModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, key, nerr := getNetworkAndEntityIDs(c, paramName)
			if nerr != nil {
				return nerr
			}

			err := model.FromBackendModels(networkID, key)
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.JSON(http.StatusOK, model)
		},
	}
}

// GetPartialUpdateEntityHandler returns a PUT obsidian handler at the specified path.
// This function updates a portion of the network entity specified by the model's ToUpdateCriteria function.
// Example:
//      (m *TierName) ToUpdateCriteria(networkID, tierID string) (configurator.EntityUpdateCriteria, error) {
// 			return configurator.EntityUpdateCriteria{
//				{
// 					Key: gatewayID,
//					Type: orc8r.MagmadGatewayType,
//					NewName: m,
//				}
//          }
// 		}
// 		updateTierNameHandler := handlers.GetPartialUpdateEntityHandler(URL, "tier_id", new(models.TierName))
//      would return a PUT handler that updates the tier name of a tier with the specified ID.
func GetPartialUpdateEntityHandler(path string, paramName string, model PartialEntityModel, serdes serde.Registry) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkID, key, nerr := getNetworkAndEntityIDs(c, paramName)
			if nerr != nil {
				return nerr
			}

			requestedUpdate, nerr := GetAndValidatePayload(c, model)
			if nerr != nil {
				return nerr
			}

			updates, err := requestedUpdate.(PartialEntityModel).ToUpdateCriteria(networkID, key)
			if err != nil {
				return obsidian.HttpError(err, http.StatusBadRequest)
			}
			_, err = configurator.UpdateEntities(networkID, updates, serdes)
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

func getNetworkAndEntityIDs(c echo.Context, paramName string) (string, string, *echo.HTTPError) {
	vals, nerr := obsidian.GetParamValues(c, "network_id", paramName)
	if nerr != nil {
		return "", "", nerr
	}
	return vals[0], vals[1], nil
}
