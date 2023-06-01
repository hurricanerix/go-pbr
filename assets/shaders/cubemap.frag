#version 410

uniform samplerCube Skybox;

in vec3 FragUV;

out vec4 FragColor;

void main() {
    FragColor = texture(Skybox, FragUV);
}