package lexoRank

import (
	"errors"
	"lexorank-go/numeralSystems"
	"strings"
)

const (
	pipe = "|"
)

var (
	defaultSystem = &numeralSystems.System36{}
)

func init() {

}

type Rank struct {
	numeralSystems.System36
	Value   string
	Bucket  *Bucket
	Decimal *Decimal
}

func (r *Rank) From(bk *Bucket, dl *Decimal) (*Rank, error) {
	if dl.GetSys().GetBase() != r.GetBase() {
		return nil, errors.New("system not match")
	}
	return &Rank{Bucket: bk, Decimal: dl}, nil
}

func (r *Rank) Min() (*Rank, error) {
	return r.From(buckets[0], minDecimal)
}

func (r *Rank) Equal(other *Rank) bool {
	if r == other {
		return true
	}
	if other == nil {
		return false
	}
	return r.Value == other.Value
}

func (r *Rank) Between(other *Rank) *Rank {
	if !r.Bucket.Equal(other.Bucket) {
		panic("between only works within the same bucket")
	}
	cmp := r.Decimal.CompareTo(other.Decimal)
	if cmp == 0 {
		panic("no enough space between two rank")
	}
	if cmp > 0 {
		return from(r.Bucket, between(other.Decimal, r.Decimal))
	}
	return from(r.Bucket, between(r.Decimal, other.Decimal))
}

func (r *Rank) Middle() *Rank {
	minRank := MinRank()
	return minRank.Between(MaxRank(minRank.Bucket))
}

func (r *Rank) isMin() bool {
	return r.Decimal.Equal(MinDecimal())
}

func (r *Rank) isMax() bool {
	return r.Decimal.Equal(MaxDecimal())
}

func (r *Rank) GenNext() *Rank {
	if r.isMin() {
		return from(r.Bucket, InitialMinDecimal())
	}
	ceilDecimal := createDecimal(r.Decimal.Ceil(), 0)
	nextDecimal := ceilDecimal.Add(EightDecimal())
	if nextDecimal.CompareTo(MaxDecimal()) >= 0 {
		nextDecimal = between(r.Decimal, MaxDecimal())
	}
	return from(r.Bucket, nextDecimal)
}

func (r *Rank) GenPrev() *Rank {
	if r.isMax() {
		return from(r.Bucket, InitialMaxDecimal())
	}
	floorDecimal := createDecimal(r.Decimal.Floor(), 0)
	prev := floorDecimal.Subtract(EightDecimal())
	if prev.CompareTo(MinDecimal()) <= 0 {
		prev = between(MinDecimal(), r.Decimal)
	}
	return from(r.Bucket, prev)
}

func (r *Rank) NextBucket() *Rank {
	return from(r.Bucket.Next(), r.Decimal)
}

func (r *Rank) PrevBucket() *Rank {
	return from(r.Bucket.Prev(), r.Decimal)
}

func (r *Rank) Format() string {
	return r.Value
}

func dmid(left, right *Decimal) *Decimal {
	// mid = (left+right) * 1/2
	sum := left.Add(right)
	mid := sum.Multiply(HalfDecimal(left.GetSys()))
	sig := left.GetSig()
	if sig < right.GetSig() {
		sig = right.GetSig()
	}
	if mid.GetSig() > sig {
		roundDown := mid.SetSig(sig, false)
		if roundDown.CompareTo(left) > 0 {
			return roundDown
		}
		roundUp := mid.SetSig(sig, true)
		if roundUp.CompareTo(right) < 0 {
			return roundUp
		}
	}
	return mid
}

func MinRank() *Rank {
	return from(buckets[0], MinDecimal())
}

func MaxRank(bucket *Bucket) *Rank {
	if bucket == nil {
		bucket = buckets[0]
	}
	return from(bucket, MaxDecimal())
}

func middle() *Rank {
	minRank := MinRank()
	return minRank.Between(MaxRank(minRank.Bucket))
}

func initial(bucket *Bucket) *Rank {
	if bucket == buckets[0] {
		return from(bucket, InitialMinDecimal())
	}
	return from(bucket, InitialMaxDecimal())
}

func middleInternal(low, high, left, right *Decimal) *Decimal {
	return checkMid(low, high, dmid(left, right))
}

func checkMid(low, high, mid *Decimal) *Decimal {
	if low.CompareTo(mid) >= 0 {
		return dmid(low, high)
	}
	if mid.CompareTo(high) >= 0 {
		return dmid(low, high)
	}
	return mid
}

func from(bucket *Bucket, decimal *Decimal) *Rank {
	if decimal.GetSys().GetBase() != defaultSystem.GetBase() {
		panic("unexpected different system")
	}
	return &Rank{Bucket: bucket, Decimal: decimal, Value: strings.Join([]string{bucket.Format(), formatDecimal(decimal)}, pipe)}
}

func NewRank(str string) *Rank {
	parts := strings.Split(str, pipe)
	if len(parts) != 2 {
		panic("invalid rank string")
	}
	bucket, err := NewBucket(parts[0])
	if err != nil {
		panic(err)
	}
	decimal := NewDecimal(parts[1], defaultSystem)
	return from(bucket, decimal)
}

func formatDecimal(decimal *Decimal) string {
	ns := strings.Builder{}
	fv := decimal.Format()
	pi := strings.Index(fv, defaultSystem.GetRadixPointChar())
	if pi == -1 {
		pi = len(fv)
		fv += defaultSystem.GetRadixPointChar()
	}
	for pi < 6 {
		ns.WriteByte(defaultSystem.ToChar(0))
		pi++
	}
	ns.WriteString(fv)
	return strings.TrimRight(ns.String(), "0")
}
