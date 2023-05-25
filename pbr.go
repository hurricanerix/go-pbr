package main

import (
	"fmt"
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

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
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

	//prog.Use()

	prog.SetUniform("projection", &mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0))
	camera := mgl32.Vec3{3, 3, 3}
	prog.SetUniform("camera", &camera)
	view := mgl32.LookAtV(camera, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	prog.SetUniform("view", &view)
	model := mgl32.Ident4()
	prog.SetUniform("model", &model)
	light := mgl32.Vec3{2, 0, 2}
	prog.SetUniform("light", &light)
	prog.SetUniform("texSampler", 0)
	prog.SetUniform("armSampler", 1)
	prog.SetUniform("dispSampler", 2)
	prog.SetUniform("norSampler", 3)

	gl.BindFragDataLocation(prog.Handle(), 0, gl.Str("FragColor\x00"))

	// Load the textureDiffuse
	//textureDiffuse, err := opengl.NewTexture("textures/concrete_brick_wall_001_diffuse_1k.png", gl.TEXTURE0)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//textureArm, err := opengl.NewTexture("textures/concrete_brick_wall_001_arm_1k.png", gl.TEXTURE1)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//textureDisp, err := opengl.NewTexture("textures/concrete_brick_wall_001_disp_1k.png", gl.TEXTURE2)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//textureNor, err := opengl.NewTexture("textures/concrete_brick_wall_001_nor_gl_1k.png", gl.TEXTURE3)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(prog.Handle(), gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	//texCoordAttrib := uint32(gl.GetAttribLocation(prog.Handle(), gl.Str("vertTexCoord\x00")))
	//gl.EnableVertexAttribArray(texCoordAttrib)
	//gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(6.0, 5.0, 1.0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render

		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.BindVertexArray(vao)

		if err := prog.Validate(); err != nil {
			panic(err)
		}

		//gl.ActiveTexture(gl.TEXTURE0)
		//gl.BindTexture(gl.TEXTURE_2D, textureDiffuse)
		//
		//gl.ActiveTexture(gl.TEXTURE1)
		//gl.BindTexture(gl.TEXTURE_2D, textureArm)
		//
		//gl.ActiveTexture(gl.TEXTURE2)
		//gl.BindTexture(gl.TEXTURE_2D, textureDisp)
		//
		//gl.ActiveTexture(gl.TEXTURE3)
		//gl.BindTexture(gl.TEXTURE_2D, textureNor)

		prog.Use()

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	prog.Destroy()
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
