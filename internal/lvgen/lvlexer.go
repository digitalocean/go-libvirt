// Copyright 2017 The go-libvirt Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lvgen

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
	"unicode/utf8"
)

// eof is returned by the lexer when there's no more input.
const eof = -1

// oneRuneTokens lists the runes the lexer will consider to be tokens when it
// finds them. These are returned to the parser using the integer value of their
// runes.
var oneRuneTokens = `{}[]<>(),=;:*`

var keywords = map[string]int{
	"hyper":    HYPER,
	"int":      INT,
	"short":    SHORT,
	"char":     CHAR,
	"bool":     BOOL,
	"case":     CASE,
	"const":    CONST,
	"default":  DEFAULT,
	"double":   DOUBLE,
	"enum":     ENUM,
	"float":    FLOAT,
	"opaque":   OPAQUE,
	"string":   STRING,
	"struct":   STRUCT,
	"switch":   SWITCH,
	"typedef":  TYPEDEF,
	"union":    UNION,
	"unsigned": UNSIGNED,
	"void":     VOID,
	"program":  PROGRAM,
	"version":  VERSION,
}

// item is a lexeme, or what the lexer returns to the parser.
type item struct {
	typ          int
	val          string
	line, column int
}

// String will display lexer items for humans to debug. There are some
// calculations here due to the way goyacc arranges token values; see the
// generated file y.go for an idea what's going on here, but the basic idea is
// that the lower token type values are reserved for single-rune tokens, which
// the lexer reports using the value of the rune itself. Everything else is
// allocated a range of type value up above all the possible single-rune values.
func (i item) String() string {
	tokType := i.typ
	if tokType >= yyPrivate {
		if tokType < yyPrivate+len(yyTok2) {
			tokType = int(yyTok2[tokType-yyPrivate])
		}
	}
	rv := fmt.Sprintf("%s %q %d:%d", yyTokname(tokType), i.val, i.line, i.column)
	return rv
}

// Lexer stores the state of this lexer.
type Lexer struct {
	input    string    // the string we're scanning.
	start    int       // start position of the item.
	pos      int       // current position in the input.
	line     int       // the current line (for error reporting).
	column   int       // current position within the current line.
	width    int       // width of the last rune scanned.
	items    chan item // channel of scanned lexer items (lexemes).
	lastItem item      // The last item the lexer handed the parser
}

// NewLexer will return a new lexer for the passed-in reader.
func NewLexer(rdr io.Reader) (*Lexer, error) {
	l := &Lexer{}

	b, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	l.input = string(b)
	l.items = make(chan item)

	return l, nil
}

// Run starts the lexer, and should be called in a goroutine.
func (l *Lexer) Run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// emit returns a token to the parser.
func (l *Lexer) emit(t int) {
	l.items <- item{t, l.input[l.start:l.pos], l.line, l.column}
	l.start = l.pos
}

// Lex gets the next token.
func (l *Lexer) Lex(st *yySymType) int {
	s := <-l.items
	l.lastItem = s
	st.val = s.val
	return int(s.typ)
}

// Error is called by the parser when it finds a problem.
func (l *Lexer) Error(s string) {
	fmt.Printf("parse error at %d:%d: %v\n", l.lastItem.line+1, l.lastItem.column+1, s)
	fmt.Printf("error at %q\n", l.lastItem.val)
}

// errorf is used by the lexer to report errors. It inserts an ERROR token into
// the items channel, and sets the state to nil, which stops the lexer's state
// machine.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{ERROR, fmt.Sprintf(format, args...), l.line, l.column}
	return nil
}

// next returns the rune at the current location, and advances to the next rune
// in the input.
func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	l.column++
	if r == '\n' {
		l.line++
		l.column = 0
	}
	return r
}

// ignore discards the current text from start to pos.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup moves back one character, but can only be called once per next() call.
func (l *Lexer) backup() {
	l.pos -= l.width
	if l.column > 0 {
		l.column--
	} else {
		l.line--
	}
	l.width = 0
}

