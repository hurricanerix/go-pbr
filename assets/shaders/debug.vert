#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

layout (location = 0) in vec3 aPos;

void main() {
    gl_Position = projection * view * model * vec4(aPos, 1.0);
}