package ean

import (
	"testing"
)

func TestEan13(t *testing.T) {
	codes := []string{
		"012345678901",
		"0123456789012",
	}
	for _, code := range codes {
		ean, err := NewEan13(code)
		if err != nil {
			t.Errorf("Error in create ean13 %s: %s", code, err)
			continue
		}
		ean.Encode()
	}
}
