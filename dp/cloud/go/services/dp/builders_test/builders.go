package builders_test

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
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
	someModel              = "some_model"
)

type DBCbsdBuilder struct {
	Cbsd *storage.DBCbsd
}

func NewDBCbsdBuilder() *DBCbsdBuilder {
	return &DBCbsdBuilder{
		Cbsd: &storage.DBCbsd{
			UserId:                  db.MakeString(someUserId),
			FccId:                   db.MakeString(someFccId),
			CbsdSerialNumber:        db.MakeString(someSerialNumber),
			PreferredBandwidthMHz:   db.MakeInt(20),
			PreferredFrequenciesMHz: db.MakeString("[3600]"),
			MinPower:                db.MakeFloat(10),
			MaxPower:                db.MakeFloat(20),
			NumberOfPorts:           db.MakeInt(2),
			CbsdCategory:            db.MakeString(catB),
			SingleStepEnabled:       db.MakeBool(false),
		},
	}
}

func (b *DBCbsdBuilder) WithId(id int64) *DBCbsdBuilder {
	b.Cbsd.Id = db.MakeInt(id)
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
	b.Cbsd.AntennaAzimuthDeg = db.MakeInt(1)
	b.Cbsd.AntennaDowntiltDeg = db.MakeInt(2)
	b.Cbsd.AntennaBeamwidthDeg = db.MakeInt(3)
	b.Cbsd.AntennaModel = db.MakeString(someModel)
	b.Cbsd.AntennaGain = db.MakeFloat(4.5)
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

func (b *DBCbsdBuilder) WithCbsdCategory(cat string) *DBCbsdBuilder {
	b.Cbsd.CbsdCategory = db.MakeString(cat)
	return b
}

func (b *DBCbsdBuilder) WithDefaulValues() *DBCbsdBuilder {
	return b.WithCbsdCategory(catB).WithSingleStepEnabled(false).WithIndoorDeployment(false)
}

type DBGrantBuilder struct {
	Grant *storage.DBGrant
}

func NewDBGrantBuilder() *DBGrantBuilder {
	return &DBGrantBuilder{
		Grant: &storage.DBGrant{
			GrantExpireTime:    db.MakeTime(time.Unix(123, 0).UTC()),
			TransmitExpireTime: db.MakeTime(time.Unix(456, 0).UTC()),
			LowFrequency:       db.MakeInt(3600 * 1e6),
			HighFrequency:      db.MakeInt(3620 * 1e6),
			MaxEirp:            db.MakeFloat(35),
			GrantId:            db.MakeString("some_grant_id"),
		},
	}
}

func (b *DBGrantBuilder) WithId(id int64) *DBGrantBuilder {
	b.Grant.Id = db.MakeInt(id)
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
			},
			DesiredState:      registered,
			CbsdCategory:      catB,
			SingleStepEnabled: false,
		},
	}
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

func NewDetailedDBCbsdBuilder(builder *DBCbsdBuilder) *DetailedDBCbsdBuilder {
	return &DetailedDBCbsdBuilder{
		Details: &storage.DetailedCbsd{
			Cbsd: builder.Cbsd,
		},
	}
}

func (b *DetailedDBCbsdBuilder) WithGrant() *DetailedDBCbsdBuilder {
	b.Details.Grant = &storage.DBGrant{
		GrantExpireTime:    db.MakeTime(time.Unix(123, 0).UTC()),
		TransmitExpireTime: db.MakeTime(time.Unix(456, 0).UTC()),
		LowFrequency:       db.MakeInt(3600 * 1e6),
		HighFrequency:      db.MakeInt(3620 * 1e6),
		MaxEirp:            db.MakeFloat(35),
	}
	return b
}

func (b *DetailedDBCbsdBuilder) WithEmptyGrant() *DetailedDBCbsdBuilder {
	b.Details.Grant = &storage.DBGrant{}
	return b
}

func (b *DetailedDBCbsdBuilder) WithEmptyGrantState() *DetailedDBCbsdBuilder {
	b.Details.GrantState = &storage.DBGrantState{}
	return b
}

