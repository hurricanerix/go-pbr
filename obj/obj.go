package obj

import (
	"bufio"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// TODO: Abstract all the GL calls out of the Obj class.

type Obj struct {
	Vao uint32

	Vertices []float32
	Normals  []float32
	UVs      []float32
	Faces    []uint32
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

func (o *Obj) Bind(prog uint32) {
	// Configure the vertex data
	gl.GenVertexArrays(1, &(o.Vao))
	gl.BindVertexArray(o.Vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(o.BufferData())*4, gl.Ptr(o.BufferData()), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(prog, gl.Str("Vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 8*4, 0)

	uvAttrib := uint32(gl.GetAttribLocation(prog, gl.Str("UV\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointerWithOffset(uvAttrib, 2, gl.FLOAT, false, 8*4, 3*4)

	normalAttrib := uint32(gl.GetAttribLocation(prog, gl.Str("aNormal\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointerWithOffset(normalAttrib, 3, gl.FLOAT, false, 8*4, 5*4)

	gl.BindVertexArray(o.Vao)
}

func (o Obj) BufferData() []float32 {
	vertexCount := len(o.Faces) / 3
	result := make([]float32, 0, vertexCount*8)
	for i := 0; i < len(o.Faces); i += 3 {
		result = append(result, o.Vertices[(o.Faces[i+0]-1)*3])
		result = append(result, o.Vertices[(o.Faces[i+0]-1)*3+1])
		result = append(result, o.Vertices[(o.Faces[i+0]-1)*3+2])
		result = append(result, o.UVs[(o.Faces[i+1]-1)*2])
		result = append(result, o.UVs[(o.Faces[i+1]-1)*2+1])
		result = append(result, o.Normals[(o.Faces[i+2]-1)*3])
		result = append(result, o.Normals[(o.Faces[i+2]-1)*3+1])
		result = append(result, o.Normals[(o.Faces[i+2]-1)*3+2])
	}
	return result
}

func (o Obj) Indices() int32 {
	return int32(len(o.Faces) / 3)
}

func (o Obj) String() string {
	s := &strings.Builder{}
	s.WriteString(fmt.Sprintf("Vertex Count: %d\n", len(o.Vertices)))
	s.WriteString(fmt.Sprintf("Normal Count: %d\n", len(o.Normals)))
	s.WriteString(fmt.Sprintf("UV Count: %d\n", len(o.UVs)))
	s.WriteString(fmt.Sprintf("Face Count: %d\n", len(o.Faces)))
	s.WriteString(fmt.Sprintf("BufferData Count: %d\n", len(o.BufferData())))
	s.WriteString(fmt.Sprintf("Indicies Count: %d\n", o.Indices()))
	s.WriteString(fmt.Sprintf("BufferData:\n"))
	bd := o.BufferData()
	for i := 0; i < len(bd); i = i + 24 {
		s.WriteString("{\n")
		s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f)\n", bd[i+0], bd[i+1], bd[i+2], bd[i+3], bd[i+4], bd[i+5], bd[i+6], bd[i+7]))
		s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f)\n", bd[i+8], bd[i+9], bd[i+10], bd[i+11], bd[i+12], bd[i+13], bd[i+14], bd[i+15]))
		s.WriteString(fmt.Sprintf("\tpos(%.1f, %.1f, %.1f) uv(%.1f, %.1f) normal(%.1f, %.1f, %.1f)\n", bd[i+16], bd[i+17], bd[i+18], bd[i+19], bd[i+20], bd[i+21], bd[i+22], bd[i+23]))
		s.WriteString("}\n")
	}
	return s.String()
}
