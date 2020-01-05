package main

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 600
const MAX_AXIS = 32767
const MIN_AXIS = -MAX_AXIS
const AXIS_RANGE = MAX_AXIS * 2
const AXIS_WIDTH_STEP = float32(SCREEN_WIDTH) / float32(AXIS_RANGE)
const AXIS_HEIGHT_STEP = float32(SCREEN_HEIGHT) / float32(AXIS_RANGE)

var grey sdl.Color = sdl.Color{0x80, 0x80, 0x80, 0xFF}
var red sdl.Color = sdl.Color{0xff, 0x00, 0x00, 0xFF}
var blue sdl.Color = sdl.Color{0x00, 0x00, 0xff, 0xFF}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

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

		gfx.CircleColor(renderer, 200, 200, 50, grey)
		var p1loc float32 = (float32(p1) + MAX_AXIS) / (AXIS_RANGE / 360)
		gfx.ArcColor(renderer, 200, 200, 70, int32(0+p1loc), int32(90+p1loc), red)

		gfx.CircleColor(renderer, 600, 400, 50, grey)
		var p2loc float32 = (float32(p2) + MAX_AXIS) / (AXIS_RANGE / 360)
		gfx.ArcColor(renderer, 600, 400, 70, int32(0+p2loc), int32(90+p2loc), blue)

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
