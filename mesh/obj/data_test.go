package obj

import "go-pbr/mesh"

var BlenderObjHeader = `# Blender 3.5.1
# www.blender.org
`

var TriangleOBJ = `o Cube
# Vertex Positions
v 1.000000 1.000000 -1.000000
v 1.000000 1.000000 1.000000
v 1.000000 -1.000000 1.000000

# Vertex Normals
vn 1.0000 -0.0000 -0.0000

# Vertex Texture Coordinates
vt 0.500000 1.000000
vt 0.000000 1.000000
vt 0.000000 0.500000

# Surface 0
s 0

# Faces
f 1/1/1 2/2/1 3/3/1
`

var TriangleConfig = mesh.Config{
	Stride: 8 * 4,
	Fields: map[mesh.VertexAttribute]mesh.VertexAttributeData{
		mesh.Position: {Size: 3, Offset: 0},
		mesh.TexCoord: {Size: 2, Offset: 3 * 4},
		mesh.Normal:   {Size: 3, Offset: 5 * 4},
	},
}

var TriangleData = []float32{
	1.0, 1.0, -1.0, 0.5, 1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.5, 1.0, 0.0, 0.0,
}

var TriangleConfigWithTangents = mesh.Config{
	Stride: 14 * 4,
	Fields: map[mesh.VertexAttribute]mesh.VertexAttributeData{
		mesh.Position:  {Size: 3, Offset: 0},
		mesh.TexCoord:  {Size: 2, Offset: 3 * 4},
		mesh.Normal:    {Size: 3, Offset: 5 * 4},
		mesh.Tangent:   {Size: 3, Offset: 8 * 4},
		mesh.Bitangent: {Size: 3, Offset: 11 * 4},
	},
}

var TriangleDataWithTangents = []float32{
	1.0, 1.0, -1.0, 0.5, 1.0, 1.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0,
}
