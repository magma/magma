package ranges

import "magma/dp/cloud/go/active_mode_controller/internal/containers"

func FindAvailable(points []Point, length int, minValue int) []Range {
	var res []Range
	i, j := -1, 0
	last := points[len(points)-1].Pos
	pos := points[0].Pos - length
	mq := &containers.MinQueue{}
	mq.Push(points[0].Value)
	for {
		moveBegin := points[i+1].Pos <= points[j].Pos-length
		delta := 0
		if moveBegin {
			i++
			delta = points[i].Pos - pos
		} else {
			delta = points[j].Pos - pos - length
			j++
		}
		if v := mq.Top(); v >= minValue {
			res = addRange(res, Range{
				Begin: pos,
				End:   pos + delta,
				Value: v,
			})
		}
		pos += delta
		if !moveBegin && pos+length == last {
			break
		}
		if moveBegin {
			mq.Pop()
		} else {
			mq.Push(points[j].Value)
		}
	}
	return res
}

func addRange(ranges []Range, r Range) []Range {
	i := len(ranges) - 1
	if i >= 0 &&
		ranges[i].End == r.Begin &&
		r.Value == ranges[i].Value {
		ranges[i].End = r.End
		return ranges
	}
	return append(ranges, r)
}
