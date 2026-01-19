package datatypes

const SEGMENT_BITS byte = 0x7F // 01111111
const CONTINUE_BIT byte = 0x80 // 10000000

type Byte int8
type UnsignedByte uint8
type Short int16
type UnsignedShort uint16
type Int int32
type Long int64
type Float float32
type Double float64
