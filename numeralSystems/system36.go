package numeralSystems

var _ NumeralSystem = &System36{}

const (
	Digits36 = "0123456789abcdefghijklmnopqrstuvwxyz"
)

type System36 struct{}

func (s System36) GetBase() int {
	return 36
}

func (s System36) GetPositiveChar() string {
	return "+"
}

func (s System36) GetNegativeChar() string {
	return "-"
}

func (s System36) GetRadixPointChar() string {
	return ":"
}

func (s System36) ToDigit(u uint8) int {
	if u >= '0' && u <= '9' {
		return int(u - '0')
	}
	if u >= 'a' && u <= 'z' {
		return int(u-'a') + 10
	}
	panic("invalid digit")
}

func (s System36) ToChar(i int) byte {
	return Digits36[i]
}
