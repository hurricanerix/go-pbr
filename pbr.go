package main

import (
	"fmt"
	"go-pbr/graphics/opengl"
	shader2 "go-pbr/graphics/opengl/shader"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"go-pbr/graphics"
	"go-pbr/obj"
)

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

	renderer := graphics.Renderer{}
	renderer.Init()

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

	brick.Load()

	model := mgl32.Ident4()

	var cameraPos = mgl32.Vec3{0.0, 0.0, 5.0}

	prog.SetUniformMatrix4fv("ProjMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0))
	prog.SetUniformMatrix4fv("ViewMatrix", mgl32.LookAtV(cameraPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}))
	prog.SetUniformMatrix4fv("ModelMatrix", model)

	prog.SetUniform3fv("ViewPos", cameraPos)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		renderer.Clear()

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		prog.SetUniformMatrix4fv("ModelMatrix", model)
		// for each material {
		brick.Use()
		// for each object using material
		//o.Bind()
		// } end each object
		// } end each material

		renderer.Draw(o.Indices())

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	prog.Destroy()
}
