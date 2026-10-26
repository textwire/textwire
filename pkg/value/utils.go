package value

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"codeberg.org/textwire/textwire/v5/pkg/fail"
)

func NativeToValue(val any) Literal {
	if val == nil {
		return new(Nil)
	}

	lit, failure := reflectToLiteral(reflect.ValueOf(val))
	if failure != nil {
		return &Error{Err: failure}
	}

	return lit
}

// reflectToLiteral converts v into a Literal. It returns an error for kinds
// Textwire can't represent, like func, chan and complex.
func reflectToLiteral(v reflect.Value) (Literal, *fail.Error) {
	if !v.IsValid() {
		return new(Nil), nil
	}

	switch v.Kind() {
	case reflect.String:
		return &Str{Val: v.String()}, nil
	case reflect.Bool:
		return &Bool{Val: v.Bool()}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Int{Val: v.Int()}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Int{Val: int64(v.Uint())}, nil
	case reflect.Float32, reflect.Float64:
		return &Float{Val: v.Float()}, nil
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return &Str{Val: t.Format("2006-01-02 15:04:05")}, nil
		}
		return nativeStructToValue(v)
	case reflect.Slice, reflect.Array:
		return nativeSliceToValue(v)
	case reflect.Map:
		return nativeMapToValue(v)
	case reflect.Pointer:
		if v.IsNil() {
			return new(Nil), nil
		}
		return reflectToLiteral(v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			return new(Nil), nil
		}
		return reflectToLiteral(v.Elem())
	}

	return nil, unsupportedTypeError(v, "")
}

func unsupportedTypeError(v reflect.Value, field string) *fail.Error {
	if field == "" {
		return fail.New(nil, "", fail.OriginTpl, fail.ErrUnsupportedType, v.Type())
	}
	return fail.New(nil, "", fail.OriginTpl, fail.ErrUnsupportedFieldType, v.Type(), field)
}

func nativeMapToValue(v reflect.Value) (Literal, *fail.Error) {
	obj := NewObj(nil)
	iter := v.MapRange()

	for iter.Next() {
		lit, failure := reflectToLiteral(iter.Value())
		if failure != nil {
			return nil, failure
		}

		obj.Pairs[iter.Key().String()] = lit
	}

	return obj, nil
}

func nativeStructToValue(v reflect.Value) (Literal, *fail.Error) {
	obj := NewObj(nil)
	valType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := valType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldVal := v.Field(i)

		if !isSupportedKind(fieldVal.Kind()) {
			return nil, unsupportedTypeError(fieldVal, field.Name)
		}

		lit, failure := reflectToLiteral(fieldVal)
		if failure != nil {
			return nil, failure
		}

		obj.Pairs[field.Name] = lit
	}

	return obj, nil
}

func nativeSliceToValue(v reflect.Value) (Literal, *fail.Error) {
	arr := new(Arr)
	arr.Elements = make([]Literal, v.Len())

	for i := 0; i < v.Len(); i++ {
		lit, failure := reflectToLiteral(v.Index(i))
		if failure != nil {
			return nil, failure
		}

		arr.Elements[i] = lit
	}

	return arr, nil
}

// isSupportedKind reports whether the kind can be converted without looking
// into the value itself.
func isSupportedKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Struct, reflect.Slice, reflect.Array, reflect.Map,
		reflect.Pointer, reflect.Interface:
		return true
	}

	return false
}

// marshalJSONString encodes s as a JSON string without HTML escaping. Textwire
// is a text engine and {!! !!} is the explicit no-escaping operator, so .json()
// doesn't escape HTML characters like <, > and &.
func marshalJSONString(s string) (string, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(s); err != nil {
		return "", err
	}

	// Encode writes a trailing newline after the value, we remove it
	return strings.TrimSuffix(buf.String(), "\n"), nil
}
