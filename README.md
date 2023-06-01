# PBR Material Test

A prototype tool leveraging Golang and OpenGL for viewing PBR materials.

## Command Line Arguments

| Argument       | Description                                | example    |
|----------------|--------------------------------------------|------------|
| resolution     | set the screen resolution                  | 1024x768   |
| skymap         | path to the directory of the skymap to use | ...        |
| model          | path to the model directory to use         | ...        |
| use-geometry   | enable enhanced geometry from normals/disp | true/false |
| lighting-model | change laighting model to use              | phong/pbr  |

## Controls

| Input | Description |
|-|-|
| Left Mouse Button | Rotate Model |
| Right Mouse Button | Rotate Camera |
| Middle Mouse Button | Rotate Light |
| Keyboard - | Decrease Ambient light |
| Keyboard + | Increase ambient light |
| Keyboard [ | Decrease Light Source |
| Keyboard ] | Incrase Light Source |

## TODO

- [X] Fix tangent and bitangent calculations.
- [X] Fix mouse input to rotate cube feature.
- [ ] Fix normal mapping.
- [ ] Add orbit camera around model feature.
- [ ] Add orbit light around model feature.
- [ ] Add keyboard controls to alter ambient light.
- [ ] Add keyboard controls to alter light source.
- [ ] Add PBR material shaders.
- [ ] Add toggle for phong/pbr shaders.
- [ ] Add Environment map reflections.
- [ ] Add Emission map.
- [ ] Add enhanced geometry via normal/disp feature.
- [ ] Add 'w' keyboard to toggle wireframe.
- [ ] Add 'n' keyboard to toggle normals.
- [ ] Add 1-7+0 keys to toggle
  - 1 - color map only
  - 2 - normal map only
  - 3 - displacement map only
  - 4 - ambient occlusion map only
  - 5 - roughness map only
  - 6 - metalic map only
  - 7 - emission map only
  - 0 - Full lighting
- [ ] Add i key to toggle onscreen info.
  - OpenGL version
  - GLSL version
  - Object Face Count
  - Model Rotation
  - Camera Rotation
  - Light Rotation

