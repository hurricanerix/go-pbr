#version 410

layout(location=0) in vec3 vColor;
layout (location=0) out vec4 FragColor;

void main() {
    FragColor = vec4(0.0, 1.0, 0.0, 1.0); // vec4(vColor, 1.0);
}