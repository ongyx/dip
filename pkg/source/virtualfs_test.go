package source

import (
	"io"
	"testing"
)

const (
	content = "# Hello World!"
)

func TestVirtualFS(t *testing.T) {
	vf := NewVirtualFile("test")
	vf.WriteString(content)

	vfs := NewVirtualFS(vf)

	f, err := vfs.Open(Root)
	if err != nil {
		t.Error("couldn't open root")
	}

	b, err := io.ReadAll(f)
	if err != nil {
		t.Error("couldn't read from vfile")
	}

	readContent := string(b)

	if readContent != content {
		t.Errorf("expected %s, got %s\n", content, readContent)
	}
}
