package action

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFormatterPrefersBiome(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "biome.json"), "{}")
	write(t, filepath.Join(root, ".prettierrc"), "{}")
	write(t, filepath.Join(root, "src", "app.ts"), "const x = 1\n")

	got := detectFormatter(filepath.Join(root, "src", "app.ts"))
	if got == nil || got.Kind != formatterBiome {
		t.Fatalf("got %#v, want Biome", got)
	}
}

func TestDetectFormatterFindsPrettier(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, ".prettierrc.json"), "{}")
	write(t, filepath.Join(root, "index.js"), "let x=1\n")

	got := detectFormatter(filepath.Join(root, "index.js"))
	if got == nil || got.Kind != formatterPrettier {
		t.Fatalf("got %#v, want Prettier", got)
	}
}

func TestDetectFormatterReadsPackageJSON(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "package.json"), `{"devDependencies":{"@biomejs/biome":"2.0.0"}}`)
	write(t, filepath.Join(root, "file.ts"), "")

	got := detectFormatter(filepath.Join(root, "file.ts"))
	if got == nil || got.Kind != formatterBiome {
		t.Fatalf("got %#v, want Biome from package.json", got)
	}
}

func TestDetectFormatterNone(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "readme.md"), "hi\n")

	if got := detectFormatter(filepath.Join(root, "readme.md")); got != nil {
		t.Fatalf("got %#v, want none", got)
	}
}

func TestDetectFormatterStopsAtGitRoot(t *testing.T) {
	outer := t.TempDir()
	write(t, filepath.Join(outer, "biome.json"), "{}")
	inner := filepath.Join(outer, "repo")
	if err := os.MkdirAll(filepath.Join(inner, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(inner, "main.go"), "package main\n")

	if got := detectFormatter(filepath.Join(inner, "main.go")); got != nil {
		t.Fatalf("got %#v, want none inside nested git repo", got)
	}
}

func write(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
