package ili948x

import (
	"machine"
)

type spiTransport struct {
	spi machine.SPI // spi bus
	buf []uint8     // spi data buffer
}

func NewSPITransport(spi machine.SPI) *spiTransport {
	return &spiTransport{
		spi: spi,
		buf: make([]uint8, 1024),
	}
}

// 8 bit
func (st *spiTransport) write8(data uint8) {
	st.buf[0] = data
	st.spi.Tx(st.buf[:1], nil)
}

func (st *spiTransport) write8sl(data []uint8) {
	j := 0
	for i := 0; i < len(data); i++ {
		st.buf[j] = data[i]
		j++
		if j >= len(st.buf) {
			st.spi.Tx(st.buf, nil)
			j = 0
		}
	}
	if j > 0 {
		st.spi.Tx(st.buf[:j], nil)
	}
}

// 16 bit
func (st *spiTransport) write16(data uint16) {
	st.buf[0] = uint8(data >> 8)
	st.buf[1] = uint8(data)
	st.spi.Tx(st.buf[:2], nil)
}

func (st *spiTransport) write16n(data uint16, n int) {
	buflen := n * 2
	if buflen > len(st.buf) {
		buflen = len(st.buf)
	}
	b0 := uint8(data >> 8)
	b1 := uint8(data)
	for i := 0; i < buflen; i++ {
		st.buf[i] = b0
		i++
		st.buf[i] = b1
	}
	bytesWritten := 0
	for ((n * 2) - bytesWritten) > len(st.buf) {
		st.spi.Tx(st.buf, nil)
		bytesWritten += len(st.buf)
	}
	st.spi.Tx(st.buf[:(n*2)-bytesWritten], nil)
}

func (st *spiTransport) writeRGB565(data []RGB565) {
	j := 0
	for i := 0; i < len(data); i++ {
		st.buf[j] = uint8(data[i] >> 8)
		j++
		st.buf[j] = uint8(data[i])
		j++
		if j >= len(st.buf) {
			st.spi.Tx(st.buf, nil)
			j = 0
		}
	}
	if j > 0 {
		st.spi.Tx(st.buf[:j], nil)
	}
}
