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
	meshName := flag.String("mesh", "xyz-cube", "")
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

	var light *obj.Obj
	if f, err := os.Open(fmt.Sprintf("%s/meshes/light.obj", *assetPath)); err == nil {
		defer f.Close()
		light = obj.Load(f)
		//light.GenNormalRequirements()
	} else {
		panic(err)
	}
	lightProgram := graphics.Program{}
	if f, err := os.Open(fmt.Sprintf("%s/shaders/unlit/shader.vert", *assetPath)); err == nil {
		defer f.Close()
		lightProgram.CompileShader(f, graphics.VertexShader)
	}
	if f, err := os.Open(fmt.Sprintf("%s/shaders/unlit/shader.frag", *assetPath)); err == nil {
		defer f.Close()
		lightProgram.CompileShader(f, graphics.FragmentShader)
	}

	lightProgram.Link()
	lightProgram.Use() // <- Important to use before loading the material.
	lightMaterial := graphics.NewMaterial(subjectProgram.Handle(), fmt.Sprintf("%s/materials/debug", *assetPath))
	lightMaterial.Load()
	graphics.Bind(light)
	if err := lightProgram.Validate(); err != nil {
		panic(err)
	}

	lightColor := mgl32.Vec3{9.922, 9.843, 8.275}
	lightDistance := float32(2)
	rotLight := mgl32.Vec3{float32(math.Cos(lightAngle)), 0, float32(math.Sin(lightAngle))}
	rotLight = rotLight.Mul(lightDistance)
	subjectProgram.SetUniform3fv(graphics.LightPosKey, rotLight)

	lightModel := mgl32.Ident4()
	lightModel = mgl32.Translate3D(rotLight.X(), rotLight.Y(), rotLight.Z())
	lightProgram.Use()
	lightProgram.SetUniformMatrix4fv(graphics.ProjectionMatrixKey, projMatrix)
	lightProgram.SetUniformMatrix4fv(graphics.ViewMatrixKey, viewMatrix)
	lightProgram.SetUniformMatrix4fv(graphics.ModelMatrixKey, lightModel)
	lightProgram.SetUniform3fv(graphics.ColorKey, lightColor)
	lightProgram.SetUniform1i("noMap", 1)
	subjectProgram.Use()
	subjectProgram.SetUniformMatrix4fv(graphics.ProjectionMatrixKey, projMatrix)
	subjectProgram.SetUniformMatrix4fv(graphics.ViewMatrixKey, viewMatrix)
	subjectProgram.SetUniformMatrix4fv(graphics.ModelMatrixKey, model)
	subjectProgram.SetUniform3fv(graphics.ViewPosKey, cameraPos)
	subjectProgram.SetUniform3fv(graphics.LightColorKey, lightColor)

	font, err := graphics.NewFont("assets/fonts/ascii.png", 16, 16, float32(12), windowWidth, windowHeight)
	if err != nil {
		panic(err)
	}

	window.SetCursorPosCallback(mousePosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	fmt.Println(light)

	currentTimer := 0
	timerResolution := 1000
	timers := map[string][]float64{
		"start":           make([]float64, timerResolution),
		"endProcessInput": make([]float64, timerResolution),
		"endUpdate":       make([]float64, timerResolution),
		"endRender":       make([]float64, timerResolution),
		"endShowLight":    make([]float64, timerResolution),
		"endShowInfo":     make([]float64, timerResolution),
		"endSwap":         make([]float64, timerResolution),
		"endPoll":         make([]float64, timerResolution),
		"end":             make([]float64, timerResolution),
	}

	for !window.ShouldClose() {
		timers["start"][currentTimer] = glfw.GetTime()
		processInput(window)
		timers["endProcessInput"][currentTimer] = glfw.GetTime()

		// Update
		rotModel := model.Mul4(mgl32.HomogRotate3DY(float32(modelAngle)))
		rotViewMatrix := viewMatrix.Mul4(mgl32.HomogRotate3DY(float32(cameraAngle)))
		rotLight := mgl32.Vec3{float32(math.Cos(lightAngle)), 0, float32(math.Sin(lightAngle))}
		rotLight = rotLight.Mul(lightDistance)
		lightModel = mgl32.Translate3D(rotLight.X(), rotLight.Y(), rotLight.Z())
		timers["endUpdate"][currentTimer] = glfw.GetTime()

		// Render
		renderer.Clear(rotViewMatrix)

		subjectProgram.Use()
		subjectProgram.SetUniformMatrix4fv(graphics.ProjectionMatrixKey, projMatrix)
		subjectProgram.SetUniformMatrix4fv(graphics.ViewMatrixKey, rotViewMatrix)
		subjectProgram.SetUniformMatrix4fv(graphics.ModelMatrixKey, rotModel)
		subjectProgram.SetUniform3fv(graphics.ViewPosKey, cameraPos)
		subjectProgram.SetUniform3fv(graphics.LightPosKey, rotLight)
		subjectProgram.SetUniform3fv(graphics.ColorKey, subjectColor)
		graphics.Use(subjectProgram.Handle(), subject)
		subjectMaterial.Use()
		graphics.Draw(subject)
		timers["endRender"][currentTimer] = glfw.GetTime()

		if showLight {
			lightProgram.Use()
			lightProgram.SetUniformMatrix4fv(graphics.ProjectionMatrixKey, projMatrix)
			lightProgram.SetUniformMatrix4fv(graphics.ViewMatrixKey, rotViewMatrix)
			lightProgram.SetUniformMatrix4fv(graphics.ModelMatrixKey, lightModel)
			lightProgram.SetUniform3fv(graphics.ColorKey, lightColor)
			graphics.Use(lightProgram.Handle(), light)
			lightProgram.Use()
			graphics.Draw(light)
		}
		timers["endShowLight"][currentTimer] = glfw.GetTime()

		if showInfo {
			font.Activate()
			font.Color(mgl32.Vec3{0.5, 0.8, 0.2})
			font.RenderText(fmt.Sprintf("OpenGL Version: %s", renderer.Version), 25.0, 41.0, 1.0)
			font.RenderText(fmt.Sprintf("GLSL Version: %s", renderer.ShadingLanguageVersion), 25.0, 25.0, 1.0)

			font.Color(mgl32.Vec3{0.8, 0.7, 0.2})
			font.RenderText(fmt.Sprintf("Model Rotation: %.2f", modelAngle), 25.0, windowHeight-25, 0.5)
			font.RenderText(fmt.Sprintf("Camera Rotation: %.2f", cameraAngle), 25.0, windowHeight-41, 0.5)
			font.RenderText(fmt.Sprintf("Light Position: %v", rotLight), 25.0, windowHeight-57, 0.5)
			font.RenderText(fmt.Sprintf("Mesh: %s", *meshName), 25.0, windowHeight-73, 0.5)
			font.RenderText(fmt.Sprintf("Material: %s", *materialName), 25.0, windowHeight-89, 0.5)
			font.RenderText(fmt.Sprintf("Shader: %s", *shaderName), 25.0, windowHeight-105, 0.5)
			font.Deactivate()
		}
		timers["endShowInfo"][currentTimer] = glfw.GetTime()

		// Maintenance
		window.SwapBuffers()
		timers["endSwap"][currentTimer] = glfw.GetTime()
		glfw.PollEvents()
		timers["endPoll"][currentTimer] = glfw.GetTime()

		timers["end"][currentTimer] = glfw.GetTime()
		currentTimer += 1
		if currentTimer >= timerResolution {
			currentTimer = 0
			//window.SetShouldClose(true)
		}
	}

	processTimers(timers)

	subjectProgram.Destroy()
}

