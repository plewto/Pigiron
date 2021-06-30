// Hex dump utility
//
package hex

import "fmt"

const LineLength int = 16

func padHex(s string) string {
	tlen := 3 * LineLength + 8
	f := fmt.Sprintf("%%-%ds ", tlen)
	return fmt.Sprintf(f, s) + "  :  "
}


func formatLine(data []byte, start int) string {
	acc := fmt.Sprintf("[%4x] ", start)
	bcc := ""
	for i, j := start, 0; i < len(data) && j < LineLength; i, j = i+1, j+1 {
		value := data[i]
		acc += fmt.Sprintf("%02x ", value)
		if 0x20 <= value && value < 0x7F {
			bcc += fmt.Sprintf("%c", value)
		} else {
			bcc += "-"
		}
	}
	return padHex(acc) + bcc
}
		

// HexDump prints formated values from byte slice.
//
func HexDump(data []byte, start int, end int) {
	if end <= 0 {
		end = len(data)
	}
	for i := start; i < end; i += LineLength {
		fmt.Println(formatLine(data, i))
	}
} 





