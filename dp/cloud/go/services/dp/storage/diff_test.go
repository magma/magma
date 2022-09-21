package storage_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestShouldENodeBDUpdateInstallationParams(t *testing.T) {
	// TODO switch to builders when available
	testData := []struct {
		name     string
		prev     *storage.DBCbsd
		next     *storage.DBCbsd
		expected bool
	}{{
		name: "should update if installation parameters have changes",
		prev: &storage.DBCbsd{
			AntennaGainDbi:    db.MakeFloat(10),
			CbsdCategory:      sql.NullString{},
			LatitudeDeg:       db.MakeFloat(50),
			LongitudeDeg:      db.MakeFloat(100),
			HeightM:           db.MakeFloat(5),
			HeightType:        db.MakeString("AGL"),
			IndoorDeployment:  db.MakeBool(false),
			SingleStepEnabled: db.MakeBool(true),
		},
		next: &storage.DBCbsd{
			AntennaGainDbi:   db.MakeFloat(20),
			CbsdCategory:     db.MakeString("A"),
			LatitudeDeg:      db.MakeFloat(50),
			LongitudeDeg:     db.MakeFloat(100),
			HeightM:          db.MakeFloat(8),
			HeightType:       db.MakeString("AGL"),
			IndoorDeployment: db.MakeBool(true),
		},
		expected: true,
	}, {
		name: "should not update all parameters are the same",
		prev: &storage.DBCbsd{
			AntennaGainDbi:    db.MakeFloat(20),
			CbsdCategory:      db.MakeString("A"),
			LatitudeDeg:       db.MakeFloat(50),
			LongitudeDeg:      db.MakeFloat(100),
			HeightM:           db.MakeFloat(8),
			HeightType:        db.MakeString("AGL"),
			IndoorDeployment:  db.MakeBool(true),
			SingleStepEnabled: db.MakeBool(true),
		},
		next: &storage.DBCbsd{
			AntennaGainDbi:   db.MakeFloat(20),
			CbsdCategory:     db.MakeString("A"),
			LatitudeDeg:      db.MakeFloat(50),
			LongitudeDeg:     db.MakeFloat(100),
			HeightM:          db.MakeFloat(8),
			HeightType:       db.MakeString("AGL"),
			IndoorDeployment: db.MakeBool(true),
		},
		expected: false,
	}, {
		name: "should update if any of coordinates are empty",
		prev: &storage.DBCbsd{
			AntennaGainDbi:    db.MakeFloat(20),
			CbsdCategory:      db.MakeString("A"),
			LatitudeDeg:       sql.NullFloat64{},
			LongitudeDeg:      sql.NullFloat64{},
			HeightM:           db.MakeFloat(8),
			HeightType:        db.MakeString("AGL"),
			SingleStepEnabled: db.MakeBool(true),
			IndoorDeployment:  db.MakeBool(true),
		},
		next: &storage.DBCbsd{
			AntennaGainDbi:   db.MakeFloat(20),
			CbsdCategory:     db.MakeString("A"),
			LatitudeDeg:      sql.NullFloat64{},
			LongitudeDeg:     sql.NullFloat64{},
			HeightM:          db.MakeFloat(8),
			HeightType:       db.MakeString("AGL"),
			IndoorDeployment: db.MakeBool(true),
		},
		expected: true,
	}, {
		name: "should not update if coordinates changed less than 10m",
		prev: &storage.DBCbsd{
			AntennaGainDbi:    db.MakeFloat(20),
			CbsdCategory:      db.MakeString("A"),
			LatitudeDeg:       db.MakeFloat(50),
			LongitudeDeg:      db.MakeFloat(100),
			HeightM:           db.MakeFloat(8),
			HeightType:        db.MakeString("AGL"),
			SingleStepEnabled: db.MakeBool(true),
			IndoorDeployment:  db.MakeBool(true),
		},
		next: &storage.DBCbsd{
			AntennaGainDbi:   db.MakeFloat(20),
			CbsdCategory:     db.MakeString("A"),
			LatitudeDeg:      db.MakeFloat(50.00006),
			LongitudeDeg:     db.MakeFloat(100.0001),
			HeightM:          db.MakeFloat(8),
			HeightType:       db.MakeString("AGL"),
			IndoorDeployment: db.MakeBool(true),
		},
		expected: false,
	}, {
		name: "should update if coordinates changed more than 10m",
		prev: &storage.DBCbsd{
			AntennaGainDbi:    db.MakeFloat(20),
			CbsdCategory:      db.MakeString("A"),
			LatitudeDeg:       db.MakeFloat(50),
			LongitudeDeg:      db.MakeFloat(100),
			HeightM:           db.MakeFloat(8),
			HeightType:        db.MakeString("AGL"),
			SingleStepEnabled: db.MakeBool(true),
			IndoorDeployment:  db.MakeBool(true),
		},
		next: &storage.DBCbsd{
			AntennaGainDbi:   db.MakeFloat(20),
			CbsdCategory:     db.MakeString("A"),
			LatitudeDeg:      db.MakeFloat(50.00007),
			LongitudeDeg:     db.MakeFloat(100.0001),
			HeightM:          db.MakeFloat(8),
			HeightType:       db.MakeString("AGL"),
			IndoorDeployment: db.MakeBool(true),
		},
		expected: true,
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := storage.ShouldEnodebdUpdateInstallationParams(tt.prev, tt.next)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
