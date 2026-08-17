package herdr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/micro-editor/micro/v2/internal/config"
)

const releaseBase = "https://github.com/herdrdev/herdr/releases/latest/download/"

// Resolve returns a herdr binary path if one is already installed.
func Resolve() string {
	if p := os.Getenv("PLEACH_HERDR"); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	if p, err := exec.LookPath("herdr"); err == nil {
		return p
	}
	home, _ := os.UserHomeDir()
	exe, _ := os.Executable()
	candidates := []string{
		filepath.Join(home, ".local", "bin", "herdr"),
		"/opt/homebrew/bin/herdr",
		"/usr/local/bin/herdr",
		filepath.Join(binDir(), "herdr"),
	}
	if exe != "" {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "herdr"))
	}
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

// Ensure returns a herdr binary, downloading the official release if needed.
func Ensure() (string, error) {
	if p := Resolve(); p != "" {
		return p, nil
	}
	asset, err := releaseAsset()
	if err != nil {
		return "", err
	}
	dir := binDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create herdr dir: %w", err)
	}
	dest := filepath.Join(dir, "herdr")
	if err := download(releaseBase+asset, dest); err != nil {
		return "", fmt.Errorf("download herdr: %w", err)
	}
	return dest, nil
}

func binDir() string {
	if config.ConfigDir != "" {
		return filepath.Join(config.ConfigDir, "bin")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pleach", "bin")
}

func releaseAsset() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "herdr-macos-aarch64", nil
		}
		return "herdr-macos-x86_64", nil
	case "linux":
		if runtime.GOARCH == "arm64" {
			return "herdr-linux-aarch64", nil
		}
		return "herdr-linux-x86_64", nil
	default:
		return "", fmt.Errorf("herdr has no official binary for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

func download(url, dest string) error {
	client := &http.Client{Timeout: 3 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "pleach")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", url, resp.Status)
	}
	tmp := dest + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		os.Remove(tmp)
		return closeErr
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}
