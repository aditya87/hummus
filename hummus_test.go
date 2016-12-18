package hummus_test

import (
	"github.com/aditya87/hummus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hummus", func() {
	Describe("Marshal", func() {
		It("marshals a simple struct into JSON", func() {
			input := struct {
				Brand string `gabs:"brand"`
				Type  string `gabs:"type"`
				Tasty bool   `gabs:"tasty"`
			}{
				Brand: "sabra",
				Type:  "jalapeno",
				Tasty: true,
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`
			{
				"brand": "sabra",
				"type": "jalapeno",
				"tasty": true
			}`))
		})

		It("omits empty values", func() {
			input := struct {
				Brand string `gabs:"brand,omitempty"`
				Type  string `gabs:"type,omitempty"`
				Price int    `gabs:"price,omitempty"`
			}{
				Brand: "whole foods",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`
			{
				"brand": "whole foods"
			}`))
		})
	})
})
