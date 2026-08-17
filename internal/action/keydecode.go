package action

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/micro-editor/tcell/v2"
)

func decodeExtendedKey(esc string) (KeyEvent, bool) {
	if ke, ok := decodeKittyCSIU(esc); ok {
		return ke, true
	}
	return decodeXterm27(esc)
}

func decodeKittyCSIU(esc string) (KeyEvent, bool) {
	if !strings.HasPrefix(esc, "\x1b[") || !strings.HasSuffix(esc, "u") {
		return KeyEvent{}, false
	}
	body := strings.TrimSuffix(esc[2:], "u")
	if body == "" || !isKittyBody(body) {
		return KeyEvent{}, false
	}
	parts := strings.Split(body, ";")
	key, err := strconv.Atoi(strings.Split(parts[0], ":")[0])
	if err != nil || key <= 0 {
		return KeyEvent{}, false
	}
	mods := 1
	if len(parts) >= 2 {
		if n, err := strconv.Atoi(strings.Split(parts[1], ":")[0]); err == nil && n > 0 {
			mods = n
		}
	}
	return kittyToKeyEvent(key, mods), true
}

func isKittyBody(body string) bool {
	for _, r := range body {
		if r != ';' && r != ':' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func decodeXterm27(esc string) (KeyEvent, bool) {
	var mods, key int
	if _, err := fmt.Sscanf(esc, "\x1b[27;%d;%d~", &mods, &key); err != nil {
		return KeyEvent{}, false
	}
	if key <= 0 || mods <= 0 {
		return KeyEvent{}, false
	}
	return kittyToKeyEvent(key, mods), true
}

func kittyToKeyEvent(key, mods int) KeyEvent {
	mod := kittyModsToTcell(mods)
	switch key {
	case 127, 8:
		return KeyEvent{code: tcell.KeyBackspace2, mod: mod}
	case 13:
		return KeyEvent{code: tcell.KeyEnter, mod: mod}
	case 9:
		return KeyEvent{code: tcell.KeyTab, mod: mod}
	case 27:
		return KeyEvent{code: tcell.KeyEsc, mod: mod}
	case 3, 57349:
		return KeyEvent{code: tcell.KeyDelete, mod: mod}
	}
	if key >= 'a' && key <= 'z' && mod&tcell.ModCtrl != 0 && mod&tcell.ModMeta == 0 && mod&tcell.ModAlt == 0 {
		return KeyEvent{code: tcell.Key(key - 'a' + 1), mod: tcell.ModCtrl}
	}
	if key >= 32 && key < 127 {
		return KeyEvent{code: tcell.KeyRune, mod: mod, r: rune(key)}
	}
	return KeyEvent{code: tcell.KeyRune, mod: mod, r: rune(key)}
}

func kittyModsToTcell(mods int) tcell.ModMask {
	bits := mods - 1
	var m tcell.ModMask
	if bits&1 != 0 {
		m |= tcell.ModShift
	}
	if bits&2 != 0 {
		m |= tcell.ModAlt
	}
	if bits&4 != 0 {
		m |= tcell.ModCtrl
	}
	if bits&8 != 0 {
		m |= tcell.ModMeta
	}
	return m
}
