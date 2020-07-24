//+build ignore
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

package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"fbc/lib/go/radius/dictionary"
	"fbc/lib/go/radius/dictionarygen"
)

func main() {
	resp, err := http.Get(`https://support.arubanetworks.com/ToolsResources/tabid/76/DMXModule/514/Command/Core_Download/Default.aspx?EntryId=156`)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	parser := dictionary.Parser{
		Opener: restrictedOpener{
			"main.dictionary": body,
		},
	}
	dict, err := parser.ParseFile("main.dictionary")
	if err != nil {
		log.Fatal(err)
	}

	gen := dictionarygen.Generator{
		Package: "aruba",
	}
	generated, err := gen.Generate(dict)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("generated.go", generated, 0644); err != nil {
		log.Fatal(err)
	}
}

type restrictedOpener map[string][]byte

func (r restrictedOpener) OpenFile(name string) (dictionary.File, error) {
	contents, ok := r[name]
	if !ok {
		return nil, errors.New("unknown file " + name)
	}
	return &restrictedFile{
		Reader:    bytes.NewReader(contents),
		NameValue: name,
	}, nil
}

type restrictedFile struct {
	io.Reader
	NameValue string
}

func (r *restrictedFile) Name() string {
	return r.NameValue
}

func (r *restrictedFile) Close() error {
	return nil
}
