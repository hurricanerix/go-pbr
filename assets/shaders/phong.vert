#version 410

uniform mat4 Projection;
uniform mat4 View;
uniform mat4 Model;
uniform vec3 LightPos;
uniform vec3 ViewPos;

layout (location = 0) in vec3 Vert; // aPos
layout (location = 1) in vec2 UV;
layout (location = 2) in vec3 aNormal;
layout (location = 3) in vec3 aTangent;
layout (location = 4) in vec3 aBitangent;

//out vec3 FragPos;
//out vec2 FragUV;
//out vec3 Normal;

out VS_OUT {
    vec3 FragPos;
    vec2 FragUV;
    vec3 TangentLightPos;
    vec3 TangentViewPos;
    vec3 TangentFragPos;
} vs_out;

void main() {
    gl_Position = Projection * View * Model * vec4(Vert, 1.0);

    vs_out.FragPos = vec3(Model * vec4(Vert, 1.0));
    vs_out.FragUV = UV;

    mat3 normalMatrix = transpose(inverse(mat3(Model)));
    vec3 T = normalize(normalMatrix * aTangent);
    vec3 N = normalize(normalMatrix * aNormal);
    T = normalize(T - dot(T, N) * N);
    vec3 B = cross(N, T);

    mat3 TBN = transpose(mat3(T, B, N));
    vs_out.TangentLightPos = TBN * LightPos;
    vs_out.TangentViewPos  = TBN * ViewPos;
    vs_out.TangentFragPos  = TBN * vs_out.FragPos;
}