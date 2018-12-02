package printer

// reverse polish notation printer

import (
	"bytes"
	"fmt"

	"github.com/bbuck/glox/tree/expr"
)

type rpnPrinter struct {
	buf *bytes.Buffer
}

// PrintRPN will walk the expression tree and print the operations in
// reverse polish notation (with the operation after the numbers)
func PrintRPN(e expr.Expr) string {
	printer := &rpnPrinter{
		buf: new(bytes.Buffer),
	}
	e.Accept(printer)

	return printer.buf.String()
}

func (p *rpnPrinter) VisitBinary(b *expr.Binary) {
	p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (p *rpnPrinter) VisitUnary(u *expr.Unary) {
	p.parenthesize(u.Operator.Lexeme, u.Right)
}

func (p *rpnPrinter) VisitLiteral(l *expr.Literal) {
	if l.Value == nil {
		p.buf.WriteString("nil")
		return
	}

	p.buf.WriteString(fmt.Sprintf("%v", l.Value))
}

func (p *rpnPrinter) VisitGrouping(g *expr.Grouping) {
	g.Expression.Accept(p)
}

func (p *rpnPrinter) VisitSequenced(s *expr.Sequenced) {
	s.Left.Accept(p)
	p.buf.WriteString(" -> ")
	s.Right.Accept(p)
}

func (p *rpnPrinter) parenthesize(name string, es ...expr.Expr) {
	for _, e := range es {
		e.Accept(p)
		p.buf.WriteRune(' ')
	}
	p.buf.WriteString(name)
}
