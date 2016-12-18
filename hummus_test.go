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

		It("deals with nested structs", func() {
			input := struct {
				Brand       string `gabs:"brand"`
				Type        string `gabs:"type"`
				Tasty       bool   `gabs:"tasty"`
				AddrStreet  string `gabs:"manufacturer_address.street"`
				AddrZipCode string `gabs:"manufacturer_address.zipcode"`
				AddrState   string `gabs:"manufacturer_address.state"`
			}{
				Brand:       "sabra",
				Type:        "jalapeno",
				Tasty:       true,
				AddrStreet:  "1234 Fake St.",
				AddrZipCode: "94040",
				AddrState:   "CA",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`
			{
				"brand": "sabra",
				"type": "jalapeno",
				"tasty": true,
				"manufacturer_address": {
					"street": "1234 Fake St.",
					"zipcode": "94040",
					"state": "CA"
				}
			}`))
		})

		It("deals with simple arrays", func() {
			input := struct {
				Brand0 string `gabs:"brands[0]"`
				Brand1 string `gabs:"brands[1]"`
				Brand2 string `gabs:"brands[2]"`
			}{
				Brand0: "sabra",
				Brand1: "athenos",
				Brand2: "whole-foods",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`
			{
				"brands": ["sabra", "athenos", "whole-foods"]
			}`))
		})
	})
})
