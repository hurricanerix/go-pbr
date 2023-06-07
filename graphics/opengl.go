package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
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

func (b *Backend) Clear() {
	// Configure global settings
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.CULL_FACE) //| gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.1, 0.1, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)
}
