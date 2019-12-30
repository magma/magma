// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"errors"
	"reflect"
	"strings"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// parseAvpTag return the avp_name and omitempty option
func parseAvpTag(tag reflect.StructTag) (string, bool) {
	if tag == "" {
		return "", false
	}
	name := string(tag)
	if strings.HasPrefix(name, "avp:\"") {
		name = name[5 : len(name)-1]
		omitEmpty := false
		if strings.HasSuffix(name, ",omitempty") {
			name = name[0 : len(name)-10]
			omitEmpty = true
		}
		if strings.IndexByte(name, '"') == -1 {
			return name, omitEmpty
		}
	}

	name = tag.Get("avp")
	if idx := strings.Index(name, ","); idx != -1 {
		return name[:idx], false
	}
	return name, true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// Marshal encodes struct into AVPs
func (m *Message) Marshal(src interface{}) error {
	v := reflect.ValueOf(src)
	if v.Kind() != reflect.Ptr {
		return errors.New("src is not a pointer to struct")
	}
	err, avps := marshalStruct(m, v)
	if err != nil {
		return err
	}
	m.AVP = avps
	m.Header.MessageLength = uint32(m.Len())
	return nil
}

func marshalStruct(m *Message, field reflect.Value) (error, []*AVP) {
	var err error
	var dictAVP *dict.AVP
	var avps []*AVP

	base := reflect.Indirect(field)
	if base.Kind() != reflect.Struct {
		return errors.New("src is not a pointer to struct"), nil
	}

	for n := 0; n < base.NumField(); n++ {
		f := base.Field(n)
		bt := base.Type().Field(n)
		avpname, omitEmpty := parseAvpTag(bt.Tag)
		if len(avpname) == 0 || (omitEmpty && isEmptyValue(f)) {
			// TODO: check the required attribute in AVP rule?
			continue
		}

		// Lookup the AVP name (tag) in the dictionary, the dictionary AVP has the code.
		// Relies on the fact that in the same app will not be AVPs with same code but different vendorId
		dictAVP, err = m.Dictionary().FindAVP(m.Header.ApplicationID, avpname)
		if err != nil {
			return err, nil
		}

		err, avp := marshal(m, f, dictAVP)
		if err != nil {
			return err, nil
		}
		avps = append(avps, avp...)
	}

	return nil, avps
}

// marshal returns a AVP type of the field
func marshal(m *Message, field reflect.Value, fieldAVP *dict.AVP) (error, []*AVP) {
	var data datatype.Type
	var avps []*AVP // avps := make([]*AVP, 0, 8)
	fieldType := field.Type()

	// log.Println(fieldAVP.Name, " begin ", field.Kind())
	// defer log.Println(fieldAVP.Name, " end")

	var t reflect.Type
	switch field.Kind() {
	case reflect.Slice:
		// 1. []byte
		//  (1) dicttype.Grouped
		//  (2) other basic type which can be a Slice. for example eg. datatype.AddressType = net.IP = []byte
		// if fieldType == reflect.TypeOf(([]byte)(nil))
		if fieldType.Elem().Kind() == reflect.Uint8 {
			goto BASIC_TYPE
		}

		// 2.  []*diam.AVP
		if fieldType == reflect.TypeOf(([]*AVP)(nil)) {
			// s := reflect.New(fieldType) // a pointer to a slice
			// s.Elem().Set(field)
			avp := field.Interface().([]*AVP)
			avps = append(avps, avp...)
			return nil, avps
		}

		// 3. real array of diameter AVPs
		// log.Print("Slice len:", field.Len())
		for n := 0; n < field.Len(); n++ {
			err, avp := marshal(m, field.Index(n), fieldAVP)
			if err != nil {
				return err, nil
			}
			avps = append(avps, avp...)
		}
		return nil, avps

	case reflect.Interface, reflect.Ptr:
		if field.IsNil() {
			return nil, avps // skip optional AVP
		}
		return marshal(m, field.Elem(), fieldAVP)
	}

BASIC_TYPE:
	switch fieldAVP.Data.Type {
	case datatype.AddressType:
		t = reflect.TypeOf((*datatype.Address)(nil)).Elem() // get Type of datatype.Address
	case datatype.DiameterIdentityType:
		t = reflect.TypeOf((*datatype.DiameterIdentity)(nil)).Elem()
	case datatype.DiameterURIType:
		t = reflect.TypeOf((*datatype.DiameterURI)(nil)).Elem()
	case datatype.EnumeratedType:
		t = reflect.TypeOf((*datatype.Enumerated)(nil)).Elem()
	case datatype.Float32Type:
		t = reflect.TypeOf((*datatype.Float32)(nil)).Elem()
	case datatype.Float64Type:
		t = reflect.TypeOf((*datatype.Float64)(nil)).Elem()
	case datatype.IPFilterRuleType:
		t = reflect.TypeOf((*datatype.IPFilterRule)(nil)).Elem()
	case datatype.IPv4Type:
		t = reflect.TypeOf((*datatype.IPv4)(nil)).Elem()
	case datatype.Integer32Type:
		t = reflect.TypeOf((*datatype.Integer32)(nil)).Elem()
	case datatype.Integer64Type:
		t = reflect.TypeOf((*datatype.Integer64)(nil)).Elem()
	case datatype.OctetStringType:
		t = reflect.TypeOf((*datatype.OctetString)(nil)).Elem()
	case datatype.TimeType:
		t = reflect.TypeOf((*datatype.Time)(nil)).Elem()
	case datatype.UTF8StringType:
		t = reflect.TypeOf((*datatype.UTF8String)(nil)).Elem()
	case datatype.Unsigned32Type:
		t = reflect.TypeOf((*datatype.Unsigned32)(nil)).Elem()
	case datatype.Unsigned64Type:
		t = reflect.TypeOf((*datatype.Unsigned64)(nil)).Elem()
	case datatype.GroupedType:
		if field.Kind() == reflect.Struct {
			// 1.  diam.AVP
			// if fieldType.String() == "diam.AVP"
			if fieldType == reflect.TypeOf(AVP{}) {
				p := reflect.New(fieldType)
				v := reflect.ValueOf(p).Elem()
				v.Set(field)
				avp := p.Interface().(*AVP)
				return nil, append(avps, avp)
			}

			// 2. GroupedAVP
			gAVP := &GroupedAVP{}
			for n := 0; n < field.NumField(); n++ {
				f := field.Field(n)
				bt := field.Type().Field(n)
				avpname, omitEmpty := parseAvpTag(bt.Tag)
				if len(avpname) == 0 || (omitEmpty && isEmptyValue(f)) {
					// TODO: check the required attribute in AVP rule?
					continue
				}
				// Lookup the AVP name (tag) in the dictionary, the dictionary AVP has the code.
				// Relies on the fact that in the same app will not be AVPs with same code but different vendorId
				d, err := m.Dictionary().FindAVP(m.Header.ApplicationID, avpname)
				if err != nil {
					return err, nil
				}
				err, avp := marshal(m, f, d)
				if err != nil {
					return err, nil
				}
				gAVP.AVP = append(gAVP.AVP, avp...) // gAVP.AddAVP()
			}
			data = gAVP
		} else if field.Kind() == reflect.Slice {
			// when code run here, we are certain that it is datatype.Grouped AVP
			// like "Failed-AVP", all we need to do is assigning the []byte slibe
			//  to a datatype.Grouped
			t = reflect.TypeOf((*datatype.Grouped)(nil)).Elem()
			break
		} else {
			return errors.New(fieldAVP.Name + " AVP's Data type is unknown."), nil
		}
	default:
		return errors.New(fieldAVP.Name + " AVP's Data type is unknown."), nil
	}

	if data == nil { // basic non-grouped AVP
		p := reflect.New(t)
		v := reflect.Indirect(p)

		if fieldType.AssignableTo(t) {
			// log.Println("assign: ", fieldAVP.Name, " ", fieldType.String(), " => ", t.String())
			v.Set(field)
		} else if fieldType.ConvertibleTo(t) {
			// log.Println("convert: ", fieldAVP.Name, " ", fieldType.String(), " => ", t.String())
			v.Set(field.Convert(t))
		} else {
			return errors.New(fieldAVP.Name + " AVP type mismatched. " + fieldType.String() + " => " + t.String()), nil
		}
		var ok bool
		data, ok = v.Interface().(datatype.Type)
		if !ok {
			return errors.New(fieldAVP.Name + ", failed to convert AVP data to datatype.Type"), nil
		}
	}

	var avpFlags uint8
	if strings.Contains(fieldAVP.Must, "M") {
		avpFlags = avp.Mbit
	}
	if fieldAVP.VendorID > 0 {
		avpFlags |= avp.Vbit
	}

	avp := &AVP{
		Code:     fieldAVP.Code,
		Flags:    avpFlags,
		VendorID: fieldAVP.VendorID,
		Data:     data,
	}

	return nil, append(avps, avp)
}

// Unmarshal stores the result of a diameter message in the struct
// pointed to by dst.
//
// Unmarshal can not only decode AVPs into the struct, but also their
// Go equivalent data types, directly.
//
// For example:
//
//	type CER struct {
//		OriginHost  AVP    `avp:"Origin-Host"`
//		.. or
//		OriginHost  *AVP   `avp:"Origin-Host"`
//		.. or
//		OriginHost  string `avp:"Origin-Host"`
//	}
//	var d CER
//	err := diam.Unmarshal(&d)
//
// This decodes the Origin-Host AVP as three different types. The first, AVP,
// makes a copy of the AVP in the message and stores in the struct. The
// second, *AVP, stores a pointer to the original AVP in the message. If you
// change the values of it, you're actually changing the message.
// The third decodes the inner contents of AVP.Data, which in this case is
// a format.DiameterIdentity, and stores the value of it in the struct.
//
// Unmarshal supports all the basic Go types, including slices, for multiple
// AVPs of the same type) and structs, for grouped AVPs.
//
// Slices:
//
//	type CER struct {
//		Vendors  []*AVP `avp:"Supported-Vendor-Id"`
//	}
//	var d CER
//	err := diam.Unmarshal(&d)
//
// Slices have the same principles of other types. If they're of type
// []*AVP it'll store references in the struct, while []AVP makes
// copies and []int (or []string, etc) decodes the AVP data for you.
//
// Grouped AVPs:
//
//	type VSA struct {
//		AuthAppID int `avp:"Auth-Application-Id"`
//		VendorID  int `avp:"Vendor-Id"`
//	}
//	type CER struct {
//		VSA VSA  `avp:"Vendor-Specific-Application-Id"`
//		.. or
//		VSA *VSA `avp:"Vendor-Specific-Application-Id"`
//		.. or
//		VSA struct {
//			AuthAppID int `avp:"Auth-Application-Id"`
//			VendorID  int `avp:"Vendor-Id"`
//		} `avp:"Vendor-Specific-Application-Id"`
//	}
//	var d CER
//	err := m.Unmarshal(&d)
//
// Other types are supported as well, such as net.IP and time.Time where
// applicable. See the format sub-package for details. Usually, you want
// to decode values to their native Go type when the AVPs don't have to be
// re-used in an answer, such as Origin-Host and friends. The ones that are
// usually added to responses, such as Origin-State-Id are better decoded to
// just AVP or *AVP, making it easier to re-use them in the answer.
//
// Note that decoding values to *AVP is much faster and more efficient than
// decoding to AVP or the native Go types.
func (m *Message) Unmarshal(dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return errors.New("dst is not a pointer to struct")
	}
	return scanStruct(m, v, m.AVP)
}

