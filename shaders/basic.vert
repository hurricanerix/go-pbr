#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
//uniform mat4 normalMat;

in vec3 vert;
//in vec3 normal;

out vec3 vertPos;
out vec3 normalInterp;

void main() {

    vec3 normal = vert;
    mat4 modelview = view * model;
    gl_Position = projection * modelview * vec4(vert, 1.0);
    mat4 normalMat = transpose(inverse(modelview));
    vec4 vertPos4 = modelview * vec4(vert, 1.0);
    vertPos = vec3(vertPos4) / vertPos4.w;
    normalInterp = vec3(normalMat * vec4(normal, 0.0));
}
