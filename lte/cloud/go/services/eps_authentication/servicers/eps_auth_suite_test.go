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

	"magma/lte/cloud/go/lte"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	utils "magma/lte/cloud/go/services/eps_authentication/servicers/test_utils"
	"magma/lte/cloud/go/services/eps_authentication/storage"
	eps_storage "magma/lte/cloud/go/services/eps_authentication/storage"
	"magma/lte/cloud/go/services/lte/obsidian/models"
	sdb_models "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"
)

// EpsAuthTestSuite is a test suite that will setup a test EPS auth server
// and config service pre-populated with a cellular config.
type EpsAuthTestSuite struct {
	suite.Suite
	Server *EPSAuthServer
}

func (suite *EpsAuthTestSuite) AuthenticationInformation(
	air *lteprotos.AuthenticationInformationRequest) (*lteprotos.AuthenticationInformationAnswer, error) {

	return suite.Server.AuthenticationInformation(getTestContext(), air)
}

func (suite *EpsAuthTestSuite) UpdateLocation(
	ulr *lteprotos.UpdateLocationRequest) (*lteprotos.UpdateLocationAnswer, error) {
	return suite.Server.UpdateLocation(getTestContext(), ulr)
}

func (suite *EpsAuthTestSuite) PurgeUE(purge *lteprotos.PurgeUERequest) (*lteprotos.PurgeUEAnswer, error) {
	return suite.Server.PurgeUE(getTestContext(), purge)
}

func (suite *EpsAuthTestSuite) SetupTest() {
	db, err := sqorc.Open("sqlite3", ":memory:")
	suite.NoError(err)
	stateStoreFactory := blobstore.NewSQLStoreFactory(eps_storage.EpsAuthStateStore, db, sqorc.GetSqlBuilder())
	if err := stateStoreFactory.InitializeFactory(); err != nil {
		suite.NoError(err)
	}
	store := storage.NewSubscriberDBStorage(stateStoreFactory)

	for _, subscriber := range utils.GetTestSubscribers() {
		err := addSubscriber(store, subscriber)
		suite.NoError(err)
	}

	server, err := NewEPSAuthServer(store)
	suite.NoError(err)
	suite.Server = server
}

// DEPRECATED -- un-skip if service is un-deprecated
func TestEpsAuthSuite(t *testing.T) {
	t.Skip("eps_authentication service temporarily deprecated")

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

	testSuite := &EpsAuthTestSuite{}
	suite.Run(t, testSuite)
}

func getTestContext() context.Context {
	return protos.NewGatewayIdentity(
		"test", "test", "test").NewContextWithIdentity(context.Background())
}

func addSubscriber(store storage.SubscriberDBStorage, sd *lteprotos.SubscriberData) error {
	ent := configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: lteprotos.SidString(sd.GetSid()),
		Config: &sdb_models.SubscriberConfig{
			Lte: &sdb_models.LteSubscription{
				AuthAlgo:   sd.GetLte().GetAuthAlgo().String(),
				AuthKey:    sd.GetLte().GetAuthKey(),
				AuthOpc:    sd.GetLte().GetAuthOpc(),
				State:      "ACTIVE",
				SubProfile: sdb_models.SubProfile(sd.GetSubProfile()),
			}},
	}
	_, err := configurator.CreateEntities(
		context.Background(), sd.GetNetworkId().GetId(), []configurator.NetworkEntity{ent}, serdes.Entity)
	return err
}
