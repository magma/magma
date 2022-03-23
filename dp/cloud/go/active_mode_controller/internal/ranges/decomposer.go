package ranges

import (
	"sort"

	"magma/dp/cloud/go/active_mode_controller/internal/containers"
)

type Range struct {
	Begin int
	End   int
	Value int
}

type Point struct {
	Value int
	Pos   int
}

func DecomposeOverlapping(ranges []Range, minValue int) []Point {
	var res []Point
	points, values := getEndsAndValues(ranges)
	pq := &containers.PriorityQueue{}
	pq.Push(minValue)
	for _, p := range points {
		res = addPoint(res, Point{
			Value: pq.Top(),
			Pos:   p.pos,
		})
		if p.isEnd {
			pq.Delete(values[p.id].item)
		} else {
			values[p.id].item = pq.Push(values[p.id].val)
		}
	}
	return res
}

func getEndsAndValues(ranges []Range) ([]*rangeEnd, []rangeValue) {
	points := make([]*rangeEnd, 2*len(ranges))
	values := make([]rangeValue, len(ranges))
	for i, r := range ranges {
		points[2*i] = &rangeEnd{id: i, pos: r.Begin, isEnd: false}
		points[2*i+1] = &rangeEnd{id: i, pos: r.End, isEnd: true}
		values[i] = rangeValue{val: r.Value}
	}
	sort.Slice(points, func(i, j int) bool {
		if points[i].pos == points[j].pos {
			return points[i].isEnd
		}
		return points[i].pos < points[j].pos
	})
	return points, values
}

func addPoint(points []Point, p Point) []Point {
	i := len(points) - 1
	if i >= 0 && p.Pos == points[i].Pos {
		return points
	}
	if i >= 0 && p.Value == points[i].Value {
		points[i].Pos = p.Pos
		return points
	}
	return append(points, p)
}

type rangeValue struct {
	val  int
	item *containers.Item
}

type rangeEnd struct {
	id    int
	pos   int
	isEnd bool
}
