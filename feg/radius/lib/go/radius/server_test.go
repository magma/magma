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

package radius_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"fbc/lib/go/radius"
	. "fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2869"
)

type VerifyResponseFunc func(response *radius.Packet) error

func TestPacketServer_basic(t *testing.T) {
	handle := radius.HandlerFunc(func(w radius.ResponseWriter, r *radius.Request) {
		username := UserName_GetString(r.Packet)
		if username == "tim" {
			w.Write(r.Response(radius.CodeAccessAccept))
		} else {
			w.Write(r.Response(radius.CodeAccessReject))
		}
	})

	verify := VerifyResponseFunc(func(response *radius.Packet) error {
		if response.Code != radius.CodeAccessAccept {
			return fmt.Errorf("expected CodeAccessAccept, got %s", response.Code)
		}
		return nil
	})

	RunTestServer(t, handle, verify)
}

func TestPacketServer_msgauth(t *testing.T) {

	verify := VerifyResponseFunc(func(response *radius.Packet) error {
		if response.Code != radius.CodeAccessAccept {
			return fmt.Errorf("expected CodeAccessAccept, got %s", response.Code)
		}
		msgAuth := response.Get(radius.Type(80))
		if msgAuth == nil {
			return fmt.Errorf("Message Authenticator was not generated")
		}
		return nil
	})

	handle := radius.HandlerFunc(func(w radius.ResponseWriter, r *radius.Request) {
		res := r.Response(radius.CodeAccessAccept)
		res.Attributes.Add(radius.Type(rfc2869.EAPMessage_Type), []byte{0x1, 0x2})
		w.Write(res)
	})

	RunTestServer(t, handle, verify)
}

func RunTestServer(t *testing.T, handler radius.HandlerFunc, verify VerifyResponseFunc) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	pc, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}

	secret := []byte("123456790")

	server := radius.PacketServer{
		SecretSource: radius.StaticSecretSource(secret),
		Handler:      handler,
	}

	var clientErr error
	go func() {
		defer server.Shutdown(context.Background())

		packet := radius.New(radius.CodeAccessRequest, secret)
		UserName_SetString(packet, "tim")
		client := radius.Client{
			Retry: time.Millisecond * 50,
		}
		response, err := client.Exchange(context.Background(), packet, pc.LocalAddr().String())
		if err != nil {
			clientErr = err
			return
		}
		clientErr = verify(response)
	}()

	if err := server.Serve(pc); err != nil {
		t.Fatal(err)
	}

	server.Shutdown(context.Background())

	if clientErr != nil {
		t.Fatal(clientErr)
	}
}
