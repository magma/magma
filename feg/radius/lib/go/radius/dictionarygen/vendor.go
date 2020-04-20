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
	"io"

	"fbc/lib/go/radius/dictionary"
)

func (g *Generator) genVendor(w io.Writer, vendor *dictionary.Vendor) {
	ident := identifier(vendor.Name)

	p(w)
	p(w, `func _`, ident, `_AddVendor(p *radius.Packet, typ byte, attr radius.Attribute) (err error) {`)
	p(w, `	var vsa radius.Attribute`)
	p(w, `	vendor := make(radius.Attribute, 2+len(attr))`)
	p(w, `	vendor[0] = typ`)
	p(w, `	vendor[1] = byte(len(vendor))`)
	p(w, `	copy(vendor[2:], attr)`)
	p(w, `	vsa, err = radius.NewVendorSpecific(_`, ident, `_VendorID, vendor)`)
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	p(w, `	p.Add(rfc2865.VendorSpecific_Type, vsa)`)
	p(w, `	return nil`)
	p(w, `}`)

	p(w)
	p(w, `func _`, ident, `_GetsVendor(p *radius.Packet, typ byte) (values []radius.Attribute) {`)
	p(w, `	for _, attr := range p.Attributes[rfc2865.VendorSpecific_Type] {`)
	p(w, `		vendorID, vsa, err := radius.VendorSpecific(attr)`)
	p(w, `		if err != nil || vendorID != _`, ident, `_VendorID {`)
	p(w, `			continue`)
	p(w, `		}`)
	p(w, `		for len(vsa) >= 3 {`)
	p(w, `			vsaTyp, vsaLen := vsa[0], vsa[1]`)
	p(w, `			if int(vsaLen) > len(vsa) || vsaLen < 3 {`) // malformed
	p(w, `				break`)
	p(w, `			}`)
	p(w, `			if vsaTyp == typ {`)
	p(w, `				values = append(values, vsa[2:int(vsaLen)])`)
	p(w, `			}`)
	p(w, `			vsa = vsa[int(vsaLen):]`)
	p(w, `		}`)
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	p(w, `func _`, ident, `_LookupVendor(p *radius.Packet, typ byte) (attr radius.Attribute, ok bool) {`)
	p(w, `	for _, a := range p.Attributes[rfc2865.VendorSpecific_Type] {`)
	p(w, `		vendorID, vsa, err := radius.VendorSpecific(a)`)
	p(w, `		if err != nil || vendorID != _`, ident, `_VendorID {`)
	p(w, `			continue`)
	p(w, `		}`)
	p(w, `		for len(vsa) >= 3 {`)
	p(w, `			vsaTyp, vsaLen := vsa[0], vsa[1]`)
	p(w, `			if int(vsaLen) > len(vsa) || vsaLen < 3 {`) // malformed
	p(w, `				break`)
	p(w, `			}`)
	p(w, `			if vsaTyp == typ {`)
	p(w, `				return vsa[2:int(vsaLen)], true`)
	p(w, `			}`)
	p(w, `			vsa = vsa[int(vsaLen):]`)
	p(w, `		}`)
	p(w, `	}`)
	p(w, `	return nil, false`)
	p(w, `}`)

	p(w)
	p(w, `func _`, ident, `_SetVendor(p *radius.Packet, typ byte, attr radius.Attribute) (err error) {`)
	p(w, `	for i := 0; i < len(p.Attributes[rfc2865.VendorSpecific_Type]); {`)
	p(w, `		vendorID, vsa, err := radius.VendorSpecific(p.Attributes[rfc2865.VendorSpecific_Type][i])`)
	p(w, `		if err != nil || vendorID != _`, ident, `_VendorID {`)
	p(w, `			i++`)
	p(w, `			continue`)
	p(w, `		}`)
	p(w, `		for j := 0; len(vsa[j:]) >= 3; {`)
	p(w, `			vsaTyp, vsaLen := vsa[0], vsa[1]`)
	p(w, `			if int(vsaLen) > len(vsa[j:]) || vsaLen < 3 {`) // malformed
	p(w, `				i++`)
	p(w, `				break`)
	p(w, `			}`)
	p(w, `			if vsaTyp == typ {`)
	p(w, `				vsa = append(vsa[:j], vsa[j+int(vsaLen):]...)`)
	p(w, `			}`)
	p(w, `			j += int(vsaLen)`)
	p(w, `		}`)
	p(w, `		if len(vsa) > 0 {`)
	p(w, `			copy(p.Attributes[rfc2865.VendorSpecific_Type][i][4:], vsa)`)
	p(w, `			i++`)
	p(w, `		} else {`)
	p(w, `			p.Attributes[rfc2865.VendorSpecific_Type] = append(p.Attributes[rfc2865.VendorSpecific_Type][:i], p.Attributes[rfc2865.VendorSpecific_Type][i+i:]...)`)
	p(w, `		}`)
	p(w, `	}`)
	p(w, `	return _`, ident, `_AddVendor(p, typ, attr)`)
	p(w, `}`)
}
