package maillatency

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/mozership/probe/actionref"
)

func TestWriteNewFile(t *testing.T) {
	// Runs in the same second get their own files instead of overwriting the
	// first one.
	dir := t.TempDir()
	var paths []string
	for i, data := range []string{"first", "second", "third"} {
		path, err := writeNewFile(dir, "20261003-120000", []byte(data))
		if err != nil {
			t.Fatalf("writeNewFile() #%d error: %v", i+1, err)
		}
		paths = append(paths, path)
	}

	want := []string{
		"mail-latency.20261003-120000.csv",
		"mail-latency.20261003-120000-2.csv",
		"mail-latency.20261003-120000-3.csv",
	}
	for i, path := range paths {
		if filepath.Base(path) != want[i] {
			t.Errorf("file #%d = %s, want %s", i+1, filepath.Base(path), want[i])
		}
	}
	got, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "first" {
		t.Errorf("first file = %q, want it left as it was", got)
	}
}

// TestManifestsDeclare checks that action.yml, and the one the e2e workflow
// uses, declare the params and the guard the action has, so that probe
// check and the guard of a run take it as it is. probe manifest keeps what
// action.yml declares in the one of a release.
func TestManifestsDeclare(t *testing.T) {
	for _, name := range []string{"action.yml", "e2e/action/action.yml"} {
		data, err := os.ReadFile(filepath.Join("..", name))
		if err != nil {
			t.Fatal(err)
		}
		m, err := actionref.ParseManifest(data)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !slices.Equal(m.Params, Params()) {
			t.Errorf("%s declares params %v, want %v", name, m.Params, Params())
		}
		if !slices.Equal(m.Guard, Keeps()) {
			t.Errorf("%s declares guard %v, want %v", name, m.Guard, Keeps())
		}
	}
}
