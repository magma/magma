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

package device

import (
	"context"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/device/protos"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func RegisterDevice(networkID, deviceType, deviceKey string, info interface{}, serdes serde.Registry) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	serializedInfo, err := serde.Serialize(info, deviceType, serdes)
	if err != nil {
		return err
	}
	entity := &protos.PhysicalEntity{
		DeviceID: deviceKey,
		Type:     deviceType,
		Info:     serializedInfo,
	}
	req := &protos.RegisterOrUpdateDevicesRequest{
		NetworkID: networkID,
		Entities:  []*protos.PhysicalEntity{entity},
	}
	_, err = client.RegisterDevices(context.Background(), req)
	return err
}

func UpdateDevice(networkID, deviceType, deviceKey string, info interface{}, serdes serde.Registry) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	serializedInfo, err := serde.Serialize(info, deviceType, serdes)
	if err != nil {
		return err
	}
	entity := &protos.PhysicalEntity{
		DeviceID: deviceKey,
		Type:     deviceType,
		Info:     serializedInfo,
	}
	req := &protos.RegisterOrUpdateDevicesRequest{
		NetworkID: networkID,
		Entities:  []*protos.PhysicalEntity{entity},
	}
	_, err = client.UpdateDevices(context.Background(), req)
	return err
}

func DeleteDevices(networkID string, ids []storage.TypeAndKey) error {
	client, err := getDeviceClient()
	if err != nil {
		return err
	}

	requestIDs := funk.Map(
		ids,
		func(id storage.TypeAndKey) *protos.DeviceID {
			return &protos.DeviceID{Type: id.Type, DeviceID: id.Key}
		},
	).([]*protos.DeviceID)

	req := &protos.DeleteDevicesRequest{NetworkID: networkID, DeviceIDs: requestIDs}
	_, err = client.DeleteDevices(context.Background(), req)
	return err
}

func DeleteDevice(networkID, deviceType, deviceKey string) error {
	return DeleteDevices(networkID, []storage.TypeAndKey{{Type: deviceType, Key: deviceKey}})
}

func GetDevice(networkID, deviceType, deviceKey string, serdes serde.Registry) (interface{}, error) {
	device, err := getDevice(networkID, deviceType, deviceKey)
	if err != nil {
		return nil, err
	}
	return serde.Deserialize(device.Info, deviceType, serdes)
}

func GetDevices(networkID string, deviceType string, deviceIDs []string, serdes serde.Registry) (map[string]interface{}, error) {
	if len(deviceIDs) == 0 {
		return map[string]interface{}{}, nil
	}
	client, err := getDeviceClient()
	if err != nil {
		return nil, err
	}

	requestIDs := funk.Map(
		deviceIDs,
		func(id string) *protos.DeviceID { return &protos.DeviceID{Type: deviceType, DeviceID: id} },
	).([]*protos.DeviceID)
	req := &protos.GetDeviceInfoRequest{NetworkID: networkID, DeviceIDs: requestIDs}
	res, err := client.GetDeviceInfo(context.Background(), req)
	if err != nil {
		return map[string]interface{}{}, err
	}

	ret := make(map[string]interface{}, len(res.DeviceMap))
	for k, val := range res.DeviceMap {
		iVal, err := serde.Deserialize(val.Info, deviceType, serdes)
		if err != nil {
			return map[string]interface{}{}, errors.Wrapf(err, "failed to deserialize device %s", k)
		}
		ret[k] = iVal
	}
	return ret, nil
}

func DoesDeviceExist(networkID, deviceType, deviceID string) (bool, error) {
	_, err := getDevice(networkID, deviceType, deviceID)
	if err == merrors.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func getDevice(networkID, deviceType, deviceKey string) (*protos.PhysicalEntity, error) {
	client, err := getDeviceClient()
	if err != nil {
		return nil, err
	}
	deviceID := &protos.DeviceID{Type: deviceType, DeviceID: deviceKey}
	req := &protos.GetDeviceInfoRequest{NetworkID: networkID, DeviceIDs: []*protos.DeviceID{deviceID}}
	res, err := client.GetDeviceInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	device, ok := res.DeviceMap[deviceKey]
	if !ok {
		return nil, merrors.ErrNotFound
	}
	return device, nil
}

func getDeviceClient() (protos.DeviceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDeviceClient(conn), err
}
