package action

import (
	"os"
	"path/filepath"

	"github.com/micro-editor/micro/v2/internal/buffer"
	"github.com/micro-editor/micro/v2/internal/clipboard"
	"github.com/micro-editor/micro/v2/internal/herdr"
	"github.com/micro-editor/tcell/v2"
)

type barButton struct {
	Label  string
	Run    func()
	Accent bool
}

func isFileManager(h *BufPane) bool {
	if h == nil || h.Buf == nil {
		return false
	}
	if h.Buf.GetName() == "filemanager" {
		return true
	}
	ft, ok := h.Buf.Settings["filetype"].(string)
	return ok && ft == "filemanager"
}

func paneForSplit() *BufPane {
	if Tabs == nil || len(Tabs.List) == 0 {
		return nil
	}
	if bp := MainTab().CurPane(); bp != nil && !isFileManager(bp) {
		return bp
	}
	for _, p := range MainTab().Panes {
		if bp, ok := p.(*BufPane); ok && !isFileManager(bp) {
			return bp
		}
	}
	for _, p := range MainTab().Panes {
		if bp, ok := p.(*BufPane); ok {
			return bp
		}
	}
	return nil
}

func runCommand(cmd string) {
	if bp := MainTab().CurPane(); bp != nil {
		bp.HandleCommand(cmd)
		return
	}
	for _, p := range MainTab().Panes {
		if bp, ok := p.(*BufPane); ok {
			bp.HandleCommand(cmd)
			return
		}
	}
}

func runCommandMode() {
	if bp := MainTab().CurPane(); bp != nil {
		bp.CommandMode()
		return
	}
	if tp, ok := MainTab().Panes[MainTab().active].(*TermPane); ok {
		tp.CommandMode()
	}
}

func defaultAgentArgs() []string {
	if p := herdr.Resolve(); p != "" {
		return []string{p}
	}
	InfoBar.Message("Downloading herdr…")
	p, err := herdr.Ensure()
	if err == nil {
		InfoBar.Message("Herdr ready")
		return []string{p}
	}
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	InfoBar.Error("Herdr unavailable: ", err)
	return []string{sh, "-c", "printf '\\n  Agent panel\\n  Could not find or download herdr.\\n  Install: curl -fsSL https://herdr.dev/install.sh | sh\\n\\n'; exec \"$SHELL\""}
}

// OpenTerminalPanel opens a shell (or args) in a split on the given side.
func OpenTerminalPanel(args []string, dir dockDir) {
	if !TermEmuSupported {
		InfoBar.Error("Terminal emulator not supported on this system")
		return
	}
	if len(args) == 0 {
		sh := os.Getenv("SHELL")
		if sh == "" {
			InfoBar.Error("Shell environment not found")
			return
		}
		args = []string{sh}
	}
	bp := paneForSplit()
	if bp == nil {
		InfoBar.Error("No editor pane to split")
		return
	}
	empty := buffer.NewBufferFromString("", "", buffer.BTScratch)
	var np *BufPane
	switch dir {
	case dockRight:
		np = bp.VSplitIndex(empty, true)
	case dockLeft:
		np = bp.VSplitIndex(empty, false)
	case dockAbove:
		np = bp.HSplitIndex(empty, false)
	case dockBelow:
		np = bp.HSplitIndex(empty, true)
	default:
		np = bp.HSplitBuf(empty)
	}
	np.openTerm(args, false)
}

func closeActivePanel() {
	if Tabs == nil || len(Tabs.List) == 0 {
		return
	}
	closeWindow(MainTab().Panes[MainTab().active])
}

// closeWindow closes the pane that was clicked or focused, like a VS Code
// editor group. The file tree is only closed when it is the target.
func closeWindow(p Pane) {
	switch pane := p.(type) {
	case *TermPane:
		closeTermWindow(pane)
	case *RawPane:
		pane.Quit()
	case *BufPane:
		closeBufWindow(pane)
	default:
		InfoBar.Message("No window to close")
	}
}

func closeTermWindow(t *TermPane) {
	if keepWorkbench(t.Tab()) {
		replaceTermWithEmptyEditor(t)
		return
	}
	t.Quit()
}

func replaceTermWithEmptyEditor(t *TermPane) {
	tab := t.Tab()
	if tab == nil {
		tab = MainTab()
	}
	empty := buffer.NewBufferFromString("", "", buffer.BTDefault)
	bp := NewBufPaneFromBuf(empty, tab)
	bp.SetID(t.ID())
	t.Close()
	for i, p := range tab.Panes {
		if p == t {
			tab.Panes[i] = bp
			tab.SetActive(i)
			tab.Resize()
			return
		}
	}
}

func closeBufWindow(h *BufPane) {
	if isFileManager(h) {
		closeTreeWindow()
		return
	}
	if keepWorkbench(h.tab) {
		replaceWithEmptyEditor(h)
		return
	}
	h.Quit()
}

func keepWorkbench(tab *Tab) bool {
	return tab != nil && len(tab.Panes) == 1 && Tabs != nil && len(Tabs.List) == 1
}

func firstEditorPane(tab *Tab) *BufPane {
	if tab == nil {
		return nil
	}
	for _, p := range tab.Panes {
		bp, ok := p.(*BufPane)
		if ok && !isFileManager(bp) {
			return bp
		}
	}
	return nil
}

func closeTreeWindow() {
	tab := MainTab()
	if firstEditorPane(tab) == nil && !hasTermPane(tab) {
		if bp := fileManagerPane(tab); bp != nil {
			empty := buffer.NewBufferFromString("", "", buffer.BTDefault)
			bp.VSplitIndex(empty, true)
		}
	}
	runCommand("tree")
}