// newIndex returns a map of AVPs indexed by their code.
// TODO: make this part of the Message.
func newIndex(avps []*AVP) map[uint32][]*AVP {
	idx := make(map[uint32][]*AVP, len(avps))
	for _, a := range avps {
		idx[a.Code] = append(idx[a.Code], a)
	}
	return idx
}

func scanStruct(m *Message, field reflect.Value, avps []*AVP) error {
	base := reflect.Indirect(field)
	if base.Kind() != reflect.Struct {
		return errors.New("dst is not a pointer to struct")
	}
	idx := newIndex(avps)
	for n := 0; n < base.NumField(); n++ {
		f := base.Field(n)
		bt := base.Type().Field(n)
		avpname, _ := parseAvpTag(bt.Tag)
		if len(avpname) == 0 {
			continue
		}
		// Lookup the AVP name (tag) in the dictionary.
		// The dictionary AVP has the code.
		d, err := m.Dictionary().FindAVP(m.Header.ApplicationID, avpname) // Relies on the fact that in the same app will not be AVPs with same code but different vendorId
		if err != nil {
			return err
		}
		// See if this AVP exist in the message.
		avps, exists := idx[d.Code]
		if !exists {
			continue
		}
		//log.Println("Handling", f, bt)
		unmarshal(m, f, avps)
	}
	return nil
}

