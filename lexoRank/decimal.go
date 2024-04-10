package lexoRank

import (
	"lexorank-go/numeralSystems"
	"strings"
)

var (
	// get decimals by func to avoid nil pointer
	zeroDecimal       *Decimal
	oneDecimal        *Decimal
	eightDecimal      *Decimal
	minDecimal        *Decimal
	maxDecimal        *Decimal
	midDecimal        *Decimal
	initialMinDecimal *Decimal
	initialMaxDecimal *Decimal
)

func OneDecimal() *Decimal {
	if oneDecimal == nil {
		oneDecimal = NewDecimal("1", defaultSystem)
	}
	return oneDecimal
}

func ZeroDecimal() *Decimal {
	if zeroDecimal == nil {
		zeroDecimal = NewDecimal("0", defaultSystem)
	}
	return zeroDecimal
}

func EightDecimal() *Decimal {
	if eightDecimal == nil {
		eightDecimal = NewDecimal("8", defaultSystem)
	}
	return eightDecimal
}

func MinDecimal() *Decimal {
	if minDecimal == nil {
		minDecimal = ZeroDecimal()
	}
	return minDecimal
}

func MaxDecimal() *Decimal {
	if maxDecimal == nil {
		maxDecimal = NewDecimal("1000000", defaultSystem).Subtract(OneDecimal())
	}
	return maxDecimal
}

func InitialMaxDecimal() *Decimal {
	if initialMaxDecimal == nil {
		ch := defaultSystem.ToChar(defaultSystem.GetBase() - 2)
		initialMaxDecimal = NewDecimal(string(ch)+"00000", defaultSystem)
	}
	return initialMaxDecimal
}

func InitialMinDecimal() *Decimal {
	if initialMinDecimal == nil {
		initialMinDecimal = NewDecimal("100000", defaultSystem)
	}
	return initialMinDecimal
}

type Decimal struct {
	Itg *Integer
	sig int // hzzzzb:xxx, sig is the length of characters after the colon
}

// Format TODO this logic may be different from the original
func (d *Decimal) Format() string {
	intStr := d.Itg.Format()
	if d.sig == 0 {
		return intStr
	}

	head := string(intStr[0])
	ns := strings.Builder{}
	itgSys := d.Itg.GetSys()
	specialHead := head == itgSys.GetPositiveChar() || head == itgSys.GetNegativeChar()
	if specialHead {
		intStr = intStr[1:]
	}
	for ns.Len()+len(intStr) <= d.sig {
		ns.WriteByte('0')
	}
	ns.WriteString(intStr)
	idx := ns.Len() - d.sig
	ret := ns.String()
	ret = strings.Join([]string{ret[:idx], ret[idx:]}, itgSys.GetRadixPointChar())
	if specialHead {
		ret = head + ret
	}
	return ret
}

func (d *Decimal) GetSys() numeralSystems.NumeralSystem {
	return d.Itg.GetSys()
}

func (d *Decimal) GetSig() int {
	return d.sig
}

func (d *Decimal) SetSig(sig int, ceiling bool) *Decimal {
	if sig >= d.sig {
		return d
	}

	if sig < 0 {
		sig = 0
	}

	diff := d.sig - sig
	newItg := d.Itg.ShiftRight(diff)
	if ceiling {
		newItg = newItg.Add(one(newItg.GetSys()))
	}
	return createDecimal(newItg, sig)
}

func (d *Decimal) Equal(other *Decimal) bool {
	if other == nil {
		return false
	}
	if other == d {
		return true
	}
	return d.Itg.Equal(other.Itg) && d.sig == other.sig
}

func (d *Decimal) CompareTo(other *Decimal) int {
	if other == nil {
		return 1
	}
	if d == other {
		return 0
	}
	tMag := d.Itg
	oMag := other.Itg
	times := d.sig - other.sig
	if times > 0 {
		oMag = oMag.ShiftLeft(times)
	} else if times < 0 {
		tMag = tMag.ShiftLeft(-times)
	}
	return tMag.CompareTo(oMag)
}

