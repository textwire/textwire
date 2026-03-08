package value

import (
	"reflect"
)

func NativeToValue(val any) Value {
	switch v := val.(type) {
	case string:
		return &Str{Val: v}
	case bool:
		return &Bool{Val: v}
	case float32:
		return &Float{Val: float64(v)}
	case float64:
		return &Float{Val: v}
	case int64:
		return &Int{Val: v}
	case int:
		return &Int{Val: int64(v)}
	case int8:
		return &Int{Val: int64(v)}
	case int16:
		return &Int{Val: int64(v)}
	case int32:
		return &Int{Val: int64(v)}
	case uint:
		return &Int{Val: int64(v)}
	case uint8:
		return &Int{Val: int64(v)}
	case uint16:
		return &Int{Val: int64(v)}
	case uint32:
		return &Int{Val: int64(v)}
	case uint64:
		return &Int{Val: int64(v)}
	case nil:
		return new(Nil)
	}

	valType := reflect.TypeOf(val)

	switch valType.Kind() {
	case reflect.Struct:
		return nativeStructToValue(val)
	case reflect.Slice:
		return nativeSliceToArrayValue(convertToInterfaceSlice(val))
	case reflect.Map:
		return nativeMapToValue(val)
	case reflect.Pointer:
		v := reflect.ValueOf(val)
		if v.IsNil() {
			return new(Nil)
		}

		// NativeToValue is used recursively here
		return NativeToValue(v.Elem().Interface())
	}

	return nil
}

func nativeMapToValue(val any) Value {
	obj := NewObj(nil)

	valValue := reflect.ValueOf(val)
	for _, key := range valValue.MapKeys() {
		obj.Pairs[key.String()] = NativeToValue(valValue.MapIndex(key).Interface())
	}

	return obj
}

func convertToInterfaceSlice(slice any) []any {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	vals := make([]any, s.Len())
	for i := 0; i < s.Len(); i++ {
		vals[i] = s.Index(i).Interface()
	}

	return vals
}

func nativeStructToValue(val any) Value {
	obj := NewObj(nil)

	valType := reflect.TypeOf(val)

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldVal := reflect.ValueOf(val).Field(i).Interface()

		obj.Pairs[field.Name] = NativeToValue(fieldVal)
	}

	return obj
}

func nativeSliceToArrayValue(slice []any) *Arr {
	arr := new(Arr)
	arr.Elements = make([]Value, len(slice))
	for i := range slice {
		arr.Elements[i] = NativeToValue(slice[i])
	}

	return arr
}
