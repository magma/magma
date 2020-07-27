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

// Package directoryd provides a client API for interacting with the
// directory cloud service, which manages the UE location information
package directoryd

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	platformregistry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/util"
)

const (
	ServiceName = "DIRECTORYD"
	ImsiPrefix  = "IMSI"

	UseCloudDirectordEnv = "USE_CLOUD_DIRECTORYD"
)

var useCloudDirectoryd = util.GetEnvBool(UseCloudDirectordEnv)

// Get a thin RPC client to the gateway directory service.
func GetGatewayDirectorydClient() (protos.GatewayDirectoryServiceClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	if useCloudDirectoryd {
		conn, err = platformregistry.Get().GetSharedCloudConnection(strings.ToLower(ServiceName))
	} else {
		conn, err = platformregistry.Get().GetConnection(ServiceName)
	}
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewGatewayDirectoryServiceClient(conn), nil
}

// UpdateRecord updates the directory record for the provided ID with the calling
// GW's HW ID and any associated identifiers
func UpdateRecord(request *protos.UpdateRecordRequest) error {
	if len(request.GetId()) == 0 {
		return fmt.Errorf("Empty ID")
	}
	client, err := GetGatewayDirectorydClient()
	if err != nil {
		return err
	}
	request.Id = PrependImsiPrefix(request.GetId())
	_, err = client.UpdateRecord(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return err
}

// DeleteRecord deletes the directory record for the provided ID
func DeleteRecord(request *protos.DeleteRecordRequest) error {
	if len(request.GetId()) == 0 {
		return fmt.Errorf("Empty ID")
	}
	client, err := GetGatewayDirectorydClient()
	if err != nil {
		return err
	}
	request.Id = PrependImsiPrefix(request.GetId())
	_, err = client.DeleteRecord(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return err
}

func PrependImsiPrefix(imsi string) string {
	if !strings.HasPrefix(imsi, ImsiPrefix) {
		imsi = ImsiPrefix + imsi
	}
	return imsi
}
