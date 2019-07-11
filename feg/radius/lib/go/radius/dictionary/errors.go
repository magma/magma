package dictionary

import "strconv"

type ParseError struct {
	Inner error
	File  File
	Line  int
}

func (e *ParseError) Error() string {
	str := "dictionary: parse error in " + e.File.Name() + ":" + strconv.Itoa(e.Line)
	if e.Inner != nil {
		str += ": " + e.Inner.Error()
	}
	return str
}

type DuplicateAttributeError struct {
	Attribute *Attribute
}

func (e *DuplicateAttributeError) Error() string {
	return `duplicate attribute "` + e.Attribute.Name + `"`
}

type UnknownAttributeTypeError struct {
	Type string
}

func (e *UnknownAttributeTypeError) Error() string {
	return `unknown attribute type "` + e.Type + "`"
}

type DuplicateAttributeFlagError struct {
	Flag string
}

func (e *DuplicateAttributeFlagError) Error() string {
	return `duplicate attribute flag "` + e.Flag + `"`
}

type UnknownAttributeFlagError struct {
	Flag string
}

func (e *UnknownAttributeFlagError) Error() string {
	return `unknown attribute flag "` + e.Flag + `"`
}

type InvalidAttributeEncryptTypeError struct {
	Type string
}

func (e *InvalidAttributeEncryptTypeError) Error() string {
	return `invalid attribute encrypt type "` + e.Type + `"`
}

type UnknownLineError struct {
	Line string
}

func (e *UnknownLineError) Error() string {
	return `unknown line`
}

type InvalidVendorFormatError struct {
	Format string
}

func (e *InvalidVendorFormatError) Error() string {
	return `invalid vendor format "` + e.Format + `"`
}

type UnknownVendorError struct {
	Vendor string
}

func (e *UnknownVendorError) Error() string {
	return `unknown vendor "` + e.Vendor + `"`
}

type UnmatchedEndVendorError struct {
}

func (e *UnmatchedEndVendorError) Error() string {
	return `unmatched END-VENDOR`
}

type InvalidEndVendorError struct {
	Vendor string
}

func (e *InvalidEndVendorError) Error() string {
	return `invalid END-VENDOR "` + e.Vendor + `"`
}

type BeginVendorIncludeError struct {
}

func (e *BeginVendorIncludeError) Error() string {
	return `invalid $INCLUDE inside BEGIN-VENDOR block`
}

type UnclosedVendorBlockError struct {
}

func (e *UnclosedVendorBlockError) Error() string {
	return `unclosed BEGIN-VENDOR block`
}

type RecursiveIncludeError struct {
	Filename string
}

func (e *RecursiveIncludeError) Error() string {
	return `file already included "` + e.Filename + `"`
}

type DuplicateVendorError struct {
	Vendor *Vendor
}

func (e *DuplicateVendorError) Error() string {
	return `duplicate vendor "` + e.Vendor.Name + `" (` + strconv.Itoa(e.Vendor.Number) + `)`
}

type NestedVendorBlockError struct {
}

func (e *NestedVendorBlockError) Error() string {
	return `invalid BEGIN-VENDOR inside vendor block`
}

type UnsupportedNestedTLVError struct {
}

func (e *UnsupportedNestedTLVError) Error() string {
	return `Nested tlv type is not supported`
}
