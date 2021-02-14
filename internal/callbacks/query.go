package callbacks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/joselitofilho/gorm/driver/arango/internal/aql"
	"github.com/joselitofilho/gorm/driver/arango/internal/transformers"
	"gorm.io/gorm"
)

func Query(db *gorm.DB) {
	if db.Error == nil {
		query, vars, err := BuildAQL(db)
		if err != nil {
			db.AddError(err)
			return
		}

		// TODO: To be continued...
		db.Statement.Vars = transformers.MapToSlice(vars)
		db.Statement.Vars = append(db.Statement.Vars, db.Statement.Dest)
		result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, query, db.Statement.Vars...)
		if err != nil {
			db.AddError(err)
			return
		}
		db.Statement.Dest = result.(aql.Result).Result()
	}
}

// BuildAQL ...
func BuildAQL(db *gorm.DB) (query string, bindingFields map[string]interface{}, err error) {
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
	filters := db.Statement.SQL.String()
	filterVars := map[string]interface{}{}
	transformers.SliceToMap(db.Statement.Vars, filterVars)
	for k, v := range filterVars {
		bindingFields[k] = v
	}

	// TODO: To be continued...
	query = fmt.Sprintf("FOR doc IN @@collection FILTER doc.DeleteAt == null %s LIMIT @offset, @limit RETURN doc", filters)
	fmt.Println(query, db.Statement.Vars)

	// TODO: sort.

	return
}
