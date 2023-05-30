package shader

var Frag = `#version 410

uniform sampler2D DiffuseSampler;
uniform sampler2D ArmSampler;
uniform sampler2D NormalSampler; // TODO: put disp into alpha of normal.
uniform sampler2D DispSampler;

uniform float AmbientStrength;
uniform vec3 AmbientColor;
uniform vec3 LightPos;
uniform vec3 LightColor;
uniform float LightPower;
uniform mat4 ModelMatrix;
uniform vec3 ViewPos;

in vec3 FragPos;
in vec2 FragUV;
in vec3 Normal;
//in vec3 LightDir;
in vec3 EyeDir;

out vec4 FragColor;

const float specularStrength = 0.5;

void main() {
	float alpha = texture(DiffuseSampler, FragUV.st).a;
	vec3 objectColor = vec3(0.5, 0.5, 0.5); //texture(DiffuseSampler, FragUV.st).rgb;
	vec3 fragNormal = texture(NormalSampler, FragUV.st).rgb;

	// Calculate Ambient Component
	vec3 ambient = AmbientStrength * AmbientColor;
	
	// Calculate Diffuse Component
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(LightPos - FragPos);
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * LightColor;

	// Calculate Specular Component
	vec3 viewDir = normalize(ViewPos - FragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
	vec3 specular = specularStrength * spec * LightColor;

	// Calculate Final Color
	FragColor = vec4((ambient + diffuse + specular) * objectColor, alpha);
}` + "\x00"
