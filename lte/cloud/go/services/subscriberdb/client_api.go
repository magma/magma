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

// Package client provides a thin client for contacting the subscriberdb service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package subscriberdb

import (
	"context"

	"magma/lte/cloud/go/services/subscriberdb/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// ListMSISDNs returns the tracked MSISDNs, keyed by their associated IMSI.
func ListMSISDNs(networkID string) (map[string]string, error) {
	msisdns := map[string]string{}

	client, err := getClient()
	if err != nil {
		return msisdns, err
	}

	res, err := client.GetMSISDNs(
		context.Background(),
		&protos.GetMSISDNsRequest{
			NetworkId: networkID,
			Msisdns:   nil, // list all
		},
	)
	if err != nil {
		return msisdns, err
	}

	return res.ImsisByMsisdn, nil
}

// GetIMSIForMSISDN returns the IMSI associated with the passed MSISDN.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/errors.
func GetIMSIForMSISDN(networkID, msisdn string) (string, error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}

	res, err := client.GetMSISDNs(
		context.Background(),
		&protos.GetMSISDNsRequest{
			NetworkId: networkID,
			Msisdns:   []string{msisdn},
		},
	)
	if err != nil {
		return "", err
	}

	msisdn, ok := res.ImsisByMsisdn[msisdn]
	if !ok {
		return "", merrors.ErrNotFound
	}

	return msisdn, nil
}

func SetIMSIForMSISDN(networkID, msisdn, imsi string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	_, err = client.SetMSISDN(
		context.Background(),
		&protos.SetMSISDNRequest{
			NetworkId: networkID,
			Msisdn:    msisdn,
			Imsi:      imsi,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMSISDN(networkID, msisdn string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	_, err = client.DeleteMSISDN(
		context.Background(),
		&protos.DeleteMSISDNRequest{
			NetworkId: networkID,
			Msisdn:    msisdn,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func getClient() (protos.SubscriberLookupClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewSubscriberLookupClient(conn), nil
}
