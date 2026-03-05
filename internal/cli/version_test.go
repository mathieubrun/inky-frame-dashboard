package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	out := b.String()
	if !strings.Contains(out, "1.0.0") {
		t.Errorf("Expected version 1.0.0, got %s", out)
	}
}
