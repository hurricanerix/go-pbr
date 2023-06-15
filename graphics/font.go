package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"os"
)

type RasterFont struct {
	ScreenWidth  float32
	ScreenHeight float32
	program      Program
	vao          uint32
	vbo          uint32
	textureID    uint32 // ID handle of the glyph texture
	characters   map[rune]Character
	color        mgl32.Vec3
}

// Holds all state information relevant to a character as loaded using FreeType
type Character struct {
	Size    mgl32.Vec2 // Size of glyph
	Bearing mgl32.Vec2 // Offset from baseline to left/top of glyph
	Advance uint32     // Horizontal offset to advance to next glyph
}

// NewFont returns
func NewFont(fontpath string, cw, ch float32, size, width, height float32) (*RasterFont, error) {
	fnt := RasterFont{
		ScreenWidth:  width,
		ScreenHeight: height,
	}

	// compile and setup the program
	// ----------------------------
	fnt.program = Program{}
	if f, err := os.Open("assets/shaders/font/shader.vert"); err == nil {
		fnt.program.CompileShader(f, VertexShader)
		f.Close()
	} else {
		return nil, err
	}
	if f, err := os.Open("assets/shaders/font/shader.frag"); err == nil {
		fnt.program.CompileShader(f, FragmentShader)
		f.Close()
	} else {
		return nil, err
	}
	fnt.program.Link()

	// configure VAO/VBO for texture quads
	// -----------------------------------
	gl.GenVertexArrays(1, &fnt.vao)
	gl.GenBuffers(1, &fnt.vbo)
	gl.BindVertexArray(fnt.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, fnt.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, 4*6*4, nil, gl.DYNAMIC_DRAW)

	// 0: = 0, 0
	// 1: = 1, 0
	// 2: = 0, 1
	// 3: = 1, 1
	//uvs := fnt.getUVs('a')
	//vertexData := []float32{
	//	0, 64, uvs[0].X(), uvs[0].Y(),
	//	0, 0, uvs[2].X(), uvs[2].Y(),
	//	64, 0, uvs[3].X(), uvs[3].Y(),
	//	0, 64, uvs[0].X(), uvs[0].Y(),
	//	64, 0, uvs[3].X(), uvs[3].Y(),
	//	64, 64, uvs[1].X(), uvs[1].Y(),
	//}
	//vertexData := []float32{
	//	0, 0, 0, 0,
	//	0, 0, 0, 0,
	//	0, 0, 0, 0,
	//	0, 0, 0, 0,
	//	0, 0, 0, 0,
	//	0, 0, 0, 0,
	//}
	//fmt.Printf("%v\n", vertexData)
	gl.BufferData(gl.ARRAY_BUFFER, 4*6*4, gl.Ptr(nil), gl.DYNAMIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)

	if err := fnt.program.Validate(); err != nil {
		return nil, err
	}
	fnt.program.Use()

	// TODO: Configure SCR_WIDTH/SCR_HEIGHT
	projection := mgl32.Ortho2D(0.0, fnt.ScreenWidth, 0.0, fnt.ScreenHeight)
	fnt.program.SetUniformMatrix4fv("projection", projection)

	textureID, err := NewTexture(fontpath, gl.TEXTURE0)
	if err != nil {
		return nil, err
	}
	fnt.textureID = textureID

	fnt.characters = make(map[rune]Character, 128)
	//for r := rune(32); r < 127; r++ {
	for r := rune(0); r < 127; r++ {
		character := Character{
			Size:    mgl32.Vec2{cw, ch},
			Bearing: mgl32.Vec2{0, 0},
			Advance: uint32(cw),
		}
		fnt.characters[r] = character
	}

	// Cleanup
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return &fnt, nil
}

func (f RasterFont) getUVs(r rune) [4]mgl32.Vec2 {
	c := int(r)
	// TODO: Fix showing "not found" rune.
	//if _, ok := f.characters[r]; !ok {
	//	c = 0
	//}
	x := float32(c % 16)
	y := float32(c / 16)

	scale := float32(256)
	xs := x * 16 / scale
	ys := y * 16 / scale
	xe := (x*16 + 16) / scale
	ye := (y*16 + 16) / scale

	result := [4]mgl32.Vec2{
		{xs, ys},
		{xe, ys},
		{xs, ye},
		{xe, ye},
	}
	return result
}

// TODO: reanme all `f` to `fnt`.
func (f *RasterFont) Activate() {
	f.program.Use()
	gl.DepthMask(false)
	gl.BindFragDataLocation(f.program.Handle(), 0, gl.Str("color\x00"))
	projection := mgl32.Ortho2D(0.0, f.ScreenWidth, 0.0, f.ScreenHeight)
	f.program.SetUniformMatrix4fv("projection", projection)
	f.program.SetUniform3fv("textColor", f.color)
	gl.BindVertexArray(f.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(f.vao)
	gl.BindTexture(gl.TEXTURE_2D, f.textureID)
	//gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//// update content of VBO memory
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)
}

func (f *RasterFont) Color(c mgl32.Vec3) {
	f.color = c
	f.program.SetUniform3fv("textColor", f.color)
}

// render line of text
// -------------------
func (f *RasterFont) RenderText(text string, x, y, scale float32) {
	// activate corresponding render state

	// iterate through all characters
	for _, r := range text {
		ch := f.characters[r]

		// TODO: implement scale
		xpos := x + ch.Bearing.X()                 //*scale
		ypos := y - (ch.Size.Y() - ch.Bearing.Y()) //*scale
		w := ch.Size.X()                           // * scale
		h := ch.Size.Y()                           // * scale

		//update VBO for each character
		uvs := f.getUVs(r)
		vertices := []float32{
			xpos, ypos + h, uvs[0].X(), uvs[0].Y(),
			xpos, ypos, uvs[2].X(), uvs[2].Y(),
			xpos + w, ypos, uvs[3].X(), uvs[3].Y(),

			xpos, ypos + h, uvs[0].X(), uvs[0].Y(),
			xpos + w, ypos, uvs[3].X(), uvs[3].Y(),
			xpos + w, ypos + h, uvs[1].X(), uvs[1].Y(),
		}

		// TODO: SIZE WAS WRONG
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*6*4, gl.Ptr(vertices)) // be sure to use glBufferSubData and not glBufferData

		//render quad
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		// now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		x += float32((ch.Advance)) // >> 6)) // * scale // bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))
	}

}

func (f RasterFont) Deactivate() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.DepthMask(true)
}
