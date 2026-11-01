package value

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"codeberg.org/textwire/textwire/v5/pkg/fail"
)

// refKey identifies a slice, map or pointer that is currently being
// converted. It is used to detect cycles in native maps, slices and pointers.
type refKey struct {
	typ reflect.Type
	ptr uintptr
	len int
}

func NativeToValue(val any) Literal {
	if val == nil {
		return new(Nil)
	}

	lit, failure := reflectToLiteral(reflect.ValueOf(val), "", map[refKey]bool{})
	if failure != nil {
		return &Error{Err: failure}
	}

	return lit
}

func circularReferenceError(v reflect.Value) *fail.Error {
	return fail.New(nil, "", fail.OriginTpl, fail.ErrCircularValue, v.Type())
}

// reflectToLiteral converts v into a Literal. The field is the struct field
// name the value came from, used for error messages. It's empty for values
// that aren't struct fields.
func reflectToLiteral(
	v reflect.Value,
	field string,
	seen map[refKey]bool,
) (Literal, *fail.Error) {
	if !v.IsValid() {
		return new(Nil), nil
	}

	if key, ok := refKeyFor(v); ok {
		if seen[key] {
			return nil, circularReferenceError(v)
		}

		seen[key] = true
		defer delete(seen, key)
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

		return nativeStructToValue(v, seen)
	case reflect.Slice, reflect.Array:
		return nativeSliceToValue(v, seen)
	case reflect.Map:
		return nativeMapToValue(v, seen)
	case reflect.Pointer:
		if v.IsNil() {
			return new(Nil), nil
		}

		return reflectToLiteral(v.Elem(), field, seen)
	case reflect.Interface:
		if v.IsNil() {
			return new(Nil), nil
		}

		return reflectToLiteral(v.Elem(), field, seen)
	}

	return nil, unsupportedTypeError(v, field)
}

// refKeyFor returns the refKey for slices, maps and non-nil pointers. Nil
// pointers are skipped because they have no element to recurse into.
func refKeyFor(v reflect.Value) (refKey, bool) {
	switch v.Kind() {
	case reflect.Slice:
		return refKey{typ: v.Type(), ptr: v.Pointer(), len: v.Len()}, true
	case reflect.Map:
		return refKey{typ: v.Type(), ptr: v.Pointer()}, true
	case reflect.Pointer:
		return refKey{typ: v.Type(), ptr: v.Pointer()}, !v.IsNil()
	}

	return refKey{}, false
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

func nativeMapToValue(v reflect.Value, seen map[refKey]bool) (Literal, *fail.Error) {
	obj := NewObj(nil)

	iter := v.MapRange()
	for iter.Next() {
		key := mapKeyToString(iter.Key())

		value, err := reflectToLiteral(iter.Value(), "", seen)
		if err != nil {
			return nil, err
		}

		obj.Pairs[key] = value
	}

	return obj, nil
}

// mapKeyToString converts any map key into the string form used by Obj pairs.
func mapKeyToString(k reflect.Value) string {
	return fmt.Sprint(k.Interface())
}

func nativeStructToValue(v reflect.Value, seen map[refKey]bool) (Literal, *fail.Error) {
	obj := NewObj(nil)
	valType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := valType.Field(i)
		if !field.IsExported() {
			continue
		}

		lit, failure := reflectToLiteral(v.Field(i), field.Name, seen)
		if failure != nil {
			return nil, failure
		}

		obj.Pairs[field.Name] = lit
	}

	return obj, nil
}

func nativeSliceToValue(v reflect.Value, seen map[refKey]bool) (Literal, *fail.Error) {
	arr := new(Arr)
	arr.Elements = make([]Literal, v.Len())

	for i := 0; i < v.Len(); i++ {
		lit, failure := reflectToLiteral(v.Index(i), "", seen)
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
