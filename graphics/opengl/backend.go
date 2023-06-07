package opengl

//
//import (
//	"errors"
//	"go-pbr/mesh"
//	"sync"
//	"sync/atomic"
//)
//
//// ErrBackend indicates that an unknown error with the backend.
//var ErrBackend = errors.New("backend: unknown error")
//
//type Backend interface {
//	Render(mesh.Mesh) error
//}
//
//// backend holds an graphics API's name, how to render with it.
//type backend struct {
//	name    string
//	backend Backend
//}
//
//// Backends is the list of registered backends.
//var (
//	backendsMu     sync.Mutex
//	atomicBackends atomic.Value
//)
//
//func RegisterBackend(name string, b Backend) {
//	backendsMu.Lock()
//	backends, _ := atomicBackends.Load().([]backend)
//	atomicBackends.Store(append(backends, backend{name, b}))
//	backendsMu.Unlock()
//}
