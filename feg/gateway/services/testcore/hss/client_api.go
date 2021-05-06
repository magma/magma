/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package hss provides a thin client for using the hss service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package hss

import (
	"context"
	"errors"
	"fmt"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// Wrapper for GRPC Client functionality
type hssClient struct {
	fegprotos.HSSConfiguratorClient
	cc *grpc.ClientConn
}

// getHSSClient is a utility function to get a RPC connection to the
// HSS service.
func getHSSClient() (*hssClient, error) {
	conn, err := registry.GetConnection(registry.MOCK_HSS)
	if err != nil {
		errMsg := fmt.Sprintf("HSS client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &hssClient{
		fegprotos.NewHSSConfiguratorClient(conn),
		conn,
	}, err
}

// AddSubscriber tries to add this subscriber to the server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func AddSubscriber(sub *lteprotos.SubscriberData) error {
	err := VerifySubscriberData(sub)
	if err != nil {
		errMsg := fmt.Errorf("Invalid AddSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.AddSubscriber(context.Background(), sub)
	return err
}

// GetSubscriberData looks up a subscriber by their Id.
// If the subscriber cannot be found, an error is returned instead.
// Input: The id of the subscriber to be looked up.
// Output: The data of the corresponding subscriber.
func GetSubscriberData(id string) (*lteprotos.SubscriberData, error) {
	err := verifyID(id)
	if err != nil {
		errMsg := fmt.Errorf("Invalid GetSubscriberDataRequest provided: %s", err)
		return nil, errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return nil, err
	}
	subID := &lteprotos.SubscriberID{
		Id: id,
	}
	return cli.GetSubscriberData(context.Background(), subID)
}

// UpdateSubscriber changes the data stored for an existing subscriber.
// If the subscriber cannot be found, an error is returned instead.
// Input: The new subscriber data to store.
func UpdateSubscriber(sub *lteprotos.SubscriberData) error {
	err := VerifySubscriberData(sub)
	if err != nil {
		errMsg := fmt.Errorf("Invalid UpdateSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.UpdateSubscriber(context.Background(), sub)
	return err
}

// DeleteSubscriber deletes a subscriber by their Id.
// If the subscriber is not found, then this call is ignored.
// Input: The id of the subscriber to be deleted.
func DeleteSubscriber(id string) error {
	err := verifyID(id)
	if err != nil {
		errMsg := fmt.Errorf("Invalid DeleteSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	subID := &lteprotos.SubscriberID{
		Id: id,
	}
	_, err = cli.DeleteSubscriber(context.Background(), subID)
	return err
}

// DeRegisterSubscriber de-registers a subscriber by their Id.
// If the subscriber is not found, an error is returned instead.
// Input: The id of the subscriber to be deleted.
func DeregisterSubscriber(id string) error {
	err := verifyID(id)
	if err != nil {
		errMsg := fmt.Errorf("Invalid DeregisterSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	subID := &lteprotos.SubscriberID{
		Id: id,
	}
	_, err = cli.DeregisterSubscriber(context.Background(), subID)
	return err
}

func VerifySubscriberData(sub *lteprotos.SubscriberData) error {
	if sub == nil {
		return fmt.Errorf("subscriber is nil")
	}
	return verifyID(sub.Sid.GetId())
}

func verifyID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("no id provided")
	}
	return nil
}
