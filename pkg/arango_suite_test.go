package arango_test

import (
	"context"
	"os"
	"testing"
	"time"

	arango "github.com/joselitofilho/gorm/driver/arango/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func setupContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(120*time.Second))
}

func newArangoDBTestConfig() *arango.Config {
	arangodbUri := os.Getenv("ARANGODB_URI")
	if arangodbUri == "" {
		arangodbUri = "http://localhost:8529"
	}
	return &arango.Config{
		URI:                  arangodbUri,
		User:                 "user",
		Password:             "password",
		Database:             "gorm-arango-test",
		Timeout:              120,
		MaxConnectionRetries: 10,
	}
}

var _ = BeforeSuite(func() {
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

	db, err := gorm.Open(arango.Open(arangodbConfig), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())
	dialector := db.Dialector.(arango.Dialector)
	Expect(dialector).NotTo(BeNil())
	err = dialector.Database.Remove(context.Background())
	Expect(err).NotTo(HaveOccurred())
})

func TestArangodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArangoDB Suite")
}
