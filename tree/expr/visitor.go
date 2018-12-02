package expr

// Visitor defines an expression visitor interface, implement this and
// you can pass it into an expressions Accept method.
type Visitor interface {
	VisitBinary(*Binary)
	VisitLiteral(*Literal)
	VisitGrouping(*Grouping)
	VisitUnary(*Unary)
}
