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

// Magma's Envoy Controller Service configures Envoy proxy with specified configuration
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"regexp"
	"sort"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/envoy_controller/control_plane"
	"magma/feg/gateway/services/envoy_controller/servicers"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	flag.Set("v", "2")
	flag.Parse()

	logPath := "/var/log/"
	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		glog.Infof("Failed reading %s dir: %s", logPath, err)
		return
	}

	maxLogFiles := 5
	levels := []string{"FATAL", "ERROR", "WARNING", "INFO"}
	for _, logSeverity := range levels {
		logFiles := []string{}
		for _, f := range files {
			match, _ := regexp.MatchString(".*envoy_controller.*log."+logSeverity+".*", f.Name())
			if match {
				logFiles = append(logFiles, f.Name())
			}
		}
		if len(logFiles) < maxLogFiles {
			continue
		}
		sort.Sort(sort.Reverse(sort.StringSlice(logFiles)))
		for _, file := range logFiles[maxLogFiles:] {
			err := os.Remove(logPath + file)
			if err != nil {
				glog.Infof("Failed removing file %s: %s", file, err)
			}
		}
	}
}

func main() {
	// Create the service
	glog.Infof("Creating '%s' Service", registry.ENVOY_CONTROLLER)
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.ENVOY_CONTROLLER)
	if err != nil {
		glog.Fatalf("Error creating Envoy Controller service: %s", err)
	}

	cli := control_plane.GetControllerClient()

	// Create servicers
	servicer := servicers.NewEnvoyControllerService(cli)

	// Register services
	protos.RegisterEnvoyControllerServer(srv.GrpcServer, servicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
	glog.Infof("Starting '%s' Service", registry.ENVOY_CONTROLLER)
}
