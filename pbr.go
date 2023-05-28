package main

import (
	"fmt"
	"go-pbr/obj"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"go-pbr/opengl"
)

//import (
//	"fmt"
//	"image"
//	"image/draw"
//	_ "image/png"
//	"io"
//	"log"
//	"os"
//	"runtime"
//	"strings"
//
//	"github.com/go-gl/gl/v4.1-core/gl"
//	"github.com/go-gl/glfw/v3.3/glfw"
//	"github.com/go-gl/mathgl/mgl32"
//)

const windowWidth = 960
const windowHeight = 540

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	var o *obj.Obj
	if f, err := os.Open("assets/models/cube.obj"); err == nil {
		defer f.Close()
		o = obj.Load(f)
	}

	fmt.Printf("Object Data:\n%s\n", o)

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "PBR", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// TODO: abstract START
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	slversion := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Println("OpenGL SL version", slversion)
	// TODO: abstract STOP

	prog := opengl.Program{}
	if f, err := os.Open("shaders/basic.vert"); err == nil {
		defer f.Close()
		prog.CompileShader(f, opengl.VertexShader)
	}

	if f, err := os.Open("shaders/basic.frag"); err == nil {
		defer f.Close()
		prog.CompileShader(f, opengl.FragmentShader)
	}

	prog.Link()

	o.Bind(prog.Handle())

	if err := prog.Validate(); err != nil {
		panic(err)
	}

	prog.Use()

	model := mgl32.Ident4()

	prog.SetUniformMatrix4fv("projection", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0))
	prog.SetUniformMatrix4fv("view", mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}))
	prog.SetUniformMatrix4fv("model", model)

	gl.BindFragDataLocation(prog.Handle(), 0, gl.Str("FragColor\x00"))

	// Configure global settings
	gl.Enable(gl.CULL_FACE) //| gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.1, 0.1, 0.2, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		prog.SetUniformMatrix4fv("model", model)
		gl.BindVertexArray(o.Vao)

		gl.DrawArrays(gl.TRIANGLES, 0, o.Indices())

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	prog.Destroy()
}
