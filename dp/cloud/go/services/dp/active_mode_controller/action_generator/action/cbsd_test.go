package action

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestSetAvailableFrequencies(t *testing.T) {
	cbsd := &storage.DBCbsd{
		Id: db.MakeInt(1),
		Channels: []storage.Channel{{
			LowFrequencyHz:  3590e6,
			HighFrequencyHz: 3610e6,
			MaxEirp:         37,
		}},
		AvailableFrequencies: nil,
	}
	expected := &UpdateCbsd{
		Data: &storage.DBCbsd{
			Id:                   db.MakeInt(1),
			AvailableFrequencies: []uint32{3584, 3584, 1024, 1024},
		},
		Mask: db.NewIncludeMask("available_frequencies"),
	}

	newFrequencies := SetAvailableFrequences(cbsd)
	assert.Equal(t, expected, newFrequencies)
}

func TestRemoveIdleGrants(t *testing.T) {
	testData := []struct {
		name     string
		cbsd     *storage.DetailedCbsd
		expected []Action
	}{{
		name: "Should do nothing when no idle grants",
		cbsd: &storage.DetailedCbsd{
			Cbsd: &storage.DBCbsd{},
			Grants: []*storage.DetailedGrant{{
				GrantState: &storage.DBGrantState{Name: db.MakeString("granted")},
				Grant:      &storage.DBGrant{Id: db.MakeInt(1)},
			}},
		},
		expected: nil,
	}, {
		name: "Should delete idle grant without affecting frequencies",
		cbsd: &storage.DetailedCbsd{
			Cbsd: &storage.DBCbsd{},
			Grants: []*storage.DetailedGrant{{
				GrantState: &storage.DBGrantState{Name: db.MakeString("idle")},
				Grant: &storage.DBGrant{
					Id:              db.MakeInt(1),
					LowFrequencyHz:  db.MakeInt(3590e6),
					HighFrequencyHz: db.MakeInt(3610e6),
				}}},
		},
		expected: []Action{&DeleteGrant{Id: 1}},
	}, {
		name: "Should delete idle grant without affecting frequencies",
		cbsd: &storage.DetailedCbsd{
			Cbsd: &storage.DBCbsd{AvailableFrequencies: []uint32{0, 0, 0, 0}},
			Grants: []*storage.DetailedGrant{{
				GrantState: &storage.DBGrantState{Name: db.MakeString("idle")},
				Grant:      &storage.DBGrant{Id: db.MakeInt(1)}},
			}},
		expected: []Action{&DeleteGrant{Id: 1}},
	}, {
		name: "Should delete idle grant and update frequencies",
		cbsd: &storage.DetailedCbsd{
			Cbsd: &storage.DBCbsd{
				Id:                   db.MakeInt(1),
				AvailableFrequencies: []uint32{0b1111, 0b110, 0b1100, 0b1010},
			},
			Grants: []*storage.DetailedGrant{{
				GrantState: &storage.DBGrantState{Name: db.MakeString("idle")},
				Grant: &storage.DBGrant{
					Id:              db.MakeInt(1),
					LowFrequencyHz:  db.MakeInt(35625e5),
					HighFrequencyHz: db.MakeInt(35675e5),
				}},
			},
		},
		expected: []Action{&UpdateCbsd{
			Data: &storage.DBCbsd{
				Id:                   db.MakeInt(1),
				AvailableFrequencies: []uint32{0b0111, 0b110, 0b1100, 0b1010},
			},
			Mask: db.NewIncludeMask("available_frequencies"),
		}, &DeleteGrant{Id: 1}},
	},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := RemoveIdleGrants(tt.cbsd)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
