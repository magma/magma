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

package servicers_test

import (
	"context"
	"net"
	"testing"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestHomeSubscriberServer_AddSubscriber(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)

	sub1 := protos.SubscriberData{Sid: &protos.SubscriberID{Id: "1"}}
	sub2 := protos.SubscriberData{Sid: &protos.SubscriberID{Id: "2"}}

	_, err := server.AddSubscriber(context.Background(), &sub1)
	assert.NoError(t, err)

	_, err = server.AddSubscriber(context.Background(), &sub1)
	assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = Subscriber '1' already exists")

	_, err = server.AddSubscriber(context.Background(), &sub2)
	assert.NoError(t, err)

	_, err = server.AddSubscriber(context.Background(), &sub1)
	assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = Subscriber '1' already exists")
}

func TestHomeSubscriberServer_GetSubscriberData(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)

	id1 := protos.SubscriberID{Id: "1"}
	sub1 := protos.SubscriberData{Sid: &id1}

	_, err := server.GetSubscriberData(context.Background(), &id1)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = Subscriber '1' not found")

	_, err = server.AddSubscriber(context.Background(), &sub1)
	assert.NoError(t, err)

	result, err := server.GetSubscriberData(context.Background(), &id1)
	assert.NoError(t, err)
	assert.Equal(t, sub1.String(), result.String())
}

func TestHomeSubscriberServer_UpdateSubscriber(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)

	_, err := server.UpdateSubscriber(context.Background(), nil)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Subscriber data cannot be nil")

	sub := &protos.SubscriberData{}
	_, err = server.UpdateSubscriber(context.Background(), sub)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Subscriber data must contain a subscriber id")

	id := &protos.SubscriberID{Id: "1"}
	sub = &protos.SubscriberData{Sid: id}
	_, err = server.UpdateSubscriber(context.Background(), sub)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = Subscriber '1' not found")

	_, err = server.AddSubscriber(context.Background(), sub)
	assert.NoError(t, err)

	updatedSub := &protos.SubscriberData{
		Sid:        id,
		SubProfile: "test",
	}
	_, err = server.UpdateSubscriber(context.Background(), updatedSub)
	assert.NoError(t, err)

	retreivedSub, err := server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, updatedSub, retreivedSub)
}

func TestHomeSubscriberServer_DeleteSubscriber(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)

	id := protos.SubscriberID{Id: "1"}
	sub := protos.SubscriberData{Sid: &id}

	_, err := server.AddSubscriber(context.Background(), &sub)
	assert.NoError(t, err)

	result, err := server.GetSubscriberData(context.Background(), &id)
	assert.NoError(t, err)
	assert.Equal(t, sub.String(), result.String())

	_, err = server.DeleteSubscriber(context.Background(), &id)
	assert.NoError(t, err)

	_, err = server.GetSubscriberData(context.Background(), &id)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = Subscriber '1' not found")
}

func getConnToTestHomeSubscriberServer(t *testing.T) *grpc.ClientConn {
	srv := test_utils.NewTestHomeSubscriberServer(t)

	s := grpc.NewServer()
	fegprotos.RegisterHSSConfiguratorServer(s, srv)

	lis, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	addr := lis.Addr()
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	assert.NoError(t, err)
	return conn
}

func TestHomeSubscriberServer_GetSubscriberDataGrpc(t *testing.T) {
	conn := getConnToTestHomeSubscriberServer(t)
	defer conn.Close()
	client := fegprotos.NewHSSConfiguratorClient(conn)

	id := protos.SubscriberID{Id: "100"}
	sub := protos.SubscriberData{Sid: &id}

	data, err := client.GetSubscriberData(context.Background(), &id)
	assert.Nil(t, data)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = Subscriber '100' not found")

	reply, err := client.AddSubscriber(context.Background(), &sub)
	assert.Equal(t, orcprotos.Void{}, *reply)
	assert.NoError(t, err)

	data, err = client.GetSubscriberData(context.Background(), &id)
	assert.Equal(t, sub.Sid.Id, data.Sid.Id)
	assert.NoError(t, err)
}
