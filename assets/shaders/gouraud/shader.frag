#version 410

uniform sampler2D diffuseMap;
uniform sampler2D normalMap; // TODO: put disp into alpha of normal.
uniform sampler2D dispMap;
uniform sampler2D armMap;

uniform float AmbientStrength;
uniform vec3 AmbientColor;
uniform vec3 lightPos;
uniform vec3 LightColor;
uniform float LightPower;
uniform mat4 Model;
uniform vec3 viewPos;

in VS_OUT {
    vec2 FragUV;
    float Diffuse;
} fs_in;

out vec4 FragColor;

void main() {
    float alpha = texture(diffuseMap, fs_in.FragUV).a;
    vec3 color = texture(diffuseMap, fs_in.FragUV).rgb;
    FragColor = vec4(color * fs_in.Diffuse, alpha);
}