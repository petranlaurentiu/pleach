package action

import (
	"testing"

	"github.com/micro-editor/tcell/v2"
)

func TestCmdBackspaceParsesAsMeta(t *testing.T) {
	ev, ok := findSingleEvent("Cmd-Backspace")
	if !ok {
		t.Fatal("Cmd-Backspace did not parse")
	}
	k, ok := ev.(KeyEvent)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if k.mod&tcell.ModMeta == 0 {
		t.Fatalf("expected Meta, got %v", k.mod)
	}
	if k.code != tcell.KeyBackspace2 {
		t.Fatalf("expected Backspace, got %v", k.code)
	}
}

func TestCmdBackspaceRawSeqParses(t *testing.T) {
	ev, ok := findSingleEvent("\x1b[27;9;127~")
	if !ok {
		t.Fatal("raw Cmd+Backspace sequence did not parse")
	}
	raw, ok := ev.(RawEvent)
	if !ok || raw.esc != "\x1b[27;9;127~" {
		t.Fatalf("got %#v", ev)
	}
}

func TestMetaAndCmdAreTheSame(t *testing.T) {
	a, _ := findSingleEvent("Cmd-Backspace")
	b, _ := findSingleEvent("Meta-Backspace")
	if a != b {
		t.Fatalf("%#v != %#v", a, b)
	}
}
