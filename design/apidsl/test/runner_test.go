package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/goa-v1/design"
	"github.com/shogo82148/goa-v1/design/apidsl"
	"github.com/shogo82148/goa-v1/dslengine"
)

// Global test definitions
const apiName = "API"
const apiDescription = "API description"
const resourceName = "R"
const resourceDescription = "R description"
const typeName = "T"
const typeDescription = "T description"
const mediaTypeIdentifier = "mt/json"
const mediaTypeDescription = "MT description"

var _ = apidsl.API(apiName, func() {
	apidsl.Description(apiDescription)
})

var _ = apidsl.Resource(resourceName, func() {
	apidsl.Description(resourceDescription)
})

var _ = apidsl.Type(typeName, func() {
	apidsl.Description(typeDescription)
	apidsl.Attribute("bar")
})

var _ = apidsl.MediaType(mediaTypeIdentifier, func() {
	apidsl.Description(mediaTypeDescription)
	apidsl.Attributes(func() { apidsl.Attribute("foo") })
	apidsl.View("default", func() { apidsl.Attribute("foo") })
})

func init() {
	dslengine.Run()

	var _ = Describe("DSL execution", func() {
		Context("with global DSL definitions", func() {
			It("runs the DSL", func() {
				Ω(dslengine.Errors).Should(BeEmpty())

				Ω(design.Design).ShouldNot(BeNil())
				Ω(design.Design.Name).Should(Equal(apiName))
				Ω(design.Design.Description).Should(Equal(apiDescription))

				Ω(design.Design.Resources).Should(HaveKey(resourceName))
				Ω(design.Design.Resources[resourceName]).ShouldNot(BeNil())
				Ω(design.Design.Resources[resourceName].Name).Should(Equal(resourceName))
				Ω(design.Design.Resources[resourceName].Description).Should(Equal(resourceDescription))

				Ω(design.Design.Types).Should(HaveKey(typeName))
				Ω(design.Design.Types[typeName]).ShouldNot(BeNil())
				Ω(design.Design.Types[typeName].TypeName).Should(Equal(typeName))
				Ω(design.Design.Types[typeName].Description).Should(Equal(typeDescription))

				Ω(design.Design.MediaTypes).Should(HaveKey(mediaTypeIdentifier))
				Ω(design.Design.MediaTypes[mediaTypeIdentifier]).ShouldNot(BeNil())
				Ω(design.Design.MediaTypes[mediaTypeIdentifier].Identifier).Should(Equal(mediaTypeIdentifier))
				Ω(design.Design.MediaTypes[mediaTypeIdentifier].Description).Should(Equal(mediaTypeDescription))
			})
		})
	})
}
