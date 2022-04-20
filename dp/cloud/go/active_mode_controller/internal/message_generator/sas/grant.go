package sas

import (
	"math"
	"sort"

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
	operationParam := chooseSuitableChannel(cbsd, g.rng)
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

	defaultBandwidthMHz = 20
)

func chooseSuitableChannel(cbsd *active_mode.Cbsd, rng RNG) *operationParam {
	calc := newEirpCalculator(cbsd.GetEirpCapabilities())
	pts := channelsToPoints(cbsd.GetChannels())
	preferred, other := getCandidates(cbsd.GetPreferences(), pts, calc)

	left := cbsd.GrantAttempts
	frequencies := cbsd.GetPreferences().GetFrequenciesMhz()
	for _, f := range frequencies {
		for i := range preferred {
			if len(preferred[i]) == 0 || preferred[i][0].Pos != toPos(f) {
				continue
			}
			if left == 0 {
				return newOperationParam(preferred[i][0], bandwidths[i], calc)
			}
			left--
			preferred[i] = preferred[i][1:]
		}
	}
	for i := range other {
		p, ok := ranges.Select(other[i], rng.Int(), 5*1e6)
		if !ok {
			continue
		}
		if left == 0 {
			return newOperationParam(p, bandwidths[i], calc)
		}
		left--
	}
	return nil
}

func channelsToPoints(channels []*active_mode.Channel) []ranges.Point {
	asRanges := make([]ranges.Range, len(channels))
	for i, c := range channels {
		val := maxSASEirp
		if c.MaxEirp != nil {
			val = int(math.Floor(float64(c.MaxEirp.Value)))
		}
		asRanges[i] = ranges.Range{
			Begin: int(c.LowFrequencyHz - lowestFrequencyHz),
			End:   int(c.HighFrequencyHz - lowestFrequencyHz),
			Value: val,
		}
	}
	return ranges.DecomposeOverlapping(asRanges, minSASEirp-1)
}

func getCandidates(preferences *active_mode.FrequencyPreferences, points []ranges.Point, calc *eirpCalculator) ([][]ranges.Point, [][]ranges.Range) {
	preferred := make([][]ranges.Point, len(bandwidths))
	other := make([][]ranges.Range, len(bandwidths))
	bandwidth := getBandwidthHZOrDefault(preferences.GetBandwidthMhz())
	frequencies := newOrderedFrequencies(preferences.GetFrequenciesMhz())
	for i, b := range bandwidths {
		if b > bandwidth {
			continue
		}
		low := int(calc.calcLowerBound(b, minSASEirp))
		all := ranges.ComposeForMidpoints(points, b, low)
		other[i], preferred[i] = ranges.Split(all, frequencies.points)
		frequencies.sort(preferred[i])
	}
	return preferred, other
}

func getBandwidthHZOrDefault(bandwidthMHz int32) int {
	b := int(bandwidthMHz)
	if b == 0 {
		b = defaultBandwidthMHz
	}
	return b * 1e6
}

func newOperationParam(p ranges.Point, band int, calc *eirpCalculator) *operationParam {
	return &operationParam{
		MaxEirp: calc.calcUpperBound(band, p.Value),
		OperationFrequencyRange: &frequencyRange{
			LowFrequency:  int64(p.Pos-band/2) + lowestFrequencyHz,
			HighFrequency: int64(p.Pos+band/2) + lowestFrequencyHz,
		},
	}
}

type orderedFrequencies struct {
	points  []int
	pointId map[int]int
}

func newOrderedFrequencies(frequencies []int32) *orderedFrequencies {
	p := &orderedFrequencies{
		points:  make([]int, len(frequencies)),
		pointId: make(map[int]int, len(frequencies)),
	}
	for i, f := range frequencies {
		pos := toPos(f)
		p.points = append(p.points, pos)
		p.pointId[pos] = i
	}
	sort.Ints(p.points)
	return p
}

func toPos(freqMHz int32) int {
	return int(freqMHz-int32(lowestFrequencyHz/1e6)) * 1e6
}

func (p orderedFrequencies) sort(pts []ranges.Point) {
	sort.Slice(pts, func(i, j int) bool {
		return p.pointId[pts[i].Pos] < p.pointId[pts[j].Pos]
	})
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
