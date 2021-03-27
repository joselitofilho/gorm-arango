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

	It("updates user's name field", func() {
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
		Expect(getUser.UpdatedAt.After(userJoselito.UpdatedAt)).To(BeTrue())
	})

	When("passing a user instance with multiple fields set", func() {
		It("updates all fields, including non-zero fields", func() {
			newName := "Ze"
			newAge := 33
			tx := gormDB.Model(&userJoselito).Updates(User{Name: newName, Age: uint(newAge)})
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.RowsAffected).To(BeEquivalentTo(1))

			var getUser User
			tx = gormDB.First(&getUser)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(getUser.ID).To(Equal(userJoselito.ID))
			Expect(getUser.Email).To(BeEmpty())
			Expect(getUser.Name).To(Equal(newName))
			Expect(getUser.Age).To(BeEquivalentTo(newAge))
			Expect(getUser.CreatedAt).To(Equal(userJoselito.CreatedAt))
			Expect(getUser.UpdatedAt.After(userJoselito.UpdatedAt)).To(BeTrue())
			Expect(getUser.DeletedAt).To(Equal(userJoselito.DeletedAt))
		})
	})

	When("passing a map with multiple user's fields", func() {
		It("updates the fields just passed", func() {
			newName := "Ze"
			newAge := 33
			tx := gormDB.Model(&userJoselito).Updates(map[string]interface{}{"Name": newName, "Age": newAge})
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.RowsAffected).To(BeEquivalentTo(1))

			var getUser User
			tx = gormDB.First(&getUser)
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(getUser.ID).To(Equal(userJoselito.ID))
			Expect(getUser.Email).To(Equal(userJoselito.Email))
			Expect(getUser.Name).To(Equal(newName))
			Expect(getUser.Age).To(BeEquivalentTo(newAge))
			Expect(getUser.CreatedAt).To(Equal(userJoselito.CreatedAt))
			Expect(getUser.UpdatedAt.After(userJoselito.UpdatedAt)).To(BeTrue())
			Expect(getUser.DeletedAt).To(Equal(userJoselito.DeletedAt))
		})
	})
})
