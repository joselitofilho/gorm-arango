package clause

import (
	"gorm.io/gorm"
	gormClause "gorm.io/gorm/clause"
)

func RegisterDefaultClauses(db *gorm.DB) {
	db.ClauseBuilders["WHERE"] = func(c gormClause.Clause, builder gormClause.Builder) {
		if c.Expression != nil {
			filterClause := Filter{Exprs: c.Expression.(gormClause.Where).Exprs}

			if c.BeforeExpression != nil {
				c.BeforeExpression.Build(builder)
				builder.WriteByte(' ')
			}

			if c.Name != "" {
				builder.WriteString(filterClause.Name())
				builder.WriteByte(' ')
			}

			if c.AfterNameExpression != nil {
				c.AfterNameExpression.Build(builder)
				builder.WriteByte(' ')
			}

			c.Build(builder)
			// TODO: To be continued...
			// filterClause.Build(builder)

			if c.AfterExpression != nil {
				builder.WriteByte(' ')
				c.AfterExpression.Build(builder)
			}
		}
	}
}
