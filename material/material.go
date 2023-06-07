package material

type Texture string

const (
	DiffuseMap = "diffuseMap"
	NormalMap  = "normalMap"
	DispMap    = "dispMap"
	ARMMap     = "armMap"
)

type Material interface {
	Textures() map[Texture]string
	Shader() string
	Use() error
}
