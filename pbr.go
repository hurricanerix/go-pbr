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

var angleX = 0.0
var angleY = 0.0

func main() {
	objectDir := "assets/objects/brick_cube"
	var cubeMesh *obj.Obj
	if f, err := os.Open(objectDir + "/geometry.obj"); err == nil {
		defer f.Close()
		cubeMesh = obj.Load(f)
		cubeMesh.GenNormalRequirements()
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

	var cameraPos = mgl32.Vec3{0.0, 2.0, 5.0}

	projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0)
	viewMatrix := mgl32.LookAtV(cameraPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	phongShader.SetUniformMatrix4fv("ProjMatrix", projMatrix)
	phongShader.SetUniformMatrix4fv("ViewMatrix", viewMatrix)
	phongShader.SetUniformMatrix4fv("ModelMatrix", model)

	phongShader.SetUniform3fv("ViewPos", cameraPos)

	window.SetCursorPosCallback(mousePosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	for !window.ShouldClose() {
		// Update
		xrot := mgl32.HomogRotate3D(float32(angleX), mgl32.Vec3{0, 1, 0})
		//yrot := mgl32.HomogRotate3D(float32(angleY), mgl32.Vec3{1, 0, 0})
		model = xrot //xrot.Mul4(yrot)

		// Render
		renderer.Clear(viewMatrix)
		// for each material {
		phongShader.Use()
		phongShader.SetUniformMatrix4fv("Projection", projMatrix)
		phongShader.SetUniformMatrix4fv("View", viewMatrix)
		phongShader.SetUniformMatrix4fv("Model", model)
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

var currentX float64
var previousX float64
var currentY float64
var previousY float64
var rotateCube bool

func mousePosCallback(w *glfw.Window, xpos float64, ypos float64) {
	previousX = currentX
	currentX = xpos
	previousY = currentY
	currentY = ypos

	if rotateCube {
		speed := 0.05
		dirX := 0.0
		if previousX < currentX {
			dirX = 1.0
		} else if previousX > currentX {
			dirX = -1.0
		}
		angleX += dirX * speed / 2

		dirY := 0.0
		if previousY < currentY {
			dirY = 1.0
		} else if previousY > currentY {
			dirY = -1.0
		}
		angleY += dirY * speed
	}
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		rotateCube = true
	}

	if button == glfw.MouseButtonLeft && action == glfw.Release {
		rotateCube = false
	}
}
