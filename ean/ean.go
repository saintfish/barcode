package ean

import (
	"image"
)

func NewEan13(code string) (*Ean13, error) {
	if len(code) == 12 {
		newCode, err := addChecksumEan13(code)
		if err != nil {
			return nil, err
		}
		code = newCode
	}
	if err := validateEan13(code); err != nil {
		return nil, err
	}
	result := Ean13(code)
	return &result, nil
}

type Ean13 string

func (code Ean13) Encode() image.Image {
	return encodeEan13(string(code))
}

func (code Ean13) String() string {
	return string(code)
}
