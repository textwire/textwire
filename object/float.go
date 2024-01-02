package object

import "fmt"

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

func (f *Float) String() string {
	return fmt.Sprintf("%f", f.Value)
}
