package amp

import "testing"

func TestAMP(t *testing.T) {

	args := [][]byte{[]byte("Hello"), []byte("World")}

	bin := Encode(args)
	msg := Decode(bin)

	if want, have := "Hello", string(msg[0]); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	if want, have := "World", string(msg[1]); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func TestAMPNoArgs(t *testing.T) {

	// it should support no args
	args := [][]byte{}

	bin := Encode(args)
	msg := Decode(bin)

	if len(msg) != 0 {
		t.Errorf("expected %s to be empty", msg)
	}
}

func BenchmarkEncode(b *testing.B) {
	args := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(args)
	}
}

func BenchmarkDecode(b *testing.B) {
	args := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	bin := Encode(args)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(bin)
	}
}
