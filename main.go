package main

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"

	d2 "github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
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
var red sdl.Color = sdl.Color{0xFF, 0x00, 0x00, 0xFF}
var blue sdl.Color = sdl.Color{0x00, 0x00, 0xFF, 0xFF}
var white sdl.Color = sdl.Color{0xFF, 0xFF, 0xFF, 0xFF}

// physics
const LEVEL_WIDTH_METERS = WINDOW_WIDTH / 10                         // 800 px window becomes 80 meters in world
const PIXEL_SIZE_METERS = float64(LEVEL_WIDTH_METERS) / SCREEN_WIDTH // width/height of a single pixel in meters

const baseRadius = SCREEN_WIDTH / 16
const bulletRadius = baseRadius / 10
const shieldRadius = SCREEN_HEIGHT / 10

const p1centerX int32 = SCREEN_WIDTH * 0.25
const p1centerY int32 = SCREEN_HEIGHT / 3
const p2centerX int32 = SCREEN_WIDTH * 0.75
const p2centerY int32 = SCREEN_HEIGHT * 2 / 3

const deg2rad = math.Pi / 180

var space *d2.Space

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err := sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	if err != nil {
		panic(err)
	}
	err = sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 16)
	if err != nil {
		panic(err)
	}

	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	_, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}

	// 	gl.Enable(gl.MULTISAMPLE)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	// renderer.SetScale(0.5, 0.5)
	renderer.SetLogicalSize(SCREEN_WIDTH, SCREEN_HEIGHT)

	joystick := sdl.JoystickOpen(0)
	if joystick == nil {
		// panic(sdl.GetError())
	}

	/* Physics */
	space = d2.NewSpace()
	space.Gravity = vect.Vect{X: 0, Y: 0}

	castle1 := d2.NewCircle(vect.Vector_Zero, float32(baseRadius*PIXEL_SIZE_METERS))
	castle1.SetElasticity(0.6)
	staticBody := d2.NewBodyStatic()
	staticBody.SetPosition(vect.Vect{
		X: vect.Float(float64(p1centerX) * PIXEL_SIZE_METERS),
		Y: vect.Float(float64(p1centerY) * PIXEL_SIZE_METERS),
	})
	staticBody.AddShape(castle1)
	space.AddBody(staticBody)

	castle2 := d2.NewCircle(vect.Vector_Zero, float32(baseRadius*PIXEL_SIZE_METERS))
	castle2.SetElasticity(0.6)
	staticBody = d2.NewBodyStatic()
	staticBody.SetPosition(vect.Vect{
		X: vect.Float(float64(p2centerX) * PIXEL_SIZE_METERS),
		Y: vect.Float(float64(p2centerY) * PIXEL_SIZE_METERS),
	})
	staticBody.AddShape(castle2)
	space.AddBody(staticBody)

	bullet := d2.NewCircle(vect.Vector_Zero, float32(bulletRadius*PIXEL_SIZE_METERS))
	bullet.SetElasticity(1.0)
	body := d2.NewBody(vect.Float(1), bullet.Moment(float32(1)))
	body.SetPosition(vect.Vect{
		X: vect.Float(SCREEN_WIDTH / 2 * PIXEL_SIZE_METERS),
		Y: vect.Float(SCREEN_HEIGHT / 2 * PIXEL_SIZE_METERS),
	})
	body.AddShape(bullet)
	space.AddBody(body)
	body.AddVelocity(-10, -10)

	// game loop
	running := true
	var p1 int16 = 0
	var p2 int16 = 0
	ticker := time.NewTicker(time.Second / 60)
	for running {

		space.Step(vect.Float(1.0 / 60))
		bulletPosition := bullet.Body.Position()

		renderer.SetDrawColor(0, 0, 0, 0xFF)
		renderer.Clear()

		var p1loc float32 = ((float32(p1) + MAX_AXIS) / (AXIS_RANGE / 360) * AXIS_SENSITIVITY)
		gfx.FilledPieColor(renderer, int32(float64(castle1.Body.Position().X)/PIXEL_SIZE_METERS), int32(float64(castle1.Body.Position().Y)/PIXEL_SIZE_METERS), shieldRadius, int32(0+p1loc), int32(90+p1loc), red)
		gfx.FilledCircleColor(renderer, int32(float64(castle1.Body.Position().X)/PIXEL_SIZE_METERS), int32(float64(castle1.Body.Position().Y)/PIXEL_SIZE_METERS), baseRadius, grey)

		var p2loc float32 = ((float32(p2) + MAX_AXIS) / (AXIS_RANGE / 360) * AXIS_SENSITIVITY)
		gfx.FilledPieColor(renderer, int32(float64(castle2.Body.Position().X)/PIXEL_SIZE_METERS), int32(float64(castle2.Body.Position().Y)/PIXEL_SIZE_METERS), shieldRadius, int32(0+p2loc), int32(90+p2loc), blue)
		gfx.FilledCircleColor(renderer, int32(float64(castle2.Body.Position().X)/PIXEL_SIZE_METERS), int32(float64(castle2.Body.Position().Y)/PIXEL_SIZE_METERS), baseRadius, grey)

		gfx.FilledCircleColor(renderer, int32(float64(bulletPosition.X)/PIXEL_SIZE_METERS), int32(float64(bulletPosition.Y)/PIXEL_SIZE_METERS), bulletRadius, white)

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
		<-ticker.C
	}
}
