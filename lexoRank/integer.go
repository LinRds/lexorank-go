package lexoRank

import (
	"lexorank-go/numeralSystems"
	"strings"
)

const (
	less    = -1
	greater = 1
	equal   = 0
)

type Integer struct {
	sys  numeralSystems.NumeralSystem
	sign int // 0, 1, -1
	mag  []int
}

func (itg *Integer) IsZero() bool {
	return itg.sign == 0 && len(itg.mag) == 1 && itg.mag[0] == 0
}

func (itg *Integer) IsOne() bool {
	return itg.sign == 1 && itg.IsOneish()
}

func (itg *Integer) IsOneish() bool {
	return len(itg.mag) == 1 && itg.mag[0] == 1
}

func (itg *Integer) GetSys() numeralSystems.NumeralSystem {
	return itg.sys
}

func (itg *Integer) GetSign() int {
	return itg.sign
}

func (itg *Integer) GetMag(index int) int {
	return itg.mag[index]
}

func (itg *Integer) Format() string {
	if itg.IsZero() {
		return "0"
	}

	ns := strings.Builder{}
	for i := len(itg.mag) - 1; i >= 0; i-- {
		ns.WriteByte(itg.sys.ToChar(itg.mag[i]))
	}
	if itg.sign == -1 {
		ns.WriteString(itg.sys.GetNegativeChar())
	}
	return ns.String()
}

// ShiftRight caution: number is stored in reverse order
func (itg *Integer) ShiftRight(times int) *Integer {
	if len(itg.mag) <= times {
		return zero(itg.sys)
	}

	newMag := make([]int, len(itg.mag)-times)
	copy(newMag, itg.mag[times:])
	return createInteger(itg.sys, itg.sign, newMag)
}

func (itg *Integer) ShiftLeft(times int) *Integer {
	if times == 0 {
		return itg
	}
	if times < 0 {
		return itg.ShiftRight(-times)
	}
	newMag := make([]int, times+len(itg.mag))
	copy(newMag[times:], itg.mag)
	return createInteger(itg.sys, itg.sign, newMag)
}

func (itg *Integer) Equal(other *Integer) bool {
	if other == nil {
		return false
	}
	if itg == other {
		return true
	}
	return itg.sys.GetBase() == other.sys.GetBase() && itg.CompareTo(other) == 0
}

func (itg *Integer) CompareTo(other *Integer) int {
	if itg == other {
		return 0
	}
	if other == nil {
		return 1
	}
	switch itg.sign {
	case -1:
		if other.sign == -1 {
			return other.sign * itg.compare(other)
		}
		return -1
	case 1:
		if other.sign == 1 {
			return itg.compare(other)
		}
		return 1
	default:
		return -other.sign
	}
}

func (itg *Integer) compare(other *Integer) int {
	sub := len(itg.mag) - len(other.mag)
	if sub < 0 {
		return -1
	} else if sub > 0 {
		return 1
	}
	for j := len(itg.mag) - 1; j >= 0; j-- {
		if itg.mag[j] < other.mag[j] {
			return -1
		}
		if itg.mag[j] > other.mag[j] {
			return 1
		}
	}
	return 0
}

func (itg *Integer) Negate() *Integer {
	if itg.IsZero() {
		return itg
	}
	sign := -1
	if itg.sign != 1 {
		sign = 1
	}
	return createInteger(itg.sys, sign, itg.mag)
}

func (itg *Integer) Subtract(other *Integer) *Integer {
	itg.checkSystem(other)
	if itg.IsZero() {
		return other.Negate()
	}
	if other.IsZero() {
		return itg
	}
	if itg.sign != other.sign {
		// -3-5=-(-(-3)+5)
		if itg.sign == -1 {
			return itg.Negate().Add(other).Negate()
		}
		return itg.Add(other.Negate())
	}
	cmp := compare(itg.mag, other.mag)
	if cmp == 0 {
		return zero(itg.sys)
	}
	if cmp < 0 {
		return createInteger(itg.sys, -itg.sign, subtract(itg.sys, other.mag, itg.mag))
	} else {
		return createInteger(itg.sys, itg.sign, subtract(itg.sys, itg.mag, other.mag))
	}
}

func (itg *Integer) Multiply(other *Integer) *Integer {
	itg.checkSystem(other)
	if itg.IsZero() {
		return itg
	}
	if other.IsZero() {
		return other
	}
	sign := 1
	if itg.sign != other.sign {
		sign = -1
	}
	multiOne := func(a *Integer) *Integer {
		return createInteger(itg.sys, sign, a.mag)
	}
	if itg.IsOneish() {
		return multiOne(other)
	}
	if other.IsOneish() {
		return multiOne(itg)
	}
	newMag := multiply(itg.sys, itg.mag, other.mag)
	return createInteger(itg.sys, sign, newMag)
}

func (itg *Integer) Add(other *Integer) *Integer {
	itg.checkSystem(other)
	if itg.IsZero() {
		return other
	}
	if other.IsZero() {
		return itg
	}
	if itg.sign != other.sign {
		if itg.sign == -1 {
			return other.Subtract(itg.Negate())
		} else {
			return itg.Add(other.Negate())
		}
	}
	ret := add(itg.sys, itg.mag, other.mag)
	return createInteger(itg.sys, itg.sign, ret)
}

func (itg *Integer) checkSystem(other *Integer) {
	if itg.sys != other.sys {
		panic("numeral sys not matched")
	}
}

func NewInteger(strFull string, system numeralSystems.NumeralSystem) *Integer {
	if len(strFull) == 0 {
		panic("empty string to parse")
	}
	sign := 1
	if strings.Index(strFull, system.GetPositiveChar()) == 0 {
		strFull = strFull[1:]
	} else if strings.Index(strFull, system.GetNegativeChar()) == 0 {
		strFull = strFull[1:]
		sign = -1
	}

	n := len(strFull)
	mag := make([]int, n)
	// reverse order
	for i := 0; i < len(mag); i++ {
		mag[i] = system.ToDigit(strFull[n-1-i])
	}
	return createInteger(system, sign, mag)
}

func createInteger(system numeralSystems.NumeralSystem, sign int, mag []int) *Integer {
	lastNoneZero := len(mag) - 1
	for lastNoneZero >= 0 {
		if mag[lastNoneZero] != 0 {
			break
		}
		lastNoneZero--
	}

	if lastNoneZero == -1 {
		return zero(system)
	}

	if lastNoneZero == len(mag)-1 {
		return &Integer{
			sys:  system,
			sign: sign,
			mag:  mag,
		}
	}

	newMag := make([]int, lastNoneZero+1)
	for i := range newMag {
		newMag[i] = mag[i]
	}
	return &Integer{sys: system, sign: sign, mag: newMag}
}

func zero(system numeralSystems.NumeralSystem) *Integer {
	return &Integer{sys: system, sign: 0, mag: []int{0}}
}

func one(sys numeralSystems.NumeralSystem) *Integer {
	return createInteger(sys, 1, []int{1})
}
