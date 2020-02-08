/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin/models"
	lteprotos "magma/lte/cloud/go/protos"
	utils "magma/lte/cloud/go/services/eps_authentication/servicers/test_utils"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/test_utils"
	orc8rprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

// EpsAuthTestSuite is a test suite that will setup a test EPS auth server
// and config service pre-populated with a cellular config.
type EpsAuthTestSuite struct {
	suite.Suite
	Server *EPSAuthServer
}

func (suite *EpsAuthTestSuite) AuthenticationInformation(air *lteprotos.AuthenticationInformationRequest) (*lteprotos.AuthenticationInformationAnswer, error) {
	return suite.Server.AuthenticationInformation(getTestContext(), air)
}

func (suite *EpsAuthTestSuite) UpdateLocation(ulr *lteprotos.UpdateLocationRequest) (*lteprotos.UpdateLocationAnswer, error) {
	return suite.Server.UpdateLocation(getTestContext(), ulr)
}

func (suite *EpsAuthTestSuite) PurgeUE(purge *lteprotos.PurgeUERequest) (*lteprotos.PurgeUEAnswer, error) {
	return suite.Server.PurgeUE(getTestContext(), purge)
}

func (suite *EpsAuthTestSuite) SetupTest() {
	store, err := storage.NewSubscriberDBStorage(test_utils.NewMockDatastore())
	suite.NoError(err)

	for _, subscriber := range utils.GetTestSubscribers() {
		_, err := store.AddSubscriber(subscriber)
		suite.NoError(err)
	}

	server, err := NewEPSAuthServer(store)
	suite.NoError(err)
	suite.Server = server
}

func TestEpsAuthSuite(t *testing.T) {
	test_init.StartTestService(t)
	err := serde.RegisterSerdes(configurator.NewNetworkConfigSerde(lte.CellularNetworkType, &models.NetworkCellularConfigs{}))
	assert.NoError(t, err)

	cellularConfig := &models.NetworkCellularConfigs{
		Ran: &models.NetworkRanConfigs{},
		Epc: &models.NetworkEpcConfigs{
			Mcc:        "123",
			Mnc:        "123",
			Tac:        1,
			LteAuthOp:  []byte("\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18"),
			LteAuthAmf: []byte("\x80\x00"),
			SubProfiles: map[string]models.NetworkEpcConfigsSubProfilesAnon{
				"default": {
					MaxUlBitRate: 1000,
					MaxDlBitRate: 2000,
				},
				"test_profile": {
					MaxUlBitRate: 7000,
					MaxDlBitRate: 5000,
				},
			},
		},
	}
	err = configurator.CreateNetwork(configurator.Network{
		ID:      "test",
		Type:    "lte",
		Configs: map[string]interface{}{lte.CellularNetworkType: cellularConfig},
	})
	assert.NoError(t, err)

	testSuite := &EpsAuthTestSuite{}
	suite.Run(t, testSuite)
}

func getTestContext() context.Context {
	return orc8rprotos.NewGatewayIdentity("test", "test", "test").NewContextWithIdentity(context.Background())
}
