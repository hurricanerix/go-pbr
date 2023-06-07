package material

type PBR struct {
	TexturePaths map[Texture]string
}

func (m PBR) Textures() map[Texture]string {
	return m.TexturePaths
}

func (m PBR) Shader() string {
	return "pbr"
}

func (m PBR) Use() error {
	return nil
}
