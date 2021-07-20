package midi

import (
	"testing"
	"fmt"
)


func TestReadSMFHeader(t *testing.T) {
	fmt.Print()
	file, filename := openTestFile(t, "a.mid")
	_, err := readSMFHeader(file)
	if err != nil {
		errmsg := "\nreadSMFHeader(\"%s\") returnd unexpected error"
		errmsg += "\n%s\n"
		t.Fatalf(errmsg, filename, err)
	}
	file, filename = openTestFile(t, "bad1.mid")
	_, err = readSMFHeader(file)
	if err == nil {
		errmsg := "\nreadSMFHeader did not return error for known bad file."
		t.Fatalf(errmsg)
	}
	
}
