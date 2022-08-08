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
	"fmt"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/logs_pusher"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/lib/go/merrors"
)

type cbsdManager struct {
	protos.UnimplementedCbsdManagementServer
	store                  storage.CbsdManager
	cbsdInactivityInterval time.Duration
	logConsumerUrl         string
	logPusher              logs_pusher.LogPusher
}

func NewCbsdManager(store storage.CbsdManager, cbsdInactivityInterval time.Duration, logConsumerUrl string, logPusher logs_pusher.LogPusher) protos.CbsdManagementServer {
	return &cbsdManager{
		store:                  store,
		cbsdInactivityInterval: cbsdInactivityInterval,
		logConsumerUrl:         logConsumerUrl,
		logPusher:              logPusher,
	}
}

func (c *cbsdManager) CreateCbsd(_ context.Context, request *protos.CreateCbsdRequest) (*protos.CreateCbsdResponse, error) {
	err := c.store.CreateCbsd(request.NetworkId, cbsdToDatabase(request.Data))
	if err != nil {
		return nil, makeErr(err, "create cbsd")
	}
	return &protos.CreateCbsdResponse{}, nil
}

func (c *cbsdManager) UserUpdateCbsd(_ context.Context, request *protos.UpdateCbsdRequest) (*protos.UpdateCbsdResponse, error) {
	err := c.store.UpdateCbsd(request.NetworkId, request.Id, cbsdToDatabase(request.Data))
	if err != nil {
		return nil, makeErr(err, "update cbsd")
	}
	return &protos.UpdateCbsdResponse{}, nil
}
func (c *cbsdManager) EnodebdUpdateCbsd(ctx context.Context, request *protos.EnodebdUpdateCbsdRequest) (*protos.UpdateCbsdResponse, error) {
	data := requestToDbCbsd(request)
	cbsd, err := c.store.EnodebdUpdateCbsd(data)
	if err != nil {
		return nil, makeErr(err, "update cbsd")
	}
	msg, _ := json.Marshal(request)
	log := &logs_pusher.DPLog{
		EventTimestamp:   clock.Now().Unix(),
		LogFrom:          "CBSD",
		LogTo:            "DP",
		LogName:          "EnodebdUpdateCbsd",
		LogMessage:       string(msg),
		CbsdSerialNumber: cbsd.CbsdSerialNumber.String,
		NetworkId:        cbsd.NetworkId.String,
		FccId:            cbsd.FccId.String,
	}
	if err := c.logPusher(ctx, log, c.logConsumerUrl); err != nil {
		glog.Warningf("Failed to log Enodebd Update. Details: %s", err)
	}

	return &protos.UpdateCbsdResponse{}, nil
}

