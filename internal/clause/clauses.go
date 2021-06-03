package clause

import (
	"gorm.io/gorm"
	gormClause "gorm.io/gorm/clause"
)

func RegisterDefaultClauses(db *gorm.DB) {
	db.ClauseBuilders["WHERE"] = filterClauseBuilder
	db.ClauseBuilders["ORDER BY"] = sortClauseBuilder
	db.ClauseBuilders["LIMIT"] = limitClauseBuilder
}

func filterClauseBuilder(c gormClause.Clause, builder gormClause.Builder) {
	if c.Expression != nil {
		filterClause := Filter{Exprs: c.Expression.(gormClause.Where).Exprs}

		if c.BeforeExpression != nil {
			c.BeforeExpression.Build(builder)
			builder.WriteByte(' ')
		}

		builder.WriteString(filterClause.Name())
		builder.WriteByte(' ')

		if c.AfterNameExpression != nil {
			c.AfterNameExpression.Build(builder)
			builder.WriteByte(' ')
		}

		filterClause.Build(builder)

		if c.AfterExpression != nil {
			builder.WriteByte(' ')
			c.AfterExpression.Build(builder)
		}
	}
}

func sortClauseBuilder(c gormClause.Clause, builder gormClause.Builder) {
	if c.Expression != nil {
		gormOrderByClause := c.Expression.(gormClause.OrderBy)
		sortClause := Sort{Columns: gormOrderByClause.Columns, Expression: gormOrderByClause.Expression}

		if c.BeforeExpression != nil {
			c.BeforeExpression.Build(builder)
			builder.WriteByte(' ')
		}

		builder.WriteString(sortClause.Name())
		builder.WriteByte(' ')

		if c.AfterNameExpression != nil {
			c.AfterNameExpression.Build(builder)
			builder.WriteByte(' ')
		}

		sortClause.Build(builder)

		if c.AfterExpression != nil {
			builder.WriteByte(' ')
			c.AfterExpression.Build(builder)
		}
	}
}

func limitClauseBuilder(c gormClause.Clause, builder gormClause.Builder) {
	if c.Expression != nil {
		gormLimitClause := c.Expression.(gormClause.Limit)
		limitClause := Limit{Limit: gormLimitClause.Limit, Offset: gormLimitClause.Offset}

		if c.BeforeExpression != nil {
			c.BeforeExpression.Build(builder)
			builder.WriteByte(' ')
		}

		if c.AfterNameExpression != nil {
			c.AfterNameExpression.Build(builder)
			builder.WriteByte(' ')
		}

		limitClause.Build(builder)

		if c.AfterExpression != nil {
			builder.WriteByte(' ')
			c.AfterExpression.Build(builder)
		}
	}
}
