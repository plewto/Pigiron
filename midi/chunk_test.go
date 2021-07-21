package midi

import (
	"testing"
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigpath"
)

func openTestFile(t *testing.T, name string) (*os.File, string) {
	filename := pigpath.ResourceFilename("testFiles", name)
	file, err := os.Open(filename)
	if err != nil {
		errmsg := "\nCan not open test file: '%s'"
		errmsg += "\n%s\n"
		t.Fatalf(errmsg, filename, err)
	}
	return file, filename
}
	


func TestReadChuckPreamble(t *testing.T) {
	file, filename := openTestFile(t, "a1.mid")
	defer file.Close()
	id, length, err := readChunkPreamble(file)
	if err != nil {
		errmsg := "\nreadChunkPreamble returned unexpected error for know good file."
		errmsg += "\nfilename was '%s'"
		errmsg += "\n%s\n"
		t.Fatalf(errmsg, filename, err)
	}
	if !chunkID(id).eq(headerID) {
		errmsg := "\nreadChunkPreamble expected id %s, got %s"
		t.Fatalf(fmt.Sprintf(errmsg, headerID), chunkID(id))
	}
	if length != 6 {
		errmsg := "\nreadChunkPreamble expected header length 6, got %d"
		t.Fatalf(fmt.Sprintf(errmsg, length))
	}
	
}


func TestReadRawChunk(t *testing.T) {
	file, filename := openTestFile(t, "a1.mid")
	defer file.Close()
	// read header
	id, data, err := readRawChunk(file)
	if err != nil {
		errmsg := "\nreadRawChunk returned unexpected error for know good file."
		errmsg += "\nfilename was '%s'"
		errmsg += "\n%s\n"
		t.Fatalf(errmsg, filename, err)
	}
	if !id.eq(headerID) {
		errmsg := "n\readRawChunk expected headerID, got %s"
		t.Fatalf(fmt.Sprintf(errmsg, id))
	}
	if len(data) != 6 {
		errmsg := "\nreadRawChunk expected header data length 6, got %d"
		t.Fatalf(fmt.Sprintf(errmsg, len(data)))
	}
	// read track
	id, data, err = readRawChunk(file)
	if err != nil {
		errmsg := "\nreadRawChunk could not read track chunk\n%s\n"
		t.Fatalf(errmsg, err)
	}
	if !id.eq(trackID) {
		errmsg := "\nreadRawChuck expected track id '%s', got '%s'\n"
		t.Fatalf(errmsg, trackID, id)
	}
	// file contents should be exhausted.
	buffer := make([]byte, 10)
	n, _ := file.Read(buffer)
	if n != 0 {
		errmsg := "\nreadRawChunk expected to have read all bytes, found additional values"
		t.Fatalf(errmsg)
	}
		

}

	
