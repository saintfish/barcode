package ean

import (
	"testing"
)

func TestValidateEan13(t *testing.T) {
	valid := []string{
		"0123456789012",
	}
	for _, code := range valid {
		if err := validateEan13(code); err != nil {
			t.Errorf("Code %s: expected valid, actually error %s", code, err)
		}
	}
}
