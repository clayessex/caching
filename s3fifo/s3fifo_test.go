package s3fifo

import (
	"fmt"
	"testing"
)

func TestGeneral(t *testing.T) {
	x := S3fifo{}
	fmt.Printf("%v\n", x)
}

func (s *S3fifo) GoString() string {
	return fmt.Sprintf("Map:%v, s:%v, g:%v, size:%v, Ks: %v, Km: %v\n", s.data, s.s, s.g, s.size, s.Ks, s.Km)
}
