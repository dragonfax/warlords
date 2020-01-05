package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 600
const MAX_AXIS = 32767
const MIN_AXIS = -MAX_AXIS
const AXIS_RANGE = MAX_AXIS * 2
const AXIS_WIDTH_STEP = float32(SCREEN_WIDTH) / float32(AXIS_RANGE)
const AXIS_HEIGHT_STEP = float32(SCREEN_HEIGHT) / float32(AXIS_RANGE)

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
	var x int32 = 0
	var y int32 = 0
	for running {

		renderer.SetDrawColor(0, 0, 0, 0xFF)
		renderer.Clear()

		rect := sdl.Rect{X: x, Y: y, W: 200, H: 200}
		renderer.SetDrawColor(0xFF, 0, 0, 0xFF)
		renderer.FillRect(&rect)

		renderer.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.JoyAxisEvent:
				fmt.Printf("axis %d %06d\n", e.Axis, e.Value)
				if e.Axis == 0 {
					x = int32(float32(int32(e.Value)+MAX_AXIS) * AXIS_WIDTH_STEP)
					fmt.Printf("new x %d\n", x)
				} else if e.Axis == 1 {
					y = int32(float32(int32(e.Value)+MAX_AXIS) * AXIS_HEIGHT_STEP)
					fmt.Printf("new y %d\n", y)
				}
			}
		}
	}
}
