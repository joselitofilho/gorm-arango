package transformers

func MapToSlice(m map[string]interface{}) (slice []interface{}) {
	slice = make([]interface{}, len(m)*2)
	i := 0
	for k, v := range m {
		slice[i] = k
		slice[i+1] = v
		i += 2
	}
	return
}

func SliceToMap(input []interface{}, output map[string]interface{}) {
	var key string
	for i, v := range input {
		if i%2 == 0 {
			key = v.(string)
		} else {
			output[key] = v
		}
	}
	return
}
