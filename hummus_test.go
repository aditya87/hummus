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
				Brand string `hummus:"brand"`
				Type  string `hummus:"type"`
				Tasty bool   `hummus:"tasty"`
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
				Brand string `hummus:"brand,omitempty"`
				Type  string `hummus:"type,omitempty"`
				Price int    `hummus:"price,omitempty"`
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
				Brand       string `hummus:"brand"`
				Type        string `hummus:"type"`
				Tasty       bool   `hummus:"tasty"`
				AddrStreet  string `hummus:"manufacturer_address.street"`
				AddrZipCode string `hummus:"manufacturer_address.zipcode"`
				AddrState   string `hummus:"manufacturer_address.state"`
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
				Brand0 string `hummus:"brands[0]"`
				Brand1 string `hummus:"brands[1]"`
				Brand2 string `hummus:"brands[2]"`
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

		It("deals with simple nested arrays", func() {
			input := struct {
				Brand0 string `hummus:"safeway.brands[0]"`
				Brand1 string `hummus:"traderjoes.brands[0]"`
				Brand2 string `hummus:"traderjoes.brands[1]"`
			}{
				Brand0: "sabra",
				Brand1: "athenos",
				Brand2: "cedars",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`
			{
				"safeway": {
					"brands": ["sabra"]
				},
				"traderjoes": {
					"brands": ["athenos", "cedars"]
				}
			}`))
		})

		It("deals with objects inside arrays", func() {
			input := struct {
				Brand0Name string `hummus:"brands[0].name"`
				Brand0Addr string `hummus:"brands[0].address"`
				Brand1Name string `hummus:"brands[1].name"`
				Brand1Addr string `hummus:"brands[1].address"`
			}{
				Brand0Name: "sabra",
				Brand0Addr: "1234 Fake St",
				Brand1Name: "cedars",
				Brand1Addr: "567 Other St",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`{
				"brands": [
					{
						"name": "sabra",
						"address": "1234 Fake St"
					},
					{
						"name": "cedars",
						"address": "567 Other St"
					}
				]
			}`))
		})

		It("handles arrays inside arrays", func() {
			input := struct {
				Brand0Name0 string `hummus:"brands[0].name[0]"`
				Brand0Name1 string `hummus:"brands[0].name[1]"`
				Brand0Addr0 string `hummus:"brands[0].address[0]"`
				Brand0Addr1 string `hummus:"brands[0].address[1]"`
				Brand1Name0 string `hummus:"brands[1].name[0]"`
				Brand1Name1 string `hummus:"brands[1].name[1]"`
			}{
				Brand0Name0: "sabra",
				Brand0Name1: "eatwell",
				Brand0Addr0: "1234 Fake St",
				Brand0Addr1: "1234 Fake2 St",
				Brand1Name0: "cedars",
				Brand1Name1: "pitapal",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`{
				"brands": [
					{
						"name": ["sabra", "eatwell"],
						"address": ["1234 Fake St", "1234 Fake2 St"]
					},
					{
						"name": ["cedars", "pitapal"]
					}
				]
			}`))
		})

		It("handles all kinds of cray-cray", func() {
			input := struct {
				Company           string `hummus:"company"`
				Address           string `hummus:"address"`
				Brand0Name        string `hummus:"brands[0].name"`
				Brand0Flavor      string `hummus:"brands[0].flavor"`
				Brand0Store0Name  string `hummus:"brands[0].stores[0].name"`
				Brand0Store0Price int    `hummus:"brands[0].stores[0].price,omitempty"`
				Brand0Store1Name  string `hummus:"brands[0].stores[1].name"`
				Brand0Store1Price int    `hummus:"brands[0].stores[1].price,omitempty"`
				Brand1Name        string `hummus:"brands[1].name,omitempty"`
				Brand1Flavor      string `hummus:"brands[1].flavor,omitempty"`
				Brand1Store0Name  string `hummus:"brands[1].stores[0].name"`
				Brand1Store0Price int    `hummus:"brands[1].stores[0].price,omitempty"`
				Reputation        string `hummus:"reputation.type"`
			}{
				Company:           "hello foods",
				Address:           "338 New St",
				Brand0Name:        "sabra",
				Brand0Flavor:      "jalapeno",
				Brand0Store0Name:  "safeway",
				Brand0Store0Price: 5,
				Brand0Store1Name:  "wholefoods",
				Brand0Store1Price: 10,
				Brand1Name:        "cedars",
				Brand1Store0Name:  "safeway",
				Brand1Store0Price: 10,
				Reputation:        "good",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`{
				"company": "hello foods",
				"address": "338 New St",
				"brands": [
					{
						"name": "sabra",
						"flavor": "jalapeno",
						"stores": [
							{
								"name": "safeway",
								"price": 5
							},
							{
								"name": "wholefoods",
								"price": 10
							}
						]
					},
					{
						"name": "cedars",
						"stores": [
						  {
							  "name": "safeway",
								"price": 10
							}
						]
					}
				],
				"reputation": {
					"type": "good"
				}
			}`))
		})

		It("allows one to escape dots", func() {
			input := struct {
				A string `hummus:"outer.inner#notchild#notchild2.value.name"`
				B string `hummus:"outer.inner"`
				C string `hummus:"outer.inner#notchild[0].name"`
				D string `hummus:"outer.inner#notchild#notchild2.value.name2"`
				E string `hummus:"outer.inner#notchild[1].name"`
			}{
				A: "A_val",
				B: "B_val",
				C: "C_val",
				D: "D_val",
				E: "E_val",
			}

			outJSON, err := hummus.Marshal(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(outJSON).To(MatchJSON(`{
				"outer": {
					"inner.notchild": [{
							"name": "C_val"
					},
					{
							"name": "E_val"
					}],
					"inner.notchild.notchild2": {
						"value": {
							"name": "A_val",
							"name2": "D_val"
						}
					},
					"inner": "B_val"
				}
			}`))
		})

		Context("special/failure cases", func() {
			Context("when passed an invalid struct tag", func() {
				It("skips the field", func() {
					input := struct {
						Brand0 string `foo:"safeway.brands[0]"`
					}{
						Brand0: "sabra",
					}

					outJSON, err := hummus.Marshal(input)
					Expect(err).NotTo(HaveOccurred())
					Expect(outJSON).To(MatchJSON(`{}`))
				})
			})

			Context("when passed extra struct tag fields", func() {
				It("returns an error", func() {
					input := struct {
						Brand0 string `hummus:"safeway.brands[0],omitempty,blah"`
					}{
						Brand0: "sabra",
					}

					_, err := hummus.Marshal(input)
					Expect(err).To(MatchError("error: invalid number of struct tag fields"))
				})
			})
		})
	})
})
