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

type User struct {
	gorm.Model
	Name  string
	Email string
}

var gormDB *gorm.DB

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

	By("Connecting to the ArangoDB", func() {
		db, err := gorm.Open(arango.Open(arangodbConfig), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())
		gormDB = db
	})
})

var _ = AfterSuite(func() {
	dialector := gormDB.Dialector.(arango.Dialector)
	Expect(dialector).NotTo(BeNil())
	err := dialector.Database.Remove(context.Background())
	Expect(err).NotTo(HaveOccurred())
})

func TestArangodb(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArangoDB Suite")
}
