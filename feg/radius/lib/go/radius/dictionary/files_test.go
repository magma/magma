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

package dictionary_test

import (
	"errors"
	"io"
	"strings"

	dict "fbc/lib/go/radius/dictionary"
)

type MemoryFile struct {
	Filename string
	Contents string

	r io.Reader
}

func (m *MemoryFile) Read(p []byte) (n int, err error) {
	if m.r == nil {
		m.r = strings.NewReader(m.Contents)
	}
	return m.r.Read(p)
}

func (m *MemoryFile) Close() error {
	return nil
}

func (m *MemoryFile) Name() string {
	return m.Filename
}

type MemoryOpener []MemoryFile

func (m MemoryOpener) OpenFile(name string) (dict.File, error) {
	for _, file := range m {
		if file.Filename == name {
			return &file, nil
		}
	}
	return nil, errors.New("unknown file " + name)
}

var files = MemoryOpener{
	{
		Filename: "simple.dict",
		Contents: `
ATTRIBUTE User-Name 1 string
ATTRIBUTE User-Password 2 octets encrypt=1

ATTRIBUTE Mode 127 integer
VALUE Mode Full 1
VALUE Mode Half 2

ATTRIBUTE ARAP-Challenge-Response 84 octets[8]
`,
	},

	{
		Filename: "recursive_1.dict",
		Contents: `
$INCLUDE recursive_2.dict
`,
	},
	{
		Filename: "recursive_2.dict",
		Contents: `
$INCLUDE recursive_1.dict
`,
	},
	{
		Filename: "tlv.dict",
		Contents: `
ATTRIBUTE  Struct-Name  4  tlv
ATTRIBUTE  Field1       4.1 string
ATTRIBUTE  Field2       4.2 integer64
`,
	},
}
