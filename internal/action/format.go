package action

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type formatterKind int

const (
	formatterNone formatterKind = iota
	formatterBiome
	formatterPrettier
)

type projectFormatter struct {
	Kind formatterKind
	Name string
	Root string
}

var prettierConfigNames = []string{
	".prettierrc",
	".prettierrc.json",
	".prettierrc.json5",
	".prettierrc.yml",
	".prettierrc.yaml",
	".prettierrc.js",
	".prettierrc.cjs",
	".prettierrc.mjs",
	".prettierrc.toml",
	"prettier.config.js",
	"prettier.config.cjs",
	"prettier.config.mjs",
	"prettier.config.ts",
}

func (f *projectFormatter) label() string {
	if f == nil {
		return ""
	}
	return "Format with " + f.Name
}

func detectFormatter(start string) *projectFormatter {
	dir := startDir(start)
	if dir == "" {
		return nil
	}

	var biome, prettier *projectFormatter
	walkProjectDirs(dir, func(d string) bool {
		if biome == nil && hasBiome(d) {
			biome = &projectFormatter{Kind: formatterBiome, Name: "Biome", Root: d}
		}
		if prettier == nil && hasPrettier(d) {
			prettier = &projectFormatter{Kind: formatterPrettier, Name: "Prettier", Root: d}
		}
		return biome != nil
	})
	if biome != nil {
		return biome
	}
	return prettier
}

func startDir(start string) string {
	if start == "" {
		wd, err := os.Getwd()
		if err != nil {
			return ""
		}
		return wd
	}
	info, err := os.Stat(start)
	if err == nil && info.IsDir() {
		return start
	}
	return filepath.Dir(start)
}

func walkProjectDirs(dir string, visit func(string) bool) {
	for {
		if visit(dir) {
			return
		}
		if fileExists(filepath.Join(dir, ".git")) {
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func hasBiome(dir string) bool {
	return fileExists(filepath.Join(dir, "biome.json")) ||
		fileExists(filepath.Join(dir, "biome.jsonc")) ||
		fileExists(filepath.Join(dir, "node_modules", ".bin", "biome")) ||
		packageJSONHas(dir, "@biomejs/biome")
}

func hasPrettier(dir string) bool {
	for _, name := range prettierConfigNames {
		if fileExists(filepath.Join(dir, name)) {
			return true
		}
	}
	return fileExists(filepath.Join(dir, "node_modules", ".bin", "prettier")) ||
		packageJSONHas(dir, "prettier")
}

func packageJSONHas(dir string, dep string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return false
	}
	return bytes.Contains(data, []byte(`"`+dep+`"`))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func resolveTool(root, bin, npmPkg string) []string {
	local := filepath.Join(root, "node_modules", ".bin", bin)
	if fileExists(local) {
		return []string{local}
	}
	if p, err := exec.LookPath(bin); err == nil {
		return []string{p}
	}
	return []string{"npx", "--no-install", npmPkg}
}

func formatterArgv(f *projectFormatter, filePath string) []string {
	switch f.Kind {
	case formatterBiome:
		return append(resolveTool(f.Root, "biome", "@biomejs/biome"), "format", "--stdin-file-path", filePath)
	case formatterPrettier:
		return append(resolveTool(f.Root, "prettier", "prettier"), "--stdin-filepath", filePath)
	default:
		return nil
	}
}

func bufferFilePath(h *BufPane) string {
	if h == nil || h.Buf == nil {
		return ""
	}
	if h.Buf.AbsPath != "" {
		return h.Buf.AbsPath
	}
	return h.Buf.Path
}

func stdinFilePath(h *BufPane) string {
	if p := bufferFilePath(h); p != "" {
		return p
	}
	ft, _ := h.Buf.Settings["filetype"].(string)
	if ft != "" && ft != "unknown" {
		return "untitled." + ft
	}
	return "untitled.txt"
}

func formatDocument(h *BufPane) {
	if h == nil || h.Buf == nil || isFileManager(h) {
		InfoBar.Message("No file to format")
		return
	}
	if h.Buf.Type.Readonly {
		InfoBar.Error("Buffer is read-only")
		return
	}

	f := detectFormatter(bufferFilePath(h))
	if f == nil {
		InfoBar.Message("No Biome or Prettier in this project")
		return
	}

	argv := formatterArgv(f, stdinFilePath(h))
	if len(argv) == 0 {
		InfoBar.Error("No formatter available")
		return
	}

	start, end := h.Buf.Start(), h.Buf.End()
	text := string(h.Buf.Substr(start, end))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = f.Root
	cmd.Stdin = strings.NewReader(text)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	formatted := stdout.String()
	if formatted == "" && err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		InfoBar.Error(f.Name + ": " + msg)
		return
	}
	if formatted == text {
		InfoBar.Message("Already formatted")
		return
	}

	h.Buf.Replace(start, end, formatted)
	h.Relocate()
	InfoBar.Message("Formatted with " + f.Name)
}

func (h *BufPane) FormatDocument() bool {
	formatDocument(h)
	return true
}

func (h *BufPane) FormatCmd(args []string) {
	formatDocument(h)
}

func editorCut(h *BufPane) {
	if !h.Cut() {
		h.CutLine()
	}
}

func editorCopy(h *BufPane) {
	if !h.Copy() {
		h.CopyLine()
	}
}
