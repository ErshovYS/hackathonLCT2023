package string_generator

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"strings"
)

// pdf requestpdf struct
type RequestPdf struct {
	body string
}

// new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

// parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func (r *RequestPdf) GeneratePDF(pdfPath string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pageReader := wkhtmltopdf.NewPageReader(strings.NewReader(r.body))
	pageReader.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(pageReader)

	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginTop.Set(0)

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeB6)
	//pdfg.PageHeightUnit.Set("2269px")
	//pdfg.PageWidthUnit.Set("1600px")
	//125 x 176
	pdfg.PageHeight.Set(176)
	pdfg.PageWidth.Set(125)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
