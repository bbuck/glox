package expr

// Sequenced represents an expression that should be executed first,
// the result discarded, and then a following expression should be
// executed. Like `do_the_first_thing(), do_the_second(), clean_up()`
type Sequenced struct {
	Left  Expr
	Right Expr
}

// NewSequence builds a new Sequenced expression from the left and
// right expressions provided.
func NewSequenced(left, right Expr) *Sequenced {
	return &Sequenced{
		Left:  left,
		Right: right,
	}
}

// Accept for Sequenced calls VisitSequenced on the visitor.
func (s *Sequenced) Accept(v Visitor) {
	v.VisitSequenced(s)
}
