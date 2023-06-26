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

package builders

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/logs_pusher"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/clock"
)

const (
	cbsdId           int64 = 123
	someSerialNumber       = "some_serial_number"
	someCbsdId             = "some_cbsd_id"
	registered             = "registered"
	authorized             = "authorized"
	someFccId              = "some_fcc_id"
	someUserId             = "some_user_id"
	catB                   = "b"
	catA                   = "a"
)

type DBCbsdBuilder struct {
	Cbsd *storage.DBCbsd
}

func NewDBCbsdBuilder() *DBCbsdBuilder {
	return &DBCbsdBuilder{
		Cbsd: &storage.DBCbsd{
			UserId:                    db.MakeString(someUserId),
			FccId:                     db.MakeString(someFccId),
			CbsdSerialNumber:          db.MakeString(someSerialNumber),
			PreferredBandwidthMHz:     db.MakeInt(20),
			PreferredFrequenciesMHz:   []int64{3600},
			MinPower:                  db.MakeFloat(10),
			MaxPower:                  db.MakeFloat(20),
			NumberOfPorts:             db.MakeInt(2),
			CbsdCategory:              db.MakeString(catB),
			SingleStepEnabled:         db.MakeBool(false),
			CarrierAggregationEnabled: db.MakeBool(false),
			GrantRedundancy:           db.MakeBool(true),
			MaxIbwMhx:                 db.MakeInt(150),
		},
	}
}

func (b *DBCbsdBuilder) Empty() *DBCbsdBuilder {
	b.Cbsd = &storage.DBCbsd{}
	return b
}

func (b *DBCbsdBuilder) WithId(id int64) *DBCbsdBuilder {
	b.Cbsd.Id = db.MakeInt(id)
	return b
}

func (b *DBCbsdBuilder) WithIsDeleted(deleted bool) *DBCbsdBuilder {
	b.Cbsd.IsDeleted = db.MakeBool(deleted)
	return b
}

func (b *DBCbsdBuilder) WithCbsdId(id string) *DBCbsdBuilder {
	b.Cbsd.CbsdId = db.MakeString(id)
	return b
}

func (b *DBCbsdBuilder) WithNetworkId(id string) *DBCbsdBuilder {
	b.Cbsd.NetworkId = db.MakeString(id)
	return b
}

func (b *DBCbsdBuilder) WithFccId(id string) *DBCbsdBuilder {
	b.Cbsd.FccId = db.MakeString(id)
	return b
}

func (b *DBCbsdBuilder) WithUserId(id string) *DBCbsdBuilder {
	b.Cbsd.UserId = db.MakeString(id)
	return b
}

func (b *DBCbsdBuilder) WithAntennaGain(gain float64) *DBCbsdBuilder {
	b.Cbsd.AntennaGainDbi = db.MakeFloat(gain)
	return b
}

func (b *DBCbsdBuilder) WithLatitude(lat float64) *DBCbsdBuilder {
	b.Cbsd.LatitudeDeg = db.MakeFloat(lat)
	return b
}

func (b *DBCbsdBuilder) WithEirpCapabilities(minPower float64, maxPower float64, numberOfPorts int64) *DBCbsdBuilder {
	b.Cbsd.MinPower = db.MakeFloat(minPower)
	b.Cbsd.MaxPower = db.MakeFloat(maxPower)
	b.Cbsd.NumberOfPorts = db.MakeInt(numberOfPorts)
	return b
}

func (b *DBCbsdBuilder) WithLongitude(lon float64) *DBCbsdBuilder {
	b.Cbsd.LongitudeDeg = db.MakeFloat(lon)
	return b
}

func (b *DBCbsdBuilder) WithLastSeen(t int64) *DBCbsdBuilder {
	b.Cbsd.LastSeen = db.MakeTime(time.Unix(t, 0).UTC())
	return b
}

func (b *DBCbsdBuilder) WithStateId(t int64) *DBCbsdBuilder {
	b.Cbsd.StateId = db.MakeInt(t)
	return b
}

