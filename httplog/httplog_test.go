package httplog

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

func capture(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	old, flags := log.Writer(), log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() { log.SetOutput(old); log.SetFlags(flags) })
	return &buf
}

func TestLineHasMethodPathStatusDuration(t *testing.T) {
	buf := capture(t)
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusTeapot)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/v1/thing", nil))
	if !regexp.MustCompile(`^http POST /v1/thing -> 418 \d+ms\n$`).MatchString(buf.String()) {
		t.Fatalf("line = %q", buf.String())
	}
}

func TestDefaultStatusIs200AndHealthzSkipped(t *testing.T) {
	buf := capture(t)
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("x")) }))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/healthz", nil))
	if buf.Len() != 0 {
		t.Fatalf("healthz logged: %q", buf.String())
	}
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/a", nil))
	if !strings.Contains(buf.String(), "http GET /a -> 200 ") {
		t.Fatalf("line = %q", buf.String())
	}
}

func TestNoSecretsLogged(t *testing.T) {
	buf := capture(t)
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(401) }))
	req := httptest.NewRequest("POST", "/login?password=hunter2&token=abc", strings.NewReader("secretbody"))
	req.Header.Set("Authorization", "Bearer sekrit")
	req.Header.Set("Cookie", "sid=cookieval")
	req.Header.Set("X-Api-Key", "keyval")
	h.ServeHTTP(httptest.NewRecorder(), req)
	out := buf.String()
	for _, s := range []string{"hunter2", "token", "abc", "sekrit", "cookieval", "keyval", "secretbody", "?"} {
		if strings.Contains(out, s) {
			t.Fatalf("%q leaked in %q", s, out)
		}
	}
}

func TestInformationalStatusIgnored(t *testing.T) {
	buf := capture(t)
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(103)
		w.WriteHeader(204)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/e", nil))
	if !strings.Contains(buf.String(), "-> 204 ") {
		t.Fatalf("line = %q", buf.String())
	}
}

func TestFlusherAndUnwrapPreserved(t *testing.T) {
	capture(t)
	rec := httptest.NewRecorder()
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, ok := w.(http.Flusher)
		if !ok {
			t.Error("Flusher lost")
			return
		}
		f.Flush()
		if http.NewResponseController(w).Flush() != nil {
			t.Error("ResponseController cannot reach the writer")
		}
	}))
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/s", nil))
	if !rec.Flushed {
		t.Fatal("flush did not reach the underlying writer")
	}
}

func TestHijackPassesThroughOrErrors(t *testing.T) {
	capture(t)
	Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, _, err := w.(http.Hijacker).Hijack(); err == nil {
			t.Error("recorder cannot hijack; expected error")
		}
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/h", nil))
}

func TestPathCannotForgeLines(t *testing.T) {
	buf := capture(t)
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req := httptest.NewRequest("GET", "/a", nil)
	req.URL.Path = "/a\nhttp GET /forged -> 200 0ms"
	h.ServeHTTP(httptest.NewRecorder(), req)
	if strings.Count(buf.String(), "\n") != 1 {
		t.Fatalf("forged line: %q", buf.String())
	}
}

func TestExtraAndCommand(t *testing.T) {
	buf := capture(t)
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("X-S", "7") }),
		WithExtra(func(w http.ResponseWriter) string { return "s=" + w.Header().Get("X-S") }))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/x", nil))
	if !regexp.MustCompile(`-> 200 \d+ms s=7\n$`).MatchString(buf.String()) {
		t.Fatalf("line = %q", buf.String())
	}
	buf.Reset()
	Command("AS_REQ", "ok", time.Now())
	if !regexp.MustCompile(`^cmd AS_REQ -> ok \d+ms\n$`).MatchString(buf.String()) {
		t.Fatalf("line = %q", buf.String())
	}
}
