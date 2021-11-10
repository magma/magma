package message_generator

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

const (
	defaultMaxEirp float32 = 37
	eirpOffset     float32 = 10
)

type grantRequest struct {
	CbsdId         string          `json:"cbsdId"`
	OperationParam *OperationParam `json:"operationParam"`
}

func (*grantRequest) name() string {
	return "grant"
}

type OperationParam struct {
	MaxEirp                 float32         `json:"maxEirp"`
	OperationFrequencyRange *frequencyRange `json:"operationFrequencyRange"`
}

type grantMessageGenerator struct{}

func (*grantMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	cbsd := config.GetCbsd()
	channel := cbsd.GetChannels()[0]
	maxEirp := choseMaxEirp(cbsd.EirpCapability, channel)
	if maxEirp <= 0 {
		return nil
	}
	req := &grantRequest{
		CbsdId: cbsd.Id,
		OperationParam: &OperationParam{
			MaxEirp: maxEirp,
			OperationFrequencyRange: &frequencyRange{
				LowFrequency:  channel.GetFrequencyRange().GetLow(),
				HighFrequency: channel.GetFrequencyRange().GetHigh(),
			},
		},
	}
	return []message{req}
}

func choseMaxEirp(eirpCapability *float32, channel *active_mode.Channel) float32 {
	if channel.LastEirp != nil {
		return channel.GetLastEirp() - 1
	}
	maxEirp := defaultMaxEirp
	if channel.MaxEirp != nil {
		maxEirp = min(maxEirp, channel.GetMaxEirp())
	}
	if eirpCapability != nil {
		maxEirp = min(maxEirp, *eirpCapability-eirpOffset)
	}
	return maxEirp
}

func min(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
