package smf

import (
	// "fmt"
	"testing"
	gomidi "gitlab.com/gomidi/midi/v2"
)


func TestIsTextType(t *testing.T) {
	if isTextType(0x00) || isTextType(0x08) {
		errmsg := "isTextType returns false positive for invalid meta type."
		t.Fatalf(errmsg)
	}
	if !isTextType(0x03) {
		errmsg := "expected true for isTextType(0x03), got false"
		t.Fatalf(errmsg)
	}
}

func TestIsTextMessage(t *testing.T) {
	good := gomidi.NewMessage([]byte{0xFF, 0x01, 0x03, 0x41, 0x42, 0x43})
	bad1 := gomidi.NewMessage([]byte{0x80, 0x00, 0x00})
	bad2 := gomidi.NewMessage([]byte{0xFF, 0x2F, 0x00})
	if !IsTextMessage(good) {
		errmsg := "IsTextMessage returnd false for meta text message"
		t.Fatalf(errmsg)
	}
	if IsTextMessage(bad1) {
		errmsg := "IsTextMessage returned true for non-meta message."
		t.Fatalf(errmsg)
	}
	if IsTextMessage(bad2) {
		errmsg := "IsTextMessage returned true for non-text meta message."
		t.Fatalf(errmsg)
	}
}

func TestText(t *testing.T) {
	msg, _ := MakeTextMessage(0x02, "ABC")
	text, txType, err := ExtractMetaText(msg)
	if err != nil {
		errmsg := "Text returned unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	if text != "ABC" {
		errmsg := "Expected text \"%s\", got \"%s\""
		t.Fatalf(errmsg, "ABC", text)
	}
	if txType != 0x02 {
		errmsg := "Expected text type 0x02, got 0x%02X"
		t.Fatalf(errmsg, txType)
	}
}
	
// func TestSplitTime(t *testing.T) {
// 	fmt.Println()
// 	fmt.Println(FormatTime(3600.0 + 72.123))
// }
