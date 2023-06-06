package main

import (
	"flag"
	"os"

	"go-pbr/mesh"
	_ "go-pbr/mesh/obj"
)

func main() {
	meshPath := flag.String("mesh", "assets/objects/cube.obj", "")
	//materialPath := flag.String("material", "assets/materials/lgl_pbr_rusted_iron", "")
	//skymapPath := flag.String("skymap", "assets/cubemaps/castle-zavelstein-cellar", "")

	flag.Parse()

	//app := &pbr.App{}
	//if err := app.Init(); err != nil {
	//	log.Fatalln("failed to initialize app:", err)
	//}
	//defer func() {
	//	if err := app.Free(); err != nil {
	//		log.Fatalln("failed to free app:", err)
	//	}
	//}()

	meshFile, err := os.Open(*meshPath)
	if err != nil {
		panic(err)
	}

	subject, _, err := mesh.Decode(meshFile, mesh.DecodeOptions{SkipHeaderCheck: true})
	if err != nil {
		panic(err)
	}

	println(subject)

	//mesh := obj.Mesh{
	//	Path: *meshPath,
	//	Material: graphics.LitMaterial{
	//		Path: *materialPath,
	//	},
	//}
	//
	//skymap := mesh.Mesh{
	//	Path: filepath.Join(*skymapPath, "skymap.mesh"),
	//	Material: graphics.CubeMaterial{
	//		XPosPath: filepath.Join(*skymapPath, "xp.png"),
	//		XNegPath: filepath.Join(*skymapPath, "xn.png"),
	//		YPosPath: filepath.Join(*skymapPath, "yp.png"),
	//		YNegPath: filepath.Join(*skymapPath, "yn.png"),
	//		ZPosPath: filepath.Join(*skymapPath, "zp.png"),
	//		ZNegPath: filepath.Join(*skymapPath, "zn.png"),
	//	},
	//}
	//
	//camera := graphics.Camera{}
	//
	//if err := mesh.Init(); err != nil {
	//	log.Fatalln("failed to initialize mesh:", err)
	//}
	//
	//if err := skymap.Init(); err != nil {
	//	log.Fatalln("failed to initialize skymap:", err)
	//}
	//
	//if err := camera.Init(); err != nil {
	//	log.Fatalln("failed to initialize camera:", err)
	//}
	//
	//for !app.ShouldClose() {
	//	dt := app.DeltaTime()
	//	if err := mesh.Update(dt); err != nil {
	//		log.Fatalln("failed to update mesh:", err)
	//	}
	//	if err := skymap.Update(dt); err != nil {
	//		log.Fatalln("failed to update skymap:", err)
	//	}
	//	if err := camera.Update(dt); err != nil {
	//		log.Fatalln("failed to initialize camera:", err)
	//	}
	//
	//	if err := app.Renderer.Draw(mesh.Data); err != nil {
	//		log.Fatalln("failed to draw mesh:", err)
	//	}
	//	if err := app.Renderer.Draw(skymap.Data); err != nil {
	//		log.Fatalln("failed to draw skymap:", err)
	//	}
	//
	//	if err := app.PostLoop(); err != nil {
	//		log.Fatalln("failed to postloop app:", err)
	//	}
	//}
}
