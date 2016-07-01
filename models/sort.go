package models

import "sort"

// Int64Sorter - joins int64 slice and sorting function
type Int64Sorter struct {
	items []int64
	by    func(a, b int64) bool
}

// By is the type of a "less" function that defines the ordering of its int64 arguments
type By func(a, b int64) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(items []int64) {
	sorter := &Int64Sorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(sorter)
}

func (a Int64Sorter) Len() int { return len(a.items) }

func (a Int64Sorter) Swap(i, j int) {
	a.items[i], a.items[j] = a.items[j], a.items[i]
}

func (a Int64Sorter) Less(i, j int) bool {
	return a.by(a.items[i], a.items[j])
}

// GameSorter - joins Game slice and sorting function
type GameSorter struct {
	items []Game
	by    func(a, b Game) bool
}

// GameBy is the type of a "less" function that defines the ordering of its arguments
type GameBy func(a, b Game) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by GameBy) Sort(items []Game) {
	sorter := &GameSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(sorter)
}

func (a GameSorter) Len() int { return len(a.items) }

func (a GameSorter) Swap(i, j int) {
	a.items[i], a.items[j] = a.items[j], a.items[i]
}

func (a GameSorter) Less(i, j int) bool {
	return a.by(a.items[i], a.items[j])
}
