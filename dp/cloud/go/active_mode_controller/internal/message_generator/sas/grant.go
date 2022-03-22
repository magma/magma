package sas

import (
	"math"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type grantRequestGenerator struct {
	rng RNG
}

type RNG interface {
	Int() int
}

func NewGrantRequestGenerator(rng RNG) *grantRequestGenerator {
	return &grantRequestGenerator{
		rng: rng,
	}
}

func (g *grantRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	operationParam := chooseSuitableChannel(
		cbsd.GetChannels(),
		cbsd.GetEirpCapabilities(),
		int(cbsd.GetGrantAttempts()),
		g.rng,
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

var bandwidths = [...]int{20 * 1e6, 15 * 1e6, 10 * 1e6, 5 * 1e6}

const (
	minSASEirp = -137
	maxSASEirp = 37
)

func chooseSuitableChannel(
	channels []*active_mode.Channel,
	capabilities *active_mode.EirpCapabilities,
	attempts int,
	rng RNG,
) *operationParam {
	// More sophisticated check will be implemented
	// together with frequency preference
	if attempts > 0 {
		return nil
	}
	calc := newEirpCalculator(capabilities)
	rs := toRanges(channels)
	pts := ranges.DecomposeOverlapping(rs, minSASEirp-1)
	for _, band := range bandwidths {
		res := tryToGetChannelForBandwidth(calc, band, pts, rng)
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
			Begin: int(c.FrequencyRange.Low - lowestFrequencyHz),
			End:   int(c.FrequencyRange.High - lowestFrequencyHz),
			Value: val,
		}
	}
	return res
}

func tryToGetChannelForBandwidth(
	calc *eirpCalculator,
	band int,
	pts []ranges.Point,
	rng RNG,
) *operationParam {
	low := int(calc.calcLowerBound(band, minSASEirp))
	midpoints := ranges.ComposeForMidpoints(pts, band, low)
	p, ok := ranges.Select(midpoints, rng.Int(), 5*1e6)
	if !ok {
		return nil
	}
	return &operationParam{
		MaxEirp: calc.calcUpperBound(band, p.Value),
		OperationFrequencyRange: &frequencyRange{
			LowFrequency:  int64(p.Pos-band/2) + lowestFrequencyHz,
			HighFrequency: int64(p.Pos+band/2) + lowestFrequencyHz,
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
	bwMHz := float64(bandwidth / 1e6)
	return power + e.antennaGain - 10*math.Log10(bwMHz/e.noPorts)
}
