package ean

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/gif"
)

const (
	kEan13WeigthPx = 13*7 + 5 + 3*2
	kEan13HeightPx = 70
)

var (
	lg = []string{
		"LLLLLL", "LLGLGG", "LLGGLG", "LLGGGL", "LGLLGG",
		"LGGLLG", "LGGGLL", "LGLGLG", "LGLGGL", "LGGLGL",
	}
)

var (
	l, g, r, blank, special image.Image
)

func imgFromFileOrDie(filename string) image.Image {
	data, err := Asset(filename)
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

func init() {
	l = imgFromFileOrDie("l.gif")
	g = imgFromFileOrDie("g.gif")
	r = imgFromFileOrDie("r.gif")
	blank = imgFromFileOrDie("blank.gif")
	special = imgFromFileOrDie("special.gif")
}

type imageWriter struct {
	image   draw.Image
	xOffset int
}

func newImageWriter() *imageWriter {
	i := image.NewRGBA(image.Rect(0, 0, kEan13WeigthPx, kEan13HeightPx))
	return &imageWriter{
		image:   i,
		xOffset: 0,
	}
}

type numBarType int

const (
	spaceBar numBarType = iota
	lBar
	gBar
	rBar
)

type specBarType int

const (
	boundBar specBarType = iota
	midBar
)

func (w *imageWriter) writeImage(img image.Image, x, width int) {
	draw.Draw(
		w.image, image.Rect(w.xOffset, 0, w.xOffset+width, kEan13HeightPx),
		img, image.Point{x, 0},
		draw.Src)
	w.xOffset += width
}

func (w *imageWriter) WriteNumber(t numBarType, num int) {
	var pattern image.Image
	switch t {
	case spaceBar:
		pattern = blank
	case lBar:
		pattern = l
	case gBar:
		pattern = g
	case rBar:
		pattern = r
	default:
		panic("Invalid type")
	}
	w.writeImage(pattern, num*7, 7)
}

func (w *imageWriter) WriteSpecial(t specBarType) {
	var x, width int
	switch t {
	case boundBar:
		x, width = 0, 3
	case midBar:
		x, width = 4, 5
	}
	w.writeImage(special, x, width)
}

func (w *imageWriter) GetImage() image.Image {
	return w.image
}

func encodeEan13(code string) image.Image {
	writer := newImageWriter()
	n := c2i(code[0])
	pattern := lg[n]
	writer.WriteNumber(spaceBar, n)
	writer.WriteSpecial(boundBar)
	for i, c := range code[1:7] {
		n := c2i(byte(c))
		if pattern[i] == 'L' {
			writer.WriteNumber(lBar, n)
		} else {
			writer.WriteNumber(rBar, n)
		}
	}
	writer.WriteSpecial(midBar)
	for _, c := range code[7:13] {
		n := c2i(byte(c))
		writer.WriteNumber(rBar, n)
	}
	writer.WriteSpecial(boundBar)
	return writer.GetImage()
}
