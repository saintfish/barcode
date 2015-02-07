package barcode

import (
	"bitbucket.org/zombiezen/gopdf/pdf"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"strconv"
)

type EAN13 struct {
	code13 uint64
}

func (ean EAN13) Code13() uint64 {
	return ean.code13
}

func (ean EAN13) Code12() uint64 {
	return ean.code13 / 10
}

func (ean EAN13) Checksum() uint8 {
	return uint8(ean.code13 % 10)
}

func (ean EAN13) String() string {
	return fmt.Sprintf("%013d", ean.code13)
}

func (ean EAN13) RenderImage(img draw.Image, bound image.Rectangle, padding int) error {
	r, err := newBitmapRenderer(img, bound, padding)
	if err != nil {
		return err
	}
	renderCode13(ean.code13, r)
	return nil
}

func (ean EAN13) RenderPdf(canvas *pdf.Canvas, bound pdf.Rectangle, padding pdf.Unit) error {
	r, err := newPdfRenderer(canvas, bound, padding)
	if err != nil {
		return err
	}
	renderCode13(ean.code13, r)
	return nil
}

var checksumWeights = [2]uint64{3, 1}

func computeEANChecksum(code uint64) uint64 {
	var checksum uint64 = 0
	var i = 0
	for c := code; c > 0; c /= 10 {
		checksum += (c % 10) * checksumWeights[i%2]
		checksum %= 10
		i++
	}
	return (10 - checksum) % 10
}

var (
	errInvalidEAN         = errors.New("Invalid EAN code")
	errInvalidEANChecksum = errors.New("Invalid EAN checksum")
)

func EAN13FromCode12(code12 uint64) (EAN13, error) {
	if code12 > 999999999999 {
		return EAN13{}, errInvalidEAN
	}
	checksum := computeEANChecksum(code12)
	return EAN13{
		code13: code12*10 + checksum,
	}, nil
}

func EAN13FromCode13(code13 uint64) (EAN13, error) {
	if code13 > 9999999999999 {
		return EAN13{}, errInvalidEAN
	}
	checksum := computeEANChecksum(code13 / 10)
	if code13%10 != checksum {
		return EAN13{}, errInvalidEANChecksum
	}
	return EAN13{
		code13: code13,
	}, nil
}

func allDigits(code string) bool {
	for _, c := range code {
		if !(c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func EAN13FromString12(code string) (EAN13, error) {
	if len(code) != 12 || !allDigits(code) {
		return EAN13{}, errInvalidEAN
	}
	code12, err := strconv.ParseUint(code, 10, 64)
	if err != nil {
		return EAN13{}, err
	}
	return EAN13FromCode12(code12)
}

func EAN13FromString13(code string) (EAN13, error) {
	if len(code) != 13 || !allDigits(code) {
		return EAN13{}, errInvalidEAN
	}
	code13, err := strconv.ParseUint(code, 10, 64)
	if err != nil {
		return EAN13{}, err
	}
	return EAN13FromCode13(code13)
}

func EAN13FromString(code string) (EAN13, error) {
	switch len(code) {
	case 12:
		return EAN13FromString12(code)
	case 13:
		return EAN13FromString13(code)
	default:
		return EAN13{}, errInvalidEAN
	}
}
