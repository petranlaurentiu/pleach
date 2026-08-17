package action

import (
	"github.com/micro-editor/micro/v2/internal/buffer"
	"github.com/micro-editor/micro/v2/internal/screen"
)

type dockDir int

const (
	dockRight dockDir = iota
	dockLeft
	dockBelow
	dockAbove
)

func isFileManagerPane(p Pane) bool {
	bp, ok := p.(*BufPane)
	return ok && isFileManager(bp)
}

func paneIndex(p Pane) int {
	for i, q := range MainTab().Panes {
		if q == p {
			return i
		}
	}
	return -1
}

func arrangeTarget(p Pane) Pane {
	var tree, other Pane
	for _, q := range MainTab().Panes {
		if q == p {
			continue
		}
		if isFileManagerPane(q) {
			tree = q
			continue
		}
		if other == nil {
			other = q
		}
	}
	if other != nil {
		return other
	}
	return tree
}

func splitNode(id uint64, dir dockDir) uint64 {
	node := MainTab().GetNode(id)
	if node == nil {
		return 0
	}
	switch dir {
	case dockRight:
		return node.VSplit(true)
	case dockLeft:
		return node.VSplit(false)
	case dockAbove:
		return node.HSplit(false)
	case dockBelow:
		return node.HSplit(true)
	default:
		return node.HSplit(true)
	}
}

func moveActivePane(dir dockDir) {
	if Tabs == nil || len(Tabs.List) == 0 {
		return
	}
	tab := MainTab()
	if len(tab.Panes) < 2 {
		InfoBar.Message("Open another pane first, then move this one")
		return
	}

	p := tab.Panes[tab.active]
	if isFileManagerPane(p) {
		InfoBar.Message("The file tree stays on the left")
		return
	}

	target := arrangeTarget(p)
	if target == nil {
		InfoBar.Message("No other pane to arrange against")
		return
	}

	node := tab.GetNode(p.ID())
	if node == nil || !node.Unsplit() {
		InfoBar.Error("Could not move this pane")
		return
	}

	newID := splitNode(target.ID(), dir)
	if newID == 0 {
		InfoBar.Error("Could not place this pane")
		return
	}
	p.SetID(newID)

	from := paneIndex(p)
	to := paneIndex(target)
	if from >= 0 && to >= 0 && from != to {
		tab.RemovePane(from)
		if from < to {
			to--
		}
		switch dir {
		case dockRight, dockBelow:
			to++
		case dockLeft, dockAbove:
		default:
		}
		if to < 0 {
			to = 0
		}
		if to > len(tab.Panes) {
			to = len(tab.Panes)
		}
		tab.AddPane(p, to)
	}

	tab.Resize()
	if i := paneIndex(p); i >= 0 {
		tab.SetActive(i)
	}
}

func newSplit(dir dockDir) {
	bp := paneForSplit()
	if bp == nil {
		InfoBar.Error("No editor pane to split")
		return
	}
	empty := buffer.NewBufferFromString("", "", buffer.BTDefault)
	switch dir {
	case dockRight:
		bp.VSplitIndex(empty, true)
	case dockLeft:
		bp.VSplitIndex(empty, false)
	case dockAbove:
		bp.HSplitIndex(empty, false)
	case dockBelow:
		bp.HSplitIndex(empty, true)
	default:
		bp.HSplitBuf(empty)
	}
}

func showLayoutMenu() {
	items := []menuItem{
		{Label: "New split right", Run: func() { newSplit(dockRight) }},
		{Label: "New split below", Run: func() { newSplit(dockBelow) }},
		{Label: "Move this pane right", Run: func() { moveActivePane(dockRight) }},
		{Label: "Move this pane left", Run: func() { moveActivePane(dockLeft) }},
		{Label: "Move this pane below", Run: func() { moveActivePane(dockBelow) }},
		{Label: "Move this pane above", Run: func() { moveActivePane(dockAbove) }},
		{Label: "Terminal on the right", Run: func() { OpenTerminalPanel(nil, dockRight) }},
		{Label: "Terminal below", Run: func() { OpenTerminalPanel(nil, dockBelow) }},
		{Label: "Agent on the right", Run: func() { OpenTerminalPanel(defaultAgentArgs(), dockRight) }},
	}
	_, h := screen.Screen.Size()
	x, y := lastMouseX, lastMouseY
	if x == 0 && y == 0 {
		x, y = 1, h-12
	}
	ShowMenuAt("Layout", items, x, y)
}

func (h *BufPane) LayoutCmd(args []string) {
	if len(args) == 0 {
		showLayoutMenu()
		return
	}
	switch args[0] {
	case "right":
		moveActivePane(dockRight)
	case "left":
		moveActivePane(dockLeft)
	case "below", "bottom", "down":
		moveActivePane(dockBelow)
	case "above", "top", "up":
		moveActivePane(dockAbove)
	case "vsplit", "side":
		newSplit(dockRight)
	case "hsplit", "stack":
		newSplit(dockBelow)
	case "menu":
		showLayoutMenu()
	default:
		InfoBar.Error("Usage: layout [right|left|below|above|vsplit|hsplit]")
	}
}

func (h *BufPane) MovePaneRight() bool {
	moveActivePane(dockRight)
	return true
}

func (h *BufPane) MovePaneLeft() bool {
	moveActivePane(dockLeft)
	return true
}

func (h *BufPane) MovePaneBelow() bool {
	moveActivePane(dockBelow)
	return true
}

func (h *BufPane) MovePaneAbove() bool {
	moveActivePane(dockAbove)
	return true
}

func (h *BufPane) ShowLayoutMenu() bool {
	showLayoutMenu()
	return true
}
