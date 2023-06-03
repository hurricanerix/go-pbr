#version 410

uniform sampler2D diffuseMap;
uniform sampler2D normalMap; // TODO: put disp into alpha of normal.
uniform sampler2D dispMap;
uniform sampler2D armMap;

uniform float AmbientStrength;
uniform vec3 AmbientColor;
uniform vec3 lightPos;
uniform vec3 LightColor;
uniform float LightPower;
uniform mat4 Model;
uniform vec3 viewPos;
uniform int enableDisplacement; // move to bit mask
uniform int discardOutOfBounds;

uniform float heightScale;

in VS_OUT {
    vec3 FragPos;
    vec2 FragUV;
    vec3 TangentLightPos;
    vec3 TangentViewPos;
    vec3 TangentFragPos;
} fs_in;

out vec4 FragColor;

const float specularStrength = 0.5;


vec2 ParallaxMapping(vec2 texCoords, vec3 viewDir) {
    float height = texture(dispMap, texCoords).r;
    vec2 p = viewDir.xy / viewDir.z * (height * heightScale);
    return texCoords - p;
}

void main() {
    vec2 texCoords = fs_in.FragUV;
    vec3 viewDir = normalize(fs_in.TangentViewPos - fs_in.TangentFragPos);

    if (enableDisplacement == 1 ) {
        // TODO: possibly replace with steep parallax mapping.
        texCoords = ParallaxMapping(texCoords, viewDir);

        while(currentLayerDepth < currentDepthMapValue)
        {
            // shift texture coordinates along direction of P
            currentTexCoords -= deltaTexCoords;
            // get depthmap value at current texture coordinates
            currentDepthMapValue = texture(depthMap, currentTexCoords).r;
            // get depth of next layer
            currentLayerDepth += layerDepth;
        }

        if (discardOutOfBounds == 1 &&
        (texCoords.x < 0.0 || texCoords.y < 0.0 ||
        texCoords.x > 1.0 || texCoords.y > 1.0)) {
            discard;
        }
    }

    float alpha = texture(diffuseMap, texCoords).a;
    vec3 color = texture(diffuseMap, texCoords).rgb;
    vec3 normal = texture(normalMap, texCoords).rgb;
    normal = normalize(normal * 2.0 - 1.0);

    // Calculate Ambient Component
    vec3 ambient = AmbientStrength * color;

    // Calculate Diffuse Component
    vec3 lightDir = normalize(fs_in.TangentLightPos - fs_in.TangentFragPos);
    float diff = max(dot(lightDir, normal), 0.0);
    vec3 diffuse = diff * color;

    // Calculate Specular Component
//    vec3 viewDir = normalize(fs_in.TangentViewPos - fs_in.TangentFragPos);
    vec3 reflectDir = reflect(-lightDir, normal);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), 32.0);
    vec3 specular = vec3(0.2) * spec;

    // Calculate Final Color
    FragColor = vec4(clamp(ambient + diffuse + specular, 0.0, 1.0), alpha);
    //    FragColor = vec4(normal, 1.0);
}