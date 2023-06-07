package main

import (
	"flag"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"go-pbr/graphics"
	"go-pbr/material"
	"go-pbr/mesh"
	_ "go-pbr/mesh/obj"
	"go-pbr/scene"
	"log"
	"os"
)

func main() {
	meshPath := flag.String("mesh", "assets/objects/cube.obj", "")
	materialPath := flag.String("material", "assets/materials/lgl_pbr_rusted_iron", "")
	//skymapPath := flag.String("skymap", "assets/cubemaps/castle-zavelstein-cellar", "")

	flag.Parse()

	app := &app{
		Title:        "PBR Viewer",
		WindowWidth:  960,
		WindowHeight: 540,
	}
	if err := app.Init(); err != nil {
		log.Fatalln("failed to initialize app:", err)
	}
	defer func() {
		if err := app.Free(); err != nil {
			log.Fatalln("failed to free app:", err)
		}
	}()

	fmt.Printf("Graphics Version: %s\n", app.backend.Version)
	fmt.Printf("Graphics Shading Language Version: %s\n", app.backend.ShadingLanguageVersion)

	meshFile, err := os.Open(*meshPath)
	if err != nil {
		panic(err)
	}

	subject, _, err := mesh.Decode(meshFile, mesh.DecodeOptions{SkipHeaderCheck: true})
	if err != nil {
		panic(err)
	}

	app.Scene = scene.Scene{
		0: scene.Object{
			Transform: graphics.Transform{Scale: mgl32.Vec3{1, 1, 1}},
			Renderer: graphics.Renderer{
				Program: graphics.NewProgram("assets/shaders/debug.vert", "assets/shaders/debug.frag"),
				Mesh:    subject,
				Material: material.PBR{
					TexturePaths: map[material.Texture]string{
						material.DiffuseMap: fmt.Sprintf("%s/diffuse_1k.png", materialPath),
						material.NormalMap:  fmt.Sprintf("%s/nor_gl_1k.png", materialPath),
						material.DispMap:    fmt.Sprintf("%s/disp_1k.png", materialPath),
						material.ARMMap:     fmt.Sprintf("%s/arm_1k.png", materialPath),
					},
				},
			},
		},
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
