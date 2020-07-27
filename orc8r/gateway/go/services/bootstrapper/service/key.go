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

// package service implements the core of bootstrapper

package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"

	"magma/gateway/config"
	"magma/orc8r/lib/go/security/key"
)

// GetChallengeKey reads and returns the bootstrapper challenge key if present
// or creates and returns the key if the key file doesn't exist
func GetChallengeKey() (privKey interface{}, err error) {
	challengeKeyFile := config.GetMagmadConfigs().BootstrapConfig.ChallengeKey
	privKey, err = key.ReadKey(challengeKeyFile)
	if err == nil {
		return // all good, return the key
	}
	glog.Warningf("Bootstrapper ReadKey(%s) error: %v", challengeKeyFile, err)

	if os.IsNotExist(err) { // file doesn't exist, check default location or try to create it
		dir := filepath.Dir(challengeKeyFile)
		if len(dir) > 3 {
			os.MkdirAll(dir, 0755)
		}
		// try default location first
		if challengeKeyFile != config.DefaultChallengeKeyFile {
			if privKey, err = key.ReadKey(config.DefaultChallengeKeyFile); err == nil {
				if err = key.WriteKey(challengeKeyFile, privKey); err == nil {
					glog.Warningf(
						"default challenge key '%s' was copied into configured location: %s",
						config.DefaultChallengeKeyFile, challengeKeyFile)
					if privKey, err = key.ReadKey(challengeKeyFile); err == nil {
						return // copied & verified challenge key
					}
				}
			}
		}
		// create new
		privKey, err = key.GenerateKey(PrivateKeyType, 0)
		if err != nil {
			err = fmt.Errorf("Bootstrapper Generate Key error: %v", err)
			return
		}
		err = key.WriteKey(challengeKeyFile, privKey)
		if err != nil {
			err = fmt.Errorf("Bootstrapper Write Key (%s) error: %v", challengeKeyFile, err)
			return
		}
		privKey, err = key.ReadKey(challengeKeyFile)
		if err != nil {
			err = fmt.Errorf(
				"Bootstrapper Failed to read recently created key from (%s) error: %v", challengeKeyFile, err)
			return
		}
		glog.Infof("successfuly created new challenge key file: %s", challengeKeyFile)
	}
	return
}
