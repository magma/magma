package ranges

import "sort"

func CheckIfContain(ranges []Range, points []int) []bool {
	res := make([]bool, len(points))
	ordered := make([]indexedPoint, len(points))
	for i, p := range points {
		ordered[i] = indexedPoint{id: i, pos: p}
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].pos < ordered[j].pos
	})
	j := 0
	for _, p := range ordered {
		for ; j < len(ranges) && p.pos > ranges[j].End; j++ {
		}
		if j == len(ranges) {
			break
		}
		if p.pos >= ranges[j].Begin {
			res[p.id] = true
		}
	}
	return res
}

type indexedPoint struct {
	id  int
	pos int
}
