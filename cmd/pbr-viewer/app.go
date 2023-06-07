package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"go-pbr/graphics"
	"go-pbr/scene"
	"runtime"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var modelAngleX = 0.0
var cameraAngleX = 0.0

type app struct {
	Title string

	Monitor      int
	WindowWidth  int
	WindowHeight int

	Scene scene.Scene

	window  *glfw.Window
	backend graphics.Backend

	currentTime  float64
	previousTime float64
}

func (a *app) Init() error {
	if err := glfw.Init(); err != nil {
		return err
	}
	// TODO: defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	//monitors := glfw.GetMonitors()
	var err error
	a.window, err = glfw.CreateWindow(a.WindowWidth, a.WindowHeight, a.Title, nil, nil)
	if err != nil {
		panic(err)
	}
	a.window.MakeContextCurrent()

	a.backend = graphics.Backend{}
	if err := a.backend.Init(); err != nil {
		return err
	}

	a.window.SetCursorPosCallback(mousePosCallback)
	a.window.SetMouseButtonCallback(mouseButtonCallback)

	return nil
}

func (a *app) Run() error {
	for _, o := range a.Scene {
		if err := o.Bind(); err != nil {
			return err
		}
	}

	for !a.window.ShouldClose() {
		a.backend.Clear()
		dt := a.DeltaTime()

		for _, o := range a.Scene {
			if err := o.Update(dt); err != nil {
				return err
			}
		}

		for _, o := range a.Scene {
			if err := o.Draw(); err != nil {
				return err
			}
		}

		a.window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

var currentX float64
var previousX float64
var rotateCube bool

var rotateCamera bool

func mousePosCallback(w *glfw.Window, xpos float64, ypos float64) {
	previousX = currentX
	currentX = xpos

	if rotateCube {
		speed := 0.05
		dirX := 0.0
		if previousX < currentX {
			dirX = 1.0
		} else if previousX > currentX {
			dirX = -1.0
		}
		modelAngleX += dirX * speed / 2
	}

	if rotateCamera {
		speed := 0.05
		dirX := 0.0
		if previousX < currentX {
			dirX = 1.0
		} else if previousX > currentX {
			dirX = -1.0
		}
		cameraAngleX += dirX * speed / 2
	}
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		rotateCube = true
	}

	if button == glfw.MouseButtonLeft && action == glfw.Release {
		rotateCube = false
	}

	if button == glfw.MouseButtonRight && action == glfw.Press {
		rotateCamera = true
	}

	if button == glfw.MouseButtonRight && action == glfw.Release {
		rotateCamera = false
	}
}

func (a *app) ShouldClose() bool {
	return a.window.ShouldClose()
}

func (a *app) DeltaTime() float32 {
	a.currentTime = glfw.GetTime()
	dt := float32(a.currentTime - a.previousTime)
	a.previousTime = a.currentTime
	return dt
}

func (a *app) PostLoop() error {
	a.window.SwapBuffers()
	glfw.PollEvents()
	return nil
}

func (a *app) Free() error {
	glfw.Terminate()
	return nil
}
