package main

import (
	"fmt"
	"go-pbr/graphics"
	"go-pbr/obj"
	shader2 "go-pbr/opengl/shader"
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
	dir := "assets/brick_cube"

	var o *obj.Obj
	if f, err := os.Open(dir + "/geometry.obj"); err == nil {
		defer f.Close()
		o = obj.Load(f)
	} else {
		panic(err)
	}

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

	fmt.Printf("Object Data:\n%s\n", o)

	prog := opengl.Program{}
	prog.CompileShader(shader2.Vert, opengl.VertexShader)
	prog.CompileShader(shader2.Frag, opengl.FragmentShader)

	prog.Link()

	o.Bind(prog.Handle())

	if err := prog.Validate(); err != nil {
		panic(err)
	}

	brick := graphics.NewMaterial(prog.Handle(), dir)

	prog.Use()

	ambiantColor := mgl32.Vec3{1.0, 1.0, 1.0}
	ambiantColorUniform := gl.GetUniformLocation(prog.Handle(), gl.Str("AmbiantColor\x00"))
	gl.Uniform3fv(ambiantColorUniform, 1, &ambiantColor[0])

	lightPos := mgl32.Vec3{0.0, 1.0, 1.0}
	lightPosUniform := gl.GetUniformLocation(prog.Handle(), gl.Str("LightPos\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])

	// 255,244,161
	lightColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	lightColorUniform := gl.GetUniformLocation(prog.Handle(), gl.Str("LightColor\x00"))
	gl.Uniform3fv(lightColorUniform, 1, &lightColor[0])

	lightPower := float32(10)
	lightPowerUniform := gl.GetUniformLocation(prog.Handle(), gl.Str("LightPower\x00"))
	gl.Uniform1f(lightPowerUniform, lightPower)

	brick.Load()

	model := mgl32.Ident4()

	prog.SetUniformMatrix4fv("ProjMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0))
	prog.SetUniformMatrix4fv("ViewMatrix", mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}))
	prog.SetUniformMatrix4fv("ModelMatrix", model)

	gl.BindFragDataLocation(prog.Handle(), 0, gl.Str("FragColor\x00"))

	// Configure global settings
	gl.Enable(gl.CULL_FACE) //| gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

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
		prog.SetUniformMatrix4fv("ModelMatrix", model)
		gl.BindVertexArray(o.Vao)

		brick.Use()

		gl.DrawArrays(gl.TRIANGLES, 0, o.Indices())

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	prog.Destroy()
}
