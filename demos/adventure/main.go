package main

import (
	"github.com/aykevl/tinygl/pixel"
)

const numLeds = 18
const maxFrame = 36

var leds [numLeds]pixel.LinearGRB888

var enemyFrame uint8
var attackingFrame uint8

var attacking bool
var oldPressed bool
var pressed bool
var playerPos uint8
var enemyPos uint8

const (
	RED = iota
	ORANGE
	YELLOW
	GREEN
	BLUE
	PURPLE
	WHITE
	BLACK
	ATTACK
)

var colors = [9]pixel.LinearGRB888{
	{R: 0xff / 3, G: 0x00 / 3, B: 0x00 / 3}, // red
	{R: 0xff / 3, G: 0x22 / 3, B: 0x00 / 3}, // orange
	{R: 0x88 / 3, G: 0xff / 3, B: 0x00 / 3}, // yellow
	{R: 0x00 / 3, G: 0xff / 3, B: 0x00 / 3}, // green
	{R: 0x00 / 3, G: 0x00 / 3, B: 0xff / 3}, // blue
	{R: 0x80 / 3, G: 0x00 / 3, B: 0x80 / 3}, // purple
	{R: 0x33 / 3, G: 0xcc / 3, B: 0xff / 3}, // white
	{R: 0x00 / 3, G: 0x00 / 3, B: 0x00 / 3}, // black
	{R: 0x00 / 3, G: 0xff / 3, B: 0x88 / 3}, // green + blue
}

func main() {
	initHardware()
	attackingFrame = 0
	attacking = false
	enemyPos = numLeds - 1

	//ledIndex := uint8(0)
	for {
		gameFrame()
		updateLEDs()

		// Update 2 LEDs.
		//for i := uint8(0); i < 2; i++ {

		//			ledIndex++
		//			if ledIndex >= 18 {
		//				ledIndex = 0
		pressed = isButtonPressed()

		if pressed && pressed != oldPressed && !attacking {
			attacking = true
		}

		oldPressed = pressed
		//			}
		//}
	}
}

func gameFrame() {

	for i := uint8(0); i < numLeds; i++ {
		leds[i] = colors[BLACK]
	}

	leds[playerPos] = colors[GREEN]
	leds[16] = colors[YELLOW]
	leds[17] = colors[YELLOW]

	if enemyFrame > 2*maxFrame {
		enemyPos--
		enemyFrame = 0
		if enemyPos <= 0 {
			enemyPos = numLeds - 1

		}
	}
	leds[enemyPos] = colors[RED]

	if attacking {
		attackingFrame++
		leds[playerPos+1] = colors[ATTACK]

		if enemyPos == (playerPos + 1) {
			// ENEMY DIE
			playerPos++
			enemyFrame = 0
			enemyPos = numLeds - 1
		}

		if attackingFrame > maxFrame {
			attackingFrame = 0
			attacking = false
		}
	}

	if enemyPos == playerPos {
		// PLAYER DIE
		playerPos = 0
		enemyFrame = 0
		enemyPos = numLeds - 1

		gameOver()
	}

	if playerPos == 16 {
		gameWin()
	}

	enemyFrame++
}

func gameOver() {
	for k := 0; k < 4; k++ {
		for i := uint8(0); i < numLeds; i++ {
			leds[i] = colors[BLACK]
		}
		for i := 0; i < 360; i++ {
			updateLEDs()
		}
		for i := uint8(0); i < numLeds; i++ {
			leds[i] = colors[RED]
		}
		for i := 0; i < 360; i++ {
			updateLEDs()
		}
	}
}

func gameWin() {
	for k := 0; k < 4; k++ {
		for i := uint8(0); i < numLeds; i++ {
			leds[i] = colors[BLACK]
		}
		for i := 0; i < 360; i++ {
			updateLEDs()
		}
		for i := uint8(0); i < numLeds; i++ {
			leds[i] = colors[GREEN]
		}
		for i := 0; i < 360; i++ {
			updateLEDs()
		}
	}
}

var xorshift32State uint32 = 1

func random() uint32 {
	x := xorshift32State
	x = xorshift32(x)
	xorshift32State = x
	return x
}

// Xorshift32 RNG. The usual algorithm uses the shift amounts [13, 17, 5], but
// [7, 1, 9] as used below are a much better fit for AVR. It is "only" 37
// instructions (excluding return) when compiled with Clang 16 while the usual
// algorithm uses 57 instructions.
// On the other hand, avr-gcc (tested most versions starting with 5.4.0) is just
// terrible in both cases, using loops for these shifts.
//
// The [7, 1, 9] algorithm is mentioned on page 9 of this paper:
// http://www.iro.umontreal.ca/~lecuyer/myftp/papers/xorshift.pdf
//
// Godbolt reference (both algorithms in avr-gcc and Clang 16):
// https://godbolt.org/z/KdeKhP54d
func xorshift32(x uint32) uint32 {
	x ^= x << 7
	x ^= x >> 1
	x ^= x << 9
	return x
}
