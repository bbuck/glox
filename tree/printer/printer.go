package printer

import (
	"bytes"
	"fmt"

	"github.com/bbuck/glox/tree/expr"
)

type astPrinter struct {
	buf *bytes.Buffer
}

// Print will walk the expression tree and convert each expression
// into a viewable string value.
func Print(e expr.Expr) string {
	printer := &astPrinter{
		buf: new(bytes.Buffer),
	}
	e.Accept(printer)

	return printer.buf.String()
}

func (p *astPrinter) VisitBinary(b *expr.Binary) {
	p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (p *astPrinter) VisitUnary(u *expr.Unary) {
	p.parenthesize(u.Operator.Lexeme, u.Right)
}

func (p *astPrinter) VisitLiteral(l *expr.Literal) {
	if l.Value == nil {
		p.buf.WriteString("nil")
		return
	}

	p.buf.WriteString(fmt.Sprintf("%v", l.Value))
}

func (p *astPrinter) VisitGrouping(g *expr.Grouping) {
	p.parenthesize("group", g.Expression)
}

func (p *astPrinter) VisitSequenced(s *expr.Sequenced) {
	s.Left.Accept(p)
	p.buf.WriteString(" -> ")
	s.Right.Accept(p)
}

func (p *astPrinter) VisitTernary(t *expr.Ternary) {
	p.parenthesize("if", t.Condition, t.Positive, t.Negative)
}

func (p *astPrinter) parenthesize(name string, es ...expr.Expr) {
	p.buf.WriteRune('(')
	p.buf.WriteString(name)
	for _, e := range es {
		p.buf.WriteRune(' ')
		e.Accept(p)
	}
	p.buf.WriteRune(')')
}
