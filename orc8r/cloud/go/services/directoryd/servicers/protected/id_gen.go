/*
Copyright 2021 The Magma Authors.

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
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/glog"

	"magma/orc8r/lib/go/merrors"
)

const (
	DEFAULTMAXATTEMPTS = 1000
)

// IdGenerator provides with Unique IDs using a rando generator and making sure the id is not
// currently being used by the database
type IdGenerator struct {
	mu          sync.Mutex
	rand        *rand.Rand
	maxAttempts int
}

// doesExistInDatabaseFunc defines a function to find if the index really exist in the database
type doesExistInDatabaseFunc func(networkId string, id string) (string, error)

func NewIdGeneratorWithAttempts(attempts int) *IdGenerator {
	return &IdGenerator{
		maxAttempts: attempts,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func NewIdGenerator() *IdGenerator {
	return NewIdGeneratorWithAttempts(DEFAULTMAXATTEMPTS)
}

// GetUniqueId finds a unique Uint32 ID making sure it does not exist on the database
// It uses doesExistInDatabaseFunc to make sure the random id is not being used
func (g *IdGenerator) GetUniqueId(network string, doesExistFunc doesExistInDatabaseFunc) (uint32, error) {
	for i := 0; i < g.maxAttempts; i++ {
		newID := g.getRandomFromOneToMaxUint32()
		_, err := doesExistFunc(network, fmt.Sprint(newID))
		if err == merrors.ErrNotFound {
			return newID, nil
		}
		if err != nil {
			glog.Errorf("GetNewSgwCTeid could not get unique TEID: %s", err)
		}
	}
	return 0, fmt.Errorf("GetNewSgwCTeid couldnt get a unique teid after MAXATTEMPTS (%d)", g.maxAttempts)
}

func (g *IdGenerator) getRandomFromOneToMaxUint32() uint32 {
	r := GetRandomInt63(g.rand, 1, math.MaxUint32)
	return uint32(r)
}

// GetRandomInt63 generates a random int64 between min and max
func GetRandomInt63(r *rand.Rand, min int64, max int64) int64 {
	if min > max {
		panic(fmt.Sprintf("GetUniqueId got overlapped arguments min(%d)>max(%d)", min, max))
	}
	return r.Int63n(max-min) + min
}
