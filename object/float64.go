package object

import "fmt"

type Float64 struct {
	Value float64
}

func (f *Float64) Type() ObjectType {
	return FLOAT64_OBJ
}

func (f *Float64) String() string {
	return fmt.Sprintf("%f", f.Value)
}
