package graphics

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics/opengl"
	"go-pbr/obj"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {
	CubeMapTex    uint32
	CubeMapVao    uint32
	cubemapShader opengl.Program
	WindowWidth   float32
	WindowHeight  float32
	cubemapMesh   *obj.Obj
}

func (r *Renderer) Init() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	slversion := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Println("OpenGL SL version", slversion)

	// Configure global settings
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.CULL_FACE) //| gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (r *Renderer) SetCubemap(path string) {
	gl.GenTextures(1, &r.CubeMapTex)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.CubeMapTex)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	files := map[uint32]string{
		gl.TEXTURE_CUBE_MAP_POSITIVE_X: "px.png",
		gl.TEXTURE_CUBE_MAP_NEGATIVE_X: "nx.png",
		gl.TEXTURE_CUBE_MAP_POSITIVE_Y: "py.png",
		gl.TEXTURE_CUBE_MAP_NEGATIVE_Y: "ny.png",
		gl.TEXTURE_CUBE_MAP_POSITIVE_Z: "pz.png",
		gl.TEXTURE_CUBE_MAP_NEGATIVE_Z: "nz.png",
	}

	for target, filename := range files {
		fullpath := path + "/" + filename

		imgFile, err := os.Open(fullpath)
		if err != nil {
			panic(fmt.Errorf("texture %q not found on disk: %v", fullpath, err))
		}
		img, _, err := image.Decode(imgFile)
		if err != nil {
			panic(err)
		}

		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			panic(fmt.Errorf("unsupported stride"))
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		gl.TexImage2D(
			target,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix))
	}

	r.cubemapShader = opengl.Program{}
	if f, err := os.Open("assets/shaders/cubemap/shader.vert"); err == nil {
		defer f.Close()
		r.cubemapShader.CompileShader(f, opengl.VertexShader)
	}

	if f, err := os.Open("assets/shaders/cubemap/shader.frag"); err == nil {
		defer f.Close()
		r.cubemapShader.CompileShader(f, opengl.FragmentShader)
	}

	r.cubemapShader.Link()

	if f, err := os.Open(path + "/geometry.obj"); err == nil {
		defer f.Close()
		r.cubemapMesh = obj.Load(f)
	} else {
		panic(err)
	}

	r.cubemapMesh.Bind()

	if err := r.cubemapShader.Validate(); err != nil {
		panic(err)
	}
}

func (r *Renderer) Clear(view mgl32.Mat4) {
	gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	if r.CubeMapTex <= 0 {
		return
	}

	ProjMatrix := mgl32.Perspective(mgl32.DegToRad(50.0), 1, 0.1, 200.0)
	gl.DepthMask(false)
	r.cubemapShader.Use()
	r.cubemapShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, ProjMatrix)
	r.cubemapShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, view)
	gl.ActiveTexture(gl.TEXTURE_CUBE_MAP)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.CubeMapTex)
	r.cubemapMesh.Use(r.cubemapShader.Handle())
	r.cubemapMesh.Draw()
	gl.DepthMask(true)
}

func (r *Renderer) Destroy() {

}
