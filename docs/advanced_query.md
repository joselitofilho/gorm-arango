# Advanced Query

## Group Conditions

Easier to write complicated AQL query with Group Conditions

```Go
type User struct {
	gorm.Model
	Name  string
	Email string
    Age   uint
}

var user User

db.Where(`{"id": {"$gt": 0}}"`).Where(db.Where(`{"age": {"$gte": 18}}`).Where(`{"age": {"$lt": 40}}`)).First(&getUser)
// FOR u IN users FILTER u.id > 0 AND (u.age >= 18 OR u.age < 40) RETURN u
```