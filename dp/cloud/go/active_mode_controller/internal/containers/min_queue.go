package containers

type MinQueue []*minQueueItem

func (m *MinQueue) Push(value int) {
	i := len(*m) - 1
	cnt := 1
	for ; i >= 0 && (*m)[i].value >= value; i-- {
		cnt += (*m)[i].count
	}
	item := &minQueueItem{value: value, count: cnt}
	*m = append((*m)[:i+1], item)
}

func (m *MinQueue) Pop() {
	(*m)[0].count--
	if (*m)[0].count == 0 {
		*m = (*m)[1:]
	}
}

func (m MinQueue) Top() int {
	return m[0].value
}

type minQueueItem struct {
	value int
	count int
}
