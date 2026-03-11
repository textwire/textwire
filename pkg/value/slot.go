package value

type Slot struct {
	Name    string
	Content Value
}

func (*Slot) Type() ValueType {
	return SLOT_VAL
}

func (s *Slot) String() string {
	if s.Content == nil {
		return ""
	}

	return s.Content.String()
}

func (s *Slot) Is(t ValueType) bool {
	return t == s.Type()
}
