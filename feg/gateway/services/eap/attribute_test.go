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

package eap

import (
	"io"
	"testing"
)

var testEAP string = "\x01\x02\x00\xbc\x17\x01\x00\x00\x01\x05\x00\x00\x01\x23\x45\x67" +
	"\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef\x02\x05\x00\x00" +
	"\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c" +
	"\x81\x05\x00\x00\xd1\xef\x2a\xdf\x8a\xf9\x74\xf1\xe2\x5f\xac\x28" +
	"\x58\xbc\xe4\x9e\x82\x19\x00\x00\x22\x24\x46\x74\xd6\x10\x1b\x1e" +
	"\xd7\xc8\xfa\x8d\x8c\x43\x87\x37\xd0\x49\x72\xac\x8a\x7a\x28\x64" +
	"\xb6\x39\x20\xb0\x7c\x25\xc4\xbf\xd4\x69\x2e\x88\xe2\x18\xd9\xd6" +
	"\xdf\x20\xe3\x05\x94\x5c\x25\x97\x23\xd4\x6a\x59\x5b\xf7\x1b\x25" +
	"\x2e\x8a\x47\xe1\x45\x0f\xb2\x3f\x40\xc1\x1b\x22\xeb\xf3\x69\x86" +
	"\xd6\x61\xb1\xa9\x98\xf1\xb8\x16\x50\xe6\x5c\x73\xd5\x66\xf1\xea" +
	"\x31\xd6\x68\x5d\x87\x36\x7d\xb4\x0b\x05\x00\x00\x55\x1f\xec\x03" +
	"\xe0\xb1\xcc\x85\x31\x48\xb7\x5d\xf2\x57\x93\x65"

var expectedAttrs = []attribute{
	NewAttribute(1, []byte("\x00\x00\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef")),
	NewAttribute(2, []byte("\x00\x00\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c")),
	NewAttribute(129, []byte("\x00\x00\xd1\xef\x2a\xdf\x8a\xf9\x74\xf1\xe2\x5f\xac\x28\x58\xbc\xe4\x9e")),
	NewAttribute(130, []byte("\x00\x00\x22\x24\x46\x74\xd6\x10\x1b\x1e\xd7\xc8\xfa\x8d\x8c\x43"+
		"\x87\x37\xd0\x49\x72\xac\x8a\x7a\x28\x64\xb6\x39\x20\xb0\x7c\x25\xc4\xbf\xd4\x69\x2e\x88\xe2\x18\xd9\xd6"+
		"\xdf\x20\xe3\x05\x94\x5c\x25\x97\x23\xd4\x6a\x59\x5b\xf7\x1b\x25\x2e\x8a\x47\xe1\x45\x0f\xb2\x3f\x40\xc1"+
		"\x1b\x22\xeb\xf3\x69\x86\xd6\x61\xb1\xa9\x98\xf1\xb8\x16\x50\xe6\x5c\x73\xd5\x66\xf1\xea\x31\xd6\x68\x5d"+
		"\x87\x36\x7d\xb4")),
	NewAttribute(11, []byte("\x00\x00\x55\x1f\xec\x03\xe0\xb1\xcc\x85\x31\x48\xb7\x5d\xf2\x57\x93\x65")),
}

func TestAttributeScanner(t *testing.T) {
	a1 := NewAttribute(11, []byte("\x00\x00\x55\x1f\xec\x03\xe0\xb1\xcc\x85\x31\x48\xb7\x5d\xf2\x57\x93\x65"))
	a2 := NewAttribute(11, []byte("\x00\x00\x55\x1f\xec\x03\xe0\xb1\xcc\x85\x31\x48\xb7\x5d\xf2\x57\x93\x65\x00"))
	if a2.Len() != a1.Len()+4 {
		t.Fatalf("Invalid Attr padding %d != %d + 4", a2.Len(), a1.Len())
	}
	if a1.Type() != 11 {
		t.Fatalf("Invalid Attr Type %d != %d", a1.Type(), 11)
	}
	if a2.Type() != 11 {
		t.Fatalf("Invalid Attr Type %d != %d", a2.Type(), 11)
	}
	if string(a2.Value()) != "\x00\x00\x55\x1f\xec\x03\xe0\xb1\xcc\x85\x31\x48\xb7\x5d\xf2\x57\x93\x65\x00\x00\x00\x00" {
		t.Fatalf("EAP Attr Value Mismatch: %v", a2.Value())
	}

	scanner, err := NewAttributeScanner([]byte(testEAP))
	if err != nil {
		t.Fatal(err)
	}
	if scanner == nil {
		t.Fatal("Nil Attribute Scanner")
	}
	var (
		attr       Attribute
		attributes []Attribute
	)
	i := 0
	for ; ; i++ {
		attr, err = scanner.Next()
		if err != nil {
			break
		}
		if attr == nil {
			t.Fatal("Nil Attribute")
		}
		if attr.Type() != expectedAttrs[i].Type() {
			t.Fatalf("EAP Attr Type Mismatch for attr #%d: expected %d got %d", i, expectedAttrs[i].Type(), attr.Type())
		}
		if string(attr.Value()) != string(expectedAttrs[i].Value()) {
			t.Fatalf("EAP Attr Value Mismatch for attr #%d: expected %v got %v", i, expectedAttrs[i].Value(), attr.Value())
		}
		attributes = append(attributes, attr)
	}
	if err != io.EOF {
		t.Fatal(err)
	}
	if i != len(expectedAttrs) {
		t.Fatalf("EAP Attrributes # Mismatch: expected %d got %d", len(expectedAttrs), i)
	}
	// Check that reset works
	scanner.Reset()
	attr, err = scanner.Next()
	if err != nil {
		t.Fatal(err)
	}
	if attr.Type() != 1 {
		t.Fatalf("EAP Attr Type Mismatch: expected 1 got %d", attr.Type())
	}
	expected := "\x00\x00\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef"
	if string(attr.Value()) != expected {
		t.Fatalf("EAP Attr Value Mismatch: expected %v got %v", []byte(expected), attr.Value())
	}
	// build EAP with reserved capacity (optimization)
	p := NewPacket(1, 2, []byte{23, 1, 0, 0}, 180)
	for _, a := range attributes {
		p, err = p.Append(a)
		if err != nil {
			t.Fatal(err)
		}
	}
	if testEAP != string(p) {
		t.Fatalf("EAP Mismatch 1\nexpected: %v\n     got: %v", []byte(testEAP), p)
	}
	// build EAP without reserved capacity
	p = NewPacket(1, 2, []byte{23, 1, 0, 0})
	for _, a := range attributes {
		p, err = p.Append(a)
		if err != nil {
			t.Fatal(err)
		}
	}
	if testEAP != string(p) {
		t.Fatalf("EAP Mismatch 2\nexpected: %v\n     got: %v", []byte(testEAP), p)
	}
}
