package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics/opengl"
	"log"
)

type Material interface {
	Use() error
	Free() error
}

type LitMaterial struct {
	Path string
}

func (m *LitMaterial) Load() error {
	return nil
}

func (m LitMaterial) Use() error {
	return nil
}

func (m LitMaterial) Free() error {
	return nil
}

type OldMaterial struct {
	Prog    uint32
	Path    string
	diffuse uint32
	arm     uint32
	disp    uint32
	norm    uint32
}

func NewLitMaterial(materialPath string) *OldMaterial {
	return nil
}

func NewMaterial(prog uint32, path string) *OldMaterial {
	return &OldMaterial{Prog: prog, Path: path}
}

func (m *OldMaterial) Load() {
	diffuseUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.DiffuseSamplerKey+"\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	normUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.NormalSamplerKey+"\x00"))
	gl.Uniform1i(normUniform, 1)

	dispUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.DispSamplerKey+"\x00"))
	gl.Uniform1i(dispUniform, 2)

	armUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.ARMSamplerKey+"\x00"))
	gl.Uniform1i(armUniform, 3)

	var err error
	m.diffuse, err = opengl.NewTexture(m.Path+"/diffuse_1k.png", gl.TEXTURE0)
	if err != nil {
		log.Fatalln(err)
	}

	m.norm, err = opengl.NewTexture(m.Path+"/nor_gl_1k.png", gl.TEXTURE1)
	if err != nil {
		log.Fatalln(err)
	}

	m.disp, err = opengl.NewTexture(m.Path+"/disp_1k.png", gl.TEXTURE2)
	if err != nil {
		log.Fatalln(err)
	}

	m.arm, err = opengl.NewTexture(m.Path+"/arm_1k.png", gl.TEXTURE3)
	if err != nil {
		log.Fatalln(err)
	}

	gl.BindFragDataLocation(m.Prog, 0, gl.Str(opengl.FragDataLocation+"\x00"))
}

func (m *OldMaterial) Use() {
	ambientStrength := float32(0.2)
	ambientStrengthUniform := gl.GetUniformLocation(m.Prog, gl.Str("AmbientStrength\x00"))
	gl.Uniform1f(ambientStrengthUniform, ambientStrength)

	//ambiantColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	ambiantColor := mgl32.Vec3{0.5568, 0.4039, 0.3529}
	ambiantColorUniform := gl.GetUniformLocation(m.Prog, gl.Str("AmbientColor\x00"))
	gl.Uniform3fv(ambiantColorUniform, 1, &ambiantColor[0])

	lightPos := mgl32.Vec3{3.0, 0.0, 0.0}
	lightPosUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.LightPosKey+"\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])

	// 255,244,161
	//lightColor := mgl32.Vec3{1.0, 0.9569, 0.6314}
	//lightColor := mgl32.Vec3{0.6568, 0.5039, 0.4529}
	lightColor := mgl32.Vec3{1, 1, 1}
	lightColorUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightColor\x00"))
	gl.Uniform3fv(lightColorUniform, 1, &lightColor[0])

	lightPower := float32(0.2)
	lightPowerUniform := gl.GetUniformLocation(m.Prog, gl.Str("LightPower\x00"))
	gl.Uniform1f(lightPowerUniform, lightPower)

	enableDisplacement := gl.GetUniformLocation(m.Prog, gl.Str(opengl.EnableDisplacementKey+"\x00"))
	gl.Uniform1i(enableDisplacement, 1)

	dispDiscardUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.DispDiscardOutOfBoundsKey+"\x00"))
	gl.Uniform1i(dispDiscardUniform, 0)

	hightScale := float32(0.1)
	hightScaleUniform := gl.GetUniformLocation(m.Prog, gl.Str(opengl.HeightScaleKey+"\x00"))
	gl.Uniform1f(hightScaleUniform, hightScale)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.diffuse)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, m.norm)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, m.disp)

	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, m.arm)
}

func (m *OldMaterial) Free() {
	gl.DeleteTextures(1, &m.diffuse)
	gl.DeleteTextures(1, &m.arm)
	gl.DeleteTextures(1, &m.disp)
	gl.DeleteTextures(1, &m.norm)
}