func (b *DBCbsdBuilder) WithDesiredStateId(t int64) *DBCbsdBuilder {
	b.Cbsd.DesiredStateId = db.MakeInt(t)
	return b
}

func (b *DBCbsdBuilder) WithChannels(channels []storage.Channel) *DBCbsdBuilder {
	b.Cbsd.Channels = channels
	return b
}

func (b *DBCbsdBuilder) WithPreferences(bandwidthMhz int64, frequenciesMhz []int64) *DBCbsdBuilder {
	b.Cbsd.PreferredBandwidthMHz = db.MakeInt(bandwidthMhz)
	b.Cbsd.PreferredFrequenciesMHz = frequenciesMhz
	return b
}

func (b *DBCbsdBuilder) WithAvailableFrequencies(frequenciesMhz []uint32) *DBCbsdBuilder {
	b.Cbsd.AvailableFrequencies = frequenciesMhz
	return b
}

func (b *DBCbsdBuilder) WithSerialNumber(serial string) *DBCbsdBuilder {
	b.Cbsd.CbsdSerialNumber = db.MakeString(serial)
	return b
}

func (b *DBCbsdBuilder) WithFullInstallationParam() *DBCbsdBuilder {
	b.Cbsd.LatitudeDeg = db.MakeFloat(10.5)
	b.Cbsd.LongitudeDeg = db.MakeFloat(11.5)
	b.Cbsd.IndoorDeployment = db.MakeBool(true)
	b.Cbsd.HeightM = db.MakeFloat(12.5)
	b.Cbsd.HeightType = db.MakeString("agl")
	b.Cbsd.AntennaGainDbi = db.MakeFloat(4.5)
	return b
}

func (b *DBCbsdBuilder) WithFullEnodebdAllowedInstallationParam() *DBCbsdBuilder {
	b.Cbsd.LatitudeDeg = db.MakeFloat(10.5)
	b.Cbsd.LongitudeDeg = db.MakeFloat(11.5)
	b.Cbsd.IndoorDeployment = db.MakeBool(true)
	b.Cbsd.HeightM = db.MakeFloat(12.5)
	b.Cbsd.HeightType = db.MakeString("agl")
	return b
}

func (b *DBCbsdBuilder) WithIncompleteInstallationParam() *DBCbsdBuilder {
	b.Cbsd.LatitudeDeg = db.MakeFloat(10.5)
	b.Cbsd.LongitudeDeg = db.MakeFloat(11.5)
	b.Cbsd.IndoorDeployment = db.MakeBool(true)
	return b
}

func (b *DBCbsdBuilder) WithIndoorDeployment(indoor bool) *DBCbsdBuilder {
	b.Cbsd.IndoorDeployment = db.MakeBool(indoor)
	return b
}

func (b *DBCbsdBuilder) WithSingleStepEnabled(enabled bool) *DBCbsdBuilder {
	b.Cbsd.SingleStepEnabled = db.MakeBool(enabled)
	return b
}

func (b *DBCbsdBuilder) WithCarrierAggregationEnabled(enabled bool) *DBCbsdBuilder {
	b.Cbsd.CarrierAggregationEnabled = db.MakeBool(enabled)
	return b
}

func (b *DBCbsdBuilder) WithGrantRedundancy(enabled bool) *DBCbsdBuilder {
	b.Cbsd.GrantRedundancy = db.MakeBool(enabled)
	return b
}

func (b *DBCbsdBuilder) WithMaxIbwMhx(ibw int64) *DBCbsdBuilder {
	b.Cbsd.MaxIbwMhx = db.MakeInt(ibw)
	return b
}

func (b *DBCbsdBuilder) WithShouldDeregister(should bool) *DBCbsdBuilder {
	b.Cbsd.ShouldDeregister = db.MakeBool(should)
	return b
}

func (b *DBCbsdBuilder) WithShouldRelinquish(should bool) *DBCbsdBuilder {
	b.Cbsd.ShouldRelinquish = db.MakeBool(should)
	return b
}

