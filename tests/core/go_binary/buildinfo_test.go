package main

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestBuildInfo(t *testing.T) {
	bin, ok := bazel.FindBinary("tests/core/go_binary", "buildinfo_bin")
	if !ok {
		t.Fatal("could not find buildinfo_bin")
	}
	out, err := exec.Command(bin).CombinedOutput()
	if err != nil {
		t.Fatalf("buildinfo_bin failed: %v\n%s", err, out)
	}

	lines := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			lines[key] = value
		}
	}

	if got := lines["ok"]; got != "true" {
		t.Fatalf("debug.ReadBuildInfo ok = %q, want true\noutput:\n%s", got, out)
	}
	if got := lines["path"]; got != "example.com/buildinfo" {
		t.Errorf("BuildInfo.Path = %q, want example.com/buildinfo", got)
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
}
