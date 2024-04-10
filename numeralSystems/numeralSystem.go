package numeralSystems

type NumeralSystem interface {
	GetBase() int
	GetPositiveChar() string
	GetNegativeChar() string
	GetRadixPointChar() string
	ToDigit(uint8) int
	ToChar(int) byte
}
