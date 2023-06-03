package main

import (
	"go-pbr/graphics/opengl"
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

	objectDir := "assets/objects"
	var cubeMesh *obj.Obj
	if f, err := os.Open(objectDir + "/cube.obj"); err == nil {
		defer f.Close()
		cubeMesh = obj.Load(f)
		cubeMesh.GenNormalRequirements()
	} else {
		panic(err)
	}

	cubemapDir := "assets/cubemaps/castle-zavelstein-cellar"
	renderer.SetCubemap(cubemapDir)

	//fmt.Printf("Object Data:\n%s\n", cubeMesh)

	phongShader := opengl.Program{}
	//if f, err := os.Open("assets/shaders/lgl5.4.vert"); err == nil {
	if f, err := os.Open("assets/shaders/phong.vert"); err == nil {
		defer f.Close()
		phongShader.CompileShader(f, opengl.VertexShader)
	}

	//if f, err := os.Open("assets/shaders/lgl5.4.frag"); err == nil {
	if f, err := os.Open("assets/shaders/phong.frag"); err == nil {
		defer f.Close()
		phongShader.CompileShader(f, opengl.FragmentShader)
	}

	phongShader.Link()
	phongShader.Use() // <- Important to use before loading the material.
	matDir := "assets/materials"
	//brickMat := graphics.NewMaterial(phongShader.Handle(), matDir+"/lgl_brickwall")
	brickMat := graphics.NewMaterial(phongShader.Handle(), matDir+"/stone_wall")

	brickMat.Load()

	cubeMesh.Bind()
	if err := phongShader.Validate(); err != nil {
		panic(err)
	}

	model := mgl32.Ident4()

	var cameraPos = mgl32.Vec3{0.0, 2.0, 5.0}

	projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0)
	viewMatrix := mgl32.LookAtV(cameraPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	phongShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, projMatrix)
	phongShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, viewMatrix)
	phongShader.SetUniformMatrix4fv(opengl.ModelMatrixKey, model)

	phongShader.SetUniform3fv(opengl.ViewPosKey, cameraPos)

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
		phongShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, projMatrix)
		phongShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, viewMatrix)
		phongShader.SetUniformMatrix4fv(opengl.ModelMatrixKey, model)
		phongShader.SetUniform3fv(opengl.ViewPosKey, cameraPos)
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
