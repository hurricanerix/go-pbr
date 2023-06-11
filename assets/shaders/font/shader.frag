#version 410 core
in vec2 TexCoords;
out vec4 color;

uniform sampler2D text;
uniform vec3 textColor;

void main()
{
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(text, TexCoords).r);
    color = vec4(textColor, 1.0) * sampled;
//    if (color.a > 0.5) color = vec4(1.0, 0.0, 0.0, 1.0);
//    else color = vec4(0.0, 1.0, 0.0, 1.0);
//    color = vec4(sampled.a, sampled.a, sampled.a, 1.0);
}