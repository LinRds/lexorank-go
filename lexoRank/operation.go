package lexoRank

import "lexorank-go/numeralSystems"

// Array-based arithmetic operations in base-n numeral system
// digits are stored in reverse order
func subtract(sys numeralSystems.NumeralSystem, left, right []int) []int {
	cp := complement(sys, right, len(left))
	rs := add(sys, left, cp)
	rs[len(rs)-1] = 0
	return add(sys, rs, []int{1})
}

func multiply(sys numeralSystems.NumeralSystem, left, right []int) []int {
	ret := make([]int, len(left)+len(right))
	base := sys.GetBase()
	for i := range left {
		for j := range right {
			k := i + j
			ret[k] += left[i] * right[j]
			ret[k+1] += ret[k] / base
			ret[k] %= base
		}
	}
	return ret
}

func complement(sys numeralSystems.NumeralSystem, arr []int, digits int) []int {
	if digits <= 0 {
		panic("expected at least 1 digit")
	}
	bm := sys.GetBase() - 1
	newArr := make([]int, max(digits, len(arr)))
	for i := range newArr {
		newArr[i] = bm
	}
	for i := 0; i < len(arr); i++ {
		newArr[i] = bm - arr[i]
	}
	return newArr
}

func compare(left, right []int) int {
	ls := len(left)
	rs := len(right)
	if ls < rs {
		return -1
	} else if rs < ls {
		return 1
	}
	for i := ls - 1; i >= 0; i-- {
		if left[i] < right[i] {
			return -1
		} else if left[i] > right[i] {
			return 1
		}
	}
	return 0
}

func add(sys numeralSystems.NumeralSystem, left, right []int) []int {
	ls := len(left)
	rs := len(right)
	if ls < rs {
		return add(sys, right, left)
	}
	ret := make([]int, ls)
	carry := 0
	base := sys.GetBase()
	var i int
	for i = range left {
		r := 0
		if i < rs {
			r = right[i]
		}
		carry = r + left[i] + carry
		ret[i] = carry % base
		carry /= base
	}
	if carry != 0 {
		ret = append(ret, carry)
	}
	return ret
}
