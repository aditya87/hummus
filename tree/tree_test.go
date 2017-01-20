package tree_test

import (
	"github.com/aditya87/hummus/tree"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {
	Describe("Insert", func() {
		Context("when provided a simple path and child", func() {
			It("inserts a node into the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brand"`, "sabra", false)
				Expect(t.NodeMap["brand"]).To(Equal(tree.Node{
					Path:        "brand",
					IsArray:     false,
					SingleChild: "sabra",
				}))
			})
		})

		Context("when provided simple array paths", func() {
			It("inserts/updates the right node in the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[0]"`, "sabra", false)
				t.Insert(`hummus:"brands[1]"`, "cedars", false)
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:          "brands",
					IsArray:       true,
					ArrayChildren: []interface{}{"sabra", "cedars"},
				}))
			})
		})

		Context("when provided array paths out of order", func() {
			It("inserts/updates the right node in the tree in the right order", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[1]"`, "sabra", false)
				t.Insert(`hummus:"brands[0]"`, "cedars", false)
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:          "brands",
					IsArray:       true,
					ArrayChildren: []interface{}{"cedars", "sabra"},
				}))
			})
		})

		Context("when provided array paths with tree children", func() {
			It("inserts/updates the right node in the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[1].name"`, "cedars", false)
				t.Insert(`hummus:"brands[0].name"`, "sabra", false)
				t.Insert(`hummus:"brands[0].location"`, "california", false)
				t.Insert(`hummus:"brands[1].location"`, "washington", false)
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:    "brands",
					IsArray: true,
					ArrayChildren: []interface{}{
						tree.Tree{
							NodeMap: map[string]tree.Node{
								"name": tree.Node{
									Path:        "name",
									IsArray:     false,
									SingleChild: "sabra",
								},
								"location": tree.Node{
									Path:        "location",
									IsArray:     false,
									SingleChild: "california",
								},
							},
						},
						tree.Tree{
							NodeMap: map[string]tree.Node{
								"name": tree.Node{
									Path:        "name",
									IsArray:     false,
									SingleChild: "cedars",
								},
								"location": tree.Node{
									Path:        "location",
									IsArray:     false,
									SingleChild: "washington",
								},
							},
						},
					},
				}))
			})
		})

		Context("when provided nested array paths", func() {
			It("inserts/updates the right node in the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[0].name"`, "sabra", false)
				t.Insert(`hummus:"brands[0].company[0].name"`, "cool foods", false)
				t.Insert(`hummus:"brands[0].company[1].name"`, "hipster foods", false)
				t.Insert(`hummus:"brands[1].name"`, "cedars", false)
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:    "brands",
					IsArray: true,
					ArrayChildren: []interface{}{
						tree.Tree{
							NodeMap: map[string]tree.Node{
								"name": tree.Node{
									Path:        "name",
									IsArray:     false,
									SingleChild: "sabra",
								},
								"company": tree.Node{
									Path:    "company",
									IsArray: true,
									ArrayChildren: []interface{}{
										tree.Tree{
											NodeMap: map[string]tree.Node{
												"name": tree.Node{
													Path:        "name",
													IsArray:     false,
													SingleChild: "cool foods",
												},
											},
										},
										tree.Tree{
											NodeMap: map[string]tree.Node{
												"name": tree.Node{
													Path:        "name",
													IsArray:     false,
													SingleChild: "hipster foods",
												},
											},
										},
									},
								},
							},
						},
						tree.Tree{
							NodeMap: map[string]tree.Node{
								"name": tree.Node{
									Path:        "name",
									IsArray:     false,
									SingleChild: "cedars",
								},
							},
						},
					},
				}))
			})
		})
	})

	Describe("BuildJSON", func() {
		Context("when given a simple tree", func() {
			It("builds a json from the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brand"`, "sabra", false)
				jsonObj := t.BuildJSON()
				Expect([]byte(jsonObj.String())).To(MatchJSON(`{
					"brand": "sabra"
				}`))
			})
		})

		Context("when given array paths with children", func() {
			It("builds a json from the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[1].name"`, "cedars", false)
				t.Insert(`hummus:"brands[0].name"`, "sabra", false)
				t.Insert(`hummus:"brands[0].location"`, "california", false)
				t.Insert(`hummus:"brands[1].location"`, "washington", false)
				jsonObj := t.BuildJSON()
				Expect([]byte(jsonObj.String())).To(MatchJSON(`{
					"brands": [
					  {
							"name": "sabra",
							"location": "california"
						},
						{
							"name": "cedars",
							"location": "washington"
						}
					]
				}`))
			})
		})

		Context("when given array paths with children", func() {
			It("builds a json from the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[0].name"`, "sabra", false)
				t.Insert(`hummus:"brands[0].company[0].name"`, "cool foods", false)
				t.Insert(`hummus:"brands[0].company[1].name"`, "hipster foods", false)
				t.Insert(`hummus:"brands[1].name"`, "cedars", false)
				jsonObj := t.BuildJSON()
				Expect([]byte(jsonObj.String())).To(MatchJSON(`{
					"brands": [
						{
							"name": "sabra",
							"company": [
								{
									"name": "cool foods"
								},
								{
									"name": "hipster foods"
								}
							]
						},
						{
							"name": "cedars"
						}
					]
				}`))
			})
		})
	})
})
