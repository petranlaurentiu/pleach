package action

import (
	"github.com/micro-editor/micro/v2/internal/config"
	"github.com/micro-editor/micro/v2/internal/screen"
)

func commandBarRow(row int) []barButton {
	switch row {
	case 0:
		return []barButton{
			{Label: themeButtonLabel(), Run: showThemeMenu, Accent: true},
			{Label: "File", Run: showFileMenu},
			{Label: "Edit", Run: showEditMenu},
			{Label: "View", Run: showViewMenu},
			{Label: "Panels", Run: showPanelsMenu},
			{Label: "Options", Run: showOptionsMenu},
			{Label: "Help", Run: showHelpMenu},
		}
	case 1:
		return []barButton{
			{Label: "New", Run: func() { withPane(func(bp *BufPane) { bp.AddTab() }) }},
			{Label: "Open", Run: func() { withPane(func(bp *BufPane) { bp.OpenFile() }) }},
			{Label: "Save", Run: func() { withPane(func(bp *BufPane) { bp.Save() }) }},
			{Label: "Find", Run: func() { withPane(func(bp *BufPane) { bp.Find() }) }},
			{Label: "Undo", Run: func() { withPane(func(bp *BufPane) { bp.Undo() }) }},
			{Label: "Redo", Run: func() { withPane(func(bp *BufPane) { bp.Redo() }) }},
			{Label: "Quit", Run: quitPleach},
		}
	default:
		return []barButton{
			{Label: "Tree", Run: func() { runCommand("tree") }},
			{Label: "Terminal", Run: func() { OpenTerminalPanel(nil, dockBelow) }},
			{Label: "Agent", Run: func() { OpenTerminalPanel(defaultAgentArgs(), dockBelow) }},
			{Label: "Close", Run: closeActivePanel},
			{Label: "Side by side", Run: func() { newSplit(dockRight) }},
			{Label: "Stacked", Run: func() { newSplit(dockBelow) }},
			{Label: "Command", Run: runCommandMode},
		}
	}
}

func withPane(fn func(*BufPane)) {
	if Tabs == nil || len(Tabs.List) == 0 {
		return
	}
	if bp := MainTab().CurPane(); bp != nil && !isFileManager(bp) {
		fn(bp)
		return
	}
	if bp := paneForSplit(); bp != nil {
		fn(bp)
	}
}

func quitPleach() {
	if Tabs == nil || len(Tabs.List) == 0 {
		return
	}
	if bp := MainTab().CurPane(); bp != nil {
		bp.QuitAll()
		return
	}
	if bp := paneForSplit(); bp != nil {
		bp.QuitAll()
	}
}

func showNamedMenu(title string, items []menuItem) {
	_, h := screen.Screen.Size()
	x, y := lastMouseX, lastMouseY
	if x == 0 && y == 0 {
		x, y = 1, h-len(items)-4
	}
	ShowMenuAt(title, items, x, y)
}

func boolOpt(name string) bool {
	if bp := paneForSplit(); bp != nil {
		if v, ok := bp.Buf.Settings[name].(bool); ok {
			return v
		}
	}
	v, ok := config.GetGlobalOption(name).(bool)
	return ok && v
}

func marked(name, label string) string {
	if boolOpt(name) {
		return label + "  [on]"
	}
	return label + "  [off]"
}

func toggleOpt(name string) {
	withPane(func(bp *BufPane) {
		bp.HandleCommand("toggle " + name)
	})
}

func startReplace() {
	withPane(func(bp *BufPane) {
		InfoBar.Prompt("> ", "replace ", "Command", nil, func(resp string, canceled bool) {
			if !canceled {
				bp.HandleCommand(resp)
			}
		})
	})
}

func openHelpTopic(topic string) {
	withPane(func(bp *BufPane) {
		if topic == "" {
			bp.ToggleHelp()
			return
		}
		bp.HandleCommand("help " + topic)
	})
}

func showFileMenu() {
	showNamedMenu("File", []menuItem{
		{Label: "New tab", Run: func() { withPane(func(bp *BufPane) { bp.AddTab() }) }},
		{Label: "Open…", Run: func() { withPane(func(bp *BufPane) { bp.OpenFile() }) }},
		{Label: "Save", Run: func() { withPane(func(bp *BufPane) { bp.Save() }) }},
		{Label: "Save as…", Run: func() { withPane(func(bp *BufPane) { bp.SaveAs() }) }},
		{Label: "Reload", Run: func() { runCommand("reopen") }},
		{Label: "Close", Run: closeActivePanel},
		{Label: "Quit", Run: quitPleach},
	})
}

func showEditMenu() {
	items := []menuItem{
		{Label: "Undo", Run: func() { withPane(func(bp *BufPane) { bp.Undo() }) }},
		{Label: "Redo", Run: func() { withPane(func(bp *BufPane) { bp.Redo() }) }},
		{Label: "Cut", Run: func() { withPane(editorCut) }},
		{Label: "Copy", Run: func() { withPane(editorCopy) }},
		{Label: "Paste", Run: func() { withPane(func(bp *BufPane) { bp.Paste() }) }},
		{Label: "Find…", Run: func() { withPane(func(bp *BufPane) { bp.Find() }) }},
		{Label: "Find next", Run: func() { withPane(func(bp *BufPane) { bp.FindNext() }) }},
		{Label: "Replace…", Run: startReplace},
		{Label: "Select all", Run: func() { withPane(func(bp *BufPane) { bp.SelectAll() }) }},
	}
	if f := detectFormatter(bufferFilePath(paneForSplit())); f != nil {
		items = append(items, menuItem{Label: f.label(), Run: func() { withPane(formatDocument) }})
	}
	showNamedMenu("Edit", items)
}