func (b *DBCbsdBuilder) WithCbsdCategory(cat string) *DBCbsdBuilder {
	b.Cbsd.CbsdCategory = db.MakeString(cat)
	return b
}

type DBGrantBuilder struct {
	Grant *storage.DBGrant
}

func NewDBGrantBuilder() *DBGrantBuilder {
	return &DBGrantBuilder{
		Grant: &storage.DBGrant{},
	}
}

func (b *DBGrantBuilder) WithDefaultTestValues() *DBGrantBuilder {
	b.Grant = &storage.DBGrant{
		GrantExpireTime:      db.MakeTime(time.Unix(123, 0).UTC()),
		TransmitExpireTime:   db.MakeTime(time.Unix(456, 0).UTC()),
		LowFrequencyHz:       db.MakeInt(3590 * 1e6),
		HighFrequencyHz:      db.MakeInt(3610 * 1e6),
		MaxEirp:              db.MakeFloat(35),
		GrantId:              db.MakeString("some_grant_id"),
		HeartbeatIntervalSec: db.MakeInt(1),
	}
	return b
}

func (b *DBGrantBuilder) WithId(id int64) *DBGrantBuilder {
	b.Grant.Id = db.MakeInt(id)
	return b
}

func (b *DBGrantBuilder) WithMaxEirp(e float64) *DBGrantBuilder {
	b.Grant.MaxEirp = db.MakeFloat(e)
	return b
}

func (b *DBGrantBuilder) WithCbsdId(id int64) *DBGrantBuilder {
	b.Grant.CbsdId = db.MakeInt(id)
	return b
}

func (b *DBGrantBuilder) WithStateId(id int64) *DBGrantBuilder {
	b.Grant.StateId = db.MakeInt(id)
	return b
}

func (b *DBGrantBuilder) WithGrantId(id string) *DBGrantBuilder {
	b.Grant.GrantId = db.MakeString(id)
	return b
}

func (b *DBGrantBuilder) WithFrequency(frequencyMHz int64) *DBGrantBuilder {
	b.Grant.LowFrequencyHz = db.MakeInt((frequencyMHz - 10) * 1e6)
	b.Grant.HighFrequencyHz = db.MakeInt((frequencyMHz + 10) * 1e6)
	return b
}

func (b *DBGrantBuilder) WithGrantExpireTime(t time.Time) *DBGrantBuilder {
	b.Grant.GrantExpireTime = db.MakeTime(t)
	return b
}

func (b *DBGrantBuilder) WithTransmitExpireTime(t time.Time) *DBGrantBuilder {
	b.Grant.TransmitExpireTime = db.MakeTime(t)
	return b
}

func (b *DBGrantBuilder) WithLastHeartbeatTime(t time.Time) *DBGrantBuilder {
	b.Grant.LastHeartbeatRequestTime = db.MakeTime(t)
	return b
}

type CbsdStateResultBuilder struct {
	Result *protos.CBSDStateResult
}

func NewCbsdStateResultBuilder(radioEnabled bool, carrierAggregation bool) *CbsdStateResultBuilder {
	return &CbsdStateResultBuilder{
		Result: &protos.CBSDStateResult{
			RadioEnabled:              radioEnabled,
			CarrierAggregationEnabled: carrierAggregation,
		},
	}
}

func (b *CbsdStateResultBuilder) WithChannels(c []*protos.LteChannel) *CbsdStateResultBuilder {
	b.Result.Channels = c
	b.Result.Channel = c[0]
	return b
}

type CbsdProtoPayloadBuilder struct {
	Payload *protos.CbsdData
}

func NewCbsdProtoPayloadBuilder() *CbsdProtoPayloadBuilder {
	return &CbsdProtoPayloadBuilder{
		Payload: &protos.CbsdData{
			UserId:       someUserId,
			FccId:        someFccId,
			SerialNumber: someSerialNumber,
			Preferences: &protos.FrequencyPreferences{
				BandwidthMhz:   20,
				FrequenciesMhz: []int64{3600},
			},
			Capabilities: &protos.Capabilities{
				MinPower:         10,
				MaxPower:         20,
				NumberOfAntennas: 2,
				MaxIbwMhz:        150,
			},
			DesiredState:              registered,
			CbsdCategory:              catB,
			SingleStepEnabled:         false,
			InstallationParam:         &protos.InstallationParam{},
			CarrierAggregationEnabled: false,
			GrantRedundancy:           true,
		},
	}
}

