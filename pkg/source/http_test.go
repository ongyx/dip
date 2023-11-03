package source

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type testHandler struct{}

func (th *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Disposition", "inline")
	h.Set("Content-Type", "text/plain")

	io.WriteString(w, content)
}

func TestHTTP(t *testing.T) {
	ts := httptest.NewServer(&testHandler{})
	defer ts.Close()

	us, err := url.Parse(ts.URL)
	if err != nil {
		t.Error("could not parse test server URL:", err)
	}

	src, err := NewHTTP(us)
	if err != nil {
		t.Error("failed to initalize http source:", err)
	}

	compare(t, src, Root, content)
}
