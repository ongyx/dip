package document

import (
	"bytes"
	"io/fs"
	"time"
)

// VirtualFile is a fs.File implementation that wraps a buffer.
type VirtualFile struct {
	bytes.Buffer

	name    string
	modtime time.Time
}

// NewVirtualFile creates a new virtual file with the given name.
func NewVirtualFile(name string) *VirtualFile {
	return &VirtualFile{name: name}
}

func (vf *VirtualFile) Stat() (fs.FileInfo, error) {
	return &VirtualFileInfo{
		name:    vf.name,
		size:    int64(vf.Buffer.Len()),
		modtime: vf.modtime,
	}, nil
}

func (vf *VirtualFile) Close() error {
	return nil
}
