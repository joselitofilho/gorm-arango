package arango_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Find", func() {
	var user User

	var _ = BeforeEach(func() {
		By("preparing collection", func() {
			err := gormDB.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		By("creating registers", func() {
			newUser := &User{
				Name:  "Joselito",
				Email: "joselitofilhoo@gmail.com",
			}
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

	It("retrieves the user by id", func() {
		var getUser User
		tx := gormDB.Find(&getUser, fmt.Sprintf("{\"ID\": %d}", user.ID))
		Expect(tx).NotTo(BeNil())
		Expect(tx.Error).To(BeNil())
		Expect(getUser.ID).To(Equal(user.ID))
		Expect(getUser.Name).To(Equal(user.Name))
		Expect(getUser.Email).To(Equal(user.Email))
		Expect(getUser.CreatedAt).NotTo(BeZero())
		Expect(getUser.UpdatedAt).NotTo(BeZero())
		Expect(getUser.DeletedAt).To(BeZero())
	})
})
