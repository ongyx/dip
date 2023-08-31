package source

import (
	"io/fs"
	"time"
)

// VirtualFileInfo is the description of a virtual file.
type VirtualFileInfo struct {
	name    string
	size    int64
	modtime time.Time
}

func (vfi *VirtualFileInfo) Name() string {
	return vfi.name
}

func (vfi *VirtualFileInfo) Size() int64 {
	return vfi.size
}

func (vfi *VirtualFileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

func (vfi *VirtualFileInfo) ModTime() time.Time {
	return vfi.modtime
}

func (vfi *VirtualFileInfo) IsDir() bool {
	return false
}

func (vfi *VirtualFileInfo) Sys() any {
	return nil
}
