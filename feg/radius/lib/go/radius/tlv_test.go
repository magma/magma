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
	"reflect"
	"testing"
	"time"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/testdata/testtlv"
)

var expectedResults = [][]testtlv.TestStruct{
	{
		{
			Field1: "First",
			Field2: 1,
		},
		{
			Field1: "First",
			Field2: 2,
		},
	},
	{
		{
			Field1: "Second",
			Field2: 1,
		},
		{
			Field1: "Second",
			Field2: 2,
		},
	},
}

func testStructToString(st []testtlv.TestStruct) (res string) {
	for i, _st := range st {
		res += fmt.Sprintf("%d:\n TestStruct, field1: %s field2: %d\n", i, _st.Field1, _st.Field2)
	}
	return
}

func testStructsPointerToString(st [][]testtlv.TestStruct) (res string) {
	for _, _st := range st {
		res += testStructToString(_st)
	}
	return
}

func serverHandler(w radius.ResponseWriter, r *radius.Request, t *testing.T) error {
	lookupRes, lookupErr := testtlv.TestStruct_Lookup(r.Packet)
	if lookupErr != nil {
		t.Log("Lookup error", lookupErr)
		return lookupErr
	}
	if !reflect.DeepEqual(lookupRes, expectedResults[0]) {
		return fmt.Errorf("\ngot:\n%s\nexpected:\n%s", testStructToString(lookupRes), testStructToString(expectedResults[0]))
	}

	getsRes, getsErr := testtlv.TestStruct_Gets(r.Packet)
	if getsErr != nil {
		return getsErr
	}
	if !reflect.DeepEqual(getsRes, expectedResults) {
		return fmt.Errorf("\ngot:\n%s\nexpected:\n%s", testStructsPointerToString(getsRes), testStructsPointerToString(expectedResults))
	}

	return nil
}

func TestPacketTLV(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:1812")
	if err != nil {
		t.Fatal(err)
	}
	pc, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}

	secret := []byte("123456790")

	var serverErr, clientErr error

	server := radius.PacketServer{
		SecretSource: radius.StaticSecretSource(secret),
		Handler: radius.HandlerFunc(func(w radius.ResponseWriter, r *radius.Request) {
			serverErr = serverHandler(w, r, t)
			w.Write(r.Response(radius.CodeAccessAccept))
		}),
	}

	go func() {
		defer server.Shutdown(context.Background())
		packet := radius.New(radius.CodeAccessRequest, secret)
		testtlv.TestStruct_Add(packet, expectedResults[0])
		testtlv.TestStruct_Add(packet, expectedResults[1])
		client := radius.Client{
			Retry: time.Millisecond * 50,
		}
		response, err := client.Exchange(context.Background(), packet, pc.LocalAddr().String())
		if err != nil {
			t.Log("Done")
			clientErr = err
			return
		}
		if response.Code != radius.CodeAccessAccept {
			clientErr = fmt.Errorf("expected CodeAccessAccept, got %s\n", response.Code)
		}
		t.Log("Done")
	}()

	if err := server.Serve(pc); err != nil {
		t.Fatal(err)
	}

	server.Shutdown(context.Background())

	if clientErr != nil {
		t.Fatal(err)
	}

	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
