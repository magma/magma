package ranges

func Split(ranges []Range, points []int) ([]Range, []Point) {
	var newRanges []Range
	var newPoints []Point
	i := 0
	for _, r := range ranges {
		end := r.End
		for ; i < len(points) && points[i] <= end; i++ {
			if points[i] < r.Begin {
				continue
			}
			newPoints = append(newPoints, Point{
				Pos:   points[i],
				Value: r.Value,
			})
			r.End = points[i] - 1
			newRanges = addRangeIfValid(newRanges, r)
			r.Begin = points[i] + 1
		}
		r.End = end
		newRanges = addRangeIfValid(newRanges, r)
	}
	return newRanges, newPoints
}

func addRangeIfValid(ranges []Range, r Range) []Range {
	if r.Begin > r.End {
		return ranges
	}
	return append(ranges, r)
}
