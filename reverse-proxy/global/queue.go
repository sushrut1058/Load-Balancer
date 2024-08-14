package global

func (q *Queue) Push(element interface{}) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.items = append(q.items, element)
}

func (q *Queue) Pop() (interface{}, bool) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	if len(q.items) != 0 {
		item := q.items[0]
		q.items = q.items[1:]
		return item, true
	}
	return nil, false
}

func (q *Queue) Size() int {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return len(q.items)
}
