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

package serde_test

import (
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serde/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSerialize(t *testing.T) {
	s := getSerde("type0", "hello world")
	registry := serde.NewRegistry(s)

	actual, err := serde.Serialize("some_val", "type0", registry)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), actual)

	_, err = serde.Serialize("some_val", "typeXXX", registry)
	assert.EqualError(t, err, "no serde in registry for type typeXXX")
}

func TestDeserialize(t *testing.T) {
	s := getSerde("type0", "hello world")
	registry := serde.NewRegistry(s)

	actual, err := serde.Deserialize([]byte("some_val"), "type0", registry)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", actual)

	_, err = serde.Serialize("some_val", "typeXXX", registry)
	assert.EqualError(t, err, "no serde in registry for type typeXXX")
}

func TestRegistryIface(t *testing.T) {
	s0 := getSerde("type0", "val0")
	s1 := getSerde("type1", "val1")
	s2 := getSerde("type2", "val2")
	s3 := getSerde("type3", "val3")
	sX := getSerde("type3", "valXXX")

	reg0 := serde.NewRegistry(s0, s1)

	actualDes, err := serde.Deserialize([]byte("some_val"), "type0", reg0)
	assert.NoError(t, err)
	assert.Equal(t, "val0", actualDes)
	actualSer, err := serde.Serialize("some_val", "type1", reg0)
	assert.NoError(t, err)
	assert.Equal(t, []byte("val1"), actualSer)

	r := serde.NewRegistry(s2, s3)
	reg1 := r.MustMerge(reg0)

	actualDes, err = serde.Deserialize([]byte("some_val"), "type0", reg1)
	assert.NoError(t, err)
	assert.Equal(t, "val0", actualDes)
	actualSer, err = serde.Serialize("some_val", "type3", reg1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("val3"), actualSer)
	_, err = serde.Serialize("some_val", "typeXXX", reg1)
	assert.EqualError(t, err, "no serde in registry for type typeXXX")

	r = serde.NewRegistry(sX)
	assert.Panics(t, func() { r.MustMerge(reg1) })
}

func getSerde(typ, ret string) serde.Serde {
	s := &mocks.Serde{}
	s.On("GetDomain").Return("domain0")
	s.On("GetType").Return(typ)
	s.On("Serialize", mock.Anything).Return([]byte(ret), nil)
	s.On("Deserialize", mock.Anything).Return(ret, nil)
	return s
}
