package action

import (
	"fmt"
	"strings"

	"github.com/micro-editor/micro/v2/internal/buffer"
	"github.com/micro-editor/micro/v2/internal/util"
)

// Sequences Ghostty/kitty/xterm may send for Ctrl/Cmd chords when an
// extended keyboard protocol is on. We do not enable that protocol
// (it broke Ctrl+Z); these are a fallback if a leftover mode is active.
func bindTerminalCmdKeys() {
	type bind struct {
		key, mods int
		action    string
	}
	for _, item := range []bind{
		{127, 9, "DeleteToStartOfLine"},
		{8, 9, "DeleteToStartOfLine"},
		{127, 5, "DeleteToStartOfLine"},
		{8, 5, "DeleteToStartOfLine"},
		{117, 5, "DeleteToStartOfLine"},
		{122, 5, "Undo"},
		{122, 9, "Undo"},
		{3, 9, "DeleteToEndOfLine"},
	} {
		for _, seq := range extendedKeySeqs(item.key, item.mods) {
			BindKey(seq, item.action, BufMapEvent)
		}
	}
}

func extendedKeySeqs(key, mods int) []string {
	return []string{
		fmt.Sprintf("\x1b[%d;%du", key, mods),
		fmt.Sprintf("\x1b[%d;%d:1u", key, mods),
		fmt.Sprintf("\x1b[27;%d;%d~", mods, key),
	}
}

func (h *BufPane) deleteSelectionOr(run func()) bool {
	if h.Cursor.HasSelection() {
		h.Cursor.DeleteSelection()
		h.Cursor.ResetSelection()
		h.Relocate()
		return true
	}
	run()
	h.Relocate()
	return true
}

// DeleteToStartOfLine deletes from the cursor to column 0 (VS Code Cmd+Backspace).
func (h *BufPane) DeleteToStartOfLine() bool {
	return h.deleteSelectionOr(func() {
		start := buffer.Loc{X: 0, Y: h.Cursor.Y}
		if h.Cursor.Loc.GreaterThan(start) {
			h.Buf.Remove(start, h.Cursor.Loc)
		}
	})
}

// DeleteToEndOfLine deletes from the cursor to the end of the line (VS Code Cmd+Delete).
func (h *BufPane) DeleteToEndOfLine() bool {
	return h.deleteSelectionOr(func() {
		end := buffer.Loc{X: util.CharacterCount(h.Buf.LineBytes(h.Cursor.Y)), Y: h.Cursor.Y}
		if h.Cursor.Loc.LessThan(end) {
			h.Buf.Remove(h.Cursor.Loc, end)
		}
	})
}

// InsertLineAfter inserts a new line below, like VS Code Cmd+Enter.
func (h *BufPane) InsertLineAfter() bool {
	y := h.Cursor.Y
	if h.Cursor.HasSelection() {
		end := h.Cursor.CurSelection[1]
		if h.Cursor.CurSelection[0].GreaterThan(end) {
			end = h.Cursor.CurSelection[0]
		}
		y = end.Y
		if end.X == 0 && y > 0 {
			y--
		}
		h.Cursor.Deselect(false)
	}
	h.Cursor.Y = y
	h.Cursor.End()
	return h.InsertNewline()
}

// InsertLineBefore inserts a new line above, like VS Code Shift+Cmd+Enter.
func (h *BufPane) InsertLineBefore() bool {
	y := h.Cursor.Y
	if h.Cursor.HasSelection() {
		start := h.Cursor.CurSelection[0]
		if h.Cursor.CurSelection[1].LessThan(start) {
			start = h.Cursor.CurSelection[1]
		}
		y = start.Y
		h.Cursor.Deselect(true)
	}
	h.Cursor.GotoLoc(buffer.Loc{X: 0, Y: y})
	h.Buf.Insert(h.Cursor.Loc, "\n")
	if h.Buf.Settings["autoindent"].(bool) {
		ws := util.GetLeadingWhitespace(h.Buf.LineBytes(h.Cursor.Y + 1))
		if len(ws) > 0 {
			h.Buf.Insert(h.Cursor.Loc, string(ws))
		}
	}
	h.Cursor.StoreVisualX()
	h.Relocate()
	return true
}

// DuplicateLineUp copies the current line(s) above, like VS Code Shift+Option+Up.
func (h *BufPane) DuplicateLineUp() bool {
	startY, endY := h.Cursor.Y, h.Cursor.Y
	if h.Cursor.HasSelection() {
		start, end := h.Cursor.CurSelection[0], h.Cursor.CurSelection[1]
		if start.GreaterThan(end) {
			start, end = end, start
		}
		startY, endY = start.Y, end.Y
		if end.X == 0 && endY > startY {
			endY--
		}
	}
	var b strings.Builder
	for y := startY; y <= endY; y++ {
		b.Write(h.Buf.LineBytes(y))
		b.WriteByte('\n')
	}
	h.Buf.Insert(buffer.Loc{X: 0, Y: startY}, b.String())
	h.Relocate()
	return true
}
