# Command bar

The command bar is opened by pressing `Ctrl-e`. It is a single-line buffer,
meaning that all keybindings from a normal buffer are supported (as well
as mouse and selection).

When running a command, you can use extra syntax that micro will expand before
running the command. To use an argument with a space in it, put it in
quotes. The command bar parser uses the same rules for parsing arguments that
`/bin/sh` would use (single quotes, double quotes, escaping). The command bar
does not look up environment variables.

# Commands

Micro provides the following commands that can be executed at the command-bar
by pressing `Ctrl-e` and entering the command. Arguments are placed in single
quotes here but these are not necessary when entering the command in micro.

* `bind 'key' 'action'`: creates a keybinding from key to action. See the
   `keybindings` documentation for more information about binding keys.
   This command will modify `bindings.json` and overwrite any bindings to
   `key` that already exist.

* `help ['topic'] ['flags']`: opens the corresponding help topics.
   If no topic is provided opens the default help screen. If multiple topics are
   provided (separated via ` `) they are opened all as splits.
   Help topics are stored as `.md` files in the `runtime/help` directory of
   the source tree, which is embedded in the final binary.
   The `flags` are optional.
   * `-hsplit`: Opens the help topic in a horizontal split
   * `-vsplit`: Opens the help topic in a vertical split

   The default split type is defined by the global `helpsplit` option.

* `save ['filename']`: saves the current buffer. If the file is provided it
   will 'save as' the filename.

* `quit`: quits micro.

* `goto 'line[:col]'`: goes to the given absolute line (and optional column)
   number.
   A negative number can be passed to go inward from the end of the file.
   Example: -5 goes to the 5th-last line in the file.

* `jump 'line[:col]'`: goes to the given relative number from the current
   line (and optional absolute column) number.
   Example: -5 jumps 5 lines up in the file, while (+)3 jumps 3 lines down.

* `replace 'search' 'value' ['flags']`: This will replace `search` with `value`.
   The `flags` are optional. Possible flags are:
   * `-a`: Replace all occurrences at once
   * `-l`: Do a literal search instead of a regex search

   Note that `search` must be a valid regex (unless `-l` is passed). If one
   of the arguments does not have any spaces in it, you may omit the quotes.

   In case the search is done non-literal (without `-l`), the 'value'
   is interpreted as a template:
   * `$3` or `${3}` substitutes the submatch of the 3rd (capturing group)
   * `$foo` or `${foo}` substitutes the submatch of the (?P<foo>named group)
   * You have to write `$$` to substitute a literal dollar.

* `replaceall 'search' 'value'`: this will replace all occurrences of `search`
   with `value` without user confirmation.

   See `replace` command for more information.

* `set 'option' 'value'`: sets the option to value. See the `options` help
   topic for a list of options you can set. This will modify your
   `settings.json` with the new value.

* `setlocal 'option' 'value'`: sets the option to value locally (only in the
   current buffer). This will *not* modify `settings.json`.

* `toggle 'option'`: toggles the option. Only works with options that accept
   exactly two values. This will modify your `settings.json` with the new value.

* `togglelocal 'option'`: toggles the option locally (only in the
   current buffer). Only works with options that accept exactly two values.
   This will *not* modify `settings.json`.

* `reset 'option'`: resets the given option to its default value.

* `show 'option'`: shows the current value of the given option.

* `showkey 'key'`: Show the action(s) bound to a given key. For example
   running `> showkey Ctrl-c` will display `Copy`.

* `run 'sh-command'`: runs the given shell command in the background. The
   command's output will be displayed in one line when it finishes running.

* `vsplit ['filename']`: opens a vertical split with `filename`. If no filename
   is provided, a vertical split is opened with an empty buffer. If multiple
   files are provided (separated via ` `) they are opened all as splits.

* `hsplit ['filename']`: same as `vsplit` but opens a horizontal split instead
   of a vertical split.

* `tab ['filename']`: opens the given file in a new tab. If no filename
   is provided, a tab is opened with an empty buffer. If multiple files are
   provided (separated via ` `) they are opened all as tabs.

* `tabmove '[-+]n'`: Moves the active tab to another slot. `n` is an integer.
   If `n` is prefixed with `-` or `+`, then it represents a relative position
   (e.g. `tabmove +2` moves the tab to the right by `2`). If `n` has no prefix,
   it represents an absolute position (e.g. `tabmove 2` moves the tab to slot `2`).

* `tabswitch 'tab'`: This command will switch to the specified tab. The `tab`
   can either be a tab number, or a name of a tab.

* `textfilter 'sh-command'`: filters the current selection through a shell
   command as standard input and replaces the selection with the stdout of
   the shell command.  For example, to sort a list of numbers, first select
   them, and then execute `> textfilter sort -n`.

* `log`: opens a log of all messages and debug statements.

* `plugin list`: lists all installed plugins.

* `plugin install 'pl'`: install a plugin.

* `plugin remove 'pl'`: remove a plugin.

* `plugin update ['pl']`: update a plugin (if no arguments are provided
   updates all plugins).

* `plugin search 'pl'`: search available plugins for a keyword.

* `plugin available`: show available plugins that can be installed.

* `reload`: reloads all runtime files (settings, keybindings, syntax files,
   colorschemes, plugins). All plugins will be unloaded by running their
   `deinit()` function (if it exists), and then loaded again by calling the
   `preinit()`, `init()` and `postinit()` functions (if they exist).

* `cd 'path'`: Change the working directory to the given `path`.

* `pwd`: Print the current working directory.

* `open 'filename'`: Open a file in the current buffer.

* `reopen`: Reopens the current file from disk.

* `retab`: Replaces all leading tabs with spaces or leading spaces with tabs
   depending on the value of `tabstospaces`.

* `raw`: micro will open a new tab and show the escape sequence for every event
   it receives from the terminal. This shows you what micro actually sees from
   the terminal and helps you see which bindings aren't possible and why. This
   is most useful for debugging keybindings.

* `term ['exec']`: Open a terminal emulator running the given executable. If no
   executable is given, this will open the default shell in the terminal
   emulator. This replaces the current pane (or opens a new tab when only one
   pane is open).

* `terminal ['exec']`: Open a terminal in a horizontal split without replacing
   the current file. Click **Terminal** on the bottom bar, press `F9`, or
   right-click the editor.

* `agent ['exec']`: Open an agent panel as a horizontal split. Runs `herdr` if
   it is on `PATH`, otherwise a labeled shell. Click **Agent** or press `F6`.

* `closepanel`: Close the focused extra panel (terminal, agent, or file tree).
   Click **Close panel** on the bottom bar.

* `theme ['name']`: Open the Theme menu, or apply a colorscheme such as
   `catppuccin` or `tokyo-night`.

* `options`: Open the Options menu (same as the toolbar **Options** button).

* `layout ['right'|'left'|'below'|'above'|'vsplit'|'hsplit']`: Arrange panes.
   With no argument, opens the Layout menu. `right`/`left`/`below`/`above`
   move the focused pane. `vsplit` and `hsplit` create a new empty split.

---

The following commands are provided by the default plugins:

* `lint`: Lint the current file for errors.
* `comment`: automatically comment or uncomment current selection or line.
