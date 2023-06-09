#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform vec3 lightPos;
uniform vec3 viewPos;

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoords;
//layout (location = 3) in vec3 aTangent;
//layout (location = 4) in vec3 aBitangent;

out VS_OUT {
    vec2 FragUV;
    float Diffuse;
} vs_out;

void main() {
    // TODO: Fix, not working
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    vs_out.FragUV = aTexCoords;

    vec3 modelViewVertex = vec3(view * model * vec4(aPos, 1.0));
    vec3 modelViewNormal = vec3(view * model * vec4(aNormal, 0.0));
    vec3 lightEye = (view * vec4(lightPos, 0.0)).xyz;
    float distance = length(lightEye - modelViewVertex);
    vec3 lightVector = normalize(lightEye - modelViewVertex);
    float diffuse = max(dot(modelViewNormal, lightVector), 1.0);
    diffuse = diffuse * (1.0 / (1.0 + (0.25 * distance * distance)));
    vs_out.Diffuse = diffuse;
}