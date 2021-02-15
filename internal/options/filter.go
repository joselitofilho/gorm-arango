package options

import (
	"errors"
	"reflect"
)

const (
	DefaultOperator = "$eq"
)

// filterOperators supported operators in the query filters
var filterOperators = map[string]string{
	"$eq":  "==",
	"$gt":  ">",
	"$gte": ">=",
	"$in":  "IN",
	"$lt":  "<",
	"$lte": "<=",
	"$ne":  "!=",
	"$nin": "NOT IN",
}

// Filter structure to represents a query filter
type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

// GetOperator returns the ArangoDB operator based on the filter operator.
func (f *Filter) GetOperator() (string, error) {
	switch reflect.TypeOf(f.Value).Kind() {
	case reflect.Slice:
		if f.Operator == "$in" {
			return "ANY IN", nil
		} else if f.Operator == "$nin" {
			return "ALL NOT IN", nil
		}
	}

	operator, ok := filterOperators[f.Operator]
	if !ok {
		return "", errors.New("Invalid argument")
	}

	return operator, nil
}
