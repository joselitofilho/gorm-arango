package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Update", func() {
	var userJoselito User

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

			userJoselito = *tx.Statement.Dest.(*User)
			Expect(userJoselito.ID).NotTo(BeZero())
		})
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("update user's name field", func() {
		newName := "Ze"
		tx := gormDB.Model(&userJoselito).Update("Name", newName)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(1))

		var getUser User
		tx = gormDB.First(&getUser)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(getUser.Email).To(Equal(userJoselito.Email))
		Expect(getUser.Name).To(Equal(newName))
	})
})
