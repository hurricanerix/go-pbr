#version 410

uniform mat4 projection;
uniform mat4 view;

layout (location = 0) in vec3 aPos;

out vec3 FragUV;

void main() {
    gl_Position = projection * view * vec4(aPos, 1.0);
    FragUV = aPos;
}