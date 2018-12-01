package token

import "fmt"

// T is a token value type contain data about the token such as what
// type the token is, what string value does the token contain,
// what converted Go value does the token represent and what line
// was the token seen on.
type T struct {
	Type
	Lexeme  string
	Literal interface{}
	Line    uint
}

// New creates a new *T and returns it containing the information
// provided.
func New(t Type, lex string, lit interface{}, line uint) *T {
	return &T{
		Type:    t,
		Lexeme:  lex,
		Literal: lit,
		Line:    line,
	}
}

// String returns a frinedly printable value representing the T
func (t *T) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}
