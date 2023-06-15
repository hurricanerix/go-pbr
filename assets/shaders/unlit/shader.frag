#version 410

uniform sampler2D diffuseMap;
uniform vec3 color;
uniform int noMap;

in VS_OUT {
    vec2 TexCoords;
} fs_in;

out vec4 FragColor;

void main() {
    if (noMap == 1) {
        FragColor = vec4(normalize(color), 1.0);
    } else {
        vec3 diffuse = texture(diffuseMap, fs_in.TexCoords).rgb;
        vec3 result = clamp(color * diffuse, 0.0, 1.0);
        FragColor = vec4(result, 1.0);
    }
}