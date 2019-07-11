package dictionary

import (
	"bytes"
	"fmt"
	"strconv"
)

type Dictionary struct {
	Attributes []*Attribute
	Values     []*Value
	Vendors    []*Vendor
}

func (d *Dictionary) GoString() string {
	var b bytes.Buffer
	b.WriteString("&dictionary.Dictionary{")

	if len(d.Attributes) > 0 {
		b.WriteString("Attributes:[]*dictionary.Attribute{")
		for _, attr := range d.Attributes {
			fmt.Fprintf(&b, "%#v,", attr)
		}
		b.WriteString("},")
	}

	if len(d.Values) > 0 {
		b.WriteString("Values:[]*dictionary.Value{")
		for _, value := range d.Values {
			fmt.Fprintf(&b, "%#v,", value)
		}
		b.WriteString("},")
	}

	if len(d.Vendors) > 0 {
		b.WriteString("Vendors:[]*dictionary.Vendor{")
		for _, vendor := range d.Vendors {
			fmt.Fprintf(&b, "%#v,", vendor)
		}
		b.WriteString("},")
	}

	b.WriteString("}")
	return b.String()
}

type AttributeType int

const (
	AttributeString AttributeType = iota + 1
	AttributeOctets
	AttributeIPAddr
	AttributeDate
	AttributeInteger
	AttributeIPv6Addr
	AttributeIPv6Prefix
	AttributeIFID
	AttributeInteger64

	AttributeVSA
	AttributeTLV
	// TODO: non-standard types?
)

func (t AttributeType) String() string {
	switch t {
	case AttributeString:
		return "string"
	case AttributeOctets:
		return "octets"
	case AttributeIPAddr:
		return "ipaddr"
	case AttributeDate:
		return "date"
	case AttributeInteger:
		return "integer"
	case AttributeIPv6Addr:
		return "ipv6addr"
	case AttributeIPv6Prefix:
		return "ipv6prefix"
	case AttributeIFID:
		return "ifid"
	case AttributeInteger64:
		return "integer64"
	case AttributeVSA:
		return "vsa"
	case AttributeTLV:
		return "tlv"
	}
	return "AttributeType(" + strconv.Itoa(int(t)) + ")"
}

func (t AttributeType) TypeDef() string {
	switch t {
	case AttributeString:
		return "string"
	case AttributeOctets:
		return "[]byte"
	case AttributeIPAddr:
		return "net.IP"
	case AttributeDate:
		return "time.Time"
	case AttributeInteger:
		return "uint32"
	case AttributeIPv6Addr:
		return "net.IP"
	case AttributeIFID:
		return "net.HardwareAddr"
	case AttributeInteger64:
		return "uint64"
	}
	return ""
}

type Attribute struct {
	Name       string
	OID        string
	Type       AttributeType
	Attributes []*Attribute

	Size *int

	FlagEncrypt *int
	FlagHasTag  *bool
	FlagConcat  *bool
}

func (a *Attribute) HasTag() bool {
	return a.FlagHasTag != nil && *a.FlagHasTag
}

func (a *Attribute) Equals(o *Attribute) bool {
	if a == o {
		return true
	}
	if a == nil || o == nil {
		return false
	}

	if a.Name != o.Name || a.OID != o.OID || a.Type != o.Type {
		return false
	}

	if (a.Size == nil && o.Size != nil) || (a.Size != nil && o.Size == nil) || (a.Size != nil && o.Size != nil && *a.Size != *o.Size) {
		return false
	}

	if (a.FlagEncrypt == nil && o.FlagEncrypt != nil) || (a.FlagEncrypt != nil && o.FlagEncrypt == nil) || (a.FlagEncrypt != nil && o.FlagEncrypt != nil && *a.FlagEncrypt != *o.FlagEncrypt) {
		return false
	}
	if (a.FlagHasTag == nil && o.FlagHasTag != nil) || (a.FlagHasTag != nil && o.FlagHasTag == nil) || (a.FlagHasTag != nil && o.FlagHasTag != nil && *a.FlagHasTag != *o.FlagHasTag) {
		return false
	}
	if (a.FlagConcat == nil && o.FlagConcat != nil) || (a.FlagConcat != nil && o.FlagConcat == nil) || (a.FlagConcat != nil && o.FlagConcat != nil && *a.FlagConcat != *o.FlagConcat) {
		return false
	}

	return true
}

func (a *Attribute) GoString() string {
	var b bytes.Buffer
	b.WriteString("&dictionary.Attribute{")

	fmt.Fprintf(&b, "Name:%#v,", a.Name)
	fmt.Fprintf(&b, "OID:%#v,", a.OID)
	fmt.Fprintf(&b, "Type:%#v,", a.Type)

	if a.Size != nil {
		fmt.Fprintf(&b, "Size:dictionary.Int(%#v),", *(a.Size))
	}

	if a.FlagEncrypt != nil {
		fmt.Fprintf(&b, "FlagEncrypt:dictionary.Int(%#v),", *(a.FlagEncrypt))
	}
	if a.FlagHasTag != nil {
		fmt.Fprintf(&b, "FlagHasTag:dictionary.Bool(%#v),", *(a.FlagHasTag))
	}
	if a.FlagConcat != nil {
		fmt.Fprintf(&b, "FlagConcat:dictionary.Bool(%#v),", *(a.FlagConcat))
	}
	if a.Attributes != nil {
		for _, attr := range a.Attributes {
			fmt.Fprintf(&b, attr.GoString())
		}
	}

	b.WriteString("}")
	return b.String()
}

type Value struct {
	Attribute string
	Name      string
	Number    int
}

type Vendor struct {
	Name   string
	Number int

	TypeOctets   *int
	LengthOctets *int

	Attributes []*Attribute
	Values     []*Value
}

func (v *Vendor) GetTypeOctets() int {
	if v.TypeOctets == nil {
		return 1
	}
	return *v.TypeOctets
}

func (v *Vendor) GetLengthOctets() int {
	if v.LengthOctets == nil {
		return 1
	}
	return *v.LengthOctets
}

func (v *Vendor) GoString() string {
	var b bytes.Buffer
	b.WriteString("&dictionary.Vendor{")

	fmt.Fprintf(&b, "Name:%#v,", v.Name)
	fmt.Fprintf(&b, "Number:%#v,", v.Number)

	fmt.Fprintf(&b, "TypeOctets:%#v,", v.TypeOctets)
	fmt.Fprintf(&b, "LengthOctets:%#v,", v.LengthOctets)

	if len(v.Attributes) > 0 {
		b.WriteString("Attributes:[]*dictionary.Attribute{")
		for _, attr := range v.Attributes {
			fmt.Fprintf(&b, "%#v,", attr)
		}
		b.WriteString("},")
	}
	if len(v.Values) > 0 {
		b.WriteString("Values:[]*dictionary.Value{")
		for _, value := range v.Values {
			fmt.Fprintf(&b, "%#v,", value)
		}
		b.WriteString("},")
	}

	b.WriteString("}")
	return b.String()
}
