package arango_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArangodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArangoDB Suite")
}
