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

package dictionary

import "fmt"

func Merge(d1, d2 *Dictionary) (*Dictionary, error) {
	for _, attr := range d2.Attributes {
		existingAttr := AttributeByName(d1.Attributes, attr.Name)
		if existingAttr == nil {
			existingAttr = AttributeByOID(d1.Attributes, attr.OID)
		}

		if existingAttr != nil {
			return nil, fmt.Errorf("duplicate attribute %s (%s)", attr.Name, attr.OID)
		}
	}

	for _, vendor := range d2.Vendors {
		existingVendorByName := VendorByName(d1.Vendors, vendor.Name)
		existingVendorByNumber := VendorByNumber(d1.Vendors, vendor.Number)
		if existingVendorByName != existingVendorByNumber {
			// TODO: make sure vendor flags, etc. match?
			return nil, fmt.Errorf("conflicting vendor: %s (%d)", vendor.Name, vendor.Number)
		}
		if existingVendorByName == nil {
			continue
		}

		for _, attr := range vendor.Attributes {
			existingAttr := AttributeByName(existingVendorByName.Attributes, attr.Name)
			if existingAttr == nil {
				existingAttr = AttributeByOID(existingVendorByName.Attributes, attr.OID)
			}

			if existingAttr != nil {
				return nil, fmt.Errorf("duplicate vendor attrbute %s (%s)", attr.Name, attr.OID)
			}
		}
	}

	newDict := new(Dictionary)

	if size := len(d1.Attributes) + len(d2.Attributes); size > 0 {
		newDict.Attributes = make([]*Attribute, 0, len(d1.Attributes)+len(d2.Attributes))
		newDict.Attributes = append(newDict.Attributes, d1.Attributes...)
		newDict.Attributes = append(newDict.Attributes, d2.Attributes...)
	}

	if size := len(d1.Values) + len(d2.Values); size > 0 {
		newDict.Values = make([]*Value, 0, len(d1.Values)+len(d2.Values))
		newDict.Values = append(newDict.Values, d1.Values...)
		newDict.Values = append(newDict.Values, d2.Values...)
	}

	if size := len(d1.Vendors) + len(d2.Vendors); size > 0 {
		newDict.Vendors = make([]*Vendor, 0, len(d1.Vendors)+len(d2.Vendors))
		newDict.Vendors = append(newDict.Vendors, d1.Vendors...)
		for _, vendor := range d2.Vendors {
			existingVendor := VendorByNumber(newDict.Vendors, vendor.Number)
			if existingVendor != nil {
				existingVendor.Attributes = append(existingVendor.Attributes, vendor.Attributes...)
				existingVendor.Values = append(existingVendor.Values, vendor.Values...)
			} else {
				newDict.Vendors = append(newDict.Vendors, vendor)
			}
		}
	}

	return newDict, nil
}

func AttributeByName(attrs []*Attribute, name string) *Attribute {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr
		}
	}
	return nil
}

func AttributeByOID(attrs []*Attribute, oid string) *Attribute {
	for _, attr := range attrs {
		if attr.OID == oid {
			return attr
		}
	}
	return nil
}

func ValuesByAttribute(values []*Value, attribute string) []*Value {
	var matched []*Value
	for _, value := range values {
		if value.Attribute == attribute {
			matched = append(matched, value)
		}
	}
	return matched
}

func VendorByName(vendors []*Vendor, name string) *Vendor {
	for _, vendor := range vendors {
		if vendor.Name == name {
			return vendor
		}
	}
	return nil
}

func VendorByNumber(vendors []*Vendor, number int) *Vendor {
	for _, vendor := range vendors {
		if vendor.Number == number {
			return vendor
		}
	}
	return nil
}

func vendorByNameOrNumber(vendors []*Vendor, name string, number int) *Vendor {
	for _, vendor := range vendors {
		if vendor.Name == name || vendor.Number == number {
			return vendor
		}
	}
	return nil
}
