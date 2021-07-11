package smf

import (
	"fmt"
	"testing"
)

var (
	mockBuffer []byte = []byte {
		byte('A'), byte('B'), byte('C'), byte('D'),  // [ 0] chunk ID 'ABCD'
		0x00, 0x00, 0x00, 0x00,                      // [ 4] long 0
		0x04, 0x03, 0x02, 0x01,                      // [ 8] long 0x4030201
		0x00, 0x00,                                  // [12] short 0x0
		0x02, 0x01,                                  // [14] short 0x201
		0x00, 0xff,                                  // [16] byte 0, 0xff
		
		0x00,                                        // [18] vlq 0x00
		0x7f,                                        // [19] vlq 0x7f
		0x81, 0x00,                                  // [20] vlq 0x80

		0xFF, 0x7f,                                  // [22] vlq 0x3fff
		0x81, 0x80, 0x00,                            // [24] vlq 0x4000
		0xc0, 0x80, 0x80, 0x00,                      // [27] vlq 0x08000000
		0xff, 0xff, 0xff, 0x7f,                      // [31] vlq 0x0FFFFFFF
	}
)


func TestExpectID(t *testing.T) {
	fmt.Print()
	var target1 = [4]byte{'A','B','C','X'}
	var target2 = [4]byte{'A','B','C','D'}
	var err error
	err = expectChunkID(mockBuffer, 0, target1)
	if err == nil {
		t.Fatal("expectChunckID false positive")
	}
	err = expectChunkID(mockBuffer, 0, target2)
	if err != nil {
		t.Fatal("expectChunckID false negative")
	}
	err = expectChunkID(mockBuffer, 33, target1)
	if err == nil {
		t.Fatal("expectChunckID did not detect index out of bounds")
	}
}



func TestGetLong(t *testing.T) {
	var err error
	var n int
	n, err = getLong(mockBuffer, 4)
	if err != nil || n != 0 {
		msg := "getLong faild at index 4,  n = %d, err = %s"
		t.Fatal(msg, n, err)
	}
	n, err = getLong(mockBuffer, 8)
	if err != nil || n != 0x4030201 {
		msg := "getLong fails at index 8, n = 0x%x, err = %s"
		t.Fatal(msg, n, err)
	}
	_, err = getLong(mockBuffer, 32) 
	if err == nil {
		t.Fatal("getLong faild to detect index out of bounds at index 32")
	}
}

func TestGetShort(t *testing.T) {
	var err error
	var n int
	n, err = getShort(mockBuffer, 12)
	if err != nil || n != 0 {
		msg := "getShort fails at index 12, n = 0x%x, err = %s"
		t.Fatal(msg, n, err)
	}
	n, err = getShort(mockBuffer, 14)
	if err != nil || n != 0x201 {
		msg := "getShort fails at index 14, n = 0x%x, err = %s"
		t.Fatal(msg, n, err)
	}
	_, err = getShort(mockBuffer, 34)
	if err == nil {
		msg := "getShort faild to detect index out of bounds, index = 34"
		t.Fatal(msg)
	}	
}

func TestGetByte(t *testing.T) {
	var err error
	var n byte
	n, err = getByte(mockBuffer, 16)
	if err != nil || n != 0x0 {
		msg := "getByte fails at index 16, n = 0x%x, err = %s"
		t.Fatal(msg, n, err)
	}
	n, err = getByte(mockBuffer, 17)
	if err != nil || n != 0xff {
		msg := "getByte fails at index 17, n = 0x%x, err = %s"
		t.Fatal(msg, n, err)
	}
	_, err = getByte(mockBuffer, 40)
	if err == nil {
		t.Fatal("getByte fails to detect index out of bounds, index = 40")
	}
}

func TestVLQ(t *testing.T) {

	type expect struct {
		index int
		value int
		length int
	}

	var test = []expect{
		expect{18, 0x00, 1},
		expect{19, 0x7f, 1},
		expect{20, 0x80, 2},
		expect{22, 0x3fff, 2},
		expect{24, 0x4000, 3},
		expect{27, 0x08000000, 4},
		expect{31, 0x0FFFFFFF, 4},
	}

	var vlq *VLQ
	var err error
	var n int
	for _, ex := range test {
		index := ex.index
		vlq, err = getVLQ(mockBuffer, index)
		if err != nil {
			msg := "getVLQ returned false err at index %d, err was %s"
			s := fmt.Sprintf(msg, index, err)
			t.Fatal(s)
		}
		n = vlq.Value()
		if n != ex.value {
			msg := "getVLQ returnd incorrect value at index %d, expected %d, got %d"
			s := fmt.Sprintf(msg, index, ex.value, n)
			t.Fatal(s)
		}
	}
	_, err = getVLQ(mockBuffer, 35)
	if err == nil {
		t.Fatal("getVLQ did not detect index out of bounds, index = 35")
	}
}
