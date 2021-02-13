package transformers

import (
	"encoding/json"
	"reflect"

	"github.com/joselitofilho/gorm/driver/arango/internal/options"
)

// GetFiltersByQuery ...
func GetFiltersByQuery(query string, enumFieldValues map[string]map[string]int32) ([]options.Filter, error) {
	queryUnmarshalled := map[string]interface{}{}

	if err := json.Unmarshal([]byte(query), &queryUnmarshalled); err != nil {
		return nil, err
	}

	filterList := []options.Filter{}
	for field, data := range queryUnmarshalled {
		switch reflect.ValueOf(data).Kind() {
		case reflect.Map:
			for operator, value := range data.(map[string]interface{}) {
				if enumFieldValues != nil {
					if enumValues, ok := enumFieldValues[field]; ok {
						value = enumValues[value.(string)]
					}
				}
				filterList = append(filterList, options.Filter{Field: field, Operator: operator, Value: value})
			}
		default:
			if enumFieldValues != nil {
				if enumValues, ok := enumFieldValues[field]; ok {
					data = enumValues[data.(string)]
				}
			}
			filterList = append(filterList, options.Filter{Field: field, Operator: options.DefaultOperator, Value: data})
		}
	}

	return filterList, nil
}
