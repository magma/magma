/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/lte"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	utils "magma/lte/cloud/go/services/eps_authentication/servicers/test_utils"
	"magma/lte/cloud/go/services/eps_authentication/storage"
	"magma/lte/cloud/go/services/lte/obsidian/models"
	sdb_models "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	orc8r_storage "magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

// EpsAuthTestSuite is a test suite that will setup a test EPS auth server
// and config service pre-populated with a cellular config.
type EpsAuthTestSuite struct {
	suite.Suite
	Server *EPSAuthServer
}

var (
	qosProfileClassIDi32                  int32  = 9
	qosProfilePriorityLevelU32            uint32 = 15
	maxUlBitRateU32                       uint32 = 100000000
	maxDlBitRateU32                       uint32 = 200000000
	qosProfilePreemptionCapabilityBool    bool   = true
	qosProfilePreemptionVulnerabilityBool bool   = false
)

func (suite *EpsAuthTestSuite) AuthenticationInformation(
	air *fegprotos.AuthenticationInformationRequest) (*fegprotos.AuthenticationInformationAnswer, error) {

	return suite.Server.AuthenticationInformation(getTestContext(), air)
}

func (suite *EpsAuthTestSuite) UpdateLocation(
	ulr *fegprotos.UpdateLocationRequest) (*fegprotos.UpdateLocationAnswer, error) {
	return suite.Server.UpdateLocation(getTestContext(), ulr)
}

func (suite *EpsAuthTestSuite) PurgeUE(purge *fegprotos.PurgeUERequest) (*fegprotos.PurgeUEAnswer, error) {
	return suite.Server.PurgeUE(getTestContext(), purge)
}

func (*EpsAuthTestSuite) SetupTest() {
}

func TestEpsAuthSuite(t *testing.T) {
	test_init.StartTestService(t)

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
	err := configurator.CreateNetwork(context.Background(), configurator.Network{
		ID:      "test",
		Type:    "lte",
		Configs: map[string]interface{}{lte.CellularNetworkConfigType: cellularConfig},
	}, serdes.Network)
	assert.NoError(t, err)

	gw, err := configurator.CreateEntity(
		context.Background(),
		"test",
		configurator.NetworkEntity{Type: lte.CellularGatewayEntityType, Key: "test"},
		serdes.Entity)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(
		context.Background(), "test", configurator.NetworkEntity{
			Type: lte.APNEntityType, Key: "apn",
			Config: &models.ApnConfiguration{
				Ambr: &models.AggregatedMaximumBitrate{
					MaxBandwidthUl: &maxUlBitRateU32,
					MaxBandwidthDl: &maxDlBitRateU32,
				},
				QosProfile: &models.QosProfile{
					ClassID:                 &qosProfileClassIDi32,
					PreemptionCapability:    &qosProfilePreemptionCapabilityBool,
					PreemptionVulnerability: &qosProfilePreemptionVulnerabilityBool,
					PriorityLevel:           &qosProfilePriorityLevelU32,
				},
			},
		}, serdes.Entity)
	assert.NoError(t, err)

	var writes []configurator.EntityWriteOperation
	writes = append(writes, configurator.NetworkEntity{
		NetworkID: "test",
		Type:      lte.APNResourceEntityType,
		Key:       "resource",
		Config: &models.ApnResource{
			ApnName:    "apn",
			GatewayIP:  "172.16.254.1",
			GatewayMac: "00:0a:95:9d:68:16",
			ID:         "resource",
			VlanID:     42,
		},
		Associations: orc8r_storage.TKs{{Type: lte.APNEntityType, Key: "apn"}},
	})
	write := configurator.EntityUpdateCriteria{
		Type:              lte.CellularGatewayEntityType,
		Key:               gw.Key,
		AssociationsToAdd: orc8r_storage.TKs{{Type: lte.APNResourceEntityType, Key: "resource"}},
	}
	writes = append(writes, write)
	err = configurator.WriteEntities(context.Background(), "test", writes, serdes.Entity)
	assert.NoError(t, err)

	testSuite := &EpsAuthTestSuite{}

	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	stateStoreFactory := blobstore.NewSQLStoreFactory(storage.EpsAuthStateStore, db, sqorc.GetSqlBuilder())
	if err := stateStoreFactory.InitializeFactory(); err != nil {
		assert.NoError(t, err)
	}
	store := storage.NewSubscriberDBStorage(stateStoreFactory)

	for _, subscriber := range utils.GetTestSubscribers() {
		err := addSubscriber(store, subscriber, map[string]string{"apn": "1.2.3.4"})
		assert.NoError(t, err)
	}

	server, err := NewEPSAuthServer(store)
	assert.NoError(t, err)
	testSuite.Server = server

	suite.Run(t, testSuite)
}

func getTestContext() context.Context {
	return protos.NewGatewayIdentity(
		"test", "test", "test").NewContextWithIdentity(context.Background())
}

func addSubscriber(store storage.SubscriberDBStorage, sd *lteprotos.SubscriberData, staticIps map[string]string) error {
	ent := configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: lteprotos.SidString(sd.GetSid()),
		Config: &sdb_models.SubscriberConfig{},
	}
	if sd.GetLte() != nil {
		subProfile := sdb_models.SubProfile(sd.GetSubProfile())
		ent.Config = &sdb_models.SubscriberConfig{
			Lte: &sdb_models.LteSubscription{
				AuthAlgo:   sd.GetLte().GetAuthAlgo().String(),
				AuthKey:    sd.GetLte().GetAuthKey(),
				AuthOpc:    sd.GetLte().GetAuthOpc(),
				State:      "ACTIVE",
				SubProfile: &subProfile,
			},
			StaticIps: staticIps,
		}
	}
	if acfg := sd.GetNon_3Gpp().GetApnConfig(); len(acfg) > 0 && len(acfg[0].GetServiceSelection()) > 1 {
		ent.Associations = orc8r_storage.TKs{{Type: lte.APNEntityType, Key: acfg[0].GetServiceSelection()}}
	}
	_, err := configurator.CreateEntities(
		context.Background(), sd.GetNetworkId().GetId(), []configurator.NetworkEntity{ent}, serdes.Entity)
	return err
}
