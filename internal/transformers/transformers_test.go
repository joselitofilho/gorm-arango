package transformers_test

import (
	"github.com/joselitofilho/gorm-arango/internal/transformers"
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
	It("converts entity to map", func() {
		obj := struct {
			Name  string
			Email string
		}{
			Name:  "Joselito",
			Email: "joselitofilhoo@gmail.com",
		}
		m, err := transformers.EntityToMap(&obj)
		Expect(err).To(BeNil())
		Expect(m).To(HaveLen(2))
		Expect(m["Name"]).To(Equal("Joselito"))
		Expect(m["Email"]).To(Equal("joselitofilhoo@gmail.com"))
	})
})
