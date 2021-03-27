package transformers

import "encoding/json"

func MapToSlice(m map[string]interface{}) []interface{} {
	slice := make([]interface{}, len(m)*2)
	i := 0
	for k, v := range m {
		slice[i] = k
		slice[i+1] = v
		i += 2
	}
	return slice
}

func SliceToMap(input []interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(input)/2)
	var key string
	for i, v := range input {
		if i%2 == 0 {
			key = v.(string)
		} else {
			output[key] = v
		}
	}
	return output
}

func EntityToMap(entity interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	modelMap := map[string]interface{}{}
	if err := json.Unmarshal(data, &modelMap); err != nil {
		return nil, err
	}

	return modelMap, nil
}
