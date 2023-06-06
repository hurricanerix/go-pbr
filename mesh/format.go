// Copyright 2010 The Go Authors. All rights reserved.
// Copyright 2023 Richard Hawkins
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file borrows code from image/format.go in the image package.

package mesh

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

// ErrFormat indicates that decoding encountered an unknown format.
var ErrFormat = errors.New("mesh: unknown format")

// A format holds an image format's name, magic header and how to decode it.
type format struct {
	name, magic string
	decode      func(io.Reader, DecodeOptions) (Mesh, error)
}

// Formats is the list of registered formats.
var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

// RegisterFormat registers an geometry format for use by Decode.
// Name is the name of the format, like "obj".
// Magic is the magic prefix that identifies the format's encoding. The magic
// string can contain "?" wildcards that each match any one byte.
// Decode is the function that decodes the encoded mesh.
func RegisterFormat(name, magic string, decode func(io.Reader, DecodeOptions) (Mesh, error)) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{name, magic, decode}))
	formatsMu.Unlock()
}

// A reader is an io.Reader that can also peek ahead.
type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

// asReader converts an io.Reader to a reader.
func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

// match reports whether magic matches b. Magic may contain "?" wildcards.
func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

// sniff determines the format of r's data.
func sniff(r reader) format {
	formats, _ := atomicFormats.Load().([]format)
	for _, f := range formats {
		b, err := r.Peek(len(f.magic))
		if err == nil && match(f.magic, b) {
			return f
		}
	}
	return format{}
}

type DecodeOptions struct {
	SkipHeaderCheck            bool
	CalculateTangentsIfMissing bool
}

// Decode decodes an image that has been encoded in a registered format.
// The string returned is the format name used during format registration.
// Format registration is typically done by an init function in the codec-
// specific package.
func Decode(r io.Reader, options DecodeOptions) (Mesh, string, error) {
	rr := asReader(r)
	f := sniff(rr)
	if f.decode == nil {
		return nil, "", ErrFormat
	}
	m, err := f.decode(rr, options)
	return m, f.name, err
}