func unmarshal(m *Message, f reflect.Value, avps []*AVP) {
	fieldType := f.Type()
	switch f.Kind() {
	case reflect.Slice:
		// Copy byte arrays.
		dv := reflect.ValueOf(avps[0].Data)
		if dv.Type().ConvertibleTo(fieldType) {
			f.Set(dv.Convert(fieldType))
			break
		}

		// Allocate new slice and copy all items.
		f.Set(reflect.MakeSlice(fieldType, len(avps), len(avps)))
		// TODO: optimize?
		for n := 0; n < len(avps); n++ {
			unmarshal(m, f.Index(n), avps[n:])
		}

	case reflect.Interface, reflect.Ptr:
		if f.IsNil() {
			f.Set(reflect.New(fieldType.Elem()))
		}
		unmarshal(m, f.Elem(), avps)

	case reflect.Struct:
		// Test for *AVP
		at := reflect.TypeOf(avps[0])
		if fieldType.AssignableTo(at) {
			f.Set(reflect.ValueOf(avps[0]))
			break
		}
		// Test for AVP
		at = reflect.TypeOf(*avps[0])
		if fieldType.ConvertibleTo(at) {
			f.Set(reflect.ValueOf(*avps[0]))
			break
		}

		// Used for unmarshalling time datatype
		if fieldType.AssignableTo(reflect.TypeOf(avps[0].Data)) {
			f.Set(reflect.ValueOf(avps[0].Data))
			break
		}

		// Used for unmarshalling time type
		if fieldType.ConvertibleTo(reflect.TypeOf(avps[0].Data)) {
			timeStamp := reflect.ValueOf(avps[0].Data).Convert(fieldType)
			f.Set(timeStamp)
			break
		}

		// Handle grouped AVPs.
		if group, ok := avps[0].Data.(*GroupedAVP); ok {
			scanStruct(m, f, group.AVP)
		}

	default:
		// Test for AVP.Data (e.g. format.UTF8String, string)
		dv := reflect.ValueOf(avps[0].Data)
		if dv.Type().ConvertibleTo(fieldType) {
			f.Set(dv.Convert(fieldType))
		}
	}
}
