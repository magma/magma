/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// NOTE: to run these tests outside the testing environment, e.g. from IntelliJ,
// ensure postgres_test and maria_test containers are running, and use the
// following environment variables to point to the relevant DB endpoints:
//	- TEST_DATABASE_HOST=localhost
//	- TEST_DATABASE_PORT_POSTGRES=5433
//	- TEST_DATABASE_PORT_MARIA=3307

// integ_test.go tests indexing and reindexing using remote indexers.

package indexer_test

import (
	"context"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_test "magma/orc8r/cloud/go/services/state/test_utils"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
)

const (
	nid0  = "some_networkid_0"
	hwid0 = "some_hwid_0"

	indexTimeout = 5 * time.Second
)

var (
	sid0     = state_types.ID{Type: orc8r.GatewayStateType, DeviceID: "some_imsi"}
	hwidByID = map[state_types.ID]string{
		sid0: hwid0,
	}
	statusByID = map[state_types.ID]*models.GatewayStatus{
		sid0: {Meta: map[string]string{"foo": "bar"}},
	}
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestStateIndexing(t *testing.T) {
	const (
		serviceName                 = "SOME_SERVICE_NAME"
		zero        indexer.Version = 0
		version0    indexer.Version = 100
	)
	var (
		types     = []string{orc8r.GatewayStateType}
		prepare0  = make(chan mock.Arguments)
		complete0 = make(chan mock.Arguments)
		index0    = make(chan mock.Arguments)
	)

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	idx0 := mocks.NewMockIndexer("some_id", version0, types, prepare0, complete0, index0)
	remoteIdx0 := indexer.NewRemoteIndexer(serviceName, version0, types...)
	state_test_init.StartNewTestIndexer(t, serviceName, idx0)

	dbName := "state___integ_test"
	r, q := initTestServices(t, dbName)

	// Register remote, to be called by this test code (which forwards over the network to the locally-registered idx)
	indexer.DeregisterAllForTest(t)
	err := indexer.RegisterIndexers(remoteIdx0)
	assert.NoError(t, err)

	t.Run("index", func(t *testing.T) {
		reportGatewayStatusForID(t, sid0)

		// Index args: (networkID string, states state_types.StatesByID)
		recv := recvArgs(t, index0, "index0")
		assertEqualStr(t, nid0, recv[0])
		assertEqualStatus(t, recv[1], sid0)
	})

	_, err = q.PopulateJobs()
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	go r.Run(ctx)
	defer cancel()

	t.Run("reindex", func(t *testing.T) {
		// Prepare args: (from, to Version, isFirstReindex bool)
		recvPrepare0 := recvArgs(t, prepare0, "prepare0")
		assertEqualVersion(t, zero, recvPrepare0[0])
		assertEqualVersion(t, version0, recvPrepare0[1])
		assertEqualBool(t, true, recvPrepare0[2])

		// Index args: (networkID string, states state_types.StatesByID)
		recvIndex0 := recvArgs(t, index0, "index0")
		assertEqualStr(t, nid0, recvIndex0[0])
		assertEqualStatus(t, recvIndex0[1], sid0)

		// Complete args: (from, to Version)
		recvComplete0 := recvArgs(t, complete0, "complete0")
		assertEqualVersion(t, zero, recvComplete0[0])
		assertEqualVersion(t, version0, recvComplete0[1])
	})
}

func initTestServices(t *testing.T, dbName string) (reindex.Reindexer, reindex.JobQueue) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	indexer.DeregisterAllForTest(t)

	device_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)
	configurator_test.RegisterNetwork(t, nid0, "Network 0 for indexer integ test")
	configurator_test.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})

	return state_test_init.StartTestServiceInternal(t, dbName, sqorc.PostgresDriver)
}

func reportGatewayStatusForID(t *testing.T, id state_types.ID) {
	ctx := state_test.GetContextWithCertificate(t, hwidByID[id])
	status := statusByID[id]

	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serialized, err := serde.Serialize(state.SerdeDomain, orc8r.GatewayStateType, status)
	assert.NoError(t, err)
	pState := &protos.State{
		Type:     orc8r.GatewayStateType,
		DeviceID: id.DeviceID,
		Value:    serialized,
	}

	_, err = client.ReportStates(ctx, &protos.ReportStatesRequest{States: []*protos.State{pState}})
	assert.NoError(t, err)
}

func recvArgs(t *testing.T, ch chan mock.Arguments, chName string) mock.Arguments {
	select {
	case args := <-ch:
		return args
	case <-time.After(indexTimeout):
		t.Fatalf("Timeout waiting for args on channel %v", chName)
		return nil
	}
}

func assertEqualStr(t *testing.T, expected string, recv interface{}) {
	recvVal, ok := recv.(string)
	assert.True(t, ok)
	assert.Equal(t, expected, recvVal)
}

func assertEqualVersion(t *testing.T, expected indexer.Version, recv interface{}) {
	recvVal, ok := recv.(indexer.Version)
	assert.True(t, ok)
	assert.Equal(t, expected, recvVal)
}

func assertEqualBool(t *testing.T, expected bool, recv interface{}) {
	recvVal, ok := recv.(bool)
	assert.True(t, ok)
	assert.Equal(t, expected, recvVal)
}

func assertEqualStatus(t *testing.T, recv interface{}, sid state_types.ID) {
	hwid := hwidByID[sid]
	reported := statusByID[sid]
	recvStates := recv.(state_types.StatesByID)
	assert.Len(t, recvStates, 1)
	assert.Equal(t, orc8r.GatewayStateType, recvStates[sid].Type)
	assert.Equal(t, hwid, recvStates[sid].ReporterID)
	assert.Equal(t, reported, recvStates[sid].ReportedState)
}
