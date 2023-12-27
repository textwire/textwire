package object

import "fmt"

type Float32 struct {
	Value float32
}

func (f *Float32) Type() ObjectType {
	return FLOAT32_OBJ
}

func (f *Float32) String() string {
	return fmt.Sprintf("%f", f.Value)
}
