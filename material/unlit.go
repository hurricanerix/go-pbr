package material

type Unlit struct {
}

func (m Unlit) Shader() string {
	return "unlit"
}

func (m Unlit) Use() error {
	return nil
}
