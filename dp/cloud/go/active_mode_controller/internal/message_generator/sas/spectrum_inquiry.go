package sas

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type spectrumInquiryRequestGenerator struct{}

func NewSpectrumInquiryRequestGenerator() *spectrumInquiryRequestGenerator {
	return &spectrumInquiryRequestGenerator{}
}

func (*spectrumInquiryRequestGenerator) GenerateRequests(config *active_mode.ActiveModeConfig) []*Request {
	req := &spectrumInquiryRequest{
		CbsdId: config.GetCbsd().GetId(),
		InquiredSpectrum: []*frequencyRange{{
			LowFrequency:  lowestFrequencyHz,
			HighFrequency: highestFrequencyHz,
		}},
	}
	return []*Request{asRequest(SpectrumInquiry, req)}
}

const (
	lowestFrequencyHz  int64 = 3550 * 1e6
	highestFrequencyHz int64 = 3700 * 1e6
)

type spectrumInquiryRequest struct {
	CbsdId           string            `json:"cbsdId"`
	InquiredSpectrum []*frequencyRange `json:"inquiredSpectrum"`
}

type frequencyRange struct {
	LowFrequency  int64 `json:"lowFrequency"`
	HighFrequency int64 `json:"highFrequency"`
}
