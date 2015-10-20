package amp

/**
 * Protocol version.
 */

var version = 1

// Encode `args`.
func Encode(args [][]byte) []byte {
	argc := len(args)
	bufl := 1

	// data length
	for i := 0; i < argc; i++ {
		bufl += 4 + len(args[i])
	}

	// buffer
	buf := make([]byte, bufl)

	// pack meta
	buf[0] = byte((version << 4) | argc)
	buf = buf[1:]

	// pack args
	for i := 0; i < argc; i++ {
		arg := args[i]
		argl := uint32(len(arg))

		buf[0] = byte(argl >> 24)
		buf[1] = byte(argl >> 16)
		buf[2] = byte(argl >> 8)
		buf[3] = byte(argl)
		buf = buf[4:]

		copy(buf, arg)
		buf = buf[argl:]
	}

	return buf
}
