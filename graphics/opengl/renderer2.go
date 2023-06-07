package opengl

import (
	"go-pbr/scene"
	_ "image/png"
)

type Renderer interface {
	Render(o scene.Object) error
}

//type Renderer struct {
//	CubeMapTex    uint32
//	CubeMapVao    uint32
//	cubemapShader opengl.Program
//	WindowWidth   float32
//	WindowHeight  float32
//	//cubemapMesh   *mesh.Obj
//	RenderBack bool
//}
//

//
//func (r *Renderer) Draw(data []float32) {
//
//}
//
//func (r *Renderer) Clear(view mgl32.Mat4) {
//	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
//	gl.UseProgram(0)
//	gl.BindTexture(gl.TEXTURE_2D, 0)
//
//	if r.CubeMapTex <= 0 {
//		return
//	}
//
//	ProjMatrix := mgl32.Perspective(mgl32.DegToRad(10.0), r.WindowWidth/r.WindowHeight, 0.1, 20.0)
//	gl.DepthMask(false)
//
//	if r.RenderBack {
//		gl.Enable(gl.CULL_FACE)
//		gl.FrontFace(gl.CCW)
//		gl.CullFace(gl.BACK)
//	}
//
//	r.cubemapShader.Use()
//	r.cubemapShader.SetUniformMatrix4fv(opengl.ProjectionMatrixKey, ProjMatrix)
//	r.cubemapShader.SetUniformMatrix4fv(opengl.ViewMatrixKey, view)
//	gl.ActiveTexture(gl.TEXTURE_CUBE_MAP)
//	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.CubeMapTex)
//	//r.cubemapMesh.Use(r.cubemapShader.Handle())
//	//r.cubemapMesh.Draw()
//
//	gl.DepthMask(true)
//	if !r.RenderBack {
//		gl.Disable(gl.CULL_FACE)
//	}
//}
//
//func (r *Renderer) Destroy() {
//
//}

//func Render(o scene.Object) error {
//
//}
