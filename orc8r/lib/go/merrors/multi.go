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

package errors

import (
	"bytes"
	"fmt"
)

// Multi error implementation
type Multi struct {
	errors []error
}

// NewMulti returns new Multi error populated with errs, if any
func NewMulti(errs ...error) *Multi {
	var me *Multi
	return me.Set(errs...)
}

// Error returns a formatted string for Multi error list
func (me *Multi) Error() string {
	if me == nil {
		return "<nil>"
	}
	switch len(me.errors) {
	case 0:
		return ""
	case 1:
		return me.errors[0].Error()
	default:
		var b bytes.Buffer
		fmtStr := "errors: [%d: %v"
		for i, e := range me.errors {
			fmt.Fprintf(&b, fmtStr, i, e)
			fmtStr = "; %d: %v"
		}
		b.Write([]byte("]"))
		return b.String()
	}
}

// Set - sets multi set to errs
func (me *Multi) Set(errs ...error) *Multi {
	var errors []error
	for _, e := range errs {
		if e != nil {
			errors = append(errors, e)
		}
	}
	if len(errors) == 0 {
		return me
	}
	if me == nil {
		return &Multi{errors: errors}
	}
	me.errors = errors
	return me
}

// Get - returns a list of errors encapsulated by the Multi error
func (me *Multi) Get() []error {
	if me != nil {
		return me.errors
	}
	return nil
}

// Add appends errs to the existing MultiError set
func (me *Multi) Add(errs ...error) *Multi {
	if me == nil {
		return me.Set(errs...)
	}
	for _, e := range errs {
		if e != nil {
			me.errors = append(me.errors, e)
		}
	}
	return me
}

// AsError returns error cast of Multi,
// the returned error is safe to use in any std error evaluations, such as if err == nil, etc.
// Functions returning error should always return Multi.AsError() instead of Multi directly
func (me *Multi) AsError() error {
	if len(me.Get()) == 0 { // nil me or empty errors list is equivalent to no error
		return nil
	}
	return me
}

// AddFmt adds a new formatted error if err is not nil, it's a noop if err == nil & returns unchanged 'me' in this case
func (me *Multi) AddFmt(err error, format string, args ...interface{}) *Multi {
	if err == nil {
		return me
	}
	return me.Add(fmt.Errorf(format+fmt.Sprintf(" %%[%d]v", len(args)+1), append(args, err)...))
}

// Cast casts Multi error (if any) to error and returns it
// if err is already nil or not Multi type Cast will just return it
// the returned error is safe to use in any std error evaluations, such as if err == nil, etc.
func Cast(err error) error {
	if err != nil {
		if me, ok := err.(*Multi); ok {
			return me.AsError()
		}
	}
	return err
}
