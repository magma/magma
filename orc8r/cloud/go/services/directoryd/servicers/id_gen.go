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
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/golang/glog"

	magmaerrors "magma/orc8r/lib/go/errors"
)

const (
	DEFAULTMAXATTEMPTS = 1000
)

// IdGenerator provides with Unique IDs using a rando generator and making sure the id is not
// currently being used by the database
type IdGenerator struct {
	mu          sync.Mutex
	maxAttempts int
}

// doesExistInDatabaseFunc defines a function to find if the index really exist in the database
type doesExistInDatabaseFunc func(networkId string, id string) (string, error)

func NewIdGeneratorWithAttempts(attempts int) *IdGenerator {
	return &IdGenerator{maxAttempts: attempts}
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{maxAttempts: DEFAULTMAXATTEMPTS}
}

// GetUniqueUint32Id finds a unique Uint32 ID making sure it does not exist on the database
// It uses doesExistInDatabaseFunc to make sure the random id is not being used
func (g *IdGenerator) GetUniqueUint32Id(network string, doesExistFunc doesExistInDatabaseFunc) (uint32, error) {
	for i := 0; i < g.maxAttempts; i++ {
		newIDuint32 := g.getRandomUintFromOneToMaxUint32()
		_, err := doesExistFunc(network, fmt.Sprint(newIDuint32))
		if err == magmaerrors.ErrNotFound {
			return newIDuint32, nil
		}
		if err != nil {
			glog.Errorf("GetNewSgwCTeid could not get unique TEID: %s", err)
		}
	}
	return 0, fmt.Errorf("GetNewSgwCTeid couldnt get a unique teid after MAXATTEMPTS (%d)", g.maxAttempts)
}

func (g *IdGenerator) getRandomUintFromOneToMaxUint32() uint32 {
	randRes := g.getNonZeroRandomInt(big.NewInt(math.MaxUint32)).Int64()
	if randRes > math.MaxUint32 {
		panic(fmt.Errorf("getRandomUintFromOneToMax returned a value bigger than MaxUint32"))
	}
	return uint32(randRes)
}

// getNonZeroRandomInt returns a uniform random value in (0, max). It panics if max <= 0.
func (g *IdGenerator) getNonZeroRandomInt(nMax *big.Int) *big.Int {
	var (
		n   *big.Int
		err error
	)
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := 0; i < 2; i++ {
		n, err = rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
		if err != nil {
			panic(err)
		}
		if n.Int64() != 0 {
			break
		}
	}
	return n
}
