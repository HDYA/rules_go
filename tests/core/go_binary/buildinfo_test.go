package main

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestBuildInfo(t *testing.T) {
	for _, tc := range []struct {
		name     string
		wantPath string
	}{
		{
			name:     "buildinfo_bin",
			wantPath: "example.com/buildinfo",
		},
		{
			name:     "buildinfo_importmap_bin",
			wantPath: "example.com/buildinfo_importmap",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lines, out := runBuildInfoBinary(t, tc.name)

			if got := lines["ok"]; got != "true" {
				t.Fatalf("debug.ReadBuildInfo ok = %q, want true\noutput:\n%s", got, out)
			}
			if got := lines["path"]; got != tc.wantPath {
				t.Errorf("BuildInfo.Path = %q, want %q", got, tc.wantPath)
			}

			wantSettings := map[string]string{
				"setting -buildmode": "exe",
				"setting -compiler":  "gc",
				"setting GOARCH":     runtime.GOARCH,
				"setting GOOS":       runtime.GOOS,
			}
			for key, want := range wantSettings {
				if got := lines[key]; got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
			if got := lines["setting CGO_ENABLED"]; got != "0" && got != "1" {
				t.Errorf("setting CGO_ENABLED = %q, want 0 or 1", got)
			}
		})
	}
}

func runBuildInfoBinary(t *testing.T, name string) (map[string]string, []byte) {
	t.Helper()

	bin, ok := bazel.FindBinary("tests/core/go_binary", name)
	if !ok {
		t.Fatalf("could not find %s", name)
	}
	out, err := exec.Command(bin).CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", name, err, out)
	}

	lines := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			lines[key] = value
		}
	}
	return lines, out
}