func (b *CbsdProtoPayloadBuilder) Empty() *CbsdProtoPayloadBuilder {
	b.Payload = &protos.CbsdData{}
	return b
}

func (b *CbsdProtoPayloadBuilder) WithEmptyPreferences() *CbsdProtoPayloadBuilder {
	b.Payload.Preferences = &protos.FrequencyPreferences{}
	return b
}

func (b *CbsdProtoPayloadBuilder) WithSingleStepEnabled() *CbsdProtoPayloadBuilder {
	b.Payload.SingleStepEnabled = true
	return b
}

func (b *CbsdProtoPayloadBuilder) WithCbsdCategory(c string) *CbsdProtoPayloadBuilder {
	b.Payload.CbsdCategory = c
	return b
}

func (b *CbsdProtoPayloadBuilder) WithEmptyInstallationParam() *CbsdProtoPayloadBuilder {
	b.Payload.InstallationParam = &protos.InstallationParam{}
	return b
}

func (b *CbsdProtoPayloadBuilder) WithAntennaGain(gain float64) *CbsdProtoPayloadBuilder {
	b.Payload.InstallationParam.AntennaGain = wrapperspb.Double(gain)
	return b
}

func (b *CbsdProtoPayloadBuilder) WithEnodebdInstallationParams() *CbsdProtoPayloadBuilder {
	b.Payload.InstallationParam = &protos.InstallationParam{
		LatitudeDeg:      wrapperspb.Double(10.5),
		LongitudeDeg:     wrapperspb.Double(11.5),
		IndoorDeployment: wrapperspb.Bool(true),
		HeightM:          wrapperspb.Double(12.5),
		HeightType:       wrapperspb.String("agl"),
	}
	return b
}

func (b *CbsdProtoPayloadBuilder) WithFullInstallationParam() *CbsdProtoPayloadBuilder {
	b.Payload.InstallationParam = &protos.InstallationParam{
		LatitudeDeg:      wrapperspb.Double(10.5),
		LongitudeDeg:     wrapperspb.Double(11.5),
		IndoorDeployment: wrapperspb.Bool(true),
		HeightM:          wrapperspb.Double(12.5),
		HeightType:       wrapperspb.String("agl"),
		AntennaGain:      wrapperspb.Double(4.5),
	}
	return b
}

func (b *CbsdProtoPayloadBuilder) WithIncompleteInstallationParam() *CbsdProtoPayloadBuilder {
	b.Payload.InstallationParam = &protos.InstallationParam{
		LatitudeDeg:      wrapperspb.Double(10.5),
		LongitudeDeg:     wrapperspb.Double(11.5),
		IndoorDeployment: wrapperspb.Bool(true),
	}
	return b
}

type DetailedDBCbsdBuilder struct {
	Details *storage.DetailedCbsd
}

func NewDetailedDBCbsdBuilder() *DetailedDBCbsdBuilder {
	return &DetailedDBCbsdBuilder{Details: &storage.DetailedCbsd{}}
}

func (b *DetailedDBCbsdBuilder) WithCbsd(cbsd *storage.DBCbsd, state string, desiredState string) *DetailedDBCbsdBuilder {
	b.Details.Cbsd = cbsd
	b.Details.CbsdState = &storage.DBCbsdState{Name: db.MakeString(state)}
	b.Details.DesiredState = &storage.DBCbsdState{Name: db.MakeString(desiredState)}
	return b
}

