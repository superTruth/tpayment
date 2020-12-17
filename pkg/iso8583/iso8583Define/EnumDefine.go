package iso8583Define

// field value type
type FieldValueType int

const (
	Alpha FieldValueType = iota
	Number
	Hex
	EBCDIC
)

// field align type
type FieldAlignType int

const (
	Right FieldAlignType = iota
	Left
)
