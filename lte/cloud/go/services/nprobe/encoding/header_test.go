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

package encoding

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
)

type attributeTC struct {
	name string
	at   []Attribute
	bs   []byte
	ok   bool
}

type iriHeaderTC struct {
	name string
	hdr  EpsIRIHeader
	bs   []byte
	ok   bool
}

var (
	tc1 = []attributeTC{
		{
			name: "valid attribute",
			at: []Attribute{
				{
					Tag:   1,
					Len:   1,
					Value: []byte{0x0a},
				},
				{
					Tag:   1,
					Len:   2,
					Value: []byte{0x0b, 0x0c},
				},
			},
			bs: []byte{0x00, 0x01, 0x00, 0x01, 0x0a, 0x00, 0x01, 0x00, 0x02, 0x0b, 0x0c},
			ok: true,
		},
		{
			name: "invalid attribute",
			at: []Attribute{
				{
					Tag:   1,
					Len:   2,
					Value: []byte{0x0a},
				},
			},
			ok: false,
		},
	}

	uid = uuid.Must(uuid.NewV4())
	tc2 = []iriHeaderTC{
		{
			name: "valid header",
			hdr: EpsIRIHeader{
				Version:       1,
				PduType:       2,
				HeaderLength:  HeaderFixLen,
				XID:           uid,
				PayloadLength: 0,
			},
			bs: []byte{},
			ok: true,
		},
		{
			name: "valid header with conditional attributes",
			hdr: EpsIRIHeader{
				Version:               1,
				PduType:               2,
				PayloadLength:         0,
				HeaderLength:          HeaderFixLen + 11, // size of tc1 attributes
				XID:                   uid,
				ConditionalAttributes: tc1[0].at,
			},
			bs: []byte{},
			ok: true,
		},
	}
)

func TestConditionalAttributesMarshaling(t *testing.T) {
	for _, test := range tc1 {
		t.Run(test.name, func(t *testing.T) {
			b := marshalAttributes(test.at)
			if len(test.bs) > 0 && cmp.Diff(test.bs, b) != "" {
				t.Fatalf("unexpected marsheled attribute :\n%v", b)
			}

			got, err := parseAttributes(b)
			if err != nil && test.ok {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && !test.ok {
				t.Fatal("expected an error, but none occurred")
			}
			if err != nil {
				t.Logf("expected failure: %v", err)
				return
			}
			if diff := cmp.Diff(test.at, got); diff != "" {
				t.Fatalf("unexpected attribute (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEpsIRIHeaderMarshaling(t *testing.T) {
	for _, test := range tc2 {
		t.Run(test.name, func(t *testing.T) {
			b := test.hdr.Marshal()
			if len(test.bs) > 0 && cmp.Diff(test.bs, b) != "" {
				t.Fatalf("unexpected marsheled attribute :\n%v", b)
			}

			var got EpsIRIHeader
			err := got.Unmarshal(b)
			if err != nil && test.ok {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && !test.ok {
				t.Fatal("expected an error, but none occurred")
			}
			if err != nil {
				t.Logf("expected error: %v", err)
				return
			}
			if diff := cmp.Diff(test.hdr, got); diff != "" {
				t.Fatalf("unexpected header (-want +got):\n%s", diff)
			}
		})
	}
}
