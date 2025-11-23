// Package I2cKbd creates an interface to the keyboard of the
// picocalc
package i2ckbd

import (
	"fmt"
	"machine"
)

// IMPORTANT:
//
// 1. The PicoCalc must be powered on for the i2c keyboard chip to be active.
// I wasted a bit of time discovering this!
//
// 2. Do not use batteries in the PicoCalc while it's pluggen in via USB mini and turned on
// becuase an electrical path is opened that causes the 18650 batteries to be charged
// beyond 4.2 volts.  Hopefully they fix this hardware flaw.
var i2cKbdAddr uint16 = 0x1F

const i2cGetKey = 0x09

const (
	ALT_KEY       byte = 0xA1
	BACKSPACE_KEY      = 0x08
	CTRL_KEY      byte = 0xA5
	DEL_KEY       byte = 0xd4
	END_KEY       byte = 0xd5
	ESC_KEY       byte = 0xb1
	F1_KEY        byte = 0x81
	F2_KEY        byte = 0x82
	F3_KEY        byte = 0x83
	F4_KEY        byte = 0x84
	F5_KEY        byte = 0x85
	F6_KEY        byte = 0x86
	F7_KEY        byte = 0x87
	F8_KEY        byte = 0x88
	F9_KEY        byte = 0x89
	F10_KEY       byte = 0x90 // odd it's not 0x8A
	HOME_KEY      byte = 0xd2
	INS_KEY       byte = 0xd1

	LEFT_KEY  byte = 0xb4
	RIGHT_KEY byte = 0xb7
	UP_KEY    byte = 0xb5
	DOWN_KEY  byte = 0xb6
)

type I2CKbd struct {
	i2c      *machine.I2C
	write    []byte
	read     []byte
	AltDown  bool
	CtrlDown bool
}

// Init initialized the i2c driver.  It may be necessary to add the ability to
// provided an i2c driver if the bus is shared (I don't believe it is currently).
func (kbd *I2CKbd) Init() error {
	kbd.write = make([]byte, 1)
	kbd.write[0] = i2cGetKey
	kbd.read = make([]byte, 2)
	kbd.i2c = machine.I2C1
	return kbd.i2c.Configure(machine.I2CConfig{
		SCL: machine.GP7,
		SDA: machine.GP6,
	})
}

// GetChar returns a keppress using the key driver codes.  Returns zero if
// nothing was pressed (which will be most of the time).
//
// You need to call this often.  Calling it in a gorouting seems like a good
// plan, but should be done after the basics are fully sorted.
func (kbd *I2CKbd) GetChar() (byte, error) {
	err := kbd.i2c.Tx(i2cKbdAddr, kbd.write, kbd.read)
	if err != nil {
		return 0, err
	}
	if (kbd.read[0] == 0) && (kbd.read[1] == 0) {
		return 0, nil
	}
	switch kbd.read[0] {
	case 0x01:
		return kbd.keyDown()
	case 0x02:
		return kbd.keyHeld()
	case 0x03:
		return kbd.keyUp()
	default:
		return 0, fmt.Errorf("unknown key response: %v", kbd.read[0])
	}
}

// called when a key is depressed
func (kbd *I2CKbd) keyDown() (byte, error) {
	k := kbd.read[1]
	switch k {
	case ALT_KEY:
		kbd.AltDown = true
		return 0, nil
	case CTRL_KEY:
		kbd.CtrlDown = true
		return 0, nil
	default:
		return k, nil
	}
}

// Sometimes called when a key is held.  Usually just for modifier keys.
func (kbd *I2CKbd) keyHeld() (byte, error) {
	switch kbd.read[1] {
	case ALT_KEY:
		// likely not needed, but doesn't hurt anything either
		kbd.AltDown = true
		return 0, nil
	case CTRL_KEY:
		// likely not needed, but doesn't hurt anything either
		kbd.CtrlDown = true
		return 0, nil
	default:
		return 0, nil
	}
}

// Called when a key is released.  We mostly don't care outside of modifier keys
func (kbd *I2CKbd) keyUp() (byte, error) {
	switch kbd.read[1] {
	case ALT_KEY:
		kbd.AltDown = false
		return 0, nil
	case CTRL_KEY:
		kbd.CtrlDown = false
		return 0, nil
	default:
		return 0, nil
	}
}
