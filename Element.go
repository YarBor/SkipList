package SkipList

import (
	"math/rand"
)

const MaxLevel = 32
const Ratio = 4

type skipListLevel struct {
	prev *Element
	span int
}

type Element struct {
	Value    Interface
	backward *Element
	level    []*skipListLevel
}

func (e *Element) Next() *Element {
	return e.level[0].prev
}

func (e *Element) Prev() *Element {
	return e.backward
}

func newElement(level int, v Interface) *Element {
	slLevels := make([]*skipListLevel, level)
	for i := 0; i < level; i++ {
		slLevels[i] = new(skipListLevel)
	}

	return &Element{
		Value:    v,
		backward: nil,
		level:    slLevels,
	}
}

func randomLevel() int {
	level := 1
	for (rand.Int31()&0xFFFF)%Ratio == 0 {
		level += 1
	}

	if level < MaxLevel {
		return level
	} else {
		return MaxLevel
	}
}
