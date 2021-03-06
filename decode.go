package amp

// Decode the given `buf`.
func Decode(buf []byte) [][]byte {
	// unpack meta
	meta := uint8(buf[0])
	// version := meta >> 4
	// if version != ProtocolVersion {
	// 	panic("AMP protocol version does not equal 1.")
	// }
	argv := int(meta & 0xf)
	args := make([][]byte, argv)
	buf = buf[1:]

	for i := 0; i < argv; i++ {
		argl := uint32(buf[3]) | uint32(buf[2])<<8 | uint32(buf[1])<<16 | uint32(buf[0])<<24
		buf = buf[4:]

		args[i] = buf[0:argl]

		buf = buf[argl:]
	}

	return args
}
