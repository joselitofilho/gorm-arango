package arango_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Where", func() {
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
			tx := gormDB.Create(&User{Name: "Joselito", Email: "joselitofilhoo@gmail.com", Age: 33})
			Expect(tx).NotTo(BeNil())
			Expect(tx.Error).To(BeNil())
			Expect(tx.Statement.Dest).NotTo(BeNil())

			user1 = *tx.Statement.Dest.(*User)
			Expect(user1.ID).NotTo(BeZero())

			tx = gormDB.Create(&User{Name: "Ze", Email: "joselito@gmail.com", Age: 26})
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

	It("retrieves the first record by: id > 0 AND age >= 18 AND age < 40", func() {
		var getUser User
		tx := gormDB.Where(`{"ID": {"$gt": 0}}`).Where(gormDB.Where(`{"age": {"$gte": 18}}`).Where(`{"age": {"$lt": 40}}`)).First(&getUser)
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(tx.RowsAffected).To(BeEquivalentTo(1))
		Expect(getUser.Name).To(Equal(user1.Name))
	})
})
