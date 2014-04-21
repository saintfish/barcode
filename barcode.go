package barcode

import (
	"github.com/saintfish/barcode/ean"
	"image"
)

type Barcode interface {
	Encode() image.Image
	String() string
}

func NewEan13(code string) (Barcode, error) {
	return ean.NewEan13(code)
}