func (b *DetailedDBCbsdBuilder) WithGrant(state string, frequencyMHz int64, grantExpireTime time.Time, transmitExpireTime time.Time) *DetailedDBCbsdBuilder {
	grant := &storage.DetailedGrant{
		Grant: &storage.DBGrant{
			GrantExpireTime:    db.MakeTime(grantExpireTime),
			TransmitExpireTime: db.MakeTime(transmitExpireTime),
			LowFrequencyHz:     db.MakeInt((frequencyMHz - 10) * 1e6),
			HighFrequencyHz:    db.MakeInt((frequencyMHz + 10) * 1e6),
			MaxEirp:            db.MakeFloat(35),
		},
		GrantState: &storage.DBGrantState{
			Name: db.MakeString(state),
		},
	}
	b.Details.Grants = append(b.Details.Grants, grant)
	return b
}

func (b *DetailedDBCbsdBuilder) WithAmcGrant(state string, frequencyMHz int64, lastHeartbeatTime time.Time, grantId string, heartbeatInterval int64) *DetailedDBCbsdBuilder {
	grant := &storage.DetailedGrant{
		Grant: &storage.DBGrant{
			GrantId:                  db.MakeString(grantId),
			LowFrequencyHz:           db.MakeInt((frequencyMHz - 10) * 1e6),
			HighFrequencyHz:          db.MakeInt((frequencyMHz + 10) * 1e6),
			LastHeartbeatRequestTime: db.MakeTime(lastHeartbeatTime),
			HeartbeatIntervalSec:     db.MakeInt(heartbeatInterval),
		},
		GrantState: &storage.DBGrantState{
			Name: db.MakeString(state),
		},
	}
	b.Details.Grants = append(b.Details.Grants, grant)
	return b
}

type DetailedProtoCbsdBuilder struct {
	Details *protos.CbsdDetails
}

func NewDetailedProtoCbsdBuilder(builder *CbsdProtoPayloadBuilder) *DetailedProtoCbsdBuilder {
	return &DetailedProtoCbsdBuilder{
		Details: &protos.CbsdDetails{
			Data:     builder.Payload,
			State:    registered,
			IsActive: false,
			Grants:   []*protos.GrantDetails{},
		},
	}
}

func (b *DetailedProtoCbsdBuilder) WithId(id int64) *DetailedProtoCbsdBuilder {
	b.Details.Id = id
	return b
}

func (b *DetailedProtoCbsdBuilder) WithCbsdId(id string) *DetailedProtoCbsdBuilder {
	b.Details.CbsdId = id
	return b
}

func (b *DetailedProtoCbsdBuilder) WithState(state string) *DetailedProtoCbsdBuilder {
	b.Details.State = state
	return b
}

func (b *DetailedProtoCbsdBuilder) Active() *DetailedProtoCbsdBuilder {
	b.Details.IsActive = true
	return b
}

func (b *DetailedProtoCbsdBuilder) WithGrant() *DetailedProtoCbsdBuilder {
	b.Details.Grants = append(b.Details.Grants, &protos.GrantDetails{
		BandwidthMhz:            20,
		FrequencyMhz:            3610,
		MaxEirp:                 35,
		State:                   authorized,
		TransmitExpireTimestamp: 456,
		GrantExpireTimestamp:    123,
	})
	return b
}

func (b *DetailedProtoCbsdBuilder) WithDefaultTestData() *DetailedProtoCbsdBuilder {
	return b.WithCbsdId(someCbsdId).WithId(cbsdId).WithState(registered).WithGrant()
}

func GetDetailedProtoCbsdList(builder *DetailedProtoCbsdBuilder) *protos.ListCbsdResponse {
	return &protos.ListCbsdResponse{
		Details:    []*protos.CbsdDetails{builder.Details},
		TotalCount: 1,
	}
}

func GetMutableDBCbsd(cbsd *storage.DBCbsd, state string) *storage.MutableCbsd {
	return &storage.MutableCbsd{
		Cbsd: cbsd,
		DesiredState: &storage.DBCbsdState{
			Name: db.MakeString(state),
		},
	}
}

func GetDetailedDBCbsdList(cbsd *storage.DetailedCbsd) *storage.DetailedCbsdList {
	return &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{cbsd},
		Count: 1,
	}
}

