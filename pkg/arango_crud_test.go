package arango_test

import (
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

})
