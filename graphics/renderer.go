package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics/opengl"
	_ "image/png"
)

type Renderer struct {
	CubeMapTex    uint32
	CubeMapVao    uint32
	cubemapShader opengl.Program
	WindowWidth   float32
	WindowHeight  float32
	//cubemapMesh   *mesh.Obj
	RenderBack bool
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
	if !r.RenderBack {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.DepthRange(0, 1)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// move above

	gl.GenVertexArrays(1, &(m.vao))
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(1, &(m.vbo))
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.data)*4, gl.Ptr(m.data), gl.STATIC_DRAW)
}

func (r *Renderer) Draw(data []float32) {

}

func (r *Renderer) Clear(view mgl32.Mat4) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	if r.CubeMapTex <= 0 {
		return
	}

	ProjMatrix := mgl32.Perspective(mgl32.DegToRad(10.0), r.WindowWidth/r.WindowHeight, 0.1, 20.0)
	gl.DepthMask(false)

	if r.RenderBack {
		gl.Enable(gl.CULL_FACE)
		gl.FrontFace(gl.CCW)
		gl.CullFace(gl.BACK)
	}

	r.cubemapShader.Use()
	r.cubemapShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, ProjMatrix)
	r.cubemapShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, view)
	gl.ActiveTexture(gl.TEXTURE_CUBE_MAP)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.CubeMapTex)
	//r.cubemapMesh.Use(r.cubemapShader.Handle())
	//r.cubemapMesh.Draw()

	gl.DepthMask(true)
	if !r.RenderBack {
		gl.Disable(gl.CULL_FACE)
	}
}

func (r *Renderer) Destroy() {

}
