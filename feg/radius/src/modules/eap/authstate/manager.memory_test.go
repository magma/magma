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

package authstate

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"fbc/cwf/radius/modules/eap/packet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

type TestEapMethodState struct {
}

func TestBasicInsertGet(t *testing.T) {
	// Arrange
	manager := NewMemoryManager()
	authReq := createRadiusPacket("called", "calling")

	// Act and Assert
	performSingleReadWriteDeleteReadTest(t, manager, authReq)
}

func performSingleReadWriteDeleteReadTest(t *testing.T, manager Manager, authReq radius.Packet) {
	// Arrange (randomize state)
	correlationID := rand.Intn(9999999)
	eapType := packet.EAPTypeAKA
	protocolState := strconv.Itoa(rand.Intn(999999))

	// Act
	stateBeforeWrite, errBeforeWrite := manager.Get(&authReq, packet.EAPTypeAKA)
	err := manager.Set(&authReq, eapType, Container{
		LogCorrelationID: uint64(correlationID),
		EapType:          eapType,
		ProtocolState:    protocolState,
	})
	require.Equal(t, nil, err)

	stateAfterWrite, errAfterWrite := manager.Get(&authReq, packet.EAPTypeAKA)
	manager.Reset(&authReq, eapType)
	stateAfterReset, errAfterReset := manager.Get(&authReq, packet.EAPTypeAKA)

	// Assert
	assert.Nil(t, stateBeforeWrite)
	assert.NotNil(t, errBeforeWrite)
	require.Equal(t, errBeforeWrite.Error(), "eap state not found")

	assert.NotNil(t, stateAfterWrite)
	assert.Nil(t, errAfterWrite)
	assert.NotNil(t, stateAfterWrite)
	assert.Equal(t, stateAfterWrite.LogCorrelationID, uint64(correlationID))
	assert.Equal(t, stateAfterWrite.EapType, eapType)
	assert.Equal(t, stateAfterWrite.ProtocolState, protocolState)

	assert.Nil(t, stateAfterReset)
	assert.NotNil(t, errAfterReset)
	assert.Equal(t, errAfterReset.Error(), "eap state not found")
}

func TestMultipleConcurrentInsertDeleteGet(t *testing.T) {
	// Arrange
	degOfParallelism := 101
	reqPerConcurrentContext := 100
	var wg sync.WaitGroup
	wg.Add(degOfParallelism)
	manager := NewMemoryManager()

	// Act
	for i := 0; i < degOfParallelism; i++ {
		go func(called string, calling string) {
			defer wg.Done()
			authReq := createRadiusPacket(called, calling)
			for i := 0; i < reqPerConcurrentContext; i++ {
				performSingleReadWriteDeleteReadTest(t, manager, authReq)
			}
		}(fmt.Sprintf("called%d", i), fmt.Sprintf("calling%d", i))
	}
	wg.Wait()

	// Assert
	// nothing to do (assert will happen in the go routines spawned above)
}

func createRadiusPacket(called string, calling string) radius.Packet {
	p := radius.Packet{}
	p.Attributes.Add(rfc2865.CallingStationID_Type, radius.Attribute(calling))
	p.Attributes.Add(rfc2865.CalledStationID_Type, radius.Attribute(called))
	return p
}
