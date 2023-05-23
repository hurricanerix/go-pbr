// pbr.frag
#version 410

uniform sampler2D texSampler;
uniform sampler2D armSampler; // R/G/B -> AO/Rough/Metal
uniform sampler2D dispSampler;
uniform sampler2D norSampler;
uniform vec3 light;

in vec2 fragTexCoord;
in vec3 lightDir;
in vec3 cameraDir;
in vec3 pos;

out vec4 outputColor;

#define AO (0)
#define ROUGHNESS (1)
#define METALIC (2)

void main() {
    vec4 tex = texture(texSampler, fragTexCoord);
    vec4 arm = texture(armSampler, fragTexCoord);
    vec4 nor = texture(norSampler, fragTexCoord);

    vec3 diffuse = texture(texSampler, fragTexCoord.st).rgb;

    vec3 ambient = vec3(0.1, 0.1, 0.1) * diffuse;
    vec3 specular = diffuse/8;

    vec3 normal = texture(norSampler, fragTexCoord.st).rgb * 2 - 1;
    float distance = length(light - pos);

    vec3 n = normalize(nor.rgb);

    float cosTheta = clamp(dot(n, lightDir), 0.0, 1.0);

    vec3 r = reflect(-lightDir, n);

    float cosAlpha = clamp(dot(cameraDir, r), 0.0, 1.0);

    vec3 lightColor = vec3(1.0, 1.0, 1.0);
    float lightPower = 5;

    outputColor = vec4(
        ambient +
        diffuse * lightColor.rgb * lightPower * cosTheta / (distance * distance) +
        specular * lightColor.rgb * lightPower * pow(int(cosAlpha), 5) / (distance * distance),
    1.0);
    outputColor = outputColor * arm[AO];
}