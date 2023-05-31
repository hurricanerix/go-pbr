package shader

var VertPhong = `#version 410

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
}` + "\x00"

var FragPhong = `#version 410

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
	// obtain normal from normal map in range [0,1]
    vec3 normal = texture(NormalSampler, fs_in.FragUV).rgb;
    // transform normal vector to range [-1,1]
    normal = normalize(normal * 2.0 - 1.0);  // this normal is in tangent space
   
    // get diffuse color
    vec3 color = texture(DiffuseSampler, fs_in.FragUV).rgb;
    // ambient
    vec3 ambient = 0.1 * color;
    // diffuse
    vec3 lightDir = normalize(fs_in.TangentLightPos - fs_in.TangentFragPos);
    float diff = max(dot(lightDir, normal), 0.0);
    vec3 diffuse = diff * color;
    // specular
    vec3 viewDir = normalize(fs_in.TangentViewPos - fs_in.TangentFragPos);
    vec3 reflectDir = reflect(-lightDir, normal);
    vec3 halfwayDir = normalize(lightDir + viewDir);  
    float spec = pow(max(dot(normal, halfwayDir), 0.0), 32.0);

    vec3 specular = vec3(0.2) * spec;
    FragColor = vec4(ambient + diffuse + specular, 1.0);
	FragColor = vec4(color, 1.0);
}` + "\x00"
