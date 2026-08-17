package action

import (
	"fmt"
	"regexp"

	"github.com/micro-editor/micro/v2/internal/buffer"
)

type findSession struct {
	pane     *BufPane
	useRegex bool
	active   bool
}

var currentFind findSession

func isFindPrompt(ptype string) bool {
	return ptype == "Find" || ptype == "FindRegex"
}

func startFind(h *BufPane, useRegex bool) bool {
	if h == nil || isFileManager(h) {
		return false
	}

	h.searchOrig = h.Cursor.Loc
	currentFind = findSession{pane: h, useRegex: useRegex, active: true}

	pattern := string(h.Cursor.GetSelection())
	if useRegex && pattern != "" {
		pattern = regexp.QuoteMeta(pattern)
	}
	if pattern == "" {
		pattern = h.Buf.LastSearch
	}

	eventCallback := func(resp string) {
		applyFindQuery(h, resp, useRegex, false, true)
	}
	findCallback := func(resp string, canceled bool) {
		currentFind.active = false
		if canceled {
			if currentFind.pane != nil && currentFind.pane.Buf != nil {
				currentFind.pane.Buf.HighlightSearch = false
			}
			return
		}
		applyFindQuery(h, resp, useRegex, true, false)
	}

	if pattern != "" {
		eventCallback(pattern)
	} else {
		InfoBar.Msg = findPromptLabel(0, 0, useRegex)
	}

	ptype := "Find"
	if useRegex {
		ptype = "FindRegex"
	}
	InfoBar.Prompt(findPromptLabel(0, 0, useRegex), pattern, ptype, eventCallback, findCallback)
	if pattern != "" {
		InfoBar.SelectAll()
		updateFindPrompt(h, pattern, useRegex)
	}
	return true
}

func applyFindQuery(h *BufPane, resp string, useRegex bool, moveNext bool, fromOrig bool) {
	if h == nil || h.Buf == nil {
		return
	}
	if resp == "" {
		h.Cursor.ResetSelection()
		h.Buf.HighlightSearch = false
		h.Buf.LastSearch = ""
		InfoBar.Msg = findPromptLabel(0, 0, useRegex)
		return
	}

	h.Buf.LastSearch = resp
	h.Buf.LastSearchRegex = useRegex
	h.Buf.HighlightSearch = true

	from := h.Cursor.Loc
	if fromOrig {
		from = h.searchOrig
	} else if h.Cursor.HasSelection() {
		from = h.Cursor.CurSelection[1]
	}

	match, found, err := h.Buf.FindNext(resp, h.Buf.Start(), h.Buf.End(), from, true, useRegex)
	if err != nil {
		InfoBar.Msg = "Find: " + err.Error() + " "
		return
	}
	if !found && !fromOrig && !moveNext {
		match, found, _ = h.Buf.FindNext(resp, h.Buf.Start(), h.Buf.End(), h.searchOrig, true, useRegex)
	}
	if found {
		selectMatch(h, match)
	} else if !moveNext {
		h.GotoLoc(h.searchOrig)
		h.Cursor.ResetSelection()
	}
	updateFindPrompt(h, resp, useRegex)
}

func commitFind(down bool) {
	if !currentFind.active || currentFind.pane == nil {
		return
	}
	h := currentFind.pane
	resp := string(InfoBar.LineBytes(0))
	h.Buf.LastSearch = resp
	h.Buf.LastSearchRegex = currentFind.useRegex
	if resp == "" {
		updateFindPrompt(h, resp, currentFind.useRegex)
		return
	}
	h.Buf.HighlightSearch = true
	if down {
		h.FindNext()
	} else {
		h.FindPrevious()
	}
	updateFindPrompt(h, resp, currentFind.useRegex)
}

func updateFindPrompt(h *BufPane, resp string, useRegex bool) {
	if h == nil || h.Buf == nil {
		return
	}
	cur, total := matchPosition(h, resp, useRegex)
	if InfoBar != nil && InfoBar.HasPrompt && isFindPrompt(InfoBar.PromptType) {
		InfoBar.Msg = findPromptLabel(cur, total, useRegex)
	}
}

func findPromptLabel(cur, total int, useRegex bool) string {
	kind := "Find"
	if useRegex {
		kind = "Find (regex)"
	}
	if total <= 0 {
		return kind + " 0: "
	}
	return fmt.Sprintf("%s %d/%d: ", kind, cur, total)
}

func matchPosition(h *BufPane, resp string, useRegex bool) (int, int) {
	matches, err := h.Buf.FindAll(resp, useRegex)
	if err != nil || len(matches) == 0 {
		return 0, 0
	}
	if h.Cursor.HasSelection() {
		start, end := h.Cursor.CurSelection[0], h.Cursor.CurSelection[1]
		for i, m := range matches {
			if m[0] == start && m[1] == end {
				return i + 1, len(matches)
			}
		}
	}
	return 1, len(matches)
}

func selectMatch(h *BufPane, match [2]buffer.Loc) {
	h.Cursor.SetSelectionStart(match[0])
	h.Cursor.SetSelectionEnd(match[1])
	h.Cursor.OrigSelection[0] = h.Cursor.CurSelection[0]
	h.Cursor.OrigSelection[1] = h.Cursor.CurSelection[1]
	h.GotoLoc(h.Cursor.CurSelection[1])
}
