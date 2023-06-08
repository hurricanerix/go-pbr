package main

import (
	"go-pbr/graphics/opengl"
	"log"
	"math"
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

var modelAngle = 0.0
var cameraAngle = 0.0
var lightAngle = 0.0

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

	// TODO: Move the shader creation into the material initialization so that the shader is always in use
	//       when the textures are loaded.q
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

	lightDistance := float32(5)
	rotLight := mgl32.Vec3{float32(math.Cos(lightAngle)), 0, float32(math.Sin(lightAngle))}
	rotLight = rotLight.Mul(lightDistance)
	phongShader.SetUniform3fv(opengl.LightPosKey, rotLight)

	for !window.ShouldClose() {
		// Update
		rotModel := model.Mul4(mgl32.HomogRotate3DY(float32(modelAngle)))
		rotViewMatrix := viewMatrix.Mul4(mgl32.HomogRotate3DY(float32(cameraAngle)))
		rotLight := mgl32.Vec3{float32(math.Cos(lightAngle)), 0, float32(math.Sin(lightAngle))}
		rotLight = rotLight.Mul(lightDistance)
		// Render
		renderer.Clear(rotViewMatrix)
		// for each material {
		phongShader.Use()
		phongShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, projMatrix)
		phongShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, rotViewMatrix)
		phongShader.SetUniformMatrix4fv(opengl.ModelMatrixKey, rotModel)
		phongShader.SetUniform3fv(opengl.ViewPosKey, cameraPos)
		phongShader.SetUniform3fv(opengl.LightPosKey, rotLight)
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
var rotateCube bool
var rotateCamera bool
var rotateLight bool

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
		modelAngle += dirX * speed / 2
	}

	if rotateCamera {
		speed := 0.05
		dirX := 0.0
		if previousX < currentX {
			dirX = -1.0
		} else if previousX > currentX {
			dirX = 1.0
		}
		cameraAngle += dirX * speed / 2
	}

	if rotateLight {
		speed := 0.05
		dirX := 0.0
		if previousX < currentX {
			dirX = -1.0
		} else if previousX > currentX {
			dirX = 1.0
		}
		lightAngle += dirX * speed / 2
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
		rotateLight = true
	}

	if button == glfw.MouseButtonRight && action == glfw.Release {
		rotateLight = false
	}

	if button == glfw.MouseButtonMiddle && action == glfw.Press {
		rotateCamera = true
	}

	if button == glfw.MouseButtonMiddle && action == glfw.Release {
		rotateCamera = false
	}
}
