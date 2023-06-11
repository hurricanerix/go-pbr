package main

import (
	"flag"
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics"
	"go-pbr/obj"
	"log"
	"math"
	"os"
	"runtime"
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
	assetPath := flag.String("assets", "assets", "")
	meshName := flag.String("mesh", "sphere", "")
	materialName := flag.String("material", "metal_plate", "")
	shaderName := flag.String("shader", "lgl-pbr", "")
	skymapName := flag.String("skymap", "lgl", "")

	lightAngle += 90

	flag.Parse()

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

	var subject *obj.Obj
	if f, err := os.Open(fmt.Sprintf("%s/meshes/%s.obj", *assetPath, *meshName)); err == nil {
		defer f.Close()
		subject = obj.Load(f)
		subject.GenNormalRequirements()
	} else {
		panic(err)
	}
	subjectColor := mgl32.Vec3{1.0, 1.0, 1.0}

	renderer.SetCubemap(fmt.Sprintf("%s/cubemaps/%s", *assetPath, *skymapName))

	subjectProgram := graphics.Program{}
	if f, err := os.Open(fmt.Sprintf("%s/shaders/%s/shader.vert", *assetPath, *shaderName)); err == nil {
		defer f.Close()
		subjectProgram.CompileShader(f, graphics.VertexShader)
	}

	if f, err := os.Open(fmt.Sprintf("%s/shaders/%s/shader.frag", *assetPath, *shaderName)); err == nil {
		defer f.Close()
		subjectProgram.CompileShader(f, graphics.FragmentShader)
	}

	// TODO: Move the shader creation into the material initialization so that the shader is always in use
	//       when the textures are loaded.q
	subjectProgram.Link()
	subjectProgram.Use() // <- Important to use before loading the material.
	subjectMaterial := graphics.NewMaterial(subjectProgram.Handle(), fmt.Sprintf("%s/materials/%s", *assetPath, *materialName))

	subjectMaterial.Load()

	graphics.Bind(subject)
	if err := subjectProgram.Validate(); err != nil {
		panic(err)
	}

	model := mgl32.Ident4()

	var cameraPos = mgl32.Vec3{0.0, 2.0, 5.0}
	projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0)
	viewMatrix := mgl32.LookAtV(cameraPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	subjectProgram.SetUniformMatrix4fv(graphics.ProjectionMatrixKey, projMatrix)
	subjectProgram.SetUniformMatrix4fv(graphics.ViewMatrixKey, viewMatrix)
	subjectProgram.SetUniformMatrix4fv(graphics.ModelMatrixKey, model)

	subjectProgram.SetUniform3fv(graphics.ViewPosKey, cameraPos)

	window.SetCursorPosCallback(mousePosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	lightColor := mgl32.Vec3{150, 150, 150}
	subjectProgram.SetUniform3fv(graphics.LightColorKey, lightColor)

	lightDistance := float32(10)
	rotLight := mgl32.Vec3{float32(math.Cos(lightAngle)), 0, float32(math.Sin(lightAngle))}
	rotLight = rotLight.Mul(lightDistance)
	subjectProgram.SetUniform3fv(graphics.LightPosKey, rotLight)

	font, err := graphics.NewFont("assets/fonts/ascii.png", 16, 16, float32(12), windowWidth, windowHeight)
	if err != nil {
		panic(err)
	}
	for !window.ShouldClose() {
		processInput(window)

		// Update
		rotModel := model.Mul4(mgl32.HomogRotate3DY(float32(modelAngle)))
		rotViewMatrix := viewMatrix.Mul4(mgl32.HomogRotate3DY(float32(cameraAngle)))
		rotLight := mgl32.Vec3{float32(math.Cos(lightAngle)), 0, float32(math.Sin(lightAngle))}
		rotLight = rotLight.Mul(lightDistance)

		// Render
		renderer.Clear(rotViewMatrix)
		subjectProgram.Use()
		subjectProgram.SetUniformMatrix4fv(graphics.ProjectionMatrixKey, projMatrix)
		subjectProgram.SetUniformMatrix4fv(graphics.ViewMatrixKey, rotViewMatrix)
		subjectProgram.SetUniformMatrix4fv(graphics.ModelMatrixKey, rotModel)
		subjectProgram.SetUniform3fv(graphics.ViewPosKey, cameraPos)
		subjectProgram.SetUniform3fv(graphics.LightPosKey, rotLight)
		subjectProgram.SetUniform3fv(graphics.Color, subjectColor)
		graphics.Use(subjectProgram.Handle(), subject)
		subjectMaterial.Use()
		graphics.Draw(subject)

		if showInfo {
			font.Activate()
			font.Color(mgl32.Vec3{0.5, 0.8, 0.2})
			font.RenderText(fmt.Sprintf("OpenGL Version: %s", renderer.Version), 25.0, 41.0, 1.0)
			font.RenderText(fmt.Sprintf("GLSL Version: %s", renderer.ShadingLanguageVersion), 25.0, 25.0, 1.0)

			font.Color(mgl32.Vec3{0.8, 0.7, 0.2})
			font.RenderText(fmt.Sprintf("Model Rotation: %.2f", modelAngle), 25.0, windowHeight-25, 0.5)
			font.RenderText(fmt.Sprintf("Camera Rotation: %.2f", cameraAngle), 25.0, windowHeight-41, 0.5)
			font.RenderText(fmt.Sprintf("Light Position: %v", rotLight), 25.0, windowHeight-57, 0.5)
			font.Deactivate()
		}

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	subjectProgram.Destroy()
}

var currentX float64
var previousX float64
var rotateCube bool
var rotateCamera bool
var rotateLight bool
var showInfo bool

// process all input: query GLFW whether relevant keys are pressed/released this frame and react accordingly
// ---------------------------------------------------------------------------------------------------------
func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyI) == glfw.Press {
		showInfo = !showInfo
	}
}

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

//// glfw: whenever the window size changed (by OS or user resize) this callback function executes
//// ---------------------------------------------------------------------------------------------
//func framebufferSizeCallback(w *glfw.Window, width int, height int) {
//	// make sure the viewport matches the new window dimensions; note that width and
//	// height will be significantly larger than specified on retina displays.
//	gl.Viewport(0, 0, int32(width), int32(height))
//}
