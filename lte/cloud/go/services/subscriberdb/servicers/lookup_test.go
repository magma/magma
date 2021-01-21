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
	"testing"

	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/lte/cloud/go/services/subscriberdb/servicers"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func TestLookupServicer_MSISDNs(t *testing.T) {
	ctx := context.Background()
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory(subscriberdb.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	l := servicers.NewLookupServicer(fact, nil)

	t.Run("initially empty", func(t *testing.T) {
		got, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   []string{"msisdn0"},
		})
		assert.NoError(t, err)
		assert.Empty(t, got.ImsisByMsisdn)

		gotAll0, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   []string{},
		})
		assert.NoError(t, err)
		assert.Empty(t, gotAll0.ImsisByMsisdn)

		gotAll1, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   nil,
		})
		assert.NoError(t, err)
		assert.Empty(t, gotAll1.ImsisByMsisdn)
	})

	t.Run("add msisdn", func(t *testing.T) {
		_, err := l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "msisdn0",
			Imsi:      "imsi0",
		})
		assert.NoError(t, err)

		_, err = l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "msisdn1",
			Imsi:      "imsi1",
		})
		assert.NoError(t, err)

		_, err = l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "nid1",
			Msisdn:    "msisdn2",
			Imsi:      "imsi2",
		})
		assert.NoError(t, err)

		got0, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   []string{"msisdn0"},
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"msisdn0": "imsi0"}, got0.ImsisByMsisdn)

		got1, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   []string{"msisdn1"},
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"msisdn1": "imsi1"}, got1.ImsisByMsisdn)

		gotAll, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"msisdn0": "imsi0", "msisdn1": "imsi1"}, gotAll.ImsisByMsisdn)
	})

	t.Run("validate requests", func(t *testing.T) {
		// Can't overwrite existing mapping
		_, err := l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "msisdn0",
			Imsi:      "imsiXXX",
		})
		assert.Error(t, err)

		// Empty network ID
		_, err = l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "",
			Msisdn:    "msisdnA",
			Imsi:      "imsiA",
		})
		assert.Error(t, err)

		// Empty MSISDN
		_, err = l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "",
			Imsi:      "imsiB",
		})
		assert.Error(t, err)

		// Empty IMSI
		_, err = l.SetMSISDN(ctx, &protos.SetMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "msisdnC",
			Imsi:      "",
		})
		assert.Error(t, err)

		// Empty network ID
		_, err = l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "",
			Msisdns:   nil,
		})
		assert.Error(t, err)

		// Empty network ID
		_, err = l.DeleteMSISDN(ctx, &protos.DeleteMSISDNRequest{
			NetworkId: "",
			Msisdn:    "msisdn0",
		})
		assert.Error(t, err)

		// Empty MSISDN
		_, err = l.DeleteMSISDN(ctx, &protos.DeleteMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "",
		})
		assert.Error(t, err)

		// All mutations should have failed
		gotAll, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"msisdn0": "imsi0", "msisdn1": "imsi1"}, gotAll.ImsisByMsisdn)
	})

	t.Run("delete msisdns", func(t *testing.T) {
		// Empty network ID
		_, err := l.DeleteMSISDN(ctx, &protos.DeleteMSISDNRequest{
			NetworkId: "",
			Msisdn:    "msisdn0",
		})
		assert.Error(t, err)

		_, err = l.DeleteMSISDN(ctx, &protos.DeleteMSISDNRequest{
			NetworkId: "nid0",
			Msisdn:    "msisdn1",
		})
		assert.NoError(t, err)

		gotAll0, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid0",
			Msisdns:   nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"msisdn0": "imsi0"}, gotAll0.ImsisByMsisdn)

		gotAll1, err := l.GetMSISDNs(ctx, &protos.GetMSISDNsRequest{
			NetworkId: "nid1",
			Msisdns:   nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"msisdn2": "imsi2"}, gotAll1.ImsisByMsisdn)
	})
}

func TestLookupServicer_IPs(t *testing.T) {
	ctx := context.Background()
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	store := storage.NewIPLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, err)
	err = store.Initialize()
	assert.NoError(t, err)
	l := servicers.NewLookupServicer(nil, store)

	t.Run("basic", func(t *testing.T) {
		got, err := l.GetIPs(ctx, &protos.GetIPsRequest{
			NetworkId: "nid0",
			Ips:       nil,
		})
		assert.NoError(t, err)
		assert.Empty(t, got.IpMappings)

		want := []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipB", Imsi: "IMSI0", Apn: "apn1"},
			{Ip: "ipC", Imsi: "IMSI1", Apn: "apn0"},
		}
		_, err = l.SetIPs(ctx, &protos.SetIPsRequest{
			NetworkId:  "nid0",
			IpMappings: want,
		})
		assert.NoError(t, err)

		got, err = l.GetIPs(ctx, &protos.GetIPsRequest{
			NetworkId: "nid0",
			Ips:       []string{"ipA", "ipB", "ipC"},
		})
		assert.NoError(t, err)
		assert.Equal(t, want, got.IpMappings)
	})

	t.Run("validate requests", func(t *testing.T) {
		// Empty network ID
		_, err = l.SetIPs(ctx, &protos.SetIPsRequest{
			NetworkId:  "",
			IpMappings: []*protos.IPMapping{{Ip: "ipA", Imsi: "imsi0", Apn: "apn0"}},
		})
		assert.Error(t, err)

		// Empty IP
		_, err = l.SetIPs(ctx, &protos.SetIPsRequest{
			NetworkId:  "nid0",
			IpMappings: []*protos.IPMapping{{Ip: "", Imsi: "imsi0", Apn: "apn0"}},
		})
		assert.Error(t, err)

		// Empty IMSI
		_, err = l.SetIPs(ctx, &protos.SetIPsRequest{
			NetworkId:  "nid0",
			IpMappings: []*protos.IPMapping{{Ip: "ipA", Imsi: "", Apn: "apn0"}},
		})
		assert.Error(t, err)

		// Empty APN
		_, err = l.SetIPs(ctx, &protos.SetIPsRequest{
			NetworkId:  "nid0",
			IpMappings: []*protos.IPMapping{{Ip: "ipA", Imsi: "imsi0", Apn: ""}},
		})
		assert.Error(t, err)

		// Empty network ID
		_, err = l.GetIPs(ctx, &protos.GetIPsRequest{
			NetworkId: "",
			Ips:       nil,
		})
		assert.Error(t, err)

		// All mutations should have failed
		got, err := l.GetIPs(ctx, &protos.GetIPsRequest{
			NetworkId: "nid0",
			Ips:       []string{"ipA", "ipB", "ipC"},
		})
		assert.NoError(t, err)
		want := []*protos.IPMapping{
			{Ip: "ipA", Imsi: "IMSI0", Apn: "apn0"},
			{Ip: "ipB", Imsi: "IMSI0", Apn: "apn1"},
			{Ip: "ipC", Imsi: "IMSI1", Apn: "apn0"},
		}
		assert.Equal(t, want, got.IpMappings)
	})
}
