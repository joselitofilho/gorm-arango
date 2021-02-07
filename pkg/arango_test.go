package arango_test

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	arango "github.com/joselitofilho/gorm/driver/arango/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArangoDB Dialector", func() {
	logrus.SetLevel(logrus.DebugLevel)

	It("It connects to a ArangoDB server", func() {
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

		db, err := gorm.Open(arango.Open(arangodbConfig), &gorm.Config{})
		Expect(db).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())
	})
})
