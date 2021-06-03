package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type APIUser struct {
	ID   uint
	Name string
}

var _ = Describe("ArangoDB Smart Select Fields", func() {
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
			tx := gormDB.Create(&User{Name: "Joselito", Email: "joselitofilhoo@gmail.com", Age: 33})
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

	It("retrieves specific fields automatically", func() {
		var apiUser APIUser
		tx := gormDB.Model(&User{}).Limit(1).Find(&apiUser)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(1))
		Expect(apiUser.ID).To(Equal(user.ID))
		Expect(apiUser.Name).To(Equal(user.Name))
	})
})
