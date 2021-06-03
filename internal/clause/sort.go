package clause

import (
	"strings"

	gormClause "gorm.io/gorm/clause"
)

type Sort struct {
	Columns    []gormClause.OrderByColumn
	Expression gormClause.Expression
}

// Name where clause name
func (sort Sort) Name() string {
	return "SORT"
}

// Build build where clause
func (sort Sort) Build(builder gormClause.Builder) {
	if sort.Expression != nil {
		sort.Expression.Build(builder)
	} else {
		for idx, column := range sort.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}

			internalColumns := strings.Split(column.Column.Name, ",")
			for idy, internalCol := range internalColumns {
				if idy > 0 {
					builder.WriteString(", ")
				}
				// TODO: We should create a field to customizer it.
				builder.WriteString("doc.")
				builder.WriteString(strings.TrimSpace(internalCol))
				if column.Desc {
					builder.WriteString(" DESC")
				}
			}
		}
	}
}

// MergeClause merge order by clauses
func (sort Sort) MergeClause(clause *gormClause.Clause) {
	if v, ok := clause.Expression.(Sort); ok {
		for i := len(sort.Columns) - 1; i >= 0; i-- {
			if sort.Columns[i].Reorder {
				sort.Columns = sort.Columns[i:]
				clause.Expression = sort
				return
			}
		}

		copiedColumns := make([]gormClause.OrderByColumn, len(v.Columns))
		copy(copiedColumns, v.Columns)
		sort.Columns = append(copiedColumns, sort.Columns...)
	}

	clause.Expression = sort
}