type CbsdModelPayloadBuilder struct {
	Payload *models.Cbsd
}

func NewCbsdModelPayloadBuilder() *CbsdModelPayloadBuilder {
	return &CbsdModelPayloadBuilder{Payload: &models.Cbsd{
		Capabilities: models.Capabilities{
			MaxPower:         to_pointer.Float(20),
			MinPower:         to_pointer.Float(10),
			NumberOfAntennas: 2,
			MaxIbwMhz:        150,
		},
		SingleStepEnabled: false,
		CbsdCategory:      catB,
		DesiredState:      registered,
		FrequencyPreferences: models.FrequencyPreferences{
			BandwidthMhz:   20,
			FrequenciesMhz: []int64{3600},
		},
		FccID:                     someFccId,
		SerialNumber:              someSerialNumber,
		UserID:                    someUserId,
		CbsdID:                    someCbsdId,
		State:                     registered,
		CarrierAggregationEnabled: false,
		GrantRedundancy:           true,
	}}
}

func (b *CbsdModelPayloadBuilder) WithSingleStepEnabled() *CbsdModelPayloadBuilder {
	b.Payload.SingleStepEnabled = true
	return b
}

func (b *CbsdModelPayloadBuilder) WithCbsdCategory(c string) *CbsdModelPayloadBuilder {
	b.Payload.CbsdCategory = c
	return b
}

func (b *CbsdModelPayloadBuilder) WithGrant() *CbsdModelPayloadBuilder {
	b.Payload.Grants = append(b.Payload.Grants, &models.Grant{
		BandwidthMhz:       20,
		FrequencyMhz:       3610,
		GrantExpireTime:    to_pointer.TimeToDateTime(123),
		MaxEirp:            35,
		State:              authorized,
		TransmitExpireTime: to_pointer.TimeToDateTime(456),
	})
	return b
}

func GetPaginatedCbsds(builder *CbsdModelPayloadBuilder) *models.PaginatedCbsds {
	return &models.PaginatedCbsds{
		Cbsds:      []*models.Cbsd{builder.Payload},
		TotalCount: 1,
	}
}

type MutableCbsdModelBuilder struct {
	Payload *models.MutableCbsd
}

func NewMutableCbsdModelPayloadBuilder() *MutableCbsdModelBuilder {
	return &MutableCbsdModelBuilder{Payload: &models.MutableCbsd{
		Capabilities: models.Capabilities{
			MaxPower:         to_pointer.Float(20),
			MinPower:         to_pointer.Float(10),
			NumberOfAntennas: 2,
			MaxIbwMhz:        150,
		},
		DesiredState:              registered,
		SingleStepEnabled:         to_pointer.Bool(false),
		CarrierAggregationEnabled: to_pointer.Bool(false),
		GrantRedundancy:           to_pointer.Bool(true),
		CbsdCategory:              catB,
		FrequencyPreferences: models.FrequencyPreferences{
			BandwidthMhz:   20,
			FrequenciesMhz: []int64{3600},
		},
		FccID:        someFccId,
		SerialNumber: someSerialNumber,
		UserID:       someUserId,
	}}
}

func (b *MutableCbsdModelBuilder) Empty() *MutableCbsdModelBuilder {
	b.Payload = &models.MutableCbsd{}
	return b
}

func (b *MutableCbsdModelBuilder) WithSerialNumber(serial string) *MutableCbsdModelBuilder {
	b.Payload.SerialNumber = serial
	return b
}

func (b *MutableCbsdModelBuilder) WithFccId(id string) *MutableCbsdModelBuilder {
	b.Payload.FccID = id
	return b
}

func (b *MutableCbsdModelBuilder) WithEmptyInstallationParam() *MutableCbsdModelBuilder {
	b.Payload.InstallationParam = models.MutableInstallationParam{}
	return b
}

func (b *MutableCbsdModelBuilder) WithAntennaGain(gain float64) *MutableCbsdModelBuilder {
	b.Payload.InstallationParam.AntennaGain = to_pointer.Float(gain)
	return b
}