func showViewMenu() {
	showNamedMenu("View", []menuItem{
		{Label: "Theme…", Run: showThemeMenu},
		{Label: "Side by side", Run: func() { newSplit(dockRight) }},
		{Label: "Stacked", Run: func() { newSplit(dockBelow) }},
		{Label: "Layout…", Run: showLayoutMenu},
		{Label: marked("ruler", "Ruler"), Run: func() { withPane(func(bp *BufPane) { bp.ToggleRuler() }) }},
		{Label: marked("softwrap", "Soft wrap"), Run: func() { toggleOpt("softwrap") }},
		{Label: marked("syntax", "Syntax"), Run: func() { toggleOpt("syntax") }},
		{Label: marked("cursorline", "Cursor line"), Run: func() { toggleOpt("cursorline") }},
		{Label: marked("scrollbar", "Scrollbar"), Run: func() { toggleOpt("scrollbar") }},
		{Label: marked("keymenu", "Toolbar"), Run: func() { withPane(func(bp *BufPane) { bp.ToggleKeyMenu() }) }},
	})
}

func showPanelsMenu() {
	showNamedMenu("Panels", []menuItem{
		{Label: "Tree", Run: func() { runCommand("tree") }},
		{Label: "Terminal below", Run: func() { OpenTerminalPanel(nil, dockBelow) }},
		{Label: "Terminal on the right", Run: func() { OpenTerminalPanel(nil, dockRight) }},
		{Label: "Agent", Run: func() { OpenTerminalPanel(defaultAgentArgs(), dockBelow) }},
		{Label: "Close", Run: closeActivePanel},
		{Label: "Layout…", Run: showLayoutMenu},
	})
}

func themeButtonLabel() string {
	name, _ := config.GetGlobalOption("colorscheme").(string)
	for _, t := range herdrThemes() {
		if t.name == name {
			return "Theme: " + t.label
		}
	}
	if name == "" {
		return "Theme"
	}
	return "Theme: " + name
}

func herdrThemes() []struct{ label, name string } {
	return []struct{ label, name string }{
		{"Catppuccin", "catppuccin"},
		{"Catppuccin Latte", "catppuccin-latte"},
		{"Tokyo Night", "tokyo-night"},
		{"Tokyo Night Day", "tokyo-night-day"},
		{"Dracula", "dracula"},
		{"Nord", "nord"},
		{"Gruvbox", "gruvbox-tc"},
		{"Gruvbox Light", "gruvbox-light"},
		{"One Dark", "one-dark"},
		{"One Light", "one-light"},
		{"Solarized", "solarized-tc"},
		{"Solarized Light", "solarized-light"},
		{"Kanagawa", "kanagawa"},
		{"Kanagawa Lotus", "kanagawa-lotus"},
		{"Rosé Pine", "rose-pine"},
		{"Rosé Pine Dawn", "rose-pine-dawn"},
		{"Vesper", "vesper"},
	}
}

func applyTheme(name string) {
	if err := SetGlobalOptionNative("colorscheme", name, true); err != nil {
		InfoBar.Error(err)
		return
	}
	InfoBar.Message("Theme: ", name)
}

func showThemeMenu() {
	items := make([]menuItem, 0, len(herdrThemes()))
	current, _ := config.GetGlobalOption("colorscheme").(string)
	for _, t := range herdrThemes() {
		t := t
		label := t.label
		if t.name == current {
			label += "  [on]"
		}
		items = append(items, menuItem{Label: label, Run: func() { applyTheme(t.name) }})
	}
	showNamedMenu("Theme", items)
}

func showOptionsMenu() {
	showNamedMenu("Options", []menuItem{
		{Label: "Theme…", Run: showThemeMenu},
		{Label: marked("autoindent", "Auto indent"), Run: func() { toggleOpt("autoindent") }},
		{Label: marked("tabstospaces", "Spaces for tabs"), Run: func() { toggleOpt("tabstospaces") }},
		{Label: marked("ignorecase", "Ignore case"), Run: func() { toggleOpt("ignorecase") }},
		{Label: marked("hlsearch", "Highlight search"), Run: func() { toggleOpt("hlsearch") }},
		{Label: marked("matchbrace", "Match braces"), Run: func() { toggleOpt("matchbrace") }},
		{Label: marked("wordwrap", "Word wrap"), Run: func() { toggleOpt("wordwrap") }},
		{Label: marked("statusline", "Status line"), Run: func() { toggleOpt("statusline") }},
		{Label: marked("mouse", "Mouse"), Run: func() { toggleOpt("mouse") }},
		{Label: "More settings…", Run: func() {
			withPane(func(bp *BufPane) {
				InfoBar.Prompt("> ", "set ", "Command", nil, func(resp string, canceled bool) {
					if !canceled {
						bp.HandleCommand(resp)
					}
				})
			})
		}},
	})
}

func showHelpMenu() {
	showNamedMenu("Help", []menuItem{
		{Label: "Help", Run: func() { openHelpTopic("") }},
		{Label: "Pleach", Run: func() { openHelpTopic("pleach") }},
		{Label: "Commands", Run: func() { openHelpTopic("commands") }},
		{Label: "Keybindings", Run: func() { openHelpTopic("keybindings") }},
		{Label: "Tutorial", Run: func() { openHelpTopic("tutorial") }},
		{Label: "Command bar", Run: runCommandMode},
	})
}

func (h *BufPane) OptionsCmd(args []string) {
	showOptionsMenu()
}

func (h *BufPane) ThemeCmd(args []string) {
	if len(args) == 0 {
		showThemeMenu()
		return
	}
	applyTheme(args[0])
}

func (h *BufPane) ShowOptionsMenu() bool {
	showOptionsMenu()
	return true
}
