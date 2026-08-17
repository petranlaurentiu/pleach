package action

import (
	"strings"
	"testing"

	"github.com/micro-editor/tcell/v2"
)

func TestArtLinesTrims(t *testing.T) {
	lines := artLines(jarvisSmall)
	if len(lines) < 6 {
		t.Fatalf("missing JARVIS lines: %d", len(lines))
	}
	if lines[0] == "" {
		t.Fatal("leading blank line was not trimmed")
	}
}

func TestRenderBannerWelcome(t *testing.T) {
	lines := renderBanner("WELCOME", 80)
	if len(lines) < 3 {
		t.Fatalf("banner too short: %#v", lines)
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "█") {
		t.Fatalf("expected block letters, got:\n%s", joined)
	}
}

func TestBannerNameStripsJunk(t *testing.T) {
	if got := bannerName("Laurentiu.petran"); got != "LAURENTIU.PETRAN" {
		t.Fatalf("got %q", got)
	}
	if got := bannerName(""); got != "SIR" {
		t.Fatalf("got %q", got)
	}
}

func TestIsEnter(t *testing.T) {
	if !isEnterKey(tcell.KeyEnter, 0) {
		t.Fatal("Enter should dismiss")
	}
	if isEnterKey(tcell.KeyEsc, 0) {
		t.Fatal("Esc should not dismiss")
	}
}
