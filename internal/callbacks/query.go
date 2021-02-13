package callbacks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm/driver/arango/internal/conn"
	"github.com/joselitofilho/gorm/driver/arango/internal/transformers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArangoBuilder struct {
	strings.Builder
}

func (b *ArangoBuilder) WriteQuoted(field interface{}) {

}
func (b *ArangoBuilder) AddVar(clause.Writer, ...interface{}) {

}

func (b *ArangoBuilder) WriteByte(byte) error {
	return nil
}

func (b *ArangoBuilder) WriteString(string) (int, error) {
	return 0, nil
}

func Query(db *gorm.DB) {
	if db.Error == nil {
		connection := db.Statement.ConnPool.(*conn.ConnPool)
		if collection, err := connection.Database.Collection(db.Statement.Context, db.Statement.Table); err == nil {
			query, vars, err := BuildAQL(db)
			if err != nil {
				db.AddError(err)
				return
			}

			cursor, err := collection.Database().Query(db.Statement.Context, query, vars)
			if err != nil {
				db.AddError(err)
				return
			}
			defer cursor.Close()
			_, err = cursor.ReadDocument(db.Statement.Context, db.Statement.Dest)
			if driver.IsNoMoreDocuments(err) {
				db.AddError(errors.New("Document not found"))
				return
			} else if err != nil {
				db.AddError(err)
				return
			}
		}
	}
}

// BuildAQL ...
func BuildAQL(db *gorm.DB) (query string, bindingFields map[string]interface{}, err error) {
	// FOR doc IN @@collection FILTER doc.DeleteAt == null AND doc.ID == @docID RETURN doc

	bindingFields = map[string]interface{}{"@collection": db.Statement.Table}

	db.Statement.Build("OFFSET")
	offsetStr := strings.TrimSpace(strings.ReplaceAll(db.Statement.SQL.String(), "OFFSET", ""))
	db.Statement.SQL.Reset()
	offset := 0
	if len(offsetStr) > 0 {
		if offset, err = strconv.Atoi(offsetStr); err != nil {
			offset = 0
		}
	}
	bindingFields["offset"] = offset

	db.Statement.Build("LIMIT")
	limitStr := strings.TrimSpace(strings.ReplaceAll(db.Statement.SQL.String(), "LIMIT", ""))
	db.Statement.SQL.Reset()
	limit := 250
	if len(limitStr) > 0 {
		if limit, err = strconv.Atoi(limitStr); err != nil {
			limit = 250
		}
	}
	bindingFields["limit"] = limit

	db.Statement.Build("WHERE")
	whereStr := strings.TrimSpace(strings.Replace(db.Statement.SQL.String(), "WHERE", "", 1))
	filters := ""
	if len(whereStr) > 0 {
		filters, err = formattedFilter(whereStr, bindingFields)
		if err != nil {
			return
		}
	}

	query = fmt.Sprintf("FOR doc IN @@collection FILTER doc.DeleteAt == null %s LIMIT @offset, @limit RETURN doc", filters)

	// TODO: limit, offset and sort.

	return
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
		formattedFilterSlice = append(formattedFilterSlice, fmt.Sprintf("doc.@%s %s @%s", formattedFields, operator, fieldKey))
		bindingFields[fieldKey] = filter.Value
	}
	formattedFilter := ""
	if len(formattedFilterSlice) > 0 {
		formattedFilter = "FILTER " + strings.Join(formattedFilterSlice, " AND ")
	}

	return formattedFilter, nil
}
