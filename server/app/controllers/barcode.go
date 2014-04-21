package controllers

import (
	"github.com/robfig/revel"
	"github.com/saintfish/barcode"
	"image"
	"image/png"
	"net/http"
)

type pngResult struct {
	image image.Image
}

func (p *pngResult) Apply(req *revel.Request, resp *revel.Response) {
	resp.ContentType = "image/png"
	err := png.Encode(resp.Out, p.image)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError, "text/plain; charset=utf-8")
		return
	}
}

type Ean13 struct {
	*revel.Controller
}

func (c Ean13) Encode(code string) revel.Result {
	ean13, err := barcode.NewEan13(code)
	if err != nil {
		revel.WARN.Print(err)
		return c.NotFound("Invalid ean13 barcode")
	}
	if ean13.String() != code {
		return c.Redirect(Ean13.Encode, ean13.String())
	}
	img := ean13.Encode()
	return &pngResult{img}
}
