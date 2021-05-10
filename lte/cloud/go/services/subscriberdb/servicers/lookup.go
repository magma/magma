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

package servicers

import (
	"context"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	subscriberdb_storage "magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// lookupServicer translates subscriber aliases to their IMSI.
//
// MSISDN is stored as a blobstore table, MSISDN -> IMSI.
//
// IP is stored as a SQL table with two string columns:
//	- IP
//	- IMSI.APN
// Each subscriber should have at most one IP per APN, so the IP+IMSI+APN
// triplet is enforced to be unique.
type lookupServicer struct {
	factory blobstore.BlobStorageFactory
	store   subscriberdb_storage.IPLookup
}

// NewLookupServicer returns a new subscriber lookup servicer.
// Stores should be initialized by the caller.
func NewLookupServicer(msisdnFact blobstore.BlobStorageFactory, ipStore subscriberdb_storage.IPLookup) protos.SubscriberLookupServer {
	return &lookupServicer{factory: msisdnFact, store: ipStore}
}

func (l *lookupServicer) GetMSISDNs(ctx context.Context, req *protos.GetMSISDNsRequest) (*protos.GetMSISDNsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	store, err := l.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "error starting transaction: %v", err)
	}
	defer store.Rollback()

	tks := storage.MakeTKs(lte.MSISDNBlobstoreType, req.Msisdns)
	var blobs blobstore.Blobs
	if len(tks) == 0 {
		blobs, err = blobstore.GetAllOfType(store, req.NetworkId, lte.MSISDNBlobstoreType)
		if err != nil {
			return nil, makeErr(err, "get msisdns from blobstore")
		}
	} else {
		blobs, err = store.GetMany(req.NetworkId, tks)
		if err != nil {
			return nil, makeErr(err, "get msisdns from blobstore")
		}
	}

	imsisByMSISDN := map[string]string{}
	for _, blob := range blobs {
		imsisByMSISDN[blob.Key] = string(blob.Value)
	}

	res := &protos.GetMSISDNsResponse{ImsisByMsisdn: imsisByMSISDN}
	return res, store.Commit()
}

func (l *lookupServicer) SetMSISDN(ctx context.Context, req *protos.SetMSISDNRequest) (*protos.SetMSISDNResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	store, err := l.factory.StartTransaction(nil)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "error starting transaction: %v", err)
	}
	defer store.Rollback()

	// Ensure mapping doesn't exist
	blob, err := store.Get(req.NetworkId, storage.TypeAndKey{Type: lte.MSISDNBlobstoreType, Key: req.Msisdn})
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "msisdn already mapped to %s", blob.Value)
	}
	if err != merrors.ErrNotFound {
		return nil, makeErr(err, "get msisdn from blobstore")
	}

	err = store.CreateOrUpdate(req.NetworkId, blobstore.Blobs{{
		Type:  lte.MSISDNBlobstoreType,
		Key:   req.Msisdn,
		Value: []byte(req.Imsi),
	}})
	if err != nil {
		return nil, makeErr(err, "create msisdn mapping in blobstore")
	}

	return &protos.SetMSISDNResponse{}, store.Commit()
}

func (l *lookupServicer) DeleteMSISDN(ctx context.Context, req *protos.DeleteMSISDNRequest) (*protos.DeleteMSISDNResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	store, err := l.factory.StartTransaction(nil)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "error starting transaction: %v", err)
	}
	defer store.Rollback()

	err = store.Delete(req.NetworkId, []storage.TypeAndKey{{
		Type: lte.MSISDNBlobstoreType,
		Key:  req.Msisdn,
	}})
	if err != nil {
		return nil, makeErr(err, "delete msisdn from blobstore")
	}

	return &protos.DeleteMSISDNResponse{}, store.Commit()
}

func (l *lookupServicer) GetIPs(ctx context.Context, req *protos.GetIPsRequest) (*protos.GetIPsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ipMappings, err := l.store.GetIPs(req.NetworkId, req.Ips)
	if err != nil {
		return nil, makeErr(err, "get IPs from store")
	}

	res := &protos.GetIPsResponse{IpMappings: ipMappings}
	return res, nil
}

func (l *lookupServicer) SetIPs(ctx context.Context, req *protos.SetIPsRequest) (*protos.SetIPsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := l.store.SetIPs(req.NetworkId, req.IpMappings)
	if err != nil {
		return nil, makeErr(err, "set IP mappings to store")
	}
	return &protos.SetIPsResponse{}, nil
}

func makeErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	code := codes.Internal
	if err == merrors.ErrNotFound {
		code = codes.NotFound
	}
	return status.Error(code, e.Error())
}
