// Kdtree is a very simple K-D tree implementation.
// This implementation uses a fixed value for K.  The intention
// is to copy the code locally, change K to your needs, and
// change T.Data's type to suit your needs too.
package kdtree

import (
	"sort"
)

// K is the dimensionality of the points in this package's K-D trees.
const K = 2

// A Point is a location in K-dimensional space.
type Point [K]float64

// SqDist returns the square distance between two points.
func (a *Point) sqDist(b *Point) float64 {
	sqDist := 0.0
	for i, x := range a {
		diff := x - b[i]
		sqDist += diff * diff
	}
	return sqDist
}

// A T is a the node of a K-D tree.  A *T is the root of a K-D tree,
// and nil is an empty K-D tree.
type T struct {
	// Point is the K-dimensional point associated with the
	// data of this node.
	Point
	// Data is auxiliary data associated with the point of this node.
	Data interface{}

	split       int
	left, right *T
}

// Insert returns a new K-D tree with the given node inserted.
// Inserting a node that is already a member of a K-D tree
// invalidates that tree.
func (t *T) Insert(n *T) *T {
	return t.insert(0, n)
}

func (t *T) insert(depth int, n *T) *T {
	if t == nil {
		n.split = depth % K
		n.left, n.right = nil, nil
		return n
	}
	if n.Point[t.split] < t.Point[t.split] {
		t.left = t.left.insert(depth+1, n)
	} else {
		t.right = t.right.insert(depth+1, n)
	}
	return t
}

// InRange appends all nodes in the K-D tree that are within a given
// distance from the given point to the given slice, which may be nil.
// To  avoid allocation, the slice can be pre-allocated with a larger
// capacity and re-used across multiple calls to InRange.
func (t *T) InRange(pt Point, dist float64, nodes []*T) []*T {
	if dist < 0 {
		return nodes
	}
	return t.inRange(&pt, dist, nodes)
}

func (t *T) inRange(pt *Point, r float64, nodes []*T) []*T {
	if t == nil {
		return nodes
	}

	diff := pt[t.split] - t.Point[t.split]

	thisSide, otherSide := t.right, t.left
	if diff < 0 {
		thisSide, otherSide = t.left, t.right
		diff = -diff // abs
	}
	nodes = thisSide.inRange(pt, r, nodes)
	if diff <= r {
		if t.Point.sqDist(pt) < r*r {
			nodes = append(nodes, t)
		}
		nodes = otherSide.inRange(pt, r, nodes)
	}

	return nodes
}

// Height returns the height of the K-D tree.
func (t *T) Height() int {
	if t == nil {
		return 0
	}
	ht := t.left.Height()
	if rht := t.right.Height(); rht > ht {
		ht = rht
	}
	return ht + 1
}

// New returns a new K-D tree built using the given nodes.
// Building a new tree with nodes that are already members of
// K-D trees invalidates those trees.
func New(nodes []*T) *T {
	if len(nodes) == 0 {
		return nil
	}
	return buildTree(0, nodes)
}

// BuildTree returns a new tree, built up from the given slice of nodes.
func buildTree(depth int, nodes []*T) *T {
	split := depth % K
	switch len(nodes) {
	case 0:
		return nil
	case 1:
		nd := nodes[0]
		nd.split = split
		nd.left, nd.right = nil, nil
		return nd
	}
	cur, nodes := med(split, nodes)
	left, right := partition(split, cur.Point[split], nodes)
	cur.split = split
	cur.left = buildTree(depth+1, left)
	cur.right = buildTree(depth+1, right)
	return cur
}

// Partition returns two node slices, the first containing all elements
// with values less than that of the pivot on the split dimension, and the
// second containing all values greater or equal to that of the pivot
// on the splitting dimension.
func partition(split int, pivot float64, nodes []*T) (fst, snd []*T) {
	p := 0
	for i, nd := range nodes {
		if nd.Point[split] < pivot {
			nodes[p], nodes[i] = nodes[i], nodes[p]
			p++
		}
	}
	return nodes[:p], nodes[p:]
}

// Med returns the median node, compared on the split dimension
// and the remaining nodes.
func med(split int, nodes []*T) (*T, []*T) {
	if len(nodes) == 0 {
		panic("med: no nodes")
	}
	sort.Sort(nodeSorter{split, nodes})
	var m int
	for m = len(nodes) / 2; m >= 1; m-- {
		if nodes[m-1].Point[split] < nodes[m].Point[split] {
			break
		}
	}
	nodes[0], nodes[m] = nodes[m], nodes[0]
	return nodes[0], nodes[1:]
}

// A nodeSorter implements sort.Interface, sortnig the nodes
// in ascending order of their point values on the split dimension.
type nodeSorter struct {
	split int
	nodes []*T
}

func (n nodeSorter) Len() int {
	return len(n.nodes)
}

func (n nodeSorter) Swap(i, j int) {
	n.nodes[i], n.nodes[j] = n.nodes[j], n.nodes[i]
}

func (n nodeSorter) Less(i, j int) bool {
	return n.nodes[i].Point[n.split] < n.nodes[j].Point[n.split]
}
