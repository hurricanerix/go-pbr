// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2023 Richard Hawkins
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file borrows code from png/reader.go in the image/png package.

package obj

import (
	"bufio"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/mesh"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// FormatError reports that the input is not a valid OBJ.
type FormatError string

func (e FormatError) Error() string { return "obj: invalid format: " + string(e) }

// Decode the provided reader using the provided options and return the mesh.
func Decode(r io.Reader, options mesh.DecodeOptions) (mesh.Mesh, error) {
	d := &decoder{
		r:       r,
		options: options,
	}

	if options.SkipHeaderCheck == false {
		if err := d.checkHeader(); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}
	}

	if err := d.decode(); err != nil {
		return Obj{}, err
	}

	return d.mesh, nil
}

const objHeader = "# OBJ"
const blenderObjHeader = "# Blender"

type decoder struct {
	r       io.Reader
	options mesh.DecodeOptions
	mesh    mesh.Mesh

	vertices   []float32
	normals    []float32
	uvs        []float32
	faces      []uint32
	tangents   []float32
	bitangents []float32
	tmp        [3 * 256]byte
}

func (d *decoder) getConfig() mesh.Config {
	c := mesh.Config{
		Stride: 8 * 4,
		Fields: map[mesh.VertexAttribute]mesh.VertexAttributeData{
			mesh.Position: {Size: 3, Offset: 0},
			mesh.TexCoord: {Size: 2, Offset: 3 * 4},
			mesh.Normal:   {Size: 3, Offset: 5 * 4},
		},
	}

	if len(d.tangents) > 0 && len(d.bitangents) > 0 {
		c.Stride += 6 * 4
		c.Fields[mesh.Tangent] = mesh.VertexAttributeData{Size: 3, Offset: 8 * 4}
		c.Fields[mesh.Bitangent] = mesh.VertexAttributeData{Size: 3, Offset: 11 * 4}
	}

	return c
}

func (d *decoder) checkHeader() error {
	if !d.validBlenderHeader() {
		return FormatError("could not detect a Blender OBJ file")
	}
	return nil
}

func (d *decoder) validBlenderHeader() bool {
	_, err := io.ReadFull(d.r, d.tmp[:len(blenderObjHeader)])
	if rs, ok := d.r.(io.Seeker); ok {
		rs.Seek(0, 0)
	}
	if err != nil {
		return false
	}
	if string(d.tmp[:len(blenderObjHeader)]) != blenderObjHeader {
		return false
	}
	return true
}

func (d *decoder) decode() error {
	lineTypeRx := regexp.MustCompile(`^([a-zA-Z#]+)\s*`)

	scanner := bufio.NewScanner(d.r)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		lineType := lineTypeRx.FindAllString(line, 1)
		switch lineType[0] {
		case "v ":
			d.vertices = append(d.vertices, parseVec3(line)...)
		case "vn ":
			d.normals = append(d.normals, parseVec3(line)...)
		case "vt ":
			d.uvs = append(d.uvs, parseVec2(line)...)
		case "f ":
			d.faces = append(d.faces, parseFace(line)...)
		default:

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if d.options.CalculateTangentsIfMissing {
		d.calcTangents()
	}

	obj := Obj{
		data:   d.bufferData(),
		config: d.getConfig(),
	}

	d.mesh = obj

	return nil
}

func (d *decoder) calcTangents() {
	d.tangents = make([]float32, len(d.normals))
	d.bitangents = make([]float32, len(d.normals))

	for i := 0; i < len(d.faces); i += 9 {
		v1 := mgl32.Vec3{
			d.vertices[(d.faces[i+0]-1)*3],
			d.vertices[(d.faces[i+0]-1)*3+1],
			d.vertices[(d.faces[i+0]-1)*3+2],
		}

		u1 := mgl32.Vec2{
			d.uvs[(d.faces[i+1]-1)*2],
			d.uvs[(d.faces[i+1]-1)*2+1],
		}

		v2 := mgl32.Vec3{
			d.vertices[(d.faces[i+3]-1)*3],
			d.vertices[(d.faces[i+3]-1)*3+1],
			d.vertices[(d.faces[i+3]-1)*3+2],
		}

		u2 := mgl32.Vec2{
			d.uvs[(d.faces[i+4]-1)*2],
			d.uvs[(d.faces[i+4]-1)*2+1],
		}

		v3 := mgl32.Vec3{
			d.vertices[(d.faces[i+6]-1)*3],
			d.vertices[(d.faces[i+6]-1)*3+1],
			d.vertices[(d.faces[i+6]-1)*3+2],
		}

		u3 := mgl32.Vec2{
			d.uvs[(d.faces[i+7]-1)*2],
			d.uvs[(d.faces[i+7]-1)*2+1],
		}

		tangent1, bitangent1 := calculateTB(v1, v2, v3, u1, u2, u3)

		d.tangents[(d.faces[i+2]-1)*3] = tangent1.X()
		d.tangents[(d.faces[i+2]-1)*3+1] = tangent1.Y()
		d.tangents[(d.faces[i+2]-1)*3+2] = tangent1.Z()

		d.bitangents[(d.faces[i+2]-1)*3] = bitangent1.X()
		d.bitangents[(d.faces[i+2]-1)*3+1] = bitangent1.Y()
		d.bitangents[(d.faces[i+2]-1)*3+2] = bitangent1.Z()
	}
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

func (d decoder) bufferData() []float32 {
	vertexCount := len(d.faces) / 3
	result := make([]float32, 0, vertexCount*14)
	for i := 0; i < len(d.faces); i += 3 {
		result = append(result, d.vertices[(d.faces[i+0]-1)*3])
		result = append(result, d.vertices[(d.faces[i+0]-1)*3+1])
		result = append(result, d.vertices[(d.faces[i+0]-1)*3+2])
		result = append(result, d.uvs[(d.faces[i+1]-1)*2])
		result = append(result, d.uvs[(d.faces[i+1]-1)*2+1])
		result = append(result, d.normals[(d.faces[i+2]-1)*3])
		result = append(result, d.normals[(d.faces[i+2]-1)*3+1])
		result = append(result, d.normals[(d.faces[i+2]-1)*3+2])

		if len(d.tangents) > 0 {
			result = append(result, d.tangents[(d.faces[i+2]-1)*3])
			result = append(result, d.tangents[(d.faces[i+2]-1)*3+1])
			result = append(result, d.tangents[(d.faces[i+2]-1)*3+2])
		}

		if len(d.bitangents) > 0 {
			result = append(result, d.bitangents[(d.faces[i+2]-1)*3])
			result = append(result, d.bitangents[(d.faces[i+2]-1)*3+1])
			result = append(result, d.bitangents[(d.faces[i+2]-1)*3+2])
		}
	}
	return result
}

func init() {
	mesh.RegisterFormat("obj", objHeader, Decode)
	mesh.RegisterFormat("obj", blenderObjHeader, Decode)
}
