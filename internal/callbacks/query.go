package callbacks

import (
	"github.com/joselitofilho/gorm/driver/arango/internal/conn"
	"gorm.io/gorm"
)

func Query(db *gorm.DB) {
	if db.Error == nil {
		buildAQL(db)

		db.Statement.ConnPool.(*conn.ConnPool).Result = db.Statement.Dest
		_, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
		if err != nil {
			db.AddError(err)
			return
		}
	}
}

func buildAQL(db *gorm.DB) {
	db.Statement.SQL.WriteString("FOR doc IN @@collection FILTER doc.DeleteAt == null ")
	db.Statement.Vars = append(db.Statement.Vars, "@collection", db.Statement.Table)

	// TODO: sort.
	db.Statement.Build("WHERE", "LIMIT")

	// TODO: select.
	db.Statement.SQL.WriteString(" RETURN doc")
}
