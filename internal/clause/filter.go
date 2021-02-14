package clause

import (
	"fmt"

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
		fmt.Println(expr.SQL)
		// TODO: To be continued...
	}
}
