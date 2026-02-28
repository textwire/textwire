package object

import "fmt"

type Slot struct {
	Name    string
	Content Object
}

func (s *Slot) Type() ObjectType {
	return SLOT_OBJ
}

func (s *Slot) String() string {
	if s.Content == nil {
		return ""
	}

	return s.Content.String()
}

func (s *Slot) Dump(ident int) string {
	return fmt.Sprintf("@slot(%q)", s.Name)
}

func (s *Slot) JSON() (string, error) {
	return "", nil
}

func (s *Slot) Val() any {
	if s.Content == nil {
		return ""
	}
	return s.Content.Val()
}

func (s *Slot) Is(t ObjectType) bool {
	return t == s.Type()
}
