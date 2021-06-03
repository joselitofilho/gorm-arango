# Query

## Order

Specify order when retrieving records from the database, for example:

```Go
type User struct {
	gorm.Model
	Name  string
	Email string
    Age   uint
}

var users []User

db.Order("Age desc, Name").Find(&users)
// FOR u IN users SORT u.age DESC, u.name ASC RETURN u
```