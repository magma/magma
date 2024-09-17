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
package main

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/golang/glog"

	"magma/dp/cloud/go/dp"
	"magma/dp/cloud/go/protos"
	dp_service "magma/dp/cloud/go/services/dp"
	"magma/dp/cloud/go/services/dp/active_mode_controller"
	amc_time "magma/dp/cloud/go/services/dp/active_mode_controller/time"
	"magma/dp/cloud/go/services/dp/logs_pusher"
	"magma/dp/cloud/go/services/dp/obsidian/cbsd"
	dp_log "magma/dp/cloud/go/services/dp/obsidian/log"
	"magma/dp/cloud/go/services/dp/servicers"
	dp_storage "magma/dp/cloud/go/services/dp/storage"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/obsidian"
	swagger_protos "magma/orc8r/cloud/go/services/obsidian/swagger/protos"
	swagger_servicers "magma/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/service/config"
)

func main() {
	srv, err := service.NewOrchestratorService(dp.ModuleName, dp_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating %s service: %s", dp_service.ServiceName, err)
	}
	obsidian.AttachHandlers(srv.EchoServer, cbsd.GetHandlers())
	obsidian.AttachHandlers(srv.EchoServer, dp_log.NewHandlersGetter(dp_log.GetElasticClient, "").GetHandlers())
	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swagger_servicers.NewSpecServicerFromFile(dp_service.ServiceName))

	var serviceConfig dp_service.Config
	config.MustGetStructuredServiceConfig(dp.ModuleName, dp_service.ServiceName, &serviceConfig)

	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %s", err)
	}
	cbsdStore := dp_storage.NewCbsdManager(db, sqorc.GetSqlBuilder(), sqorc.GetErrorChecker(), sqorc.GetSqlLocker())

	dpCfg := serviceConfig.DpBackend
	interval := time.Second * time.Duration(dpCfg.CbsdInactivityIntervalSec)
	logConsumerUrl := dpCfg.LogConsumerUrl

	protos.RegisterCbsdManagementServer(srv.GrpcServer, servicers.NewCbsdManager(cbsdStore, interval, logConsumerUrl, logs_pusher.PushDPLog))

	cancel, errs := startAmc(db, serviceConfig.ActiveModeController)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running %s service amd echo server: %s", dp_service.ServiceName, err)
	}

	stopAmc(cancel, errs)
}

func startAmc(db *sql.DB, cfg *dp_service.AmcConfig) (context.CancelFunc, chan error) {
	clock := &amc_time.Clock{}
	seed := rand.NewSource(clock.Now().Unix())
	amcManager := dp_storage.NewAmcManager(db, sqorc.GetSqlBuilder(), sqorc.GetErrorChecker(), sqorc.GetSqlLocker())
	app := active_mode_controller.NewApp(
		active_mode_controller.WithDb(db),
		active_mode_controller.WithAmcManager(amcManager),
		active_mode_controller.WithClock(clock),
		active_mode_controller.WithRNG(rand.New(seed)),
		active_mode_controller.WithHeartbeatSendTimeout(
			secToDuration(cfg.HeartbeatSendTimeoutSec),
			secToDuration(cfg.RequestProcessingIntervalSec)),
		active_mode_controller.WithPollingInterval(secToDuration(cfg.PollingIntervalSec)),
		active_mode_controller.WithCbsdInactivityTimeout(secToDuration(cfg.CbsdInactivityTimeoutSec)),
	)
	errs := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { errs <- app.Run(ctx) }()
	return cancel, errs
}

func secToDuration(s int) time.Duration {
	return time.Second * time.Duration(s)
}

func stopAmc(cancel context.CancelFunc, errs chan error) {
	cancel()
	err := <-errs
	if err != nil && err != context.Canceled {
		glog.Fatalf("Error while shutting down amc: %s", err)
	}
}
