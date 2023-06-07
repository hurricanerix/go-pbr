package scene

import (
	"go-pbr/graphics"
)

type Object struct {
	Transform graphics.Transform

	Renderer graphics.Renderer
}

func (o *Object) Attach() error {
	if err := o.Renderer.Attach(); err != nil {
		return err
	}
	return nil
}

func (o *Object) Bind() error {
	//if err := o.Renderer.LoadMaterial(o.Material); err != nil {
	//	return err
	//}
	if err := o.Renderer.Bind(); err != nil {
		return err
	}
	return nil
}

func (o *Object) Update(dt float32) error {
	if err := o.Renderer.SetTransform(o.Transform); err != nil {
		return err
	}
	return nil
}

func (o *Object) Draw() error {

	if err := o.Renderer.Draw(); err != nil {
		return err
	}
	return nil
}
