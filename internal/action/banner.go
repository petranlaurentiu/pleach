package action

import (
	"strings"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"
)

var bannerBig = map[rune][]string{
	' ': {"     ", "     ", "     ", "     ", "     "},
	'A': {" ███ ", "█   █", "█████", "█   █", "█   █"},
	'B': {"████ ", "█   █", "████ ", "█   █", "████ "},
	'C': {" ████", "█    ", "█    ", "█    ", " ████"},
	'D': {"████ ", "█   █", "█   █", "█   █", "████ "},
	'E': {"█████", "█    ", "████ ", "█    ", "█████"},
	'F': {"█████", "█    ", "████ ", "█    ", "█    "},
	'G': {" ████", "█    ", "█ ███", "█   █", " ███ "},
	'H': {"█   █", "█   █", "█████", "█   █", "█   █"},
	'I': {"█████", "  █  ", "  █  ", "  █  ", "█████"},
	'J': {"█████", "   █ ", "   █ ", "█  █ ", " ██  "},
	'K': {"█   █", "█  █ ", "███  ", "█  █ ", "█   █"},
	'L': {"█    ", "█    ", "█    ", "█    ", "█████"},
	'M': {"█   █", "██ ██", "█ █ █", "█   █", "█   █"},
	'N': {"█   █", "██  █", "█ █ █", "█  ██", "█   █"},
	'O': {" ███ ", "█   █", "█   █", "█   █", " ███ "},
	'P': {"████ ", "█   █", "████ ", "█    ", "█    "},
	'Q': {" ███ ", "█   █", "█   █", "█  ██", " ████"},
	'R': {"████ ", "█   █", "████ ", "█  █ ", "█   █"},
	'S': {" ████", "█    ", " ███ ", "    █", "████ "},
	'T': {"█████", "  █  ", "  █  ", "  █  ", "  █  "},
	'U': {"█   █", "█   █", "█   █", "█   █", " ███ "},
	'V': {"█   █", "█   █", "█   █", " █ █ ", "  █  "},
	'W': {"█   █", "█   █", "█ █ █", "██ ██", "█   █"},
	'X': {"█   █", " █ █ ", "  █  ", " █ █ ", "█   █"},
	'Y': {"█   █", " █ █ ", "  █  ", "  █  ", "  █  "},
	'Z': {"█████", "   █ ", "  █  ", " █   ", "█████"},
	'-': {"     ", "     ", "█████", "     ", "     "},
	'.': {"     ", "     ", "     ", "     ", "  █  "},
}

var bannerMini = map[rune][]string{
	' ': {"   ", "   ", "   "},
	'A': {"█▀█", "█▀█", "█ █"},
	'B': {"█▀▄", "█▀▄", "█▀ "},
	'C': {"█▀▀", "█  ", "█▀▀"},
	'D': {"█▀▄", "█ █", "█▀ "},
	'E': {"█▀▀", "█▀ ", "█▀▀"},
	'F': {"█▀▀", "█▀ ", "█  "},
	'G': {"█▀▀", "█ █", "█▀█"},
	'H': {"█ █", "█▀█", "█ █"},
	'I': {"▀█▀", " █ ", "▄█▄"},
	'J': {"▀▀█", "  █", "▀▀ "},
	'K': {"█ █", "█▀ ", "█ █"},
	'L': {"█  ", "█  ", "█▀▀"},
	'M': {"█▄█", "█▀█", "█ █"},
	'N': {"█▄█", "█ █", "█ █"},
	'O': {"█▀█", "█ █", "█▀█"},
	'P': {"█▀█", "█▀▀", "█  "},
	'Q': {"█▀█", "█ █", "▀█▄"},
	'R': {"█▀█", "█▀▄", "█ █"},
	'S': {"█▀▀", "▀▀█", "▀▀█"},
	'T': {"▀█▀", " █ ", " █ "},
	'U': {"█ █", "█ █", "█▀█"},
	'V': {"█ █", "█ █", " ▀ "},
	'W': {"█ █", "█▄█", "▀▄▀"},
	'X': {"█ █", " █ ", "█ █"},
	'Y': {"█ █", " ▀ ", " █ "},
	'Z': {"▀▀█", " █ ", "█▀▀"},
	'-': {"   ", "▀▀▀", "   "},
	'.': {"   ", "   ", " █ "},
	',': {"   ", "   ", " █ "},
	'?': {"▀█▀", "▄▀ ", " ▄ "},
	'!': {" █ ", " █ ", " ▄ "},
	'\'': {" █ ", "   ", "   "},
}

func renderBanner(text string, cols int) []string {
	text = strings.ToUpper(text)
	if lines := composeBanner(text, bannerBig, 1); runewidth.StringWidth(lines[0]) <= cols {
		return lines
	}
	if lines := composeBanner(text, bannerMini, 1); runewidth.StringWidth(lines[0]) <= cols {
		return lines
	}
	return []string{text}
}

func renderSubtitle(text string, cols int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if cols < 8 || runewidth.StringWidth(text) <= cols {
		return []string{text}
	}
	words := strings.Fields(text)
	var out []string
	var line string
	for _, word := range words {
		if line == "" {
			line = word
			continue
		}
		trial := line + " " + word
		if runewidth.StringWidth(trial) > cols {
			out = append(out, line)
			line = word
			continue
		}
		line = trial
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}

func composeBanner(text string, glyphs map[rune][]string, gap int) []string {
	var sample []string
	for _, g := range glyphs {
		sample = g
		break
	}
	h := len(sample)
	lines := make([]string, h)
	first := true
	for _, r := range text {
		if r > unicode.MaxASCII {
			continue
		}
		g, ok := glyphs[r]
		if !ok {
			g = glyphs[' ']
		}
		if g == nil {
			continue
		}
		for i := 0; i < h && i < len(g); i++ {
			if !first {
				lines[i] += strings.Repeat(" ", gap)
			}
			lines[i] += g[i]
		}
		first = false
	}
	return lines
}
