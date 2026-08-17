package action

import (
	runewidth "github.com/mattn/go-runewidth"
	"github.com/micro-editor/micro/v2/internal/config"
	"github.com/micro-editor/micro/v2/internal/display"
	"github.com/micro-editor/micro/v2/internal/screen"
	"github.com/micro-editor/tcell/v2"
)

type menuItem struct {
	Label string
	Run   func()
}

type hitRegion struct {
	X0, X1, Y int
	Run       func()
}

type popupMenu struct {
	Title string
	Items []menuItem
	X, Y  int
	hits  []hitRegion
}

var (
	activeMenu       *popupMenu
	lastMouseX       int
	lastMouseY       int
	commandBarHits   []hitRegion
	commandBarArmed  bool
	commandBarPressY int
	commandBarPressX int
)

func init() {
	display.DrawCommandBar = drawCommandBar
}

func CloseMenu() {
	activeMenu = nil
}

// ShowMenuLua opens a clickable popup. The callback receives a 1-based index.
func ShowMenuLua(title string, labels []string, cb func(int)) {
	items := make([]menuItem, 0, len(labels))
	for i, label := range labels {
		i := i
		items = append(items, menuItem{
			Label: label,
			Run: func() {
				if cb != nil {
					cb(i + 1)
				}
			},
		})
	}
	ShowMenuAt(title, items, lastMouseX, lastMouseY)
}

func ShowMenuAt(title string, items []menuItem, x, y int) {
	if len(items) == 0 {
		return
	}
	w, h := screen.Screen.Size()
	width := runewidth.StringWidth(title) + 4
	for _, item := range items {
		n := runewidth.StringWidth(item.Label) + 6
		if n > width {
			width = n
		}
	}
	footer := "click or 1-9, Esc"
	if runewidth.StringWidth(footer)+4 > width {
		width = runewidth.StringWidth(footer) + 4
	}
	height := len(items) + 3
	if width > w {
		width = w
	}
	if height > h {
		height = h
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+width > w {
		x = w - width
	}
	if y+height > h {
		y = h - height
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	activeMenu = &popupMenu{Title: title, Items: items, X: x, Y: y}
}

func MenuHandleEvent(event tcell.Event) bool {
	if activeMenu == nil {
		return false
	}

	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			CloseMenu()
			return true
		}
		if ev.Key() == tcell.KeyRune {
			n := int(ev.Rune() - '1')
			if n >= 0 && n < len(activeMenu.Items) {
				run := activeMenu.Items[n].Run
				CloseMenu()
				if run != nil {
					run()
				}
				return true
			}
		}
		return true
	case *tcell.EventMouse:
		btn := ev.Buttons()
		x, y := ev.Position()
		if btn == tcell.Button1 {
			for _, hit := range activeMenu.hits {
				if y == hit.Y && x >= hit.X0 && x < hit.X1 {
					run := hit.Run
					CloseMenu()
					if run != nil {
						run()
					}
					return true
				}
			}
			CloseMenu()
			return false
		}
		if btn == tcell.ButtonSecondary {
			CloseMenu()
			return false
		}
		return true
	default:
		return true
	}
}

