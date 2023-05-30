package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics/opengl"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
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
	diffuseUniform := gl.GetUniformLocation(m.Prog, gl.Str("DiffuseSampler\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	armUniform := gl.GetUniformLocation(m.Prog, gl.Str("ArmSampler\x00"))
	gl.Uniform1i(armUniform, 1)

	dispUniform := gl.GetUniformLocation(m.Prog, gl.Str("DispSampler\x00"))
	gl.Uniform1i(dispUniform, 2)

	normUniform := gl.GetUniformLocation(m.Prog, gl.Str("NormalSampler\x00"))
	gl.Uniform1i(normUniform, 3)

	var err error
	m.diffuse, err = opengl.NewTexture(m.Path+"/diffuse_1k.png", gl.TEXTURE0)
	if err != nil {
		log.Fatalln(err)
	}

	m.arm, err = opengl.NewTexture(m.Path+"/arm_1k.png", gl.TEXTURE1)
	if err != nil {
		log.Fatalln(err)
	}

	m.disp, err = opengl.NewTexture(m.Path+"/disp_1k.png", gl.TEXTURE2)
	if err != nil {
		log.Fatalln(err)
	}

	m.norm, err = opengl.NewTexture(m.Path+"/nor_gl_1k.png", gl.TEXTURE3)
	if err != nil {
		log.Fatalln(err)
	}
	gl.BindFragDataLocation(m.Prog, 0, gl.Str("FragColor\x00"))
}

func (m *Material) Use() {
	ambientStrength := float32(0.3)
	ambientStrengthUniform := gl.GetUniformLocation(m.Prog, gl.Str("AmbientStrength\x00"))
	gl.Uniform1f(ambientStrengthUniform, ambientStrength)

	ambiantColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	ambiantColorUniform := gl.GetUniformLocation(m.Prog, gl.Str("AmbientColor\x00"))
	gl.Uniform3fv(ambiantColorUniform, 1, &ambiantColor[0])

	lightPos := mgl32.Vec3{1.2, 1.0, 2.0}
	lightPosUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightPos\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])

	// 255,244,161
	lightColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	lightColorUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightColor\x00"))
	gl.Uniform3fv(lightColorUniform, 1, &lightColor[0])

	lightPower := float32(10)
	lightPowerUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightPower\x00"))
	gl.Uniform1f(lightPowerUniform, lightPower)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.diffuse)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, m.arm)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, m.disp)

	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, m.norm)
}

func (m *Material) Free() {
	gl.DeleteTextures(1, &m.diffuse)
	gl.DeleteTextures(1, &m.arm)
	gl.DeleteTextures(1, &m.disp)
	gl.DeleteTextures(1, &m.norm)
}
