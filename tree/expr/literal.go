package expr

// LiteralType represents why type of data the literal contains
// letting us know what to expect in the Value field.
type LiteralType uint8

// The various kinds of values that a Literal expression will contain
const (
	StringLiteral LiteralType = iota
	NumberLiteral
	KeywordLiteral
)

// Literal represents a value found literally in the code.
type Literal struct {
	Type  LiteralType
	Value interface{}
}

// NewLiteral constructs a new Literal expression with the given type
// and value.
func NewLiteral(typ LiteralType, value interface{}) *Literal {
	return &Literal{
		Type:  typ,
		Value: value,
	}
}

// Accept for a Literal calls the VisitLiteral method on the visitor.
func (l *Literal) Accept(v Visitor) {
	v.VisitLiteral(l)
}
