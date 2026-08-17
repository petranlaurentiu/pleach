package display

import "testing"

func TestBufWindowHitCloseButton(t *testing.T) {
	w := &BufWindow{closeBtnX: 20, closeBtnY: 10, closeBtnW: 2}
	if !w.HitCloseButton(20, 10) || !w.HitCloseButton(21, 10) {
		t.Fatal("expected hits on the close control")
	}
	if w.HitCloseButton(19, 10) || w.HitCloseButton(20, 9) || w.HitCloseButton(22, 10) {
		t.Fatal("expected misses around the close control")
	}
}

func TestBufWindowHitCloseButtonDisabled(t *testing.T) {
	w := &BufWindow{}
	if w.HitCloseButton(0, 0) {
		t.Fatal("empty control should not hit")
	}
}