func (d *Decimal) Subtract(other *Decimal) *Decimal {
	tm, ts := d.Itg, d.sig
	om, os := other.Itg, other.sig
	for ; ts < os; ts++ {
		tm = tm.ShiftLeft(1)
	}
	for ts > os {
		om = om.ShiftLeft(1)
		os++
	}
	return createDecimal(tm.Subtract(om), ts)
}

func (d *Decimal) Multiply(other *Decimal) *Decimal {
	return createDecimal(d.Itg.Multiply(other.Itg), d.sig+other.sig)
}

func (d *Decimal) Floor() *Integer {
	return d.Itg.ShiftRight(d.sig)
}

func (d *Decimal) isExact() bool {
	for i := 0; i < d.sig; i++ {
		if d.Itg.GetMag(i) != 0 {
			return false
		}
	}
	return true
}
func (d *Decimal) Ceil() *Integer {
	if d.isExact() {
		return d.Itg
	}
	return d.Floor().Add(createInteger(d.GetSys(), 1, []int{1}))
}

func (d *Decimal) Add(other *Decimal) *Decimal {
	tm := d.Itg
	tsig := d.sig
	om := other.Itg
	osig := other.sig
	if osig > tsig {
		tm = tm.ShiftLeft(osig - tsig)
		tsig = osig
	} else {
		om = om.ShiftLeft(tsig - osig)
		osig = tsig
	}
	return createDecimal(tm.Add(om), tsig)
}

func between(left, right *Decimal) *Decimal {
	left.Itg.checkSystem(right.Itg)
	l, r := left, right
	if left.GetSig() < right.GetSig() {
		nr := right.SetSig(left.GetSig(), true)
		if left.CompareTo(nr) >= 0 {
			return dmid(left, right)
		}
		right = nr
	}

	if left.GetSig() > right.GetSig() {
		// round up after removing the sig part
		nl := left.SetSig(right.GetSig(), true)
		if nl.CompareTo(right) >= 0 {
			return dmid(left, right)
		}
		left = nl
	}
	var nr *Decimal
	for sig := l.GetSig(); sig > 0; right = nr {
		nl := l.SetSig(sig-1, true)
		nr = right.SetSig(sig-1, false)
		cmp := nl.CompareTo(nr)
		if cmp == 0 {
			return checkMid(left, right, nl)
		}
		if nl.CompareTo(nr) > 0 {
			break
		}
		sig -= 1
		l = nl
	}

	mid := middleInternal(left, right, l, r)
	var ns int
	for sig := mid.GetSig(); sig > 0; sig = ns {
		ns = sig - 1
		nm := mid.SetSig(ns, false)
		if left.CompareTo(nm) >= 0 || nm.CompareTo(right) >= 0 {
			break
		}
		mid = nm
	}
	return mid
}
func NewDecimal(str string, system numeralSystems.NumeralSystem) *Decimal {
	if strings.Count(str, system.GetRadixPointChar()) > 1 {
		panic("more than one radix pointer")
	}
	partialIndex := strings.Index(str, system.GetRadixPointChar())
	// not existed
	if partialIndex == -1 {
		return createDecimal(NewInteger(str, system), 0)
	}
	newItg := NewInteger(str[:partialIndex]+str[partialIndex+1:], system)
	return createDecimal(newItg, len(str)-1-partialIndex)
}

func HalfDecimal(sys numeralSystems.NumeralSystem) *Decimal {
	mid := sys.GetBase()/2 | 0
	return createDecimal(createInteger(sys, 1, []int{mid}), 1)
}

func createDecimal(itg *Integer, sig int) *Decimal {
	if itg.IsZero() {
		return &Decimal{Itg: itg, sig: sig}
	}

	zeroCount := 0
	for zeroCount < sig && itg.GetMag(zeroCount) == 0 {
		zeroCount++
	}
	if zeroCount > 0 {
		zeroCount--
	}
	return &Decimal{Itg: itg.ShiftRight(zeroCount), sig: sig - zeroCount}
}
