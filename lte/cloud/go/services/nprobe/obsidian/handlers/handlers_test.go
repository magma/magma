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

package handlers_test

import (
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/nprobe/obsidian/handlers"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/lte/cloud/go/services/nprobe/storage"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func getNProbeBlobstore(t *testing.T) storage.NProbeStorage {
	fact := test_utils.NewSQLBlobstore(t, "nprobe_handlers_test_blobstore")
	return storage.NewNProbeBlobstore(fact)
}

func TestCreateNetworkProbeTask(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/tasks"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	createNetworkProbeTask := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	payload := &models.NetworkProbeTask{
		TaskID: "test",
		TaskDetails: &models.NetworkProbeTaskDetails{
			TargetID:      "test",
			TargetType:    "imsi",
			DeliveryType:  "all",
			CorrelationID: 8674665223082154000,
			Timestamp:     strfmt.DateTime(time.Now().UTC()),
		},
	}

	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        payload,
		Handler:        createNetworkProbeTask,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.NetworkProbeTaskEntityType, "test", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.NetworkProbeTaskEntityType,
		Key:       "test",
		Config:    payload.TaskDetails,
		GraphID:   "2",
	}

	expected_task := expected.Config.(*models.NetworkProbeTaskDetails)
	actual_task := actual.Config.(*models.NetworkProbeTaskDetails)

	assert.Equal(t, expected_task.TargetID, actual_task.TargetID)
	assert.Equal(t, expected_task.TargetType, actual_task.TargetType)
	assert.Equal(t, expected_task.DeliveryType, actual_task.DeliveryType)
	assert.Equal(t, expected_task.CorrelationID, actual_task.CorrelationID)
}

