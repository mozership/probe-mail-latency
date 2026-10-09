package maillatency

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/linyows/probe/actionref"
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

// TestManifestsDeclare checks that the action.yml a release writes, and the
// one the e2e workflow uses, declare the params and the guard the action
// has, so that probe check and the guard of a run take it as it is.
func TestManifestsDeclare(t *testing.T) {
	checksums := filepath.Join(t.TempDir(), "checksums.txt")
	if err := os.WriteFile(checksums, []byte(strings.Repeat("a", 64)+"  probe-mail-latency_linux_amd64\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	released, err := exec.Command("sh", "../scripts/action-yml.sh", "v0.0.0", checksums).Output()
	if err != nil {
		t.Fatalf("scripts/action-yml.sh: %v", err)
	}
	local, err := os.ReadFile("../e2e/action/action.yml")
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{"scripts/action-yml.sh": released, "e2e/action/action.yml": local} {
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
