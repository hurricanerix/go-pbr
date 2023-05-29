package shader

var Vert = `#version 410

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;
uniform vec3 LightPos;
//uniform mat4 NormalMatrix;

in vec3 Vert;
in vec2 UV;
in vec3 Normal;

out vec3 VertPos;
out vec2 FragUV;
out vec3 NormalDir;
out vec3 LightDir;
out vec3 EyeDir;

void main() {
    // Set Position
	mat4 modelView = ViewMatrix * ModelMatrix;
	gl_Position = ProjMatrix * modelView * vec4(Vert, 1.0);

	// Set VertPos
	vec4 vertPos4 = modelView * vec4(Vert, 1.0);
	VertPos = vec3(vertPos4) / vertPos4.w;
	
	// Set FragUV
	FragUV = UV;

	// Set NormalDir
	mat4 normalMatrix = transpose(inverse(modelView));
    NormalDir = vec3(normalMatrix * vec4(Normal, 0.0));

	// Set LightDir
	//vec3 LightPos = vec3(4, 4, 4);
	LightDir = vec3(1.0, 1.0, 0.0);
	//vec3 MCTangent = tangent(MCNormal);
	//mat3 mv3Matrix = mat3x3(mvMatrix);
	//vec3 n = normalize(mv3Matrix * MCNormal);
	//vec3 t = normalize(mv3Matrix * MCTangent);
	//vec3 b = normalize(mv3Matrix * cross(n, t));
	//LightDir = vec3(ViewMatrix * vec4(LightPos, 0.0)) - vec3(ccVertex);
	//vec3 v;
	//v.x = dot(LightDir, t);
	//v.y = dot(LightDir, b);
	//v.z = dot(LightDir, n);
	//LightDir = v;
	
	// Set EyeDir
	EyeDir = vec3(3, 3, 3) - VertPos;
}` + "\x00"