func (b *MutableCbsdModelBuilder) WithUserId(id string) *MutableCbsdModelBuilder {
	b.Payload.UserID = id
	return b
}

func (b *MutableCbsdModelBuilder) WithFrequencies(frequencies []int64) *MutableCbsdModelBuilder {
	b.Payload.FrequencyPreferences.FrequenciesMhz = frequencies
	return b
}

func (b *MutableCbsdModelBuilder) WithBandwidth(bandwidth int64) *MutableCbsdModelBuilder {
	b.Payload.FrequencyPreferences.BandwidthMhz = bandwidth
	return b
}

func (b *MutableCbsdModelBuilder) WithNumberOfAntennas(number int64) *MutableCbsdModelBuilder {
	b.Payload.Capabilities.NumberOfAntennas = number
	return b
}

func (b *MutableCbsdModelBuilder) WithMaxIbwMhz(number int64) *MutableCbsdModelBuilder {
	b.Payload.Capabilities.MaxIbwMhz = number
	return b
}

func (b *MutableCbsdModelBuilder) WithSingleStepEnabled(enabled *bool) *MutableCbsdModelBuilder {
	if enabled == nil {
		b.Payload.SingleStepEnabled = nil
	} else {
		b.Payload.SingleStepEnabled = to_pointer.Bool(*enabled)
	}
	return b
}

func (b *MutableCbsdModelBuilder) WithCarrierAggregationEnabled(enabled *bool) *MutableCbsdModelBuilder {
	if enabled == nil {
		b.Payload.CarrierAggregationEnabled = nil
	} else {
		b.Payload.CarrierAggregationEnabled = to_pointer.Bool(*enabled)
	}
	return b
}

func (b *MutableCbsdModelBuilder) WithGrantRedundancy(enabled *bool) *MutableCbsdModelBuilder {
	if enabled == nil {
		b.Payload.GrantRedundancy = nil
	} else {
		b.Payload.GrantRedundancy = to_pointer.Bool(*enabled)
	}
	return b
}

func (b *MutableCbsdModelBuilder) WithMinPower(power *float64) *MutableCbsdModelBuilder {
	if power == nil {
		b.Payload.Capabilities.MinPower = nil
	} else {
		b.Payload.Capabilities.MinPower = to_pointer.Float(*power)
	}
	return b
}

func (b *MutableCbsdModelBuilder) WithMaxPower(power *float64) *MutableCbsdModelBuilder {
	if power == nil {
		b.Payload.Capabilities.MaxPower = nil
	} else {
		b.Payload.Capabilities.MaxPower = to_pointer.Float(*power)
	}
	return b
}

func (b *MutableCbsdModelBuilder) WithCbsdCategory(c string) *MutableCbsdModelBuilder {
	b.Payload.CbsdCategory = c
	return b
}

type DPLogBuilder struct {
	Log *logs_pusher.DPLog
}

func NewDPLogBuilder(from string, to string, name string) *DPLogBuilder {
	return &DPLogBuilder{Log: &logs_pusher.DPLog{
		EventTimestamp:   clock.Now().Unix(),
		LogFrom:          from,
		LogTo:            to,
		LogName:          name,
		LogMessage:       "null",
		CbsdSerialNumber: "some_serial_number",
		NetworkId:        "some_network",
		FccId:            "some_fcc_id",
	}}
}

func (b *DPLogBuilder) WithLogMessage(m string) *DPLogBuilder {
	b.Log.LogMessage = m
	return b
}

func NewRequestBuilder(id int64, cbsdId int64, typeId int64, payload string) *RequestBuilder {
	return &RequestBuilder{
		Request: &storage.DBRequest{
			Id:      db.MakeInt(id),
			TypeId:  db.MakeInt(typeId),
			CbsdId:  db.MakeInt(cbsdId),
			Payload: db.MakeString(payload),
		},
	}
}

type RequestBuilder struct {
	Request *storage.DBRequest
}
