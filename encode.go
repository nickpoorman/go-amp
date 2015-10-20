package amp

/**
 * Protocol version.
 */

var version uint32 = 1

// PutUint32 encode unint32 into binary big endian
func PutUint32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

// Encode `args`.
func Encode(args [][]byte) []byte {
	argc := len(args)
	buffLen := 1

	// data length
	for i := 0; i < argc; i++ {
		buffLen += 4 + len(args[i])
	}

	// buffer
	buff := make([]byte, buffLen)

	// pack meta
	PutUint32(buff[0:1], ((version << 4) | uint32(argc)))
	buff = buff[1:]

	// pack args
	for i := 0; i < argc; i++ {
		arg := args[i]
		argLen := uint32(len(arg))

		PutUint32(buff[0:4], argLen)
		buff = buff[4:]

		copy(buff, arg)
		buff = buff[argLen:]
	}

	return buff
}
