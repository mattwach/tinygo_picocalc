package ili948x

import (
	"image/color"
)

type Bitmap struct {
	Data   []RGB565
	Width  int16
	Height int16
}

func (bm *Bitmap) Init(w, h int16) {
	bm.Width = w
	bm.Height = h
	bm.Data = make([]RGB565, w*h)
}

// For compatibility with https://github.com/tinygo-org/drivers/blob/release/displayer.go
// Size returns the current size of the display.
func (bm *Bitmap) Size() (int16, int16) {
	return bm.Width, bm.Height
}

// For compatibility with https://github.com/tinygo-org/drivers/blob/release/displayer.go
// Use SetPixel565 for better performance
func (bm *Bitmap) SetPixel(x, y int16, c color.RGBA) {
	// This can be uncommented for debugging / checking
	//if (x < 0) || (x >= bm.Width) || (y < 0) || (y >= bm.Height) {
	//	panic(fmt.Sprintf("SetPixel x=%v y=%v", x, y))
	//}
	bm.Data[y*bm.Width+x] = RGBA2RGB565(c)
}

func (bm *Bitmap) SetPixel565(x, y int16, c RGB565) {
	bm.Data[y*bm.Width+x] = c
}

func (bm *Bitmap) FillWith(c RGB565) {
	for i := range bm.Data {
		bm.Data[i] = c
	}
}

// For compatibility with https://github.com/tinygo-org/drivers/blob/release/displayer.go
func (bm *Bitmap) Display() error {
	return nil
}
