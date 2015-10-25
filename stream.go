package amp

import "bytes"

// Parser is the underlying struct which holds the state while reading the stream
type Parser struct {
	data    chan<- bytes.Buffer
	state   string
	lenbuf  []byte // the length of the argument (4 bytes) (a buffer to hold argl)
	version uint8  // the protocol version
	argv    int    // the number of arguments
	bufs    *bytes.Buffer
	nargs   int    // number of arguments received
	leni    int    // a pointer to the current byte in the argument length header
	arglen  uint32 // the total number of bytes in the argument (argl)
	argcur  uint32 // the number of bytes in the argument received so far
}

// Write will read the bytes and then write them to the data channel
func (p *Parser) Write(chunk []byte) (n int, err error) {
	chunkl := uint32(len(chunk))
	var i uint32
	for ; i < chunkl; i++ {
		switch p.state {
		case "message":
			meta := uint8(chunk[i])
			p.version = meta >> 4
			p.argv = int(meta & 0xf)
			p.state = "arglen"
			p.bufs = &bytes.Buffer{}
			p.bufs.WriteByte(meta)
			p.nargs = 0
			p.leni = 0

		case "arglen":
			p.lenbuf[p.leni] = chunk[i]
			p.leni++

			// done, we've got the whole length of the argument now
			if 4 == p.leni {
				p.arglen = uint32(p.lenbuf[3]) | uint32(p.lenbuf[2])<<8 | uint32(p.lenbuf[1])<<16 | uint32(p.lenbuf[0])<<24
				p.bufs.Write(p.lenbuf)
				p.argcur = 0
				p.state = "arg"
			}

		case "arg":
			// bytes remaining in the argument
			rem := p.arglen - p.argcur

			// consume the chunk we need to complete
			// the argument, or the remainder of the
			// chunk if it's not mixed-boundary
			pos := min(rem+i, chunkl)

			// slice arg chunk
			part := chunk[i:pos]
			p.bufs.Write(part)

			// check if we have the complete arg
			p.argcur += pos - i
			done := p.argcur == p.arglen
			i = pos - 1

			if done {
				p.nargs++
			}

			// no more args
			if p.nargs == p.argv {
				p.state = "message"
				p.data <- *p.bufs
				break
			}

			if done {
				p.state = "arglen"
				p.leni = 0
			}
		}
	}
	return 0, nil
}

// NewParser will take in a stream and write out to a channel
func NewParser(data chan<- bytes.Buffer) *Parser {
	lenbuf := make([]byte, 4)

	return &Parser{
		data:   data,
		state:  "message",
		lenbuf: lenbuf,
	}
}

func min(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}
