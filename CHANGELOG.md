CHANGELOG
=========

0.65.1
------
- Fixed incorrect `$FZF_CLICK_HEADER_WORD` and `$FZF_CLICK_FOOTER_WORD` when the header or footer contains ANSI escape sequences and tab characters.
- Fixed a bug where you cannot unset the default `--nth` using `change-nth` action.
- Fixed a highlighting bug when using `--color fg:dim,nth:regular` pattern over ANSI-colored items.

0.65.0
------
- Added `click-footer` event that is triggered when the footer section is clicked. When the event is triggered, the following environment variables are set:
    - `$FZF_CLICK_FOOTER_COLUMN` - clicked column (1-based)
    - `$FZF_CLICK_FOOTER_LINE` - clicked line (1-based)
    - `$FZF_CLICK_FOOTER_WORD` - the word under the cursor
  ```sh
  fzf --footer $'[Edit] [View]\n[Copy to clipboard]' \
      --with-shell 'bash -c' \
      --bind 'click-footer:transform:
        [[ $FZF_CLICK_FOOTER_WORD =~ Edit ]] && echo "execute:vim \{}"
        [[ $FZF_CLICK_FOOTER_WORD =~ View ]] && echo "execute:view \{}"
        (( FZF_CLICK_FOOTER_LINE == 2 )) && (( FZF_CLICK_FOOTER_COLUMN < 20 )) &&
            echo "execute-silent(echo -n \{} | pbcopy)+bell"
      '
  ```
- Added `trigger(...)` action that triggers events bound to another key or event.
  ```sh
  # You can click on each key name to trigger the actions bound to that key
  fzf --footer 'Ctrl-E: Edit / Ctrl-V: View / Ctrl-Y: Copy to clipboard' \
      --with-shell 'bash -c' \
      --bind 'ctrl-e:execute:vim {}' \
      --bind 'ctrl-v:execute:view {}' \
      --bind 'ctrl-y:execute-silent(echo -n {} | pbcopy)+bell' \
      --bind 'click-footer:transform:
        [[ $FZF_CLICK_FOOTER_WORD =~ Ctrl ]] && echo "trigger(${FZF_CLICK_FOOTER_WORD%:})"
      '
  ```
    - You can specify a series of keys and events
      ```sh
      fzf --bind 'a:up,b:trigger(a,a,a)'
      ```
- Added support for `{*n}` and `{*nf}` placeholder.
    - `{*n}` evaluates to the zero-based ordinal index of all matched items.
    - `{*nf}` evaluates to the temporary file containing that.
