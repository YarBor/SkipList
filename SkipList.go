package SkipList

type Interface interface {
	Less(other interface{}) bool
}

type SkipList struct {
	head   *Element
	tail   *Element
	update []*Element
	rank   []int
	length int
	level  int
}

func New() *SkipList {
	return &SkipList{
		head:   newElement(MaxLevel, nil),
		tail:   nil,
		update: make([]*Element, MaxLevel),
		rank:   make([]int, MaxLevel),
		length: 0,
		level:  1,
	}
}

func (sl *SkipList) Init() *SkipList {
	sl.head = newElement(MaxLevel, nil)
	sl.tail = nil
	sl.update = make([]*Element, MaxLevel)
	sl.rank = make([]int, MaxLevel)
	sl.length = 0
	sl.level = 1
	return sl
}

func (sl *SkipList) Front() *Element {
	return sl.head.level[0].prev
}

func (sl *SkipList) Back() *Element {
	return sl.tail
}

func (sl *SkipList) Len() int {
	return sl.length
}

func (sl *SkipList) Insert(v Interface) *Element {
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {

		if i == sl.level-1 {
			sl.rank[i] = 0
		} else {
			sl.rank[i] = sl.rank[i+1]
		}
		for x.level[i].prev != nil && x.level[i].prev.Value.Less(v) {
			sl.rank[i] += x.level[i].span
			x = x.level[i].prev
		}
		sl.update[i] = x
	}

	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			sl.rank[i] = 0
			sl.update[i] = sl.head
			sl.update[i].level[i].span = sl.length
		}
		sl.level = level
	}

	x = newElement(level, v)
	for i := 0; i < level; i++ {
		x.level[i].prev = sl.update[i].level[i].prev
		sl.update[i].level[i].prev = x

		x.level[i].span = sl.update[i].level[i].span - sl.rank[0] + sl.rank[i]
		sl.update[i].level[i].span = sl.rank[0] - sl.rank[i] + 1
	}

	for i := level; i < sl.level; i++ {
		sl.update[i].level[i].span++
	}

	if sl.update[0] == sl.head {
		x.backward = nil
	} else {
		x.backward = sl.update[0]
	}
	if x.level[0].prev != nil {
		x.level[0].prev.backward = x
	} else {
		sl.tail = x
	}
	sl.length++

	return x
}

func (sl *SkipList) deleteElement(e *Element, update []*Element) {
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].prev == e {
			update[i].level[i].span += e.level[i].span - 1
			update[i].level[i].prev = e.level[i].prev
		} else {
			update[i].level[i].span -= 1
		}
	}

	if e.level[0].prev != nil {
		e.level[0].prev.backward = e.backward
	} else {
		sl.tail = e.backward
	}

	for sl.level > 1 && sl.head.level[sl.level-1].prev == nil {
		sl.level--
	}
	sl.length--
}

func (sl *SkipList) Remove(e *Element) interface{} {
	x := sl.find(e.Value)
	if x == e && !e.Value.Less(x.Value) {
		sl.deleteElement(x, sl.update)
		return x.Value
	}

	return nil
}

func (sl *SkipList) Delete(v Interface) interface{} {
	x := sl.find(v)
	if x != nil && !v.Less(x.Value) {
		sl.deleteElement(x, sl.update)
		return x.Value
	}

	return nil
}

func (sl *SkipList) Find(v Interface) *Element {
	x := sl.find(v)
	if x != nil && !v.Less(x.Value) {
		return x
	}

	return nil
}

func (sl *SkipList) find(v Interface) *Element {
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].prev != nil && x.level[i].prev.Value.Less(v) {
			x = x.level[i].prev
		}
		sl.update[i] = x
	}

	return x.level[0].prev
}

func (sl *SkipList) GetRank(v Interface) int {
	x := sl.head
	rank := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].prev != nil && x.level[i].prev.Value.Less(v) {
			rank += x.level[i].span
			x = x.level[i].prev
		}
		if x.level[i].prev != nil && !x.level[i].prev.Value.Less(v) && !v.Less(x.level[i].prev.Value) {
			rank += x.level[i].span
			return rank
		}
	}

	return 0
}

func (sl *SkipList) GetElementByRank(rank int) *Element {
	x := sl.head
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].prev != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].prev
		}
		if traversed == rank {
			return x
		}
	}

	return nil
}
