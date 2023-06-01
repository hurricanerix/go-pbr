#version 410

uniform sampler2D DiffuseSampler;
uniform sampler2D ArmSampler;
uniform sampler2D NormalSampler; // TODO: put disp into alpha of normal.
uniform sampler2D DispSampler;

uniform float AmbientStrength;
uniform vec3 AmbientColor;
uniform vec3 LightPos;
uniform vec3 LightColor;
uniform float LightPower;
uniform mat4 Model;
uniform vec3 ViewPos;

in VS_OUT {
    vec3 FragPos;
    vec2 FragUV;
    vec3 TangentLightPos;
    vec3 TangentViewPos;
    vec3 TangentFragPos;
} fs_in;

out vec4 FragColor;

const float specularStrength = 0.5;

void main() {
    float alpha = texture(DiffuseSampler, fs_in.FragUV).a;
    vec3 color = texture(DiffuseSampler, fs_in.FragUV).rgb;
    vec3 normal = texture(NormalSampler, fs_in.FragUV).rgb;
    normal = normalize(normal * 2.0 - 1.0);

    // Calculate Ambient Component
    vec3 ambient = AmbientStrength * AmbientColor;

    // Calculate Diffuse Component
    vec3 lightDir = normalize(fs_in.TangentLightPos - fs_in.TangentFragPos);
    float diff = max(dot(lightDir, normal), 0.0);
    vec3 diffuse = diff * LightColor;

    // Calculate Specular Component
    vec3 viewDir = normalize(fs_in.TangentViewPos - fs_in.TangentFragPos);
    vec3 reflectDir = reflect(-lightDir, normal);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), 32.0);
    vec3 specular = LightColor * spec;

    // Calculate Final Color
    FragColor = vec4((ambient + diffuse + specular) * color, alpha);
}