package callbacks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/joselitofilho/gorm-arango/internal/conn"
	"github.com/joselitofilho/gorm-arango/internal/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newInstanceOfSliceType(arr interface{}) reflect.Type {
	val := reflect.ValueOf(arr)
	if val.Kind() == reflect.Ptr {
		return newInstanceOfSliceType(val.Elem().Interface())
	}
	if val.Kind() != reflect.Slice {
		// TODO: Implement
		return nil
	}
	return reflect.TypeOf(arr).Elem()
}

func Query(db *gorm.DB) {
	if db.Error == nil {
		buildAQL(db)

		elemType := newInstanceOfSliceType(db.Statement.Dest)

		isSlice := db.Statement.ReflectValue.Kind() == reflect.Slice || db.Statement.ReflectValue.Kind() == reflect.Array
		db.Statement.ConnPool.(*conn.ConnPool).Return = conn.ConnPoolReturn{
			Dest:     db.Statement.Dest,
			ElemType: elemType,
			IsSlice:  isSlice,
		}

		aql := db.Statement.SQL.String()
		if _, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, aql, db.Statement.Vars...); err != nil {
			db.AddError(fmt.Errorf("AQL = %s", aql))
			db.AddError(err)
			return
		}

		if isSlice {
			model := reflect.New(elemType).Interface()
			if err := scan(db, model); err != nil {
				db.AddError(err)
				return
			}
		} else {
			db.RowsAffected = 1
		}
	}
}

func buildAQL(db *gorm.DB) {
	// TODO: Think better about that.
	db.Statement.Build("ORDER BY")
	sort := db.Statement.SQL.String()
	sort = strings.ReplaceAll(sort, "SORT .", "")
	db.Statement.SQL.Reset()

	db.Statement.SQL.WriteString("FOR doc IN @@collection FILTER doc.DeleteAt == null ")
	db.Statement.Vars = append(db.Statement.Vars, "@collection", db.Statement.Table)

	db.Statement.Build("WHERE", "LIMIT")

	db.Statement.SQL.WriteString(sort)

	// TODO: select.
	// TODO: We should create a field to customizer it.
	db.Statement.SQL.WriteString(" RETURN doc")
}

func rowColumns(model interface{}) ([]string, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	var mapModel map[string]interface{}
	err = json.Unmarshal(data, &mapModel)
	if err != nil {
		return nil, err
	}

	columns := make([]string, 0)
	for key := range mapModel {
		columns = append(columns, key)
	}

	return columns, nil
}

// This method is based on gorm.Scan() method.
func scan(db *gorm.DB, model interface{}) error {
	db.RowsAffected = 0

	columns, err := rowColumns(model)
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))

	switch dest := db.Statement.Dest.(type) {
	case map[string]interface{}, *map[string]interface{},
		*[]map[string]interface{},
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *uintptr,
		*float32, *float64,
		*bool, *string, *time.Time,
		*sql.NullInt32, *sql.NullInt64, *sql.NullFloat64,
		*sql.NullBool, *sql.NullString, *sql.NullTime:
		return errors.ErrMethodNotImplemented(dest)
	default:
		Schema := db.Statement.Schema

		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			var (
				reflectValueType = db.Statement.ReflectValue.Type().Elem()
				isPtr            = reflectValueType.Kind() == reflect.Ptr
				fields           = make([]*schema.Field, len(columns))
			)

			if isPtr {
				reflectValueType = reflectValueType.Elem()
			}

			db.Statement.ReflectValue.Set(reflect.MakeSlice(db.Statement.ReflectValue.Type(), 0, 20))

			if Schema != nil {
				if reflectValueType != Schema.ModelType && reflectValueType.Kind() == reflect.Struct {
					return fmt.Errorf("Invalid schema %v.", Schema)
				}

				for idx, column := range columns {
					if field := Schema.LookUpField(column); field != nil && field.Readable {
						fields[idx] = field
					} else {
						values[idx] = &sql.RawBytes{}
					}
				}
			}

			// pluck values into slice of data
			isPluck := false
			if len(fields) == 1 {
				if _, ok := reflect.New(reflectValueType).Interface().(sql.Scanner); ok || // is scanner
					reflectValueType.Kind() != reflect.Struct || // is not struct
					Schema.ModelType.ConvertibleTo(schema.TimeReflectType) { // is time
					isPluck = true
				}
			}

			for _, row := range db.Statement.ConnPool.(*conn.ConnPool).Return.Dest.([]interface{}) {
				db.RowsAffected++

				data, _ := json.Marshal(row)
				var mapModel map[string]interface{}
				json.Unmarshal(data, &mapModel)

				elem := reflect.New(reflectValueType)
				if isPluck {
					db.AddError(fmt.Errorf("Error scanning row: %v", row))
				} else {
					for idx, field := range fields {
						if field != nil {
							values[idx] = mapModel[field.Name]
						}
					}

					for idx, field := range fields {
						if field != nil {
							field.Set(elem, values[idx])
						}
					}
				}

				if isPtr {
					db.Statement.ReflectValue.Set(reflect.Append(db.Statement.ReflectValue, elem))
				} else {
					db.Statement.ReflectValue.Set(reflect.Append(db.Statement.ReflectValue, elem.Elem()))
				}
			}
		}
	}

	return nil
}
