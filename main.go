package main

import (
	"image/color"
	"picocalc/i2ckbd"
	"picocalc/ili948x"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/shnm"
)

var (
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
)

var font = &shnm.Shnmk12
var disp *ili948x.Ili948x

func main() {
	disp = ili948x.InitDisplay()
	tinyfont.WriteLine(disp, font, 130, 100, "Press Any Key", white)
	var kdb i2ckbd.I2CKbd
	if err := kdb.Init(); err != nil {
		showError(err)
	}
	for {
		k, err := kdb.GetChar()
		if err != nil {
			showError(err)
		}
		if k != 0 {
			disp.FillRectangle(140, 140, 40, 40, ili948x.BLACK)
			tinyfont.DrawChar(disp, font, 160, 160, rune(k), green)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func showError(err error) {
	tinyfont.WriteLine(disp, font, 0, 64, err.Error(), red)
	for {
	}
}
