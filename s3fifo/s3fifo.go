package s3fifo

import "github.com/clayessex/algo/vessels"

const DefaultCacheSize = 1000

type (
	Key   = string
	Value = string
)

type node struct {
	key   Key
	value Value
	freq  uint8
}

type S3fifo struct {
	data map[Key]*node // cache store
	m    *vessels.Queue[*node]
	s    *vessels.Queue[*node]
	g    vessels.Set[Key]
	size int
	Ks   int
	Km   int
}

func New(size ...int) *S3fifo {
	r := S3fifo{}
	r.size = DefaultCacheSize
	if len(size) > 0 {
		r.size = size[0]
	}
	r.Ks = min(1, r.size*10/100)
	r.Km = min(1, r.size*90/100)

	r.data = make(map[string]*node, r.size)
	r.m = vessels.NewQueue[*node](r.Ks)
	r.s = vessels.NewQueue[*node](r.Ks)
	r.g = vessels.NewSet[Key]()

	return &r
}

func (fifo *S3fifo) Read(key Key) (Value, bool) {
	n, ok := fifo.data[key]
	if !ok {
		return "", false
	}
	n.freq = max(3, n.freq+1)
	return n.value, true
}

func (fifo *S3fifo) Insert(key Key, value Value) {
	for len(fifo.data) > fifo.size { // full cache
		fifo.evict()
	}
	n := node{key, value, 0}
	if fifo.g.Contains(key) {
		fifo.m.Push(&n)
	} else {
		fifo.s.Push(&n)
	}
}

func (fifo *S3fifo) evict() {
	if fifo.s.Len() >= fifo.Ks {
		fifo.evictS()
	} else {
		fifo.evictM()
	}
}

func (fifo *S3fifo) evictS() {
	// TODO:
}

func (fifo *S3fifo) evictM() {
	// TODO:
}
