package arango_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB CRUD", func() {
	var _ = BeforeEach(func() {
		err := gormDB.AutoMigrate(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Create", func() {
		It("inserts a record into the collection", func() {
			user := &User{
				Name:  "Joselito",
				Email: "joselitofilhoo@gmail.com",
			}
			tx := gormDB.Create(user)
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

	Context("Find", func() {
		var user User

		BeforeEach(func() {
			newUser := &User{
				Name:  "Joselito",
				Email: "joselitofilhoo@gmail.com",
			}
			tx := gormDB.Create(newUser)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user = *tx.Statement.Dest.(*User)
			Expect(user.ID).NotTo(BeZero())
		})

		It("retrieves the user by id", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf("{\"ID\": %d}", user.ID))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(getUser.ID).To(Equal(user.ID))
			Expect(getUser.Name).To(Equal(user.Name))
			Expect(getUser.Email).To(Equal(user.Email))
			Expect(getUser.CreatedAt).NotTo(BeZero())
			Expect(getUser.UpdatedAt).NotTo(BeZero())
			Expect(getUser.DeletedAt).To(BeZero())
		})
	})

	Context("First", func() {
		var user1 User
		var user2 User

		BeforeEach(func() {
			newUser := &User{Name: "Joselito", Email: "joselitofilhoo@gmail.com"}
			tx := gormDB.Create(newUser)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user1 = *tx.Statement.Dest.(*User)
			Expect(user1.ID).NotTo(BeZero())

			newUser = &User{Name: "Joselito", Email: "joselito@gmail.com"}
			tx = gormDB.Create(newUser)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user2 = *tx.Statement.Dest.(*User)
			Expect(user2.ID).NotTo(BeZero())
		})

		It("retrieves the first record by name", func() {
			var getUser User
			tx := gormDB.First(&getUser, fmt.Sprintf("{\"Name\": \"Joselito\"}"))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(getUser.Email).To(Equal(user1.Email))
			Expect(getUser.Email).NotTo(Equal(user2.Email))
		})
	})
})