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

package parallel

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

const (
	DefaultNumWorkers = 10
)

// TODO(maxhbr): use generics after go update to 1.18
type Func func(In) (Out, error)
type In interface{}
type Out interface{}

// MapString is same as Map but for strings to strings.
func MapString(items []string, nWorkers int, f Func) ([]string, error) {
	var inputs []In
	for _, s := range items {
		inputs = append(inputs, s)
	}
	return MapInToString(inputs, nWorkers, f)
}

// MapInToString is same as Map but for []In to strings.
func MapInToString(inputs []In, nWorkers int, f Func) ([]string, error) {
	outsI, err := Map(inputs, nWorkers, f)
	if err != nil {
		return nil, err
	}
	var outs []string
	for _, s := range outsI {
		ss, ok := s.(string)
		if !ok {
			return nil, fmt.Errorf("could not convert returned item of type '%T' to string: '%+v'", s, s)
		}
		outs = append(outs, ss)
	}
	return outs, nil
}

// Map performs f on each element of items, with nWorkers in parallel.
// Out is in same order as items.
func Map(inputs []In, nWorkers int, f Func) ([]Out, error) {
	nJobs := len(inputs)
	jobs := make(chan workerIn, nJobs)
	outputs := make(chan workerOut, nJobs)

	// Workers
	work := func(ins chan workerIn, outs chan workerOut) {
		for in := range ins {
			out, err := f(in.input)
			outs <- workerOut{idx: in.idx, output: out, err: err}
		}
	}
	for i := 0; i < nWorkers; i++ {
		go work(jobs, outputs)
	}

	// Inputs
	for idx, input := range inputs {
		jobs <- workerIn{idx: idx, input: input}
	}

	// Outputs
	rets := make([]Out, nJobs)
	errs := &multierror.Error{}
	for i := 0; i < nJobs; i++ {
		ret := <-outputs
		rets[ret.idx] = ret.output
		errs = multierror.Append(errs, ret.err)
	}
	close(jobs)

	return rets, errs.ErrorOrNil()
}

type workerIn struct {
	idx   int
	input In
}

type workerOut struct {
	idx    int
	output Out
	err    error
}
