package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
)

type Material struct {
	Prog    uint32
	Path    string
	diffuse uint32
	arm     uint32
	disp    uint32
	norm    uint32
}

func NewMaterial(prog uint32, path string) *Material {
	return &Material{Prog: prog, Path: path}
}

func (m *Material) Load() {
	diffuseUniform := gl.GetUniformLocation(m.Prog, gl.Str(DiffuseSamplerKey+"\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	normUniform := gl.GetUniformLocation(m.Prog, gl.Str(NormalSamplerKey+"\x00"))
	gl.Uniform1i(normUniform, 1)

	dispUniform := gl.GetUniformLocation(m.Prog, gl.Str(DispSamplerKey+"\x00"))
	gl.Uniform1i(dispUniform, 2)

	armUniform := gl.GetUniformLocation(m.Prog, gl.Str(ARMSamplerKey+"\x00"))
	gl.Uniform1i(armUniform, 3)

	var err error
	m.diffuse, err = NewTexture(m.Path+"/diffuse_1k.png", gl.TEXTURE0)
	if err != nil {
		log.Fatalln(err)
	}

	m.norm, err = NewTexture(m.Path+"/nor_gl_1k.png", gl.TEXTURE1)
	if err != nil {
		log.Fatalln(err)
	}

	m.disp, err = NewTexture(m.Path+"/disp_1k.png", gl.TEXTURE2)
	if err != nil {
		log.Fatalln(err)
	}

	m.arm, err = NewTexture(m.Path+"/arm_1k.png", gl.TEXTURE3)
	if err != nil {
		log.Fatalln(err)
	}

	gl.BindFragDataLocation(m.Prog, 0, gl.Str(FragDataLocation+"\x00"))
}

func (m *Material) Use() {
	//ambientStrength := float32(0.2)
	//ambientStrengthUniform := gl.GetUniformLocation(m.Prog, gl.Str("AmbientStrength\x00"))
	//gl.Uniform1f(ambientStrengthUniform, ambientStrength)
	//
	////ambiantColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	//ambiantColor := mgl32.Vec3{0.5568, 0.4039, 0.3529}
	//ambiantColorUniform := gl.GetUniformLocation(m.Prog, gl.Str("AmbientColor\x00"))
	//gl.Uniform3fv(ambiantColorUniform, 1, &ambiantColor[0])

	// 255,244,161
	//lightColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	//lightColor := mgl32.Vec3{0.6568, 0.5039, 0.4529}
	//lightColor := mgl32.Vec3{1, 1, 1}
	//lightColor := mgl32.Vec3{150, 150, 150}
	//lightColorUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightColor\x00"))
	//gl.Uniform3fv(lightColorUniform, 1, &lightColor[0])

	//lightPower := float32(0.2)
	//lightPowerUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightPower\x00"))
	//gl.Uniform1f(lightPowerUniform, lightPower)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.diffuse)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, m.norm)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, m.disp)

	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, m.arm)
}

func (m *Material) Free() {
	gl.DeleteTextures(1, &m.diffuse)
	gl.DeleteTextures(1, &m.arm)
	gl.DeleteTextures(1, &m.disp)
	gl.DeleteTextures(1, &m.norm)
}
