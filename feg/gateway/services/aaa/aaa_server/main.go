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

// Package main implements WiFi AAA server
package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/servicers"
	"magma/feg/gateway/services/aaa/store"
	"magma/feg/gateway/utils"
	managed_configs "magma/gateway/mconfig"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/service"
)

const (
	AAAServiceName                 = "aaa_server"
	Version                        = "0.1"
	AccountingReportingEnabledFlag = "acct_reporting_enabled"
	AccountingReportingEnabledEnv  = "AAA_ACCT_REPORTING_ENABLED"
)

var (
	_ = flag.Bool(AccountingReportingEnabledFlag, false, "Enable base accounting reports")
)

func main() {
	flag.Parse() // for glog

	// Create a shared Session Table
	sessions := store.NewMemorySessionTable()

	// Create the EAP AKA Provider service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.AAA_SERVER)
	if err != nil {
		glog.Fatalf("Error creating AAA service: %s", err)
	}
	aaaConfigs := &mconfig.AAAConfig{}
	err = managed_configs.GetServiceConfigs(AAAServiceName, aaaConfigs)
	if err != nil {
		glog.Warningf("Error getting AAA Server service configs: %s", err)
		aaaConfigs = &mconfig.AAAConfig{}
	}
	aaaConfigs.AcctReportingEnabled = utils.GetBoolValueOrEnv(
		AccountingReportingEnabledFlag, AccountingReportingEnabledEnv, aaaConfigs.AcctReportingEnabled)

	acct, _ := servicers.NewAccountingService(sessions, proto.Clone(aaaConfigs).(*mconfig.AAAConfig))
	protos.RegisterAccountingServer(srv.GrpcServer, acct)
	lteprotos.RegisterAbortSessionResponderServer(srv.GrpcServer, acct)
	fegprotos.RegisterSwxGatewayServiceServer(srv.GrpcServer, acct)
	fegprotos.RegisterS6AGatewayServiceServer(srv.GrpcServer, acct)

	auth, _ := servicers.NewEapAuthenticator(sessions, aaaConfigs, acct)
	protos.RegisterAuthenticatorServer(srv.GrpcServer, auth)

	// Starts built in radius server if built with this option
	startBuiltInRadius(aaaConfigs, auth, acct)

	glog.Infof("Starting AAA Service v%s.", Version)
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running AAA service: %s", err)
	}
}
