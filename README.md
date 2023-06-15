# PBR Material Test

A prototype tool leveraging Golang and OpenGL for viewing PBR materials.

## TODO

- [ ] Triangulate Quads (OBJ)
- [ ] Emission
- [ ] Framebuffer
- [ ] Environment Mapping (Reflection)
- [ ] Gamma Correction
- [ ] HDR
- [ ] Bloom
- [ ] Deferred Shading

## Dev

* https://www.khronos.org/opengl/wiki/Platform_specifics:_Windows#Installing_Mesa3D_on_Windows
* https://www.glfw.org

### macOS

brew install mesa
brew install glfw

### Windows 11

* https://www.khronos.org/opengl/wiki/Platform_specifics:_Windows#Installing_Mesa3D_on_Windows
* https://www.glfw.org

### Linux

pacman -S mesa glfw

### Benchmarks


Intel(R) Core(TM) i5-3427U CPU @ 1.80GHz
MemTotal:        8041484 kB
OpenGL 4.2 (Core Profile) Mesa 23.0.4
GLSL 4.20


No Info
```
GOROOT=/home/rhawkins/bin/go #gosetup
GOPATH=/home/rhawkins/bin/gop #gosetup
/home/rhawkins/bin/go/bin/go build -o /home/rhawkins/.cache/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer . #gosetup
/home/rhawkins/.cache/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer
%!v(PANIC=String method: runtime error: index out of range [23040] with length 23040)
start
        Min: 999999.000000
        Avg: 0.000000
        Max: 0.000000
endProcessInput
        Min: 0.000002
        Avg: 0.000007
        Max: 0.000170
endUpdate
        Min: 0.000004
        Avg: 0.000014
        Max: 0.000178
endRender
        Min: 0.000074
        Avg: 0.015718
        Max: 0.017842
endShowLight
        Min: 0.000084
        Avg: 0.015888
        Max: 0.018184
endShowInfo
        Min: 0.000085
        Avg: 0.015894
        Max: 0.018194
endSwap
        Min: 0.001207
        Avg: 0.016514
        Max: 0.018889
endPoll
        Min: 0.001225
        Avg: 0.016579
        Max: 0.018951
end
        Min: 0.001225
        Avg: 0.016580
        Max: 0.018953

Process finished with the exit code 0



```

Show Info
```
GOROOT=/home/rhawkins/bin/go #gosetup
GOPATH=/home/rhawkins/bin/gop #gosetup
/home/rhawkins/bin/go/bin/go build -o /home/rhawkins/.cache/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer . #gosetup
/home/rhawkins/.cache/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer
%!v(PANIC=String method: runtime error: index out of range [23040] with length 23040)
start
        Min: 999999.000000
        Avg: 0.000000
        Max: 0.000000
endProcessInput
        Min: 0.000002
        Avg: 0.000006
        Max: 0.000124
endUpdate
        Min: 0.000004
        Avg: 0.000015
        Max: 0.000241
endRender
        Min: 0.000151
        Avg: 0.011687
        Max: 0.016025
endShowLight
        Min: 0.000164
        Avg: 0.011796
        Max: 0.017055
endShowInfo
        Min: 0.001996
        Avg: 0.015074
        Max: 0.020937
endSwap
        Min: 0.002797
        Avg: 0.016634
        Max: 0.022726
endPoll
        Min: 0.002815
        Avg: 0.016696
        Max: 0.022761
end
        Min: 0.002815
        Avg: 0.016697
        Max: 0.022762

Process finished with the exit code 0
```



-----

no info
```
GOROOT=/Users/rhawkins/bin/go1.20.4 #gosetup
GOPATH=/Users/rhawkins/bin/go #gosetup
/Users/rhawkins/bin/go1.20.4/bin/go build -o /Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer . #gosetup
/Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer
%!v(PANIC=String method: runtime error: index out of range [23040] with length 23040)
start
        Min: 999999.000000
        Avg: 0.000000
        Max: 0.000000
endProcessInput
        Min: 0.000000
        Avg: 0.000003
        Max: 0.000099
endUpdate
        Min: 0.000001
        Avg: 0.000006
        Max: 0.000103
endRender
        Min: 0.000041
        Avg: 0.000178
        Max: 0.001328
endShowLight
        Min: 0.000046
        Avg: 0.000200
        Max: 0.001420
endShowInfo
        Min: 0.000046
        Avg: 0.000201
        Max: 0.001420
endSwap
        Min: 0.000400
        Avg: 0.008191
        Max: 0.016720
endPoll
        Min: 0.000622
        Avg: 0.008455
        Max: 0.049624
end
        Min: 0.000622
        Avg: 0.008455
        Max: 0.049625

Process finished with the exit code 0



```


Info
```
GOROOT=/Users/rhawkins/bin/go1.20.4 #gosetup
GOPATH=/Users/rhawkins/bin/go #gosetup
/Users/rhawkins/bin/go1.20.4/bin/go build -o /Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer . #gosetup
/Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer
%!v(PANIC=String method: runtime error: index out of range [23040] with length 23040)
start
Min: 999999.000000
Avg: 0.000000
Max: 0.000000
endProcessInput
Min: 0.000000
Avg: 0.000001
Max: 0.000046
endUpdate
Min: 0.000001
Avg: 0.000003
Max: 0.000048
endRender
Min: 0.000031
Avg: 0.000060
Max: 0.000383
endShowLight
Min: 0.000037
Avg: 0.000071
Max: 0.000466
endShowInfo
Min: 0.034618
Avg: 0.055476
Max: 0.064545
endSwap
Min: 0.035039
Avg: 0.056032
Max: 0.066686
endPoll
Min: 0.035178
Avg: 0.056178
Max: 0.100311
end
Min: 0.035178
Avg: 0.056178
Max: 0.100312

Process finished with the exit code 0


```

GOROOT=/Users/rhawkins/bin/go1.20.4 #gosetup
GOPATH=/Users/rhawkins/bin/go #gosetup
/Users/rhawkins/bin/go1.20.4/bin/go build -o /Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer . #gosetup
/Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer
%!v(PANIC=String method: runtime error: index out of range [23040] with length 23040)

Process finished with the exit code 0


```
GOROOT=/Users/rhawkins/bin/go1.20.4 #gosetup
GOPATH=/Users/rhawkins/bin/go #gosetup
/Users/rhawkins/bin/go1.20.4/bin/go build -o /Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer . #gosetup
/Users/rhawkins/Library/Caches/JetBrains/GoLand2023.1/tmp/GoLand/___pbr_viewer
%!v(PANIC=String method: runtime error: index out of range [23040] with length 23040)
start
        Min: 999999.000000
        Avg: 0.000000
        Max: 0.000000
endProcessInput
        Min: 0.000000
        Avg: 0.000001
        Max: 0.000046
endUpdate
        Min: 0.000001
        Avg: 0.000003
        Max: 0.000048
endRender
        Min: 0.000031
        Avg: 0.000060
        Max: 0.000383
endShowLight
        Min: 0.000037
        Avg: 0.000071
        Max: 0.000466
endShowInfo
        Min: 0.034618
        Avg: 0.055476
        Max: 0.064545
endSwap
        Min: 0.035039
        Avg: 0.056032
        Max: 0.066686
endPoll
        Min: 0.035178
        Avg: 0.056178
        Max: 0.100311
end
        Min: 0.035178
        Avg: 0.056178
        Max: 0.100312

Process finished with the exit code 0



```