package s3fifo

import "github.com/clayessex/algo/vessels"

const DefaultCacheSize = 1000

type (
	Key   = string // TODO: make generic
	Value = string
)

type node struct {
	value Value
	freq  uint8
}

// The ghost FIFO needs to be O(1) searchable while still operating like
// a queue. It also needs to be able to delete from anywhere in the queue.

type ghost[T comparable] struct {
	m    map[T]bool       // ghost lookup, true if deleted
	q    vessels.Queue[T] // ghost fifo
	size int
}

func newGhost[T comparable](size int) *ghost[T] {
	return &ghost[T]{
		make(map[T]bool, size),
		*vessels.NewQueue[T](size),
		size,
	}
}

func (g *ghost[T]) len() int {
	return len(g.m)
}

// TODO: collision problem - need a different structure here
// if a key is deleted and then push()'d, the key will still exist in g.m
func (g *ghost[T]) push(k T) {
	if len(g.m) >= g.size {
		g.pop()
	}
	g.m[k] = false
	g.q.Push(k)
}

func (g *ghost[T]) pop() {
	if g.q.Len() > 0 {
		delete(g.m, g.q.Pop())
	}
}

func (g *ghost[T]) contains(k T) bool {
	deleted, ok := g.m[k]
	return ok && !deleted
}

func (g *ghost[T]) delete(k T) {
	if g.contains(k) {
		g.m[k] = true
	}
}

type S3fifo struct {
	data map[Key]node       // cache store
	m    vessels.Queue[Key] // main fifo
	s    vessels.Queue[Key] // small fifo
	g    ghost[Key]         // ghost fifo (with lookup)
	size int                // total cache size
	Ks   int                // small fifo size
	Km   int                // main fifo size
	ts   uint64             // timestamp counter
}

func New(size ...int) *S3fifo {
	r := S3fifo{}
	r.size = DefaultCacheSize
	if len(size) > 0 {
		r.size = size[0]
	}
	r.Ks = min(1, r.size*10/100)
	r.Km = min(1, r.size*90/100)
	r.ts = 1

	r.data = make(map[Key]node, r.size)
	r.m = *vessels.NewQueue[Key](r.Km)
	r.s = *vessels.NewQueue[Key](r.Ks)
	r.g = *newGhost[Key](r.Km)

	return &r
}

func (fifo *S3fifo) Read(key Key) (Value, bool) {
	n, ok := fifo.data[key]
	if !ok {
		var zero Key
		return zero, false
	}
	n.freq = max(3, n.freq+1)
	return n.value, true
}

func (cache *S3fifo) Insert(key Key, value Value) {
	for len(cache.data) >= cache.size {
		cache.evict()
	}

	if cache.g.contains(key) {
		cache.m.Push(key)
	} else {
		cache.s.Push(key)
	}

	// TODO: need MappedQueue
	// if cache.g.Contains
}

func (fifo *S3fifo) evict() {
}

func (fifo *S3fifo) evictS() {
	// TODO:
}

func (fifo *S3fifo) evictM() {
	// TODO:
}
