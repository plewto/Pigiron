package smf

import (
	"testing"
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigpath"
)


func TestReadSMFHeader(t *testing.T) {

	fmt.Println("*** EXPECT TO SEE WARNINGS ***")
	openTestFile := func(name string) (*os.File, string) {
		filename := pigpath.ResourceFilename("testFiles", name)
		file, err := os.Open(filename)
		if err != nil {
			errmsg := "\nCan not open test file: '%s'"
			errmsg += "\n%s\n"
			t.Fatalf(errmsg, filename, err)
		}
		fmt.Printf("Using test file '%s'\n", filename)
		return file, filename
	}
		
	file, filename := openTestFile("a1.mid")
	defer file.Close()
	header, err := readHeader(file)
	if err != nil {
		errmsg := "\nreadHeader(\"%s\") returned unexpected error"
		errmsg += "\n%s\n"
		t.Fatalf(errmsg, filename, err)
	}
	if header.format != 1 {
		errmsg := "\nexpected format 2, got %d"
		t.Fatalf(errmsg, header.format)
	}
	if header.trackCount != 1 {
		errmsg := "\nexpected trackCount 1, got %d"
		t.Fatalf(errmsg, header.trackCount)
	}
	if header.division != 24 {
		errmsg := "\nexpected division 24, got %d"
		t.Fatalf(errmsg, header.division)
	}

	
	// Malformed files

	file, filename = openTestFile("b1.mid")  
	defer file.Close()
	_, err = readHeader(file)
	if err == nil {
		errmsg := "\nreadHeader did not return error for non-midi file"
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}

	file, filename = openTestFile("b2.mid")
	defer file.Close()
	_, err = readHeader(file)
	if err == nil {
		errmsg := "\nreadHeader did not return error for unexpected header id."
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}

	file, filename = openTestFile("b3.mid")
	defer file.Close()
	_, err = readHeader(file)
	if err == nil {
		errmsg := "\nreadHeader did not detect wrong chunk length."
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}

	// malformed & recoverable
	file, filename = openTestFile("c1.mid")
	defer file.Close()
	header, err = readHeader(file)
	if err != nil {
		errmsg := "\nreadHeader returned error for recoverable file, invalid format."
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}
	if header.format != 0 {
		errmsg := "\nreadHeader did not correct invalid format to default"
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}

	file, filename = openTestFile("c2.mid")
	defer file.Close()
	header, err = readHeader(file)
	if err != nil {
		errmsg := "\nreadHeader returned error for recoverable file, invalid division."
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}
	if header.division != 24 {
		errmsg := "\nreadHeader did not correct invalid division to default"
		errmsg += "\nfilename was %s\n"
		t.Fatalf(errmsg, filename)
	}
}
