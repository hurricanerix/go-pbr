#version 410

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;

layout (location = 0) in vec3 Vert; // aPos

out vec3 FragUV;

void main() {
    gl_Position = ProjMatrix * ViewMatrix * vec4(Vert, 1.0);
    FragUV = Vert;
}