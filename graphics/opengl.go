package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"go-pbr/mesh"
)

type Backend struct {
	Version                string
	ShadingLanguageVersion string
}

func (b *Backend) Init() error {
	if err := gl.Init(); err != nil {
		return err
	}
	b.Version = gl.GoStr(gl.GetString(gl.VERSION))
	b.ShadingLanguageVersion = gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	return nil
}

func (b *Backend) Bind(uuid uint64, m mesh.Mesh) error {
	return nil
}

func (b *Backend) Clear() {
	// Configure global settings
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.CULL_FACE) //| gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)
}

//// Configure global settings
//gl.Enable(gl.TEXTURE_2D)
//if !r.RenderBack {
//	gl.Enable(gl.CULL_FACE)
//	gl.CullFace(gl.BACK)
//}
//gl.Enable(gl.DEPTH_TEST)
//gl.FrontFace(gl.CCW)
//gl.DepthRange(0, 1)
//gl.DepthFunc(gl.LESS)
//gl.ClearColor(0.0, 0.0, 0.0, 1.0)
//
//// move above
//
//gl.GenVertexArrays(1, &(m.vao))
//gl.BindVertexArray(m.vao)
//
//gl.GenBuffers(1, &(m.vbo))
//gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
//gl.BufferData(gl.ARRAY_BUFFER, len(m.data)*4, gl.Ptr(m.data), gl.STATIC_DRAW)
