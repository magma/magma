/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package multiplex allows to create different algorithms to select out of an index base on some
//	parameters. An example of a multiplexor can be seen on StaticMultiplexByIMSI type.
//	Multiplexors in this packages use Context tyoe to obtain dynamic parameters. Right now we have
//	only identified one parameter to be used by algorithms (IMSI). But in the future more parameter
//	can be added in case other algorithms uses them
//
//	Addition of new algorithms:
//	- if you need new dynamic (per call) parameters:
//	  > Add them to Context.
//	  > Create "With..." receivers like the one below in order to chain the calls
//	    def:     func (mp *Context) WithNewParam(newParam Type) *Context{}
//	    usage:   muxCtx := NewContext.WithIMSI("1234").WithNewParam(newParam)
//	  > Modify the context calls on magma to make sure that parameter is included where you need it
//	- Implement multiplexor interface with your new algorithm. HEre you can add some static data when
//	  you create the multiplexor (an example could be a list of predefined imsis)

package multiplex

import (
	"fmt"

	"magma/lte/cloud/go/protos"
)

// Multiplexor interface has to be implemented by any new Multiplexor.
// Multiplexor can be loaded with data during creation and we can pass some dynamic parameters
// coming from each requests (using context)
type Multiplexor interface {
	// GetIndex return an index based on the implementation of Multiplexor
	GetIndex(*Context) (int, error)
}

// Context is a type used as a way to pass dynamic parameters coming from specific calls during
// execution of service (for example CreateSessionRequest)
type Context struct {
	imsiNumeric uint64
	lastError   error
}

// NewContext creates a new context
func NewContext() *Context {
	return &Context{}
}

// GetContext returns IMSI in uint64 format
// That IMSI can come from different sources like IMSI string or SessionId
func (c *Context) GetIMSI() (uint64, error) {
	if c.lastError != nil {
		return 0, c.lastError
	}
	return c.imsiNumeric, nil
}

func (c *Context) GetError() error {
	return c.lastError
}

// WithIMSI adds imsi to context from a IMSI string (from IMSI123456789012345   123456789012345)
func (c *Context) WithIMSI(imsi string) *Context {
	if c == nil {
		c = &Context{}
	}
	if c.lastError != nil {
		return c
	}
	_, imsiNumeric, err := protos.StripPrefixFromIMSIandFormat(imsi)
	if err != nil {
		c.lastError = err
		return c
	}
	c.imsiNumeric = imsiNumeric
	return c
}

// WithSessionId adds imsi to context from sessionId (from IMSI123456789012345-54321 to 123456789012345)
func (c *Context) WithSessionId(sessionId string) *Context {
	if c == nil {
		c = &Context{}
	}
	if c.lastError != nil {
		return c
	}
	imsiWithPrefix, err := protos.GetIMSIwithPrefixFromSessionId(sessionId)
	if err != nil {
		c.lastError = err
		return c
	}
	return c.WithIMSI(imsiWithPrefix)
}

// StaticMultiplexByIMSI contains is a basic Multiplexor that distribuites each IMSI to a specific
// index based on the IMSImod(numServers)
type StaticMultiplexByIMSI struct {
	totalServers uint64
}

// NewStaticMultiplexByIMSI creates a StaticMultiplexByIMSI with a specific number of servers
// That number must be specified before service is started. Context must be used to pass dynamic
// data coming during execution of service
func NewStaticMultiplexByIMSI(numServers int) (Multiplexor, error) {
	if numServers < 1 {
		return nil, fmt.Errorf("MultiplexByIMSI needs to be configured with 1 or more than 1 servers (%d configured)", numServers)
	}
	return &StaticMultiplexByIMSI{uint64(numServers)}, nil
}

// GetIndex provides the index of the server per that IMSI
func (m *StaticMultiplexByIMSI) GetIndex(c *Context) (int, error) {
	if c.lastError != nil {
		return -1, c.lastError
	}
	index := int(c.imsiNumeric % m.totalServers)
	return index, nil
}
