package pigpath

import (
	"testing"
	"fmt"
	"path/filepath"
)


func TestPigpath(t *testing.T) {
	fmt.Print("")
	
	home := UserHomeDir()
	config := UserConfigDir()

	a := "~/foo"
	b := SubSpecialDirectories(a)
	if b[0:len(home)] != home {
		t.Fatalf("Expected filename start to be user's home: got '%s'", b)
	}

	a = "!/foo"
	b = SubSpecialDirectories(a)

	if b[0:len(config)] != config {
	 	t.Fatalf("Expected filename start to be configuration directory, got '%s'", b)
	 }
	
	b = Join("~", "alpha", "bar")
	expect := filepath.Join(home, "alpha", "bar") 
	if b != expect {
		t.Fatalf("Expected filenem.Join to return '%s', got '%s'", expect, b)
	}
	
	result := Join("~")
	if result != home {
		t.Fatalf("Expected Join(\"~\") to return '%s', got '%s'\n", home, result)
	}

	result = Join("!")
	if result != config {
		t.Fatalf("Expected Join(\"~\") to return '%s', got '%s'\n", config, result)
	}
}
