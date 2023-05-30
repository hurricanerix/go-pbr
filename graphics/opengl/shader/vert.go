package shader

var Vert = `#version 410

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;
uniform vec3 LightPos;

layout (location = 0) in vec3 Vert; // aPos
layout (location = 1) in vec2 UV;
layout (location = 2) in vec3 aNormal;

out vec3 FragPos;
out vec2 FragUV;
out vec3 Normal;

void main() {
	gl_Position = ProjMatrix * ViewMatrix * ModelMatrix * vec4(Vert, 1.0);
	FragUV = UV;
	FragPos = vec3(ModelMatrix * vec4(Vert, 1.0));
	// TODO: inverse operation is complicated, calculate this on the CPU
	// and pass it to the shader.
	Normal = mat3(transpose(inverse(ModelMatrix))) * aNormal;
}` + "\x00"
