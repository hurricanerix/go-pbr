package pbr

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"

	"go-pbr/graphics"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const windowWidth = 960
const windowHeight = 540

var modelAngleX = 0.0
var cameraAngleX = 0.0

type App struct {
	window       *glfw.Window
	renderer     graphics.Renderer
	currentTime  float64
	previousTime float64
}

func (a *App) Init() error {
	if err := glfw.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var err error
	a.window, err = glfw.CreateWindow(windowWidth, windowHeight, "PBR", nil, nil)
	if err != nil {
		panic(err)
	}
	a.window.MakeContextCurrent()

	a.renderer = graphics.Renderer{
		WindowWidth:  windowWidth,
		WindowHeight: windowHeight,
	}
	a.renderer.Init()

	//var cubeMesh *mesh.Obj
	//if f, err := os.Open(*meshPath); err == nil {
	//	defer f.Close()
	//	cubeMesh = mesh.Load(f)
	//	cubeMesh.GenNormalRequirements()
	//} else {
	//	panic(err)
	//}

	//a.renderer.SetCubemap(*skymapPath)

	//phongShader := opengl.Program{}
	//if f, err := os.Open("assets/shaders/parallax.vert"); err == nil {
	//	//if f, err := os.Open("assets/shaders/phong.vert"); err == nil {
	//	defer f.Close()
	//	phongShader.CompileShader(f, opengl.VertexShader)
	//}
	//
	//if f, err := os.Open("assets/shaders/parallax.frag"); err == nil {
	//	//if f, err := os.Open("assets/shaders/phong.frag"); err == nil {
	//	defer f.Close()
	//	phongShader.CompileShader(f, opengl.FragmentShader)
	//}
	//
	//// TODO: Move the shader creation into the material initialization so that the shader is always in use
	////       when the textures are loaded.q
	//phongShader.Link()
	//phongShader.Use() // <- Important to use before loading the material.
	//brickMat := graphics.NewMaterial(phongShader.Handle(), *materialPath)
	////brickMat := graphics.NewMaterial(phongShader.Handle(), matDir+"/lgl_brickwall")
	////brickMat := graphics.NewMaterial(phongShader.Handle(), matDir+"/stone_wall")
	//
	//brickMat.Load()
	//
	//cubeMesh.Bind()
	//if err := phongShader.Validate(); err != nil {
	//	panic(err)
	//}
	//
	//model := mgl32.Ident4()
	//projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 20.0)
	//
	//cameraX := float32(math.Cos(cameraAngleX))
	//cameraY := float32(math.Sin(cameraAngleX))
	//cameraDir := mgl32.Vec3{cameraX, 0, cameraY}.Normalize()
	//var cameraPos = cameraDir.Mul(7)
	//viewMatrix := mgl32.LookAtV(cameraPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	//
	//phongShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, projMatrix)
	//phongShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, viewMatrix)
	//phongShader.SetUniformMatrix4fv(opengl.ModelMatrixKey, model)
	//
	//phongShader.SetUniform3fv(opengl.ViewPosKey, cameraPos)

	a.window.SetCursorPosCallback(mousePosCallback)
	a.window.SetMouseButtonCallback(mouseButtonCallback)

	//material := graphics.NewLitMaterial(*materialPath)
	//mesh := mesh.NewMesh(*meshPath, material)
	//skymap := graphics.NewSkymap(*skymapPath)
	//camera := graphics.NewCamera()
	//re?tur nil

	return nil
}

var currentX float64
var previousX float64
var rotateCube bool

var rotateCamera bool

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
		modelAngleX += dirX * speed / 2
	}

	if rotateCamera {
		speed := 0.05
		dirX := 0.0
		if previousX < currentX {
			dirX = 1.0
		} else if previousX > currentX {
			dirX = -1.0
		}
		cameraAngleX += dirX * speed / 2
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
		rotateCamera = true
	}

	if button == glfw.MouseButtonRight && action == glfw.Release {
		rotateCamera = false
	}
}

func (a App) ShouldClose() bool {
	return a.window.ShouldClose()
}

func (a *App) DeltaTime() float32 {
	a.currentTime = glfw.GetTime()
	dt := float32(a.currentTime - a.previousTime)
	a.previousTime = a.currentTime
	return dt
}

func (a *App) PostLoop() error {
	a.window.SwapBuffers()
	glfw.PollEvents()
	return nil
}

func (a *App) Free() error {
	glfw.Terminate()
	return nil
}
