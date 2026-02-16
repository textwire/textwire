package object

import (
	"reflect"
)

func NativeToObject(val any) Object {
	switch v := val.(type) {
	case string:
		return &Str{Value: v}
	case bool:
		return &Bool{Value: v}
	case float32:
		return &Float{Value: float64(v)}
	case float64:
		return &Float{Value: v}
	case int64:
		return &Int{Value: v}
	case int:
		return &Int{Value: int64(v)}
	case int8:
		return &Int{Value: int64(v)}
	case int16:
		return &Int{Value: int64(v)}
	case int32:
		return &Int{Value: int64(v)}
	case uint:
		return &Int{Value: int64(v)}
	case uint8:
		return &Int{Value: int64(v)}
	case uint16:
		return &Int{Value: int64(v)}
	case uint32:
		return &Int{Value: int64(v)}
	case uint64:
		return &Int{Value: int64(v)}
	case nil:
		return new(Nil)
	}

	valType := reflect.TypeOf(val)

	switch valType.Kind() {
	case reflect.Struct:
		return nativeStructToObject(val)
	case reflect.Slice:
		return nativeSliceToArrayObject(convertToInterfaceSlice(val))
	case reflect.Map:
		return nativeMapToObject(val)
	case reflect.Pointer:
		// NativeToObject is used recursively to handle pointers
		return NativeToObject(reflect.ValueOf(val).Elem().Interface())
	}

	return nil
}

func nativeMapToObject(val any) Object {
	obj := NewObj(nil)

	valValue := reflect.ValueOf(val)
	for _, key := range valValue.MapKeys() {
		obj.Pairs[key.String()] = NativeToObject(valValue.MapIndex(key).Interface())
	}

	return obj
}

func convertToInterfaceSlice(slice any) []any {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	result := make([]any, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}

	return result
}

func nativeStructToObject(val any) Object {
	obj := NewObj(nil)

	valType := reflect.TypeOf(val)

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		if !field.IsExported() {
			continue
		}

		fieldVal := reflect.ValueOf(val).Field(i).Interface()

		obj.Pairs[field.Name] = NativeToObject(fieldVal)
	}

	return obj
}

func nativeSliceToArrayObject(slice []any) *Array {
	arr := new(Array)
	arr.Elements = make([]Object, len(slice))
	for i := range slice {
		arr.Elements[i] = NativeToObject(slice[i])
	}

	return arr
}
