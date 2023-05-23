// pbr.vert
#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform vec3 camera;
uniform vec3 light;

in vec3 vert;
in vec3 normal;
in vec2 vertTexCoord;

out vec2 fragTexCoord;
out vec3 pos;
out vec3 lightDir;
out vec3 cameraDir;

vec3 tangent(vec3 n) {
    vec3 t = vec3(0.0, 0.0, 0.0);
    vec3 c1 = cross(n, vec3(0.0, 0.0, 1.0));
    vec3 c2 = cross(n, vec3(0.0, 1.0, 0.0));
    if(length(c1) > length(c2)) {
        t = c1;
    } else {
        t = c2;
    }
    normalize(t);
    return t;
}

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * view * model * vec4(vert, 1);
    pos = (model * vec4(vert, 1)).xyz;
    lightDir = normalize(light - pos);
    cameraDir = normalize(camera - pos);

}
