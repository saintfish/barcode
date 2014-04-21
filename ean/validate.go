package ean

import (
	"fmt"
)

var checksumWeights = []int{1, 3}

func allNumbersWithLength(code string, length int) error {
	if len(code) != length {
		return fmt.Errorf("Code \"%s\" has wrong length", code)
	}
	for i, c := range code {
		if !(c >= '0' && c <= '9') {
			return fmt.Errorf("Invalid char %c at %d in code \"%s\"", c, i, code)
		}
	}
	return nil
}

func computeChecksum(code string) byte {
	l := len(code)
	checksum := 0
	for i := 1; i <= l; i++ {
		checksum += c2i(code[l-i]) * checksumWeights[i%2]
		checksum %= 10
	}
	return i2c((10 - checksum) % 10)
}

func validateEan13(code string) error {
	if err := allNumbersWithLength(code, 13); err != nil {
		return err
	}
	checksum := computeChecksum(code[0:12])
	if checksum != code[12] {
		return fmt.Errorf(
			"Checksum mismatch for code \"%s\", expected %c vs actually %c",
			code, checksum, code[12])
	}
	return nil
}

func addChecksumEan13(code string) (string, error) {
	if err := allNumbersWithLength(code, 12); err != nil {
		return code, err
	}
	checksum := computeChecksum(code)
	return fmt.Sprintf("%s%c", code, checksum), nil
}
