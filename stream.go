package amp

import "bytes"

// Parser is the underlying struct which holds the state while reading the stream
type Parser struct {
	data    chan<- bytes.Buffer
	state   string
	lenbuf  []byte // the length of the argument (4 bytes) (a buffer to hold argl)
	version uint8  // the protocol version
	argv    int    // the number of arguments
	// bufs    [][]byte // bufs (args)
	bufs   *bytes.Buffer
	nargs  int    // some sort of argument counter
	leni   int    // a pointer to the current byte in the argument length header
	arglen uint32 // the total number of bytes in the argument (argl)
	argcur uint32 // the number of bytes in the argument received so far
}

// Write will read the bytes and then write them to the data channel
func (p *Parser) Write(chunk []byte) (n int, err error) {
	chunkl := uint32(len(chunk))
	var i uint32
	for ; i < chunkl; i++ {
		switch p.state {
		case "message":
			// var meta = chunk[i];
			// this.version = meta >> 4;
			// this.argv = meta & 0xf;
			// this.state = 'arglen';
			// this._bufs = [new Buffer([meta])];
			// this._nargs = 0;
			// this._leni = 0;
			// break;
			meta := uint8(chunk[i])
			p.version = meta >> 4
			p.argv = int(meta & 0xf)
			p.state = "arglen"
			// p.bufs = make([][]byte, p.argv)
			p.bufs = &bytes.Buffer{}
			p.nargs = 0
			p.leni = 0
		case "arglen":
			// this._lenbuf[this._leni++] = chunk[i];
			//
			// // done
			// if (4 == this._leni) {
			//   this._arglen = this._lenbuf.readUInt32BE(0);
			//   var buf = new Buffer(4);
			//   buf[0] = this._lenbuf[0];
			//   buf[1] = this._lenbuf[1];
			//   buf[2] = this._lenbuf[2];
			//   buf[3] = this._lenbuf[3];
			//   this._bufs.push(buf);
			//   this._argcur = 0;
			//   this.state = 'arg';
			// }
			p.lenbuf[p.leni] = chunk[i]

			// done, we've got the whole length of the argument now
			if 4 == p.leni {
				p.arglen = uint32(p.lenbuf[3]) | uint32(p.lenbuf[2])<<8 | uint32(p.lenbuf[1])<<16 | uint32(p.lenbuf[0])<<24
				// lenbufCopy := make([]byte, 4)
				// lenbufCopy[0] = p.lenbuf[0]
				// lenbufCopy[1] = p.lenbuf[1]
				// lenbufCopy[2] = p.lenbuf[2]
				// lenbufCopy[3] = p.lenbuf[3]
				// p.bufs = append(p.bufs, lenbufCopy)
				p.bufs.Write(p.lenbuf)
				p.argcur = 0
				p.state = "arg"
			}

		case "arg":
			//     // bytes remaining in the argument
			// var rem = this._arglen - this._argcur;
			//
			// // consume the chunk we need to complete
			// // the argument, or the remainder of the
			// // chunk if it's not mixed-boundary
			// var pos = Math.min(rem + i, chunk.length);
			//
			// // slice arg chunk
			// var part = chunk.slice(i, pos);
			// this._bufs.push(part);
			//
			// // check if we have the complete arg
			// this._argcur += pos - i;
			// var done = this._argcur == this._arglen;
			// i = pos - 1;
			//
			// if (done) this._nargs++;
			//
			// // no more args
			// if (this._nargs == this.argv) {
			//   this.state = 'message';
			//   this.emit('data', Buffer.concat(this._bufs));
			//   break;
			// }
			//
			// if (done) {
			//   this.state = 'arglen';
			//   this._leni = 0;
			// }

			// bytes remaining in the argument
			rem := p.arglen - p.argcur

			// consume the chunk we need to complete
			// the argument, or the remainder of the
			// chunk if it's not mixed-boundary
			pos := min(rem+i, chunkl)

			// slice arg chunk
			part := chunk[i:pos]
			// p.bufs = append(p.bufs, part)
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
