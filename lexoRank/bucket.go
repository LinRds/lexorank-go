package lexoRank

import (
	"strconv"
)

var (
	buckets []*Bucket
)

func init() {
	buckets = make([]*Bucket, 3)
	for i := range buckets {
		buckets[i] = &Bucket{index: i}
	}
}

type Bucket struct {
	index int
}

func (b *Bucket) Format() string {
	return strconv.Itoa(b.index)
}

func (b *Bucket) Equal(other *Bucket) bool {
	if other == nil {
		return false
	}
	if b == other {
		return true
	}
	return b.index == other.index
}

func (b *Bucket) Next() *Bucket {
	return buckets[(b.index+1)%3]
}

func (b *Bucket) Prev() *Bucket {
	return buckets[(b.index+2)%3]
}

func NewBucket(str string) (*Bucket, error) {
	id, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	return &Bucket{index: id}, nil
}
