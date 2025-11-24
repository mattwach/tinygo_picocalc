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
var lcd *ili948x.Ili948x

func main() {
	lcd = ili948x.InitDisplay()
	tinyfont.WriteLine(lcd, font, 130, 100, "Press Any Key", white)
	var keyboard i2ckbd.I2CKbd
	if err := keyboard.Init(); err != nil {
		showError(err)
	}
	for {
		k, err := keyboard.GetChar()
		if err != nil {
			showError(err)
		}
		if k != 0 {
			lcd.FillRectangle(140, 140, 40, 40, ili948x.BLACK)
			tinyfont.DrawChar(lcd, font, 160, 160, rune(k), green)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func showError(err error) {
	tinyfont.WriteLine(lcd, font, 0, 64, err.Error(), red)
	for {
	}
}
