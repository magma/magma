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

package debug

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/dictionary"
)

type Config struct {
	Dictionary *dictionary.Dictionary
}

func Dump(w io.Writer, c *Config, p *radius.Packet) {
	io.WriteString(w, p.Code.String())
	io.WriteString(w, " Id ")
	io.WriteString(w, strconv.Itoa(int(p.Identifier)))
	io.WriteString(w, "\n")
	dumpAttrs(w, c, p)
}

func DumpString(c *Config, p *radius.Packet) string {
	var b bytes.Buffer
	Dump(&b, c, p)
	b.Truncate(b.Len() - 1) // remove trailing \n
	return b.String()
}

func DumpRequest(w io.Writer, c *Config, req *radius.Request) {
	io.WriteString(w, req.Code.String())
	io.WriteString(w, " Id ")
	io.WriteString(w, strconv.Itoa(int(req.Identifier)))
	io.WriteString(w, " from ")
	io.WriteString(w, req.RemoteAddr.String())
	io.WriteString(w, " to ")
	io.WriteString(w, req.LocalAddr.String())
	io.WriteString(w, "\n")
	dumpAttrs(w, c, req.Packet)
}

func DumpRequestString(c *Config, req *radius.Request) string {
	var b bytes.Buffer
	DumpRequest(&b, c, req)
	b.Truncate(b.Len() - 1) // remove trailing \n
	return b.String()
}

func dumpAttrs(w io.Writer, c *Config, p *radius.Packet) {
	for _, elem := range sortedAttributes(p.Attributes) {
		attrsType, attrs := elem.Type, elem.Attrs
		if len(attrs) == 0 {
			continue
		}

		for _, attr := range attrs {

			var attrTypeStr string
			var attrsTypeIntStr string
			var attrStr string

			searchAttrs := c.Dictionary.Attributes
			searchValues := c.Dictionary.Values

			attrsTypeIntStr = strconv.Itoa(int(attrsType))

			dictAttr := dictionary.AttributeByOID(searchAttrs, attrsTypeIntStr)
			if dictAttr != nil {
				attrTypeStr = dictAttr.Name
				switch dictAttr.Type {
				case dictionary.AttributeString, dictionary.AttributeOctets:
					if dictAttr != nil && dictAttr.FlagEncrypt != nil && *dictAttr.FlagEncrypt == 1 {
						decryptedValue, err := radius.UserPassword(radius.Attribute(attr), p.Secret, p.Authenticator[:])
						if err == nil {
							attrStr = fmt.Sprintf("%q", decryptedValue)
							break
						}
					}
					attrStr = fmt.Sprintf("%q", attr)

				case dictionary.AttributeDate:
					if len(attr) == 4 {
						t := time.Unix(int64(binary.BigEndian.Uint32(attr)), 0).UTC()
						attrStr = t.Format(time.RFC3339)
					}

				case dictionary.AttributeInteger:
					switch len(attr) {
					case 4:
						intVal := int(binary.BigEndian.Uint32(attr))
						if dictAttr != nil {
							var matchedNames []string
							for _, value := range dictionary.ValuesByAttribute(searchValues, dictAttr.Name) {
								if value.Number == intVal {
									matchedNames = append(matchedNames, value.Name)
								}
							}
							if len(matchedNames) > 0 {
								sort.Stable(sort.StringSlice(matchedNames))
								attrStr = strings.Join(matchedNames, " / ")
								break
							}
						}
						attrStr = strconv.Itoa(intVal)
					case 8:
						attrStr = strconv.Itoa(int(binary.BigEndian.Uint64(attr)))
					}

				case dictionary.AttributeIPAddr, dictionary.AttributeIPv6Addr:
					switch len(attr) {
					case net.IPv4len, net.IPv6len:
						attrStr = net.IP(attr).String()
					}

				case dictionary.AttributeIFID:
					if len(attr) == 8 {
						attrStr = net.HardwareAddr(attr).String()
					}

				}
			} else {
				attrTypeStr = "#" + attrsTypeIntStr
			}

			if len(attrStr) == 0 {
				attrStr = "0x" + hex.EncodeToString(attr)
			}

			io.WriteString(w, "  ")
			io.WriteString(w, attrTypeStr)
			io.WriteString(w, " = ")
			io.WriteString(w, attrStr)
			io.WriteString(w, "\n")
		}
	}
}

type attributesElement struct {
	Type  radius.Type
	Attrs []radius.Attribute
}

func sortedAttributes(attributes radius.Attributes) []attributesElement {
	var sortedAttrs []attributesElement
	for attrsType, attrs := range attributes {
		sortedAttrs = append(sortedAttrs, attributesElement{
			Type:  attrsType,
			Attrs: attrs,
		})
	}

	sort.Sort(sortAttributesType(sortedAttrs))

	return sortedAttrs
}

type sortAttributesType []attributesElement

func (s sortAttributesType) Len() int           { return len(s) }
func (s sortAttributesType) Less(i, j int) bool { return s[i].Type < s[j].Type }
func (s sortAttributesType) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
