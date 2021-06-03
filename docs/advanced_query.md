# Advanced Query

## Smart Select Fields

GORM allows select specific fields with Select, if you often use this in your application, maybe you want to define a smaller struct for API usage which can select specific fields automatically, for example:

```Go
type User struct {
	gorm.Model
	Name  string
	Email string
    Age   uint
}

type APIUser struct {
	ID   uint
	Name string
}

var apiUser APIUser

// Select `id`, `name` automatically when querying
db.Model(&User{}).Limit(1).Find(&apiUser)
// FOR u IN users LIMIT 1 RETURN { ID: u.ID, Name: u.Name }
```

## Group Conditions

Easier to write complicated AQL query with Group Conditions, for example:

```Go
type User struct {
	gorm.Model
	Name  string
	Email string
    Age   uint
}

var user User

db.Where(`{"ID": {"$gt": 0}}"`).Where(db.Where(`{"age": {"$gte": 18}}`).Where(`{"age": {"$lt": 40}}`)).First(&user)
// FOR u IN users FILTER u.ID > 0 AND (u.age >= 18 OR u.age < 40) RETURN u
```