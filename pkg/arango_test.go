package arango_test

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	arango "github.com/joselitofilho/gorm/driver/arango/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

var _ = Describe("ArangoDB Dialector", func() {
	logrus.SetLevel(logrus.DebugLevel)

	arangodbUri := os.Getenv("ARANGODB_URI")
	if arangodbUri == "" {
		arangodbUri = "http://localhost:8529"
	}
	arangodbConfig := &arango.Config{
		URI:                  arangodbUri,
		User:                 "user",
		Password:             "password",
		Database:             "gorm-arango-test",
		Timeout:              120,
		MaxConnectionRetries: 10,
	}

	It("connects to a ArangoDB server", func() {
		db, err := gorm.Open(arango.Open(arangodbConfig), &gorm.Config{})
		Expect(db).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())
	})

	When("ArangoDB server already connected", func() {
		var db *gorm.DB

		BeforeEach(func() {
			_db, err := gorm.Open(arango.Open(arangodbConfig), &gorm.Config{})
			Expect(err).NotTo(HaveOccurred())
			Expect(_db).NotTo(BeNil())
			db = _db

			err = db.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		FIt("Creating a collection", func() {
			tx := db.Create(&User{
				Name:  "Joselito",
				Email: "joselitofilhoo@gmail.com",
			})
			Expect(tx).NotTo(BeNil())
		})
	})
})
