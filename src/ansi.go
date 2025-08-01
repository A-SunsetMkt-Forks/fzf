package fzf

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/junegunn/fzf/src/tui"
)

type ansiOffset struct {
	offset [2]int32
	color  ansiState
}

type url struct {
	uri    string
	params string
}

type ansiState struct {
	fg   tui.Color
	bg   tui.Color
	attr tui.Attr
	lbg  tui.Color
	url  *url
}

func (s *ansiState) colored() bool {
	return s.fg != -1 || s.bg != -1 || s.attr > 0 || s.lbg >= 0 || s.url != nil
}

func (s *ansiState) equals(t *ansiState) bool {
	if t == nil {
		return !s.colored()
	}
	return s.fg == t.fg && s.bg == t.bg && s.attr == t.attr && s.lbg == t.lbg && s.url == t.url
}

func (s *ansiState) ToString() string {
	if !s.colored() {
		return ""
	}

	ret := ""
	if s.attr&tui.Bold > 0 || s.attr&tui.BoldForce > 0 {
		ret += "1;"
	}
	if s.attr&tui.Dim > 0 {
		ret += "2;"
	}
	if s.attr&tui.Italic > 0 {
		ret += "3;"
	}
	if s.attr&tui.Underline > 0 {
		ret += "4;"
	}
	if s.attr&tui.Blink > 0 {
		ret += "5;"
	}
	if s.attr&tui.Reverse > 0 {
		ret += "7;"
	}
	if s.attr&tui.StrikeThrough > 0 {
		ret += "9;"
	}
	ret += toAnsiString(s.fg, 30) + toAnsiString(s.bg, 40)

	ret = "\x1b[" + strings.TrimSuffix(ret, ";") + "m"
	if s.url != nil {
		ret = fmt.Sprintf("\x1b]8;%s;%s\x1b\\%s\x1b]8;;\x1b", s.url.params, s.url.uri, ret)
	}
	return ret
}

func toAnsiString(color tui.Color, offset int) string {
	col := int(color)
	ret := ""
	if col == -1 {
		ret += strconv.Itoa(offset + 9)
	} else if col < 8 {
		ret += strconv.Itoa(offset + col)
	} else if col < 16 {
		ret += strconv.Itoa(offset - 30 + 90 + col - 8)
	} else if col < 256 {
		ret += strconv.Itoa(offset+8) + ";5;" + strconv.Itoa(col)
	} else if col >= (1 << 24) {
		r := strconv.Itoa((col >> 16) & 0xff)
		g := strconv.Itoa((col >> 8) & 0xff)
		b := strconv.Itoa(col & 0xff)
		ret += strconv.Itoa(offset+8) + ";2;" + r + ";" + g + ";" + b
	}
	return ret + ";"
}

func isPrint(c uint8) bool {
	return '\x20' <= c && c <= '\x7e'
}

func matchOperatingSystemCommand(s string, start int) int {
	// `\x1b][0-9][;:][[:print:]]+(?:\x1b\\\\|\x07)`
	//                 ^ match starting here after the first printable character
	//
	i := start // prefix matched in nextAnsiEscapeSequence()
	for ; i < len(s) && isPrint(s[i]); i++ {
	}
	if i < len(s) {
		if s[i] == '\x07' {
			return i + 1
		}
		// `\x1b]8;PARAMS;URI\x1b\\TITLE\x1b]8;;\x1b`
		//                   ------
		if s[i] == '\x1b' && i < len(s)-1 && s[i+1] == '\\' {
			return i + 2
		}
	}

	// `\x1b]8;PARAMS;URI\x1b\\TITLE\x1b]8;;\x1b`
	//                              ------------
	if i < len(s) && s[:i+1] == "\x1b]8;;\x1b" {
		return i + 1
	}

	return -1
}

func matchControlSequence(s string) int {
	// `\x1b[\\[()][0-9;:?]*[a-zA-Z@]`
	//                     ^ match starting here
	//
	i := 2 // prefix matched in nextAnsiEscapeSequence()
	for ; i < len(s); i++ {
		c := s[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ';', ':', '?':
			// ok
		default:
			if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '@' {
				return i + 1
			}
			return -1
		}
	}
	return -1
}

func isCtrlSeqStart(c uint8) bool {
	switch c {
	case '\\', '[', '(', ')':
		return true
	}
	return false
}

