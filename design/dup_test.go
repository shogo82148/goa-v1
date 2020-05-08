package design_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/goa-v1/design"
)

var _ = Describe("Dup", func() {
	var dt design.DataType
	var dup design.DataType

	JustBeforeEach(func() {
		dup = design.Dup(dt)
	})

	Context("with a primitive type", func() {
		BeforeEach(func() {
			dt = design.Integer
		})

		It("returns the same value", func() {
			Ω(dup).Should(Equal(dt))
		})
	})

	Context("with an array type", func() {
		var elemType = design.Integer

		BeforeEach(func() {
			dt = &design.Array{
				ElemType: &design.AttributeDefinition{Type: elemType},
			}
		})

		It("returns a duplicate array type", func() {
			Ω(dup).Should(Equal(dt))
			Ω(dup == dt).Should(BeFalse())
			Ω(dup.(*design.Array).ElemType == dt.(*design.Array).ElemType).Should(BeFalse())
		})
	})

	Context("with a hash type", func() {
		var keyType = design.String
		var elemType = design.Integer

		BeforeEach(func() {
			dt = &design.Hash{
				KeyType:  &design.AttributeDefinition{Type: keyType},
				ElemType: &design.AttributeDefinition{Type: elemType},
			}
		})

		It("returns a duplicate hash type", func() {
			Ω(dup).Should(Equal(dt))
			Ω(dup == dt).Should(BeFalse())
			Ω(dup.(*design.Hash).KeyType == dt.(*design.Hash).KeyType).Should(BeFalse())
			Ω(dup.(*design.Hash).ElemType == dt.(*design.Hash).ElemType).Should(BeFalse())
		})
	})

	Context("with a user type", func() {
		const typeName = "foo"
		var att = &design.AttributeDefinition{Type: design.Integer}

		BeforeEach(func() {
			dt = &design.UserTypeDefinition{
				TypeName:            typeName,
				AttributeDefinition: att,
			}
		})

		It("returns a duplicate user type", func() {
			Ω(dup).Should(Equal(dt))
			Ω(dup == dt).Should(BeFalse())
			Ω(dup.(*design.UserTypeDefinition).AttributeDefinition == att).Should(BeFalse())
		})
	})

	Context("with a media type", func() {
		var obj = design.Object{"att": &design.AttributeDefinition{Type: design.Integer}}
		var ut = &design.UserTypeDefinition{
			TypeName:            "foo",
			AttributeDefinition: &design.AttributeDefinition{Type: obj},
		}
		const identifier = "vnd.application/test"
		var links = map[string]*design.LinkDefinition{
			"link": {Name: "att", View: "default"},
		}
		var views = map[string]*design.ViewDefinition{
			"default": {
				Name:                "default",
				AttributeDefinition: &design.AttributeDefinition{Type: obj},
			},
		}

		BeforeEach(func() {
			dt = &design.MediaTypeDefinition{
				UserTypeDefinition: ut,
				Identifier:         identifier,
				Links:              links,
				Views:              views,
			}
		})

		It("returns a duplicate media type", func() {
			Ω(dup).Should(Equal(dt))
			Ω(dup == dt).Should(BeFalse())
			Ω(dup.(*design.MediaTypeDefinition).UserTypeDefinition == ut).Should(BeFalse())
		})
	})

	Context("with two media types referring to each other", func() {
		var ut *design.UserTypeDefinition

		BeforeEach(func() {
			mt := &design.MediaTypeDefinition{Identifier: "application/mt1"}
			mt2 := &design.MediaTypeDefinition{Identifier: "application/mt2"}
			obj1 := design.Object{"att": &design.AttributeDefinition{Type: mt2}}
			obj2 := design.Object{"att": &design.AttributeDefinition{Type: mt}}

			att1 := &design.AttributeDefinition{Type: obj1}
			ut = &design.UserTypeDefinition{AttributeDefinition: att1}
			link1 := &design.LinkDefinition{Name: "att", View: "default"}
			view1 := &design.ViewDefinition{AttributeDefinition: att1, Name: "default"}
			mt.UserTypeDefinition = ut
			mt.Links = map[string]*design.LinkDefinition{"att": link1}
			mt.Views = map[string]*design.ViewDefinition{"default": view1}

			att2 := &design.AttributeDefinition{Type: obj2}
			ut2 := &design.UserTypeDefinition{AttributeDefinition: att2}
			link2 := &design.LinkDefinition{Name: "att", View: "default"}
			view2 := &design.ViewDefinition{AttributeDefinition: att2, Name: "default"}
			mt2.UserTypeDefinition = ut2
			mt2.Links = map[string]*design.LinkDefinition{"att": link2}
			mt2.Views = map[string]*design.ViewDefinition{"default": view2}

			dt = mt
		})

		It("duplicates without looping infinity", func() {
			Ω(dup).Should(Equal(dt))
			Ω(dup == dt).Should(BeFalse())
			Ω(dup.(*design.MediaTypeDefinition).UserTypeDefinition == ut).Should(BeFalse())
		})
	})
})

var _ = Describe("DupAtt", func() {
	var att *design.AttributeDefinition
	var dup *design.AttributeDefinition

	JustBeforeEach(func() {
		dup = design.DupAtt(att)
	})

	Context("with an attribute with a type which is a media type", func() {
		BeforeEach(func() {
			att = &design.AttributeDefinition{Type: &design.MediaTypeDefinition{}}
		})

		It("does not clone the type", func() {
			Ω(dup == att).Should(BeFalse())
			Ω(dup.Type == att.Type).Should(BeTrue())
		})
	})
})
