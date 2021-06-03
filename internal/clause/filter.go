package clause

import (
	"fmt"
	"strings"

	"github.com/joselitofilho/gorm-arango/internal/transformers"
	gormClause "gorm.io/gorm/clause"
)

// Filter filter clause
type Filter struct {
	Exprs []gormClause.Expression
}

// Name filter clause name
func (filter Filter) Name() string {
	return "FILTER"
}

// Build build where clause
func (filter Filter) Build(builder gormClause.Builder) {
	if len(filter.Exprs) > 0 {
		expr := filter.Exprs[0].(gormClause.Expr)
		bindingFields := map[string]interface{}{}
		if sql, err := formattedFilter(expr.SQL, bindingFields); err == nil {
			builder.WriteString(sql)
			vars := transformers.MapToSlice(bindingFields)
			for _, v := range vars {
				builder.AddVar(builder, v)
			}
		}
	}
}

func prepareFieldBindings(fieldName string, bindingFields map[string]interface{}) string {
	fields := strings.Split(fieldName, ".")
	result := strings.Join(fields, ".@")
	for _, field := range fields {
		bindingFields[field] = field
	}

	return result
}

// formattedFilter returns formatted string and bindingFields.
func formattedFilter(query string, bindingFields map[string]interface{}) (string, error) {
	filters, err := transformers.GetFiltersByQuery(query, nil)
	if err != nil {
		return "", err
	}

	formattedFilterSlice := []string{}
	for index, filter := range filters {
		fieldKey := fmt.Sprintf("field_filter_%d", index)
		operator, err := filter.GetOperator()
		if err != nil {
			return "", err
		}
		formattedFields := prepareFieldBindings(filter.Field, bindingFields)
		// TODO: We should create a field to customizer it.
		formattedFilterSlice = append(formattedFilterSlice, fmt.Sprintf("doc.@%s %s @%s", formattedFields, operator, fieldKey))
		bindingFields[fieldKey] = filter.Value
	}
	formattedFilter := ""
	if len(formattedFilterSlice) > 0 {
		formattedFilter = strings.Join(formattedFilterSlice, " AND ")
	}

	return formattedFilter, nil
}
