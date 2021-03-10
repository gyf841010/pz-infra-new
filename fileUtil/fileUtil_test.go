package fileUtil

import (
	"io/ioutil"
	"testing"
)

func TestPdf(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/测试pdf文件上传.pdf")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "pdf" {
		t.Fail()
	}
}

func TestDocx(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/sample.docx")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "docx" {
		t.Fail()
	}
}

func TestGif(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/sample.gif")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "gif" {
		t.Fail()
	}
}

func TestJPG(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/sample.jpg")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "jpg" {
		t.Fail()
	}
}

func TestPng(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/sample.png")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "png" {
		t.Fail()
	}
}

func TestTif(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/sample.tif")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "tif" {
		t.Fail()
	}
}

func TestXLSX(t *testing.T) {
	cont, err := ioutil.ReadFile("./test/sample.xlsx")
	if err != nil {
		t.Fail()
	}

	ftype := GetFileTypeNew(cont)
	if ftype != "xlsx" {
		t.Fail()
	}
}
