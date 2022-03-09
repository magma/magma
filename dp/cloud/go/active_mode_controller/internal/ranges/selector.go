package ranges

func SelectPoint(ranges []Range, provider indexProvider) Point {
	sum := 0
	for _, r := range ranges {
		sum += r.End - r.Begin + 1
	}
	n := provider.Intn(sum)
	i := 0
	for ; n > ranges[i].End-ranges[i].Begin; i++ {
		n -= ranges[i].End - ranges[i].Begin + 1
	}
	return Point{
		Value: ranges[i].Value,
		Pos:   ranges[i].Begin + n,
	}
}

type indexProvider interface {
	Intn(int) int
}
