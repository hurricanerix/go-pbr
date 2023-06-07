package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/material"
	"go-pbr/mesh"
)

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}

type Renderer struct {
	ID       int
	Program  *Program
	Mesh     mesh.Mesh
	Material material.Material

	Projection mgl32.Mat4
	View       mgl32.Mat4
	Model      mgl32.Mat4
	vao        uint32
	vbo        uint32
}

func (r *Renderer) Attach() error {
	var tmp uint32
	gl.GenVertexArrays(1, &(tmp))
	r.vao = tmp
	gl.GenBuffers(1, &(tmp))
	r.vbo = tmp
	if err := r.Bind(); err != nil {
		return nil
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(r.Mesh.Data())*4, gl.Ptr(r.Mesh.Data()), gl.STATIC_DRAW)

	return nil
}

func (r *Renderer) Detatch() error {
	r.Program.Destroy() // TODO: rename Delete
	return nil
}

func (r *Renderer) Bind() error {
	gl.BindVertexArray(r.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)

	for k, v := range r.Mesh.Config().Fields {
		var attribName string
		switch k {
		case mesh.Position:
			attribName = VertexPosKey
		case mesh.TexCoord:
			attribName = VertexUVKey
		case mesh.Normal:
			attribName = VertexNormalKey
		case mesh.Tangent:
			attribName = VertexTangentKey
		case mesh.Bitangent:
			attribName = VertexBitangentKey
		}
		attrib := uint32(gl.GetAttribLocation(r.Program.Handle(), gl.Str(attribName+"\x00")))
		gl.EnableVertexAttribArray(attrib)
		gl.VertexAttribPointerWithOffset(attrib, v.Size, gl.FLOAT, false, r.Mesh.Config().Stride, v.Offset)
	}

	return nil
}

func (r *Renderer) SetTransform(t Transform) error {
	r.Model = mgl32.Ident4()
	//r.Model = r.Model.Mul4(mgl32.HomogRotate3DX(t.Rotation.X()))
	//r.Model = r.Model.Mul4(mgl32.HomogRotate3DY(t.Rotation.Y()))
	//r.Model = r.Model.Mul4(mgl32.HomogRotate3DZ(t.Rotation.Z()))
	return nil
}

func (r *Renderer) Draw() error {
	r.Program.Use() // <- Important to use before loading the material.
	if err := r.Material.Use(); err != nil {
		return err
	}

	r.Program.SetUniformMatrix4fv(ProjectionMatrixKey, r.Projection)
	r.Program.SetUniformMatrix4fv(ViewMatrixKey, r.View)
	r.Program.SetUniformMatrix4fv(ModelMatrixKey, r.Model)

	gl.DrawArrays(gl.TRIANGLES, 0, r.Mesh.Config().Indices)

	return nil
}
