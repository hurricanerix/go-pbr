package opengl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io"
	"strings"
)

type ShaderType int

const (
	VertexShader   ShaderType = gl.VERTEX_SHADER
	FragmentShader            = gl.FRAGMENT_SHADER
)

type Program struct {
	handle uint32
	linked bool
}

func New() *Program {
	return &Program{}
}

func (p *Program) Destroy() {
	if p.handle == 0 {
		return
	}

	p.detachAndDeleteShaderObjects()
	// gl.DeleteProgram()
}

func (p *Program) Handle() uint32 {
	return p.handle
}

func (p *Program) Linked() bool {
	return p.linked
}

func (p *Program) detachAndDeleteShaderObjects() {

}

func (p *Program) deleteProgram() {

}

func (p *Program) CompileShader(file io.Reader, shaderType ShaderType) {
	if p.handle <= 0 {
		p.handle = gl.CreateProgram()
	}

	shaderHandle := gl.CreateShader(uint32(shaderType))

	source := read(file)

	csources, free := gl.Strs(source)
	defer free()

	gl.ShaderSource(shaderHandle, 1, csources, nil)

	gl.CompileShader(shaderHandle)

	var status int32
	gl.GetShaderiv(shaderHandle, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderHandle, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to compile %v: %v", source, log))
	}

	gl.AttachShader(p.handle, shaderHandle)
}

func (p *Program) Link() {

	gl.LinkProgram(p.handle)

	var status int32
	gl.GetProgramiv(p.handle, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(p.handle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p.handle, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to link program: %v", log))
	}

	p.linked = true
	//		gl.DeleteShader(vertexShader)
	//		gl.DeleteShader(fragmentShader)
}

func (p *Program) FindUniformLocations() {

}

func (p *Program) Use() {
	gl.UseProgram(p.handle)
}

func (p *Program) BindAttribLocation() {

}

func (p *Program) BindFragDataLocation() {

}

// Generics?
func (p *Program) SetUniform() {
	//camera := mgl32.Vec3{3, 3, 3}
	//cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	//gl.Uniform3fv(cameraUniform, 1, &camera[0])
}

func (p *Program) PrintActiveUniforms() {

}

func (p *Program) PrintActiveUniformBlocks() {

}

func (p *Program) PrintActiveAttribs() {

}

func (p *Program) Validate() error {
	if !p.linked {
		return fmt.Errorf("Program is not linked")
	}

	var status int32
	gl.ValidateProgram(p.handle)
	gl.GetProgramiv(p.handle, gl.VALIDATE_STATUS, &status)

	if gl.FALSE == status {
		var logLength int32
		gl.GetProgramiv(p.handle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p.handle, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}

	return nil
}

func (p *Program) getExtension(name string) string {
	return ""
}

//func NewTexture(file string, loc uint32) (uint32, error) {
//	imgFile, err := os.Open(file)
//	if err != nil {
//		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
//	}
//	img, _, err := image.Decode(imgFile)
//	if err != nil {
//		return 0, err
//	}
//
//	rgba := image.NewRGBA(img.Bounds())
//	if rgba.Stride != rgba.Rect.Size().X*4 {
//		return 0, fmt.Errorf("unsupported stride")
//	}
//	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
//
//	var texture uint32
//	gl.GenTextures(1, &texture)
//	gl.ActiveTexture(loc)
//	gl.BindTexture(gl.TEXTURE_2D, texture)
//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
//	gl.TexImage2D(
//		gl.TEXTURE_2D,
//		0,
//		gl.RGBA,
//		int32(rgba.Rect.Size().X),
//		int32(rgba.Rect.Size().Y),
//		0,
//		gl.RGBA,
//		gl.UNSIGNED_BYTE,
//		gl.Ptr(rgba.Pix))
//
//	return texture, nil
//}

func read(f io.Reader) string {
	result := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		result = append(result, buf[:n]...)
	}
	return string(result)
}
