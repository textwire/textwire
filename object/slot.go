package object

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

func (s *Slot) Is(t ObjectType) bool {
	return t == s.Type()
}
