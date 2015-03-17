package barcode

import (
	"bitbucket.org/saintfish/gopdf/pdf"
	"image"
	"image/gif"
	"os"
	"testing"
)

func TestRenderImage(t *testing.T) {
	r := image.Rect(0, 0, 600, 300)
	img := image.NewGray(r)
	code, _ := EAN13FromString12("590123412345")
	code.RenderImage(img, r, 20)

	f, _ := os.Create("test.gif")
	defer f.Close()
	gif.Encode(f, img, nil)
}

func TestRenderPdf(t *testing.T) {
	doc := pdf.New()
	p := doc.NewPage(pdf.USLetterWidth, pdf.USLetterHeight)
	p.Translate(0.5*pdf.Inch, 0.5*pdf.Inch)
	num := uint64(590123412345)
	for i := 0; i < 5; i++ {
		for j := 0; j < 10; j++ {
			code, _ := EAN13FromCode12(num)
			p1 := pdf.Point{pdf.Unit(i) * 1.5 * pdf.Inch, pdf.Unit(j) * 1 * pdf.Inch}
			p2 := pdf.Point{p1.X + 1.5*pdf.Inch, p1.Y + 1*pdf.Inch}
			rect := pdf.Rectangle{p1, p2}
			code.RenderPdf(p, rect, 0.13*pdf.Inch)
			num += 20
		}
	}
	p.Close()

	f, _ := os.Create("test.pdf")
	defer f.Close()
	doc.Encode(f)
}
