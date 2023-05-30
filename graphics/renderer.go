package graphics

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {
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
	gl.Enable(gl.CULL_FACE) //| gl.DEPTH_TEST)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (r *Renderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)
}
func (r *Renderer) Draw(count int32) {
	gl.DrawArrays(gl.TRIANGLES, 0, count)
}

func (r *Renderer) Destroy() {

}
