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
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	//mode_rgba := map[string]string{"r": "r", "g": "g", "b": "b", "a": "a"}
	mode_arm := map[string]string{"r": "ao", "g": "roughness", "b": "metallic"}

	for _, infile := range args {
		inImgFile, err := os.Open(infile)
		defer inImgFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file '%s'\n", infile)
			continue
		}

		inImg, _, err := image.Decode(inImgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not decode file '%s'\n", infile)
			continue
		}

		outImages := map[string]*image.RGBA{}
		for component, _ := range mode_arm {
			outImages[component] = image.NewRGBA(inImg.Bounds())
		}

		b := inImg.Bounds()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				col := inImg.At(x, y)
				r, g, b, _ := col.RGBA()
				if ao, ok := outImages["r"]; ok {
					ao.Set(x, y, color.RGBA{uint8(r >> 8), uint8(r >> 8), uint8(r >> 8), 255})
				}

				if roughness, ok := outImages["g"]; ok {
					roughness.Set(x, y, color.RGBA{uint8(g >> 8), uint8(g >> 8), uint8(g >> 8), 255})
				}

				if metalic, ok := outImages["b"]; ok {
					metalic.Set(x, y, color.RGBA{uint8(b >> 8), uint8(b >> 8), uint8(b >> 8), 255})
				}
			}
		}

		infileNoExt, _ := strings.CutSuffix(infile, path.Ext(infile))
		for component, suffix := range mode_arm {
			if _, ok := outImages[component]; !ok {
				continue
			}

			outfile := fmt.Sprintf("%s_%s.png", infileNoExt, suffix)
			outImgFile, err := os.Create(outfile)
			defer outImgFile.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not create '%s'\n", outfile)
				continue
			}

			//outImgFile.Write(outImages[component].)
			if err := png.Encode(outImgFile, outImages[component]); err != nil {
				fmt.Fprintf(os.Stderr, "Could not write '%s'\n", outfile)
				continue
			}
		}
	}
}
