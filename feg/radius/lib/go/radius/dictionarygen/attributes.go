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
	"net"
	"strconv"
	"strings"

	"fbc/lib/go/radius/dictionary"
)

func (g *Generator) genAttributeStringOctets(w io.Writer, attr *dictionary.Attribute, vendor *dictionary.Vendor) {
	ident := identifier(attr.Name)
	var vendorIdent string
	if vendor != nil {
		vendorIdent = identifier(vendor.Name)
	}

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Add(p *radius.Packet, value []byte) (err error) {`)
	} else {
		p(w, `func `, ident, `_Add(p *radius.Packet, tag byte, value []byte) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `	a, err = radius.NewUserPassword(value, p.Secret, p.Authenticator[:])`)
	} else {
		p(w, `	a, err = radius.NewBytes(value)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `	a, err = radius.NewTag(tag, a)`)
		p(w, `	if err != nil {`)
		p(w, `		return`)
		p(w, `	}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_AddString(p *radius.Packet, value string) (err error) {`)
	} else {
		p(w, `func `, ident, `_AddString(p *radius.Packet, tag byte, value string) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `	a, err = radius.NewUserPassword([]byte(value), p.Secret, p.Authenticator[:])`)
	} else {
		p(w, `	a, err = radius.NewString(value)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `	a, err = radius.NewTag(tag, a)`)
		p(w, `	if err != nil {`)
		p(w, `		return`)
		p(w, `	}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Get(p *radius.Packet) (value []byte) {`)
		p(w, `	value, _ = `, ident, `_Lookup(p)`)
	} else {
		p(w, `func `, ident, `_Get(p *radius.Packet) (tag byte, value []byte) {`)
		p(w, `	tag, value, _ = `, ident, `_Lookup(p)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_GetString(p *radius.Packet) (value string) {`)
		p(w, `	return string(`, ident, `_Get(p))`)
	} else {
		p(w, `func `, ident, `_GetString(p *radius.Packet) (tag byte, value string) {`)
		p(w, `	var valueBytes []byte`)
		p(w, `	tag, valueBytes = `, ident, `_Get(p)`)
		p(w, `	value = string(valueBytes)`)
		p(w, `	return`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (values [][]byte, err error) {`)
	} else {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (tags []byte, values [][]byte, err error) {`)
	}
	p(w, `	var i []byte`)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if attr.HasTag() {
		p(w, `		tag, attr, err = radius.Tag(attr)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `		i, err = radius.UserPassword(attr, p.Secret, p.Authenticator[:])`)
	} else {
		p(w, `		i = radius.Bytes(attr)`)
	}
	p(w, `		if err != nil {`)
	p(w, `			return`)
	p(w, `		}`)
	p(w, `		values = append(values, i)`)
	if attr.HasTag() {
		p(w, `		tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_GetStrings(p *radius.Packet) (values []string, err error) {`)
	} else {
		p(w, `func `, ident, `_GetStrings(p *radius.Packet) (tags []byte, values []string, err error) {`)
	}
	p(w, `	var i string`)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if attr.HasTag() {
		p(w, `		tag, attr, err = radius.Tag(attr)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `		var up radius.Attribute`)
		p(w, `		up, err = radius.UserPassword(attr, p.Secret, p.Authenticator[:])`)
		p(w, `		if err == nil {`)
		p(w, `			i = string(up)`)
		p(w, `		}`)
	} else {
		p(w, `		i = radius.String(attr)`)
	}
	p(w, `		if err != nil {`)
	p(w, `			return`)
	p(w, `		}`)
	p(w, `		values = append(values, i)`)
	if attr.HasTag() {
		p(w, `		tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (value []byte, err error) {`)
	} else {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (tag byte, value []byte, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `	value, err = radius.UserPassword(a, p.Secret, p.Authenticator[:])`)
	} else {
		p(w, `	value = radius.Bytes(a)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_LookupString(p *radius.Packet) (value string, err error) {`)
	} else {
		p(w, `func `, ident, `_LookupString(p *radius.Packet) (tag byte, value string, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `	var b []byte`)
		p(w, `	b, err = radius.UserPassword(a, p.Secret, p.Authenticator[:])`)
		p(w, `	if err == nil {`)
		p(w, `		value = string(b)`)
		p(w, `	}`)
	} else {
		p(w, `	value = radius.String(a)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Set(p *radius.Packet, value []byte) (err error) {`)
	} else {
		p(w, `func `, ident, `_Set(p *radius.Packet, tag byte, value []byte) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `	a, err = radius.NewUserPassword(value, p.Secret, p.Authenticator[:])`)
	} else {
		p(w, `	a, err = radius.NewBytes(value)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, a)`)
		p(w, `	return`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_SetString(p *radius.Packet, value string) (err error) {`)
	} else {
		p(w, `func `, ident, `_SetString(p *radius.Packet, tag byte, value string) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	if attr.FlagEncrypt != nil && *attr.FlagEncrypt == 1 {
		p(w, `	a, err = radius.NewUserPassword([]byte(value), p.Secret, p.Authenticator[:])`)
	} else {
		p(w, `	a, err = radius.NewString(value)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, a)`)
		p(w, `	return`)
	}
	p(w, `}`)
}

func (g *Generator) genAttributeIPAddr(w io.Writer, attr *dictionary.Attribute, vendor *dictionary.Vendor, length int) {
	if length != net.IPv4len && length != net.IPv6len {
		panic("invalid length")
	}

	ident := identifier(attr.Name)
	var vendorIdent string
	if vendor != nil {
		vendorIdent = identifier(vendor.Name)
	}

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Add(p *radius.Packet, value net.IP) (err error) {`)
	} else {
		p(w, `func `, ident, `_Add(p *radius.Packet, tag byte, value net.IP) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	if length == net.IPv4len {
		p(w, `	a, err = radius.NewIPAddr(value)`)
	} else {
		p(w, `	a, err = radius.NewIPv6Addr(value)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Get(p *radius.Packet) (value net.IP) {`)
		p(w, `	value, _ = `, ident, `_Lookup(p)`)
	} else {
		p(w, `func `, ident, `_Get(p *radius.Packet) (tag byte, value net.IP) {`)
		p(w, `	tag, value, _ = `, ident, `_Lookup(p)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (values []net.IP, err error) {`)
	} else {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (tags []byte, values []net.IP, err error) {`)
	}
	p(w, `	var i net.IP`)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if length == net.IPv4len {
		p(w, `		i, err = radius.IPAddr(attr)`)
	} else {
		p(w, `		i, err = radius.IPv6Addr(attr)`)
	}
	p(w, `		if err != nil {`)
	p(w, `			return`)
	p(w, `		}`)
	if attr.HasTag() {
		p(w, `		i, err = radius.NewTag(tag, i)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	p(w, `		values = append(values, i)`)
	if attr.HasTag() {
		p(w, `		tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (value net.IP, err error) {`)
	} else {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (tag byte, value net.IP, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if length == net.IPv4len {
		p(w, `	value, err = radius.IPAddr(a)`)
	} else {
		p(w, `	value, err = radius.IPv6Addr(a)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Set(p *radius.Packet, value net.IP) (err error) {`)
	} else {
		p(w, `func `, ident, `_Set(p *radius.Packet, tag byte, value net.IP) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if length == net.IPv4len {
		p(w, `	a, err = radius.NewIPAddr(value)`)
	} else {
		p(w, `	a, err = radius.NewIPv6Addr(value)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)
}

func (g *Generator) genAttributeIFID(w io.Writer, attr *dictionary.Attribute, vendor *dictionary.Vendor) {
	ident := identifier(attr.Name)
	var vendorIdent string
	if vendor != nil {
		vendorIdent = identifier(vendor.Name)
	}

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Add(p *radius.Packet, value net.HardwareAddr) (err error) {`)
	} else {
		p(w, `func `, ident, `_Add(p *radius.Packet, tag byte, value net.HardwareAddr) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	p(w, `	a, err = radius.NewIFID(value)`)
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Get(p *radius.Packet) (value net.HardwareAddr) {`)
		p(w, `	value, _ = `, ident, `_Lookup(p)`)
	} else {
		p(w, `func `, ident, `_Get(p *radius.Packet) (tag byte, value net.HardwareAddr) {`)
		p(w, `	tag, value, _ = `, ident, `_Lookup(p)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (values []net.HardwareAddr, err error) {`)
	} else {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (tags []byte, values []net.HardwareAddr, err error) {`)
	}
	p(w, `	var i net.HardwareAddr`)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if attr.HasTag() {
		p(w, `		tag, attr, err = radius.Tag(attr)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	p(w, `		i, err = radius.IFID(attr)`)
	p(w, `		if err != nil {`)
	p(w, `			return`)
	p(w, `		}`)
	p(w, `		values = append(values, i)`)
	if attr.HasTag() {
		p(w, `		tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (value net.HardwareAddr, err error) {`)
	} else {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (tag byte, value net.HardwareAddr, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	p(w, `	value, err = radius.IFID(a)`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Set(p *radius.Packet, value net.HardwareAddr) (err error) {`)
	} else {
		p(w, `func `, ident, `_Set(p *radius.Packet, tag byte, value net.HardwareAddr) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	p(w, `	a, err = radius.NewIFID(value)`)
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)
}

func (g *Generator) genAttributeDate(w io.Writer, attr *dictionary.Attribute, vendor *dictionary.Vendor) {
	ident := identifier(attr.Name)
	var vendorIdent string
	if vendor != nil {
		vendorIdent = identifier(vendor.Name)
	}

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Add(p *radius.Packet, value time.Time) (err error) {`)
	} else {
		p(w, `func `, ident, `_Add(p *radius.Packet, tag byte, value time.Time) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	p(w, `	a, err = radius.NewDate(value)`)
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Get(p *radius.Packet) (value time.Time) {`)
		p(w, `	value, _ = `, ident, `_Lookup(p)`)
	} else {
		p(w, `func `, ident, `_Get(p *radius.Packet) (tag byte, value time.Time) {`)
		p(w, `	tag, value, _ = `, ident, `_Lookup(p)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (values []time.Time, err error) {`)
	} else {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (tags []byte, values []time.Time, err error) {`)
	}
	p(w, `	var i time.Time`)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if attr.HasTag() {
		p(w, `		tag, attr, err = radius.Tag(attr)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	p(w, `		i, err = radius.Date(attr)`)
	p(w, `		if err != nil {`)
	p(w, `			return`)
	p(w, `		}`)
	p(w, `		values = append(values, i)`)
	if attr.HasTag() {
		p(w, `		tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (value time.Time, err error) {`)
	} else {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (tag byte, value time.Time, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	p(w, `	value, err = radius.Date(a)`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Set(p *radius.Packet, value time.Time) (err error) {`)
	} else {
		p(w, `func `, ident, `_Set(p *radius.Packet, tag byte, value time.Time) (err error) {`)
	}
	p(w, `	var a radius.Attribute`)
	p(w, `	a, err = radius.NewDate(value)`)
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)
}

func (g *Generator) genAttributeInteger(w io.Writer, attr *dictionary.Attribute, allValues []*dictionary.Value, bitsize int, vendor *dictionary.Vendor) {
	var values []*dictionary.Value
	for _, value := range allValues {
		if value.Attribute == attr.Name {
			if len(values) > 0 && values[len(values)-1].Number == value.Number {
				values[len(values)-1] = value
			} else {
				values = append(values, value)
			}
		}
	}

	ident := identifier(attr.Name)
	var vendorIdent string
	if vendor != nil {
		vendorIdent = identifier(vendor.Name)
	}

	p(w)
	if bitsize == 64 {
		p(w, `type `, ident, ` uint64`)
	} else { // 32
		p(w, `type `, ident, ` uint32`)
	}

	// Values
	if len(values) > 0 {
		p(w)
		p(w, `const (`)
		for _, value := range values {
			valueIdent := identifier(value.Name)
			p(w, `	`, ident, `_Value_`, valueIdent, ` `, ident, ` = `, strconv.Itoa(value.Number))
		}
		p(w, `)`)
	}

	p(w)
	p(w, `var `, ident, `_Strings = map[`, ident, `]string{`)
	for _, value := range values {
		valueIdent := identifier(value.Name)
		p(w, `	`, ident, `_Value_`, valueIdent, `: `, strconv.Quote(value.Name), `,`)
	}
	p(w, `}`)

	p(w)
	p(w, `func (a `, ident, `) String() string {`)
	p(w, `	if str, ok := `, ident, `_Strings[a]; ok {`)
	p(w, `		return str`)
	p(w, `	}`)
	p(w, `	return "`, ident, `(" + strconv.FormatUint(uint64(a), 10) + ")"`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Add(p *radius.Packet, value `, ident, `) (err error) {`)
	} else {
		p(w, `func `, ident, `_Add(p *radius.Packet, tag byte, value `, ident, `) (err error) {`)
	}
	if bitsize == 64 {
		p(w, `	a := radius.NewInteger64(uint64(value))`)
	} else { // 32
		p(w, `	a := radius.NewInteger(uint32(value))`)
	}
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Get(p *radius.Packet) (value `, ident, `) {`)
		p(w, `	value, _ = `, ident, `_Lookup(p)`)
	} else {
		p(w, `func `, ident, `_Get(p *radius.Packet) (tag byte, value `, ident, `) {`)
		p(w, `	tag, value, _ = `, ident, `_Lookup(p)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (values []`, ident, `, err error) {`)
	} else {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (tags []byte, values []`, ident, `, err error) {`)
	}
	if bitsize == 64 {
		p(w, `	var i uint64`)
	} else { // 32
		p(w, `	var i uint32`)
	}
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if attr.HasTag() {
		p(w, `		tag, attr, err = radius.Tag(attr)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if bitsize == 64 {
		p(w, `		i, err = radius.Integer64(attr)`)
	} else { // 32
		p(w, `		i, err = radius.Integer(attr)`)
	}
	p(w, `		if err != nil {`)
	p(w, `			return`)
	p(w, `		}`)
	p(w, `		values = append(values, `, ident, `(i))`)
	if attr.HasTag() {
		p(w, `		tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (value `, ident, `, err error) {`)
	} else {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (tag byte, value `, ident, `, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w, `		tag, a, err = radius.Tag(a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if bitsize == 64 {
		p(w, `	var i uint64`)
		p(w, `	i, err = radius.Integer64(a)`)
	} else { // 32
		p(w, `	var i uint32`)
		p(w, `	i, err = radius.Integer(a)`)
	}
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	p(w, `	value = `, ident, `(i)`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Set(p *radius.Packet, value `, ident, `) (err error) {`)
	} else {
		p(w, `func `, ident, `_Set(p *radius.Packet, tag byte, value `, ident, `) (err error) {`)
	}
	if bitsize == 64 {
		p(w, `	a := radius.NewInteger64(uint64(value))`)
	} else { // 32
		p(w, `	a := radius.NewInteger(uint32(value))`)
	}
	if attr.HasTag() {
		p(w, `		a, err = radius.NewTag(tag, a)`)
		p(w, `		if err != nil {`)
		p(w, `			return`)
		p(w, `		}`)
	}
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, a)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, a)`)
		p(w, `	return nil`)
	}
	p(w, `}`)
}

func subIdentifier(parentIdent string, name string) string {
	return strings.Join([]string{parentIdent, identifier(name)}, `_`)
}

func subType(parentIdent string, name string) string {
	return strings.Join([]string{subIdentifier(parentIdent, name), `Type`}, `_`)
}

func genNewTLVObject(w io.Writer, attr *dictionary.Attribute) {
	ident := identifier(attr.Name)
	p(w)
	if !attr.HasTag() {
		p(w, `func new`, ident, `(value `, ident, `) (attribute radius.Attribute, err error) {`)
	} else {
		p(w, `func new`, ident, `(tag byte, value `, ident, `) (attribute radius.Attribute, err error) {`)
	}
	p(w, `  var typedAttributes []radius.TypedAttribute`)
	p(w, `  var a radius.Attribute`)

	for _, subAttr := range attr.Attributes {
		p(w)
		switch subAttr.Type {
		case dictionary.AttributeString:
			p(w, `  a, err = radius.NewString(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeOctets:
			p(w, `  a, err = radius.NewBytes(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeIPAddr:
			p(w, `  a, err = radius.NewIPAddr(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeDate:
			p(w, `  a, err = radius.NewDate(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeInteger:
			p(w, `  a = radius.NewInteger(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeIPv6Addr:
			p(w, `  a, err = radius.NewIPv6Addr(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeIFID:
			p(w, `  a, err = radius.NewIFID(value.`, identifier(subAttr.Name), `)`)
		case dictionary.AttributeInteger64:
			p(w, `  a = radius.NewInteger64(value.`, identifier(subAttr.Name), `)`)
		}
		p(w, `  if err != nil {`)
		p(w, `	  return nil, err`)
		p(w, `	}`)
		p(w, `  typedAttributes = append(typedAttributes, radius.TypedAttribute{`, subType(ident, subAttr.Name), `, a})`)
	}
	p(w)
	p(w, `	attribute, err = radius.NewTLV(typedAttributes)`)
	p(w, `  if err != nil {`)
	p(w, `	  return nil, err`)
	p(w, `	}`)
	if attr.HasTag() {
		p(w)
		p(w, `  a, err = radius.NewTag(tag, a)`)
		p(w, `	if err != nil {`)
		p(w, `	  return = nil, err`)
		p(w, `	}`)
	}
	p(w, `return`)
	p(w, `}`)
}

func getTLVObject(w io.Writer, attr *dictionary.Attribute) {
	ident := identifier(attr.Name)
	isFirst := true
	p(w)
	if !attr.HasTag() {
		p(w, `func set`, ident, `(a radius.Attribute) (values []`, ident, `,err error) {`)
	} else {
		p(w, `func set`, ident, `(a radius.Attribute) (tag byte, values []`, ident, `,err error) {`)
	}
	p(w, `	var attributes radius.Attributes`)
	p(w, `  valuesLen := -1`)
	if attr.HasTag() {
		p(w)
		p(w, `  tag, a, err = radius.Tag(a)`)
		p(w, `	if err != nil {`)
		p(w, `	  return`)
		p(w, `	}`)
	}
	p(w)
	p(w, `	attributes, err = radius.TLV(a)`)
	p(w, `	if err != nil {`)
	p(w, `		return`)
	p(w, `	}`)
	for _, subAttr := range attr.Attributes {
		p(w)
		p(w, `  if val, ok := attributes[`, subType(ident, subAttr.Name), `]; ok {`)
		if isFirst {
			isFirst = false
			p(w, `    valuesLen = len(val)`)
			p(w, `    values = make([]`, ident, `, valuesLen)`)
		}
		p(w, `    if len(val) != valuesLen {`)
		p(w, `      err = radius.ErrTLVAttribute`)
		p(w, `    } else {`)
		p(w, `      for i := range(val) {`)
		switch subAttr.Type {
		case dictionary.AttributeString:
			p(w, `        values[i].`, identifier(subAttr.Name), `= radius.String(val[i])`)
		case dictionary.AttributeOctets:
			p(w, `        values[i].`, identifier(subAttr.Name), `= radius.Bytes(val[i])`)
		case dictionary.AttributeIPAddr:
			p(w, `        values[i].`, identifier(subAttr.Name), `, err = radius.IPAddr(val[i])`)
		case dictionary.AttributeDate:
			p(w, `        values[i].`, identifier(subAttr.Name), `, err = radius.Date(val[i])`)
		case dictionary.AttributeInteger:
			p(w, `        values[i].`, identifier(subAttr.Name), `, err = radius.Integer(val[i])`)
		case dictionary.AttributeIPv6Addr:
			p(w, `        values[i].`, identifier(subAttr.Name), `, err = radius.IPv6Addr(val[i])`)
		case dictionary.AttributeIFID:
			p(w, `        values[i].`, identifier(subAttr.Name), `, err = radius.IFID(val[i])`)
		case dictionary.AttributeInteger64:
			p(w, `        values[i].`, identifier(subAttr.Name), `, err = radius.Integer64(val[i])`)
		}
		p(w, `      }`)
		p(w, `    }`)
		p(w, `  } else {`)
		p(w, `    err = radius.ErrTLVAttribute`)
		p(w, `  }`)
		p(w, `  if err != nil {`)
		p(w, `    return`)
		p(w, `  }`)
	}
	p(w, `  return`)
	p(w, `}`)
}

func (g *Generator) genAttributeTLV(w io.Writer, attr *dictionary.Attribute, allValues []*dictionary.Value, vendor *dictionary.Vendor) {
	ident := identifier(attr.Name)
	var vendorIdent string
	if vendor != nil {
		vendorIdent = identifier(vendor.Name)
	}

	p(w)
	p(w, `type `, ident, ` struct {`)
	for _, subAttr := range attr.Attributes {
		typeStr := subAttr.Type.TypeDef()
		if len(typeStr) > 0 {
			p(w, identifier(subAttr.Name), ` `, typeStr)
		}
	}
	p(w, `}`)

	if len(attr.Attributes) > 0 {
		// create const for all the attributes
		p(w)
		p(w, `const (`)
		for _, subAttr := range attr.Attributes {
			codes := strings.Split(subAttr.OID, ".")
			p(w, `	`, subType(ident, subAttr.Name), ` radius.Type = `, codes[len(codes)-1])
		}
		p(w, `)`)
	}

	genNewTLVObject(w, attr)
	getTLVObject(w, attr)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Add(p *radius.Packet, values []`, ident, `) (error) {`)
	} else {
		p(w, `func `, ident, `_Add(p *radius.Packet, tag byte, values []`, ident, `) (error) {`)
	}
	p(w, `  var attribute radius.Attribute`)
	p(w, `  if len(values) < 1 {`)
	p(w, `    return radius.ErrEmptyStruct`)
	p(w, `  }`)
	p(w, `  for _, value := range values {`)
	if !attr.HasTag() {
		p(w, `    _attr, _err :=  new`, ident, `(value)`)
	} else {
		p(w, `    _attr, _err :=  new`, ident, `(tag, value)`)
	}
	p(w, ` 	  if _err != nil {`)
	p(w, `	    return _err`)
	p(w, `	  }`)
	p(w, `    attribute = append(attribute, _attr...)`)
	p(w, `  }`)

	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_AddVendor(p, `, attr.OID, `, attribute)`)
	} else {
		p(w, `	p.Add(`, ident, `_Type, attribute)`)
		p(w, `	return nil`)
	}
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Get(p *radius.Packet) (value []`, ident, `) {`)
		p(w, `	value, _ = `, ident, `_Lookup(p)`)
	} else {
		p(w, `func `, ident, `_Get(p *radius.Packet) (tag byte, value []`, ident, `) {`)
		p(w, `	tag, value, _ = `, ident, `_Lookup(p)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (values [][]`, ident, `, err error) {`)
	} else {
		p(w, `func `, ident, `_Gets(p *radius.Packet) (tags []byte, values [][]`, ident, `, err error) {`)
	}
	p(w, `	var value []`, ident)
	if attr.HasTag() {
		p(w, `	var tag byte`)
	}
	if vendor != nil {
		p(w, `	for _, attr := range _`, vendorIdent, `_GetsVendor(p, `, attr.OID, `) {`)
	} else {
		p(w, `	for _, attr := range p.Attributes[`, ident, `_Type] {`)
	}
	if !attr.HasTag() {
		p(w, `    value, err = set`, ident, `(attr)`)
	} else {
		p(w, `    tag, value, err = set`, ident, `(attr)`)
	}
	p(w, `    if err != nil {`)
	p(w, `      return`)
	p(w, `    }`)
	p(w, `    values = append(values, []`, ident, `(value))`)
	if attr.HasTag() {
		p(w, `    tags = append(tags, tag)`)
	}
	p(w, `	}`)
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (values []`, ident, `, err error) {`)
	} else {
		p(w, `func `, ident, `_Lookup(p *radius.Packet) (tag byte, values []`, ident, `, err error) {`)
	}
	if vendor != nil {
		p(w, `	a, ok  := _`, vendorIdent, `_LookupVendor(p, `, attr.OID, `)`)
	} else {
		p(w, `	a, ok  := p.Lookup(`, ident, `_Type)`)
	}
	p(w, `	if !ok {`)
	p(w, `		err = radius.ErrNoAttribute`)
	p(w, `		return`)
	p(w, `	}`)
	if !attr.HasTag() {
		p(w, `  values, err = set`, ident, `(a)`)
	} else {
		p(w, `  tag, values, err = set`, ident, `(a)`)
	}
	p(w, `	return`)
	p(w, `}`)

	p(w)
	if !attr.HasTag() {
		p(w, `func `, ident, `_Set(p *radius.Packet, values []`, ident, `) (error) {`)
	} else {
		p(w, `func `, ident, `_Set(p *radius.Packet, tag byte, value []`, ident, `) (error) {`)
	}
	p(w, `  var attribute radius.Attribute`)
	p(w, `  for _, value := range values {`)
	if !attr.HasTag() {
		p(w, `    _attr, _err := new`, ident, `(value)`)
	} else {
		p(w, `    _attr, _err := new`, ident, `(tag, value)`)
	}
	p(w, ` 	  if _err != nil {`)
	p(w, `	    return _err`)
	p(w, `	  }`)
	p(w, `    attribute = append(attribute, _attr...)`)
	p(w, `  }`)
	if vendor != nil {
		p(w, `	return _`, vendorIdent, `_SetVendor(p, `, attr.OID, `, attribute)`)
	} else {
		p(w, `	p.Set(`, ident, `_Type, attribute)`)
		p(w, `	return nil`)
	}
	p(w, `}`)
}