func (b *DetailedDBCbsdBuilder) WithCbsdState(state string) *DetailedDBCbsdBuilder {
	b.Details.CbsdState = &storage.DBCbsdState{
		Name: db.MakeString(state),
	}
	return b
}

func (b *DetailedDBCbsdBuilder) WithGrantState(state string) *DetailedDBCbsdBuilder {
	b.Details.GrantState = &storage.DBGrantState{
		Name: db.MakeString(state),
	}
	return b
}

func (b *DetailedDBCbsdBuilder) WithDesiredState(state string) *DetailedDBCbsdBuilder {
	b.Details.DesiredState = &storage.DBCbsdState{
		Name: db.MakeString(state),
	}
	return b
}

func (b *DetailedDBCbsdBuilder) WithDefaultTestData() *DetailedDBCbsdBuilder {
	return b.WithGrant().WithGrantState(authorized).WithCbsdState(registered).WithDesiredState(registered)
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
	b.Details.Grant = &protos.GrantDetails{
		BandwidthMhz:            20,
		FrequencyMhz:            3610,
		MaxEirp:                 35,
		State:                   authorized,
		TransmitExpireTimestamp: 456,
		GrantExpireTimestamp:    123,
	}
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

func GetMutableDBCbsd(cbsd *storage.DBCbsd) *storage.MutableCbsd {
	return &storage.MutableCbsd{
		Cbsd: cbsd,
		DesiredState: &storage.DBCbsdState{
			Name: db.MakeString(registered),
		},
	}
}

func GetDetailedDBCbsdList(builder *DetailedDBCbsdBuilder) *storage.DetailedCbsdList {
	cbsdList := &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{builder.Details},
	}
	cbsdList.Count = int64(len(cbsdList.Cbsds))
	return cbsdList
}

type CbsdModelPayloadBuilder struct {
	Payload *models.Cbsd
}

func NewCbsdModelPayloadBuilder() *CbsdModelPayloadBuilder {
	return &CbsdModelPayloadBuilder{Payload: &models.Cbsd{
		Capabilities: &models.Capabilities{
			MaxPower:         to_pointer.Float(20),
			MinPower:         to_pointer.Float(10),
			NumberOfAntennas: 2,
		},
		SingleStepEnabled: false,
		CbsdCategory:      catB,
		DesiredState:      registered,
		FrequencyPreferences: models.FrequencyPreferences{
			BandwidthMhz:   20,
			FrequenciesMhz: []int64{3600},
		},
		FccID:        someFccId,
		SerialNumber: someSerialNumber,
		UserID:       someUserId,
		CbsdID:       someCbsdId,
		State:        registered,
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
	b.Payload.Grant = &models.Grant{
		BandwidthMhz:       20,
		FrequencyMhz:       3610,
		GrantExpireTime:    to_pointer.TimeToDateTime(123),
		MaxEirp:            35,
		State:              authorized,
		TransmitExpireTime: to_pointer.TimeToDateTime(456),
	}
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
		Capabilities: &models.Capabilities{
			MaxPower:         to_pointer.Float(20),
			MinPower:         to_pointer.Float(10),
			NumberOfAntennas: 2,
		},
		DesiredState:      registered,
		SingleStepEnabled: to_pointer.Bool(false),
		CbsdCategory:      catB,
		FrequencyPreferences: models.FrequencyPreferences{
			BandwidthMhz:   20,
			FrequenciesMhz: []int64{3600},
		},
		FccID:        someFccId,
		SerialNumber: someSerialNumber,
		UserID:       someUserId,
	}}
}

func (b *MutableCbsdModelBuilder) withSingleStepEnabled() *MutableCbsdModelBuilder {
	b.Payload.SingleStepEnabled = to_pointer.Bool(true)
	return b
}

func (b *MutableCbsdModelBuilder) withCbsdCategory(c string) *MutableCbsdModelBuilder {
	b.Payload.CbsdCategory = c
	return b
}
