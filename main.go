package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"

	d2 "github.com/neguse/go-box2d-lite/box2dlite"
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

var gravity = d2.MakeB2Vec2(0.0, -10.0)
var world = d2.MakeB2World(gravity)

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

	// physics
	castle1BodyDef := d2.B2BodyDef{}
	castle1BodyDef.Position.Set(float64(p1centerX)*PIXEL_SIZE_METERS, float64(p1centerY)*PIXEL_SIZE_METERS)
	castle1Body := world.CreateBody(&castle1BodyDef)
	castle1Shape := d2.B2CircleShape{}
	castle1Shape.SetRadius(baseRadius * PIXEL_SIZE_METERS)
	castle1Body.CreateFixture(castle1Shape, 0.0)

	castle2BodyDef := d2.B2BodyDef{}
	castle2BodyDef.Position.Set(float64(p2centerX)*PIXEL_SIZE_METERS, float64(p2centerY)*PIXEL_SIZE_METERS)
	castle2Body := world.CreateBody(&castle2BodyDef)
	castle2Shape := d2.B2CircleShape{}
	castle2Shape.SetRadius(baseRadius * PIXEL_SIZE_METERS)
	castle2Body.CreateFixture(castle2Shape, 0.0)

	bulletBodyDef := d2.B2BodyDef{Type: d2.B2BodyType.B2_dynamicBody}
	bulletBodyDef.Bullet = true
	bulletBodyDef.Position.Set(SCREEN_WIDTH/2*PIXEL_SIZE_METERS, SCREEN_HEIGHT/2*PIXEL_SIZE_METERS)
	bulletBody := world.CreateBody(&bulletBodyDef)
	bulletShape := d2.B2CircleShape{}
	bulletShape.SetRadius(bulletRadius)
	bulletFixture := d2.B2FixtureDef{}
	bulletFixture.Shape = bulletShape
	bulletFixture.Density = 1.0
	bulletFixture.Friction = 0.0
	bulletBody.CreateFixtureFromDef(&bulletFixture)

	v := bulletBody.GetWorldVector(d2.B2Vec2{X: -1, Y: -1})
	v.OperatorScalarMulInplace(bulletBody.GetMass() * 10)
	fmt.Printf("mass: %f\n", bulletBody.GetMass())
	bulletBody.ApplyLinearImpulseToCenter(d2.B2Vec2{X: -10, Y: -10}, true)

	timestep := 1.0 / 60.0
	velocityIterations := 6
	positionIterations := 2

	// game loop
	running := true
	var p1 int16 = 0
	var p2 int16 = 0
	for running {

		world.Step(timestep, velocityIterations, positionIterations)
		bulletPosition := bulletBody.GetPosition()
		fmt.Printf("bullet: %v\n", bulletPosition)

		renderer.SetDrawColor(0, 0, 0, 0xFF)
		renderer.Clear()

		var p1loc float32 = ((float32(p1) + MAX_AXIS) / (AXIS_RANGE / 360) * AXIS_SENSITIVITY)
		gfx.FilledPieColor(renderer, p1centerX, p1centerY, shieldRadius, int32(0+p1loc), int32(90+p1loc), red)
		gfx.FilledCircleColor(renderer, p1centerX, p1centerY, baseRadius, grey)

		var p2loc float32 = ((float32(p2) + MAX_AXIS) / (AXIS_RANGE / 360) * AXIS_SENSITIVITY)
		gfx.FilledPieColor(renderer, p2centerX, p2centerY, shieldRadius, int32(0+p2loc), int32(90+p2loc), blue)
		gfx.FilledCircleColor(renderer, p2centerX, p2centerY, baseRadius, grey)

		gfx.FilledCircleColor(renderer, int32(bulletPosition.X/PIXEL_SIZE_METERS), int32(bulletPosition.Y/PIXEL_SIZE_METERS), bulletRadius, white)

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
