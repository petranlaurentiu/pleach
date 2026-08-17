package action

import (
	"testing"

	"github.com/micro-editor/tcell/v2"
)

func TestDecodeKittyCtrlZ(t *testing.T) {
	ke, ok := decodeExtendedKey("\x1b[122;5u")
	if !ok {
		t.Fatal("expected decode")
	}
	if ke.code != tcell.KeyCtrlZ {
		t.Fatalf("code %v", ke.code)
	}
	if ke.mod&tcell.ModCtrl == 0 {
		t.Fatalf("mod %v", ke.mod)
	}
}

func TestDecodeXtermCmdBackspace(t *testing.T) {
	ke, ok := decodeExtendedKey("\x1b[27;9;127~")
	if !ok {
		t.Fatal("expected decode")
	}
	if ke.code != tcell.KeyBackspace2 {
		t.Fatalf("code %v", ke.code)
	}
	if ke.mod&tcell.ModMeta == 0 {
		t.Fatalf("mod %v", ke.mod)
	}
}

func TestDecodeKittyCmdZ(t *testing.T) {
	ke, ok := decodeExtendedKey("\x1b[122;9u")
	if !ok {
		t.Fatal("expected decode")
	}
	if ke.r != 'z' {
		t.Fatalf("rune %q", ke.r)
	}
	if ke.mod&tcell.ModMeta == 0 {
		t.Fatalf("mod %v", ke.mod)
	}
}
