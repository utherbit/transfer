package simpleTypes

type Entity struct {
	ID   int
	name string

	intField   int
	int8Field  int8
	int16Field int16
	int32Field int32
	int64Field int64

	uintField    uint
	uint8Field   uint8
	uint16Field  uint16
	uint32Field  uint32
	uint64Field  uint64
	float32Field float32
	float64Field float64

	complex64Field  complex64
	complex128Field complex128

	uintptrField uintptr
	byteField    byte
	boolField    bool

	runeField rune
}
