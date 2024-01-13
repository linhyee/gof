package help

import (
	"testing"
)

func TestInarrayInt64(t *testing.T) {
	arr := []int64{1, 2, 3, 4, 6, 8, 9}
	if InArrayInt64(arr, 5) == true {
		t.Fatal("5 shouldn't be found")
	}
}

func TestReadFile(t *testing.T) {
	s := ReadFile("./array.go")
	if s == "" {
		t.Fatal("read empty content")
	}
}

func TestAtoi(t *testing.T) {
	s := "8"
	if Atoi(s) != 8 {
		t.Fatal("should be 8")
	}
}
