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

package dictionarygen

import (
	"bytes"
	"io"
	"strings"
	"unicode"
)

func p(w io.Writer, s ...string) {
	for _, v := range s {
		io.WriteString(w, v)
	}
	io.WriteString(w, "\n")
}

var firstCharacterReplacements = map[byte]string{
	'0': "Zero",
	'1': "One",
	'2': "Two",
	'3': "Three",
	'4': "Four",
	'5': "Five",
	'6': "Six",
	'7': "Seven",
	'8': "Eight",
	'9': "Nine",
}

func identifier(name string) string {
	if len(name) == 0 {
		return ""
	}

	if replacement, ok := firstCharacterReplacements[name[0]]; ok {
		name = replacement + name[1:]
	}

	fields := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
	var id bytes.Buffer
	for _, field := range fields {
		fieldUpper := strings.ToUpper(field)
		if commonInitialisms[fieldUpper] {
			id.WriteString(fieldUpper)
		} else {
			id.WriteString(strings.Title(field))
		}
	}
	return id.String()
}

// from golint
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}
