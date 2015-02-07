package barcode

import (
	"fmt"
	"testing"
)

func TestEAN13(t *testing.T) {
	type data struct {
		code12   uint64
		checksum uint8
	}
	var testdata = []data{
		{590123412345, 7},
		{111234567890, 1},
		{354212567890, 2},
		{12345678901, 2},
	}
	for _, d := range testdata {
		ean, err := EAN13FromCode12(d.code12)
		if err != nil {
			t.Errorf("Invalid ean %d", d.code12)
		}
		if d.checksum != ean.Checksum() {
			t.Errorf("Unexpected checksum %d v.s. %d", d.checksum, ean.Checksum)
		}
		ean, err = EAN13FromCode13(d.code12*10 + uint64(d.checksum))
		if err != nil {
			t.Errorf("Invalid ean %d", d.code12)
		}
		if d.checksum != ean.Checksum() {
			t.Errorf("Unexpected checksum %d v.s. %d", d.checksum, ean.Checksum)
		}
		ean, err = EAN13FromString12(fmt.Sprintf("%012d", d.code12))
		if err != nil {
			t.Errorf("Invalid ean %d", d.code12)
		}
		if d.checksum != ean.Checksum() {
			t.Errorf("Unexpected checksum %d v.s. %d", d.checksum, ean.Checksum)
		}
		ean, err = EAN13FromString13(fmt.Sprintf("%012d%d", d.code12, d.checksum))
		if err != nil {
			t.Errorf("Invalid ean %d", d.code12)
		}
		if d.checksum != ean.Checksum() {
			t.Errorf("Unexpected checksum %d v.s. %d", d.checksum, ean.Checksum)
		}
		ean, err = EAN13FromString(fmt.Sprintf("%012d", d.code12))
		if err != nil {
			t.Errorf("Invalid ean %d", d.code12)
		}
		if d.checksum != ean.Checksum() {
			t.Errorf("Unexpected checksum %d v.s. %d", d.checksum, ean.Checksum)
		}
		ean, err = EAN13FromString(fmt.Sprintf("%012d%d", d.code12, d.checksum))
		if err != nil {
			t.Errorf("Invalid ean %d", d.code12)
		}
		if d.checksum != ean.Checksum() {
			t.Errorf("Unexpected checksum %d v.s. %d", d.checksum, ean.Checksum)
		}
	}
}

func TestInvalidEAN13(t *testing.T) {
	var invalidCode12 = []uint64{
		1000000000000,
		9999999999999,
	}
	for _, c := range invalidCode12 {
		if _, err := EAN13FromCode12(c); err == nil {
			t.Errorf("Unexpected valid code %d", c)
		}
	}
	var invalidCode13 = []uint64{
		10000000000000, // too long
		5901234123450,  // invalid checksum
	}
	for _, c := range invalidCode13 {
		if _, err := EAN13FromCode13(c); err == nil {
			t.Errorf("Unexpected valid code %d", c)
		}
	}
	var invalidString = []string{
		"10000000000000", // too long
		"123456789012a",  // not all digits
		"5901234123450",  // invalid checksum
	}
	for _, c := range invalidString {
		if _, err := EAN13FromString(c); err == nil {
			t.Errorf("Unexpected valid code %s", c)
		}
	}
}
