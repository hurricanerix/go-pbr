package opengl

//
//import (
//	"fmt"
//	"image"
//	"os"
//)
//
//func (r *Renderer) SetCubemap(path string) {
//	gl.GenTextures(1, &r.CubeMapTex)
//	gl.BindTexture(gl.TEXTURE_CUBE_MAP, r.CubeMapTex)
//	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
//	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
//	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
//	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
//	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
//
//	files := map[uint32]string{
//		gl.TEXTURE_CUBE_MAP_POSITIVE_X: "px.png",
//		gl.TEXTURE_CUBE_MAP_NEGATIVE_X: "nx.png",
//		gl.TEXTURE_CUBE_MAP_POSITIVE_Y: "py.png",
//		gl.TEXTURE_CUBE_MAP_NEGATIVE_Y: "ny.png",
//		gl.TEXTURE_CUBE_MAP_POSITIVE_Z: "pz.png",
//		gl.TEXTURE_CUBE_MAP_NEGATIVE_Z: "nz.png",
//	}
//
//	for target, filename := range files {
//		fullpath := path + "/" + filename
//
//		imgFile, err := os.Open(fullpath)
//		if err != nil {
//			panic(fmt.Errorf("texture %q not found on disk: %v", fullpath, err))
//		}
//		img, _, err := image.Decode(imgFile)
//		if err != nil {
//			panic(err)
//		}
//
//		rgba := image.NewRGBA(img.Bounds())
//		if rgba.Stride != rgba.Rect.Size().X*4 {
//			panic(fmt.Errorf("unsupported stride"))
//		}
//		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
//
//		gl.TexImage2D(
//			target,
//			0,
//			gl.RGBA,
//			int32(rgba.Rect.Size().X),
//			int32(rgba.Rect.Size().Y),
//			0,
//			gl.RGBA,
//			gl.UNSIGNED_BYTE,
//			gl.Ptr(rgba.Pix))
//	}
//
//	r.cubemapShader = opengl.Program{}
//	if f, err := os.Open("assets/shaders/cubemap.vert"); err == nil {
//		defer f.Close()
//		r.cubemapShader.CompileShader(f, opengl.VertexShader)
//	}
//
//	if f, err := os.Open("assets/shaders/cubemap.frag"); err == nil {
//		defer f.Close()
//		r.cubemapShader.CompileShader(f, opengl.FragmentShader)
//	}
//
//	r.cubemapShader.Link()
//
//	//if f, err := os.Open(path + "/geometry.mesh"); err == nil {
//	//	defer f.Close()
//	//	r.cubemapMesh = mesh.Load(f)
//	//} else {
//	//	panic(err)
//	//}
//	//
//	//r.cubemapMesh.Bind()
//
//	if err := r.cubemapShader.Validate(); err != nil {
//		panic(err)
//	}
//}
