package design

import (
	. "github.com/shogo82148/goa-v1/design"
	. "github.com/shogo82148/goa-v1/design/apidsl"
)

var _ = Resource("health", func() {

	BasePath("/_ah")

	Action("health", func() {
		Routing(
			GET("/health"),
		)
		Description("Perform health check.")
		Response(OK, "text/plain")
	})
})
