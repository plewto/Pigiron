package smf


/*
 * Defines functions for meta text messages.
*/

import (
	"fmt"
	"strings"
	gomidi "gitlab.com/gomidi/midi/v2"
)

// isTextType returns true iff textType is a meta text type byte.
// 
func isTextType(textType byte) bool {
	return 0x01 <= textType && textType <= 0x07
}

// IsTextMessage returns true if message is a meta text message.
// All text bearing messages return true.
//
func IsTextMessage(msg gomidi.Message) bool {
	d := msg.Data
	return len(d) > 3 && d[0] == 0xFF && isTextType(d[1])
}

// MakeTextMessage creates a new meta message bearing text.
// args:
//   textType must be one of:
//       0x01 text
//       0x02 copyright
//       0x03 track name
//       0x04 instrument name
//       0x05 lyric
//       0x06 marker
//       0x07 cuepoint
//
func MakeTextMessage(textType byte, text string) (msg gomidi.Message, err error) {
	if !isTextType(textType) {
		errmsg := "0x%02X is not a valid Meta Text type"
		err = fmt.Errorf(errmsg, textType)
		return
	}
	vlq := NewVLQ(len(text))
	data := make([]byte, 2, 2 + len(text) + vlq.Length())
	data[0] = 0xFF
	data[1] = textType
	for _, b := range vlq.Bytes() {
		data = append(data, b)
	}
	for _, b := range []byte(text) {
		data = append(data, b)
	}
	msg = gomidi.NewMessage(data)
	return 
}


// ExtractMetaText returns tex contents of meta text message.
//
func ExtractMetaText(msg gomidi.Message) (text string, txType byte, err error) {
	if !IsTextMessage(msg) {
		errmsg := "Expected meta text message, got %v"
		err = fmt.Errorf(errmsg, err)
		return
	}
	var start int
	_, start, err = ExpectVLQ(msg.Data, 2)
	if err != nil {
		errmsg := "Could not read vlq for meta text message: %v\n%s"
		err = fmt.Errorf(errmsg, msg, err)
		return
	}
	d := msg.Data
	text = string(d[start: len(d)])
	txType = d[1]
	return
}
			
	
func SplitTime(seconds float64) (hr int, min int, sec int, fsec float64) {
	hr = int(seconds / 3600)
	seconds -= float64(hr * 3600)
	min = int(seconds / 60)
	seconds -= float64(min * 60)
	sec = int(seconds)
	fsec = float64(sec) - seconds

	return hr, min, sec, fsec
}
	
func FormatTime(seconds float64) string {
	hr, min, sec, fsec := SplitTime(seconds)
	f := fmt.Sprintf("%f", fsec)
	pos := strings.Index(f, ".")
	return fmt.Sprintf("%02d:%02d:%02d%s", hr, min, sec, f[pos:])
}
