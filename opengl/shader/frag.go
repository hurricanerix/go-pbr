package shader

var Frag = `
#version 410

uniform sampler2D DiffuseSampler;
uniform sampler2D ArmSampler;
uniform sampler2D NormalSampler; // TODO: put disp into alpha of normal.
uniform sampler2D DispSampler;

uniform vec3 LightPos;
uniform vec3 AmbientColor;
uniform vec3 LightColor;
uniform float LightPower;

in vec3 VertPos;
in vec2 FragUV;
in vec3 NormalDir;
in vec3 LightDir;
in vec3 EyeDir;

out vec4 FragColor;

void main() {
	float alpha = texture(DiffuseSampler, FragUV.st).a;
	vec3 diffuse = texture(DiffuseSampler, FragUV.st).rgb;

	vec3 ambient = AmbientColor.rgb * diffuse;
	vec3 specular = diffuse/8;

	vec3 normal = texture(NormalSampler, FragUV.st).rgb * 2 - 1;
	float distance = length(LightPos - VertPos);

	vec3 n = normalize(NormalDir);
	vec3 l = normalize(LightDir);

	float cosTheta = clamp(dot(n, l), 0.0, 1.0);

	vec3 e = normalize(EyeDir);
	vec3 r = reflect(-l, n);

	float cosAlpha = clamp(dot(e, r), 0.0, 1.0);

	FragColor = vec4(
    	ambient +
    	diffuse * LightColor.rgb * LightPower * cosTheta /
			(distance * distance) +
		specular * LightColor.rgb * LightPower * pow(cosAlpha, 5) /
			(distance * distance), alpha);

}` + "\x00"
