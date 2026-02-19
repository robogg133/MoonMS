package wasm

func marshalParams(alloc WasmAlloc, params ...any) []uint64 {
	result := make([]uint64, 0, len(params))

	for _, v := range params {
		switch x := v.(type) {

		case nil:
			result = append(result, 0)

		case bool:
			if x {
				result = append(result, 1)
			} else {
				result = append(result, 0)
			}

		case int:
			result = append(result, uint64(x))
		case int8:
			result = append(result, uint64(x))
		case int16:
			result = append(result, uint64(x))
		case int32:
			result = append(result, uint64(x))
		case int64:
			result = append(result, uint64(x))

		case uint:
			result = append(result, uint64(x))
		case uint8: // byte
			result = append(result, uint64(x))
		case uint16:
			result = append(result, uint64(x))
		case uint32:
			result = append(result, uint64(x))
		case uint64:
			result = append(result, x)

		case string:
			data := []byte(x)
			ptr := alloc.Alloc(uint32(len(data)))
			alloc.Write(ptr, data)

			result = append(result,
				uint64(ptr),
				uint64(len(data)),
			)

		case []byte:
			ptr := alloc.Alloc(uint32(len(x)))
			alloc.Write(ptr, x)

			result = append(result,
				uint64(ptr),
				uint64(len(x)),
			)

		default:
			panic("unsupported param type")
		}
	}

	return result
}
