package arango_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Find", func() {
	var user User

	fnCreateUser := func(name, email string) User {
		newUser := &User{
			Name:  name,
			Email: email,
		}
		tx := gormDB.Create(newUser)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.Statement.Dest).NotTo(BeNil())

		user := tx.Statement.Dest.(*User)
		Expect(user.ID).NotTo(BeZero())
		return *user
	}

	var _ = BeforeEach(func() {
		By("preparing collection", func() {
			err := gormDB.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		By("creating registers", func() {
			user = fnCreateUser("Joselito", "joselitofilhoo@gmail.com")
		})
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("retrieves the user by ID", func() {
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

	When("operator is $eq", func() {
		It("returns user to the passed user.ID", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf("{\"ID\": {\"$eq\": %d}}", user.ID))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(getUser.ID).To(Equal(user.ID))
		})
	})

	When("operator is $gt", func() {
		It("returns the first user with an ID greater than zero", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf("{\"ID\": {\"$gt\": %d}}", 0))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(getUser.ID).To(Equal(user.ID))
		})
	})

	When("user.ID passed does not exist in the database", func() {
		It("returns a 'not found' error", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf("{\"ID\": %d}", user.ID+1))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).NotTo(BeNil())
			Expect(tx.Error.Error()).To(Equal("Document not found"))
		})
	})

	When("ID passed is invalid", func() {
		It("returns a 'not found' error", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf("{\"ID\": %d}", -1))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).NotTo(BeNil())
			Expect(tx.Error.Error()).To(Equal("Document not found"))
		})
	})
})
