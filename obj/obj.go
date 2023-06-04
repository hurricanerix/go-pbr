package obj

import (
	"bufio"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics/opengl"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// TODO: Abstract all the GL calls out of the Obj class.

type Obj struct {
	Vao        uint32
	Vbo        uint32
	Vertices   []float32
	Normals    []float32
	UVs        []float32
	Tangents   []float32
	Bitangents []float32
	Faces      []uint32
}

func Load(data io.Reader) *Obj {
	o := &Obj{
		Vertices: make([]float32, 0),
		Normals:  make([]float32, 0),
		UVs:      make([]float32, 0),
		Faces:    make([]uint32, 0),
	}

	lineTypeRx := regexp.MustCompile(`^([a-zA-Z#]+)\s*`)

	scanner := bufio.NewScanner(data)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		lineType := lineTypeRx.FindAllString(line, 1)
		switch lineType[0] {
		case "v ":
			o.Vertices = append(o.Vertices, parseVec3(line)...)
		case "vn ":
			o.Normals = append(o.Normals, parseVec3(line)...)
		case "vt ":
			o.UVs = append(o.UVs, parseVec2(line)...)
		case "f ":
			o.Faces = append(o.Faces, parseFace(line)...)
		default:

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return o
}

func parseVec3(s string) []float32 {
	result := strings.Split(s, " ")
	var c0, c1, c2 float64
	c0, _ = strconv.ParseFloat(result[1], 32)
	c1, _ = strconv.ParseFloat(result[2], 32)
	c2, _ = strconv.ParseFloat(result[3], 32)
	return []float32{float32(c0), float32(c1), float32(c2)}
}

func parseVec2(s string) []float32 {
	result := strings.Split(s, " ")
	var c0, c1 float64
	c0, _ = strconv.ParseFloat(result[1], 32)
	c1, _ = strconv.ParseFloat(result[2], 32)
	return []float32{float32(c0), float32(c1)}
}

func parseFace(s string) []uint32 {
	vertices := strings.Split(s, " ")
	results := make([]uint32, 0, 9)

	for _, v := range (vertices)[1:] {
		values := strings.Split(v, "/")
		var t, u, n uint64
		t, _ = strconv.ParseUint(values[0], 10, 32)
		u, _ = strconv.ParseUint(values[1], 10, 32)
		n, _ = strconv.ParseUint(values[2], 10, 32)
		results = append(results, []uint32{uint32(t), uint32(u), uint32(n)}...)
	}

	return results
}

func (o *Obj) GenNormalRequirements() {
	o.Tangents = make([]float32, len(o.Normals))
	o.Bitangents = make([]float32, len(o.Normals))

	for i := 0; i < len(o.Faces); i += 9 {
		v1 := mgl32.Vec3{
			o.Vertices[(o.Faces[i+0]-1)*3],
			o.Vertices[(o.Faces[i+0]-1)*3+1],
			o.Vertices[(o.Faces[i+0]-1)*3+2],
		}

		u1 := mgl32.Vec2{
			o.UVs[(o.Faces[i+1]-1)*2],
			o.UVs[(o.Faces[i+1]-1)*2+1],
		}

		v2 := mgl32.Vec3{
			o.Vertices[(o.Faces[i+3]-1)*3],
			o.Vertices[(o.Faces[i+3]-1)*3+1],
			o.Vertices[(o.Faces[i+3]-1)*3+2],
		}

		u2 := mgl32.Vec2{
			o.UVs[(o.Faces[i+4]-1)*2],
			o.UVs[(o.Faces[i+4]-1)*2+1],
		}

		v3 := mgl32.Vec3{
			o.Vertices[(o.Faces[i+6]-1)*3],
			o.Vertices[(o.Faces[i+6]-1)*3+1],
			o.Vertices[(o.Faces[i+6]-1)*3+2],
		}

		u3 := mgl32.Vec2{
			o.UVs[(o.Faces[i+7]-1)*2],
			o.UVs[(o.Faces[i+7]-1)*2+1],
		}

		tangent1, bitangent1 := calculateTB(v1, v2, v3, u1, u2, u3)

		o.Tangents[(o.Faces[i+2]-1)*3] = tangent1.X()
		o.Tangents[(o.Faces[i+2]-1)*3+1] = tangent1.Y()
		o.Tangents[(o.Faces[i+2]-1)*3+2] = tangent1.Z()

		o.Bitangents[(o.Faces[i+2]-1)*3] = bitangent1.X()
		o.Bitangents[(o.Faces[i+2]-1)*3+1] = bitangent1.Y()
		o.Bitangents[(o.Faces[i+2]-1)*3+2] = bitangent1.Z()

		//o.Tangents = append(o.Tangents, []float32{
		//	tangent1[0], tangent1[1], tangent1[2],
		//	tangent1[0], tangent1[1], tangent1[2],
		//	tangent1[0], tangent1[1], tangent1[2],
		//}...)

		//o.Bitangents = append(o.Bitangents, []float32{
		//	bitangent1[0], bitangent1[1], bitangent1[2],
		//	bitangent1[0], bitangent1[1], bitangent1[2],
		//	bitangent1[0], bitangent1[1], bitangent1[2],
		//}...)
	}
}

func calculateTB(v1, v2, v3 mgl32.Vec3, u1, u2, u3 mgl32.Vec2) (mgl32.Vec3, mgl32.Vec3) {
	edge1 := v2.Sub(v1)
	edge2 := v3.Sub(v1)
	deltaUV1 := u2.Sub(u1)
	deltaUV2 := u3.Sub(u1)

	f := 1.0 / (deltaUV1.X()*deltaUV2.Y() - deltaUV2.X()*deltaUV1.Y())

	tangent := mgl32.Vec3{
		f * (deltaUV2.Y()*edge1.X() - deltaUV1.Y()*edge2.X()),
		f * (deltaUV2.Y()*edge1.Y() - deltaUV1.Y()*edge2.Y()),
		f * (deltaUV2.Y()*edge1.Z() - deltaUV1.Y()*edge2.Z()),
	}.Normalize()

	bitangent := mgl32.Vec3{
		f * (-deltaUV2.X()*edge1.X() + deltaUV1.X()*edge2.X()),
		f * (-deltaUV2.X()*edge1.Y() + deltaUV1.X()*edge2.Y()),
		f * (-deltaUV2.X()*edge1.Z() + deltaUV1.X()*edge2.Z()),
	}.Normalize()

	return tangent, bitangent
}

func (o *Obj) Bind() {
	// Configure the vertex data
	gl.GenVertexArrays(1, &(o.Vao))
	gl.BindVertexArray(o.Vao)

	gl.GenBuffers(1, &(o.Vbo))
	gl.BindBuffer(gl.ARRAY_BUFFER, o.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(o.BufferData())*4, gl.Ptr(o.BufferData()), gl.STATIC_DRAW)
}

func (o Obj) Use(prog uint32) {
	gl.BindVertexArray(o.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, o.Vbo)

	stride := int32(8 * 4)
	if len(o.Tangents) > 0 {
		stride = 14 * 4
	}

	vertAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(opengl.VertexPosKey+"\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, stride, 0)

	uvAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(opengl.VertexUVKey+"\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointerWithOffset(uvAttrib, 2, gl.FLOAT, false, stride, 3*4)

	normalAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(opengl.VertexNormalKey+"\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointerWithOffset(normalAttrib, 3, gl.FLOAT, false, stride, 5*4)

	if len(o.Tangents) > 0 {
		tangentAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(opengl.VertexTangentKey+"\x00")))
		gl.EnableVertexAttribArray(tangentAttrib)
		gl.VertexAttribPointerWithOffset(tangentAttrib, 3, gl.FLOAT, false, stride, 8*4)

		bitangentAttrib := uint32(gl.GetAttribLocation(prog, gl.Str(opengl.VertexBitangentKey+"\x00")))
		gl.EnableVertexAttribArray(bitangentAttrib)
		gl.VertexAttribPointerWithOffset(bitangentAttrib, 3, gl.FLOAT, false, stride, 11*4)
	}
}

func (o Obj) BufferData() []float32 {
	vertexCount := len(o.Faces) / 3
	result := make([]float32, 0, vertexCount*14)
	for i := 0; i < len(o.Faces); i += 3 {
		result = append(result, o.Vertices[(o.Faces[i+0]-1)*3])
		result = append(result, o.Vertices[(o.Faces[i+0]-1)*3+1])
		result = append(result, o.Vertices[(o.Faces[i+0]-1)*3+2])
		result = append(result, o.UVs[(o.Faces[i+1]-1)*2])
		result = append(result, o.UVs[(o.Faces[i+1]-1)*2+1])
		result = append(result, o.Normals[(o.Faces[i+2]-1)*3])
		result = append(result, o.Normals[(o.Faces[i+2]-1)*3+1])
		result = append(result, o.Normals[(o.Faces[i+2]-1)*3+2])

		if len(o.Tangents) > 0 {
			result = append(result, o.Tangents[(o.Faces[i+2]-1)*3])
			result = append(result, o.Tangents[(o.Faces[i+2]-1)*3+1])
			result = append(result, o.Tangents[(o.Faces[i+2]-1)*3+2])
		}

		if len(o.Bitangents) > 0 {
			result = append(result, o.Bitangents[(o.Faces[i+2]-1)*3])
			result = append(result, o.Bitangents[(o.Faces[i+2]-1)*3+1])
			result = append(result, o.Bitangents[(o.Faces[i+2]-1)*3+2])
		}
	}
	return result
}

func (o Obj) indices() int32 {
	return int32(len(o.Faces) / 3)
}

func (o Obj) Draw() {
	gl.DrawArrays(gl.TRIANGLES, 0, o.indices())
}

func (o Obj) String() string {
	s := &strings.Builder{}
	s.WriteString(fmt.Sprintf("Vertex Count: %d\n", len(o.Vertices)))
	s.WriteString(fmt.Sprintf("Normal Count: %d\n", len(o.Normals)))
	s.WriteString(fmt.Sprintf("UV Count: %d\n", len(o.UVs)))
	s.WriteString(fmt.Sprintf("Face Count: %d\n", len(o.Faces)))
	s.WriteString(fmt.Sprintf("Tangent Count: %d\n", len(o.Tangents)))
	s.WriteString(fmt.Sprintf("Bitangent Count: %d\n", len(o.Bitangents)))
	s.WriteString(fmt.Sprintf("BufferData Count: %d\n", len(o.BufferData())))
	s.WriteString(fmt.Sprintf("Indicies Count: %d\n", o.indices()))
	s.WriteString(fmt.Sprintf("BufferData:\n"))

	if len(o.Tangents) < 0 {
		bd := o.BufferData()
		for i := 0; i < len(bd); i = i + 24 {
			s.WriteString("{\n")
			s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f)\n", bd[i+0], bd[i+1], bd[i+2], bd[i+3], bd[i+4], bd[i+5], bd[i+6], bd[i+7]))
			s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f)\n", bd[i+8], bd[i+9], bd[i+10], bd[i+11], bd[i+12], bd[i+13], bd[i+14], bd[i+15]))
			s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f)\n", bd[i+16], bd[i+17], bd[i+18], bd[i+19], bd[i+20], bd[i+21], bd[i+22], bd[i+23]))
			s.WriteString("}\n")
		}
	} else {
		bd := o.BufferData()
		for i := 0; i < len(bd); i = i + 42 {
			s.WriteString("{\n")
			s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f) tangent(%.1f, %.1f, %.1f) bitangent(%.1f, %.1f, %.1f)\n", bd[i+0], bd[i+1], bd[i+2], bd[i+3], bd[i+4], bd[i+5], bd[i+6], bd[i+7], bd[i+8], bd[i+9], bd[i+10], bd[i+11], bd[i+12], bd[i+13]))
			s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f) tangent(%.1f, %.1f, %.1f) bitangent(%.1f, %.1f, %.1f)\n", bd[i+14], bd[i+15], bd[i+16], bd[i+17], bd[i+18], bd[i+19], bd[i+20], bd[i+21], bd[i+22], bd[i+23], bd[i+24], bd[i+25], bd[i+26], bd[i+27]))
			s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f) tangent(%.1f, %.1f, %.1f) bitangent(%.1f, %.1f, %.1f)\n", bd[i+28], bd[i+29], bd[i+30], bd[i+31], bd[i+32], bd[i+33], bd[i+34], bd[i+35], bd[i+36], bd[i+37], bd[i+38], bd[i+39], bd[i+40], bd[i+41]))
			s.WriteString("}\n")
		}
	}

	return s.String()
}
