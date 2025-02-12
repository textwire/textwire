package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type DotExp struct {
	Token token.Token // The dot token
	Left  Expression  // -->x.y
	Key   Expression  // x.y<--
	Pos   Position
}

func (de *DotExp) expressionNode() {
}

func (de *DotExp) TokenLiteral() string {
	return de.Token.Literal
}

func (de *DotExp) String() string {
	return fmt.Sprintf("(%s.%s)", de.Left, de.Key)
}

func (de *DotExp) Line() uint {
	return de.Token.Line
}

func (de *DotExp) Position() Position {
	return de.Pos
}
