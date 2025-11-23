// The starting point for this code was https://github.com/inindev/ili948x
// but it's been heavily modified to the point where it shares little in common
// with the original.  Some changes:
//
// 1) Use RGB565 instead of RGB666 for faster screen updates
// 2) Introduce compatbility with the tinygo displayer interface
// 3) Add RGB565 and Bitmap for efficient (basically direct) pixel operations.
package ili948x

import (
	"image/color"

	"machine"
	"time"

	"tinygo.org/x/drivers/pixel"
)

type Rotation uint8
type RGB565 uint16

// By experimantal observation, I discovered the pixel format
// from uint16 (the way it's unpacked in spi_transport) to be:
//
// BBBBBGGGGGGRRRRR
// 4321021054343210
//
// rpn:
// '11 << 1> 5 << | 1> | hex' rgb565=
const (
	BLACK RGB565 = 0x0000
	WHITE RGB565 = 0xFFFF

	RED     RGB565 = 0x001f
	ORANGE  RGB565 = 0x03ff
	YELLOW  RGB565 = 0x07ff
	GREEN   RGB565 = 0x07E0
	BLUE    RGB565 = 0xf800
	CYAN    RGB565 = 0xffe0
	MAGENTA RGB565 = 0xF81F
)

func RGBA2RGB565(c color.RGBA) RGB565 {
	return (RGB565(c.R) >> 3) |
		((RGB565(c.G) >> 2) << 5) |
		((RGB565(c.B) >> 3) << 11)
}

// r -> 0-31, g -> 0-63, b -> 0-31
func NewRGB565(r, g, b uint8) RGB565 {
	return RGB565(uint16(r) | (uint16(g) << 5) | (uint16(b) << 11))
}

const (
	TFT_WIDTH  = 320
	TFT_HEIGHT = 320
)

// Image buffer type used in the ili9341.
type Image = pixel.Image[pixel.RGB565BE]

type Ili948x struct {
	trans  *spiTransport
	cs     machine.Pin // spi chip select
	dc     machine.Pin // tft data / command
	rst    machine.Pin // tft reset
	width  int16       // tft pixel width
	height int16       // tft pixel height
	x0, x1 int16       // current address window for
	y0, y1 int16       //  CMD_PASET and CMD_CASET
}

func InitDisplay() *Ili948x {
	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP12,
		Frequency: 40000000,
	})

	display := NewIli9488(
		NewSPITransport(*machine.SPI1),
		machine.GP13, // chip select
		machine.GP14, // data / command
		machine.GP15, // reset
	)

	display.FillScreen(0)

	return display
}

func NewIli9488(trans *spiTransport, cs, dc, rst machine.Pin) *Ili948x {
	disp := &Ili948x{
		trans: trans,
		cs:    cs,
		dc:    dc,
		rst:   rst,
		x0:    0,
		x1:    0,
		y0:    0,
		y1:    0,
	}

	// chip select pin
	if cs != machine.NoPin { // cs may be implemented by hardware spi
		cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
		cs.High()
	}

	// data/command pin
	dc.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dc.High()

	// reset pin
	if rst != machine.NoPin {
		disp.rst.Configure(machine.PinConfig{Mode: machine.PinOutput})
		disp.rst.High()
	}

	// reset the display
	disp.Reset()

	// init display settings
	disp.init()

	return disp
}

// For compatibility with https://github.com/tinygo-org/drivers/blob/release/displayer.go
// Size returns the current size of the display.
func (disp *Ili948x) Size() (int16, int16) {
	return disp.width, disp.height
}

// For compatibility with https://github.com/tinygo-org/drivers/blob/release/displayer.go
// Use SetPixel565 for better performance
func (disp *Ili948x) SetPixel(x, y int16, c color.RGBA) {
	disp.setWindow(x, y, 1, 1)
	disp.writeCmd(CMD_RAMWR)
	disp.startWrite()
	disp.trans.write16(uint16(RGBA2RGB565(c)))
	disp.endWrite()
}

func (disp *Ili948x) SetPixel565(x, y int16, c RGB565) {
	disp.setWindow(x, y, 1, 1)
	disp.writeCmd(CMD_RAMWR)
	disp.startWrite()
	disp.trans.write16(uint16(c))
	disp.endWrite()
}

// DrawHLine draws a horizontal line with the specified color.
func (disp *Ili948x) DrawHLine(x0, x1, y int16, c RGB565) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if x0 < 0 {
		x0 = 0
	}
	if x1 >= TFT_WIDTH {
		x1 = TFT_WIDTH - 1
	}
	width := x1 - x0 + 1
	if width > 0 {
		disp.setWindow(x0, y, width, 1)
		disp.writeCmd(CMD_RAMWR)
		disp.startWrite()
		disp.trans.write16n(uint16(c), int(width))
		disp.endWrite()
	}
}

// DrawVLine draws a vertical line with the specified color.
func (disp *Ili948x) DrawVLine(x, y0, y1 int16, c RGB565) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if y0 < 0 {
		y0 = 0
	}
	if y1 >= TFT_HEIGHT {
		y1 = TFT_HEIGHT - 1
	}
	height := y1 - y0 + 1
	if height > 0 {
		disp.setWindow(x, y0, 1, height)
		disp.writeCmd(CMD_RAMWR)
		disp.startWrite()
		disp.trans.write16n(uint16(c), int(height))
		disp.endWrite()
	}
}

