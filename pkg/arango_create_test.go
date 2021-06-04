package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Create", func() {
	var _ = BeforeEach(func() {
		By("dropping collection", func() {
			err := gormDB.Migrator().DropTable(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		By("preparing collection", func() {
			err := gormDB.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("inserts a record into the collection", func() {
		user := User{
			Name:  "Joselito",
			Email: "joselitofilhoo@gmail.com",
		}
		tx := gormDB.Create(&user)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.Statement.RowsAffected).To(BeEquivalentTo(1))

		Expect(user.ID).NotTo(BeZero())
		Expect(user.Name).To(Equal(user.Name))
		Expect(user.Email).To(Equal(user.Email))
		Expect(user.CreatedAt).NotTo(BeZero())
		Expect(user.UpdatedAt).NotTo(BeZero())
		Expect(user.DeletedAt).To(BeZero())
	})
})
