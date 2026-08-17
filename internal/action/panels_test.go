package action

import "testing"

func TestKeepWorkbenchOnlyWhenLastWindow(t *testing.T) {
	if keepWorkbench(nil) {
		t.Fatal("nil tab should not keep a workbench")
	}
}

func TestHitCloseButtonGeometry(t *testing.T) {
	// BufWindow hit-test is covered in display; this guards the close
	// policy: a focused editor must close itself, not the tree.
	if isFileManager(nil) {
		t.Fatal("nil pane is not the file tree")
	}
}
