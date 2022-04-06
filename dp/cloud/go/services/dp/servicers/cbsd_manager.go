/*
Copyright 2022 The Magma Authors.

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
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/lib/go/merrors"
)

type cbsdManager struct {
	protos.UnimplementedCbsdManagementServer
	store                  storage.CbsdManager
	cbsdInactivityInterval time.Duration
}

func NewCbsdManager(store storage.CbsdManager, cbsdInactivityInterval time.Duration) protos.CbsdManagementServer {
	return &cbsdManager{
		store:                  store,
		cbsdInactivityInterval: cbsdInactivityInterval,
	}
}

func (c *cbsdManager) CreateCbsd(_ context.Context, request *protos.CreateCbsdRequest) (*protos.CreateCbsdResponse, error) {
	err := c.store.CreateCbsd(request.NetworkId, cbsdToDatabase(request.Data))
	if err != nil {
		return nil, makeErr(err, "create cbsd")
	}
	return &protos.CreateCbsdResponse{}, nil
}

func (c *cbsdManager) UpdateCbsd(_ context.Context, request *protos.UpdateCbsdRequest) (*protos.UpdateCbsdResponse, error) {
	err := c.store.UpdateCbsd(request.NetworkId, request.Id, cbsdToDatabase(request.Data))
	if err != nil {
		return nil, makeErr(err, "update cbsd")
	}
	return &protos.UpdateCbsdResponse{}, nil
}

func (c *cbsdManager) DeleteCbsd(_ context.Context, request *protos.DeleteCbsdRequest) (*protos.DeleteCbsdResponse, error) {
	err := c.store.DeleteCbsd(request.NetworkId, request.Id)
	if err != nil {
		return nil, makeErr(err, "delete cbsd")
	}
	return &protos.DeleteCbsdResponse{}, nil
}

func (c *cbsdManager) FetchCbsd(_ context.Context, request *protos.FetchCbsdRequest) (*protos.FetchCbsdResponse, error) {
	result, err := c.store.FetchCbsd(request.NetworkId, request.Id)
	if err != nil {
		return nil, makeErr(err, "fetch cbsd")
	}
	details := cbsdFromDatabase(result, c.cbsdInactivityInterval)
	return &protos.FetchCbsdResponse{Details: details}, nil
}

func (c *cbsdManager) ListCbsds(_ context.Context, request *protos.ListCbsdRequest) (*protos.ListCbsdResponse, error) {
	pagination := dbPagination(request.Pagination)
	result, err := c.store.ListCbsd(request.NetworkId, pagination)
	if err != nil {
		return nil, makeErr(err, "list cbsds")
	}
	resp := &protos.ListCbsdResponse{
		Details:    make([]*protos.CbsdDetails, len(result.Cbsds)),
		TotalCount: result.Count,
	}
	for i, data := range result.Cbsds {
		resp.Details[i] = cbsdFromDatabase(data, c.cbsdInactivityInterval)
	}
	return resp, nil
}

func dbPagination(pagination *protos.Pagination) *storage.Pagination {
	p := &storage.Pagination{}
	if pagination.Limit != nil {
		p.Limit = db.MakeInt(pagination.GetLimit().Value)
	}
	if pagination.Offset != nil {
		p.Offset = db.MakeInt(pagination.GetOffset().Value)
	}
	return p
}

func cbsdToDatabase(data *protos.CbsdData) *storage.DBCbsd {
	capabilities := data.Capabilities
	preferences := data.Preferences
	b, _ := json.Marshal(preferences.FrequenciesMhz)
	return &storage.DBCbsd{
		UserId:                  db.MakeString(data.UserId),
		FccId:                   db.MakeString(data.FccId),
		CbsdSerialNumber:        db.MakeString(data.SerialNumber),
		MinPower:                db.MakeFloat(capabilities.MinPower),
		MaxPower:                db.MakeFloat(capabilities.MaxPower),
		AntennaGain:             db.MakeFloat(capabilities.AntennaGain),
		NumberOfPorts:           db.MakeInt(capabilities.NumberOfAntennas),
		PreferredBandwidthMHz:   db.MakeInt(preferences.BandwidthMhz),
		PreferredFrequenciesMHz: db.MakeString(string(b)),
	}
}

func cbsdFromDatabase(data *storage.DetailedCbsd, inactivityInterval time.Duration) *protos.CbsdDetails {
	const mega int64 = 1e6
	var grant *protos.GrantDetails
	if data.GrantState.Name.Valid {
		bandwidth := (data.Grant.HighFrequency.Int64 - data.Grant.LowFrequency.Int64) / mega
		frequency := (data.Grant.HighFrequency.Int64 + data.Grant.LowFrequency.Int64) / (mega * 2)
		grant = &protos.GrantDetails{
			BandwidthMhz:            bandwidth,
			FrequencyMhz:            frequency,
			MaxEirp:                 data.Grant.MaxEirp.Float64,
			State:                   data.GrantState.Name.String,
			TransmitExpireTimestamp: data.Grant.TransmitExpireTime.Time.Unix(),
			GrantExpireTimestamp:    data.Grant.GrantExpireTime.Time.Unix(),
		}
	}
	isActive := clock.Since(data.Cbsd.LastSeen.Time) < inactivityInterval
	var frequencies []int64
	_ = json.Unmarshal([]byte(data.Cbsd.PreferredFrequenciesMHz.String), &frequencies)
	return &protos.CbsdDetails{
		Id: data.Cbsd.Id.Int64,
		Data: &protos.CbsdData{
			UserId:       data.Cbsd.UserId.String,
			FccId:        data.Cbsd.FccId.String,
			SerialNumber: data.Cbsd.CbsdSerialNumber.String,
			Capabilities: &protos.Capabilities{
				MinPower:         data.Cbsd.MinPower.Float64,
				MaxPower:         data.Cbsd.MaxPower.Float64,
				NumberOfAntennas: data.Cbsd.NumberOfPorts.Int64,
				AntennaGain:      data.Cbsd.AntennaGain.Float64,
			},
			Preferences: &protos.FrequencyPreferences{
				BandwidthMhz:   data.Cbsd.PreferredBandwidthMHz.Int64,
				FrequenciesMhz: frequencies,
			},
		},
		CbsdId:   data.Cbsd.CbsdId.String,
		State:    data.CbsdState.Name.String,
		IsActive: isActive,
		Grant:    grant,
	}
}

func makeErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	code := codes.Internal
	if err == merrors.ErrNotFound {
		code = codes.NotFound
	}
	return status.Error(code, e.Error())
}
