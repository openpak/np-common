package capture

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCaptureMasksAndKeepsBody(t *testing.T) {
	dir := t.TempDir()
	r, _ := New(dir)
	var seen string
	h := r.Wrap(http.HandlerFunc(func(w http.ResponseWriter, q *http.Request) {
		b, _ := io.ReadAll(q.Body)
		seen = string(b)
		w.Header().Set("Set-Cookie", "s=1")
		w.WriteHeader(404)
	}))
	req := httptest.NewRequest("POST", "http://nsx.np.dl.playstation.net/store/list?id=5&ticket=abc", strings.NewReader("a=1&password=hunter2"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer x")
	req.Header.Set("User-Agent", "PS3")
	h.ServeHTTP(httptest.NewRecorder(), req)
	if seen != "a=1&password=hunter2" {
		t.Fatalf("handler saw %q", seen)
	}
	req2 := httptest.NewRequest("POST", "http://auth.api.np.ac.playstation.net/nav/auth", strings.NewReader("secretbody"))
	h.ServeHTTP(httptest.NewRecorder(), req2)
	files, _ := filepath.Glob(filepath.Join(dir, "capture-*.jsonl"))
	if len(files) != 1 {
		t.Fatalf("files %v", files)
	}
	b, _ := os.ReadFile(files[0])
	s := string(b)
	for _, bad := range []string{"hunter2", "abc", "Bearer", "secretbody", "s=1"} {
		if strings.Contains(s, bad) {
			t.Errorf("leaked %q in %s", bad, s)
		}
	}
	for _, want := range []string{"nsx.np.dl.playstation.net", `"status":404`, "User-Agent", `a=1\u0026password=[masked]`, "authentication path"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
}

func TestPrune(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "capture-2020-01-01.jsonl")
	os.WriteFile(old, []byte("x"), 0o600)
	r, _ := New(dir)
	r.now = func() time.Time { return time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC) }
	r.write(entry{Kind: "tls"})
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatal("old file kept")
	}
}
