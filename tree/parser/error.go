package parser

import (
	"errors"

	"github.com/bbuck/glox/errs"
	"github.com/bbuck/glox/token"
)

// ParseError represents a failure parsing the program.
var ParseError = errors.New("Parse Error")

func parseError(tok *token.T, msg string) error {
	errs.TokenError(tok, msg)
	return ParseError
}
