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
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testdataDir     = "testdata"
	goldenFilepath1 = filepath.Join(testdataDir, "gen_time_encoded.golden")
	goldenFilepath2 = filepath.Join(testdataDir, "record_encoded.golden")

	generalizedTime = "2021-02-18T05:13:26.019519+00:00"
)

func TestGeneralizedTime(t *testing.T) {
	encodedGeneralizedTime, err := readFile(goldenFilepath1)
	assert.NoError(t, err)

	ptime, err := time.Parse(time.RFC3339Nano, generalizedTime)
	assert.NoError(t, err)

	ret := encodeGeneralizedTime(ptime)
	if !reflect.DeepEqual(encodedGeneralizedTime, ret) {
		t.Errorf("Bad result: %q → %v (expected %v)\n %v", generalizedTime, ret, encodedGeneralizedTime, hex.Dump(ret))
	}
}

func TestEpsIRIRecord(t *testing.T) {
	encodedRecord, err := readFile(goldenFilepath2)
	assert.NoError(t, err)

	var record EpsIRIRecord
	if err := record.Decode(encodedRecord); err != nil {
		t.Errorf("Decoding record failed: %v", err)
	}

	assert.Equal(t, HeaderVersion, record.Header.Version)
	assert.Equal(t, uint64(0x866cb3979084570), record.Header.CorrelationID)
	assert.Equal(t, "29f28e1c-f230-486a-a860-f5a784ab9178", record.Header.XID.String())

	assert.Equal(t, BearerActivation, record.Payload.EPSEvent)
	assert.Equal(t, InitiatorNotAvailable, record.Payload.Initiator)
	assert.Equal(t, GetOID(), record.Payload.Hi2epsDomainID)
}

func readFile(fname string) ([]byte, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return []byte{}, err
	}
	return base64.StdEncoding.DecodeString(string(content))
}
