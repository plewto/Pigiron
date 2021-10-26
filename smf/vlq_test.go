package smf

import (
	"fmt"
	"testing"
)

// compare 2 byte arrays
//
func cmp(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		t := b[i]
		if s != t {
			return false
		}
	}
	return true
}

// hex format byte array
//
func xformat(a []byte) string {
	acc := "{ "
	for i, v := range a {
		acc += fmt.Sprintf("0x%02X", v)
		if i == len(a) - 1 {
			acc += " }"
		} else {
			acc += ", "
		}
	}
	return acc
}

func TestVLQ(t *testing.T) {
	fmt.Println("TestVLQ")
	cases := make(map[int][]byte)
        cases[0x00000000] = []byte{0x00}
        cases[0x00000040] = []byte{0x40}
        cases[0x0000007F] = []byte{0x7F}
        cases[0x00000080] = []byte{0x81, 0x00}
        cases[0x00002000] = []byte{0xC0, 0x00}
        cases[0x00003FFF] = []byte{0xFF, 0x7F}
        cases[0x00004000] = []byte{0x81, 0x80, 0x00}
        cases[0x00100000] = []byte{0xC0, 0x80, 0x00}
        cases[0x001FFFFF] = []byte{0xFF, 0xFF, 0x7F}
        cases[0x00200000] = []byte{0x81, 0x80, 0x80, 0x00}
        cases[0x08000000] = []byte{0xC0, 0x80, 0x80, 0x00}
        cases[0x0FFFFFFF] = []byte{0xFF, 0xFF, 0xFF, 0x7F}

	for value, expect := range cases {
		vlq := NewVLQ(value)
		if !cmp(expect, vlq.Bytes()) {
			msg := "For VLQ value 0x%08X, expected bytes %s, got %s"
			t.Fatal(msg, value, expect, vlq.Bytes())
		}
	}

	for expect, bytes := range cases {
		vlq := NewVLQ(0)
		vlq.setBytes(bytes)
		if expect != vlq.Value() {
			msg := "For VLQ bytes %s, expected value 0x%08X, got 0x%08X"
			t.Fatal(msg, xformat(bytes), expect, vlq.Value())
		}
	}
}
	
