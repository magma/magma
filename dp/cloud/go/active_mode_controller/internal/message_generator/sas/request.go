package sas

import "encoding/json"

type Request struct {
	Data json.RawMessage
	Type RequestType
}
type RequestType uint8

const (
	Registration RequestType = iota
	SpectrumInquiry
	Grant
	Heartbeat
	Relinquishment
	Deregistration
)

func (r RequestType) String() string {
	var pref string
	switch r {
	case Registration:
		pref = "registration"
	case SpectrumInquiry:
		pref = "spectrumInquiry"
	case Grant:
		pref = "grant"
	case Heartbeat:
		pref = "heartbeat"
	case Relinquishment:
		pref = "relinquishment"
	case Deregistration:
		pref = "deregistration"
	}
	return pref + "Request"
}

func asRequest(requestType RequestType, data interface{}) *Request {
	b, _ := json.Marshal(data)
	return &Request{
		Type: requestType,
		Data: b,
	}
}
