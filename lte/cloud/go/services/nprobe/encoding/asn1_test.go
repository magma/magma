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
	"encoding/asn1"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	goldenFilepath = filepath.Join(testdataDir, "payload_asn1.golden")
)

func TestEpsIRIContent(t *testing.T) {
	encodedPayload, err := readFile(goldenFilepath)
	assert.NoError(t, err)

	var content EpsIRIContent
	if _, err := asn1.UnmarshalWithParams(encodedPayload, &content, IRIBeginRecord); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	assert.Equal(t, content.EPSEvent, BearerActivation)
	assert.Equal(t, content.Initiator, InitiatorNotAvailable)
	assert.Equal(t, content.Hi2epsDomainID, GetOID())

	imsi := encodeIMSI("IMSI001010000000001")
	assert.Equal(t, imsi, content.PartyInformation[0].PartyIdentity.IMSI)

	apn := encodeAPN("magma.ipv4")
	assert.Equal(t, apn, content.EPSSpecificParameters.APN)

	bid := []byte("IMSI001010000000001-104552")
	assert.Equal(t, bid, content.EPSSpecificParameters.EPSBearerIdentity)
	assert.Equal(t, content.EPSSpecificParameters.BearerActivationType, DefaultBearer)
	assert.Equal(t, content.EPSSpecificParameters.RATType, []byte{RatTypeEutran})

	marshaledContent, err := asn1.MarshalWithParams(content, IRIBeginRecord)
	if err != nil {
		t.Errorf("Marshal failed: %v", err)
	}
	assert.Equal(t, marshaledContent, encodedPayload)
}
