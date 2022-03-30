/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test_init

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_protos "magma/lte/cloud/go/services/lte/protos"
	"magma/lte/cloud/go/services/lte/servicers"
	protected_servicers "magma/lte/cloud/go/services/lte/servicers/protected"
	"magma/lte/cloud/go/services/lte/storage"
	"magma/orc8r/cloud/go/orc8r"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	state_protos "magma/orc8r/cloud/go/services/state/protos"
	provider_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	StartTestServiceWithConfig(t, lte_service.Config{DefaultSubscriberdbSyncInterval: 300})
}

func StartTestServiceWithConfig(t *testing.T, serviceConfig lte_service.Config) {
	streams := []string{
		lte.SubscriberStreamName,
		lte.PolicyStreamName,
		lte.BaseNameStreamName,
		lte.ApnRuleMappingsStreamName,
		lte.NetworkWideRulesStreamName,
		lte.RatingGroupStreamName,
	}
	labels := map[string]string{
		orc8r.MconfigBuilderLabel: "true",
		orc8r.StreamProviderLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StreamProviderStreamsAnnotation: strings.Join(streams, orc8r.AnnotationFieldSeparator),
	}

	srv, lis, plis := test_utils.NewTestOrchestratorService(t, lte.ModuleName, lte_service.ServiceName, labels, annotations)
	builder_protos.RegisterMconfigBuilderServer(srv.ProtectedGrpcServer, protected_servicers.NewBuilderServicer(serviceConfig))
	provider_protos.RegisterStreamProviderServer(srv.ProtectedGrpcServer, servicers.NewProviderServicer())

	// Init storage
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	enbStateStore := storage.NewEnodebStateLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, enbStateStore.Initialize())

	// Add servicers
	lte_protos.RegisterEnodebStateLookupServer(srv.ProtectedGrpcServer, protected_servicers.NewLookupServicer(enbStateStore))
	state_protos.RegisterIndexerServer(srv.ProtectedGrpcServer, protected_servicers.NewIndexerServicer())

	go srv.RunTest(lis, plis)
}
