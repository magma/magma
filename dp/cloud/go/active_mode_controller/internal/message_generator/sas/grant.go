package sas

import (
	"math"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type grantRequestGenerator struct {
	indexProvider ranges.IndexProvider
}

func NewGrantRequestGenerator(indexProvider ranges.IndexProvider) *grantRequestGenerator {
	return &grantRequestGenerator{
		indexProvider: indexProvider,
	}
}

func (g *grantRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	operationParam := chooseSuitableChannel(
		cbsd.GetChannels(),
		cbsd.GetEirpCapabilities(),
		g.indexProvider,
	)
	if operationParam == nil {
		return nil
	}
	req := &grantRequest{
		CbsdId:         cbsd.Id,
		OperationParam: operationParam,
	}
	return []*Request{asRequest(Grant, req)}
}

type grantRequest struct {
	CbsdId         string          `json:"cbsdId"`
	OperationParam *operationParam `json:"operationParam"`
}

type operationParam struct {
	MaxEirp                 float64         `json:"maxEirp"`
	OperationFrequencyRange *frequencyRange `json:"operationFrequencyRange"`
}

var bandwidths = [...]int{200, 150, 100, 50}

const (
	minSASEirp = -137
	maxSASEirp = 37
	tenthMHz   = 1e5
	deci       = 10
)

func chooseSuitableChannel(
	channels []*active_mode.Channel,
	capabilities *active_mode.EirpCapabilities,
	indexProvider ranges.IndexProvider,
) *operationParam {
	calc := newEirpCalculator(capabilities)
	rs := toRanges(channels)
	pts := ranges.Decompose(rs, minSASEirp-1)
	for _, band := range bandwidths {
		res := tryToGetChannelForBandwidth(calc, band, pts, indexProvider)
		if res != nil {
			return res
		}
	}
	return nil
}

func toRanges(channels []*active_mode.Channel) []ranges.Range {
	res := make([]ranges.Range, len(channels))
	for i, c := range channels {
		val := maxSASEirp
		if c.MaxEirp != nil {
			val = int(math.Floor(float64(c.MaxEirp.Value)))
		}
		res[i] = ranges.Range{
			Begin: int(c.FrequencyRange.Low / tenthMHz),
			End:   int(c.FrequencyRange.High / tenthMHz),
			Value: val,
		}
	}
	return res
}

func tryToGetChannelForBandwidth(
	calc *eirpCalculator,
	band int,
	pts []ranges.Point,
	provider ranges.IndexProvider,
) *operationParam {
	low := int(calc.calcLowerBound(band, minSASEirp))
	if len(pts) <= 1 {
		return nil
	}
	valid := ranges.FindAvailable(pts, band, low)
	if len(valid) == 0 {
		return nil
	}
	minFreq := ranges.SelectPoint(valid, provider)
	return &operationParam{
		MaxEirp: calc.calcUpperBound(band, minFreq.Value),
		OperationFrequencyRange: &frequencyRange{
			LowFrequency:  int64(minFreq.Pos) * tenthMHz,
			HighFrequency: int64(minFreq.Pos+band) * tenthMHz,
		},
	}
}

type eirpCalculator struct {
	minPower    float64
	maxPower    float64
	antennaGain float64
	noPorts     float64
}

func newEirpCalculator(capabilities *active_mode.EirpCapabilities) *eirpCalculator {
	return &eirpCalculator{
		minPower:    float64(capabilities.GetMinPower()),
		maxPower:    float64(capabilities.GetMaxPower()),
		antennaGain: float64(capabilities.GetAntennaGain()),
		noPorts:     float64(capabilities.GetNumberOfPorts()),
	}
}

func (e eirpCalculator) calcLowerBound(bandwidth int, min int) float64 {
	eirp := e.calcEirp(e.minPower, bandwidth)
	return math.Ceil(math.Max(eirp, float64(min)))
}

func (e eirpCalculator) calcUpperBound(bandwidth int, max int) float64 {
	eirp := e.calcEirp(e.maxPower, bandwidth)
	return math.Floor(math.Min(eirp, float64(max)))
}

func (e eirpCalculator) calcEirp(power float64, bandwidth int) float64 {
	bwMHz := float64(bandwidth / deci)
	return power + e.antennaGain - 10*math.Log10(bwMHz/e.noPorts)
}
