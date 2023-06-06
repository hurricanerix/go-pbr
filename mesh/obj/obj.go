package obj

import "go-pbr/mesh"

// Obj represents a wavefront mesh implementation.
type Obj struct {
	config mesh.Config
	data   []float32
}

// Config to be used when reconstructing the mesh.
func (o Obj) Config() mesh.Config {
	return o.config
}

// Data to be used when reconstructing the mesh.
func (o Obj) Data() []float32 {
	return o.data
}
