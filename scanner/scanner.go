package scanner

import (
	"bytes"

	"strconv"

	"unicode"

	"github.com/bbuck/glox/errs"
	"github.com/bbuck/glox/token"
)

// S is a token scanner for the Lox programming language. It will scan
// the source for tokens and build a list of tokens.
type S struct {
	Source    string
	runes     []rune
	tokens    []*token.T
	completed bool
	hadError  bool

	// locate ourselves
	start   int
	current int
	line    uint
}

// New constructs a new scanner for the provided source string
// with an empty list of tokens.
func New(source string) *S {
	return &S{
		Source: source,
		runes:  []rune(source),
		tokens: make([]*token.T, 0),
		line:   1,
	}
}

// ScanTokens will look through the source string and build a list of tokens
// that it finds. The return value is true if an error was encountered. This
// method is idempotent and safe to call multiple times. No changes will be
// made to the list of tokens.
func (s *S) ScanTokens() bool {
	if s.completed {
		return s.hadError
	}

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.addTokenRaw(token.EOF, "", nil)

	return s.hadError
}

// Tokens returns the list of scanned tokens after ScanTokens has been called.
// If you fetch the list before ScanTokens has been called then you will receive
// an empty list of tokens.
func (s *S) Tokens() []*token.T {
	return s.tokens
}

func (s *S) addTokenRaw(t token.Type, lex string, lit interface{}) {
	tok := token.New(t, lex, lit, s.line)
	s.tokens = append(s.tokens, tok)
}

func (s *S) addToken(t token.Type, lit interface{}) {

	s.addTokenRaw(t, s.currentLexeme(), lit)
}

func (s *S) addNoValueToken(t token.Type) {
	s.addToken(t, nil)
}

func (s *S) isAtEnd() bool {
	return s.current >= len(s.runes)
}

func (s *S) scanToken() {
	r := s.advance()
	switch r {
	case '(':
		s.addNoValueToken(token.LeftParen)
	case ')':
		s.addNoValueToken(token.RightParen)
	case '{':
		s.addNoValueToken(token.LeftBrace)
	case '}':
		s.addNoValueToken(token.RightBrace)
	case ',':
		s.addNoValueToken(token.Comma)
	case '.':
		s.addNoValueToken(token.Dot)
	case '-':
		s.addNoValueToken(token.Minus)
	case '+':
		s.addNoValueToken(token.Plus)
	case ';':
		s.addNoValueToken(token.Semicolon)
	case '*':
		s.addNoValueToken(token.Star)
	case '?':
		s.addNoValueToken(token.QuestionMark)
	case ':':
		s.addNoValueToken(token.Colon)
	case '!':
		s.scanEqualToken(token.BangEqual, token.Bang)
	case '=':
		s.scanEqualToken(token.EqualEqual, token.Equal)
	case '>':
		s.scanEqualToken(token.GreaterEqual, token.Greater)
	case '<':
		s.scanEqualToken(token.LessEqual, token.Less)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.scanBlockComment()
		} else {
			s.addNoValueToken(token.Slash)
		}
	case '"':
		s.scanString()
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
		// ignore these tokens
	case '\n':
		s.line++
	default:
		if isDigit(r) {
			s.scanNumber()
		} else if isAlpha(r) {
			s.scanIdentifier()
		} else {
			errs.Error(s.line, "Unexpected character.")
			s.hadError = true
		}
	}
}

func (s *S) scanEqualToken(found token.Type, notFound token.Type) {
	kind := notFound
	if s.match('=') {
		kind = found
	}
	s.addNoValueToken(kind)
}

func (s *S) advance() rune {
	s.current++

	return s.runes[s.current-1]
}

func (s *S) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}

	return s.runes[s.current]
}

func (s *S) peekNext() rune {
	if s.current+1 >= len(s.runes) {
		return rune(0)
	}

	return s.runes[s.current+1]
}

func (s *S) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.runes[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *S) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errs.Error(s.line, "Unterminated string.")
		s.hadError = true
		return
	}

	// capture the closing quote
	s.advance()

	value := s.currentLexeme()
	s.addToken(token.String, value[1:len(value)-1])
}

func (s *S) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value := s.currentLexeme()
	f64, _ := strconv.ParseFloat(value, 64)
	s.addToken(token.Number, f64)
}

func (s *S) scanIdentifier() {
	for isIdentifierRune(s.peek()) {
		s.advance()
	}

	lexeme := s.currentLexeme()
	if kind, ok := keywords[lexeme]; ok {
		s.addTokenRaw(kind, lexeme, nil)
		return
	}

	s.addTokenRaw(token.Identifier, lexeme, nil)
}

func (s *S) scanBlockComment() {
	level := 1
	for level > 0 {
		switch {
		case s.isAtEnd():
			errs.Error(s.line, "Unterminated block comment.")
			s.hadError = true
			return
		case s.peek() == '*' && s.peekNext() == '/':
			level--
			// consume *
			s.advance()
			// consume /
			s.advance()
		case s.peek() == '/' && s.peekNext() == '*':
			level++
			// consume /
			s.advance()
			// consume *
			s.advance()
		default:
			if s.peek() == '\n' {
				s.line++
			}

			s.advance()
		}
	}
}

func (s *S) currentLexeme() string {
	buf := new(bytes.Buffer)
	for i := s.start; i < s.current; i++ {
		buf.WriteRune(s.runes[i])
	}

	return buf.String()
}

// unattched helpers

func isAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsNumber(r)
}

func isIdentifierRune(r rune) bool {
	return isDigit(r) || isAlpha(r) || r == '_'
}
