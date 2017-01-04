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
				t.Insert(`hummus:"brand"`, "sabra")
				Expect(t.NodeMap["brand"]).To(Equal(tree.Node{
					Path:        "brand",
					OmitEmpty:   false,
					IsArray:     false,
					SingleChild: "sabra",
				}))
			})
		})

		Context("when provided childless array paths", func() {
			It("inserts/updates the right node in the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[0]"`, "sabra")
				t.Insert(`hummus:"brands[1]"`, "cedars")
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:          "brands",
					OmitEmpty:     false,
					IsArray:       true,
					ArrayChildren: []interface{}{"sabra", "cedars"},
				}))
			})
		})

		Context("when provided array paths out of order", func() {
			It("inserts/updates the right node in the tree in the right order", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[1]"`, "sabra")
				t.Insert(`hummus:"brands[0]"`, "cedars")
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:          "brands",
					OmitEmpty:     false,
					IsArray:       true,
					ArrayChildren: []interface{}{"cedars", "sabra"},
				}))
			})
		})

		Context("when provided array paths with children", func() {
			It("inserts/updates the right node in the tree", func() {
				t := tree.NewTree()
				t.Insert(`hummus:"brands[0].name"`, "sabra")
				t.Insert(`hummus:"brands[0].location"`, "california")
				t.Insert(`hummus:"brands[1].name"`, "cedars")
				t.Insert(`hummus:"brands[1].location"`, "washington")
				Expect(t.NodeMap["brands"]).To(Equal(tree.Node{
					Path:      "brands",
					OmitEmpty: false,
					IsArray:   true,
					ArrayChildren: []interface{}{
						tree.Tree{
							NodeMap: map[string]tree.Node{
								"name": tree.Node{
									Path:        "name",
									OmitEmpty:   false,
									IsArray:     false,
									SingleChild: "sabra",
								},
								"location": tree.Node{
									Path:        "location",
									OmitEmpty:   false,
									IsArray:     false,
									SingleChild: "california",
								},
							},
						},
						tree.Tree{
							NodeMap: map[string]tree.Node{
								"name": tree.Node{
									Path:        "name",
									OmitEmpty:   false,
									IsArray:     false,
									SingleChild: "cedars",
								},
								"location": tree.Node{
									Path:        "location",
									OmitEmpty:   false,
									IsArray:     false,
									SingleChild: "washington",
								},
							},
						},
					},
				}))
			})
		})
	})
})
