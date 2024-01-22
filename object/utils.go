package object

func nativeToObject(val interface{}) Object {
	switch v := val.(type) {
	case string:
		return &String{Value: v}
	case bool:
		return &Boolean{Value: v}
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
	case []string, []int64, []int, []int8, []int16, []int32, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []bool:
		return convertSlice(v)
	}

	return nil
}

func convertSlice(slice interface{}) *Array {
	arr := &Array{}

	switch v := slice.(type) {
	case []string:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []bool:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []float32:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []float64:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []int64:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []int32:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []int16:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []int8:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []int:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []uint64:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []uint32:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []uint16:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []uint8:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	case []uint:
		for _, val := range v {
			arr.Elements = append(arr.Elements, nativeToObject(val))
		}
	}

	return arr
}