// FillScreen fills the screen with the specified color.
func (disp *Ili948x) FillScreen(c RGB565) {
	disp.FillRectangle(0, 0, disp.height, disp.width, c)
}

// FillRectangle fills a rectangle at given coordinates and dimensions with the specified color.
func (disp *Ili948x) FillRectangle(x, y, width, height int16, c RGB565) {
	if x < 0 {
		width += x
		x = 0
	}
	if (x + width) > TFT_WIDTH {
		width = TFT_WIDTH - x
	}
	if width <= 0 {
		return
	}
	if y < 0 {
		height += y
		y = 0
	}
	if (y + height) > TFT_HEIGHT {
		height = TFT_HEIGHT - y
	}
	if height <= 0 {
		return
	}
	disp.setWindow(x, y, width, height)
	disp.writeCmd(CMD_RAMWR)
	disp.startWrite()
	disp.trans.write16n(uint16(c), int(width)*int(height))
	disp.endWrite()
}

// DrawBitmap copies an RGB565 bitmap to the internal buffer at given coordinates
func (disp *Ili948x) DrawBitmap(x, y int16, bm *Bitmap) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if (x + bm.Width) >= TFT_WIDTH {
		return
	}
	if (y + bm.Height) >= TFT_HEIGHT {
		return
	}
	disp.setWindow(x, y, bm.Width, bm.Height)
	disp.writeCmd(CMD_RAMWR)
	disp.startWrite()
	disp.trans.writeRGB565(bm.Data)
	disp.endWrite()
}

// Reset performs a hardware reset if rst pin present, otherwise performs a CMD_SWRESET software reset of the TFT display.
func (disp *Ili948x) Reset() {
	// prefer a hardware reset if there is one
	if disp.rst != machine.NoPin {
		disp.rst.Low()
		time.Sleep(time.Millisecond * 64) // datasheet says 10ms
		disp.rst.High()
	} else {
		// if no hardware reset, send software reset
		disp.writeCmd(CMD_SWRESET)
	}
	time.Sleep(time.Millisecond * 140) // datasheet says 120ms
}

// setWindow defines the output area for subsequent calls to CMD_RAMWR
func (disp *Ili948x) setWindow(x, y, w, h int16) {
	x1 := x + w - 1
	if x != disp.x0 || x1 != disp.x1 {
		disp.writeCmd(CMD_CASET,
			uint8(x>>8),
			uint8(x),
			uint8(x1>>8),
			uint8(x1),
		)
		disp.x0, disp.x1 = x, x1
	}
	y1 := y + h - 1
	if y != disp.y0 || y1 != disp.y1 {
		disp.writeCmd(CMD_PASET,
			uint8(y>>8),
			uint8(y),
			uint8(y1>>8),
			uint8(y1),
		)
		disp.y0, disp.y1 = y, y1
	}
}

// init performs base-level initialization and setup of the TFT display
func (disp *Ili948x) init() {
	disp.writeCmd(CMD_PWCTRL1,
		0x17, // VREG1OUT:  5.0000
		0x15) // VREG2OUT: -4.8750

	disp.writeCmd(CMD_PWCTRL2,
		0x41) // VGH: VCI x 6  VGL: -VCI x 4

	disp.writeCmd(CMD_VMCTRL,
		0x00, // nVM
		0x12, // VCM_REG:    -1.71875
		0x80, // VCM_REG_EN: true
		0x40) // VCM_OUT

	disp.writeCmd(CMD_PIXFMT,
		0x55) // DPI/DBI: 16 bits / pixel

	disp.writeCmd(CMD_FRMCTRL1,
		0xa0, // FRS: 60.76  DIVA: 0
		0x11) // RTNA: 17 clocks

	disp.writeCmd(CMD_DISCTRL,
		0x02, // PT: AGND
		0x22, // SS: S960 -> S1  ISC: 5 frames
		0x27) // NL: 8 * (3b + 1) = 320 lines

	disp.writeCmd(CMD_INVON) // make it actually RGB

	// need to mirror the display for picocalc
	disp.writeCmd(CMD_MADCTRL, MADCTRL_MX|MADCTRL_MH)

	disp.writeCmd(CMD_SLPOUT)
	time.Sleep(time.Millisecond * 120)
	disp.writeCmd(CMD_DISON)
}

// writeCmd issues a TFT command with optional data
func (disp *Ili948x) writeCmd(cmd uint8, data ...uint8) {
	disp.startWrite()

	disp.dc.Low() // command mode
	disp.trans.write8(cmd)

	disp.dc.High() // data mode
	disp.trans.write8sl(data)

	disp.endWrite()
}

//go:inline
func (disp *Ili948x) startWrite() {
	if disp.cs != machine.NoPin {
		disp.cs.Low()
	}
}

//go:inline
func (disp *Ili948x) endWrite() {
	if disp.cs != machine.NoPin {
		disp.cs.High()
	}
}

// For compatibility with https://github.com/tinygo-org/drivers/blob/release/displayer.go
func (disp *Ili948x) Display() error {
	return nil
}
