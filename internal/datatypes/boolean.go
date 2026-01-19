package datatypes

type Boolean byte

const TRUE_VALUE Boolean = 0x01
const FALSE_VALUE Boolean = 0x00

func NewBoolean(value bool) Boolean {
	if value {
		return TRUE_VALUE
	}
	return FALSE_VALUE
}
