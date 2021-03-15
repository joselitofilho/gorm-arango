package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB First", func() {
	var user User

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

			user = *tx.Statement.Dest.(*User)
			Expect(user.ID).NotTo(BeZero())
		})
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("update all fields", func() {
		newName := "Zelito"
		user.Name = newName
		tx := gormDB.Save(&user)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(1))

		var getUser User
		tx = gormDB.First(&getUser)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(getUser.Email).To(Equal(user.Email))
		Expect(getUser.Name).To(Equal(newName))
	})
})
