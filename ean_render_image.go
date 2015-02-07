package barcode

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
)

//go:generate sh gen-data.sh

func imgFromFileOrDie(filename string) image.Image {
	data, err := _asset(filename)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	return img
}

var digitsImage image.Image
var digitsImageWidth, digitsImageHeight int

func init() {
	digitsImage = imgFromFileOrDie("digits.gif")
	digitsImageWidth = digitsImage.Bounds().Dx() / 10
	digitsImageHeight = digitsImage.Bounds().Dy()
}

func measureBitmapFont(width int) (fontSize, fontWidth, fontHeight int) {
	fontSize = width / digitsImageWidth
	fontWidth = digitsImageWidth * fontSize
	fontHeight = digitsImageHeight * fontSize
	return
}

type bitmapRenderer struct {
	img       draw.Image
	bound     image.Rectangle
	converter *eanCoordinateConverter
}

func newBitmapRenderer(img draw.Image, bound image.Rectangle, padding int) (eanRenderer, error) {
	bound = bound.Intersect(img.Bounds())
	inner := image.Rectangle{
		Min: bound.Min.Add(image.Pt(padding, padding)),
		Max: bound.Max.Sub(image.Pt(padding, padding)),
	}.Canon()
	converter, err := newEanCoordinateConverter(inner, measureBitmapFont)
	if err != nil {
		return nil, err
	}
	return &bitmapRenderer{
		img:       img,
		bound:     bound,
		converter: converter,
	}, nil
}

func fillRect(img draw.Image, rect image.Rectangle, color color.Color) {
	draw.Draw(img, rect, image.NewUniform(color), image.ZP, draw.Src)
}

func (r *bitmapRenderer) Start() *eanCoordinateConverter {
	fillRect(r.img, r.bound, color.White)
	return r.converter
}

func (r *bitmapRenderer) DrawBar(rect image.Rectangle) {
	fillRect(r.img, rect, color.Black)
}

func (r *bitmapRenderer) DrawDigit(digit int, rect image.Rectangle, fontSize int) {
	for i := 0; i < digitsImageWidth; i++ {
		for j := 0; j < digitsImageHeight; j++ {
			dotTopLeft := rect.Min.Add(image.Pt(i*fontSize, j*fontSize))
			dotBottomRight := dotTopLeft.Add(image.Pt(fontSize, fontSize))
			fillRect(
				r.img,
				image.Rectangle{Min: dotTopLeft, Max: dotBottomRight},
				digitsImage.At(digitsImageWidth*digit+i, j))
		}
	}
}

func (r *bitmapRenderer) End() {
}
