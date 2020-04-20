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

// package service implements the core of configurator
package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/glog"

	"magma/gateway/config"
	"magma/gateway/mconfig"
)

// SaveConfig saves new gateway.configs and returns old configuration if any
func SaveConfigs(cfgJson []byte) error {
	if len(cfgJson) == 0 {
		return fmt.Errorf("empty gateway mconfigs")
	}
	err := safeSwap(mconfig.ConfigFilePath(), cfgJson)
	if err == nil {
		glog.V(1).Infof("successfully updated mconfig in %s", mconfig.ConfigFilePath())
	}
	return err
}

// updateStaticConfigs saves new gateway.configs into static mconfig location
func updateStaticConfigs(cfgJson []byte) error {
	intervalMin := config.GetMagmadConfigs().StaticMconfigUpdateIntervalMin
	if intervalMin <= 0 {
		return nil
	}
	intervalDuration := time.Duration(intervalMin) * time.Minute
	mconfigPath := mconfig.DefaultConfigFilePath()
	fi, err := os.Stat(mconfigPath)
	if (err == nil && fi.ModTime().Add(intervalDuration).Before(time.Now())) || os.IsNotExist(err) {
		return safeSwap(mconfigPath, cfgJson)
	}
	return nil
}

func safeSwap(mconfigPath string, cfgJson []byte) error {
	newMconfigPath := mconfigPath + ".new"
	oldMconfigPath := mconfigPath + ".old"
	err := ioutil.WriteFile(newMconfigPath, cfgJson, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			if os.MkdirAll(filepath.Dir(mconfigPath), 0755) == nil {
				err = ioutil.WriteFile(newMconfigPath, cfgJson, 0644)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to save mconfigs into %s: %v", newMconfigPath, err)
		}
	}
	oerr := os.Rename(mconfigPath, oldMconfigPath) // best effort, needed just for rollback on error
	err = os.Rename(newMconfigPath, mconfigPath)
	if err != nil {
		err = fmt.Errorf("failed to move mconfigs from %s to %s: %v", newMconfigPath, mconfigPath, err)
		if oerr == nil { // roll back if previous rename succeeded
			os.Rename(oldMconfigPath, mconfigPath)
		}
	}
	return err
}

func readCfg() ([]byte, error) {
	return ioutil.ReadFile(mconfig.ConfigFilePath())
}
