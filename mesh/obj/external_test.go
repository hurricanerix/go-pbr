package obj_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-pbr/mesh"
	"go-pbr/mesh/obj"
	"io"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name         string
		input        io.Reader
		inputOptions mesh.DecodeOptions
		wantConfig   mesh.Config
		wantData     []float32
	}{
		{
			name:         "simple",
			input:        strings.NewReader(obj.TriangleOBJ),
			inputOptions: mesh.DecodeOptions{SkipHeaderCheck: true},
			wantConfig:   obj.TriangleConfig,
			wantData:     obj.TriangleData,
		},
		{
			name:         "simple-with-tangents",
			input:        strings.NewReader(obj.TriangleOBJ),
			inputOptions: mesh.DecodeOptions{SkipHeaderCheck: true, CalculateTangentsIfMissing: true},
			wantConfig:   obj.TriangleConfigWithTangents,
			wantData:     obj.TriangleDataWithTangents,
		},
		{
			name:         "simple-blender",
			input:        strings.NewReader(obj.BlenderObjHeader + obj.TriangleOBJ),
			inputOptions: mesh.DecodeOptions{},
			wantConfig:   obj.TriangleConfig,
			wantData:     obj.TriangleData,
		},
		{
			name:         "simple-blender-with-tangents",
			input:        strings.NewReader(obj.BlenderObjHeader + obj.TriangleOBJ),
			inputOptions: mesh.DecodeOptions{CalculateTangentsIfMissing: true},
			wantConfig:   obj.TriangleConfigWithTangents,
			wantData:     obj.TriangleDataWithTangents,
		},
	}

	for _, tc := range tests {
		got, err := obj.Decode(tc.input, tc.inputOptions)
		require.Nil(t, err, "%s returned error", tc.name)
		require.NotNil(t, got, "%s received nil mesh", tc.name)
		assert.Equal(t, tc.wantConfig, got.Config())
		assert.Equal(t, tc.wantData, got.Data(), tc.name)
	}
}
