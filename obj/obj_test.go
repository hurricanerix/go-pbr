package obj

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVec3(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []float32
	}{
		{name: "simple", input: "v 1.000000 1.000000 -1.000000", want: []float32{1.0, 1.0, -1.0}},
	}

	for _, tc := range tests {
		got := parseVec3(tc.input)
		assert.Equal(t, tc.want, got, tc.name)
	}
}

func TestParseVec2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []float32
	}{
		{name: "simple", input: "vt 0.500000 1.000000", want: []float32{0.5, 1.0}},
	}

	for _, tc := range tests {
		got := parseVec2(tc.input)
		assert.Equal(t, tc.want, got, tc.name)
	}
}

func TestParseFace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []uint32
	}{
		{name: "simple", input: "f 1/1/1 2/2/1 3/3/1", want: []uint32{1, 1, 1, 2, 2, 1, 3, 3, 1}},
	}

	for _, tc := range tests {
		got := parseFace(tc.input)
		assert.Equal(t, tc.want, got, tc.name)
	}
}

func TestBufferData(t *testing.T) {
	tests := []struct {
		name  string
		input io.Reader
		want  []float32
	}{
		{name: "simple", input: strings.NewReader(testTriangle), want: testTriangleBufferData},
	}

	for _, tc := range tests {
		obj := Load(tc.input)
		got := obj.BufferData()
		assert.Equal(t, tc.want, got, tc.name)
	}
}

var testTriangle = `# Blender 3.5.1
# www.blender.org
o Cube
v 1.000000 1.000000 -1.000000
v 1.000000 1.000000 1.000000
v 1.000000 -1.000000 1.000000
vn 1.0000 -0.0000 -0.0000
vt 0.500000 1.000000
vt 0.000000 1.000000
vt 0.000000 0.500000
s 0
f 1/1/1 2/2/1 3/3/1
`

var testTriangleBufferData = []float32{
	1.0, 1.0, -1.0, 0.5, 1.0, 1.0, -0.0, -0.0,
	1.0, 1.0, 1.0, 0.0, 1.0, 1.0, -0.0, -0.0,
	1.0, -1.0, 1.0, 0.0, 0.5, 1.0, -0.0, -0.0,
}

var testTrianglesSharp = `# Blender 3.5.1
# www.blender.org
o Cube
v 1.000000 1.000000 -1.000000
v 1.000000 -1.000000 -1.000000
v 1.000000 1.000000 1.000000
v 1.000000 -1.000000 1.000000
v -1.000000 1.000000 -1.000000
v -1.000000 -1.000000 -1.000000
vn 1.0000 -0.0000 -0.0000
vn -0.0000 -0.0000 -1.0000
vt 0.500000 1.000000
vt 0.000000 0.500000
vt 0.500000 0.500000
vt -0.000000 0.000000
vt 0.000000 1.000000
vt 0.000000 0.500000
vt 0.500000 0.500000
vt 0.500000 -0.000000
s 0
f 1/1/1 4/6/1 2/3/1
f 5/7/2 2/4/2 6/8/2
f 1/1/1 3/5/1 4/6/1
f 5/7/2 1/2/2 2/4/2
`

var testTrianglesSmooth = `# Blender 3.5.1
# www.blender.org
o Cube
v 1.000000 1.000000 -1.000000
v 1.000000 -1.000000 -1.000000
v 1.000000 1.000000 1.000000
v 1.000000 -1.000000 1.000000
v -1.000000 1.000000 -1.000000
v -1.000000 -1.000000 -1.000000
vn 0.7071 -0.0000 -0.7071
vn 1.0000 -0.0000 -0.0000
vn -0.0000 -0.0000 -1.0000
vt 0.500000 1.000000
vt 0.000000 0.500000
vt 0.500000 0.500000
vt -0.000000 0.000000
vt 0.000000 1.000000
vt 0.000000 0.500000
vt 0.500000 0.500000
vt 0.500000 -0.000000
s 1
f 1/1/1 4/6/2 2/3/1
f 5/7/3 2/4/1 6/8/3
f 1/1/1 3/5/2 4/6/2
f 5/7/3 1/2/1 2/4/1
`
