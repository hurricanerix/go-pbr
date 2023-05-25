#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

in vec3 vert;

void main() {
    gl_Position = projection * view * model * vec4(vert, 1);
}
