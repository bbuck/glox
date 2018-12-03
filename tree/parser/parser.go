package parser

import (
	"github.com/bbuck/glox/token"
	"github.com/bbuck/glox/tree/expr"
)

// P encapsulates the parsers current state allowing further calls to parse
// to maintain positonal information within the token list.
type P struct {
	tokens  []*token.T
	current int
	Err     error
}

// New constructs a new parser with the token list and returns it ready for
// use.
func New(toks []*token.T) *P {
	return &P{
		tokens: toks,
	}
}

// Parse returns the top-most expression in the syntax tree parsed from the
// token list. If a parse error occurred this will return nil instead.
func (p *P) Parse() expr.Expr {
	ex := p.expression()
	if p.Err == nil {
		return ex
	}

	// TODO: Do something because of error
	return nil
}

func (p *P) expression() expr.Expr {
	if p.Err != nil {
		return nil
	}

	return p.sequenced()
}

func (p *P) sequenced() expr.Expr {
	if p.Err != nil {
		return nil
	}

	ex := p.ternary()

	for p.match(token.Comma) {
		right := p.ternary()
		ex = expr.NewSequenced(ex, right)
	}

	return ex
}

func (p *P) ternary() expr.Expr {
	if p.Err != nil {
		return nil
	}

	ex := p.equality()

	if p.match(token.QuestionMark) {
		pos := p.expression()
		p.consume(token.Colon, "Expected ':' separating true/false branch")
		neg := p.expression()
		ex = expr.NewTernary(ex, pos, neg)
	}

	return ex
}

func (p *P) equality() expr.Expr {
	if p.Err != nil {
		return nil
	}

	ex := p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		op := p.previous()
		right := p.comparison()
		ex = expr.NewBinary(ex, op, right)
	}

	return ex
}

func (p *P) comparison() expr.Expr {
	if p.Err != nil {
		return nil
	}

	ex := p.addition()

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		op := p.previous()
		right := p.addition()
		ex = expr.NewBinary(ex, op, right)
	}

	return ex
}

func (p *P) addition() expr.Expr {
	if p.Err != nil {
		return nil
	}

	ex := p.multiplication()

	for p.match(token.Minus, token.Plus) {
		op := p.previous()
		right := p.multiplication()
		ex = expr.NewBinary(ex, op, right)
	}

	return ex
}

func (p *P) multiplication() expr.Expr {
	if p.Err != nil {
		return nil
	}

	ex := p.unary()

	for p.match(token.Slash, token.Star) {
		op := p.previous()
		right := p.unary()
		ex = expr.NewBinary(ex, op, right)
	}

	return ex
}

func (p *P) unary() expr.Expr {
	if p.Err != nil {
		return nil
	}

	if p.match(token.Bang, token.Bang) {
		op := p.previous()
		right := p.unary()

		return expr.NewUnary(op, right)
	}

	return p.primary()
}

func (p *P) primary() expr.Expr {
	if p.Err != nil {
		return nil
	}

	switch {
	case p.match(token.False):
		return expr.NewLiteral(expr.BooleanLiteral, false)
	case p.match(token.True):
		return expr.NewLiteral(expr.BooleanLiteral, true)
	case p.match(token.Nil):
		return expr.NewLiteral(expr.NilLiteral, nil)
	case p.match(token.Number):
		return expr.NewLiteral(expr.NumberLiteral, p.previous().Literal)
	case p.match(token.String):
		return expr.NewLiteral(expr.StringLiteral, p.previous().Literal)
	case p.match(token.LeftParen):
		ex := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression")

		return expr.NewGrouping(ex)
	}

	p.Err = parseError(p.peek(), "Expected expression")

	return nil
}

// helpers

func (p *P) match(types ...token.Type) bool {
	if p.Err != nil {
		return false
	}

	for _, typ := range types {
		if p.check(typ) {
			p.advance()

			return true
		}
	}

	return false
}

func (p *P) check(typ token.Type) bool {
	if p.Err != nil {
		return false
	}

	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == typ
}

func (p *P) advance() *token.T {
	if p.Err != nil {
		return nil
	}

	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *P) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *P) peek() *token.T {
	return p.tokens[p.current]
}

func (p *P) previous() *token.T {
	return p.tokens[p.current-1]
}

func (p *P) consume(typ token.Type, msg string) *token.T {
	if p.Err != nil {
		return nil
	}

	if p.check(typ) {
		return p.advance()
	}

	p.Err = parseError(p.peek(), msg)

	return nil
}

func (p *P) synchronize() {
	p.Err = nil

	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.Semicolon {
			return
		}

		switch p.peek().Type {
		case token.Class:
			fallthrough
		case token.Fun:
			fallthrough
		case token.Var:
			fallthrough
		case token.For:
			fallthrough
		case token.If:
			fallthrough
		case token.While:
			fallthrough
		case token.Print:
			fallthrough
		case token.Return:
			return
		}

		p.advance()
	}
}
