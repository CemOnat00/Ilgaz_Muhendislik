// Package catalog ürün verisinden dinamik PDF kataloğu üretir.
package catalog

import (
	"fmt"
	"strings"
	"time"

	"github.com/cemonat00/ilgaz-backend/internal/models"
	"github.com/signintech/gopdf"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	pageW       = 595.28
	pageH       = 841.89
	marginX     = 42.0
	contentW    = pageW - 2*marginX
	footerY     = pageH - 40.0
	bottomLimit = pageH - 64.0
)

// Generate, verilen ürün listesinden katalog PDF'i üretir ve byte dizisi döndürür.
// products önceden filtrelenmiş gelir; categoryLabel başlıkta gösterilir.
func Generate(products []models.Product, categoryLabel string) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	if err := pdf.AddTTFFontData("reg", goregular.TTF); err != nil {
		return nil, err
	}
	if err := pdf.AddTTFFontData("bold", gobold.TTF); err != nil {
		return nil, err
	}

	if categoryLabel == "" || categoryLabel == "Genel" {
		categoryLabel = "Tüm Ürünler"
	}

	pageNum := 0
	y := 0.0

	drawFooter := func() {
		pdf.SetStrokeColor(229, 231, 235)
		pdf.SetLineWidth(0.5)
		pdf.Line(marginX, footerY-12, pageW-marginX, footerY-12)
		pdf.SetFont("reg", "", 8)
		pdf.SetTextColor(150, 150, 150)
		pdf.SetXY(marginX, footerY)
		pdf.Cell(nil, "Ilgaz Mühendislik  •  www.ilgazmuhendislik.com")
		pdf.SetXY(pageW-marginX-55, footerY)
		pdf.Cell(nil, fmt.Sprintf("Sayfa %d", pageNum))
	}

	newPage := func(first bool) {
		if pageNum > 0 {
			drawFooter()
		}
		pdf.AddPage()
		pageNum++

		if first {
			pdf.SetFillColor(214, 48, 49)
			pdf.RectFromUpperLeftWithStyle(0, 0, pageW, 88, "F")
			pdf.SetFont("bold", "", 22)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetXY(marginX, 28)
			pdf.Cell(nil, "ILGAZ MÜHENDİSLİK")
			pdf.SetFont("reg", "", 10)
			pdf.SetTextColor(255, 220, 220)
			pdf.SetXY(marginX, 56)
			pdf.Cell(nil, "Ürün Kataloğu")

			y = 116
			pdf.SetFont("bold", "", 17)
			pdf.SetTextColor(31, 41, 55)
			pdf.SetXY(marginX, y)
			pdf.Cell(nil, categoryLabel)
			y += 24
			pdf.SetFont("reg", "", 9)
			pdf.SetTextColor(120, 120, 120)
			pdf.SetXY(marginX, y)
			pdf.Cell(nil, fmt.Sprintf("%d ürün  •  %s tarihli", len(products), time.Now().Format("02.01.2006")))
			y += 14
			pdf.SetStrokeColor(214, 48, 49)
			pdf.SetLineWidth(2)
			pdf.Line(marginX, y, marginX+46, y)
			y += 22
		} else {
			pdf.SetFont("reg", "", 8)
			pdf.SetTextColor(170, 170, 170)
			pdf.SetXY(marginX, 30)
			pdf.Cell(nil, "ILGAZ MÜHENDİSLİK — Ürün Kataloğu / "+categoryLabel)
			y = 56
		}
	}

	ensure := func(needed float64) {
		if y+needed > bottomLimit {
			newPage(false)
		}
	}

	writeLines := func(text string, family string, size, lineH, indent float64, r, g, b uint8) {
		pdf.SetFont(family, "", size)
		lines, err := pdf.SplitTextWithWordWrap(text, contentW-indent)
		if err != nil || len(lines) == 0 {
			lines = []string{text}
		}
		for _, ln := range lines {
			ensure(lineH)
			pdf.SetFont(family, "", size)
			pdf.SetTextColor(r, g, b)
			pdf.SetXY(marginX+indent, y)
			pdf.Cell(nil, ln)
			y += lineH
		}
	}

	newPage(true)

	if len(products) == 0 {
		pdf.SetFont("reg", "", 11)
		pdf.SetTextColor(120, 120, 120)
		pdf.SetXY(marginX, y+10)
		pdf.Cell(nil, "Bu kategoride listelenecek ürün bulunmuyor.")
		drawFooter()
		return pdf.GetBytesPdf(), nil
	}

	for i, p := range products {
		ensure(64)

		name := strings.TrimSpace(p.Isim)
		if name == "" {
			name = strings.TrimSpace(p.Baslik)
		}
		writeLines(fmt.Sprintf("%d.  %s", i+1, name), "bold", 12, 16, 0, 31, 41, 55)

		stock := "Stokta"
		if !p.StockStatus {
			stock = "Tükendi"
		}
		meta := p.Kategori
		if p.Price > 0 {
			meta += "   |   " + formatPrice(p.Price)
		}
		meta += "   |   " + stock
		writeLines(meta, "reg", 9, 15, 0, 214, 48, 49)
		y += 2

		if desc := strings.TrimSpace(p.Description); desc != "" {
			writeLines(desc, "reg", 9.5, 13, 0, 80, 80, 80)
		}

		for _, s := range p.TechnicalSpecs {
			key := strings.TrimSpace(s.Key)
			val := strings.TrimSpace(s.Value)
			if key == "" && val == "" {
				continue
			}
			writeLines("•  "+key+": "+val, "reg", 8.5, 12, 10, 110, 110, 110)
		}

		y += 12
		ensure(6)
		pdf.SetStrokeColor(236, 236, 236)
		pdf.SetLineWidth(0.5)
		pdf.Line(marginX, y, pageW-marginX, y)
		y += 18
	}

	drawFooter()
	return pdf.GetBytesPdf(), nil
}

// formatPrice fiyatı binlik ayraçlı Türk Lirası biçimine çevirir (ör. 12.500 TL).
func formatPrice(v float64) string {
	n := int64(v + 0.5)
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, c)
	}
	return string(out) + " TL"
}
