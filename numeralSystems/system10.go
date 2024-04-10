package numeralSystems

const (
	Digits10 = ""
)

var _ NumeralSystem = &System10{}

type System10 struct {
}

func (s System10) GetBase() int {
	return 10
}

func (s System10) GetPositiveChar() string {
	return "+"
}

func (s System10) GetNegativeChar() string {
	return "-"
}

func (s System10) GetRadixPointChar() string {
	return "."
}

func (s System10) ToDigit(u uint8) int {
	if u >= '0' && u <= '9' {
		return int(u - '0')
	}
	panic("invalid digit")
}

func (s System10) ToChar(i int) byte {
	return byte(i + 'a')
}
