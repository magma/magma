package storage_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestShouldENodeBDUpdate(t *testing.T) {
	// TODO switch to builders when available
	testData := []struct {
		name     string
		prev     *storage.DBCbsd
		next     *storage.DBCbsd
		expected bool
	}{{
		name: "should update if installation parameters have changes",
		prev: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(10),
			CbsdCategory:        sql.NullString{},
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(5),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(false),
			CpiDigitalSignature: sql.NullString{},
		},
		next: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		expected: true,
	}, {
		name: "should not update all parameters are the same",
		prev: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		next: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		expected: false,
	}, {
		name: "should not update if cbsd had cpi signature",
		prev: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(10),
			CbsdCategory:        sql.NullString{},
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(5),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(false),
			CpiDigitalSignature: db.MakeString("some signature"),
		},
		next: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		expected: false,
	}, {
		name: "should update if any of coordinates are empty",
		prev: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         sql.NullFloat64{},
			LongitudeDeg:        sql.NullFloat64{},
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		next: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         sql.NullFloat64{},
			LongitudeDeg:        sql.NullFloat64{},
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		expected: true,
	}, {
		name: "should not update if coordinates changed less than 10m",
		prev: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		next: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50.00006),
			LongitudeDeg:        db.MakeFloat(100.0001),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		expected: false,
	}, {
		name: "should update if coordinates changed more than 10m",
		prev: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50),
			LongitudeDeg:        db.MakeFloat(100),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		next: &storage.DBCbsd{
			AntennaGain:         db.MakeFloat(20),
			CbsdCategory:        db.MakeString("A"),
			LatitudeDeg:         db.MakeFloat(50.00007),
			LongitudeDeg:        db.MakeFloat(100.0001),
			HeightM:             db.MakeFloat(8),
			HeightType:          db.MakeString("AGL"),
			IndoorDeployment:    db.MakeBool(true),
			CpiDigitalSignature: sql.NullString{},
		},
		expected: true,
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := storage.ShouldENodeBDUpdate(tt.prev, tt.next)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
