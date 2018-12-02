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
	if ex, err := p.expression(); err == nil {
		return ex
	}

	// TODO: Do something because of error
	return nil
}

func (p *P) expression() (expr.Expr, error) {
	return p.sequenced()
}

func (p *P) sequenced() (expr.Expr, error) {
	ex, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.Comma) {
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		ex = expr.NewSequenced(ex, right)
	}

	return ex, nil
}

func (p *P) equality() (expr.Expr, error) {
	var (
		ex  expr.Expr
		err error
	)
	ex, err = p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BangEqual, token.EqualEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		ex = expr.NewBinary(ex, op, right)
	}

	return ex, nil
}

func (p *P) comparison() (expr.Expr, error) {
	var (
		ex  expr.Expr
		err error
	)
	ex, err = p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		op := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		ex = expr.NewBinary(ex, op, right)
	}

	return ex, nil
}

func (p *P) addition() (expr.Expr, error) {
	var (
		ex  expr.Expr
		err error
	)
	ex, err = p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(token.Minus, token.Plus) {
		op := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		ex = expr.NewBinary(ex, op, right)
	}

	return ex, nil
}

func (p *P) multiplication() (expr.Expr, error) {
	var (
		ex  expr.Expr
		err error
	)
	ex, err = p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.Slash, token.Star) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		ex = expr.NewBinary(ex, op, right)
	}

	return ex, nil
}

func (p *P) unary() (expr.Expr, error) {
	if p.match(token.Bang, token.Bang) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expr.NewUnary(op, right), nil
	}

	return p.primary()
}

func (p *P) primary() (expr.Expr, error) {
	switch {
	case p.match(token.False):
		return expr.NewLiteral(expr.BooleanLiteral, false), nil
	case p.match(token.True):
		return expr.NewLiteral(expr.BooleanLiteral, true), nil
	case p.match(token.Nil):
		return expr.NewLiteral(expr.NilLiteral, nil), nil
	case p.match(token.Number):
		return expr.NewLiteral(expr.NumberLiteral, p.previous().Literal), nil
	case p.match(token.String):
		return expr.NewLiteral(expr.StringLiteral, p.previous().Literal), nil
	case p.match(token.LeftParen):
		ex, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return expr.NewGrouping(ex), nil
	}

	return nil, parseError(p.peek(), "Expect expression")
}

// helpers

func (p *P) match(types ...token.Type) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()

			return true
		}
	}

	return false
}

func (p *P) check(typ token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == typ
}

func (p *P) advance() *token.T {
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

func (p *P) consume(typ token.Type, msg string) (*token.T, error) {
	if p.check(typ) {
		return p.advance(), nil
	}

	return nil, parseError(p.peek(), msg)
}

func (p *P) synchronize() {
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
