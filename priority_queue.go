// This file contains a min-heap implementation of a priority queue
// using type Node as elements of the min-heap.

package main

type MinHeap []Node

// Len returns the current length of the min-heap
func (m *MinHeap) Len() int {
	return len(*m)
}

// Less defines a comparison operator for elements, uses the evaluation of G and H for that node
func (m *MinHeap) Less(i, j int) bool {
	return (*m)[i].g+(*m)[i].h < (*m)[j].g+(*m)[j].h
}

// Swap swaps elements at index i and j
func (m *MinHeap) Swap(i, j int) {
	tmp := (*m)[i]
	(*m)[i] = (*m)[j]
	(*m)[j] = tmp
}

// Parent returns the parent of i in an array implementation of a heap
func Parent(i int) int {
	return i >> 1
}

// Left returns the left child of i in an array implementation of a heap
func Left(i int) int {
	return i << 1
}

// Right returns the index of element i's right child in an array implementation of a heap
func Right(i int) int {
	return (i << 1) + 1
}

// MinHeapify enforces the min-heap property for the min-heap rooted at element i
func (m *MinHeap) MinHeapify(i int) {
	l := Left(i)
	r := Right(i)

	var smallest int
	if l < m.Len() && m.Less(l, i) {
		smallest = l
	} else {
		smallest = i
	}

	if r < m.Len() && m.Less(r, smallest) {
		smallest = r
	}

	if smallest != i {
		m.Swap(i, smallest)
		m.MinHeapify(smallest)
	}
}

// Pop removes and returns the min element of the heap
func (m *MinHeap) Pop() Node {
	min := (*m)[0]            // get min element
	(*m)[0] = (*m)[m.Len()-1] // put last element first
	*m = (*m)[:m.Len()-1]     // decrease size of array
	if m.Len() > 0 {
		m.MinHeapify(0) // sort first element back in heap
	}
	return min
}

// Insert inserts node n into the min-heap
func (m *MinHeap) Insert(n Node) {
	*m = append(*m, n)
	i := len(*m) - 1 // get index of newly inserted element
	// move newly inserted element up the heap
	for i >= 1 && m.Less(i, Parent(i)) {
		m.Swap(i, Parent(i))
		i = Parent(i)
	}

}

// Initialize sets up the min-heap invariant for the elements in m
func (m *MinHeap) Initialize() {
	for i := m.Len() / 2; i >= 0; i-- {
		m.MinHeapify(i)
	}
}
