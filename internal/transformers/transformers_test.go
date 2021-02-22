package transformers_test

import (
	"github.com/joselitofilho/gorm/driver/arango/internal/transformers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transformers", func() {
	It("converts map to slice", func() {
		m := map[string]interface{}{
			"ID":    1234,
			"Name":  "Joselito",
			"Email": "joselitofilhoo@gmail.com",
		}
		slice := transformers.MapToSlice(m)
		Expect(slice).To(HaveLen(6))
	})
	It("converts slice to map", func() {
		slice := []interface{}{"ID", 1234, "Name", "Joselito", "Email", "joselitofilhoo@gmail.com"}
		m := transformers.SliceToMap(slice)
		Expect(m).To(HaveLen(3))
	})
})
