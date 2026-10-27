package value

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"time"

	"codeberg.org/textwire/textwire/v5/pkg/fail"
)

func NativeToValue(val any) Literal {
	if val == nil {
		return new(Nil)
	}

	lit, failure := reflectToLiteral(reflect.ValueOf(val), "")
	if failure != nil {
		return &Error{Err: failure}
	}

	return lit
}

// reflectToLiteral converts v into a Literal. The field is the struct field
// name the value came from, used for error messages. It's empty for values
// that aren't struct fields.
func reflectToLiteral(v reflect.Value, field string) (Literal, *fail.Error) {
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
		return uintToInt(v, field)
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

		return reflectToLiteral(v.Elem(), field)
	case reflect.Interface:
		if v.IsNil() {
			return new(Nil), nil
		}

		return reflectToLiteral(v.Elem(), field)
	}

	return nil, unsupportedTypeError(v, field)
}

func uintToInt(v reflect.Value, field string) (Literal, *fail.Error) {
	val := v.Uint()
	if val > math.MaxInt64 {
		return nil, uintOverflowError(val, field)
	}

	return &Int{Val: int64(val)}, nil
}

func unsupportedTypeError(v reflect.Value, field string) *fail.Error {
	if field == "" {
		return fail.New(nil, "", fail.OriginTpl, fail.ErrUnsupportedType, v.Type())
	}

	return fail.New(nil, "", fail.OriginTpl, fail.ErrUnsupportedFieldType, v.Type(), field)
}

func uintOverflowError(val uint64, field string) *fail.Error {
	if field == "" {
		return fail.New(nil, "", fail.OriginTpl, fail.ErrUintOverflow, val)
	}

	return fail.New(nil, "", fail.OriginTpl, fail.ErrUintOverflowField, val, field)
}

func nativeMapToValue(v reflect.Value) (Literal, *fail.Error) {
	obj := NewObj(nil)
	iter := v.MapRange()

	for iter.Next() {
		lit, failure := reflectToLiteral(iter.Value(), "")
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

		lit, failure := reflectToLiteral(v.Field(i), field.Name)
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
		lit, failure := reflectToLiteral(v.Index(i), "")
		if failure != nil {
			return nil, failure
		}

		arr.Elements[i] = lit
	}

	return arr, nil
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
