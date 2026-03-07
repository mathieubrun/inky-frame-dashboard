package core

import (
	"testing"
)

func TestCalculateMD5(t *testing.T) {
	data := []byte("hello world")
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3"
	actual := CalculateMD5(data)
	if actual != expected {
		t.Errorf("CalculateMD5() = %s; want %s", actual, expected)
	}
}
