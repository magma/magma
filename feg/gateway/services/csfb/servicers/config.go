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

package servicers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"

	mconfProtos "magma/feg/cloud/go/protos/mconfig"
	utilsDiammeter "magma/feg/gateway/diameter"
	fegMconfig "magma/gateway/mconfig"
	"magma/lte/cloud/go/protos/mconfig"
)

const (
	DefaultMMEName  = ".mmec01.mmegi0001.mme.epc.mnc001.mcc001.3gppnetwork.org"
	MMECLength      = 2
	MMEGILength     = 4
	MNCLength       = 3
	MCCLength       = 3
	MMEServiceName  = "mme"
	CsfbServiceName = "csfb"
)

func GetCsfbConfig() *mconfProtos.CsfbConfig {
	config := &mconfProtos.CsfbConfig{}
	err := fegMconfig.GetServiceConfigs(CsfbServiceName, config)
	if err != nil || config.Client == nil {
		glog.V(2).Infof("%s Managed Configs Load Error: %v", CsfbServiceName, err)
		return envOrDefaultConfig()
	}
	glog.V(2).Infof("Loaded %s configs: %+v", CsfbServiceName, *config)
	return config
}

func envOrDefaultConfig() *mconfProtos.CsfbConfig {
	defaultVlrIpAndPort := fmt.Sprintf("%s:%d", DefaultVLRIPAddress, DefaultVLRPort)
	defaultLocalIpAndPort := fmt.Sprintf("%s:%d", LocalIPAddress, LocalPort)

	// TODO: move GetValueOrEnv out of diameter. It should belong to some kind of utils module
	return &mconfProtos.CsfbConfig{
		Client: &mconfProtos.SCTPClientConfig{
			ServerAddress: utilsDiammeter.GetValueOrEnv("", VLRAddrEnv, defaultVlrIpAndPort),
			LocalAddress:  utilsDiammeter.GetValueOrEnv("", LocalAddrEnv, defaultLocalIpAndPort),
		},
	}
}

// TODO: MME should be gathered from the gRPC request not hardcoded/fixed per feg
// ConstructMMEName constructs MME name from mconfig
func ConstructMMEName() (string, error) {
	mmeConfig, err := getMMEConfig()
	if err != nil {
		glog.V(2).Infof(
			"Failed to retrieve MME config: %s, using default MME name: %s",
			err,
			DefaultMMEName,
		)
		return DefaultMMEName, nil
	}

	mnc := mmeConfig.GetCsfbMnc()
	mnc = fieldLengthCorrection(mnc, MNCLength)

	mcc := mmeConfig.GetCsfbMcc()
	mcc = fieldLengthCorrection(mcc, MCCLength)

	mmeCode := strconv.Itoa(int(mmeConfig.GetMmeCode()))
	mmeCode = fieldLengthCorrection(mmeCode, MMECLength)

	mmeGid := strconv.Itoa(int(mmeConfig.GetMmeGid()))
	mmeGid = fieldLengthCorrection(mmeGid, MMEGILength)

	mmeName := fmt.Sprintf(
		".mmec%s.mmegi%s.mme.epc.mnc%s.mcc%s.3gppnetwork.org",
		mmeCode,
		mmeGid,
		mnc,
		mcc,
	)

	return mmeName, nil
}

func getMMEConfig() (*mconfig.MME, error) {
	mmeConfig := &mconfig.MME{}
	err := fegMconfig.GetServiceConfigs(MMEServiceName, mmeConfig)
	if err != nil {
		return nil, err
	}
	return mmeConfig, nil
}

func fieldLengthCorrection(field string, requiredLength int) string {
	prefix := strings.Repeat("0", requiredLength-len(field))
	return prefix + field
}
