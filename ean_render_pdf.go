package barcode

import (
	"bitbucket.org/zombiezen/gopdf/pdf"
	"image"
)

type pdfRenderer struct {
	canvas    *pdf.Canvas
	bound     pdf.Rectangle
	converter *eanCoordinateConverter
}

// font metrics of Helvetica
const (
	pdfFontAscender  = 718 + 100
	pdfFontDescender = -207
	pdfFontHeight    = pdfFontAscender - pdfFontDescender
	pdfDigitWidth    = 556
	pdfFontUnitScale = 1000
)

const (
	pdfCoordinateScale = 5
)

func measurePdfFont(width int) (fontSize, fontWidth, fontHeight int) {
	fontSize = width * pdfFontUnitScale / pdfDigitWidth / pdfCoordinateScale
	fontWidth = pdfDigitWidth * fontSize * pdfCoordinateScale / pdfFontUnitScale
	fontHeight = pdfFontHeight * fontSize * pdfCoordinateScale / pdfFontUnitScale
	return
}

func newPdfRenderer(canvas *pdf.Canvas, bound pdf.Rectangle, padding pdf.Unit) (*pdfRenderer, error) {
	imgRect := image.Rect(
		int(padding*pdfCoordinateScale),
		int(padding*pdfCoordinateScale),
		int((bound.Dx()-padding)*pdfCoordinateScale),
		int((bound.Dy()-padding)*pdfCoordinateScale))
	converter, err := newEanCoordinateConverter(imgRect, measurePdfFont)
	if err != nil {
		return nil, err
	}
	return &pdfRenderer{
		canvas:    canvas,
		bound:     bound,
		converter: converter,
	}, nil
}

func (r *pdfRenderer) Start() *eanCoordinateConverter {
	r.canvas.Push()
	r.canvas.SetColor(1, 1, 1) // white
	p := new(pdf.Path)
	p.Rectangle(r.bound)
	r.canvas.Fill(p)
	r.canvas.SetColor(0, 0, 0) // black
	r.canvas.Transform(1, 0, 0, -1, float32(r.bound.Min.X), float32(r.bound.Max.Y))
	return r.converter
}

func (r *pdfRenderer) DrawBar(rect image.Rectangle) {
	p := new(pdf.Path)
	p.Rectangle(pdf.Rectangle{
		Min: pdf.Point{pdf.Unit(rect.Min.X) / pdfCoordinateScale, pdf.Unit(rect.Min.Y) / pdfCoordinateScale},
		Max: pdf.Point{pdf.Unit(rect.Max.X) / pdfCoordinateScale, pdf.Unit(rect.Max.Y) / pdfCoordinateScale},
	})
	r.canvas.Fill(p)
}

var digitsString = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func (r *pdfRenderer) DrawDigit(digit int, rect image.Rectangle, fontSize int) {
	text := new(pdf.Text)
	text.SetFont(pdf.Helvetica, pdf.Unit(fontSize))
	text.Text(digitsString[digit])
	r.canvas.Push()
	r.canvas.Transform(
		1, 0, 0, -1,
		float32(rect.Min.X)/pdfCoordinateScale,
		float32(rect.Min.Y)/pdfCoordinateScale+float32(pdfFontAscender*fontSize)/pdfFontUnitScale)
	r.canvas.DrawText(text)
	r.canvas.Pop()
}

func (r *pdfRenderer) End() {
	r.canvas.Pop()
}
