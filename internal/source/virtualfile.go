package source

import (
	"bytes"
	"io/fs"
	"time"
)

// VirtualFile is a fs.File implementation that wraps a buffer.
type VirtualFile struct {
	bytes.Buffer

	name    string
	modTime time.Time
}

// NewVirtualFile creates a new virtual file with the given name.
func NewVirtualFile(name string) *VirtualFile {
	return &VirtualFile{name: name}
}

func (f *VirtualFile) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *VirtualFile) Close() error {
	return nil
}

func (f *VirtualFile) Name() string {
	return f.name
}

func (f *VirtualFile) Size() int64 {
	return int64(f.Buffer.Len())
}

func (f *VirtualFile) Mode() fs.FileMode {
	return fs.ModePerm
}

func (f *VirtualFile) ModTime() time.Time {
	return f.modTime
}

func (f *VirtualFile) IsDir() bool {
	return false
}

func (f *VirtualFile) Sys() any {
	return nil
}
