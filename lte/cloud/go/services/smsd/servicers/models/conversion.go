package models

import (
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/lte/cloud/go/services/smsd/storage"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func (m *SmsMessage) FromProto(from *storage.SMS) *SmsMessage {
	m.Pk = from.Pk
	m.Imsi = models.SubscriberID(from.Imsi)
	m.SourceMsisdn = from.SourceMsisdn
	m.AttemptCount = int64(from.AttemptCount)
	m.Message = from.Message

	m.TimeCreated = tsToDT(from.CreatedTime)
	lastAttempt := tsToDT(from.LastDeliveryAttemptTime)
	if lastAttempt != nil {
		m.TimeLastAttempted = *lastAttempt
	}

	switch from.Status {
	case storage.MessageStatus_WAITING:
		m.Status = strPtr(SmsMessageStatusWaiting)
	case storage.MessageStatus_DELIVERED:
		m.Status = strPtr(SmsMessageStatusDelivered)
	case storage.MessageStatus_FAILED:
		m.Status = strPtr(SmsMessageStatusFailed)
	default:
		m.Status = strPtr(SmsMessageStatusWaiting)
	}

	return m
}

func (m *MutableSmsMessage) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableSmsMessage) ToProto() storage.MutableSMS {
	return storage.MutableSMS{
		Imsi:         string(m.Imsi),
		SourceMsisdn: m.SourceMsisdn,
		Message:      m.Message,
	}
}

func tsToDT(ts *timestamp.Timestamp) *strfmt.DateTime {
	if ts == nil {
		return nil
	}

	ret, err := ptypes.Timestamp(ts)
	if err != nil {
		return nil
	}
	dt := strfmt.DateTime(ret)
	return &dt
}

func strPtr(s string) *string {
	return &s
}
