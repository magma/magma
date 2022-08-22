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

package coa

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"fbc/lib/go/radius/dictionaries/ruckus"
	"fbc/lib/go/radius/dictionaries/xwf"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"

	"github.com/stretchr/testify/require"
)

var coaPacketXWFCertified = Params{
	NASIdentifier:       "XWF-C-TLV",
	AcctInterimInterval: 60,
	AcctSessionID:       "00-04-56-91-C4-F8-55-E4-EA-92-87-D1-9A-5A-51-43",
	CallingStationID:    "30:07:4d:9f:33:8b",
	CaptivePortalToken:  "eyJyYWRpdXNfc2Vzc2lvbl92YW5pbGxhX2lkIjoyNTYyNDU4NTQ5MDc4MTJ9",
	TrafficClasses: []AuthorizeTrafficClasses{
		{
			AuthorizeClassName: "xwf",
			AuthorizeBytesLeft: 0,
		},
		{
			AuthorizeClassName: "fbs",
			AuthorizeBytesLeft: 0,
		},
		{
			AuthorizeClassName: "internet",
			AuthorizeBytesLeft: 104109374,
		},
	},
	VendorName: VendorXWFCertified,
}

var coaPacketRuckus = Params{
	NASIdentifier:       "XWF-C-TLV",
	AcctInterimInterval: 60,
	AcctSessionID:       "00-04-56-91-C4-F8-55-E4-EA-92-87-D1-9A-5A-51-43",
	CallingStationID:    "30:07:4d:9f:33:8b",
	CaptivePortalToken:  "eyJyYWRpdXNfc2Vzc2lvbl92YW5pbGxhX2lkIjoyNTYyNDU4NTQ5MDc4MTJ9",
	TrafficClasses: []AuthorizeTrafficClasses{
		{
			AuthorizeClassName: "xwf",
			AuthorizeBytesLeft: 0,
		},
		{
			AuthorizeClassName: "fbs",
			AuthorizeBytesLeft: 0,
		},
		{
			AuthorizeClassName: "internet",
			AuthorizeBytesLeft: 103905273,
		},
	},
	VendorName: VendorRuckus,
}

var disconnectPacket = DisconnectParams{
	NASIdentifier:    "XWF-C-TLV",
	AcctSessionID:    "00-04-56-91-C4-F8-55-E4-EA-92-87-D1-9A-5A-51-43",
	CallingStationID: "30:07:4d:9f:33:8b",
}

const XWFAuthorizeClassInternet = "internet"

func serverCoAHandlerXWFCertified(traffic []xwf.XWFAuthorizeTrafficClasses, t *testing.T) {
	for _, val := range traffic {
		if val.XWFAuthorizeClassName == XWFAuthorizeClassInternet {
			require.NotEqual(t, val.XWFAuthorizeBytesLeft, coaPacketRuckus.TrafficClasses[2].AuthorizeBytesLeft)
			require.Equal(t, val.XWFAuthorizeBytesLeft, coaPacketXWFCertified.TrafficClasses[2].AuthorizeBytesLeft)
			return
		}
	}
	require.Fail(t, "Missing CoA traffic class internt")
}

func serverCoAHandlerRuckus(traffic []ruckus.RuckusTCAttrIdsWithQuota, t *testing.T) {
	for _, val := range traffic {
		if val.RuckusTCNameQuota == XWFAuthorizeClassInternet {
			require.NotEqual(t, val.RuckusTCQuota, coaPacketXWFCertified.TrafficClasses[2].AuthorizeBytesLeft)
			require.Equal(t, val.RuckusTCQuota, coaPacketRuckus.TrafficClasses[2].AuthorizeBytesLeft)
			return
		}
	}
	require.Fail(t, "Missing CoA traffic class internt")
}

func serverCoAHandler(w radius.ResponseWriter, r *radius.Request, t *testing.T) {
	xwfTCAttr, err := xwf.XWFAuthorizeTrafficClasses_Gets(r.Packet)
	require.NoError(t, err, "Failed to get XWF traffic classes")
	ruckesTCAttr, err := ruckus.RuckusTCAttrIdsWithQuota_Gets(r.Packet)
	require.NoError(t, err, "Failed to get Ruckus traffic classes")
	switch {
	case len(xwfTCAttr) > 0:
		require.Equal(t, 1, len(xwfTCAttr), "Number of authorized traffic classes don't fit")
		serverCoAHandlerXWFCertified(xwfTCAttr[0], t)
	case len(ruckesTCAttr) > 0:
		require.Equal(t, 1, len(ruckesTCAttr), "Number of authorized traffic classes don't fit")
		serverCoAHandlerRuckus(ruckesTCAttr[0], t)
	default:
		require.Fail(t, "Unexpected CoA package")
		err := w.Write(r.Response(radius.CodeCoANAK))
		require.NoError(t, err)
		return
	}
	err = w.Write(r.Response(radius.CodeCoAACK))
	require.NoError(t, err)
}

func serverDisconnectHandler(w radius.ResponseWriter, r *radius.Request, t *testing.T) {
	res, err := rfc2865.CallingStationID_LookupString(r.Packet)
	require.NoError(t, err)
	require.Equal(t, res, disconnectPacket.CallingStationID, "Wrong calling station ID")
	err = w.Write(r.Response(radius.CodeDisconnectACK))
	require.NoError(t, err)
}

func serverHandler(w radius.ResponseWriter, r *radius.Request, t *testing.T) {
	switch {
	case r.Packet.Code == radius.CodeCoARequest:
		serverCoAHandler(w, r, t)
	case r.Packet.Code == radius.CodeDisconnectRequest:
		serverDisconnectHandler(w, r, t)
	default:
		require.Fail(t, "Unexpected package", r.Packet.Code)
		err := w.Write(r.Response(radius.CodeCoANAK))
		require.NoError(t, err)
	}
}

func sendPacket(ctx context.Context, t *testing.T, p Packet, client Client, addr string) {
	cod, err := client.Send(ctx, p, addr)
	t.Logf("get code: %d", cod)
	require.NoError(t, err)
	require.Equal(t, cod, CodeACK, "expecting ack")
}

func TestCoA(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:1812")
	require.NoError(t, err)
	pc, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)

	secret := []byte("123456790")

	firstCoAPacket, err := CreateCoARequest(coaPacketXWFCertified, secret)
	require.NoError(t, err)
	secondCoAPacket, err := CreateCoARequest(coaPacketRuckus, secret)
	require.NoError(t, err)
	disconnectPacket, err := CreateCoADisconnect(disconnectPacket, secret)
	require.NoError(t, err)

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second)
	defer ctxCancel()

	client := CreateClient(time.Millisecond*50, 0)

	server := radius.PacketServer{
		SecretSource: radius.StaticSecretSource(secret),
		Handler: radius.HandlerFunc(func(w radius.ResponseWriter, r *radius.Request) {
			serverHandler(w, r, t)
		}),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer server.Shutdown(context.Background())
		sendPacket(ctx, t, firstCoAPacket, client, pc.LocalAddr().String())
		sendPacket(ctx, t, secondCoAPacket, client, pc.LocalAddr().String())
		sendPacket(ctx, t, disconnectPacket, client, pc.LocalAddr().String())
		t.Log("Done")
	}()

	err = server.Serve(pc)
	require.NoError(t, err)
	wg.Wait()
}
