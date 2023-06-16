# PBR Material Viewer + Tools

* A prototype tool leveraging Golang and OpenGL for viewing PBR materials.
* A tool to pack multiple grayscale PNG files into a single PNG using different color channels.
* A tool to unpack a single PNG's color channels into multiple grayscale PNGs.


## Dev

Instructions for what is needed to build the apps for macOS/Windows/Linux.


### macOS

1. Install [HomeBrew](https://brew.sh)
2. Open Terminal and run `brew install mesa`
3. Then run `brew install glfw`


### Windows

1. Install [7zip](https://www.7-zip.org)
2. Download [x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/sjlj/x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z)
3. Extract the contents to `C:\mingw-w64`
4. Open a command prompt and enter `setx path “%PATH%;C:\mingw-w64\bin”`
5. Reboot your computer

More details here: https://medium.com/@martin.kunc/golang-compile-go-gl-example-on-windows-in-mingw-64-bfb6eb66a143


### Linux

1. Open a terminal
2. Use the package manager to install `mesa` and `glfw`
   3. Arch Based Systems: `pacman -S mesa glfw`
   4. Debian Based Systems: `apt-get install mesa glfw`


## TODO

- [ ] Triangulate Quads (OBJ)
- [ ] Emission
- [ ] Framebuffer
- [ ] Environment Mapping (Reflection)
- [ ] Gamma Correction
- [ ] HDR
- [ ] Bloom
- [ ] Deferred Shading