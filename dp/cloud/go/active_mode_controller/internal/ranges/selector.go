package ranges

func Select(ranges []Range, n int, mod int) (Point, bool) {
	sum := 0
	for i := range ranges {
		sum += countDivisible(ranges[i].Begin, ranges[i].End, mod)
	}
	if sum == 0 {
		return Point{}, false
	}
	i := 0
	for n %= sum; n >= 0; i++ {
		n -= countDivisible(ranges[i].Begin, ranges[i].End, mod)
	}
	i--
	return Point{
		Value: ranges[i].Value,
		Pos:   (ranges[i].End/mod + n + 1) * mod,
	}, true
}

func countDivisible(a int, b int, mod int) int {
	x := b/mod - a/mod
	if a%mod == 0 {
		x++
	}
	return x
}
