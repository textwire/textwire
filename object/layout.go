package object

import "github.com/textwire/textwire/ast"

type Layout struct {
	Name    *String
	Program *ast.Program
}

func (l *Layout) Type() ObjectType {
	return LAYOUT_OBJ
}

func (l *Layout) String() string {
	return l.Program.String()
}
