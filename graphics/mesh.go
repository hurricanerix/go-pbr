package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"go-pbr/obj"
)

func Draw(o *obj.Obj) {
	gl.DrawArrays(gl.TRIANGLES, 0, o.Indices())
}

func Bind(o *obj.Obj) {
	// Configure the vertex data
	gl.GenVertexArrays(1, &(o.Vao))
	gl.BindVertexArray(o.Vao)

	gl.GenBuffers(1, &(o.Vbo))
	gl.BindBuffer(gl.ARRAY_BUFFER, o.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(o.BufferData())*4, gl.Ptr(o.BufferData()), gl.STATIC_DRAW)
}

func Use(prog uint32, o *obj.Obj) {
	gl.BindVertexArray(o.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, o.Vbo)

	stride := int32(8 * 4)
	if len(o.Tangents) > 0 {
		stride = 14 * 4
	}

	vertAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(VertexPosKey+"\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, stride, 0)

	uvAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(VertexUVKey+"\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointerWithOffset(uvAttrib, 2, gl.FLOAT, false, stride, 3*4)

	normalAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(VertexNormalKey+"\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointerWithOffset(normalAttrib, 3, gl.FLOAT, false, stride, 5*4)

	if len(o.Tangents) > 0 {
		tangentAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(VertexTangentKey+"\x00")))
		gl.EnableVertexAttribArray(tangentAttrib)
		gl.VertexAttribPointerWithOffset(tangentAttrib, 3, gl.FLOAT, false, stride, 8*4)

		bitangentAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(VertexBitangentKey+"\x00")))
		gl.EnableVertexAttribArray(bitangentAttrib)
		gl.VertexAttribPointerWithOffset(bitangentAttrib, 3, gl.FLOAT, false, stride, 11*4)
	}
}
