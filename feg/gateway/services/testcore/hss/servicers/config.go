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
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	configs "magma/gateway/mconfig"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

// HSS Flag Variables to overwrite default configs
const (
	hssServiceName      = "hss"
	hssDefaultProtocol  = "tcp"
	hssDefaultHost      = "magma.com"
	hssDefaultRealm     = "magma.com"
	maxUlBitRateFlag    = "max_ul_bit_rate"
	maxDlBitRateFlag    = "max_dl_bit_rate"
	defaultMaxUlBitRate = uint64(100000000)
	defaultMaxDlBitRate = uint64(200000000)
)

var (
	hssDefaultLteAuthAmf  = []byte("\x80\x00")
	hssDefaultLteAuthOp   = []byte("\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18")
	streamSubscribersFlag = flag.Bool("stream_subscribers", false, "Whether to stream subscribers from the cloud")
)

func init() {
	flag.Uint64(maxUlBitRateFlag, defaultMaxUlBitRate, "Maximum uplink bit rate (AMBR-UL)")
	flag.Uint64(maxDlBitRateFlag, defaultMaxDlBitRate, "Maximum downlink bit rate (AMBR-DL)")
}

// GetHSSConfig returns the server config for an HSS based on the input flags
func GetHSSConfig() (*mconfig.HSSConfig, error) {
	serviceBaseName := filepath.Base(os.Args[0])
	serviceBaseName = strings.TrimSuffix(serviceBaseName, filepath.Ext(serviceBaseName))
	if hssServiceName != serviceBaseName {
		glog.Errorf(
			"NOTE: HSS Service name: %s does not match its managed configs key: %s\n",
			serviceBaseName, hssServiceName)
	}

	configsPtr := &mconfig.HSSConfig{}
	err := configs.GetServiceConfigs(hssServiceName, configsPtr)
	if err != nil || configsPtr.Server == nil || configsPtr.DefaultSubProfile == nil {
		glog.Errorf("%s Managed Configs Load Error: %v\n", hssServiceName, err)
		return &mconfig.HSSConfig{
			Server: &mconfig.DiamServerConfig{
				Address:      diameter.GetValue(diameter.AddrFlag, ""),
				Protocol:     diameter.GetValue(diameter.NetworkFlag, hssDefaultProtocol),
				LocalAddress: diameter.GetValue(diameter.LocalAddrFlag, ""),
				DestHost:     diameter.GetValue(diameter.DestHostFlag, hssDefaultHost),
				DestRealm:    diameter.GetValue(diameter.DestRealmFlag, hssDefaultRealm),
			},
			LteAuthOp:  hssDefaultLteAuthOp,
			LteAuthAmf: hssDefaultLteAuthAmf,
			DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
				MaxUlBitRate: diameter.GetValueUint64(maxUlBitRateFlag, defaultMaxUlBitRate),
				MaxDlBitRate: diameter.GetValueUint64(maxDlBitRateFlag, defaultMaxDlBitRate),
			},
			SubProfiles:       make(map[string]*mconfig.HSSConfig_SubscriptionProfile),
			StreamSubscribers: *streamSubscribersFlag,
		}, err
	}

	glog.V(2).Infof("Loaded %s configs: %+v\n", hssServiceName, *configsPtr)
	return &mconfig.HSSConfig{
		Server: &mconfig.DiamServerConfig{
			Address:      diameter.GetValue(diameter.AddrFlag, configsPtr.Server.Address),
			Protocol:     diameter.GetValue(diameter.NetworkFlag, configsPtr.Server.Protocol),
			LocalAddress: diameter.GetValue(diameter.LocalAddrFlag, configsPtr.Server.LocalAddress),
			DestHost:     diameter.GetValue(diameter.DestHostFlag, configsPtr.Server.DestHost),
			DestRealm:    diameter.GetValue(diameter.DestRealmFlag, configsPtr.Server.DestRealm),
		},
		LteAuthOp:  configsPtr.LteAuthOp,
		LteAuthAmf: configsPtr.LteAuthAmf,
		DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
			MaxUlBitRate: diameter.GetValueUint64(maxUlBitRateFlag, configsPtr.DefaultSubProfile.MaxUlBitRate),
			MaxDlBitRate: diameter.GetValueUint64(maxDlBitRateFlag, configsPtr.DefaultSubProfile.MaxDlBitRate),
		},
		SubProfiles:       configsPtr.SubProfiles,
		StreamSubscribers: configsPtr.StreamSubscribers || *streamSubscribersFlag,
	}, nil
}

// GetConfiguredSubscribers returns a slice of subscribers configured in hss.yml
func GetConfiguredSubscribers() ([]*protos.SubscriberData, error) {
	hsscfg, err := config.GetServiceConfig("", hssServiceName)
	if err != nil {
		return nil, err
	}
	subscribers, ok := hsscfg.RawMap["subscribers"]
	if !ok {
		return nil, fmt.Errorf("Could not find 'subscribers' in config file")
	}
	rawMap, ok := subscribers.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to convert %T to map %v", subscribers, rawMap)
	}
	var subscriberData []*protos.SubscriberData
	for k, v := range rawMap {
		imsi, ok := k.(string)
		if !ok {
			continue
		}
		rawMap, ok := v.(map[interface{}]interface{})
		if !ok {
			continue
		}
		configMap := &config.ConfigMap{RawMap: rawMap}

		// If auth_key is incorrect, skip subscriber
		authKey, err := configMap.GetString("auth_key")
		if err != nil {
			glog.Errorf("Could not add subscriber due to missing auth_key: %s", err)
			continue
		}
		authKeyBytes, err := hex.DecodeString(authKey)
		if err != nil {
			glog.Errorf("Could not add subscriber due to incorrect auth key format: %s", err)
			continue
		}
		non3gppEnabled, err := configMap.GetBool("non_3gpp_enabled")
		if err != nil {
			non3gppEnabled = true
		}
		subscriberData = append(subscriberData, createSubscriber(imsi, authKeyBytes, non3gppEnabled))
	}
	return subscriberData, err
}

func createSubscriber(imsi string, authKey []byte, non3gppEnabled bool) *protos.SubscriberData {
	var non3gppProfile *protos.Non3GPPUserProfile
	if non3gppEnabled {
		non3gppProfile = &protos.Non3GPPUserProfile{
			Msisdn:              msisdn,
			Non_3GppIpAccess:    protos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_ALLOWED,
			Non_3GppIpAccessApn: protos.Non3GPPUserProfile_NON_3GPP_APNS_ENABLE,
			ApnConfig:           []*protos.APNConfiguration{{}},
		}
	} else {
		non3gppProfile = &protos.Non3GPPUserProfile{
			Msisdn:              msisdn,
			Non_3GppIpAccess:    protos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_BARRED,
			Non_3GppIpAccessApn: protos.Non3GPPUserProfile_NON_3GPP_APNS_DISABLE,
			ApnConfig:           []*protos.APNConfiguration{{}},
		}
	}
	return &protos.SubscriberData{
		Sid: &protos.SubscriberID{Id: imsi},
		Gsm: &protos.GSMSubscription{State: protos.GSMSubscription_ACTIVE},
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthKey:  authKey,
			AuthAlgo: protos.LTESubscription_MILENAGE,
		},
		State: &protos.SubscriberState{
			TgppAaaServerRegistered: false,
		},
		Non_3Gpp: non3gppProfile,
	}
}
