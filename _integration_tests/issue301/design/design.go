package design

import (
	. "github.com/shogo82148/goa-v1/design"
	. "github.com/shogo82148/goa-v1/design/apidsl"
)

var _ = API("issue301", func() {
	Title("This API has user definition type default values")
	Host("localhost:8080")
	Scheme("https")
})

var _ = Resource("issue301", func() {
	Action("test", func() {
		Routing(GET("user-definition"))
		Payload(Issue301Payload)
		Response(OK)
	})
})

var Issue301Payload = Type("Issue301Type", func() {
	Attribute("user-definition-type", Integer, func() {
		Default(10)
		Metadata("struct:field:type", "design.SecuritySchemeKind", "github.com/shogo82148/goa-v1/design")
	})

	Attribute("primitive-type-number", Number, func() {
		Default(3.14)
	})

	Attribute("primitive-type-time", DateTime, func() {
		Default("2006-01-02T15:04:05Z")
	})

	Required("user-definition-type")
	Required("primitive-type-number")
	Required("primitive-type-time")
})
