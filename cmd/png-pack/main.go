package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
	"path"
)

func main() {
	rSuffix := flag.String("r", "ao", "")
	gSuffix := flag.String("g", "roughness", "")
	bSuffix := flag.String("b", "metallic", "")
	aSuffix := flag.String("a", "", "")

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		panic("too many files")
	}

	outFile, err := os.Create(args[0])
	defer outFile.Close()
	if err != nil {
		panic(err)
	}

	//mode_rgba := map[string]string{"r": "r", "g": "g", "b": "b", "a": "a"}
	mode := map[string]string{}
	if *rSuffix != "" {
		mode["r"] = *rSuffix
	}
	if *gSuffix != "" {
		mode["g"] = *gSuffix
	}
	if *bSuffix != "" {
		mode["b"] = *bSuffix
	}
	if *aSuffix != "" {
		mode["a"] = *aSuffix
	}

	outfileDir := path.Dir(args[0])
	inImages := map[string]image.Image{}

	b := image.Rectangle{}
	for component, suffix := range mode {
		inFile := fmt.Sprintf("%s/%s.png", outfileDir, suffix)
		inImgFile, err := os.Open(inFile)
		defer inImgFile.Close()
		if err != nil {
			panic(err)
		}
		inImages[component], err = png.Decode(inImgFile)
		if err != nil {
			panic(err)
		}

		b = inImages[component].Bounds()
		// TODO: check all image bounds to match
	}

	outImgFile := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r := uint32(0)
			g := uint32(0)
			b := uint32(0)
			a := uint32(255 << 8)

			if ao, ok := inImages["r"]; ok {
				r, _, _, _ = ao.At(x, y).RGBA()
				//r, g, b, a = ao.At(x, y).RGBA()
			}

			if roughness, ok := inImages["g"]; ok {
				_, g, _, _ = roughness.At(x, y).RGBA()
				//r, g, b, a = roughness.At(x, y).RGBA()
			}

			if metalic, ok := inImages["b"]; ok {
				_, _, b, _ = metalic.At(x, y).RGBA()
			}

			//if ao, ok := inImages["r"]; ok {
			//	r, _, _, _ = ao.At(x, y).RGBA()
			//}
			//tmp := uint8(r >> 8)
			//println(tmp)
			outImgFile.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	if err := png.Encode(outFile, outImgFile); err != nil {
		panic(err)
	}
}