// peek looks ahead at the next rune in the stream without consuming it.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept will advance to the next rune if it's contained in the string of valid
// runes passed in by the caller.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun advances over a number of valid runes, stopping as soon as it hits
// one not on the list.
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// procIdent checks whether an identifier matches the pattern for a procedure
// enum.
func procIdent(ident string) bool {
	// The pattern we're looking for is "<PROGRAM>_PROC_<NAME>", like
	// "REMOTE_PROC_DOMAIN_OPEN_CONSOLE"
	if ix := strings.Index(ident, "_PROC_"); ix != -1 {
		if strings.Index(ident[:ix], "_") == -1 {
			return true
		}
	}
	return false
}

// keyword checks whether the current lexeme is a keyword or not. If so it
// returns the keyword's token id, otherwise it returns IDENTIFIER.
func (l *Lexer) keyword() int {
	ident := l.input[l.start:l.pos]
	tok, ok := keywords[ident]
	if ok == true {
		return int(tok)
	}
	if procIdent(ident) {
		return PROCIDENTIFIER
	}
	return IDENTIFIER
}

// oneRuneToken determines whether a rune is a token. If so it returns the token
// id and true, otherwise it returns false.
func (l *Lexer) oneRuneToken(r rune) (int, bool) {
	if strings.IndexRune(oneRuneTokens, r) >= 0 {
		return int(r), true
	}

	return 0, false
}

// State functions
type stateFn func(*Lexer) stateFn

// lexText is the master lex routine. The lexer is started in this state.
func lexText(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], "/*") {
			return lexBlockComment
		}
		r := l.next()
		if r == eof {
			break
		}
		if unicode.IsSpace(r) {
			l.ignore()
			return lexText
		}
		if l.column == 1 && r == '%' {
			l.backup()
			return lexDirective
		}
		if unicode.IsLetter(r) {
			l.backup()
			return lexIdent
		}
		if unicode.IsNumber(r) || r == '-' {
			l.backup()
			return lexNumber
		}
		if t, isToken := l.oneRuneToken(r); isToken == true {
			l.emit(t)
		}
	}

	return nil
}

// lexBlockComment is used when we find a comment marker '/*' in the input.
func lexBlockComment(l *Lexer) stateFn {
	// Double star is used only at the start of metadata comments
	metadataComment := strings.HasPrefix(l.input[l.pos:], "/**")
	for {
		if strings.HasPrefix(l.input[l.pos:], "*/") {
			// Found the end. Advance past the '*/' and discard the comment body
			// unless it's a metadata comment
			l.next()
			l.next()
			if metadataComment {
				l.emit(METADATACOMMENT)
			} else {
				l.ignore()
			}
			return lexText
		}
		if l.next() == eof {
			return l.errorf("unterminated block comment")
		}
	}
}

// lexIdent handles identifiers.
func lexIdent(l *Lexer) stateFn {
	for {
		r := l.next()
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			continue
		}
		l.backup()
		break
	}
	// We may have a keyword, so check for that before emitting.
	l.emit(l.keyword())

	return lexText
}

// lexNumber handles decimal and hexadecimal numbers. Decimal numbers may begin
// with a '-'; hex numbers begin with '0x' and do not accept leading '-'.
func lexNumber(l *Lexer) stateFn {
	// Leading '-' is ok
	digits := "0123456789"
	neg := l.accept("-")
	if !neg {
		// allow '0x' for hex numbers, as long as there's not a leading '-'.
		r := l.peek()
		if r == '0' {
			l.next()
			if l.accept("x") {
				digits = "0123456789ABCDEFabcdef"
			}
		}
	}
	// followed by any number of digits
	l.acceptRun(digits)
	r := l.peek()
	if unicode.IsLetter(r) {
		l.next()
		return l.errorf("invalid number: %q", l.input[l.start:l.pos])
	}
	l.emit(CONSTANT)
	return lexText
}

// lexDirective handles lines beginning with '%'. These are used to emit C code
// directly to the output file. For now we're ignoring them, but some of the
// constants in the protocol file do depend on values from #included header
// files, so that may need to change.
func lexDirective(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == '\n' {
			l.ignore()
			return lexText
		}
		if r == eof {
			return l.errorf("unterminated directive")
		}
	}
}
