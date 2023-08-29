package document

import (
	"bytes"
	"io/fs"
	"time"
)

// VirtualFile is a fs.File implementation that wraps a byte slice.
type VirtualFile struct {
	buf *bytes.Buffer
	vfi *VirtualFileInfo
}

// NewVirtualFile creates a new virtual file with the given byte slice and name.
// The byte slice must not be modified afterwards.
func NewVirtualFile(data []byte, name string) *VirtualFile {
	return &VirtualFile{
		buf: bytes.NewBuffer(data),
		vfi: &VirtualFileInfo{
			name:    name,
			size:    int64(len(data)),
			modtime: time.Now(),
		},
	}
}

func (f *VirtualFile) Stat() (fs.FileInfo, error) {
	return f.vfi, nil
}

func (f *VirtualFile) Read(p []byte) (int, error) {
	return f.buf.Read(p)
}

func (f *VirtualFile) Close() error {
	return nil
}
