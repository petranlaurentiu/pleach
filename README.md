# Pleach

A mouse-first terminal IDE. Start it like Neovim, click like VS Code.

```bash
cd pleach
make build
./pleach .
```

- `pleach .` opens the current folder with a file tree on the left
- Click a folder to expand it, or a file to open it in the editor; `Ctrl+S` saves; `Ctrl+B` or `F8` toggles the tree
- The bottom toolbar has File / Edit / View / Panels / Options / Help, plus New, Open, Save, Find, Undo, Redo, Quit
- **Tree** brings the explorer back; **Quit** closes Pleach
- **Side by side** splits left/right; **Stacked** splits top/bottom
- **Panels** or right-click a pane to move it left, right, above, or below
- Right-click a file or folder for Open / New file / New folder / Rename / Delete / Copy path
- `F9` opens a terminal split; **Agent** / `F6` opens [Herdr](https://herdr.dev)
- The first toolbar button is **Theme: …** — click it to change colorschemes
- Config lives in `~/.config/pleach` (or `$PLEACH_CONFIG_HOME`)

Pleach is a rebranded fork of [Micro](https://github.com/micro-editor/micro) (MIT). See [NOTICE](NOTICE) and [LICENSE](LICENSE).
