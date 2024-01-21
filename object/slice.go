package object

import "fmt"

type Slice struct {
	Elements []Object
}

func (s *Slice) Type() ObjectType {
	return SLICE_OBJ
}

func (s *Slice) String() string {
	// todo: think about this implementation
	return fmt.Sprintf("%v", s.Elements)
}