func TestListNetworkProbeTasks(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/tasks"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	listNetworkProbeTasks := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listNetworkProbeTasks,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.NetworkProbeTask{}),
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Key:  "IMSI1234",
				Type: lte.NetworkProbeTaskEntityType,
				Config: &models.NetworkProbeTaskDetails{
					TargetID:      "IMSI1234",
					TargetType:    "imsi",
					DeliveryType:  "events_only",
					CorrelationID: 8674665223082154000,
				},
			},
			{
				Key:  "IMSI1235",
				Type: lte.NetworkProbeTaskEntityType,
				Config: &models.NetworkProbeTaskDetails{
					TargetID:      "IMSI1235",
					TargetType:    "imsi",
					DeliveryType:  "all",
					CorrelationID: 8674665223082154099,
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listNetworkProbeTasks,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.NetworkProbeTask{
			"IMSI1234": {
				TaskID: "IMSI1234",
				TaskDetails: &models.NetworkProbeTaskDetails{
					TargetID:      "IMSI1234",
					TargetType:    "imsi",
					DeliveryType:  "events_only",
					CorrelationID: 8674665223082154000,
				},
			},
			"IMSI1235": {
				TaskID: "IMSI1235",
				TaskDetails: &models.NetworkProbeTaskDetails{
					TargetID:      "IMSI1235",
					TargetType:    "imsi",
					DeliveryType:  "all",
					CorrelationID: 8674665223082154099,
				},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetNetworkProbeTask(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/tasks/:task_id"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	getNetworkProbeTask := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getNetworkProbeTask,
		ParamNames:     []string{"network_id", "task_id"},
		ParamValues:    []string{"n1", "IMSI1234"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Key:  "IMSI1234",
			Type: lte.NetworkProbeTaskEntityType,
			Config: &models.NetworkProbeTaskDetails{
				TargetID:      "IMSI1234",
				TargetType:    "imsi",
				DeliveryType:  "events_only",
				CorrelationID: 8674665223082154000,
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getNetworkProbeTask,
		ParamNames:     []string{"network_id", "task_id"},
		ParamValues:    []string{"n1", "IMSI1234"},
		ExpectedStatus: 200,
		ExpectedResult: &models.NetworkProbeTask{
			TaskID: "IMSI1234",
			TaskDetails: &models.NetworkProbeTaskDetails{
				TargetID:      "IMSI1234",
				TargetType:    "imsi",
				DeliveryType:  "events_only",
				CorrelationID: 8674665223082154000,
			},
		},
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetworkProbeTask(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/tasks/:task_id"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	updateNetworkProbeTask := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	// 404
	payload := &models.NetworkProbeTask{
		TaskID: "IMSI1234",
		TaskDetails: &models.NetworkProbeTaskDetails{
			TargetID:      "IMSI1234",
			TargetType:    "imsi",
			DeliveryType:  "events_only",
			CorrelationID: 8674665223082154000,
		},
	}

	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateNetworkProbeTask,
		Payload:        payload,
		ParamNames:     []string{"network_id", "task_id"},
		ParamValues:    []string{"n1", "IMSI1234"},
		ExpectedStatus: 500,
		ExpectedError:  "failed to load entity being updated: expected to load 1 ent for update, got 0",
	}
	tests.RunUnitTest(t, e, tc)

	// Add the NetworkProbeTask
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Key:  "IMSI1234",
			Type: lte.NetworkProbeTaskEntityType,
			Config: &models.NetworkProbeTaskDetails{
				TargetID:      "IMSI1234",
				TargetType:    "imsi",
				DeliveryType:  "events_only",
				CorrelationID: 8674665223082154000,
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateNetworkProbeTask,
		Payload:        payload,
		ParamNames:     []string{"network_id", "task_id"},
		ParamValues:    []string{"n1", "IMSI1234"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.NetworkProbeTaskEntityType, "IMSI1234", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.NetworkProbeTaskEntityType,
		Key:       "IMSI1234",
		Config:    payload.TaskDetails,
		GraphID:   "2",
		Version:   1,
	}
	assert.Equal(t, expected, actual)
}

func TestDeleteNetworkProbeTask(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/tasks/:task_id"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	deleteNetworkProbeTask := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Key:  "IMSI1234",
				Type: lte.NetworkProbeTaskEntityType,
				Config: &models.NetworkProbeTaskDetails{
					TargetID:      "IMSI1234",
					TargetType:    "imsi",
					DeliveryType:  "events_only",
					CorrelationID: 8674665223082154000,
				},
			},
			{
				Key:  "IMSI1235",
				Type: lte.NetworkProbeTaskEntityType,
				Config: &models.NetworkProbeTaskDetails{
					TargetID:      "IMSI1235",
					TargetType:    "imsi",
					DeliveryType:  "all",
					CorrelationID: 8674665223082154099,
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteNetworkProbeTask,
		ParamNames:     []string{"network_id", "task_id"},
		ParamValues:    []string{"n1", "IMSI1234"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, _, err := configurator.LoadAllEntitiesOfType("n1", lte.NetworkProbeTaskEntityType, configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actual))
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.NetworkProbeTaskEntityType,
		Key:       "IMSI1235",
		Config: &models.NetworkProbeTaskDetails{
			TargetID:      "IMSI1235",
			TargetType:    "imsi",
			DeliveryType:  "all",
			CorrelationID: 8674665223082154099,
		},
		GraphID: "4",
		Version: 0,
	}
	assert.Equal(t, expected, actual[0])
}

func TestCreateNetworkProbeDestination(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/destinations"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	createNetworkProbeDestination := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	payload := &models.NetworkProbeDestination{
		DestinationID: "test",
		DestinationDetails: &models.NetworkProbeDestinationDetails{
			DeliveryAddress: "127.0.0.1:4000",
			DeliveryType:    "all",
		},
	}

	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        payload,
		Handler:        createNetworkProbeDestination,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.NetworkProbeDestinationEntityType, "test", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.NetworkProbeDestinationEntityType,
		Key:       "test",
		Config:    payload.DestinationDetails,
		GraphID:   "2",
	}
	assert.Equal(t, expected, actual)
}

func TestListNetworkProbeDestinations(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/destinations"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	listNetworkProbeDestinations := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listNetworkProbeDestinations,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.NetworkProbeDestination{}),
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Key:  "1111-2222-3333",
				Type: lte.NetworkProbeDestinationEntityType,
				Config: &models.NetworkProbeDestinationDetails{
					DeliveryAddress: "127.0.0.1:4000",
					DeliveryType:    "all",
				},
			},
			{
				Key:  "2222-3333-4444",
				Type: lte.NetworkProbeDestinationEntityType,
				Config: &models.NetworkProbeDestinationDetails{
					DeliveryAddress: "127.0.0.1:4001",
					DeliveryType:    "all",
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listNetworkProbeDestinations,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.NetworkProbeDestination{
			"1111-2222-3333": {
				DestinationID: "1111-2222-3333",
				DestinationDetails: &models.NetworkProbeDestinationDetails{
					DeliveryAddress: "127.0.0.1:4000",
					DeliveryType:    "all",
				},
			},
			"2222-3333-4444": {
				DestinationID: "2222-3333-4444",
				DestinationDetails: &models.NetworkProbeDestinationDetails{
					DeliveryAddress: "127.0.0.1:4001",
					DeliveryType:    "all",
				},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetNetworkProbeDestination(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/destinations/:destination_id"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	getNetworkProbeDestination := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getNetworkProbeDestination,
		ParamNames:     []string{"network_id", "destination_id"},
		ParamValues:    []string{"n1", "1111-2222-3333"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Key:  "1111-2222-3333",
			Type: lte.NetworkProbeDestinationEntityType,
			Config: &models.NetworkProbeDestinationDetails{
				DeliveryAddress: "127.0.0.1:4000",
				DeliveryType:    "all",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getNetworkProbeDestination,
		ParamNames:     []string{"network_id", "destination_id"},
		ParamValues:    []string{"n1", "1111-2222-3333"},
		ExpectedStatus: 200,
		ExpectedResult: &models.NetworkProbeDestination{
			DestinationID: "1111-2222-3333",
			DestinationDetails: &models.NetworkProbeDestinationDetails{
				DeliveryAddress: "127.0.0.1:4000",
				DeliveryType:    "all",
			},
		},
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateNetworkProbeDestination(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/destinations/:destination_id"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	updateNetworkProbeDestination := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	// 404
	payload := &models.NetworkProbeDestination{
		DestinationID: "1111-2222-3333",
		DestinationDetails: &models.NetworkProbeDestinationDetails{
			DeliveryAddress: "127.0.0.1:4000",
			DeliveryType:    "all",
		},
	}

	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateNetworkProbeDestination,
		Payload:        payload,
		ParamNames:     []string{"network_id", "destination_id"},
		ParamValues:    []string{"n1", "1111-2222-3333"},
		ExpectedStatus: 500,
		ExpectedError:  "failed to load entity being updated: expected to load 1 ent for update, got 0",
	}
	tests.RunUnitTest(t, e, tc)

	// Add the NetworkProbeDestination
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Key:  "1111-2222-3333",
			Type: lte.NetworkProbeDestinationEntityType,
			Config: &models.NetworkProbeDestinationDetails{
				DeliveryAddress: "127.0.0.1:4000",
				DeliveryType:    "all",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateNetworkProbeDestination,
		Payload:        payload,
		ParamNames:     []string{"network_id", "destination_id"},
		ParamValues:    []string{"n1", "1111-2222-3333"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.NetworkProbeDestinationEntityType, "1111-2222-3333", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.NetworkProbeDestinationEntityType,
		Key:       "1111-2222-3333",
		Config:    payload.DestinationDetails,
		GraphID:   "2",
		Version:   1,
	}
	assert.Equal(t, expected, actual)
}

func TestDeleteNetworkProbeDestination(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/network_probe/destinations/:destination_id"
	handlers := handlers.GetHandlers(getNProbeBlobstore(t))
	deleteNetworkProbeDestination := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Key:  "1111-2222-3333",
				Type: lte.NetworkProbeDestinationEntityType,
				Config: &models.NetworkProbeDestinationDetails{
					DeliveryAddress: "127.0.0.1:4000",
					DeliveryType:    "events_only",
				},
			},
			{
				Key:  "2222-3333-4444",
				Type: lte.NetworkProbeDestinationEntityType,
				Config: &models.NetworkProbeDestinationDetails{
					DeliveryAddress: "127.0.0.1:4001",
					DeliveryType:    "all",
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteNetworkProbeDestination,
		ParamNames:     []string{"network_id", "destination_id"},
		ParamValues:    []string{"n1", "1111-2222-3333"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, _, err := configurator.LoadAllEntitiesOfType("n1", lte.NetworkProbeDestinationEntityType, configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actual))
	expected := configurator.NetworkEntity{
		NetworkID: "n1",
		Type:      lte.NetworkProbeDestinationEntityType,
		Key:       "2222-3333-4444",
		Config: &models.NetworkProbeDestinationDetails{
			DeliveryAddress: "127.0.0.1:4001",
			DeliveryType:    "all",
		},
		GraphID: "4",
		Version: 0,
	}
	assert.Equal(t, expected, actual[0])
}