func DisplayMenu() {
	if activeMenu == nil {
		return
	}

	style := config.DefStyle.Reverse(true)
	if s, ok := config.Colorscheme["statusline"]; ok {
		style = s
	}
	hintStyle := config.DefStyle
	if s, ok := config.Colorscheme["message"]; ok {
		hintStyle = s
	}

	w, h := screen.Screen.Size()
	width := runewidth.StringWidth(activeMenu.Title) + 4
	for _, item := range activeMenu.Items {
		n := runewidth.StringWidth(item.Label) + 6
		if n > width {
			width = n
		}
	}
	footer := "click or 1-9, Esc"
	if runewidth.StringWidth(footer)+4 > width {
		width = runewidth.StringWidth(footer) + 4
	}
	if width > w {
		width = w
	}
	height := len(activeMenu.Items) + 3
	if height > h {
		height = h
	}

	x, y := activeMenu.X, activeMenu.Y
	fillLine := func(row int, text string, st tcell.Style) {
		col := x
		for _, r := range text {
			if col >= x+width || col >= w {
				break
			}
			rw := runewidth.RuneWidth(r)
			screen.SetContent(col, row, r, nil, st)
			for i := 1; i < rw; i++ {
				screen.SetContent(col+i, row, ' ', nil, st)
			}
			col += rw
		}
		for col < x+width && col < w {
			screen.SetContent(col, row, ' ', nil, st)
			col++
		}
	}

	fillLine(y, " "+activeMenu.Title, style)
	activeMenu.hits = activeMenu.hits[:0]
	for i, item := range activeMenu.Items {
		row := y + 1 + i
		if row >= h {
			break
		}
		label := item.Label
		if i < 9 {
			label = string(rune('1'+i)) + " " + item.Label
		}
		fillLine(row, " "+label, style)
		activeMenu.hits = append(activeMenu.hits, hitRegion{
			X0:  x,
			X1:  x + width,
			Y:   row,
			Run: item.Run,
		})
	}
	fillLine(y+height-1, " "+footer, hintStyle)
}

func commandBarContainsY(y int) bool {
	infoY := InfoBar.GetView().Y
	top := infoY - config.CommandBarRows
	return y >= top && y < infoY
}

func drawCommandBar(width, infoY int, style tcell.Style) {
	btnStyle := config.DefStyle.Reverse(true)
	if s, ok := config.Colorscheme["statusline"]; ok {
		btnStyle = s
	}
	accentStyle := btnStyle.Reverse(true)
	if s, ok := config.Colorscheme["statusline.suggestions"]; ok {
		accentStyle = s
	} else if s, ok := config.Colorscheme["match-brace"]; ok {
		accentStyle = s
	}

	clearRow := func(y int) {
		for x := 0; x < width; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}

	commandBarHits = commandBarHits[:0]
	for i := 0; i < config.CommandBarRows; i++ {
		y := infoY - config.CommandBarRows + i
		clearRow(y)
		x := 1
		for _, btn := range commandBarRow(i) {
			label := " " + btn.Label + " "
			lw := runewidth.StringWidth(label)
			if x+lw >= width {
				break
			}
			st := btnStyle
			if btn.Accent {
				st = accentStyle
			}
			col := x
			for _, r := range label {
				screen.SetContent(col, y, r, nil, st)
				col += runewidth.RuneWidth(r)
			}
			commandBarHits = append(commandBarHits, hitRegion{
				X0:  x,
				X1:  x + lw,
				Y:   y,
				Run: btn.Run,
			})
			x += lw + 1
		}
	}
}

func CommandBarHandleEvent(event tcell.Event) bool {
	if !config.GetGlobalOption("keymenu").(bool) {
		return false
	}
	me, ok := event.(*tcell.EventMouse)
	if !ok {
		return false
	}

	x, y := me.Position()
	onBar := commandBarContainsY(y)

	if me.Buttons() == tcell.ButtonNone {
		if commandBarArmed {
			commandBarArmed = false
			if onBar && commandBarPressY == y {
				for _, hit := range commandBarHits {
					if y == hit.Y && commandBarPressX >= hit.X0 && commandBarPressX < hit.X1 &&
						x >= hit.X0 && x < hit.X1 {
						if hit.Run != nil {
							hit.Run()
						}
						return true
					}
				}
			}
			return onBar
		}
		return false
	}

	if me.Buttons() != tcell.Button1 || !onBar {
		return false
	}
	if !commandBarArmed {
		commandBarArmed = true
		commandBarPressX = x
		commandBarPressY = y
	}
	return true
}

// HandleOverlayEvent consumes clicks on the command bar or an open popup.
// Resize events are left for the rest of the editor after the menu is closed.
func HandleOverlayEvent(event tcell.Event, hasPrompt bool) bool {
	if _, ok := event.(*tcell.EventResize); ok {
		CloseMenu()
		return false
	}
	if SplashHandleEvent(event) {
		return true
	}
	if MenuHandleEvent(event) {
		return true
	}
	if hasPrompt {
		return false
	}
	return CommandBarHandleEvent(event)
}
