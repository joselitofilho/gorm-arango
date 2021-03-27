# GORM Arango Driver

[Arango](https://www.arangodb.com/) driver for GORM library. Checkout [gorm.io](https://gorm.io) for details.

## USAGE

```go
import (
    "gorm.io/gorm"
    arango "github.com/joselitofilho/gorm/driver/arango/pkg"
)

type User struct {
	gorm.Model
	Name  string
	Email string
	Age   uint
}

func main() {
    db, err := gorm.Open(arango.Open(&arango.Config{}), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    db.AutoMigrate(&User{})

    // Create
    db.Create(&User{Name: "Joselito", Email: "joselitofilhoo@gmail.com", Age: 32})

    // Read
    var user User
    db.Find(&user, "{\"ID\": 1}") // find user with ID = 1
    db.First(&user, "{\"Name\": Joselito}") // find first user with Name is Joselito

    // Update - update user's name to Zelito.
    user.Name = "Zelito"
    db.Save(&user)

    // Update - update user's name to Ze.
    db.Model(&user).Update("Name", "Ze")
    // Update - update multiple fields
    db.Model(&user).Updates(User{Name: "Ze", Age: 33}) // including non-zero fields. Updates user's name to Ze, age to 33 and email to empty.
    db.Model(&user).Updates(map[string]interface{}{"Name": "Ze", "Age": 33}) // updates just user's name to Ze and age to 33.
}
```

## Contributors

Checkout [Contribute](docs/CONTRIBUTING.md) for details.