// nextAnsiEscapeSequence returns the ANSI escape sequence and is equivalent to
// calling FindStringIndex() on the below regex (which was originally used):
//
// "(?:\x1b[\\[()][0-9;:?]*[a-zA-Z@]|\x1b][0-9]+[;:][[:print:]]+(?:\x1b\\\\|\x07)|\x1b.|[\x0e\x0f]|.\x08|\n)"
func nextAnsiEscapeSequence(s string) (int, int) {
	// fast check for ANSI escape sequences
	i := 0
	for ; i < len(s); i++ {
		switch s[i] {
		case '\x0e', '\x0f', '\x1b', '\x08', '\n':
			// We ignore the fact that '\x08' cannot be the first char
			// in the string and be an escape sequence for the sake of
			// speed and simplicity.
			goto Loop
		}
	}
	return -1, -1

Loop:
	for ; i < len(s); i++ {
		switch s[i] {
		case '\n':
			// match: `\n`
			return i, i + 1
		case '\x08':
			// backtrack to match: `.\x08`
			if i > 0 && s[i-1] != '\n' {
				if s[i-1] < utf8.RuneSelf {
					return i - 1, i + 1
				}
				_, n := utf8.DecodeLastRuneInString(s[:i])
				return i - n, i + 1
			}
		case '\x1b':
			// match: `\x1b[\\[()][0-9;:?]*[a-zA-Z@]`
			if i+2 < len(s) && isCtrlSeqStart(s[i+1]) {
				if j := matchControlSequence(s[i:]); j != -1 {
					return i, i + j
				}
			}

			// match: `\x1b][0-9]+[;:][[:print:]]+(?:\x1b\\\\|\x07)`
			if i+5 < len(s) && s[i+1] == ']' {
				j := 2
				// \x1b][0-9]+[;:][[:print:]]+(?:\x1b\\\\|\x07)
				//      ------
				for ; i+j < len(s) && isNumeric(s[i+j]); j++ {
				}

				// \x1b][0-9]+[;:][[:print:]]+(?:\x1b\\\\|\x07)
				//            ---------------
				if j > 2 && i+j+1 < len(s) && (s[i+j] == ';' || s[i+j] == ':') && isPrint(s[i+j+1]) {
					if k := matchOperatingSystemCommand(s[i:], j+2); k != -1 {
						return i, i + k
					}
				}
			}

			// match: `\x1b.`
			if i+1 < len(s) && s[i+1] != '\n' {
				if s[i+1] < utf8.RuneSelf {
					return i, i + 2
				}
				_, n := utf8.DecodeRuneInString(s[i+1:])
				return i, i + n + 1
			}
		case '\x0e', '\x0f':
			// match: `[\x0e\x0f]`
			return i, i + 1
		}
	}
	return -1, -1
}

func extractColor(str string, state *ansiState, proc func(string, *ansiState) bool) (string, *[]ansiOffset, *ansiState) {
	// We append to a stack allocated variable that we'll
	// later copy and return, to save on allocations.
	offsets := make([]ansiOffset, 0, 32)

	if state != nil {
		offsets = append(offsets, ansiOffset{[2]int32{0, 0}, *state})
	}

	var (
		pstate    *ansiState // lazily allocated
		output    strings.Builder
		prevIdx   int
		runeCount int
	)
	for idx := 0; idx < len(str); {
		// Make sure that we found an ANSI code
		start, end := nextAnsiEscapeSequence(str[idx:])
		if start == -1 {
			break
		}
		start += idx
		idx += end

		// Check if we should continue
		prev := str[prevIdx:start]
		if proc != nil && !proc(prev, state) {
			return "", nil, nil
		}
		prevIdx = idx

		if len(prev) != 0 {
			runeCount += utf8.RuneCountInString(prev)
			// Grow the buffer size to the maximum possible length (string length
			// containing ansi codes) to avoid repetitive allocation
			if output.Cap() == 0 {
				output.Grow(len(str))
			}
			output.WriteString(prev)
		}

		code := str[start:idx]
		newState := interpretCode(code, state)
		if code == "\n" || !newState.equals(state) {
			if state != nil {
				// Update last offset
				(&offsets[len(offsets)-1]).offset[1] = int32(runeCount)
			}

			if code == "\n" {
				output.WriteRune('\n')
				runeCount++
				// Full-background marker
				if newState.lbg >= 0 {
					marker := newState
					marker.attr |= tui.FullBg
					offsets = append(offsets, ansiOffset{
						[2]int32{int32(runeCount), int32(runeCount)},
						marker,
					})
					// Reset the full-line background color
					newState.lbg = -1
				}
			}

			if newState.colored() {
				// Append new offset
				if pstate == nil {
					pstate = &ansiState{}
				}
				*pstate = newState
				state = pstate
				offsets = append(offsets, ansiOffset{
					[2]int32{int32(runeCount), int32(runeCount)},
					newState,
				})
			} else {
				// Discard state
				state = nil
			}
		}
	}

	var rest string
	var trimmed string
	if prevIdx == 0 {
		// No ANSI code found
		rest = str
		trimmed = str
	} else {
		rest = str[prevIdx:]
		output.WriteString(rest)
		trimmed = output.String()
	}
	if proc != nil {
		proc(rest, state)
	}
	if len(offsets) > 0 {
		if len(rest) > 0 && state != nil {
			// Update last offset
			runeCount += utf8.RuneCountInString(rest)
			(&offsets[len(offsets)-1]).offset[1] = int32(runeCount)
		}
		// Return a copy of the offsets slice
		a := make([]ansiOffset, len(offsets))
		copy(a, offsets)
		return trimmed, &a, state
	}
	return trimmed, nil, state
}

