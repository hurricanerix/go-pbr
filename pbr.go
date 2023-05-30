package main

import (
	"fmt"
	"go-pbr/graphics/opengl"
	"go-pbr/graphics/opengl/shader"
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
	objectDir := "assets/objects/brick_cube"
	var cubeMesh *obj.Obj
	if f, err := os.Open(objectDir + "/geometry.obj"); err == nil {
		defer f.Close()
		cubeMesh = obj.Load(f)
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

	renderer := graphics.Renderer{
		WindowWidth:  windowWidth,
		WindowHeight: windowHeight,
	}
	renderer.Init()
	cubemapDir := "assets/cubemaps/castle-zavelstein-cellar"
	renderer.SetCubemap(cubemapDir)

	fmt.Printf("Object Data:\n%s\n", cubeMesh)

	phongShader := opengl.Program{}
	phongShader.CompileShader(shader.VertPhong, opengl.VertexShader)
	phongShader.CompileShader(shader.FragPhong, opengl.FragmentShader)
	phongShader.Link()
	cubeMesh.Bind()
	if err := phongShader.Validate(); err != nil {
		panic(err)
	}
	brickMat := graphics.NewMaterial(phongShader.Handle(), objectDir)

	//phongShader.Use()

	brickMat.Load()

	model := mgl32.Ident4()

	var cameraPos = mgl32.Vec3{0.0, 0.0, 5.0}

	projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0)
	viewMatrix := mgl32.LookAtV(cameraPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	phongShader.SetUniformMatrix4fv("ProjMatrix", projMatrix)
	phongShader.SetUniformMatrix4fv("ViewMatrix", viewMatrix)
	phongShader.SetUniformMatrix4fv("ModelMatrix", model)

	phongShader.SetUniform3fv("ViewPos", cameraPos)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		renderer.Clear(viewMatrix)
		// for each material {
		phongShader.Use()
		phongShader.SetUniformMatrix4fv("ProjMatrix", projMatrix)
		phongShader.SetUniformMatrix4fv("ViewMatrix", viewMatrix)
		phongShader.SetUniformMatrix4fv("ModelMatrix", model)
		phongShader.SetUniform3fv("ViewPos", cameraPos)
		// for each object using material
		cubeMesh.Use(phongShader.Handle())
		brickMat.Use()
		cubeMesh.Draw()
		// } end each object
		// } end each material

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	phongShader.Destroy()
}
