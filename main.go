package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600
const SCREEN_WIDTH = WINDOW_WIDTH
const SCREEN_HEIGHT = WINDOW_HEIGHT

const MAX_AXIS = math.MaxInt16
const MIN_AXIS = math.MinInt16
const AXIS_RANGE = MAX_AXIS - MIN_AXIS
const AXIS_SENSITIVITY = 2

var grey sdl.Color = sdl.Color{0x80, 0x80, 0x80, 0xFF}
var red sdl.Color = sdl.Color{0xff, 0x00, 0x00, 0xFF}
var blue sdl.Color = sdl.Color{0x00, 0x00, 0xff, 0xFF}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	// renderer.SetScale(0.5, 0.5)
	renderer.SetLogicalSize(SCREEN_WIDTH, SCREEN_HEIGHT)

	joystick := sdl.JoystickOpen(0)
	if joystick == nil {
		panic(sdl.GetError())
	}

	running := true
	var p1 int16 = 0
	var p2 int16 = 0
	for running {

		renderer.SetDrawColor(0, 0, 0, 0xFF)
		renderer.Clear()

		var baseRadius int32 = SCREEN_WIDTH / 16
		var shieldRadius int32 = SCREEN_HEIGHT / 10

		var p1centerX int32 = SCREEN_WIDTH * 0.25
		var p1centerY int32 = SCREEN_HEIGHT / 3
		var p1loc float32 = ((float32(p1) + MAX_AXIS) / (AXIS_RANGE / 360) * AXIS_SENSITIVITY)
		gfx.FilledPieColor(renderer, p1centerX, p1centerY, shieldRadius, int32(0+p1loc), int32(90+p1loc), red)
		gfx.FilledCircleColor(renderer, p1centerX, p1centerY, baseRadius, grey)

		var p2centerX int32 = SCREEN_WIDTH * 0.75
		var p2centerY int32 = SCREEN_HEIGHT * 2 / 3
		var p2loc float32 = ((float32(p2) + MAX_AXIS) / (AXIS_RANGE / 360) * AXIS_SENSITIVITY)
		gfx.FilledPieColor(renderer, p2centerX, p2centerY, shieldRadius, int32(0+p2loc), int32(90+p2loc), blue)
		gfx.FilledCircleColor(renderer, p2centerX, p2centerY, baseRadius, grey)

		renderer.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.JoyAxisEvent:
				if e.Axis == 0 {
					p1 = e.Value
				} else if e.Axis == 1 {
					p2 = e.Value
				}
			}
		}
	}
}
