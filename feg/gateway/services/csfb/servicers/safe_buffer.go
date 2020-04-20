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

package servicers

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/message"

	"github.com/golang/protobuf/ptypes/any"
)

const (
	TimeOut = 1 // TimeOut in second for reading a chunk
)

// SafeBuffer has a list used like a buffer with protection of mutex
type SafeBuffer struct {
	chunkBuffer list.List
	sync.Mutex
	writeSignal chan bool
}

// NewSafeBuffer returns a SafeBuffer with writeSignal channel initialized
func NewSafeBuffer() (*SafeBuffer, error) {
	return &SafeBuffer{writeSignal: make(chan bool)}, nil
}

// GetNextMessage reads the next message in the buffer and decode the message
func (safeBuffer *SafeBuffer) GetNextMessage(timeOutInSec int) (decode.SGsMessageType, *any.Any, error) {
	// read the chunk which is a complete message
	chunk, err := readNextChunk(safeBuffer, TimeOut)
	if err != nil {
		return decode.SGsMessageType(0x00), &any.Any{}, err
	}

	// according to the type of the message, decide how to decode
	return message.SGsMessageDecoder(chunk)
}

// WriteChunk writes a chunk of bytes to the list
func (safeBuffer *SafeBuffer) WriteChunk(chunk []byte) (int, error) {
	safeBuffer.Lock()
	// write the chunk to the buffer
	safeBuffer.chunkBuffer.PushBack(chunk)
	// notify the reader if there is any reader waiting
	select {
	case safeBuffer.writeSignal <- true:
	default:
	}
	safeBuffer.Unlock()
	return len(chunk), nil
}

// readNextChunk tries to read the next chunk before timeout
func readNextChunk(safeBuffer *SafeBuffer, timeOutInSec int) ([]byte, error) {
	safeBuffer.Lock()
	// wait until there are sufficient bytes in the buffer
	for safeBuffer.chunkBuffer.Len() == 0 {
		// unlock the buffer and wait for incoming data
		safeBuffer.Unlock()
		select {
		case <-safeBuffer.writeSignal:
			// something is written to the buffer, lock the buffer again
			safeBuffer.Lock()
		case <-time.After(time.Duration(timeOutInSec) * time.Second):
			return []byte{}, errors.New("buffer read timeout")
		}
	}
	// now there are available chunk in the buffer which is locked by readBytes
	var chunk = safeBuffer.chunkBuffer.Front().Value
	safeBuffer.chunkBuffer.Remove(safeBuffer.chunkBuffer.Front())
	defer safeBuffer.Unlock()
	return chunk.([]byte), nil
}