- Bug fixes and improvements
    - [neovim] Fixed margin background color when `&winborder` is used (#4453)
    - Fixed rendering error when hiding a preview window without border (#4465)
    - fix(shell): check for mawk existence before version check (#4468)
        - Thanks to @LangLangBart and @akinomyoga
    - Fixed `--no-header-lines-border` behavior (08027e7a)

0.64.0
------
- Added `multi` event that is triggered when the multi-selection has changed.
  ```sh
  fzf --multi \
      --bind 'ctrl-a:select-all,ctrl-d:deselect-all' \
      --bind 'multi:transform-footer:(( FZF_SELECT_COUNT )) && echo "Selected $FZF_SELECT_COUNT item(s)"'
  ```
- [Halfwidth and fullwidth alphanumeric and punctuation characters](https://en.wikipedia.org/wiki/Halfwidth_and_Fullwidth_Forms_(Unicode_block)) are now internally normalized to their ASCII equivalents to allow matching with ASCII queries.
  ```sh
  echo ＡＢＣ| fzf -q abc
  ```
- Renamed `clear-selection` action to `clear-multi` for consistency.
    - `clear-selection` remains supported as an alias for backward compatibility.
- Bug fixes
    - Fixed a bug that could cause fzf to abort due to incorrect update ordering.
    - Fixed a bug where some multi-selections were lost when using `exclude` or `change-nth`.

0.63.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.63.0/_

- Added footer. The default border style for footer is `line`, which draws a single separator line.
  ```sh
  fzf --reverse --footer "fzf: friend zone forever"
  ```
  - Options
      - `--footer[=STRING]`
      - `--footer-border[=STYLE]`
      - `--footer-label=LABEL`
      - `--footer-label-pos=COL[:bottom]`
  - Colors
      - `footer`
      - `footer-bg`
      - `footer-border`
      - `footer-label`
  - Actions
      - `change-footer`
      - `transform-footer`
      - `bg-transform-footer`
      - `change-footer-label`
      - `transform-footer-label`
      - `bg-transform-footer-label`
- `line` border style is now allowed for all types of border except for `--list-border`.
  ```sh
  fzf --height 50% --style full:line --preview 'cat {}' \
      --bind 'focus:bg-transform-header(file {})+bg-transform-footer(wc {})'
  ```
- Added `{*}` placeholder flag that evaluates to all matched items.
  ```bash
  seq 10000 | fzf --preview "awk '{sum += \$1} END {print sum}' {*f}"
  ```
  - Use this with caution, as it can make fzf sluggish for large lists.
- Added asynchronous transform actions with `bg-` prefix that run asynchronously in the background, along with `bg-cancel` action to cancel currently running `bg-transform` actions.
  ```sh
  # Implement popup that disappears after 1 second
  #   * Use footer as the popup
  #   * Use `bell` to ring the terminal bell
  #   * Use `bg-transform-footer` to clear the footer after 1 second
  #   * Use `bg-cancel` to cancel currently running background transform actions
  fzf --multi --list-border \
      --bind 'enter:execute-silent(echo -n {+} | pbcopy)+bell' \
      --bind 'enter:+transform-footer(echo Copied {} to clipboard)' \
      --bind 'enter:+bg-cancel+bg-transform-footer(sleep 1)'

  # It's okay for the commands to take a little while because they run in the background
  GETTER='curl -s http://metaphorpsum.com/sentences/1'
  fzf --style full --border --preview : \
      --bind "focus:bg-transform-header:$GETTER" \
      --bind "focus:+bg-transform-footer:$GETTER" \
      --bind "focus:+bg-transform-border-label:$GETTER" \
      --bind "focus:+bg-transform-preview-label:$GETTER" \
      --bind "focus:+bg-transform-input-label:$GETTER" \
      --bind "focus:+bg-transform-list-label:$GETTER" \
      --bind "focus:+bg-transform-header-label:$GETTER" \
      --bind "focus:+bg-transform-footer-label:$GETTER" \
      --bind "focus:+bg-transform-ghost:$GETTER" \
      --bind "focus:+bg-transform-prompt:$GETTER"
  ```
- Added support for full-line background color in the list section
  ```sh
  for i in $(seq 16 255); do
    echo -e "\x1b[48;5;${i}m\x1b[0Khello"
  done | fzf --ansi
  ```
- SSH completion enhancements by @akinomyoga
- Bug fixes and improvements

0.62.0
------
- Relaxed the `--color` option syntax to allow whitespace-separated entries (in addition to commas), making multi-line definitions easier to write and read
  ```sh
  # seoul256-light
  fzf --style full --color='
    fg:#616161 fg+:#616161
    bg:#ffffff bg+:#e9e9e9 alt-bg:#f1f1f1
    hl:#719872 hl+:#719899
    pointer:#e12672 marker:#e17899
    header:#719872
    spinner:#719899 info:#727100
    prompt:#0099bd query:#616161
    border:#e1e1e1
  '
  ```
- Added `alt-bg` color to create striped lines to visually separate rows
  ```sh
  fzf --color bg:237,alt-bg:238,current-bg:236 --highlight-line

  declare -f | perl -0777 -pe 's/^}\n/}\0/gm' |
    bat --plain --language bash --color always |
    fzf --read0 --ansi --reverse --multi \
        --color bg:237,alt-bg:238,current-bg:236 --highlight-line
  ```
- [fish] Improvements in CTRL-R binding (@bitraid)
    - You can trigger CTRL-R in the middle of a command to insert the selected item
    - You can delete history items with SHIFT-DEL
- Bug fixes and improvements
    - Fixed unnecessary 100ms delay after `reload` (#4364)
    - Fixed `selected-bg` not applied to colored items (#4372)

0.61.3
------
- Reverted #4351 as it caused `tmux run-shell 'fzf --tmux'` to fail (#4559 #4560)
- More environment variables for child processes (#4356)

0.61.2
------
- Fixed panic when using header border without pointer/marker (@phanen)
- Fixed `--tmux` option when already inside a tmux popup (@peikk0)
- Bug fixes and improvements in CTRL-T binding of fish (#4334) (@bitraid)
- Added `--no-tty-default` option to make fzf search for the current TTY device instead of defaulting to `/dev/tty` (#4242)

0.61.1
------
- Disable bracketed-paste mode on exit. This fixes issue where pasting breaks after running fzf on old bash versions that don't support the mode.

0.61.0
------
- Added `--ghost=TEXT` to display a ghost text when the input is empty
  ```sh
  # Display "Type to search" when the input is empty
  fzf --ghost "Type to search"
  ```
- Added `change-ghost` and `transform-ghost` actions for dynamically changing the ghost text
- Added `change-pointer` and `transform-pointer` actions for dynamically changing the pointer sign
- Added `r` flag for placeholder expression (raw mode) for unquoted output
- Bug fixes and improvements

0.60.3
------
- Bug fixes and improvements
    - [fish] Enable multiple history commands insertion (#4280) (@bitraid)
    - [walker] Append '/' to directory entries on MSYS2 (#4281)
    - Trim trailing whitespaces after processing ANSI sequences (#4282)
    - Remove temp files before `become` when using `--tmux` option (#4283)
    - Fix condition for using item numlines cache (#4285) (@alex-huff)
    - Make `--accept-nth` compatible with `--select-1` (#4287)
    - Increase the query length limit from 300 to 1000 (#4292)
    - [windows] Prevent fzf from consuming user input while paused (#4260)

0.60.2
------
- Template for `--with-nth` and `--accept-nth` now supports `{n}` which evaluates to the zero-based ordinal index of the item
- Fixed a regression that caused the last field in the "nth" expression to be trimmed when a regular expression delimiter is used
    - Thanks to @phanen for the fix
- Fixed 'jump' action when the pointer is an empty string

0.60.1
------
- Bug fixes and minor improvements
    - Built-in walker now prints directory entries with a trailing slash
    - Fixed a bug causing unexpected behavior with [fzf-tab](https://github.com/Aloxaf/fzf-tab). Please upgrade if you use it.
- Thanks to @alexeisersun, @bitraid, @Lompik, and @fsc0 for the contributions

0.60.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.60.0/_

- Added `--accept-nth` for choosing output fields
  ```sh
  ps -ef | fzf --multi --header-lines 1 | awk '{print $2}'
  # Becomes
  ps -ef | fzf --multi --header-lines 1 --accept-nth 2

  git branch | fzf | cut -c3-
  # Can be rewritten as
  git branch | fzf --accept-nth -1
  ```
- `--accept-nth` and `--with-nth` now support a template that includes multiple field index expressions in curly braces
  ```sh
  echo foo,bar,baz | fzf --delimiter , --accept-nth '{1}, {3}, {2}'
    # foo, baz, bar

  echo foo,bar,baz | fzf --delimiter , --with-nth '{1},{3},{2},{1..2}'
    # foo,baz,bar,foo,bar
  ```
- Added `exclude` and `exclude-multi` actions for dynamically excluding items
  ```sh
  seq 100 | fzf --bind 'ctrl-x:exclude'

  # 'exclude-multi' will exclude the selected items or the current item
  seq 100 | fzf --multi --bind 'ctrl-x:exclude-multi'
  ```
- Preview window now prints wrap indicator when wrapping is enabled
  ```sh
  seq 100 | xargs | fzf --wrap --preview 'echo {}' --preview-window wrap
  ```
- Bug fixes and improvements

0.59.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.59.0/_

- Prioritizing file name matches (#4192)
    - Added a new tiebreak option `pathname` for prioritizing file name matches
    - `--scheme=path` now sets `--tiebreak=pathname,length`
    - fzf will automatically choose `path` scheme
        * when the input is a TTY device, where fzf would start its built-in walker or run `$FZF_DEFAULT_COMMAND` which is usually a command for listing files,
        * but not when `reload` or `transform` action is bound to `start` event, because in that case, fzf can't be sure of the input type.
- Added `--header-lines-border` to display header from `--header-lines` with a separate border
  ```sh
  # Use --header-lines-border to separate two headers
  ps -ef | fzf --style full --layout reverse --header-lines 1 \
               --bind 'ctrl-r:reload(ps -ef)' --header 'Press CTRL-R to reload' \
               --header-lines-border bottom --no-list-border
  ```
- `click-header` event now sets `$FZF_CLICK_HEADER_WORD` and `$FZF_CLICK_HEADER_NTH`. You can use them to implement a clickable header for changing the search scope using the new `transform-nth` action.
  ```sh
  # Click on the header line to limit search scope
  ps -ef | fzf --style full --layout reverse --header-lines 1 \
               --header-lines-border bottom --no-list-border \
               --color fg:dim,nth:regular \
               --bind 'click-header:transform-nth(
                         echo $FZF_CLICK_HEADER_NTH
                       )+transform-prompt(
                         echo "$FZF_CLICK_HEADER_WORD> "
                       )'
  ```
    - `$FZF_KEY` was updated to expose the type of the click. e.g. `click`, `ctrl-click`, etc. You can use it to implement a more sophisticated behavior.
    - `kill` completion for bash and zsh were updated to use this feature
- Added `--no-input` option to completely disable and hide the input section
  ```sh
  # Click header to trigger search
  fzf --header '[src] [test]' --no-input --layout reverse \
      --header-border bottom --input-border \
      --bind 'click-header:transform-search:echo ${FZF_CLICK_HEADER_WORD:1:-1}'

  # Vim-like mode switch
  fzf --layout reverse-list --no-input \
      --bind 'j:down,k:up,/:show-input+unbind(j,k,/)' \
      --bind 'enter,esc,ctrl-c:transform:
        if [[ $FZF_INPUT_STATE = enabled ]]; then
          echo "rebind(j,k,/)+hide-input"
        elif [[ $FZF_KEY = enter ]]; then
          echo accept
        else
          echo abort
        fi
      '
  ```
    - You can later show the input section using `show-input` or `toggle-input` action, and hide it again using `hide-input`, or `toggle-input`.
- Extended `{q}` placeholder to support ranges. e.g. `{q:1}`, `{q:2..}`, etc.
- Added `search(...)` and `transform-search(...)` action to trigger an fzf search with an arbitrary query string. This can be used to extend the search syntax of fzf. In the following example, fzf will use the first word of the query to trigger ripgrep search, and use the rest of the query to perform fzf search within the result.
  ```sh
  export TEMP=$(mktemp -u)
  trap 'rm -f "$TEMP"' EXIT

  TRANSFORMER='
    rg_pat={q:1}      # The first word is passed to ripgrep
    fzf_pat={q:2..}   # The rest are passed to fzf

    if ! [[ -r "$TEMP" ]] || [[ $rg_pat != $(cat "$TEMP") ]]; then
      echo "$rg_pat" > "$TEMP"
      printf "reload:sleep 0.1; rg --column --line-number --no-heading --color=always --smart-case %q || true" "$rg_pat"
    fi
    echo "+search:$fzf_pat"
  '
  fzf --ansi --disabled \
    --with-shell 'bash -c' \
    --bind "start,change:transform:$TRANSFORMER"
  ```
- You can now bind actions to multiple keys and events at once by writing a comma-separated list of keys and events before the colon
  ```sh
  # Load 'ps -ef' output on start and reload it on CTRL-R
  fzf --bind 'start,ctrl-r:reload:ps -ef'
  ```
- `--min-height` option now takes a number followed by `+`, which tells fzf to show at least that many items in the list section. The default value is now changed to `10+`.
  ```sh
  # You will only see the input section which takes 3 lines
  fzf --style=full --height 1% --min-height 3

  # You will see 3 items in the list section
  fzf --style full --height 1% --min-height 3+
  ```
    - Shell integration scripts were updated to use `--min-height 20+` by default
- `--header-lines` will be displayed at the top in `reverse-list` layout
- Added `bell` action to ring the terminal bell
  ```sh
  # Press CTRL-Y to copy the current line to the clipboard and ring the bell
  fzf --bind 'ctrl-y:execute-silent(echo -n {} | pbcopy)+bell'
  ```
- Added `toggle-bind` action
- Bug fixes and improvements
- Fixed fish script to support fish 3.1.2 or later (@bitraid)

0.58.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.58.0/_

This version introduces three new border types, `--list-border`, `--input-border`, and `--header-border`, offering much greater flexibility for customizing the user interface.

<img src="https://raw.githubusercontent.com/junegunn/i/master/fzf-4-borders.png" />

Also, fzf now offers "style presets" for quick customization, which can be activated using the `--style` option.

| Preset    | Screenshot                                                                             |
| :---      | :---                                                                                   |
| `default` | <img src="https://raw.githubusercontent.com/junegunn/i/master/fzf-style-default.png"/> |
| `full`    | <img src="https://raw.githubusercontent.com/junegunn/i/master/fzf-style-full.png"/>    |
| `minimal` | <img src="https://raw.githubusercontent.com/junegunn/i/master/fzf-style-minimal.png"/> |

- Style presets (#4160)
    - `--style=full[:BORDER_STYLE]`
    - `--style=default`
    - `--style=minimal`
- Border and label for the list section (#4148)
    - Options
        - `--list-border[=STYLE]`
        - `--list-label=LABEL`
        - `--list-label-pos=COL[:bottom]`
    - Colors
        - `list-fg`
        - `list-bg`
        - `list-border`
        - `list-label`
    - Actions
        - `change-list-label`
        - `transform-list-label`
- Border and label for the input section (prompt line and info line) (#4154)
    - Options
        - `--input-border[=STYLE]`
        - `--input-label=LABEL`
        - `--input-label-pos=COL[:bottom]`
    - Colors
        - `input-fg` (`query`)
        - `input-bg`
        - `input-border`
        - `input-label`
    - Actions
        - `change-input-label`
        - `transform-input-label`
- Border and label for the header section (#4159)
    - Options
        - `--header-border[=STYLE]`
        - `--header-label=LABEL`
        - `--header-label-pos=COL[:bottom]`
    - Colors
        - `header-fg` (`header`)
        - `header-bg`
        - `header-border`
        - `header-label`
    - Actions
        - `change-header-label`
        - `transform-header-label`
- Added `--preview-border[=STYLE]` as short for `--preview-window=border[-STYLE]`
- Added new preview border style `line` which draws a single separator line between the preview window and the rest of the interface
- fzf will now render a dashed line (`┈┈`) in each `--gap` for better visual separation.
  ```sh
  # All bash/zsh functions, highlighted
  declare -f |
    perl -0 -pe 's/^}\n/}\0/gm' |
    bat --plain --language bash --color always |
    fzf --read0 --ansi --layout reverse --multi --highlight-line --gap
  ```
    * You can customize the line using `--gap-line[=STR]`.
- You can specify `border-native` to `--tmux` so that native tmux border is used instead of `--border`. This can be useful if you start a different program from inside the popup.
  ```sh
  fzf --tmux border-native --bind 'enter:execute:less {}'
  ```
- Added `toggle-multi-line` action
- Added `toggle-hscroll` action
- Added `change-nth` action for dynamically changing the value of the `--nth` option
  ```sh
  # Start with --nth 1, then 2, then 3, then back to the default, 1
  echo 'foo foobar foobarbaz' | fzf --bind 'space:change-nth(2|3|)' --nth 1 -q foo
  ```
- `--nth` parts of each line can now be rendered in a different text style
  ```sh
  # nth in a different style
  ls -al | fzf --nth -1 --color nth:italic
  ls -al | fzf --nth -1 --color nth:reverse
  ls -al | fzf --nth -1 --color nth:reverse:bold

  # Dim the other parts
  ls -al | fzf --nth -1 --color nth:regular,fg:dim

  # With 'change-nth'. The current nth option is exported as $FZF_NTH.
  ps -ef | fzf --reverse --header-lines 1 --header-border bottom --input-border \
             --color nth:regular,fg:dim \
             --bind 'ctrl-n:change-nth(8..|1|2|3|4|5|6|7|)' \
             --bind 'result:transform-prompt:echo "${FZF_NTH}> "'
  ```
- A single-character delimiter is now treated as a plain string delimiter rather than a regular expression delimiter, even if it's a regular expression meta-character.
    - This means you can just write `--delimiter '|'` instead of escaping it as `--delimiter '\|'`
- Bug fixes
- Bug fixes and improvements in fish scripts (thanks to @bitraid)

0.57.0
------
- You can now resize the preview window by dragging the border
- Built-in walker improvements
    - `--walker-root` can take multiple directory arguments. e.g. `--walker-root include src lib`
    - `--walker-skip` can handle multi-component patterns. e.g. `--walker-skip target/build`
- Removed long processing delay when displaying images in the preview window
- `FZF_PREVIEW_*` environment variables are exported to all child processes (#4098)
- Bug fixes in fish scripts

0.56.3
------
- Bug fixes in zsh scripts
    - fix(zsh): handle backtick trigger edge case (#4090)
    - revert(zsh): remove 'fc -RI' call in the history widget (#4093)
    - Thanks to @LangLangBart for the contributions

0.56.2
------
- Bug fixes
    - Fixed abnormal scrolling behavior when `--wrap` is set (#4083)
    - [zsh] Fixed warning message when `ksh_arrays` is set (#4084)

0.56.1
------
- Bug fixes and improvements
    - Fixed a race condition which would cause fzf to present stale results after `reload` (#4070)
    - `page-up` and `page-down` actions now work correctly with multi-line items (#4069)
    - `{n}` is allowed in `SCROLL` expression in `--preview-window` (#4079)
    - [zsh] Fixed regression in history loading with shared option (#4071)
    - [zsh] Better command extraction in zsh completion (#4082)
- Thanks to @LangLangBart, @jaydee-coder, @alex-huff, and @vejkse for the contributions

0.56.0
------
- Added `--gap[=N]` option to display empty lines between items.
    - This can be useful to visually separate adjacent multi-line items.
      ```sh
      # All bash functions, highlighted
      declare -f | perl -0777 -pe 's/^}\n/}\0/gm' |
        bat --plain --language bash --color always |
        fzf --read0 --ansi --reverse --multi --highlight-line --gap
      ```
    - Or just to make the list easier to read. For single-line items, you probably want to set `--color gutter:-1` as well to hide the gutter.
      ```sh
      fzf --info inline-right --gap --color gutter:-1
      ```
- Added `noinfo` option to `--preview-window` to hide the scroll indicator in the preview window
- Bug fixes
    - Thanks to @LangLangBart, @akinomyoga, and @charlievieth for fixing the bugs

0.55.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.55.0/_

- Added `exact-boundary-match` type to the search syntax. When a search term is single-quoted, fzf will search for the exact occurrences of the string with both ends at word boundaries.
  ```sh
  fzf --query "'here'" << EOF
  come here
  not there
  EOF
  ```
- [bash] Fuzzy path completion is enabled for all commands
    - 1. If the default completion is not already set
    - 2. And if the current bash supports `complete -D` option
    - However, fuzzy completion for some commands can be "dynamically" disabled by the dynamic completion loader
    - See the comment in `__fzf_default_completion` function for more information
- Comments are now allowed in `$FZF_DEFAULT_OPTS` and `$FZF_DEFAULT_OPTS_FILE`
  ```sh
  export FZF_DEFAULT_OPTS='
    # Layout options
    --layout=reverse
    --info=inline-right   # Show info on the right side of the prompt line
    # ...
  '
  ```
- Hyperlinks (OSC 8) are now supported in the preview window and in the main window
  ```sh
  printf '<< \e]8;;http://github.com/junegunn/fzf\e\\Link to \e[32mfz\e[0mf\e]8;;\e\\ >>' | fzf --ansi

  fzf --preview "printf '<< \e]8;;http://github.com/junegunn/fzf\e\\Link to \e[32mfz\e[0mf\e]8;;\e\\ >>'"
  ```
- The default `--ellipsis` is now `··` instead of `..`.
- [vim] A spec can have `exit` callback that is called with the exit status of fzf
    - This can be used to clean up temporary resources or restore the original state when fzf is closed without a selection
- Fixed `--tmux bottom` when the status line is not at the bottom
- Fixed extra scroll offset in multi-line mode (`--read0` or `--wrap`)
- Added fallback `ps` command for `kill` completion on Cygwin

0.54.3
------
- Fixed incompatibility of adaptive height specification and 'start:reload'
  ```sh
  # A regression in 0.54.0 would cause this to fail
  fzf --height '~100%' --bind 'start:reload:seq 10'
  ```
- Environment variables are now available to `$FZF_DEFAULT_COMMAND`
  ```sh
  FZF_DEFAULT_COMMAND='echo $FZF_QUERY' fzf --query foo
  ```

0.54.2
------
- Fixed incorrect syntax highlighting of truncated multi-line entries
- Updated GoReleaser to 2.1.0 to simplify notarization of macOS binaries
    - macOS archives will be in `tar.gz` format instead of `zip` format since we no longer notarize the zip files but binaries
- (Windows) Reverted a mintty fix in 0.54.0
    - As a result, mouse may not work on mintty in fullscreen mode. However, fzf will correctly read non-ASCII input in fullscreen mode (`--no-height`).
    - fzf unfortunately cannot read non-ASCII input when not in fullscreen mode on Windows. So if you need to input non-ASCII characters, add `--no-height` to your `$FZF_DEFAULT_OPTS`.
    - Any help in fixing this issue will be appreciated (#3799, #3847).

0.54.1
------
- Updated [fastwalk](https://github.com/charlievieth/fastwalk) dependency for built-in directory walker
    - [fastwalk: add optional sorting and improve documentation](https://github.com/charlievieth/fastwalk/pull/27)
    - [fastwalk: only check if MSYSTEM is set during MSYS/MSYS2](https://github.com/charlievieth/fastwalk/pull/28)
    - Thanks to @charlievieth
- Reverted ALT-C binding of fish to use `cd` instead of `builtin cd`
    - `builtin cd` was introduced to work around a bug of `cd` coming from `zoxide init --cmd cd fish` where it cannot handle `--` argument.
    - However, the default `cd` of fish is actually a wrapper function for supporting `cd -`, so we want to use it instead.
    - See [#3928](https://github.com/junegunn/fzf/pull/3928) for more information and consider helping zoxide fix the bug.

0.54.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.54.0/_

- Implemented line wrap of long items
    - `--wrap` option enables line wrap
    - `--wrap-sign` customizes the sign for wrapped lines (default: `↳ `)
    - `toggle-wrap` action toggles line wrap
      ```sh
      history | fzf --tac --wrap --bind 'ctrl-/:toggle-wrap' --wrap-sign $'\t↳ '
      ```
    - fzf by default binds `CTRL-/` and `ALT-/` to `toggle-wrap`
- Updated shell integration scripts to leverage line wrap
    - CTRL-R binding includes `--wrap-sign $'\t↳ '` to indent wrapped lines
    - `kill **` completion uses `--wrap` to show the whole line by default
      instead of showing it in the preview window
- Added `--info-command` option for customizing the info line
  ```sh
  # Prepend the current cursor position in yellow
  fzf --info-command='echo -e "\x1b[33;1m$FZF_POS\x1b[m/$FZF_INFO 💛"'
  ```
    - `$FZF_INFO` is set to the original info text
    - ANSI color codes are supported
- Pointer and marker signs can be set to empty strings
  ```sh
  # Minimal style
  fzf --pointer '' --marker '' --prompt '' --info hidden
  ```
- Better cache management and improved rendering for `--tail`
- Improved `--sync` behavior
    - When `--sync` is provided, fzf will not render the interface until the initial filtering and the associated actions (bound to any of `start`, `load`, `result`, or `focus`) are complete.
      ```sh
      # fzf will not render intermediate states
      (sleep 1; seq 1000000; sleep 1) |
        fzf --sync --query 5 --listen --bind start:up,load:up,result:up,focus:change-header:Ready
      ```
- GET endpoint is now available from `execute` and `transform` actions (it used to timeout due to lock conflict)
  ```sh
  fzf --listen --sync --bind 'focus:transform-header:curl -s localhost:$FZF_PORT?limit=0 | jq .'
  ```
- Added `offset-middle` action to place the current item is in the middle of the screen
- fzf will not start the initial reader when `reload` or `reload-sync` is bound to `start` event. `fzf < /dev/null` or `: | fzf` are no longer required and extraneous `load` event will not fire due to the empty list.
  ```sh
  # Now this will work as expected. Previously, this would print an invalid header line.
  # `fzf < /dev/null` or `: | fzf` would fix the problem, but then an extraneous
  # `load` event would fire and the header would be prematurely updated.
  fzf --header 'Loading ...' --header-lines 1 \
      --bind 'start:reload:sleep 1; ps -ef' \
      --bind 'load:change-header:Loaded!'
  ```
- Fixed mouse support on Windows
- Fixed crash when using `--tiebreak=end` with very long items
- zsh 5.0 compatibility (thanks to @LangLangBart)
- Fixed `--walker-skip` to also skip symlinks to directories
- Fixed `result` event not fired when input stream is not complete
- New tags will have `v` prefix so that they are available on https://proxy.golang.org/

0.53.0
------
_Release highlights: https://junegunn.github.io/fzf/releases/0.53.0/_

- Multi-line display
    - See [Processing multi-line items](https://junegunn.github.io/fzf/tips/processing-multi-line-items/)
    - fzf can now display multi-line items
      ```sh
      # All bash functions, highlighted
      declare -f | perl -0777 -pe 's/^}\n/}\0/gm' |
        bat --plain --language bash --color always |
        fzf --read0 --ansi --reverse --multi --highlight-line

      # Ripgrep multi-line output
      rg --pretty bash | perl -0777 -pe 's/\n\n/\n\0/gm' |
        fzf --read0 --ansi --multi --highlight-line --reverse --tmux 70%
      ```
        - To disable multi-line display, use `--no-multi-line`
    - CTRL-R bindings of bash, zsh, and fish have been updated to leverage multi-line display
    - The default `--pointer` and `--marker` have been changed from `>` to Unicode bar characters as they look better with multi-line items
    - Added `--marker-multi-line` to customize the select marker for multi-line entries with the default set to `╻┃╹`
      ```
      ╻First line
      ┃...
      ╹Last line
      ```
- Native tmux integration
    - Added `--tmux` option to replace fzf-tmux script and simplify distribution
      ```sh
      # --tmux [center|top|bottom|left|right][,SIZE[%]][,SIZE[%]]
      # Center, 100% width and 70% height
      fzf --tmux 100%,70% --border horizontal --padding 1,2

      # Left, 30% width
      fzf --tmux left,30%

      # Bottom, 50% height
      fzf --tmux bottom,50%
      ```
        - To keep the implementation simple, it only uses popups. You need tmux 3.3 or later.
    - To use `--tmux` in Vim plugin:
      ```vim
      let g:fzf_layout = { 'tmux': '100%,70%' }
      ```
- Added support for endless input streams
    - See [Browsing log stream with fzf](https://junegunn.github.io/fzf/tips/browsing-log-streams/)
    - Added `--tail=NUM` option to limit the number of items to keep in memory. This is useful when you want to browse an endless stream of data (e.g. log stream) with fzf while limiting memory usage.
      ```sh
      # Interactive filtering of a log stream
      tail -f *.log | fzf --tail 100000 --tac --no-sort --exact
      ```
- Better Windows Support
    - fzf now works on Git bash (mintty) out of the box via winpty integration
    - Many fixes and improvements for Windows
- man page is now embedded in the binary; `fzf --man` to see it
- Changed the default `--scroll-off` to 3, as we think it's a better default
- Process started by `execute` action now directly writes to and reads from `/dev/tty`. Manual `/dev/tty` redirection for interactive programs is no longer required.
  ```sh
  # Vim will work fine without /dev/tty redirection
  ls | fzf --bind 'space:execute:vim {}' > selected
  ```
- Added `print(...)` action to queue an arbitrary string to be printed on exit. This was mainly added to work around the limitation of `--expect` where it's not compatible with `--bind` on the same key and it would ignore other actions bound to it.
  ```sh
  # This doesn't work as expected because --expect is not compatible with --bind
  fzf --multi --expect ctrl-y --bind 'ctrl-y:select-all'

  # This is something you can do instead
  fzf --multi --bind 'enter:print()+accept,ctrl-y:select-all+print(ctrl-y)+accept'
  ```
    - We also considered making them compatible, but realized that some users may have been relying on the current behavior.
- [`NO_COLOR`](https://no-color.org/) environment variable is now respected. If the variable is set, fzf defaults to `--no-color` unless otherwise specified.

0.52.1
------
- Fixed a critical bug in the Windows version
    - Windows users are strongly encouraged to upgrade to this version

0.52.0
------
- Added `--highlight-line` to highlight the whole current line (à la `set cursorline` of Vim)
- Added color names for selected lines: `selected-fg`, `selected-bg`, and `selected-hl`
  ```sh
  fzf --border --multi --info inline-right --layout reverse --marker ▏ --pointer ▌ --prompt '▌ '  \
      --highlight-line --color gutter:-1,selected-bg:238,selected-fg:146,current-fg:189
  ```
- Added `click-header` event that is triggered when the header section is clicked. When the event is triggered, `$FZF_CLICK_HEADER_COLUMN` and `$FZF_CLICK_HEADER_LINE` are set.
  ```sh
  fd --type f |
    fzf --header $'[Files] [Directories]' --header-first \
        --bind 'click-header:transform:
          (( FZF_CLICK_HEADER_COLUMN <= 7 )) && echo "reload(fd --type f)"
          (( FZF_CLICK_HEADER_COLUMN >= 9 )) && echo "reload(fd --type d)"
        '
  ```
- Add `$FZF_COMPLETION_{DIR,PATH}_OPTS` for separately customizing the behavior of fuzzy completion
  ```sh
  # Set --walker options without 'follow' not to follow symbolic links
  FZF_COMPLETION_PATH_OPTS="--walker=file,dir,hidden"
  FZF_COMPLETION_DIR_OPTS="--walker=dir,hidden"
  ```
- Fixed Windows argument escaping
- Bug fixes and improvements
- The code was heavily refactored to allow using fzf as a library in Go programs. The API is still experimental and subject to change.
    - https://gist.github.com/junegunn/193990b65be48a38aac6ac49d5669170

0.51.0
------
- Added a new environment variable `$FZF_POS` exported to the child processes. It's the vertical position of the cursor in the list starting from 1.
  ```sh
  # Toggle selection to the top or to the bottom
  seq 30 | fzf --multi --bind 'load:pos(10)' \
    --bind 'shift-up:transform:for _ in $(seq $FZF_POS $FZF_MATCH_COUNT); do echo -n +toggle+up; done' \
    --bind 'shift-down:transform:for _ in $(seq 1 $FZF_POS); do echo -n +toggle+down; done'
  ```
- Added `--with-shell` option to start child processes with a custom shell command and flags
  ```sh
  gem list | fzf --with-shell 'ruby -e' \
    --preview 'pp Gem::Specification.find_by_name({1})' \
    --bind 'ctrl-o:execute-silent:
        spec = Gem::Specification.find_by_name({1})
        [spec.homepage, *spec.metadata.filter { _1.end_with?("uri") }.values].uniq.each do
          system "open", _1
        end
    '
  ```
- Added `change-multi` action for dynamically changing `--multi` option
    - `change-multi` - enable multi-select mode with no limit
    - `change-multi(NUM)` - enable multi-select mode with a limit
    - `change-multi(0)` - disable multi-select mode
- Windows improvements
    - `become` action is now supported on Windows
        - Unlike in *nix, this does not use `execve(2)`. Instead it spawns a new process and waits for it to finish, so the exact behavior may differ.
    - Fixed argument escaping for Windows cmd.exe. No redundant escaping of backslashes.
- Bug fixes and improvements

0.50.0
------
- Search performance optimization. You can observe 50%+ improvement in some scenarios.
  ```
  $ rg --line-number --no-heading --smart-case . > $DATA

  $ wc < $DATA
   5520118 26862362 897487793

  $ hyperfine -w 1 -L bin fzf-0.49.0,fzf-7ce6452,fzf-a5447b8,fzf '{bin} --filter "///" < $DATA | head -30'
  Summary
    fzf --filter "///" < $DATA | head -30 ran
      1.16 ± 0.03 times faster than fzf-a5447b8 --filter "///" < $DATA | head -30
      1.23 ± 0.03 times faster than fzf-7ce6452 --filter "///" < $DATA | head -30
      1.52 ± 0.03 times faster than fzf-0.49.0 --filter "///" < $DATA | head -30
  ```
- Added `jump` and `jump-cancel` events that are triggered when leaving `jump` mode
  ```sh
  # Default behavior
  fzf --bind space:jump

  # Same as jump-accept action
  fzf --bind space:jump,jump:accept

  # Accept on jump, abort on cancel
  fzf --bind space:jump,jump:accept,jump-cancel:abort

  # Change header on jump-cancel
  fzf --bind 'space:change-header(Type jump label)+jump,jump-cancel:change-header:Jump cancelled'
  ```
- Added a new environment variable `$FZF_KEY` exported to the child processes. It's the name of the last key pressed.
  ```sh
  fzf --bind 'space:jump,jump:accept,jump-cancel:transform:[[ $FZF_KEY =~ ctrl-c ]] && echo abort'
  ```
- fzf can be built with profiling options. See [BUILD.md](BUILD.md) for more information.
- Bug fixes

0.49.0
------
- Ingestion performance improved by around 40% (more or less depending on options)
- `--info=hidden` and `--info=inline-right` will no longer hide the horizontal separator by default. This gives you more flexibility in customizing the layout.
    ```sh
    fzf --border --info=inline-right
    fzf --border --info=inline-right --separator ═
    fzf --border --info=inline-right --no-separator
    fzf --border --info=hidden
    fzf --border --info=hidden --separator ━
    fzf --border --info=hidden --no-separator
    ```
- Added two environment variables exported to the child processes
    - `FZF_PREVIEW_LABEL`
    - `FZF_BORDER_LABEL`
    ```sh
    # Use the current value of $FZF_PREVIEW_LABEL to determine which actions to perform
    git ls-files |
      fzf --header 'Press CTRL-P to change preview mode' \
          --bind='ctrl-p:transform:[[ $FZF_PREVIEW_LABEL =~ cat ]] \
          && echo "change-preview(git log --color=always \{})+change-preview-label([[ log ]])" \
          || echo "change-preview(bat --color=always \{})+change-preview-label([[ cat ]])"'
    ```
- Renamed `track` action to `track-current` to highlight the difference between the global tracking state set by `--track` and a one-off tracking action
    - `track` is still available as an alias
- Added `untrack-current` and `toggle-track-current` actions
    - `*-current` actions are no-op when the global tracking state is set
- Bug fixes and minor improvements

0.48.1
------
- CTRL-T and ALT-C bindings can be disabled by setting `FZF_CTRL_T_COMMAND` and `FZF_ALT_C_COMMAND` to empty strings respectively when sourcing the script
    ```sh
    # bash
    FZF_CTRL_T_COMMAND= FZF_ALT_C_COMMAND= eval "$(fzf --bash)"

    # zsh
    FZF_CTRL_T_COMMAND= FZF_ALT_C_COMMAND= eval "$(fzf --zsh)"

    # fish
    fzf --fish | FZF_CTRL_T_COMMAND= FZF_ALT_C_COMMAND= source
    ```
    - Setting the variables after sourcing the script will have no effect
- Bug fixes

0.48.0
------
- Shell integration scripts are now embedded in the fzf binary. This simplifies the distribution, and the users are less likely to have problems caused by using incompatible scripts and binaries.
    - bash
      ```sh
      # Set up fzf key bindings and fuzzy completion
      eval "$(fzf --bash)"
      ```
    - zsh
      ```sh
      # Set up fzf key bindings and fuzzy completion
      eval "$(fzf --zsh)"
      ```
    - fish
      ```fish
      # Set up fzf key bindings
      fzf --fish | source
      ```
- Added options for customizing the behavior of the built-in walker
    | Option               | Description                                       | Default              |
    | ---                  | ---                                               | ---                  |
    | `--walker=OPTS`      | Walker options (`[file][,dir][,follow][,hidden]`) | `file,follow,hidden` |
    | `--walker-root=DIR`  | Root directory from which to start walker         | `.`                  |
    | `--walker-skip=DIRS` | Comma-separated list of directory names to skip   | `.git,node_modules`  |
    - Examples
        ```sh
        # Built-in walker is only used by standalone fzf when $FZF_DEFAULT_COMMAND is not set
        unset FZF_DEFAULT_COMMAND

        fzf # default: --walker=file,follow,hidden --walker-root=. --walker-skip=.git,node_modules
        fzf --walker=file,dir,hidden,follow --walker-skip=.git,node_modules,target

        # Walker options in $FZF_DEFAULT_OPTS
        export FZF_DEFAULT_OPTS="--walker=file,dir,hidden,follow --walker-skip=.git,node_modules,target"
        fzf

        # Reading from STDIN; --walker is ignored
        seq 100 | fzf --walker=dir

        # Reading from $FZF_DEFAULT_COMMAND; --walker is ignored
        export FZF_DEFAULT_COMMAND='seq 100'
        fzf --walker=dir
        ```
- Shell integration scripts have been updated to use the built-in walker with these new options and they are now much faster out of the box.

0.47.0
------
- Replaced ["the default find command"][find] with a built-in directory walker to simplify the code and to achieve better performance and consistent behavior across platforms.
  This doesn't affect you if you have `$FZF_DEFAULT_COMMAND` set.
    - Breaking changes:
        - Unlike [the previous "find" command][find], the new traversal code will list hidden files, but hidden directories will still be ignored
        - No filtering of `devtmpfs` or `proc` types
        - Traversal is parallelized, so the order of the entries will be different each time
    - You may wonder why fzf implements directory walker anyway when it's a filter program following the [Unix philosophy][unix].
      But fzf has had [the walker code for years][walker] to tackle the performance problem on Windows. And I decided to use the same approach on different platforms as well for the benefits listed above.
    - Built-in walker is using the excellent [charlievieth/fastwalk][fastwalk] library, which easily outperforms its competitors and supports safely following symlinks.
- Added `$FZF_DEFAULT_OPTS_FILE` to allow managing default options in a file
    - See [#3618](https://github.com/junegunn/fzf/pull/3618)
    - Option precedence from lower to higher
        1. Options read from `$FZF_DEFAULT_OPTS_FILE`
        1. Options from `$FZF_DEFAULT_OPTS`
        1. Options from command-line arguments
- Bug fixes and improvements

[find]: https://github.com/junegunn/fzf/blob/0.46.1/src/constants.go#L60-L64
[walker]: https://github.com/junegunn/fzf/pull/1847
[fastwalk]: https://github.com/charlievieth/fastwalk
[unix]: https://en.wikipedia.org/wiki/Unix_philosophy

0.46.1
------
- Bug fixes and improvements
- Fixed Windows binaries
- Downgraded Go version to 1.20 to support older versions of Windows
    - https://tip.golang.org/doc/go1.21#windows
- Updated [rivo/uniseg](https://github.com/rivo/uniseg) dependency to v0.4.6

0.46.0
------
- Added two new events
    - `result` - triggered when the filtering for the current query is complete and the result list is ready
    - `resize` - triggered when the terminal size is changed
- fzf now exports the following environment variables to the child processes
  | Variable           | Description                                                 |
  | ---                | ---                                                         |
  | `FZF_LINES`        | Number of lines fzf takes up excluding padding and margin   |
  | `FZF_COLUMNS`      | Number of columns fzf takes up excluding padding and margin |
  | `FZF_TOTAL_COUNT`  | Total number of items                                       |
  | `FZF_MATCH_COUNT`  | Number of matched items                                     |
  | `FZF_SELECT_COUNT` | Number of selected items                                    |
  | `FZF_QUERY`        | Current query string                                        |
  | `FZF_PROMPT`       | Prompt string                                               |
  | `FZF_ACTION`       | The name of the last action performed                       |
  - This allows you to write sophisticated transformations like so
    ```sh
    # Script to dynamically resize the preview window
    transformer='
      # 1 line for info, another for prompt, and 2 more lines for preview window border
      lines=$(( FZF_LINES - FZF_MATCH_COUNT - 4 ))
      if [[ $FZF_MATCH_COUNT -eq 0 ]]; then
        echo "change-preview-window:hidden"
      elif [[ $lines -gt 3 ]]; then
        echo "change-preview-window:$lines"
      elif [[ $FZF_PREVIEW_LINES -ne 3 ]]; then
        echo "change-preview-window:3"
      fi
    '
    seq 10000 | fzf --preview 'seq {} 10000' --preview-window up \
                    --bind "result:transform:$transformer" \
                    --bind "resize:transform:$transformer"
    ```
  - And we're phasing out `{fzf:prompt}` and `{fzf:action}`
- Changed [mattn/go-runewidth](https://github.com/mattn/go-runewidth) dependency to [rivo/uniseg](https://github.com/rivo/uniseg) for accurate results
    - Set `--ambidouble` if your terminal displays ambiguous width characters (e.g. box-drawing characters for borders) as 2 columns
    - `RUNEWIDTH_EASTASIAN=1` is still respected for backward compatibility, but it's recommended that you use this new option instead
- Bug fixes

0.45.0
------
- Added `transform` action to conditionally perform a series of actions
  ```sh
  # Disallow selecting an empty line
  echo -e "1. Hello\n2. Goodbye\n\n3. Exit" |
    fzf --height '~100%' --reverse --header 'Select one' \
        --bind 'enter:transform:[[ -n {} ]] && echo accept || echo "change-header:Invalid selection"'

  # Move cursor past the empty line
  echo -e "1. Hello\n2. Goodbye\n\n3. Exit" |
    fzf --height '~100%' --reverse --header 'Select one' \
        --bind 'enter:transform:[[ -n {} ]] && echo accept || echo "change-header:Invalid selection"' \
        --bind 'focus:transform:[[ -n {} ]] && exit; [[ {fzf:action} =~ up$ ]] && echo up || echo down'

  # A single key binding to toggle between modes
  fd --type file |
    fzf --prompt 'Files> ' \
        --header 'CTRL-T: Switch between Files/Directories' \
        --bind 'ctrl-t:transform:[[ ! {fzf:prompt} =~ Files ]] &&
                  echo "change-prompt(Files> )+reload(fd --type file)" ||
                  echo "change-prompt(Directories> )+reload(fd --type directory)"'
  ```
- Added placeholder expressions
    - `{fzf:action}` - The name of the last action performed
    - `{fzf:prompt}` - Prompt string (including ANSI color codes)
    - `{fzf:query}` - Synonym for `{q}`
- Added support for negative height
  ```sh
  # Terminal height minus 1, so you can still see the command line
  fzf --height=-1
  ```
  - This handles a terminal resize better than `--height=$(($(tput lines) - 1))`
- Added `accept-or-print-query` action that acts like `accept` but prints the
  current query when there's no match for the query
  ```sh
  # You can make CTRL-R paste the current query when there's no match
  export FZF_CTRL_R_OPTS='--bind enter:accept-or-print-query'
  ```
  - Note that there are alternative ways to implement the same strategy
    ```sh
    # 'become' is apparently more versatile but it's not available on Windows.
    export FZF_CTRL_R_OPTS='--bind "enter:become:if [ -z {} ]; then echo {q}; else echo {}; fi"'

    # Using the new 'transform' action
    export FZF_CTRL_R_OPTS='--bind "enter:transform:[ -z {} ] && echo print-query || echo accept"'
    ```
- Added `show-header` and `hide-header` actions
- Bug fixes

0.44.1
------
- Fixed crash when preview window is hidden on `focus` event

0.44.0
------
- (Experimental) Sixel image support in preview window (not available on Windows)
    - [bin/fzf-preview.sh](bin/fzf-preview.sh) is added to demonstrate how to
      display an image using Kitty image protocol or Sixel. You can use it
      like so:
      ```sh
      fzf --preview='fzf-preview.sh {}'
      ```
- (Experimental) iTerm2 inline image protocol support in preview window (not available on Windows)
  ```sh
  # Using https://iterm2.com/utilities/imgcat
  fzf --preview 'imgcat -W $FZF_PREVIEW_COLUMNS -H $FZF_PREVIEW_LINES {}'
  ```
- HTTP server can be configured to accept remote connections
  ```sh
  # FZF_API_KEY is required for a non-localhost listen address
  export FZF_API_KEY="$(head -c 32 /dev/urandom | base64)"
  fzf --listen 0.0.0.0:6266
  ```
    - To allow remote process execution, use `--listen-unsafe` instead
      (`execute*`, `reload*`, `become`, `preview`, `change-preview`, `transform-*`)
      ```sh
      fzf --listen-unsafe 0.0.0.0:6266
      ```
- Bug fixes

0.43.0
------
- (Experimental) Added support for Kitty image protocol in the preview window
  (not available on Windows)
  ```sh
  fzf --preview='
    if file --mime-type {} | grep -qF image/; then
      # --transfer-mode=memory is the fastest option but if you want fzf to be able
      # to redraw the image on terminal resize or on 'change-preview-window',
      # you need to use --transfer-mode=stream.
      kitty icat --clear --transfer-mode=memory --unicode-placeholder --stdin=no --place=${FZF_PREVIEW_COLUMNS}x${FZF_PREVIEW_LINES}@0x0 {} | sed \$d
    else
      bat --color=always {}
    fi
  '
  ```
- (Experimental) `--listen` server can report program state in JSON format (`GET /`)
  ```sh
  # fzf server started in "headless" mode
  fzf --listen 6266 2> /dev/null

  # Get program state
  curl localhost:6266 | jq .

  # Increase the number of items returned (default: 100)
  curl localhost:6266?limit=1000 | jq .
  ```
- `--listen` server can be secured by setting `$FZF_API_KEY` environment
  variable.
  ```sh
  export FZF_API_KEY="$(head -c 32 /dev/urandom | base64)"

  # Server
  fzf --listen 6266

  # Client
  curl localhost:6266 -H "x-api-key: $FZF_API_KEY" -d 'change-query(yo)'
  ```
- Added `toggle-header` action
- Added mouse events for `--bind`
    - `scroll-up` (bound to `up`)
    - `scroll-down` (bound to `down`)
    - `shift-scroll-up` (bound to `toggle+up`)
    - `shift-scroll-down` (bound to `toggle+down`)
    - `shift-left-click` (bound to `toggle`)
    - `shift-right-click` (bound to `toggle`)
    - `preview-scroll-up` (bound to `preview-up`)
    - `preview-scroll-down` (bound to `preview-down`)
    ```sh
    # Twice faster scrolling both in the main window and the preview window
    fzf --bind 'scroll-up:up+up,scroll-down:down+down' \
        --bind 'preview-scroll-up:preview-up+preview-up' \
        --bind 'preview-scroll-down:preview-down+preview-down' \
        --preview 'cat {}'
    ```
- Added `offset-up` and `offset-down` actions
  ```sh
  # Scrolling will behave similarly to CTRL-E and CTRL-Y of vim
  fzf --bind scroll-up:offset-up,scroll-down:offset-down \
      --bind ctrl-y:offset-up,ctrl-e:offset-down \
      --scroll-off=5
  ```
- Shell extensions
    - Updated bash completion for fzf options
    - bash key bindings no longer requires perl; it will use awk or mawk
      instead if perl is not found
    - Basic context-aware completion for ssh command
    - Applied `--scheme=path` for better ordering of the result
- Bug fixes and improvements

0.42.0
------
- Added new info style: `--info=right`
- Added new info style: `--info=inline-right`
- Added new border style `thinblock` which uses symbols for legacy computing
  [one eighth block elements](https://en.wikipedia.org/wiki/Symbols_for_Legacy_Computing)
    - Similarly to `block`, this style is suitable when using a different
      background color because the window is completely contained within the border.
      ```sh
      BAT_THEME=GitHub fzf --info=right --border=thinblock --preview-window=border-thinblock \
          --margin=3 --scrollbar=▏▕ --preview='bat --color=always --style=numbers {}' \
          --color=light,query:238,fg:238,bg:251,bg+:249,gutter:251,border:248,preview-bg:253
      ```
    - This style may not render correctly depending on the font and the
      terminal emulator.

0.41.1
------
- Fixed a bug where preview window is not updated when `--disabled` is set and
  a reload is triggered by `change:reload` binding

0.41.0
------
- Added color name `preview-border` and `preview-scrollbar`
- Added new border style `block` which uses [block elements](https://en.wikipedia.org/wiki/Block_Elements)
- `--scrollbar` can take two characters, one for the main window, the other
  for the preview window
- Putting it altogether:
  ```sh
  fzf-tmux -p 80% --padding 1,2 --preview 'bat --style=plain --color=always {}' \
      --color 'bg:237,bg+:235,gutter:237,border:238,scrollbar:236' \
      --color 'preview-bg:235,preview-border:236,preview-scrollbar:234' \
      --preview-window 'border-block' --border block --scrollbar '▌▐'
  ```
- Bug fixes and improvements

0.40.0
------
- Added `zero` event that is triggered when there's no match
  ```sh
  # Reload the candidate list when there's no match
  echo $RANDOM | fzf --bind 'zero:reload(echo $RANDOM)+clear-query' --height 3
  ```
- New actions
    - Added `track` action which makes fzf track the current item when the
      search result is updated. If the user manually moves the cursor, or the
      item is not in the updated search result, tracking is automatically
      disabled. Tracking is useful when you want to see the surrounding items
      by deleting the query string.
      ```sh
      # Narrow down the list with a query, point to a command,
      # and hit CTRL-T to see its surrounding commands.
      export FZF_CTRL_R_OPTS="
        --preview 'echo {}' --preview-window up:3:hidden:wrap
        --bind 'ctrl-/:toggle-preview'
        --bind 'ctrl-t:track+clear-query'
        --bind 'ctrl-y:execute-silent(echo -n {2..} | pbcopy)+abort'
        --color header:italic
        --header 'Press CTRL-Y to copy command into clipboard'"
      ```
    - Added `change-header(...)`
    - Added `transform-header(...)`
    - Added `toggle-track` action
- Fixed `--track` behavior when used with `--tac`
    - However, using `--track` with `--tac` is not recommended. The resulting
      behavior can be very confusing.
- Bug fixes and improvements

0.39.0
------
- Added `one` event that is triggered when there's only one match
  ```sh
  # Automatically select the only match
  seq 10 | fzf --bind one:accept
  ```
- Added `--track` option that makes fzf track the current selection when the
  result list is updated. This can be useful when browsing logs using fzf with
  sorting disabled.
  ```sh
  git log --oneline --graph --color=always | nl |
      fzf --ansi --track --no-sort --layout=reverse-list
  ```
- If you use `--listen` option without a port number fzf will automatically
  allocate an available port and export it as `$FZF_PORT` environment
  variable.
  ```sh
  # Automatic port assignment
  fzf --listen --bind 'start:execute-silent:echo $FZF_PORT > /tmp/fzf-port'

  # Say hello
  curl "localhost:$(cat /tmp/fzf-port)" -d 'preview:echo Hello, fzf is listening on $FZF_PORT.'
  ```
- A carriage return and a line feed character will be rendered as dim ␍ and
  ␊ respectively.
  ```sh
  printf "foo\rbar\nbaz" | fzf --read0 --preview 'echo {}'
  ```
- fzf will stop rendering a non-displayable characters as a space. This will
  likely cause less glitches in the preview window.
  ```sh
  fzf --preview 'head -1000 /dev/random'
  ```
- Bug fixes and improvements

0.38.0
------
- New actions
    - `become(...)` - Replace the current fzf process with the specified
      command using `execve(2)` system call.
      See https://github.com/junegunn/fzf#turning-into-a-different-process for
      more information.
      ```sh
      # Open selected files in Vim
      fzf --multi --bind 'enter:become(vim {+})'

      # Open the file in Vim and go to the line
      git grep --line-number . |
          fzf --delimiter : --nth 3.. --bind 'enter:become(vim {1} +{2})'
      ```
        - This action is not supported on Windows
    - `show-preview`
    - `hide-preview`
- Bug fixes
    - `--preview-window 0,hidden` should not execute the preview command until
      `toggle-preview` action is triggered

0.37.0
------
- Added a way to customize the separator of inline info
  ```sh
  fzf --info 'inline: ╱ ' --prompt '╱ ' --color prompt:bright-yellow
  ```
- New event
    - `focus` - Triggered when the focus changes due to a vertical cursor
      movement or a search result update
      ```sh
      fzf --bind 'focus:transform-preview-label:echo [ {} ]' --preview 'cat {}'

      # Any action bound to the event runs synchronously and thus can make the interface sluggish
      # e.g. lolcat isn't one of the fastest programs, and every cursor movement in
      #      fzf will be noticeably affected by its execution time
      fzf --bind 'focus:transform-preview-label:echo [ {} ] | lolcat -f' --preview 'cat {}'

      # Beware not to introduce an infinite loop
      seq 10 | fzf --bind 'focus:up' --cycle
      ```
- New actions
    - `change-border-label`
    - `change-preview-label`
    - `transform-border-label`
    - `transform-preview-label`
- Bug fixes and improvements

0.36.0
------
- Added `--listen=HTTP_PORT` option to start HTTP server. It allows external
  processes to send actions to perform via POST method.
  ```sh
  # Start HTTP server on port 6266
  fzf --listen 6266

  # Send actions to the server
  curl -XPOST localhost:6266 -d 'reload(seq 100)+change-prompt(hundred> )'
  ```
- Added draggable scrollbar to the main search window and the preview window
  ```sh
  # Hide scrollbar
  fzf --no-scrollbar

  # Customize scrollbar
  fzf --scrollbar ┆ --color scrollbar:blue
  ```
- New event
    - Added `load` event that is triggered when the input stream is complete
      and the initial processing of the list is complete.
      ```sh
      # Change the prompt to "loaded" when the input stream is complete
      (seq 10; sleep 1; seq 11 20) | fzf --prompt 'Loading> ' --bind 'load:change-prompt:Loaded> '

      # You can use it instead of 'start' event without `--sync` if asynchronous
      # trigger is not an issue.
      (seq 10; sleep 1; seq 11 20) | fzf --bind 'load:last'
      ```
- New actions
    - Added `pos(...)` action to move the cursor to the numeric position
        - `first` and `last` are equivalent to `pos(1)` and `pos(-1)` respectively
      ```sh
      # Put the cursor on the 10th item
      seq 100 | fzf --sync --bind 'start:pos(10)'

      # Put the cursor on the 10th to last item
      seq 100 | fzf --sync --bind 'start:pos(-10)'
      ```
    - Added `reload-sync(...)` action which replaces the current list only after
      the reload process is complete. This is useful when the command takes
      a while to produce the initial output and you don't want fzf to run against
      an empty list while the command is running.
      ```sh
      # You can still filter and select entries from the initial list for 3 seconds
      seq 100 | fzf --bind 'load:reload-sync(sleep 3; seq 1000)+unbind(load)'
      ```
    - Added `next-selected` and `prev-selected` actions to move between selected
      items
      ```sh
      # `next-selected` will move the pointer to the next selected item below the current line
      # `prev-selected` will move the pointer to the previous selected item above the current line
      seq 10 | fzf --multi --bind ctrl-n:next-selected,ctrl-p:prev-selected

      # Both actions respect --layout option
      seq 10 | fzf --multi --bind ctrl-n:next-selected,ctrl-p:prev-selected --layout reverse
      ```
    - Added `change-query(...)` action that simply changes the query string to the
      given static string. This can be useful when used with `--listen`.
      ```sh
      curl localhost:6266 -d "change-query:$(date)"
      ```
    - Added `transform-prompt(...)` action for transforming the prompt string
      using an external command
      ```sh
      # Press space to change the prompt string using an external command
      # (only the first line of the output is taken)
      fzf --bind 'space:reload(ls),load:transform-prompt(printf "%s> " "$(date)")'
      ```
    - Added `transform-query(...)` action for transforming the query string using
      an external command
      ```sh
      # Press space to convert the query to uppercase letters
      fzf --bind 'space:transform-query(tr "[:lower:]" "[:upper:]" <<< {q})'

      # Bind it to 'change' event for automatic conversion
      fzf --bind 'change:transform-query(tr "[:lower:]" "[:upper:]" <<< {q})'

      # Can only type numbers
      fzf --bind 'change:transform-query(sed "s/[^0-9]//g" <<< {q})'
      ```
    - `put` action can optionally take an argument string
      ```sh
      # a will put 'alpha' on the prompt, ctrl-b will put 'bravo'
      fzf --bind 'a:put+put(lpha),ctrl-b:put(bravo)'
      ```
- Added color name `preview-label` for `--preview-label` (defaults to `label`
  for `--border-label`)
- Better support for (Windows) terminals where each box-drawing character
  takes 2 columns. Set `RUNEWIDTH_EASTASIAN` environment variable to `0` or `1`.
    - On Vim, the variable will be automatically set if `&ambiwidth` is `double`
- Behavior changes
    - fzf will always execute the preview command if the command template
      contains `{q}` even when it's empty. If you prefer the old behavior,
      you'll have to check if `{q}` is empty in your command.
      ```sh
      # This will show // even when the query is empty
      : | fzf --preview 'echo /{q}/'

      # But if you don't want it,
      : | fzf --preview '[ -n {q} ] || exit; echo /{q}/'
      ```
    - `double-click` will behave the same as `enter` unless otherwise specified,
      so you don't have to repeat the same action twice in `--bind` in most cases.
      ```sh
      # No need to bind 'double-click' to the same action
      fzf --bind 'enter:execute:less {}' # --bind 'double-click:execute:less {}'
      ```
    - If the color for `separator` is not specified, it will default to the
      color for `border`. Same holds true for `scrollbar`. This is to reduce
      the number of configuration items required to achieve a consistent color
      scheme.
    - If `follow` flag is specified in `--preview-window` option, fzf will
      automatically scroll to the bottom of the streaming preview output. But
      when the user manually scrolls the window, the following stops. With
      this version, fzf will resume following if the user scrolls the window
      to the bottom.
    - Default border style on Windows is changed to `sharp` because some
      Windows terminals are not capable of displaying `rounded` border
      characters correctly.
- Minor bug fixes and improvements

0.35.1
------
- Fixed a bug where fzf with `--tiebreak=chunk` crashes on inverse match query
- Fixed a bug where clicking above fzf would paste escape sequences

0.35.0
------
- Added `start` event that is triggered only once when fzf finder starts.
  Since fzf consumes the input stream asynchronously, the input list is not
  available unless you use `--sync`.
  ```sh
  seq 100 | fzf --multi --sync --bind 'start:last+select-all+preview(echo welcome)'
  ```
- Added `--border-label` and `--border-label-pos` for putting label on the border
  ```sh
  # ANSI color codes are supported
  # (with https://github.com/busyloop/lolcat)
  label=$(curl -s http://metaphorpsum.com/sentences/1 | lolcat -f)

  # Border label at the center
  fzf --height=10 --border --border-label="╢ $label ╟" --color=label:italic:black

  # Left-aligned (positive integer)
  fzf --height=10 --border --border-label="╢ $label ╟" --border-label-pos=3 --color=label:italic:black

  # Right-aligned (negative integer) on the bottom line (:bottom)
  fzf --height=10 --border --border-label="╢ $label ╟" --border-label-pos=-3:bottom --color=label:italic:black
  ```
- Also added `--preview-label` and `--preview-label-pos` for the border of the
  preview window
  ```sh
  fzf --preview 'cat {}' --border --preview-label=' Preview ' --preview-label-pos=2
  ```
- Info panel (match counter) will be followed by a horizontal separator by
  default
    - Use `--no-separator` or `--separator=''` to hide the separator
    - You can specify an arbitrary string that is repeated to form the
      horizontal separator. e.g. `--separator=╸`
    - The color of the separator can be customized via `--color=separator:...`
    - ANSI color codes are also supported
  ```sh
  fzf --separator=╸ --color=separator:green
  fzf --separator=$(lolcat -f -F 1.4 <<< ▁▁▂▃▄▅▆▆▅▄▃▂▁▁) --info=inline
  ```
- Added `--border=bold` and `--border=double` along with
  `--preview-window=border-bold` and `--preview-window=border-double`

0.34.0
------
- Added support for adaptive `--height`. If the `--height` value is prefixed
  with `~`, fzf will automatically determine the height in the range according
  to the input size.
  ```sh
  seq 1 | fzf --height ~70% --border --padding 1 --margin 1
  seq 10 | fzf --height ~70% --border --padding 1 --margin 1
  seq 100 | fzf --height ~70% --border --padding 1 --margin 1
  ```
    - There are a few limitations
        - Not compatible with percent top/bottom margin/padding
          ```sh
          # This is not allowed (top/bottom margin in percent value)
          fzf --height ~50% --border --margin 5%,10%

          # This is allowed (top/bottom margin in fixed value)
          fzf --height ~50% --border --margin 2,10%
          ```
        - fzf will not start until it can determine the right height for the input
          ```sh
          # fzf will open immediately
          (sleep 2; seq 10) | fzf --height 50%

          # fzf will open after 2 seconds
          (sleep 2; seq 10) | fzf --height ~50%
          (sleep 2; seq 1000) | fzf --height ~50%
          ```
- Fixed tcell renderer used to render full-screen fzf on Windows
- ~~`--no-clear` is deprecated. Use `reload` action instead.~~

0.33.0
------
- Added `--scheme=[default|path|history]` option to choose scoring scheme
    - (Experimental)
    - We updated the scoring algorithm in 0.32.0, however we have learned that
      this new scheme (`default`) is not always giving the optimal result
    - `path`: Additional bonus point is only given to the characters after
      path separator. You might want to choose this scheme if you have many
      files with spaces in their paths.
    - `history`: No additional bonus points are given so that we give more
      weight to the chronological ordering. This is equivalent to the scoring
      scheme before 0.32.0. This also sets `--tiebreak=index`.
- ANSI color sequences with colon delimiters are now supported.
  ```sh
  printf "\e[38;5;208mOption 1\e[m\nOption 2" | fzf --ansi
  printf "\e[38:5:208mOption 1\e[m\nOption 2" | fzf --ansi
  ```
- Support `border-{up,down}` as the synonyms for `border-{top,bottom}` in
  `--preview-window`
- Added support for ANSI `strikethrough`
  ```sh
  printf "\e[9mdeleted" | fzf --ansi
  fzf --color fg+:strikethrough
  ```

0.32.1
------
- Fixed incorrect ordering of `--tiebreak=chunk`
- fzf-tmux will show fzf border instead of tmux popup border (requires tmux 3.3)
  ```sh
  fzf-tmux -p70%
  fzf-tmux -p70% --color=border:bright-red
  fzf-tmux -p100%,60% --color=border:bright-yellow --border=horizontal --padding 1,5 --margin 1,0
  fzf-tmux -p70%,100% --color=border:bright-green --border=vertical

  # Key bindings (CTRL-T, CTRL-R, ALT-C) will use these options
  export FZF_TMUX_OPTS='-p100%,60% --color=border:green --border=horizontal --padding 1,5 --margin 1,0'
  ```

0.32.0
------
- Updated the scoring algorithm
    - Different bonus points to different categories of word boundaries
      (listed higher to lower bonus point)
        - Word after whitespace characters or beginning of the string
        - Word after common delimiter characters (`/,:;|`)
        - Word after other non-word characters
      ```sh
      # foo/bar.sh` is preferred over `foo-bar.sh` on `bar`
      fzf --query=bar --height=4 << EOF
      foo-bar.sh
      foo/bar.sh
      EOF
      ```
- Added a new tiebreak `chunk`
    - Favors the line with shorter matched chunk. A chunk is a set of
      consecutive non-whitespace characters.
    - Unlike the default `length`, this scheme works well with tabular input
      ```sh
      # length prefers item #1, because the whole line is shorter,
      # chunk prefers item #2, because the matched chunk ("foo") is shorter
      fzf --height=6 --header-lines=2 --tiebreak=chunk --reverse --query=fo << "EOF"
      N | Field1 | Field2 | Field3
      - | ------ | ------ | ------
      1 | hello  | foobar | baz
      2 | world  | foo    | bazbaz
      EOF
      ```
    - If the input does not contain any spaces, `chunk` is equivalent to
      `length`. But we're not going to set it as the default because it is
      computationally more expensive.
- Bug fixes and improvements

0.31.0
------
- Added support for an alternative preview window layout that is activated
  when the size of the preview window is smaller than a certain threshold.
  ```sh
  # If the width of the preview window is smaller than 50 columns,
  # it will be displayed above the search window.
  fzf --preview 'cat {}' --preview-window 'right,50%,border-left,<50(up,30%,border-bottom)'

  # Or you can just hide it like so
  fzf --preview 'cat {}' --preview-window '<50(hidden)'
  ```
- fzf now uses SGR mouse mode to properly support mouse on larger terminals
- You can now use characters that do not satisfy `unicode.IsGraphic` constraint
  for `--marker`, `--pointer`, and `--ellipsis`. Allows Nerd Fonts and stuff.
  Use at your own risk.
- Bug fixes and improvements
- Shell extension
    - `kill` completion now requires trigger sequence (`**`) for consistency

0.30.0
------
- Fixed cursor flickering over the screen by hiding it during rendering
- Added `--ellipsis` option. You can take advantage of it to make fzf
  effectively search non-visible parts of the item.
  ```sh
  # Search against hidden line numbers on the far right
  nl /usr/share/dict/words                  |
    awk '{printf "%s%1000s\n", $2, $1}'     |
    fzf --nth=-1 --no-hscroll --ellipsis='' |
    awk '{print $2}'
  ```
- Added `rebind` action for restoring bindings after `unbind`
- Bug fixes and improvements

0.29.0
------
- Added `change-preview(...)` action to change the `--preview` command
    - cf. `preview(...)` is a one-off action that doesn't change the default
      preview command
- Added `change-preview-window(...)` action
    - You can rotate through the different options separated by `|`
      ```sh
      fzf --preview 'cat {}' --preview-window right:40% \
          --bind 'ctrl-/:change-preview-window(right,70%|down,40%,border-top|hidden|)'
      ```
- Fixed rendering of the prompt line when overflow occurs with `--info=inline`

0.28.0
------
- Added `--header-first` option to print header before the prompt line
  ```sh
  fzf --header $'Welcome to fzf\n▔▔▔▔▔▔▔▔▔▔▔▔▔▔' --reverse --height 30% --border --header-first
  ```
- Added `--scroll-off=LINES` option (similar to `scrolloff` option of Vim)
    - You can set it to a very large number so that the cursor stays in the
      middle of the screen while scrolling
      ```sh
      fzf --scroll-off=5
      fzf --scroll-off=999
      ```
- Fixed bug where preview window is not updated on `reload` (#2644)
- fzf on Windows will also use `$SHELL` to execute external programs
    - See #2638 and #2647
    - Thanks to @rashil2000, @vovcacik, and @janlazo

0.27.3
------
- Preview window is `hidden` by default when there are `preview` bindings but
  `--preview` command is not given
- Fixed bug where `{n}` is not properly reset on `reload`
- Fixed bug where spinner is not displayed on `reload`
- Enhancements in tcell renderer for Windows (#2616)
- Vim plugin
    - `sinklist` is added as a synonym to `sink*` so that it's easier to add
      a function to a spec dictionary
      ```vim
      let spec = { 'source': 'ls', 'options': ['--multi', '--preview', 'cat {}'] }
      function spec.sinklist(matches)
        echom string(a:matches)
      endfunction

      call fzf#run(fzf#wrap(spec))
      ```
    - Vim 7 compatibility

0.27.2
------
- 16 base ANSI colors can be specified by their names
  ```sh
  fzf --color fg:3,fg+:11
  fzf --color fg:yellow,fg+:bright-yellow
  ```
- Fix bug where `--read0` not properly displaying long lines

0.27.1
------
- Added `unbind` action. In the following Ripgrep launcher example, you can
  use `unbind(reload)` to switch to fzf-only filtering mode.
    - See https://github.com/junegunn/fzf/blob/master/ADVANCED.md#switching-to-fzf-only-search-mode
- Vim plugin
    - Vim plugin will stop immediately even when the source command hasn't finished
      ```vim
      " fzf will read the stream file while allowing other processes to append to it
      call fzf#run({'source': 'cat /dev/null > /tmp/stream; tail -f /tmp/stream'})
      ```
    - It is now possible to open popup window relative to the current window
      ```vim
      let g:fzf_layout = { 'window': { 'width': 0.9, 'height': 0.6, 'relative': v:true, 'yoffset': 1.0 } }
      ```

0.27.0
------
- More border options for `--preview-window`
  ```sh
  fzf --preview 'cat {}' --preview-window border-left
  fzf --preview 'cat {}' --preview-window border-left --border horizontal
  fzf --preview 'cat {}' --preview-window top:border-bottom
  fzf --preview 'cat {}' --preview-window top:border-horizontal
  ```
- Automatically set `/dev/tty` as STDIN on execute action
  ```sh
  # Redirect /dev/tty to suppress "Vim: Warning: Input is not from a terminal"
  # ls | fzf --bind "enter:execute(vim {} < /dev/tty)"

  # "< /dev/tty" part is no longer needed
  ls | fzf --bind "enter:execute(vim {})"
  ```
- Bug fixes and improvements
- Signed and notarized macOS binaries
  (Huge thanks to [BACKERS.md](https://github.com/junegunn/junegunn/blob/main/BACKERS.md)!)

0.26.0
------
- Added support for fixed header in preview window
  ```sh
  # Display top 3 lines as the fixed header
  fzf --preview 'bat --style=header,grid --color=always {}' --preview-window '~3'
  ```
- More advanced preview offset expression to better support the fixed header
  ```sh
  # Preview with bat, matching line in the middle of the window below
  # the fixed header of the top 3 lines
  #
  #   ~3    Top 3 lines as the fixed header
  #   +{2}  Base scroll offset extracted from the second field
  #   +3    Extra offset to compensate for the 3-line header
  #   /2    Put in the middle of the preview area
  #
  git grep --line-number '' |
    fzf --delimiter : \
        --preview 'bat --style=full --color=always --highlight-line {2} {1}' \
        --preview-window '~3:+{2}+3/2'
  ```
- Added `select` and `deselect` action for unconditionally selecting or
  deselecting a single item in `--multi` mode. Complements `toggle` action.
- Significant performance improvement in ANSI code processing
- Bug fixes and improvements
- Built with Go 1.16

0.25.1
------
- Added `close` action
    - Close preview window if open, abort fzf otherwise
- Bug fixes and improvements

0.25.0
------
- Text attributes set in `--color` are not reset when fzf sees another
  `--color` option for the same element. This allows you to put custom text
  attributes in your `$FZF_DEFAULT_OPTS` and still have those attributes
  even when you override the colors.

  ```sh
  # Default colors and attributes
  fzf

  # Apply custom text attributes
  export FZF_DEFAULT_OPTS='--color fg+:italic,hl:-1:underline,hl+:-1:reverse:underline'

  fzf

  # Different colors but you still have the attributes
  fzf --color hl:176,hl+:177

  # Write "regular" if you want to clear the attributes
  fzf --color hl:176:regular,hl+:177:regular
  ```
- Renamed `--phony` to `--disabled`
- You can dynamically enable and disable the search functionality using the
  new `enable-search`, `disable-search`, and `toggle-search` actions
- You can assign a different color to the query string for when search is disabled
  ```sh
  fzf --color query:#ffffff,disabled:#999999 --bind space:toggle-search
  ```
- Added `last` action to move the cursor to the last match
    - The opposite action `top` is renamed to `first`, but `top` is still
      recognized as a synonym for backward compatibility
- Added `preview-top` and `preview-bottom` actions
- Extended support for alt key chords: alt with any case-sensitive single character
  ```sh
  fzf --bind alt-,:first,alt-.:last
  ```

0.24.4
------
- Added `--preview-window` option `follow`
  ```sh
  # Preview window will automatically scroll to the bottom
  fzf --preview-window follow --preview 'for i in $(seq 100000); do
    echo "$i"
    sleep 0.01
    (( i % 300 == 0 )) && printf "\033[2J"
  done'
  ```
- Added `change-prompt` action
  ```sh
  fzf --prompt 'foo> ' --bind $'a:change-prompt:\x1b[31mbar> '
  ```
- Bug fixes and improvements

0.24.3
------
- Added `--padding` option
  ```sh
  fzf --margin 5% --padding 5% --border --preview 'cat {}' \
      --color bg:#222222,preview-bg:#333333
  ```

0.24.2
------
- Bug fixes and improvements

0.24.1
------
- Fixed broken `--color=[bw|no]` option

0.24.0
------
- Real-time rendering of preview window
  ```sh
  # fzf can render preview window before the command completes
  fzf --preview 'sleep 1; for i in $(seq 100); do echo $i; sleep 0.01; done'

  # Preview window can process ANSI escape sequence (CSI 2 J) for clearing the display
  fzf --preview 'for i in $(seq 100000); do
    (( i % 200 == 0 )) && printf "\033[2J"
    echo "$i"
    sleep 0.01
  done'
  ```
- Updated `--color` option to support text styles
  - `regular` / `bold` / `dim` / `underline` / `italic` / `reverse` / `blink`
    ```sh
    # * Set -1 to keep the original color
    # * Multiple style attributes can be combined
    # * Italic style may not be supported by some terminals
    rg --line-number --no-heading --color=always "" |
      fzf --ansi --prompt "Rg: " \
          --color fg+:italic,hl:underline:-1,hl+:italic:underline:reverse:-1 \
          --color pointer:reverse,prompt:reverse,input:159 \
          --pointer '  '
    ```
- More `--border` options
  - `vertical`, `top`, `bottom`, `left`, `right`
  - Updated Vim plugin to use these new `--border` options
    ```vim
    " Floating popup window in the center of the screen
    let g:fzf_layout = { 'window': { 'width': 0.9, 'height': 0.6 } }

    " Popup with 100% width
    let g:fzf_layout = { 'window': { 'width': 1.0, 'height': 0.5, 'border': 'horizontal' } }

    " Popup with 100% height
    let g:fzf_layout = { 'window': { 'width': 0.5, 'height': 1.0, 'border': 'vertical' } }

    " Similar to 'down' layout, but it uses a popup window and doesn't affect the window layout
    let g:fzf_layout = { 'window': { 'width': 1.0, 'height': 0.5, 'yoffset': 1.0, 'border': 'top' } }

    " Opens on the right;
    "   'highlight' option is still supported but it will only take the foreground color of the group
    let g:fzf_layout = { 'window': { 'width': 0.5, 'height': 1.0, 'xoffset': 1.0, 'border': 'left', 'highlight': 'Comment' } }
    ```
- To indicate if `--multi` mode is enabled, fzf will print the number of
  selected items even when no item is selected
  ```sh
  seq 100 | fzf
    # 100/100
  seq 100 | fzf --multi
    # 100/100 (0)
  seq 100 | fzf --multi 5
    # 100/100 (0/5)
  ```
- Since 0.24.0, release binaries will be uploaded to https://github.com/junegunn/fzf/releases

0.23.1
------
- Added `--preview-window` options for disabling flags
    - `nocycle`
    - `nohidden`
    - `nowrap`
    - `default`
- Built with Go 1.14.9 due to performance regression
    - https://github.com/golang/go/issues/40727

0.23.0
------
- Support preview scroll offset relative to window height
  ```sh
  git grep --line-number '' |
    fzf --delimiter : \
        --preview 'bat --style=numbers --color=always --highlight-line {2} {1}' \
        --preview-window +{2}-/2
  ```
- Added `--preview-window` option for sharp edges (`--preview-window sharp`)
- Added `--preview-window` option for cyclic scrolling (`--preview-window cycle`)
- Reduced vertical padding around the preview window when `--preview-window
  noborder` is used
- Added actions for preview window
    - `preview-half-page-up`
    - `preview-half-page-down`
- Vim
    - Popup width and height can be given in absolute integer values
    - Added `fzf#exec()` function for getting the path of fzf executable
        - It also downloads the latest binary if it's not available by running
          `./install --bin`
- Built with Go 1.15.2
    - We no longer provide 32-bit binaries

0.22.0
------
- Added more options for `--bind`
    - `backward-eof` event
      ```sh
      # Aborts when you delete backward when the query prompt is already empty
      fzf --bind backward-eof:abort
      ```
    - `refresh-preview` action
      ```sh
      # Rerun preview command when you hit '?'
      fzf --preview 'echo $RANDOM' --bind '?:refresh-preview'
      ```
    - `preview` action
      ```sh
      # Default preview command with an extra preview binding
      fzf --preview 'file {}' --bind '?:preview:cat {}'

      # A preview binding with no default preview command
      # (Preview window is initially empty)
      fzf --bind '?:preview:cat {}'

      # Preview window hidden by default, it appears when you first hit '?'
      fzf --bind '?:preview:cat {}' --preview-window hidden
      ```
- Added preview window option for setting the initial scroll offset
  ```sh
  # Initial scroll offset is set to the line number of each line of
  # git grep output *minus* 5 lines
  git grep --line-number '' |
    fzf --delimiter : --preview 'nl {1}' --preview-window +{2}-5
  ```
- Added support for ANSI colors in `--prompt` string
- Smart match of accented characters
    - An unaccented character in the query string will match both accented and
      unaccented characters, while an accented character will only match
      accented characters. This is similar to how "smart-case" match works.
- Vim plugin
    - `tmux` layout option for using fzf-tmux
      ```vim
      let g:fzf_layout = { 'tmux': '-p90%,60%' }
      ```

0.21.1
------
- Shell extension
    - CTRL-R will remove duplicate commands
- fzf-tmux
    - Supports tmux popup window (require tmux 3.2 or above)
        - ```sh
          # 50% width and height
          fzf-tmux -p

          # 80% width and height
          fzf-tmux -p 80%

          # 80% width and 40% height
          fzf-tmux -p 80%,40%
          fzf-tmux -w 80% -h 40%

          # Window position
          fzf-tmux -w 80% -h 40% -x 0 -y 0
          fzf-tmux -w 80% -h 40% -y 1000

          # Write ordinary fzf options after --
          fzf-tmux -p -- --reverse --info=inline --margin 2,4 --border
          ```
        - On macOS, you can build the latest tmux from the source with
          `brew install tmux --HEAD`
- Bug fixes
    - Fixed Windows file traversal not to include directories
    - Fixed ANSI colors with `--keep-right`
    - Fixed _fzf_complete for zsh
- Built with Go 1.14.1

0.21.0
------
- `--height` option is now available on Windows as well (@kelleyma49)
- Added `--pointer` and `--marker` options
- Added `--keep-right` option that keeps the right end of the line visible
  when it's too long
- Style changes
    - `--border` will now print border with rounded corners around the
      finder instead of printing horizontal lines above and below it.
      The previous style is available via `--border=horizontal`
    - Unicode spinner
- More keys and actions for `--bind`
- Added PowerShell script for downloading Windows binary
- Vim plugin: Built-in floating windows support
  ```vim
  let g:fzf_layout = { 'window': { 'width': 0.9, 'height': 0.6 } }
  ```
- bash: Various improvements in key bindings (CTRL-T, CTRL-R, ALT-C)
    - CTRL-R will start with the current command-line as the initial query
    - CTRL-R properly supports multi-line commands
- Fuzzy completion API changed
  ```sh
  # Previous: fzf arguments given as a single string argument
  # - This style is still supported, but it's deprecated
  _fzf_complete "--multi --reverse --prompt=\"doge> \"" "$@" < <(
    echo foo
  )

  # New API: multiple fzf arguments before "--"
  # - Easier to write multiple options
  _fzf_complete --multi --reverse --prompt="doge> " -- "$@" < <(
    echo foo
  )
  ```
- Bug fixes and improvements

0.20.0
------
- Customizable preview window color (`preview-fg` and `preview-bg` for `--color`)
  ```sh
  fzf --preview 'cat {}' \
      --color 'fg:#bbccdd,fg+:#ddeeff,bg:#334455,preview-bg:#223344,border:#778899' \
      --border --height 20 --layout reverse --info inline
  ```
- Removed the immediate flicking of the screen on `reload` action.
  ```sh
  : | fzf --bind 'change:reload:seq {q}' --phony
  ```
- Added `clear-query` and `clear-selection` actions for `--bind`
- It is now possible to split a composite bind action over multiple `--bind`
  expressions by prefixing the later ones with `+`.
  ```sh
  fzf --bind 'ctrl-a:up+up'

  # Can be now written as
  fzf --bind 'ctrl-a:up' --bind 'ctrl-a:+up'

  # This is useful when you need to write special execute/reload form (i.e. `execute:...`)
  # to avoid parse errors and add more actions to the same key
  fzf --multi --bind 'ctrl-l:select-all+execute:less {+f}' --bind 'ctrl-l:+deselect-all'
  ```
- Fixed parse error of `--bind` expression where concatenated execute/reload
  action contains `+` character.
  ```sh
  fzf --multi --bind 'ctrl-l:select-all+execute(less {+f})+deselect-all'
  ```
- Fixed bugs of reload action
    - Not triggered when there's no match even when the command doesn't have
      any placeholder expressions
    - Screen not properly cleared when `--header-lines` not filled on reload

0.19.0
------

- Added `--phony` option which completely disables search functionality.
  Useful when you want to use fzf only as a selector interface. See below.
- Added "reload" action for dynamically updating the input list without
  restarting fzf. See https://github.com/junegunn/fzf/issues/1750 to learn
  more about it.
  ```sh
  # Using fzf as the selector interface for ripgrep
  RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
  INITIAL_QUERY="foo"
  FZF_DEFAULT_COMMAND="$RG_PREFIX '$INITIAL_QUERY' || true" \
    fzf --bind "change:reload:$RG_PREFIX {q} || true" \
        --ansi --phony --query "$INITIAL_QUERY"
  ```
- `--multi` now takes an optional integer argument which indicates the maximum
  number of items that can be selected
  ```sh
  seq 100 | fzf --multi 3 --reverse --height 50%
  ```
- If a placeholder expression for `--preview` and `execute` action (and the
  new `reload` action) contains `f` flag, it is replaced to the
  path of a temporary file that holds the evaluated list. This is useful
  when you multi-select a large number of items and the length of the
  evaluated string may exceed [`ARG_MAX`][argmax].
  ```sh
  # Press CTRL-A to select 100K items and see the sum of all the numbers
  seq 100000 | fzf --multi --bind ctrl-a:select-all \
                   --preview "awk '{sum+=\$1} END {print sum}' {+f}"
  ```
- `deselect-all` no longer deselects unmatched items. It is now consistent
  with `select-all` and `toggle-all` in that it only affects matched items.
- Due to the limitation of bash, fuzzy completion is enabled by default for
  a fixed set of commands. A helper function for easily setting up fuzzy
  completion for any command is now provided.
  ```sh
  # usage: _fzf_setup_completion path|dir COMMANDS...
  _fzf_setup_completion path git kubectl
  ```
- Info line style can be changed by `--info=STYLE`
    - `--info=default`
    - `--info=inline` (same as old `--inline-info`)
    - `--info=hidden`
- Preview window border can be disabled by adding `noborder` to
  `--preview-window`.
- When you transform the input with `--with-nth`, the trailing white spaces
  are removed.
- `ctrl-\`, `ctrl-]`, `ctrl-^`, and `ctrl-/` can now be used with `--bind`
- See https://github.com/junegunn/fzf/milestone/15?closed=1 for more details

[argmax]: https://unix.stackexchange.com/questions/120642/what-defines-the-maximum-size-for-a-command-single-argument

0.18.0
------

- Added placeholder expression for zero-based item index: `{n}` and `{+n}`
    - `fzf --preview 'echo {n}: {}'`
- Added color option for the gutter: `--color gutter:-1`
- Added `--no-unicode` option for drawing borders in non-Unicode, ASCII
  characters
- `FZF_PREVIEW_LINES` and `FZF_PREVIEW_COLUMNS` are exported to preview process
    - fzf still overrides `LINES` and `COLUMNS` as before, but they may be
      reset by the default shell.
- Bug fixes and improvements
    - See https://github.com/junegunn/fzf/milestone/14?closed=1
- Built with Go 1.12.1

0.17.5
------

- Bug fixes and improvements
    - See https://github.com/junegunn/fzf/milestone/13?closed=1
- Search query longer than the screen width is allowed (up to 300 chars)
- Built with Go 1.11.1

0.17.4
------

- Added `--layout` option with a new layout called `reverse-list`.
    - `--layout=reverse` is a synonym for `--reverse`
    - `--layout=default` is a synonym for `--no-reverse`
- Preview window will be updated even when there is no match for the query
  if any of the placeholder expressions (e.g. `{q}`, `{+}`) evaluates to
  a non-empty string.
- More keys for binding: `shift-{up,down}`, `alt-{up,down,left,right}`
- fzf can now start even when `/dev/tty` is not available by making an
  educated guess.
- Updated the default command for Windows.
- Fixes and improvements on bash/zsh completion
- install and uninstall scripts now supports generating files under
  `XDG_CONFIG_HOME` on `--xdg` flag.

See https://github.com/junegunn/fzf/milestone/12?closed=1 for the full list of
changes.

0.17.3
------
- `$LINES` and `$COLUMNS` are exported to preview command so that the command
  knows the exact size of the preview window.
- Better error messages when the default command or `$FZF_DEFAULT_COMMAND`
  fails.
- Reverted #1061 to avoid having duplicate entries in the list when find
  command detected a file system loop (#1120). The default command now
  requires that find supports `-fstype` option.
- fzf now distinguishes mouse left click and right click (#1130)
    - Right click is now bound to `toggle` action by default
    - `--bind` understands `left-click` and `right-click`
- Added `replace-query` action (#1137)
    - Replaces query string with the current selection
- Added `accept-non-empty` action (#1162)
    - Same as accept, except that it prevents fzf from exiting without any
      selection

0.17.1
------

- Fixed custom background color of preview window (#1046)
- Fixed background color issues of Windows binary
- Fixed Windows binary to execute command using cmd.exe with no parsing and
  escaping (#1072)
- Added support for `window` layout on Vim 8 using Vim 8 terminal (#1055)

0.17.0-2
--------

A maintenance release for auxiliary scripts. fzf binaries are not updated.

- Experimental support for the builtin terminal of Vim 8
    - fzf can now run inside GVim
- Updated Vim plugin to better handle `&shell` issue on fish
- Fixed a bug of fzf-tmux where invalid output is generated
- Fixed fzf-tmux to work even when `tput` does not work

0.17.0
------
- Performance optimization
- One can match literal spaces in extended-search mode with a space prepended
  by a backslash.
- `--expect` is now additive and can be specified multiple times.

0.16.11
-------
- Performance optimization
- Fixed missing preview update

0.16.10
-------
- Fixed invalid handling of ANSI colors in preview window
- Further improved `--ansi` performance

0.16.9
------
- Memory and performance optimization
    - Around 20% performance improvement for general use cases
    - Up to 5x faster processing of `--ansi`
    - Up to 50% reduction of memory usage
- Bug fixes and usability improvements
    - Fixed handling of bracketed paste mode
    - [ERROR] on info line when the default command failed
    - More efficient rendering of preview window
    - `--no-clear` updated for repetitive relaunching scenarios

0.16.8
------
- New `change` event and `top` action for `--bind`
    - `fzf --bind change:top`
        - Move cursor to the top result whenever the query string is changed
    - `fzf --bind 'ctrl-w:unix-word-rubout+top,ctrl-u:unix-line-discard+top'`
        - `top` combined with `unix-word-rubout` and `unix-line-discard`
- Fixed inconsistent tiebreak scores when `--nth` is used
- Proper display of tab characters in `--prompt`
- Fixed not to `--cycle` on page-up/page-down to prevent overshoot
- Git revision in `--version` output
- Basic support for Cygwin environment
- Many fixes in Vim plugin on Windows/Cygwin (thanks to @janlazo)

0.16.7
------
- Added support for `ctrl-alt-[a-z]` key chords
- CTRL-Z (SIGSTOP) now works with fzf
- fzf will export `$FZF_PREVIEW_WINDOW` so that the scripts can use it
- Bug fixes and improvements in Vim plugin and shell extensions

0.16.6
------
- Minor bug fixes and improvements
- Added `--no-clear` option for scripting purposes

0.16.5
------
- Minor bug fixes
- Added `toggle-preview-wrap` action
- Built with Go 1.8

0.16.4
------
- Added `--border` option to draw border above and below the finder
- Bug fixes and improvements

0.16.3
------
- Fixed a bug where fzf incorrectly display the lines when straddling tab
  characters are trimmed
- Placeholder expression used in `--preview` and `execute` action can
  optionally take `+` flag to be used with multiple selections
    - e.g. `git log --oneline | fzf --multi --preview 'git show {+1}'`
- Added `execute-silent` action for executing a command silently without
  switching to the alternate screen. This is useful when the process is
  short-lived and you're not interested in its output.
    - e.g. `fzf --bind 'ctrl-y:execute!(echo -n {} | pbcopy)'`
- `ctrl-space` is allowed in `--bind`

0.16.2
------
- Dropped ncurses dependency
- Binaries for freebsd, openbsd, arm5, arm6, arm7, and arm8
- Official 24-bit color support
- Added support for composite actions in `--bind`. Multiple actions can be
  chained using `+` separator.
    - e.g. `fzf --bind 'ctrl-y:execute(echo -n {} | pbcopy)+abort'`
- `--preview-window` with size 0 is allowed. This is used to make fzf execute
  preview command in the background without displaying the result.
- Minor bug fixes and improvements

0.16.1
------
- Fixed `--height` option to properly fill the window with the background
  color
- Added `half-page-up` and `half-page-down` actions
- Added `-L` flag to the default find command

0.16.0
------
- *Added `--height HEIGHT[%]` option*
    - fzf can now display finder without occupying the full screen
- Preview window will truncate long lines by default. Line wrap can be enabled
  by `:wrap` flag in `--preview-window`.
- Latin script letters will be normalized before matching so that it's easier
  to match against accented letters. e.g. `sodanco` can match `Só Danço Samba`.
    - Normalization can be disabled via `--literal`
- Added `--filepath-word` to make word-wise movements/actions (`alt-b`,
  `alt-f`, `alt-bs`, `alt-d`) respect path separators

0.15.9
------
- Fixed rendering glitches introduced in 0.15.8
- The default escape delay is reduced to 50ms and is configurable via
  `$ESCDELAY`
- Scroll indicator at the top-right corner of the preview window is always
  displayed when there's overflow
- Can now be built with ncurses 6 or tcell to support extra features
    - *ncurses 6*
        - Supports more than 256 color pairs
        - Supports italics
    - *tcell*
        - 24-bit color support
    - See https://github.com/junegunn/fzf/blob/master/BUILD.md

0.15.8
------
- Updated ANSI processor to handle more VT-100 escape sequences
- Added `--no-bold` (and `--bold`) option
- Improved escape sequence processing for WSL
- Added support for `alt-[0-9]`, `f11`, and `f12` for `--bind` and `--expect`

0.15.7
------
- Fixed panic when color is disabled and header lines contain ANSI colors

0.15.6
------
- Windows binaries! (@kelleyma49)
- Fixed the bug where header lines are cleared when preview window is toggled
- Fixed not to display ^N and ^O on screen
- Fixed cursor keys (or any key sequence that starts with ESC) on WSL by
  making fzf wait for additional keystrokes after ESC for up to 100ms

0.15.5
------
- Setting foreground color will no longer set background color to black
    - e.g. `fzf --color fg:153`
- `--tiebreak=end` will consider relative position instead of absolute distance
- Updated `fzf#wrap` function to respect `g:fzf_colors`

0.15.4
------
- Added support for range expression in preview and execute action
    - e.g. `ls -l | fzf --preview="echo user={3} when={-4..-2}; cat {-1}" --header-lines=1`
    - `{q}` will be replaced to the single-quoted string of the current query
- Fixed to properly handle unicode whitespace characters
- Display scroll indicator in preview window
- Inverse search term will use exact matcher by default
    - This is a breaking change, but I believe it makes much more sense. It is
      almost impossible to predict which entries will be filtered out due to
      a fuzzy inverse term. You can still perform inverse-fuzzy-match by
      prepending `!'` to the term.

0.15.3
------
- Added support for more ANSI attributes: dim, underline, blink, and reverse
- Fixed race condition in `toggle-preview`

0.15.2
------
- Preview window is now scrollable
    - With mouse scroll or with bindable actions
        - `preview-up`
        - `preview-down`
        - `preview-page-up`
        - `preview-page-down`
- Updated ANSI processor to support high intensity colors and ignore
  some VT100-related escape sequences

0.15.1
------
- Fixed panic when the pattern occurs after 2^15-th column
- Fixed rendering delay when displaying extremely long lines

0.15.0
------
- Improved fuzzy search algorithm
    - Added `--algo=[v1|v2]` option so one can still choose the old algorithm
      which values the search performance over the quality of the result
- Advanced scoring criteria
- `--read0` to read input delimited by ASCII NUL character
- `--print0` to print output delimited by ASCII NUL character

0.13.5
------
- Memory and performance optimization
    - Up to 2x performance with half the amount of memory

0.13.4
------
- Performance optimization
    - Memory footprint for ascii string is reduced by 60%
    - 15 to 20% improvement of query performance
    - Up to 45% better performance of `--nth` with non-regex delimiters
- Fixed invalid handling of `hidden` property of `--preview-window`

0.13.3
------
- Fixed duplicate rendering of the last line in preview window

0.13.2
------
- Fixed race condition where preview window is not properly cleared

0.13.1
------
- Fixed UI issue with large `--preview` output with many ANSI codes

0.13.0
------
- Added preview feature
    - `--preview CMD`
    - `--preview-window POS[:SIZE][:hidden]`
- `{}` in execute action is now replaced to the single-quoted (instead of
  double-quoted) string of the current line
- Fixed to ignore control characters for bracketed paste mode

0.12.2
------

- 256-color capability detection does not require `256` in `$TERM`
- Added `print-query` action
- More named keys for binding; <kbd>F1</kbd> ~ <kbd>F10</kbd>,
  <kbd>ALT-/</kbd>, <kbd>ALT-space</kbd>, and <kbd>ALT-enter</kbd>
- Added `jump` and `jump-accept` actions that implement [EasyMotion][em]-like
  movement
  ![][jump]

[em]: https://github.com/easymotion/vim-easymotion
[jump]: https://cloud.githubusercontent.com/assets/700826/15367574/b3999dc4-1d64-11e6-85da-28ceeb1a9bc2.png

0.12.1
------

- Ranking algorithm introduced in 0.12.0 is now universally applied
- Fixed invalid cache reference in exact mode
- Fixes and improvements in Vim plugin and shell extensions

0.12.0
------

- Enhanced ranking algorithm
- Minor bug fixes

0.11.4
------

- Added `--hscroll-off=COL` option (default: 10) (#513)
- Some fixes in Vim plugin and shell extensions

0.11.3
------

- Graceful exit on SIGTERM (#482)
- `$SHELL` instead of `sh` for `execute` action and `$FZF_DEFAULT_COMMAND` (#481)
- Changes in fuzzy completion API
    - [`_fzf_compgen_{path,dir}`](https://github.com/junegunn/fzf/commit/9617647)
    - [`_fzf_complete_COMMAND_post`](https://github.com/junegunn/fzf/commit/8206746)
      for post-processing

0.11.2
------

- `--tiebreak` now accepts comma-separated list of sort criteria
    - Each criterion should appear only once in the list
    - `index` is only allowed at the end of the list
    - `index` is implicitly appended to the list when not specified
    - Default is `length` (or equivalently `length,index`)
- `begin` criterion will ignore leading whitespaces when calculating the index
- Added `toggle-in` and `toggle-out` actions
    - Switch direction depending on `--reverse`-ness
    - `export FZF_DEFAULT_OPTS="--bind tab:toggle-out,shift-tab:toggle-in"`
- Reduced the initial delay when `--tac` is not given
    - fzf defers the initial rendering of the screen up to 100ms if the input
      stream is ongoing to prevent unnecessary redraw during the initial
      phase. However, 100ms delay is quite noticeable and might give the
      impression that fzf is not snappy enough. This commit reduces the
      maximum delay down to 20ms when `--tac` is not specified, in which case
      the input list quickly fills the entire screen.

0.11.1
------

- Added `--tabstop=SPACES` option

0.11.0
------

- Added OR operator for extended-search mode
- Added `--execute-multi` action
- Fixed incorrect cursor position when unicode wide characters are used in
  `--prompt`
- Fixes and improvements in shell extensions

0.10.9
------

- Extended-search mode is now enabled by default
    - `--extended-exact` is deprecated and instead we have `--exact` for
      orthogonally controlling "exactness" of search
- Fixed not to display non-printable characters
- Added `double-click` for `--bind` option
- More robust handling of SIGWINCH

0.10.8
------

- Fixed panic when trying to set colors after colors are disabled (#370)

0.10.7
------

- Fixed unserialized interrupt handling during execute action which often
  caused invalid memory access and crash
- Changed `--tiebreak=length` (default) to use trimmed length when `--nth` is
  used

0.10.6
------

- Replaced `--header-file` with `--header` option
- `--header` and `--header-lines` can be used together
- Changed exit status
    - 0: Okay
    - 1: No match
    - 2: Error
    - 130: Interrupted
- 64-bit linux binary is statically-linked with ncurses to avoid
  compatibility issues.

0.10.5
------

- `'`-prefix to unquote the term in `--extended-exact` mode
- Backward scan when `--tiebreak=end` is set

0.10.4
------

- Fixed to remove ANSI code from output when `--with-nth` is set

0.10.3
------

- Fixed slow performance of `--with-nth` when used with `--delimiter`
    - Regular expression engine of Golang as of now is very slow, so the fixed
      version will treat the given delimiter pattern as a plain string instead
      of a regular expression unless it contains special characters and is
      a valid regular expression.
    - Simpler regular expression for delimiter for better performance

0.10.2
------

### Fixes and improvements

- Improvement in perceived response time of queries
    - Eager, efficient rune array conversion
- Graceful exit when failed to initialize ncurses (invalid $TERM)
- Improved ranking algorithm when `--nth` option is set
- Changed the default command not to fail when there are files whose names
  start with dash

0.10.1
------

### New features

- Added `--margin` option
- Added options for sticky header
    - `--header-file`
    - `--header-lines`
- Added `cancel` action which clears the input or closes the finder when the
  input is already empty
    - e.g. `export FZF_DEFAULT_OPTS="--bind esc:cancel"`
- Added `delete-char/eof` action to differentiate `CTRL-D` and `DEL`

### Minor improvements/fixes

- Fixed to allow binding colon and comma keys
- Fixed ANSI processor to handle color regions spanning multiple lines

0.10.0
------

### New features

- More actions for `--bind`
    - `select-all`
    - `deselect-all`
    - `toggle-all`
    - `ignore`
- `execute(...)` action for running arbitrary command without leaving fzf
    - `fzf --bind "ctrl-m:execute(less {})"`
    - `fzf --bind "ctrl-t:execute(tmux new-window -d 'vim {}')"`
    - If the command contains parentheses, use any of the follows alternative
      notations to avoid parse errors
        - `execute[...]`
        - `execute~...~`
        - `execute!...!`
        - `execute@...@`
        - `execute#...#`
        - `execute$...$`
        - `execute%...%`
        - `execute^...^`
        - `execute&...&`
        - `execute*...*`
        - `execute;...;`
        - `execute/.../`
        - `execute|...|`
        - `execute:...`
            - This is the special form that frees you from parse errors as it
              does not expect the closing character
            - The catch is that it should be the last one in the
              comma-separated list
- Added support for optional search history
    - `--history HISTORY_FILE`
        - When used, `CTRL-N` and `CTRL-P` are automatically remapped to
          `next-history` and `previous-history`
    - `--history-size MAX_ENTRIES` (default: 1000)
- Cyclic scrolling can be enabled with `--cycle`
- Fixed the bug where the spinner was not spinning on idle input stream
    - e.g. `sleep 100 | fzf`

### Minor improvements/fixes

- Added synonyms for key names that can be specified for `--bind`,
  `--toggle-sort`, and `--expect`
- Fixed the color of multi-select marker on the current line
- Fixed to allow `^pattern$` in extended-search mode


0.9.13
------

### New features

- Color customization with the extended `--color` option

### Bug fixes

- Fixed premature termination of Reader in the presence of a long line which
  is longer than 64KB

0.9.12
------

### New features

- Added `--bind` option for custom key bindings

### Bug fixes

- Fixed to update "inline-info" immediately after terminal resize
- Fixed ANSI code offset calculation

0.9.11
------

### New features

- Added `--inline-info` option for saving screen estate (#202)
     - Useful inside Neovim
     - e.g. `let $FZF_DEFAULT_OPTS = $FZF_DEFAULT_OPTS.' --inline-info'`

### Bug fixes

- Invalid mutation of input on case conversion (#209)
- Smart-case for each term in extended-search mode (#208)
- Fixed double-click result when scroll offset is positive

0.9.10
------

### Improvements

- Performance optimization
- Less aggressive memoization to limit memory usage

### New features

- Added color scheme for light background: `--color=light`

0.9.9
-----

### New features

- Added `--tiebreak` option (#191)
- Added `--no-hscroll` option (#193)
- Visual indication of `--toggle-sort` (#194)

0.9.8
-----

### Bug fixes

- Fixed Unicode case handling (#186)
- Fixed to terminate on RuneError (#185)

0.9.7
-----

### New features

- Added `--toggle-sort` option (#173)
    - `--toggle-sort=ctrl-r` is applied to `CTRL-R` shell extension

### Bug fixes

- Fixed to print empty line if `--expect` is set and fzf is completed by
  `--select-1` or `--exit-0` (#172)
- Fixed to allow comma character as an argument to `--expect` option

0.9.6
-----

### New features

#### Added `--expect` option (#163)

If you provide a comma-separated list of keys with `--expect` option, fzf will
allow you to select the match and complete the finder when any of the keys is
pressed. Additionally, fzf will print the name of the key pressed as the first
line of the output so that your script can decide what to do next based on the
information.

```sh
fzf --expect=ctrl-v,ctrl-t,alt-s,f1,f2,~,@
```

The updated vim plugin uses this option to implement
[ctrlp](https://github.com/kien/ctrlp.vim)-compatible key bindings.

### Bug fixes

- Fixed to ignore ANSI escape code `\e[K` (#162)

0.9.5
-----

### New features

#### Added `--ansi` option (#150)

If you give `--ansi` option to fzf, fzf will interpret ANSI color codes from
the input, display the item with the ANSI colors (true colors are not
supported), and strips the codes from the output. This option is off by
default as it entails some overhead.

### Improvements

#### Reduced initial memory footprint (#151)

By removing unnecessary copy of pointers, fzf will use significantly smaller
amount of memory when it's started. The difference is hugely noticeable when
the input is extremely large. (e.g. `locate / | fzf`)

### Bug fixes

- Fixed panic on `--no-sort --filter ''` (#149)

0.9.4
-----

### New features

#### Added `--tac` option to reverse the order of the input.

One might argue that this option is unnecessary since we can already put `tac`
or `tail -r` in the command pipeline to achieve the same result. However, the
advantage of `--tac` is that it does not block until the input is complete.

### *Backward incompatible changes*

#### Changed behavior on `--no-sort`

`--no-sort` option will no longer reverse the display order within finder. You
may want to use the new `--tac` option with `--no-sort`.

```
history | fzf +s --tac
```

### Improvements

#### `--filter` will not block when sort is disabled

When fzf works in filtering mode (`--filter`) and sort is disabled
(`--no-sort`), there's no need to block until input is complete. The new
version of fzf will print the matches on-the-fly when the following condition
is met:

    --filter TERM --no-sort [--no-tac --no-sync]

or simply:

    -f TERM +s

This change removes unnecessary delay in the use cases like the following:

    fzf -f xxx +s | head -5

However, in this case, fzf processes the lines sequentially, so it cannot
utilize multiple cores, and fzf will run slightly slower than the previous
mode of execution where filtering is done in parallel after the entire input
is loaded. If the user is concerned about this performance problem, one can
add `--sync` option to re-enable buffering.

0.9.3
-----

### New features
- Added `--sync` option for multi-staged filtering

### Improvements
- `--select-1` and `--exit-0` will start finder immediately when the condition
  cannot be met
