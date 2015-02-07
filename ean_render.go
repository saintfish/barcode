package barcode

import (
	"errors"
	"image"
)

// b2s converts string to stripes
func b2s(b string) []int {
	lastZero := true
	r := []int{}
	for i, c := range b {
		if lastZero {
			if c == '1' {
				r = append(r, i)
				lastZero = false
			}
		} else {
			if c == '0' {
				r = append(r, i)
				lastZero = true
			}
		}
	}
	if !lastZero {
		r = append(r, len(b))
	}
	return r
}

var barTable = [10][3][]int{
	{b2s("0001101"), b2s("0100111"), b2s("1110010")},
	{b2s("0011001"), b2s("0110011"), b2s("1100110")},
	{b2s("0010011"), b2s("0011011"), b2s("1101100")},
	{b2s("0111101"), b2s("0100001"), b2s("1000010")},
	{b2s("0100011"), b2s("0011101"), b2s("1011100")},
	{b2s("0110001"), b2s("0111001"), b2s("1001110")},
	{b2s("0101111"), b2s("0000101"), b2s("1010000")},
	{b2s("0111011"), b2s("0010001"), b2s("1000100")},
	{b2s("0110111"), b2s("0001001"), b2s("1001000")},
	{b2s("0001011"), b2s("0010111"), b2s("1110100")},
}

var dispatchTable = [10][6]int{
	{0, 0, 0, 0, 0, 0},
	{0, 0, 1, 0, 1, 1},
	{0, 0, 1, 1, 0, 1},
	{0, 0, 1, 1, 1, 0},
	{0, 1, 0, 0, 1, 1},
	{0, 1, 1, 0, 0, 1},
	{0, 1, 1, 1, 0, 0},
	{0, 1, 0, 1, 0, 1},
	{0, 1, 0, 1, 1, 0},
	{0, 1, 1, 0, 1, 0},
}

var (
	startMarker  = b2s("101")
	endMarker    = startMarker
	centerMarker = b2s("01010")
)

const (
	digitBarSize     int = 7
	startMarkerSize      = 3
	endMarkerSize        = 3
	centerMarkerSize     = 5
)

// eanCoordinateConverter maps a logical coordinate to the coordinate in target system.
type eanCoordinateConverter struct {
	bound    image.Rectangle
	fontDim  image.Rectangle
	scale    int
	fontSize int
}

var errAreaTooSmall = errors.New("Bound area too small")
var errFontTooBig = errors.New("Font is too big to fit in the small barcode")

type fontMeasurer func(width int) (fontSize, fontWidth, fontHeight int)

func newEanCoordinateConverter(outerBound image.Rectangle, fm fontMeasurer) (*eanCoordinateConverter, error) {
	const rightMargin = digitBarSize
	const logicalWidth = 13*digitBarSize + startMarkerSize + endMarkerSize + centerMarkerSize + rightMargin
	scale := outerBound.Dx() / logicalWidth
	if scale <= 0 {
		return nil, errAreaTooSmall
	}
	xMargin := (outerBound.Dx() % logicalWidth) / 2
	fSize, fWidth, fHeight := fm(digitBarSize * scale)
	if digitBarSize*scale < fWidth {
		return nil, errFontTooBig
	}
	if outerBound.Dy() < fHeight*2 {
		return nil, errAreaTooSmall
	}
	return &eanCoordinateConverter{
		bound: image.Rectangle{
			Min: outerBound.Min.Add(image.Pt(xMargin, 0)),
			Max: outerBound.Max.Sub(image.Pt(xMargin, 0)),
		},
		fontDim:  image.Rect(0, 0, fWidth, fHeight),
		scale:    scale,
		fontSize: fSize,
	}, nil
}

// x0, x1 are logical pixel
func (c *eanCoordinateConverter) translateBar(x0, x1 int, long bool) image.Rectangle {
	h := c.bound.Dy() - c.fontDim.Dy()
	if long {
		h += c.fontDim.Dy() / 2
	}
	return image.Rectangle{
		Min: c.bound.Min.Add(image.Pt(x0*c.scale, 0)),
		Max: c.bound.Min.Add(image.Pt(x1*c.scale, h)),
	}
}

func (c *eanCoordinateConverter) translateFont(x int) (rect image.Rectangle, fontSize int) {
	fontCellWidth := c.scale * digitBarSize
	fontXOffset := (fontCellWidth - c.fontDim.Dx()) / 2
	fontYOffset := c.bound.Dy() - c.fontDim.Dy()
	topLeft := c.bound.Min.Add(image.Pt(x*c.scale+fontXOffset, fontYOffset))
	bottomRight := topLeft.Add(image.Pt(c.fontDim.Dx(), c.fontDim.Dy()))
	return image.Rectangle{Min: topLeft, Max: bottomRight}, c.fontSize
}

type eanRenderer interface {
	// Start to render a new barcode. Return the coordinate converter for the render logic to decide the coordination
	Start() *eanCoordinateConverter
	// Draw a black bar
	DrawBar(rect image.Rectangle)
	// Draw a digit at x, y as top left corner
	DrawDigit(digit int, rect image.Rectangle, fontSize int)
	// End of rendering the barcode. Clean up can be done in this function
	End()
}

func drawStripe(cx int, s []int, long bool, r eanRenderer, c *eanCoordinateConverter) {
	for i := 0; i+1 < len(s); i += 2 {
		r.DrawBar(c.translateBar(cx+s[i], cx+s[i+1], long))
	}
}

func drawDigit(cx int, digit int, r eanRenderer, c *eanCoordinateConverter) {
	rect, fontSize := c.translateFont(cx)
	r.DrawDigit(digit, rect, fontSize)
}

func renderCode13(code13 uint64, r eanRenderer) {
	c := r.Start()

	var digits [13]int
	for i, c := 12, code13; i >= 0; i-- {
		digits[i] = int(c % 10)
		c = c / 10
	}

	cx := 0
	first := digits[0]
	// Draw first digit
	drawDigit(cx, first, r, c)
	cx += digitBarSize
	// Draw start marker
	drawStripe(cx, startMarker, true, r, c)
	cx += startMarkerSize
	// Draw fist 6 digits
	for i := 1; i <= 6; i++ {
		stripe := barTable[digits[i]][dispatchTable[first][i-1]]
		drawStripe(cx, stripe, false, r, c)
		drawDigit(cx, digits[i], r, c)
		cx += digitBarSize
	}
	// Draw center marker
	drawStripe(cx, centerMarker, true, r, c)
	cx += centerMarkerSize
	// Draw last 6 digits
	for i := 7; i <= 12; i++ {
		drawStripe(cx, barTable[digits[i]][2], false, r, c)
		drawDigit(cx, digits[i], r, c)
		cx += digitBarSize
	}
	// Draw end marker
	drawStripe(cx, endMarker, true, r, c)
	cx += endMarkerSize
	r.End()
}
