package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenError      TokenType = "ERROR"
	TokenEOF        TokenType = "EOF"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenString     TokenType = "STRING"
	TokenIndent     TokenType = "INDENT"
	TokenDedent     TokenType = "DEDENT"
	TokenNewline    TokenType = "NEWLINE"
	TokenComment    TokenType = "COMMENT"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

type Lexer struct {
	input       string
	pos         int
	line        int
	indentStack []int
	pendingToks []Token
	isLineStart bool
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:       input,
		line:        1,
		indentStack: []int{0},
		isLineStart: true,
	}
}

func (l *Lexer) NextToken() Token {
	if len(l.pendingToks) > 0 {
		tok := l.pendingToks[0]
		l.pendingToks = l.pendingToks[1:]
		return tok
	}

	if l.isLineStart {
		l.isLineStart = false
		indent := l.consumeIndentation()

		blankLine := false
		for i := l.pos; i < len(l.input); i++ {
			c := l.input[i]
			if c == '\n' || c == '#' || (c == '/' && i+1 < len(l.input) && l.input[i+1] == '/') {
				blankLine = true
				break
			}
			if c != ' ' && c != '\t' {
				break
			}
		}

		if blankLine {
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			if l.pos < len(l.input) && l.input[l.pos] == '\n' {
				l.line++
				l.pos++
			}
			l.isLineStart = true
			return Token{Type: TokenNewline, Line: l.line - 1}
		}

		currentIndent := l.indentStack[len(l.indentStack)-1]

		if indent > currentIndent {
			if (indent-currentIndent)%2 != 0 {
				return Token{Type: TokenError, Value: fmt.Sprintf("indentation must be in multiples of 2 spaces, got %d spaces at level %d", indent-currentIndent, currentIndent), Line: l.line}
			}
			l.indentStack = append(l.indentStack, indent)
			return Token{Type: TokenIndent, Line: l.line}
		}

		if indent < currentIndent {
			for len(l.indentStack) > 0 && l.indentStack[len(l.indentStack)-1] > indent {
				l.indentStack = l.indentStack[:len(l.indentStack)-1]
				l.pendingToks = append(l.pendingToks, Token{Type: TokenDedent, Line: l.line})
			}
			if l.indentStack[len(l.indentStack)-1] != indent {
				return Token{Type: TokenError, Value: "indentation compilation alignment tracking error", Line: l.line}
			}
			if len(l.pendingToks) > 0 {
				tok := l.pendingToks[0]
				l.pendingToks = l.pendingToks[1:]
				return tok
			}
		}
	}

	l.skipWhitespace()

	if l.pos >= len(l.input) {
		if len(l.indentStack) > 1 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return Token{Type: TokenDedent, Line: l.line}
		}
		return Token{Type: TokenEOF, Line: l.line}
	}

	ch := l.input[l.pos]

	if ch == '\n' {
		l.line++
		l.pos++
		l.isLineStart = true
		return Token{Type: TokenNewline, Line: l.line - 1}
	}

	if ch == '#' || (ch == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/') {
		return l.readComment()
	}

	if ch == '"' {
		return l.readString()
	}

	if isIdentifierStart(ch) {
		return l.readIdentifier()
	}

	l.pos++
	return Token{Type: TokenError, Value: fmt.Sprintf("unexpected character '%c'", ch), Line: l.line}
}

func (l *Lexer) consumeIndentation() int {
	count := 0
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' {
			count++
			l.pos++
		} else if ch == '\t' {
			count += 2
			l.pos++
		} else {
			break
		}
	}
	return count
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\r' || ch == '\t' {
			l.pos++
		} else {
			break
		}
	}
}

func (l *Lexer) readComment() Token {
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
	}
	return Token{Type: TokenComment, Value: l.input[start:l.pos], Line: l.line}
}

func (l *Lexer) readString() Token {
	l.pos++
	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			l.pos++
			break
		}
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.pos++
			switch l.input[l.pos] {
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case 'r':
				buf.WriteByte('\r')
			case '\\':
				buf.WriteByte('\\')
			case '"':
				buf.WriteByte('"')
			default:
				buf.WriteByte('\\')
				buf.WriteByte(l.input[l.pos])
			}
		} else {
			if ch == '\n' {
				l.line++
			}
			buf.WriteByte(ch)
		}
		l.pos++
	}
	val := buf.String()
	val = strings.ReplaceAll(val, "\n", "")
	val = strings.ReplaceAll(val, "\r", "")
	val = strings.ReplaceAll(val, "\t", "")
	return Token{Type: TokenString, Value: val, Line: l.line}
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if isIdentBase(ch) {
			l.pos++
		} else {
			break
		}
	}
	if l.pos < len(l.input) && l.input[l.pos] == '[' {
		l.pos++
		depth := 1
		for l.pos < len(l.input) && depth > 0 {
			ch := l.input[l.pos]
			if ch == '[' {
				depth++
			} else if ch == ']' {
				depth--
			}
			l.pos++
		}
	}
	return Token{Type: TokenIdentifier, Value: l.input[start:l.pos], Line: l.line}
}

func isIdentifierStart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_'
}

func isIdentBase(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' || ch == '-' || ch == '.'
}