func requestToDbCbsd(request *protos.EnodebdUpdateCbsdRequest) *storage.DBCbsd {
	cbsd := storage.DBCbsd{
		CbsdSerialNumber: db.MakeString(request.SerialNumber),
		CbsdCategory:     db.MakeString(request.CbsdCategory),
	}
	params := request.GetInstallationParam()
	setInstallationParam(&cbsd, params)
	return &cbsd
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
	filter := dbFilter(request.Filter)
	result, err := c.store.ListCbsd(request.NetworkId, pagination, filter)
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

func (c *cbsdManager) DeregisterCbsd(_ context.Context, request *protos.DeregisterCbsdRequest) (*protos.DeregisterCbsdResponse, error) {
	err := c.store.DeregisterCbsd(request.NetworkId, request.Id)
	if err != nil {
		return nil, makeErr(err, "deregister cbsd")
	}
	return &protos.DeregisterCbsdResponse{}, nil
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

func dbFilter(filter *protos.CbsdFilter) *storage.CbsdFilter {
	p := &storage.CbsdFilter{}
	if filter != nil && filter.SerialNumber != "" {
		p.SerialNumber = filter.GetSerialNumber()
	}
	return p
}

func cbsdToDatabase(data *protos.CbsdData) *storage.MutableCbsd {
	cbsd := buildCbsd(data)
	return &storage.MutableCbsd{
		Cbsd: cbsd,
		DesiredState: &storage.DBCbsdState{
			Name: db.MakeString(data.DesiredState),
		},
	}
}

func buildCbsd(data *protos.CbsdData) *storage.DBCbsd {
	capabilities := data.GetCapabilities()
	preferences := data.GetPreferences()
	installationParam := data.GetInstallationParam()
	b, _ := json.Marshal(preferences.GetFrequenciesMhz())
	cbsd := &storage.DBCbsd{
		UserId:                    db.MakeString(data.GetUserId()),
		FccId:                     db.MakeString(data.GetFccId()),
		CbsdSerialNumber:          db.MakeString(data.GetSerialNumber()),
		MinPower:                  db.MakeFloat(capabilities.GetMinPower()),
		MaxPower:                  db.MakeFloat(capabilities.GetMaxPower()),
		NumberOfPorts:             db.MakeInt(capabilities.GetNumberOfAntennas()),
		PreferredBandwidthMHz:     db.MakeInt(preferences.GetBandwidthMhz()),
		PreferredFrequenciesMHz:   db.MakeString(string(b)),
		SingleStepEnabled:         db.MakeBool(data.GetSingleStepEnabled()),
		CbsdCategory:              db.MakeString(data.GetCbsdCategory()),
		CarrierAggregationEnabled: db.MakeBool(data.GetCarrierAggregationEnabled()),
		GrantRedundancy:           db.MakeBool(data.GetGrantRedundancy()),
		MaxIbwMhx:                 db.MakeInt(data.Capabilities.GetMaxIbwMhz()),
	}
	setInstallationParam(cbsd, installationParam)
	return cbsd
}

func setInstallationParam(cbsd *storage.DBCbsd, params *protos.InstallationParam) {
	if params != nil {
		cbsd.LatitudeDeg = dbFloat64OrNil(params.LatitudeDeg)
		cbsd.LongitudeDeg = dbFloat64OrNil(params.LongitudeDeg)
		cbsd.HeightM = dbFloat64OrNil(params.HeightM)
		cbsd.HeightType = dbStringOrNil(params.HeightType)
		cbsd.IndoorDeployment = dbBoolOrNil(params.IndoorDeployment)
		cbsd.AntennaGain = dbFloat64OrNil(params.AntennaGain)
	}
}

func cbsdFromDatabase(data *storage.DetailedCbsd, inactivityInterval time.Duration) *protos.CbsdDetails {
	isActive := clock.Since(data.Cbsd.LastSeen.Time) < inactivityInterval
	var frequencies []int64
	_ = json.Unmarshal([]byte(data.Cbsd.PreferredFrequenciesMHz.String), &frequencies)
	return &protos.CbsdDetails{
		Id: data.Cbsd.Id.Int64,
		Data: &protos.CbsdData{
			UserId:            data.Cbsd.UserId.String,
			FccId:             data.Cbsd.FccId.String,
			SerialNumber:      data.Cbsd.CbsdSerialNumber.String,
			CbsdCategory:      data.Cbsd.CbsdCategory.String,
			SingleStepEnabled: data.Cbsd.SingleStepEnabled.Bool,
			Capabilities: &protos.Capabilities{
				MinPower:         data.Cbsd.MinPower.Float64,
				MaxPower:         data.Cbsd.MaxPower.Float64,
				NumberOfAntennas: data.Cbsd.NumberOfPorts.Int64,
				MaxIbwMhz:        data.Cbsd.MaxIbwMhx.Int64,
			},
			Preferences: &protos.FrequencyPreferences{
				BandwidthMhz:   data.Cbsd.PreferredBandwidthMHz.Int64,
				FrequenciesMhz: frequencies,
			},
			DesiredState:              data.DesiredState.Name.String,
			InstallationParam:         getInstallationParam(data.Cbsd),
			CarrierAggregationEnabled: data.Cbsd.CarrierAggregationEnabled.Bool,
			GrantRedundancy:           data.Cbsd.GrantRedundancy.Bool,
		},
		CbsdId:   data.Cbsd.CbsdId.String,
		State:    data.CbsdState.Name.String,
		IsActive: isActive,
		Grants:   grantsFromDatabase(data.Grants),
	}
}

func getInstallationParam(c *storage.DBCbsd) *protos.InstallationParam {
	p := &protos.InstallationParam{}
	p.LatitudeDeg = protoDoubleOrNil(c.LatitudeDeg)
	p.LongitudeDeg = protoDoubleOrNil(c.LongitudeDeg)
	p.IndoorDeployment = protoBoolOrNil(c.IndoorDeployment)
	p.HeightM = protoDoubleOrNil(c.HeightM)
	p.HeightType = protoStringOrNil(c.HeightType)
	p.AntennaGain = protoDoubleOrNil(c.AntennaGain)
	return p
}

func grantsFromDatabase(grants []*storage.DetailedGrant) []*protos.GrantDetails {
	const mega int64 = 1e6
	res := make([]*protos.GrantDetails, len(grants))
	for i, g := range grants {
		bw := (g.Grant.HighFrequency.Int64 - g.Grant.LowFrequency.Int64) / mega
		freq := (g.Grant.HighFrequency.Int64 + g.Grant.LowFrequency.Int64) / (mega * 2)
		res[i] = &protos.GrantDetails{
			BandwidthMhz:            bw,
			FrequencyMhz:            freq,
			MaxEirp:                 g.Grant.MaxEirp.Float64,
			State:                   g.GrantState.Name.String,
			TransmitExpireTimestamp: g.Grant.TransmitExpireTime.Time.Unix(),
			GrantExpireTimestamp:    g.Grant.GrantExpireTime.Time.Unix(),
		}
	}
	return res
}

func makeErr(err error, wrap string) error {
	e := fmt.Errorf(wrap+": %w", err)
	code := codes.Internal
	if err == merrors.ErrNotFound {
		code = codes.NotFound
	} else if err == merrors.ErrAlreadyExists {
		code = codes.AlreadyExists
	}
	return status.Error(code, e.Error())
}
