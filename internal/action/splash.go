package action

import (
	"math/rand"
	"strings"
	"sync"
	"time"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/micro-editor/micro/v2/internal/config"
	"github.com/micro-editor/micro/v2/internal/screen"
	"github.com/micro-editor/tcell/v2"
)

// Compact JARVIS wordmark from https://ascii.co.uk/art/jarvis
const jarvisSmall = `
   _                  _
  (_)                (_)
   _  __ _ _ ____   ___ ___
  | |/ _` + "`" + ` | '__\ \ / / / __|
  | | (_| | |   \ V /| \__ \
  | |\__,_|_|    \_/ |_|___/
 _/ |
|__/`

const matrixGlyphs = "01ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ"

var matrixRunes = []rune(matrixGlyphs)

var (
	splashMu     sync.Mutex
	splashActive bool
	splashStop   chan struct{}
	rainCols     []matrixCol
	rainTick     int
	rainWidth    int
	rainHeight   int
	splashStart  time.Time
)

type matrixCol struct {
	head  int
	speed int
	trail int
}

func artLines(art string) []string {
	lines := strings.Split(strings.Trim(art, "\n"), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return lines
}

func artSize(lines []string) (w, h int) {
	h = len(lines)
	for _, line := range lines {
		if n := runewidth.StringWidth(line); n > w {
			w = n
		}
	}
	return w, h
}

func splashStyle(group string, fallback tcell.Style) tcell.Style {
	if s, ok := config.Colorscheme[group]; ok {
		return s
	}
	return fallback
}

func CloseSplash() {
	splashMu.Lock()
	if !splashActive {
		splashMu.Unlock()
		return
	}
	splashActive = false
	stop := splashStop
	splashStop = nil
	splashMu.Unlock()
	if stop != nil {
		close(stop)
	}
	screen.Redraw()
}

// ShowWelcome opens a full-screen intro and keeps it until Enter.
func ShowWelcome() {
	splashMu.Lock()
	splashActive = true
	splashStart = time.Now()
	if splashStop != nil {
		close(splashStop)
	}
	stop := make(chan struct{})
	splashStop = stop
	rainCols = nil
	rainTick = 0
	splashMu.Unlock()

	go func() {
		t := time.NewTicker(70 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				screen.Redraw()
			}
		}
	}()
}

func SplashHandleEvent(event tcell.Event) bool {
	splashMu.Lock()
	active := splashActive
	splashMu.Unlock()
	if !active {
		return false
	}
	if _, ok := event.(*tcell.EventResize); ok {
		return false
	}
	if isEnter(event) {
		CloseSplash()
	}
	return true
}

func isEnter(event tcell.Event) bool {
	k, ok := event.(*tcell.EventKey)
	if !ok {
		return false
	}
	return isEnterKey(k.Key(), k.Rune())
}

func isEnterKey(key tcell.Key, r rune) bool {
	return key == tcell.KeyEnter || r == '\n' || r == '\r'
}

func resetRain(w, h int) {
	rainWidth, rainHeight = w, h
	rainCols = make([]matrixCol, w)
	for x := range rainCols {
		rainCols[x] = matrixCol{
			head:  rand.Intn(h + 8),
			speed: 1 + rand.Intn(3),
			trail: 6 + rand.Intn(14),
		}
	}
}

func stepRain(w, h int) {
	if len(rainCols) != w || rainHeight != h {
		resetRain(w, h)
	}
	rainTick++
	for x := range rainCols {
		col := &rainCols[x]
		if rainTick%col.speed != 0 {
			continue
		}
		col.head++
		if col.head-col.trail > h {
			col.head = -rand.Intn(h/2 + 1)
			col.speed = 1 + rand.Intn(3)
			col.trail = 6 + rand.Intn(14)
		}
	}
}

func matrixRune(x, y, tick int) rune {
	n := len(matrixRunes)
	i := (x*31 + y*17 + tick/3) % n
	if i < 0 {
		i += n
	}
	return matrixRunes[i]
}

func drawRain(w, h int, head, mid, tail tcell.Style) {
	for x, col := range rainCols {
		for d := 0; d < col.trail; d++ {
			y := col.head - d
			if y < 0 || y >= h {
				continue
			}
			st := mid
			switch {
			case d == 0:
				st = head
			case d > col.trail/2:
				st = tail
			}
			screen.SetContent(x, y, matrixRune(x, y, rainTick), nil, st)
		}
	}
}

func drawCenteredLines(lines []string, y, w, h int, st tcell.Style) int {
	for i, text := range lines {
		row := y + i
		if row < 0 || row >= h {
			continue
		}
		x := (w - runewidth.StringWidth(text)) / 2
		if x < 0 {
			x = 0
		}
		col := x
		for _, r := range text {
			if col >= w {
				break
			}
			if r != ' ' {
				screen.SetContent(col, row, r, nil, st)
			}
			col += runewidth.RuneWidth(r)
		}
	}
	return y + len(lines)
}

func DisplaySplash() {
	splashMu.Lock()
	active := splashActive
	splashMu.Unlock()
	if !active || screen.Screen == nil {
		return
	}

	w, h := screen.Screen.Size()
	bg := config.DefStyle.Foreground(tcell.ColorDarkGreen).Background(tcell.ColorBlack)
	if s, ok := config.Colorscheme["default"]; ok {
		_, bgColor, _ := s.Decompose()
		bg = s.Foreground(tcell.ColorDarkGreen).Background(bgColor)
	}
	head := config.DefStyle.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	mid := config.DefStyle.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	tail := config.DefStyle.Foreground(tcell.ColorDarkGreen).Background(tcell.ColorBlack)
	artStyle := splashStyle("comment", config.DefStyle.Foreground(tcell.ColorGray))
	bannerStyle := splashStyle("constant.number", config.DefStyle.Foreground(tcell.ColorGold))
	nameStyle := splashStyle("todo", config.DefStyle.Foreground(tcell.ColorYellow))
	msgStyle := splashStyle("identifier", config.DefStyle).Bold(true)
	hintStyle := splashStyle("comment", config.DefStyle.Dim(true))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			screen.SetContent(x, y, ' ', nil, bg)
		}
	}

	splashMu.Lock()
	stepRain(w, h)
	drawRain(w, h, head, mid, tail)
	splashMu.Unlock()

	art := artLines(jarvisSmall)
	welcome := renderBanner("WELCOME", w-2)
	name := renderBanner(bannerName(welcomeName()), w-2)
	splashMu.Lock()
	started := splashStart
	splashMu.Unlock()
	if started.IsZero() {
		started = time.Now()
	}
	subtitle := renderSubtitle(welcomeMessage(time.Now(), welcomeName(), time.Since(started)), w-4)
	hint := "press Enter"

	blockH := len(art) + 1 + len(welcome) + 1 + len(name) + 1 + len(subtitle) + 3
	y := (h - blockH) / 2
	if y < 1 {
		y = 1
	}

	y = drawCenteredLines(art, y, w, h, artStyle)
	y += 1
	y = drawCenteredLines(welcome, y, w, h, bannerStyle)
	y += 1
	y = drawCenteredLines(name, y, w, h, nameStyle)
	y += 1
	drawCenteredLines(subtitle, y, w, h, msgStyle)
	drawCenteredLines([]string{hint}, h-2, w, h, hintStyle)
}

func bannerName(name string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '.' || r == ' ' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "SIR"
	}
	return b.String()
}
