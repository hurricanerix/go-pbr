package mesh

// Mesh is responsible for providing the information needed to reconstruct the geometry.
type Mesh interface {
	Config() Config
	Data() []float32
}

// Config is responsible for providing the data needed to decode the data.
type Config struct {
	// Specifies the byte offset between consecutive generic vertex attributes. If stride is 0, the generic vertex
	// attributes are understood to be tightly packed in the array. The initial value is 0.
	Stride uint32

	Fields map[VertexAttribute]VertexAttributeData
}

// VertexAttribute that can be provided.
type VertexAttribute int

const (
	// Position of the vertex in 3D space.
	Position VertexAttribute = iota
	// Normal for the vertex, this can be assumed to be normalized.
	Normal
	// TexCoord associated with the vertex.
	TexCoord
	// Tangent for the vertex, in relation to the Normal & Bitangent.
	Tangent
	// Bitangent for the vertex, in relation to the Normal & Tangent.
	Bitangent
)

// VertexAttributeData represents how the components are stored in the Data slice.
type VertexAttributeData struct {
	// Specifies the number of components for the vertex attribute.
	Size int32
	// Offset of the first component of the first vertex attribute in the Data slice.
	Offset uintptr
}
