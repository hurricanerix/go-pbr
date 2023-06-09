#version 410

uniform sampler2D diffuseMap;

uniform vec3 lightPos;
uniform vec3 viewPos;
uniform vec3 LightColor;

in VS_OUT {
    vec3 FragPos;
    vec2 TexCoords;
    vec3 Normal;
} fs_in;

out vec4 FragColor;

void main() {
    float alpha = texture(diffuseMap, fs_in.TexCoords).a;
    vec3 color = texture(diffuseMap, fs_in.TexCoords).rgb;

    // ambient
    float ambientStrength = 0.1;
    vec3 ambient = ambientStrength * LightColor;

    // diffuse
    vec3 norm = normalize(fs_in.Normal);
    vec3 lightDir = normalize(lightPos - fs_in.FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * LightColor;

    // specular
    float specularStrength = 0.5;
    vec3 viewDir = normalize(viewPos - fs_in.FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
    vec3 specular = specularStrength * spec * LightColor;

    vec3 result = (ambient + diffuse + specular) * color;
    FragColor = vec4(result, alpha);
}