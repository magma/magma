package message_generator

import (
	"math"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

const (
	minSASEirp = -137
	maxSASEirp = 37
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
	operationParam := chooseSuitableChannel(cbsd.GetChannels(), cbsd.GetEirpCapabilities())
	if operationParam == nil {
		return nil
	}
	req := &grantRequest{
		CbsdId:         cbsd.Id,
		OperationParam: operationParam,
	}
	return []message{req}
}

func chooseSuitableChannel(
	channels []*active_mode.Channel,
	capabilities *active_mode.EirpCapabilities,
) *OperationParam {
	for _, channel := range channels {
		maxEirp, ok := choseMaxEirp(channel, capabilities)
		if !ok {
			continue
		}
		frequency := channel.GetFrequencyRange()
		return &OperationParam{
			MaxEirp: float32(maxEirp),
			OperationFrequencyRange: &frequencyRange{
				LowFrequency:  frequency.GetLow(),
				HighFrequency: frequency.GetHigh(),
			},
		}
	}
	return nil
}

func choseMaxEirp(channel *active_mode.Channel, capabilities *active_mode.EirpCapabilities) (float64, bool) {
	minEirp, maxEirp := calculateEirpBounds(channel, capabilities)
	v := maxEirp
	if channel.LastEirp != nil {
		v = float64(channel.GetLastEirp().Value - 1)
	}
	if v < minEirp {
		return 0, false
	}
	return v, true
}

func calculateEirpBounds(channel *active_mode.Channel, capabilities *active_mode.EirpCapabilities) (float64, float64) {
	frequencyRange := channel.GetFrequencyRange()
	partialPower := calculatePartialPower(
		frequencyRange.GetLow(), frequencyRange.GetHigh(),
		capabilities.GetAntennaGain(), capabilities.GetNumberOfPorts(),
	)
	minCapableEirp := calculateEirp(partialPower, capabilities.GetMinPower())
	maxCapableEirp := calculateEirp(partialPower, capabilities.GetMaxPower())
	minEirp := math.Max(minSASEirp, minCapableEirp)
	maxEirp := math.Min(maxSASEirp, maxCapableEirp)
	if channel.MaxEirp != nil {
		maxEirp = math.Min(float64(channel.GetMaxEirp().Value), maxEirp)
	}
	return minEirp, maxEirp
}

func calculatePartialPower(minFreqHz int64, maxFreqHz int64, antennaGain float32, numberOfPorts int32) float64 {
	bandwidthMHz := float64((maxFreqHz - minFreqHz) / 1e6)
	ports64 := float64(numberOfPorts)
	gain64 := float64(antennaGain)
	return gain64 - 10*math.Log10(bandwidthMHz/ports64)
}

func calculateEirp(partialPower float64, power float32) float64 {
	power64 := float64(power)
	return math.Floor(power64 + partialPower)
}
