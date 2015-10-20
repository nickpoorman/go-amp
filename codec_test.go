package amp_test

import (
	"fmt"
	"testing"

	. "github.com/nickpoorman/go-amp"
)

func TestAMP(t *testing.T) {

	args := [][]byte{[]byte("Hello"), []byte("World")}

	bin := Encode(args)
	msg := Decode(bin)

	if fmt.Sprintf("%s", msg[0]) != "Hello" {
		t.Errorf("expected %s to equal Hello", msg[0])
	}

	if fmt.Sprintf("%s", msg[1]) != "World" {
		t.Errorf("expected %s to equal World", msg[1])
	}
}
