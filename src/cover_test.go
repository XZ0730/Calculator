package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func TestCover(t *testing.T) {
	a := app.New()
	_ = a.NewWindow("Calculator")
	entry = widget.NewEntry()
	entry.MultiLine = true
	entry.Resize(fyne.NewSize(150, 150))
	funcs := equals()
	entry.Text = "1+2+3+2^3+rsin(sin(30))+rtan(tan(30))+rcos(cos(60))"
	t.Log(entry.Text)
	funcs()
	entry.Text = "error:xxx"
	funcs()
	entry.Text = "Inf"
	funcs()
	entry.Text = "rsin(10)"
	funcs()
	entry.Text = "rcos(10)"
	funcs()
	ifunc := input("+")
	ifunc()
	ifunc()
	entry.Text = "3"
	sfunc := sign()
	sfunc()
	entry.Text = "3.3"
	sfunc()
	entry.Text = "8^9999"
	bfunc := back()
	bfunc()
	funcs()
	entry.Text = "\n"
	bfunc()
	entry.Text = ""
	bfunc()
	t.Log(entry.Text)
}
