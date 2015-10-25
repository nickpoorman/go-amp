package amp

import (
	"bytes"
	"sync"
	"testing"
)

func TestAMPStream(t *testing.T) {

	// it should emit data events
	data := make(chan bytes.Buffer)

	parser := NewParser(data)

	a := Encode([][]byte{[]byte("tobi")})
	b := Encode([][]byte{[]byte("loki"), []byte("abby")})
	c := Encode([][]byte{[]byte("manny"), []byte("luna"), []byte("ewald")})

	// wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		n := 0

		for buf := range data {
			msgBytes := Decode(buf.Bytes())
			msg := mapToString(msgBytes)

			switch n {
			case 0:
				n++
				if want, have := len(msg), 1; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				if want, have := "tobi", msg[0]; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
			case 1:
				n++
				if want, have := len(msg), 2; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				if want, have := "loki", msg[0]; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				if want, have := "abby", msg[1]; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
			case 2:
				n++
				if want, have := len(msg), 3; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				if want, have := "manny", msg[0]; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				if want, have := "luna", msg[1]; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				if want, have := "ewald", msg[2]; want != have {
					t.Errorf("want %#v, have %#v", want, have)
				}
				return
			}
		}
	}()

	write(a, parser)
	write(b, parser)
	parser.Write(c[0:5])
	parser.Write(c[5:])

	wg.Wait()
}

func mapToString(msgBytes [][]byte) []string {
	var msg []string
	for _, argBytes := range msgBytes {
		msg = append(msg, string(argBytes))
	}
	return msg
}

func write(from []byte, to *Parser) {
	for i := 0; i < len(from); i++ {
		to.Write(from[i : i+1])
	}

}