func fileManagerPane(tab *Tab) *BufPane {
	if tab == nil {
		return nil
	}
	for _, p := range tab.Panes {
		if bp, ok := p.(*BufPane); ok && isFileManager(bp) {
			return bp
		}
	}
	return nil
}

func hasTermPane(tab *Tab) bool {
	if tab == nil {
		return false
	}
	for _, p := range tab.Panes {
		if _, ok := p.(*TermPane); ok {
			return true
		}
	}
	return false
}

func replaceWithEmptyEditor(h *BufPane) {
	if h == nil {
		return
	}
	openEmpty := func() {
		h.OpenBuffer(buffer.NewBufferFromString("", "", buffer.BTDefault))
	}
	if h.Buf != nil && h.Buf.Modified() && !h.Buf.Shared() {
		h.closePrompt("Close", openEmpty)
		return
	}
	openEmpty()
}

func (h *BufPane) TerminalCmd(args []string) {
	OpenTerminalPanel(args, dockBelow)
}

func (h *BufPane) AgentCmd(args []string) {
	if len(args) == 0 {
		args = defaultAgentArgs()
	}
	OpenTerminalPanel(args, dockBelow)
}

func (h *BufPane) ClosePanelCmd(args []string) {
	closeActivePanel()
}

func (h *BufPane) OpenTerminalAction() bool {
	OpenTerminalPanel(nil, dockBelow)
	return true
}

func (h *BufPane) OpenAgentAction() bool {
	OpenTerminalPanel(defaultAgentArgs(), dockBelow)
	return true
}

func (h *BufPane) ClosePanelAction() bool {
	closeActivePanel()
	return true
}

// ContextMenu opens a right-click menu. A click inside an existing
// selection keeps it (so Cut/Copy work); otherwise the cursor moves.
func (h *BufPane) ContextMenu(e *tcell.EventMouse) bool {
	if e != nil {
		lastMouseX, lastMouseY = e.Position()
		if isFileManager(h) || !clickInsideSelection(h, e) {
			h.MousePress(e)
		}
	} else {
		v := h.GetView()
		lastMouseX = v.X + 2
		lastMouseY = v.Y + 2
	}
	if isFileManager(h) {
		return true
	}
	showEditorContextMenu(h)
	return true
}

func clickInsideSelection(h *BufPane, e *tcell.EventMouse) bool {
	if h == nil || h.Cursor == nil || !h.Cursor.HasSelection() {
		return false
	}
	mx, my := e.Position()
	loc := h.LocFromVisual(buffer.Loc{X: mx, Y: my})
	start, end := h.Cursor.CurSelection[0], h.Cursor.CurSelection[1]
	if start.GreaterThan(end) {
		start, end = end, start
	}
	return loc.GreaterEqual(start) && loc.LessEqual(end)
}

func showEditorContextMenu(h *BufPane) {
	items := []menuItem{
		{Label: "Cut", Run: func() { withPane(editorCut) }},
		{Label: "Copy", Run: func() { withPane(editorCopy) }},
		{Label: "Paste", Run: func() { withPane(func(bp *BufPane) { bp.Paste() }) }},
		{Label: "Select all", Run: func() { withPane(func(bp *BufPane) { bp.SelectAll() }) }},
	}
	if f := detectFormatter(bufferFilePath(h)); f != nil {
		items = append(items, menuItem{Label: f.label(), Run: func() { withPane(formatDocument) }})
	}
	ShowMenuAt("Editor", items, lastMouseX, lastMouseY)
}

func showTermContextMenu(t *TermPane) {
	items := []menuItem{
		{Label: "Move pane right", Run: func() { moveActivePane(dockRight) }},
		{Label: "Move pane left", Run: func() { moveActivePane(dockLeft) }},
		{Label: "Move pane below", Run: func() { moveActivePane(dockBelow) }},
		{Label: "Move pane above", Run: func() { moveActivePane(dockAbove) }},
		{Label: "Close", Run: func() { closeWindow(t) }},
		{Label: "Command", Run: t.CommandMode},
		{Label: "Next split", Run: t.NextSplit},
	}
	ShowMenuAt("Terminal", items, lastMouseX, lastMouseY)
}

// OpenPathInEditor opens a file in the existing editor pane, like VS Code.
func OpenPathInEditor(path string) {
	if Tabs == nil || len(Tabs.List) == 0 {
		return
	}
	tab := MainTab()
	abs := path
	if resolved, err := filepath.Abs(path); err == nil {
		abs = resolved
	}

	for i, p := range tab.Panes {
		bp, ok := p.(*BufPane)
		if !ok || isFileManager(bp) || bp.Buf == nil {
			continue
		}
		if bp.Buf.AbsPath == abs || bp.Buf.Path == path {
			tab.SetActive(i)
			return
		}
	}

	buf, err := buffer.NewBufferFromFile(path, buffer.BTDefault)
	if err != nil {
		InfoBar.Error(err)
		return
	}

	for i, p := range tab.Panes {
		bp, ok := p.(*BufPane)
		if !ok || isFileManager(bp) {
			continue
		}
		if bp.Buf != nil && bp.Buf.Modified() && !bp.Buf.Shared() {
			bp.VSplitBuf(buf)
			return
		}
		bp.OpenBuffer(buf)
		tab.SetActive(i)
		return
	}

	for _, p := range tab.Panes {
		if bp, ok := p.(*BufPane); ok && isFileManager(bp) {
			bp.VSplitIndex(buf, true)
			return
		}
	}
}

// CopyToClipboard copies text to the system clipboard (used by Lua plugins).
func CopyToClipboard(text string) {
	_ = clipboard.Write(text, clipboard.ClipboardReg)
}