func processTimers(t map[string][]float64) {
	max := map[string]float64{
		"start":           0.0,
		"endProcessInput": 0.0,
		"endUpdate":       0.0,
		"endRender":       0.0,
		"endShowLight":    0.0,
		"endShowInfo":     0.0,
		"endSwap":         0.0,
		"endPoll":         0.0,
		"end":             0.0,
	}
	min := map[string]float64{
		"start":           999999.0,
		"endProcessInput": 999999.0,
		"endUpdate":       999999.0,
		"endRender":       999999.0,
		"endShowLight":    999999.0,
		"endShowInfo":     999999.0,
		"endSwap":         999999.0,
		"endPoll":         999999.0,
		"end":             999999.0,
	}
	avg := map[string]float64{
		"start":           0.0,
		"endProcessInput": 0.0,
		"endUpdate":       0.0,
		"endRender":       0.0,
		"endShowLight":    0.0,
		"endShowInfo":     0.0,
		"endSwap":         0.0,
		"endPoll":         0.0,
		"end":             0.0,
	}

	for _, k := range []string{"endProcessInput", "endUpdate", "endRender", "endShowLight", "endShowInfo", "endSwap", "endPoll", "end"} {
		v := t[k]
		for i := range v {
			d := t[k][i] - t["start"][i]
			if d > max[k] {
				max[k] = d
			}
			if d < min[k] {
				min[k] = d
			}
			avg[k] += d
		}
	}

	for _, k := range []string{"start", "endProcessInput", "endUpdate", "endRender", "endShowLight", "endShowInfo", "endSwap", "endPoll", "end"} {
		fmt.Printf("%s\n", k)
		fmt.Printf("\tMin: %f\n", min[k])
		fmt.Printf("\tAvg: %f\n", avg[k]/float64(len(t["start"])))
		fmt.Printf("\tMax: %f\n", max[k])
	}

}

var currentX float64
var previousX float64
var rotateCube bool
var rotateCamera bool
var rotateLight bool

// var showInfo bool
// var showLight bool
var showLight = true
var showInfo = true

// process all input: query GLFW whether relevant keys are pressed/released this frame and react accordingly
// ---------------------------------------------------------------------------------------------------------
func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyI) == glfw.Press {
		showInfo = !showInfo
	}
	if window.GetKey(glfw.KeyL) == glfw.Press {
		showLight = !showLight
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
