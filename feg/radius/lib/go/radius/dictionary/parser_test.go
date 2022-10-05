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
	"bytes"
	"fmt"
	"reflect"
	"testing"

	dict "fbc/lib/go/radius/dictionary"
)

func TestParser(t *testing.T) {
	parser := dict.Parser{
		Opener: files,
	}

	d, err := parser.ParseFile("simple.dict")
	if err != nil {
		t.Fatal(err)
	}

	expected := &dict.Dictionary{
		Attributes: []*dict.Attribute{
			{
				Name: "User-Name",
				OID:  "1",
				Type: dict.AttributeString,
			},
			{
				Name:        "User-Password",
				OID:         "2",
				Type:        dict.AttributeOctets,
				FlagEncrypt: newIntPtr(1),
			},
			{
				Name: "Mode",
				OID:  "127",
				Type: dict.AttributeInteger,
			},
			{
				Name: "ARAP-Challenge-Response",
				OID:  "84",
				Type: dict.AttributeOctets,
				Size: newIntPtr(8),
			},
		},
		Values: []*dict.Value{
			{
				Attribute: "Mode",
				Name:      "Full",
				Number:    1,
			},
			{
				Attribute: "Mode",
				Name:      "Half",
				Number:    2,
			},
		},
	}

	if !reflect.DeepEqual(d, expected) {
		t.Fatalf("got %s, expected %s", dictString(d), dictString(expected))
	}
}

func TestParser_recursiveinclude(t *testing.T) {
	parser := dict.Parser{
		Opener: files,
	}

	d, err := parser.ParseFile("recursive_1.dict")
	pErr, ok := err.(*dict.ParseError)
	if !ok || pErr == nil || d != nil {
		t.Fatalf("got %v, expected *ParseError", pErr)
	}
	if _, ok := pErr.Inner.(*dict.RecursiveIncludeError); !ok {
		t.Fatalf("got %v, expected *RecursiveIncludeError", pErr.Inner)
	}
}

func TestParser_TLV(t *testing.T) {
	parser := dict.Parser{
		Opener: files,
	}

	d, err := parser.ParseFile("tlv.dict")
	if err != nil {
		t.Fatal(err)
	}

	expected := &dict.Dictionary{
		Attributes: []*dict.Attribute{
			{
				Name: "Struct-Name",
				OID:  "4",
				Type: dict.AttributeTLV,
				Attributes: []*dict.Attribute{
					{
						Name: "Field1",
						OID:  "1",
						Type: dict.AttributeString,
					},
					{
						Name: "Field2",
						OID:  "2",
						Type: dict.AttributeInteger64,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(d, expected) {
		t.Fatalf("got %s expected %s", dictString(d), dictString(expected))
	}
}
func newIntPtr(i int) *int {
	return &i
}

func dictString(d *dict.Dictionary) string {
	var b bytes.Buffer
	b.WriteString("dictionary.Dictionary\n")

	b.WriteString("\tAttributes:\n")
	for _, attr := range d.Attributes {
		b.WriteString(fmt.Sprintf("\t\t%q %q %q %#v %#v\n", attr.Name, attr.OID, attr.Type, attr.FlagHasTag, attr.FlagEncrypt))
	}

	b.WriteString("\tValues:\n")
	for _, value := range d.Values {
		b.WriteString(fmt.Sprintf("\t\t%q %q %d\n", value.Attribute, value.Name, value.Number))
	}

	b.WriteString("\tVendors:\n")
	for _, vendor := range d.Vendors {
		b.WriteString(fmt.Sprintf("\t\t%q %d\n", vendor.Name, vendor.Number))

		b.WriteString("\t\tAttributes:\n")
		for _, attr := range vendor.Attributes {
			b.WriteString(fmt.Sprintf("\t\t%q %q %q %#v %#v\n", attr.Name, attr.OID, attr.Type, attr.FlagHasTag, attr.FlagEncrypt))
		}

		b.WriteString("\t\tValues:\n")
		for _, value := range vendor.Values {
			b.WriteString(fmt.Sprintf("\t\t%q %q %d\n", value.Attribute, value.Name, value.Number))
		}
	}

	return b.String()
}
