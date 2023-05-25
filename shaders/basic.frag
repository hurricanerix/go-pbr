#version 410

precision mediump float;

in vec3 vertPos;
in vec3 normalInterp;

out vec4 FragColor;

const vec3 lightPos = vec3(2.0, 0.0, 2.0);
const vec3 lightColor = vec3(1.0, 1.0, 1.0);
const float lightPower = 40.0;
const vec3 ambientColor = vec3(0.7, 0.35, 0.0);
const vec3 diffuseColor = vec3(0.5, 0.0, 0.0);
const vec3 specColor = vec3(1.0, 1.0, 1.0);
const float shininess = 32.0;
const float screenGamma = 1; // Assume the monitor is calibrated to the sRGB color space
const int mode = 2;

void main() {
    vec3 normal = normalize(normalInterp);
    vec3 lightDir = lightPos - vertPos;
    float distance = length(lightDir);
    distance = distance * distance;
    lightDir = normalize(lightDir);

    float lambertian = max(dot(lightDir, normal), 0.0);
    float specular = 0.0;

    if (lambertian > 0.0) {

        vec3 viewDir = normalize(-vertPos);

        // this is blinn phong
        vec3 halfDir = normalize(lightDir + viewDir);
        float specAngle = max(dot(halfDir, normal), 0.0);
        specular = pow(specAngle, shininess);

        // this is phong (for comparison)
        if (mode == 2) {
            vec3 reflectDir = reflect(-lightDir, normal);
            specAngle = max(dot(reflectDir, viewDir), 0.0);
            // note that the exponent is different here
            specular = pow(specAngle, shininess/4.0);
        }
    }
    vec3 colorLinear = ambientColor +
    diffuseColor * lambertian * lightColor * lightPower / distance +
    specColor * specular * lightColor * lightPower / distance;
    // apply gamma correction (assume ambientColor, diffuseColor and specColor
    // have been linearized, i.e. have no gamma correction in them)
    vec3 colorGammaCorrected = pow(colorLinear, vec3(1.0 / screenGamma));
    // use the gamma corrected color in the fragment
    FragColor = vec4(colorGammaCorrected, 1.0);
}
