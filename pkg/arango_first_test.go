package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB First", func() {
	var user1 User
	var user2 User

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
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("retrieves the first record by name", func() {
		var getUser User
		tx := gormDB.First(&getUser, `{"Name": "Joselito"}`)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(1))
		Expect(getUser.Email).To(Equal(user1.Email))
		Expect(getUser.Email).NotTo(Equal(user2.Email))
	})

	When("not exists the user for the informed condition", func() {
		It("returns a 'not found' error", func() {
			var getUser User
			tx := gormDB.First(&getUser, `{"Name": "Ze"}`)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).NotTo(BeNil())
			Expect(tx.Error.Error()).To(ContainSubstring("not found"))
			Expect(tx.RowsAffected).To(BeZero())
			Expect(getUser.Name).To(Equal(""))
		})
	})
})