func parseAnsiCode(s string) (int, string) {
	var remaining string
	var i int
	// Faster than strings.IndexAny(";:")
	i = strings.IndexByte(s, ';')
	if i < 0 {
		i = strings.IndexByte(s, ':')
	}
	if i >= 0 {
		remaining = s[i+1:]
		s = s[:i]
	}

	if len(s) > 0 {
		// Inlined version of strconv.Atoi() that only handles positive
		// integers and does not allocate on error.
		code := 0
		for _, ch := range stringBytes(s) {
			ch -= '0'
			if ch > 9 {
				return -1, remaining
			}
			code = code*10 + int(ch)
		}
		return code, remaining
	}

	return -1, remaining
}

func interpretCode(ansiCode string, prevState *ansiState) ansiState {
	if ansiCode == "\n" {
		if prevState != nil {
			return *prevState
		}
		return ansiState{-1, -1, 0, -1, nil}
	}

	var state ansiState
	if prevState == nil {
		state = ansiState{-1, -1, 0, -1, nil}
	} else {
		state = ansiState{prevState.fg, prevState.bg, prevState.attr, prevState.lbg, prevState.url}
	}
	if ansiCode[0] != '\x1b' || ansiCode[1] != '[' || ansiCode[len(ansiCode)-1] != 'm' {
		if prevState != nil && (strings.HasSuffix(ansiCode, "0K") || strings.HasSuffix(ansiCode, "[K")) {
			state.lbg = prevState.bg
		} else if strings.HasPrefix(ansiCode, "\x1b]8;") && (strings.HasSuffix(ansiCode, "\x1b\\") || strings.HasSuffix(ansiCode, "\a")) {
			stLen := 2
			if strings.HasSuffix(ansiCode, "\a") {
				stLen = 1
			}
			// "\x1b]8;;\x1b\\" or "\x1b]8;;\a"
			if len(ansiCode) == 5+stLen && ansiCode[4] == ';' {
				state.url = nil
			} else if paramsEnd := strings.IndexRune(ansiCode[4:], ';'); paramsEnd >= 0 {
				params := ansiCode[4 : 4+paramsEnd]
				uri := ansiCode[5+paramsEnd : len(ansiCode)-stLen]
				state.url = &url{uri: uri, params: params}
			}
		}
		return state
	}

	reset := func() {
		state.fg = -1
		state.bg = -1
		state.attr = 0
	}

	if len(ansiCode) <= 3 {
		reset()
		return state
	}
	ansiCode = ansiCode[2 : len(ansiCode)-1]

	state256 := 0
	ptr := &state.fg

	count := 0
	for len(ansiCode) != 0 {
		var num int
		if num, ansiCode = parseAnsiCode(ansiCode); num != -1 {
			count++
			switch state256 {
			case 0:
				switch num {
				case 38:
					ptr = &state.fg
					state256++
				case 48:
					ptr = &state.bg
					state256++
				case 39:
					state.fg = -1
				case 49:
					state.bg = -1
				case 1:
					state.attr = state.attr | tui.Bold
				case 2:
					state.attr = state.attr | tui.Dim
				case 3:
					state.attr = state.attr | tui.Italic
				case 4:
					state.attr = state.attr | tui.Underline
				case 5:
					state.attr = state.attr | tui.Blink
				case 7:
					state.attr = state.attr | tui.Reverse
				case 9:
					state.attr = state.attr | tui.StrikeThrough
				case 22:
					state.attr = state.attr &^ tui.Bold
					state.attr = state.attr &^ tui.Dim
				case 23: // tput rmso
					state.attr = state.attr &^ tui.Italic
				case 24: // tput rmul
					state.attr = state.attr &^ tui.Underline
				case 25:
					state.attr = state.attr &^ tui.Blink
				case 27:
					state.attr = state.attr &^ tui.Reverse
				case 29:
					state.attr = state.attr &^ tui.StrikeThrough
				case 0:
					reset()
					state256 = 0
				default:
					if num >= 30 && num <= 37 {
						state.fg = tui.Color(num - 30)
					} else if num >= 40 && num <= 47 {
						state.bg = tui.Color(num - 40)
					} else if num >= 90 && num <= 97 {
						state.fg = tui.Color(num - 90 + 8)
					} else if num >= 100 && num <= 107 {
						state.bg = tui.Color(num - 100 + 8)
					}
				}
			case 1:
				switch num {
				case 2:
					state256 = 10 // MAGIC
				case 5:
					state256++
				default:
					state256 = 0
				}
			case 2:
				*ptr = tui.Color(num)
				state256 = 0
			case 10:
				*ptr = tui.Color(1<<24) | tui.Color(num<<16)
				state256++
			case 11:
				*ptr = *ptr | tui.Color(num<<8)
				state256++
			case 12:
				*ptr = *ptr | tui.Color(num)
				state256 = 0
			}
		}
	}

	// Empty sequence: reset
	if count == 0 {
		reset()
	}

	if state256 > 0 {
		*ptr = -1
	}
	return state
}
