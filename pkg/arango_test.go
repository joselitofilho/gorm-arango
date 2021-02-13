package arango_test

import (
	"gorm.io/gorm"

	arango "github.com/joselitofilho/gorm/driver/arango/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

var _ = Describe("ArangoDB Dialector", func() {
	arangoDBTestConfig := newArangoDBTestConfig()

	It("connects to a ArangoDB server", func() {
		db, err := gorm.Open(arango.Open(arangoDBTestConfig), &gorm.Config{})
		Expect(db).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())
	})

	When("ArangoDB server already connected", func() {
		var db *gorm.DB

		BeforeEach(func() {
			_db, err := gorm.Open(arango.Open(arangoDBTestConfig), &gorm.Config{})
			Expect(err).NotTo(HaveOccurred())
			Expect(_db).NotTo(BeNil())
			db = _db

			err = db.Migrator().DropTable(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a collection", func() {
			err := db.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("inserts a record into the collection", func() {
			err := db.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())

			user := &User{
				Name:  "Joselito",
				Email: "joselitofilhoo@gmail.com",
			}
			tx := db.Create(user)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			newUser := tx.Statement.Dest.(*User)
			Expect(newUser.ID).NotTo(BeZero())
			Expect(newUser.Name).To(Equal(user.Name))
			Expect(newUser.Email).To(Equal(user.Email))
			Expect(newUser.CreatedAt).NotTo(BeZero())
			Expect(newUser.UpdatedAt).NotTo(BeZero())
			Expect(newUser.DeletedAt).To(BeZero())
		})
	})
})
