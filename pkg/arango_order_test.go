package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Order", func() {
	var user1 User
	var user2 User
	var user3 User

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
			tx := gormDB.Create(&User{Name: "Joselito", Email: "joselitofilhoo@gmail.com", Age: 26})
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user1 = *tx.Statement.Dest.(*User)
			Expect(user1.ID).NotTo(BeZero())

			tx = gormDB.Create(&User{Name: "Ze", Email: "joselito@gmail.com", Age: 33})
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user2 = *tx.Statement.Dest.(*User)
			Expect(user2.ID).NotTo(BeZero())

			tx = gormDB.Create(&User{Name: "Saraiva", Email: "lucassaraiva@gmail.com", Age: 26})
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user3 = *tx.Statement.Dest.(*User)
			Expect(user3.ID).NotTo(BeZero())
		})
	})

	var _ = AfterEach(func() {
		err := gormDB.Migrator().DropTable(&User{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("retrieves records in a specify order", func() {
		var users []User
		tx := gormDB.Order("Age desc, Name").Find(&users)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(3))
		Expect(users).To(HaveLen(3))
		Expect(users[0].Name).To(Equal(user2.Name))
		Expect(users[1].Name).To(Equal(user1.Name))
		Expect(users[2].Name).To(Equal(user3.Name))
	})
})
