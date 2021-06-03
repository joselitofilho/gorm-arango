package arango_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Find", func() {
	var joselitoUser User
	var lucasUser User

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
		By("dropping collection", func() {
			err := gormDB.Migrator().DropTable(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		By("preparing collection", func() {
			err := gormDB.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		By("creating users", func() {
			joselitoUser = fnCreateUser("Joselito", "joselitofilhoo@gmail.com")
			Expect(joselitoUser).NotTo(BeNil())
			Expect(joselitoUser.Name).To(Equal("Joselito"))
			Expect(joselitoUser.DeletedAt).To(BeZero())

			lucasUser = fnCreateUser("Lucas Saraiva", "lucas.saraiva019@gmail.com")
			Expect(lucasUser).NotTo(BeNil())
			Expect(lucasUser.Name).To(Equal("Lucas Saraiva"))
			Expect(lucasUser.DeletedAt).To(BeZero())
		})
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("retrieves the user by ID", func() {
		var getUser User
		tx := gormDB.Find(&getUser, fmt.Sprintf(`{"ID": %d}`, joselitoUser.ID))
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(1))
		Expect(getUser.ID).To(Equal(joselitoUser.ID))
		Expect(getUser.Name).To(Equal(joselitoUser.Name))
		Expect(getUser.Email).To(Equal(joselitoUser.Email))
		Expect(getUser.CreatedAt).NotTo(BeZero())
		Expect(getUser.UpdatedAt).NotTo(BeZero())
		Expect(getUser.DeletedAt).To(BeZero())
	})

	When("there is more than one user in the result query", func() {
		It("retrieves the two users with ID > 0", func() {
			var users []User
			tx := gormDB.Find(&users, fmt.Sprintf(`{"ID": {"$gt": %d}}`, 0))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.RowsAffected).To(BeEquivalentTo(2))
			Expect(users).To(HaveLen(2))
			Expect(users[0].Name).To(Equal(joselitoUser.Name))
		})
	})

	When("operator is $eq", func() {
		It("returns user to the passed user.ID", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf(`{"ID": {"$eq": %d}}`, joselitoUser.ID))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.RowsAffected).To(BeEquivalentTo(1))
			Expect(getUser.ID).To(Equal(joselitoUser.ID))
		})
	})

	When("operator is $gt", func() {
		It("returns the first user with an ID greater than zero", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf(`{"ID": {"$gt": %d}}`, 0))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.RowsAffected).To(BeEquivalentTo(1))
			Expect(getUser.ID).To(Equal(joselitoUser.ID))
		})
	})

	When("user.ID passed does not exist in the database", func() {
		It("returns a 'not found' error", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf(`{"ID": %d}`, joselitoUser.ID+1))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).NotTo(BeNil())
			Expect(tx.Error.Error()).To(Equal("Document not found"))
			Expect(tx.RowsAffected).To(BeZero())
		})
	})

	When("ID passed is invalid", func() {
		It("returns a 'not found' error", func() {
			var getUser User
			tx := gormDB.Find(&getUser, fmt.Sprintf(`{"ID": %d}`, -1))
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).NotTo(BeNil())
			Expect(tx.Error.Error()).To(Equal("Document not found"))
			Expect(tx.RowsAffected).To(BeZero())
		})
	})
})
