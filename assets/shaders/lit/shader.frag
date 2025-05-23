#version 410

uniform sampler2D diffuseMap;
uniform sampler2D normalMap; // TODO: put disp into alpha of normal.
uniform sampler2D dispMap;
uniform sampler2D armMap;

//uniform float AmbientStrength;
//uniform vec3 AmbientColor;
uniform vec3 lightPos;
uniform vec3 LightColor;
uniform float LightPower;
uniform mat4 Model;
uniform vec3 viewPos;

in VS_OUT {
    vec3 FragPos;
    vec2 TexCoords;
    vec3 TangentLightPos;
    vec3 TangentViewPos;
    vec3 TangentFragPos;
} fs_in;

out vec4 FragColor;

const float specularStrength = 0.5;

void main() {
    vec3 normal = texture(normalMap, fs_in.TexCoords).rgb;
    normal = normalize(normal * 2.0 - 1.0);

    vec3 color = texture(diffuseMap, fs_in.TexCoords).rgb;

    // Calculate Ambient Component
    vec3 ambient = 0.1 * color;

    // Calculate Diffuse Component
    vec3 lightDir = normalize(fs_in.TangentLightPos - fs_in.TangentFragPos);
    float diff = max(dot(lightDir, normal), 0.0);
    vec3 diffuse = diff * color;

    // Calculate Specular Component
    vec3 viewDir = normalize(fs_in.TangentViewPos - fs_in.TangentFragPos);
    vec3 reflectDir = reflect(-lightDir, normal);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), 32.0);
    vec3 specular = vec3(0.2) * spec;

    // Calculate Final Color
    FragColor = vec4(clamp(ambient + diffuse + specular, 0.0, 1.0), 1.0);
    FragColor = vec4(diffuse, 1.0);
}