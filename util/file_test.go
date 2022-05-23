package util

import (
	"testing"
)

func TestCase_file(t *testing.T) {
	if err := DeleteFile("./test.log"); err == nil {
		t.Fatal(err)
	}
}
