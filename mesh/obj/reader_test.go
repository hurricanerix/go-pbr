package obj

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestCalculateTB(t *testing.T) {
	tests := []struct {
		name      string
		inputPos  [3]mgl32.Vec3
		inputUV   [3]mgl32.Vec2
		wantTan   mgl32.Vec3
		wantBiTan mgl32.Vec3
	}{
		{
			name:      "triangle1",
			inputPos:  [3]mgl32.Vec3{{-1, 1, 0}, {-1, -1, 0}, {1, -1, 0}},
			inputUV:   [3]mgl32.Vec2{{0, 1}, {0, 0}, {1, 0}},
			wantTan:   mgl32.Vec3{1, 0, 0},
			wantBiTan: mgl32.Vec3{0, 1, 0},
		},
		{
			name:      "triangle2",
			inputPos:  [3]mgl32.Vec3{{-1, 1, 0}, {1, -1, 0}, {1, 1, 0}},
			inputUV:   [3]mgl32.Vec2{{0, 1}, {1, 0}, {1, 1}},
			wantTan:   mgl32.Vec3{1, 0, 0},
			wantBiTan: mgl32.Vec3{0, 1, 0},
		},
		{
			name:      "triangle3",
			inputPos:  [3]mgl32.Vec3{{-1, 1, 1}, {-1, -1, -1}, {-1, -1, 1}},
			inputUV:   [3]mgl32.Vec2{{1, 1}, {0, 0}, {1, 0}},
			wantTan:   mgl32.Vec3{0, 0, 1},
			wantBiTan: mgl32.Vec3{0, 1, 0},
		},
	}

	for _, tc := range tests {
		gotTan, gotBiTan := calculateTB(
			tc.inputPos[0], tc.inputPos[1], tc.inputPos[2],
			tc.inputUV[0], tc.inputUV[1], tc.inputUV[2])
		assert.Equal(t, tc.wantTan, gotTan, tc.name+" tan")
		assert.Equal(t, tc.wantBiTan, gotBiTan, tc.name+" bitan")
	}
}
