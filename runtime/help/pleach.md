# Pleach

Pleach is mouse-first. Most layout actions are on the bottom command bar or
on a right-click menu. You do not need vim chords.

## Toolbar

The three lines above the status/info line are a toolbar. Click them.

The first button is **Theme: …** and always shows the current colorscheme.
Click it to switch among the Herdr palettes.

Menus:

* **File** — New, Open, Save, Save as, Reload, Quit
* **Edit** — Undo, Redo, Cut, Copy, Paste, Find, Replace, Select all
* **View** — splits, layout, ruler, wrap, syntax, toolbar
* **Panels** — tree, terminal, agent, close panel
* **Options** — indent, tabs, search, mouse, and more settings
* **Help** — help topics and the command bar

Shortcuts on the next rows: New, Open, Save, Find, Undo, Redo, Quit, Tree,
Terminal, Agent, Close panel, Side by side, Stacked, Command.

Hide the toolbar with `Alt-g`, or `> set keymenu false`.

## File tree

The explorer stays on the left, like VS Code:

* Click a folder to expand or collapse it
* Click a file to open it in the editor (not a new split)
* Click **↑ ..** to go up one folder
* Right-click a folder for Expand / Open folder / New file / Rename / Delete

Right-click a file or folder for:

* Open
* New file
* New folder
* Rename
* Delete
* Copy path
* Close tree

`Ctrl-b` or `F8` also toggles the tree.

## Arranging panes

Click a pane, then:

* **Side by side** or **Stacked** to open a new split
* **Layout** or right-click to **move this pane** left, right, above, or below
* Drag the `|` or `-` divider between panes to resize

A terminal can sit below the editor or beside it: use **Layout → Terminal on the right**, or right-click the terminal and move it.

## Terminal and agent

* Click **Terminal** or press `F9` (opens below)
* Click **Agent** or press `F6` — this opens [Herdr](https://herdr.dev)
* If herdr is not installed, Pleach downloads the official binary to
  `~/.config/pleach/bin/herdr`
* Right-click a terminal to move it, close it, or switch splits

`> terminal` and `> agent` do the same thing from the command prompt.
`> term` still replaces the current pane (Micro's original behavior).

## Editor right-click

Right-click in a file to move this pane, open a new split, or start a terminal
on the right or below.
