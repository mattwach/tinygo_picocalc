# Introduction

This project demonstrates how to use the PicoCalc with TinyGo. Specifically, it
writes a string to the LCD, then displays the characters you type.  It
intentionally uses as little code as possible so that it's easy to use as
a starting point for experimentation.

check out [main.go](main.go) to see the code.

# Prerequisites

- Install TinyGo using their [official instructions](https://tinygo.org/getting-started/install/).
- Optional (recommended): Do the [blinking light tutorial](https://tinygo.org/tour/blink/onboard/)
  with a pico on a breadboard.  The flash command is `tinygo flash -target=pico` (or `-target=pico2`)
- If you are new to Go or want a refresher, [Tour of Go](https://go.dev/tour/) can help.

# Flash

```bash
tinygo flash -target=pico
```

Use `-target=pico2` if you are using a Pico 2.  If it works, you'll see this:

![picocalc](img/picocalc.jpg) 

# More Examples

Here is variant that prints "hello world" with no error checking:

```golang
package main

import (
	"image/color"
	"picocalc/ili948x"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/shnm"
)


func main() {
	lcd := ili948x.InitDisplay()
	tinyfont.WriteLine(
		lcd,
		&shnm.Shnmk12,
		130,
		100,
		"hello world",
		color.RGBA{255, 255, 255, 255})
	for {
	}
}
```

and here is an minimal version that echos PicoCalc keystrokes to the serial console.
Compile it with `-serial=uart`, then use 'tinygo monitor' to view the output (via USB-C).
The PicoCalc must be "on" for the Keyboard to function.

```golang
package main

import (
	"picocalc/i2ckbd"
	"time"
)

func main() {
	var keyboard i2ckbd.I2CKbd
	_ = keyboard.Init()
	for {
		k, _ := keyboard.GetChar()
		if k != 0 {
			print(string(rune(k)))
		}
		time.Sleep(20 * time.Millisecond)
	}
}

```


# More Tips

## Flashing New Code

Flashing the Pico inside the PicoCalc was not made as easy as it could be.
Some people created some advanced solutions with 3D printing, etc. Here I
present a simple hack.

First, I superglued a SMD button to the Pico and soldered a jumper wire to
the reset pin, like this:

![picocalc](img/hacked_pico.jpg)

Now the Pico has a reset button like it arguably always should have and you
can press reset while holding boot to go into programming mode.

For the next step, I simply drilled holes in the PicoCalc case where the buttons
are.

![picocalc](img/hacked_picocalc.jpg)

Now to program, I do this:

1. Turn the picocalc upside down.  Do not turn it on (see below).
2. Plug the micro USB programming cable between the PC and Pico
3. Using two tools (chopsticks, hex wrench, etc), hold down the boot button and press reset
4. Your Pico should mount as a USB drive on the PC.
5. Run `tingo flash -target=pico` (or whatever programming command you need)
6. Unplug the micro USB before turning on the pico.

The last step is important.  My PicoCalc hardware revision has a bug where if
the PicoCalc is on with the micro USB attached, it will feed 5V to the 18650
batteries - this can overcharge these batteries and is best avoided.

## Serial Communications

As said above, you have to use the micro USB to program the PI Pico but should
not use it as a serial console, due to the 5V charging hardware issue.  The
work-around is to compile like this:

    tinygo flash -target=pico -serial=uart

Now you can use the USB-C port for serial communications.  This link shows
the basics of how.  Note that the Pico Hardware tries to change
the 18650 batteries when you use this port so the current draw can be high.
Running with no batteries is probably the safest option.

## Debugging

If your go code panics (illegal array access, out of memory), it will dump the
panic address over the UART. Once you get the address, you may wonder what to do
with it. The answer is to dissasemble the firmware so you can see what the
address to pointing to. The steps are:

    # replace target with pico2 if needed
    tinygo build -target=pico

    # your elf target might have a different name
    objdump -d picocalc.elf > picocalc.asm

Now you can see addresses in the `asm` file and find out what function threw
the panic.

If you want to try attaching a full debugger (which I have not gotten to yet),
instructions are [here](https://tinygo.org/docs/guides/debugging/).

## Cross compilation

I have a bigger TinyGo project which is a programmable scientific calculator
which can run on either PC (compiled with traditional go) or the PicoCalc (using TinyGo):
[RPNGO]()

![pc calc](img/rpngo_pc.png)
![picocalc calc](img/rpngo_picocalc.png)

You can check out the project sources for more in-depth go usage examples.
The main thing I'll talk about here is Go's use of build tags. The basic
pattern that you use [go build tags](https://pkg.go.dev/go/build) to define
files that will compile differently on PC and PicoCalc.  For example,
say you want to print to the screen or LCD.  You could make a PC version,
`printpc.go`:

```golang
//go:build !pico && !pico2

package console

func Print(msg string) {
	print(msg)
}
```

and a picocalc version in the same directory, `printpicocalc.go`:

```golang
//go:build pico && pico2

package console

func Print(msg string) {
	// do it the PicoCalc way...
}
```

You can also consider using [go interfaces](https://gobyexample.com/interfaces).
These are useful if you want to create objects with state and have those
objects implemented differently on PC vs PicoCalc.


