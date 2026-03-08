package value

import "fmt"

type Slot struct {
	Name    string
	Content Value
}

func (s *Slot) Type() ValueType {
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

func (s *Slot) Native() any {
	if s.Content == nil {
		return ""
	}
	return s.Content.Native()
}

func (s *Slot) Is(t ValueType) bool {
	return t == s.Type()
}
